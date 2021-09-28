package clickhouse

import (
	redis_factory "goapm/redis"
	"context"
	"github.com/stretchr/testify/mock"
)

type ClickhouseHealthCheck interface {
	CheckHealthyClickhouse(key string, redisFactory redis_factory.RedisFactoryInterface, ctx context.Context) error
}

type MockClickhouseHealthCheck struct {
	mock.Mock
}

func (service *MockClickhouseHealthCheck) CheckHealthyClickhouse(key string, redisFactory redis_factory.RedisFactoryInterface, ctx context.Context) error {
	args := service.Called(mock.Anything, mock.Anything, mock.Anything)
	return args.Error(0)
}
