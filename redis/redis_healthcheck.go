package redis_factory

import (
	"goapm/logger"
	"goapm/utils"
	"goapm/web"
	"context"
	"sync"
)

var redisHealthCheckOnce sync.Once
var redisHealthCheckService *RedisHealthCheckImpl

type RedisHealthCheckImpl struct {
}

func NewRedisHealthCheckService() *RedisHealthCheckImpl {
	redisHealthCheckOnce.Do(func() {
		redisHealthCheckService = &RedisHealthCheckImpl{
		}
		web.HealthChecksToRun["redisMs"] = redisHealthCheckService
	})
	return redisHealthCheckService
}


// CheckService perform health check(ping) for first MS connector which exists on connectors map -
// specific_redis_connection - var redisConnMap = ds_utils.NewSyncedMap()
func (check *RedisHealthCheckImpl) CheckService(ctx context.Context) utils.ServiceHealth {
	serviceHealth := utils.ServiceHealth{
		Name:          "MsRedisConnectorHealthCheck",
		Status:        "UP",
		StatusCode:    200,
		ErrorMessages: make(map[string]interface{}),
	}
	if redisConnMap.Size() > 0 {
		connector := redisConnMap.Values()[0].(RedisFactoryInterface)
		_, err := connector.Ping(ctx)
		if err != nil {
			logger.LOGGER.Error("Unable to perform redis healthcheck, err - ", err, ", connector - ", connector)
			serviceHealth.Status = "DOWN"
			serviceHealth.StatusCode = 503
			serviceHealth.ErrorMessages["error"] = err.Error()
		}
	}

	return serviceHealth
}
