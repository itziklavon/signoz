package clickhouse

import (
	"goapm/ds_utils"
	logger2 "goapm/logger"
	"goapm/utils"
	"database/sql"
	"errors"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

var clickhousConnectionOnce sync.Once
var clickhouseConnectionService *ClickhouseConnectionServiceImpl
var logger = logger2.LOGGER
var connectionsMap = ds_utils.NewSyncedMap()
var dataSources = ds_utils.NewSyncedHashSet()

const emptyResult = "sql: no rows in result set"

type stmt func(stmt *sql.Stmt) error

type ClickhouseConnectionServiceImpl struct {
}

func NewClickhouseConnectionService() *ClickhouseConnectionServiceImpl {
	clickhousConnectionOnce.Do(func() {
		clickhouseConnectionService = &ClickhouseConnectionServiceImpl{
		}
		NewClickhouseHealthCheckService(clickhouseConnectionService)
	})
	return clickhouseConnectionService
}

func Connect() {
	url := viper.GetString("GO_CLICKHOUSE_CONNECTION_URL")
	userName := viper.GetString("CLICKHOUSE_USER_NAME")
	password := viper.GetString("CLICKHOUSE_PASSWORD")
	dateSources := viper.GetString("CLICKHOUSE_DATA_SOURCES")

	dateSourcesArr := strings.Split(dateSources, ",")
	var strSlice []string

	for i, s := range dateSourcesArr {
		strSlice = append(strSlice, s)
		logger.Info(fmt.Sprintf("data source at index = %d is = %s", i, s))
	}
	urls := GetDataSourcesUrls(strSlice, url, userName, password)

	for _, url := range urls {
		connect, err := sqlx.Open("clickhouse", url)
		if err != nil {
			logger.Error("can't connect clickhouse, error = ", err, ", stack trace = ", string(debug.Stack()))
		}
		if err = connect.Ping(); err != nil {
			if exception, ok := err.(*clickhouse.Exception); ok {
				logger.Debug(fmt.Sprintf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace))
			} else {
				logger.Error("can't connect clickhouse, error = ", err, ", stack trace = ", string(debug.Stack()))
			}
			return
		}
		connect.SetMaxOpenConns(utils.GetOrDefaultInt(viper.GetString("CLICKHOUSE_MAX_CONNECTIONS"), 10))
		connect.SetConnMaxLifetime(1 * time.Hour)
		connectionsMap.Put(url, connect)
		dataSources.Add(url)
	}
	logger.Debug("clickhouse connection map = ", connectionsMap)
}

func (services *ClickhouseConnectionServiceImpl) GetConnectionMap() *ds_utils.ConcurrentHashMap {
	if connectionsMap.Size() == 0 {
		Connect()
	}
	return connectionsMap
}

func (services *ClickhouseConnectionServiceImpl) ExecuteInsertFunction(query string, statement stmt) error {
	if dataSources.Size() == 0 {
		Connect()
	}
	for _, s := range dataSources.Values() {
		connectVal, _ := connectionsMap.Get(s)
		connect := connectVal.(*sqlx.DB)
		var (
			tx, _             = connect.Begin()
			preparedStmt, err = tx.Prepare(query)
		)
		if err != nil {
			logger.Error("an error occurred while insert query = "+query+", error = ", err, ", stack trace = ", string(debug.Stack()))
		}
		defer preparedStmt.Close()
		err = statement(preparedStmt)
		if err != nil {
			logger.Error("an error occurred while insert query = "+query+", error = ", err, ", stack trace = ", string(debug.Stack()))
			_ = tx.Rollback()
		} else if err = tx.Commit(); err != nil {
			logger.Error("an error occurred while insert query = "+query+", error = ", err, ", stack trace = ", string(debug.Stack()))
		} else {
			return nil
		}
	}
	logger.Error("couldn't execute method for all data sources, query = "+query, ", stack trace = ", string(debug.Stack()))
	return errors.New("couldn't insert data, query = " + query)
}

func (services *ClickhouseConnectionServiceImpl) ExecuteSelectFunction(isArray bool, dest interface{}, query string, args []interface{}) error {
	if dataSources.Size() == 0 {
		Connect()
	}
	var err error
	dataSourcesReversed := Reverse(dataSources)
	for _, s := range dataSourcesReversed {
		connectVal, _ := connectionsMap.Get(s)
		connect := connectVal.(*sqlx.DB)
		if args != nil {
			if isArray {
				err = connect.Select(dest, query, args...)
			} else {
				err = connect.Get(dest, query, args...)
			}
		} else {
			if isArray {
				err = connect.Select(dest, query)
			} else {
				err = connect.Get(dest, query)
			}
		}
		if err != nil && !strings.EqualFold(err.Error(), emptyResult) {
			logger.Error("an error occurred while select query = "+query+", error = ", err, ", stack trace = ", string(debug.Stack()))
		} else {
			return nil
		}
	}
	logger.Error("couldn't execute method for all data sources", ", stack trace = ", string(debug.Stack()))
	return errors.New("couldn't execute select method, " + err.Error())
}

func Reverse(inputSet *ds_utils.ConcurrentHashSet) []interface{} {
	var reverseList []interface{}
	if inputSet.Size() == 0 {
		return reverseList
	}

	input := inputSet.Values()
	for i := len(input) - 1; i >= 0; i-- {
		reverseList = append(reverseList, input[i])
	}
	return reverseList
}

func GetDataSourcesUrls(strSlice []string, url string, userName string, password string) []string {
	var urls []string
	for i, s := range strSlice {
		logger.Debug("date source is = "+s, i)
		dateSourcesSplit := strings.Split(s, "/")
		clickhouseHost := strings.Split(dateSourcesSplit[2], ":")[0]
		database := strings.Split(dateSourcesSplit[3], "?")[0]

		url1 := strings.Replace(url, "{clickhouse_host}", clickhouseHost, 1)
		url1 = strings.Replace(url1, "{username}", userName, 1)
		url1 = strings.Replace(url1, "{password}", password, 1)
		url1 = strings.Replace(url1, "{database}", database, 1)
		urls = append(urls, url1)
	}
	logger.Debug("url list = ", urls)
	return urls
}
