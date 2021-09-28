package controllers

import (
	"goapm/clickhouse"
	redis_factory "goapm/redis"
	"goapm/dao"
	"goapm/services"
)

var apmService *services.ApmServiceImpl
var traceFilterJob *services.TraceFilterJob

func InitHealthCheck() {
	clickhouse.NewClickhouseConnectionService()
	redisService := redis_factory.NewSpecificRedisService()
	_ = redisService.GetSpecificRedis()
}

func InitServices() {
	apmDao := dao.NewApmDao(clickhouse.NewClickhouseConnectionService())
	apmService = services.NewApmServiceImpl(apmDao)

	traceFilterJob = services.NewTraceFilterJob(clickhouse.NewClickhouseConnectionService(),
		redis_factory.NewSpecificRedisService())
}
