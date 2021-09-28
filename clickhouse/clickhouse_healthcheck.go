package clickhouse

import (
	"goapm/utils"
	"goapm/web"
	"context"
	"github.com/ClickHouse/clickhouse-go"
	"github.com/jmoiron/sqlx"
	"strconv"
	"sync"
	"time"
)

var lastRunTime = time.Now().Add(-1 * time.Hour).UTC()

var clickhouseHealthCheckOnce sync.Once
var clickhouseHealthCheckService *ClickhouseHealthCheckImpl

type ClickhouseHealthCheckImpl struct {
	ClickhouseConnectionService ClickhouseConnectionService
}

func NewClickhouseHealthCheckService(service ClickhouseConnectionService) *ClickhouseHealthCheckImpl {
	clickhouseHealthCheckOnce.Do(func() {
		clickhouseHealthCheckService = &ClickhouseHealthCheckImpl{
			ClickhouseConnectionService: service,
		}
		web.HealthChecksToRun["clickhouseMs"] = clickhouseHealthCheckService
	})
	return clickhouseHealthCheckService
}

func (check *ClickhouseHealthCheckImpl) CheckService(ctx context.Context) utils.ServiceHealth {
	serviceHealth := utils.ServiceHealth{
		Name:          "MsClickhouseServer",
		Status:        "UP",
		StatusCode:    200,
		ErrorMessages: make(map[string]interface{}),
	}

	if time.Since(lastRunTime).Minutes() > 10 {
		connections := check.ClickhouseConnectionService.GetConnectionMap()

		if connections.Size() > 0 {
			connectionVal := connections.Values()[0]
			connect := connectionVal.(*sqlx.DB)
			if err := connect.Ping(); err != nil {
				if exception, ok := err.(*clickhouse.Exception); ok {
					logger.Errorf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
					serviceHealth.Status = "DOWN"
					serviceHealth.StatusCode = 503
					serviceHealth.ErrorMessages["error_code"] = strconv.Itoa(int(exception.Code))
					serviceHealth.ErrorMessages["error_message"] = exception.Message
				} else {
					serviceHealth.Status = "DOWN"
					serviceHealth.StatusCode = 503
					serviceHealth.ErrorMessages["error"] = err.Error()
				}
			}
		}
		lastRunTime = time.Now().UTC()
	}
	return serviceHealth
}
