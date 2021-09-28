package services

import (
	"goapm/clickhouse"
	"goapm/logger"
	redis_factory "goapm/redis"
	"context"
	"database/sql"
	"go.uber.org/zap"
	"sync"
	"time"
)

var timeLayout = "2006-01-02 15:04:05"

var traceFilterJobOnce sync.Once
var traceFilterJob *TraceFilterJob

type TraceFilterJob struct {
	Logger            *zap.SugaredLogger
	ClickhouseService clickhouse.ClickhouseConnectionService
	RedisService      redis_factory.SpecificRedisService
}

func NewTraceFilterJob(ClickhouseService clickhouse.ClickhouseConnectionService,
	RedisService redis_factory.SpecificRedisService) *TraceFilterJob {
	traceFilterJobOnce.Do(func() {
		traceFilterJob = &TraceFilterJob{
			Logger:            logger.LOGGER,
			ClickhouseService: ClickhouseService,
			RedisService:      RedisService,
		}
	})
	return traceFilterJob
}

func (service *TraceFilterJob) Run() {
	ctx := context.Background()

	var args []interface{}
	args = append(args, "TRACE_FILTER")
	var lastSuccess string
	err := service.ClickhouseService.ExecuteSelectFunction(false, &lastSuccess, "SELECT last_success_date FROM last_success FINAL WHERE last_success_key = ?", args)

	if err != nil {
		service.Logger.Error("unable to get last success", err)
		return
	}
	timeParsed, err := time.Parse(timeLayout, lastSuccess)
	if err != nil {
		service.Logger.Error("unable to parse last success", err)
		return
	}

	if time.Now().UTC().Add(-5 * time.Minute).Before(timeParsed) {
		return
	}
	currTime := time.Now()
	redisFactory := service.RedisService.GetSpecificRedis()
	inserted, err := redisFactory.SetNX(ctx, "TRACE_FILTER", currTime.Format(timeLayout), 120)
	if err != nil || !inserted {
		service.Logger.Error("unable to set value into redis", err)
		return
	}
	fromDate := timeParsed.Add(-20 * time.Minute).Format(timeLayout)

	query := "INSERT INTO signoz_index_final SELECT * FROM signoz_index_tmp " +
		"WHERE timestamp >= ? AND traceID IN (SELECT DISTINCT(traceID) FROM signoz_index_tmp WHERE (kind = 2 AND durationNano >= 2000000000) OR statusCode >= 400) " +
		"AND spanID NOT IN (SELECT DISTINCT(spanID) FROM signoz_index_final)"

	err = service.ClickhouseService.ExecuteInsertFunction(query, func(stmt *sql.Stmt) error {
		_, err = stmt.Exec(fromDate)
		return err
	})

	if err != nil {
		service.Logger.Error("unable to get problematic traces", err)
		return
	}

	query = "INSERT INTO last_success (last_success_key, last_success_date, event_time) VALUES (?, ?, ?)"
	err = service.ClickhouseService.ExecuteInsertFunction(query, func(stmt *sql.Stmt) error {
		_, err = stmt.Exec("TRACE_FILTER", currTime.Format(timeLayout), currTime.Format(timeLayout))
		return err
	})
	if err != nil {
		service.Logger.Error("unable to insert to last success", err)
	}

}
