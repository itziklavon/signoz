package redis_factory

import (
	"goapm/ds_utils"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)


func TestHealthyRedisSuccess(t *testing.T) {
	redisConnMap = ds_utils.NewSyncedMap()
	redisMock := new(MockRedisFactory)
	redisConnMap.Put("ms", redisMock)
	redisMock.On("Ping", mock.Anything).Return("Pong", nil)
	classUnderTest := NewRedisHealthCheckService()
	serviceHealthResponse := classUnderTest.CheckService(context.Background())
	assert.True(t, serviceHealthResponse.StatusCode == 200)
}

func TestHealthyRedisSuccessErr(t *testing.T) {
	redisConnMap = ds_utils.NewSyncedMap()
	redisMock := new(MockRedisFactory)
	redisConnMap.Put("ms", redisMock)
	redisMock.On("Ping", mock.Anything).Return("Pong", errors.New("error"))
	classUnderTest := NewRedisHealthCheckService()
	serviceHealthResponse := classUnderTest.CheckService(context.Background())
	assert.False(t, serviceHealthResponse.StatusCode == 200)
}
