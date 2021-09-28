package redis_factory

import (
	"goapm/ds_utils"
	"github.com/spf13/viper"
	"strconv"
	"sync"
)

var redisConnMap = ds_utils.NewSyncedMap()
var once sync.Once
var specificRedisService *SpecificRedisServiceImpl

type SpecificRedisServiceImpl struct {
}

// NewSpecificRedisService init redis service singleton
func NewSpecificRedisService() *SpecificRedisServiceImpl {
	once.Do(func() {
		specificRedisService = &SpecificRedisServiceImpl{}
		NewRedisHealthCheckService()
	})
	return specificRedisService
}

// GetSpecificRedis if only host is being sent, or nothing, use default values for redis connection
func (service SpecificRedisServiceImpl) GetSpecificRedis(host ...string) RedisFactoryInterface {
	dbHost := viper.GetString("REDIS_HOST")
	port := viper.GetString("REDIS_PORT")
	if len(port) == 0 {
		port = "6379"
	}
	if len(host) > 0 {
		dbHost = host[0]
	}

	if val, ok := redisConnMap.Get(dbHost + "_" + port); ok {
		return val.(RedisFactoryInterface)
	}
	factory := service.GetSpecificRedisWithPort(dbHost, port)
	return factory
}

// GetSpecificRedisWithPort if only host is being sent, or with port, init redis connection
func (service SpecificRedisServiceImpl) GetSpecificRedisWithPort(host string, port ...string) RedisFactoryInterface {
	dbPort := "6379"
	if len(port) > 0 {
		dbPort = port[0]
	}

	if val, ok := redisConnMap.Get(host + "_" + dbPort + "_0"); ok {
		return val.(RedisFactoryInterface)
	}
	factory := service.GetSpecificRedisWithParams(host, dbPort, 0)
	return factory
}

// GetSpecificRedisWithParams if only host is being sent, with port, and database, init connection
func (service SpecificRedisServiceImpl) GetSpecificRedisWithParams(host string, port string, database ...int) RedisFactoryInterface {
	connectionDb := 0
	if len(database) > 0 {
		connectionDb = database[0]
	}
	if val, ok := redisConnMap.Get(host + "_" + port + "_" + strconv.Itoa(connectionDb)); ok {
		return val.(RedisFactoryInterface)
	}
	factoryService := RedisFactoryServiceImpl{}
	factory := factoryService.GetConnection(host, port, connectionDb)
	redisConnMap.Put(host+"_"+port+"_"+strconv.Itoa(connectionDb), factory)
	return factory

}
