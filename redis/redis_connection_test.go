package redis_factory

import (
	"goapm/ds_utils"
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

var host string
var port string

var classUnderTest RedisFactoryInterface
var classUnderTestSpecificRedis SpecificRedisService

func TestMain(m *testing.M) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	redisConnMap = ds_utils.NewSyncedMap()

	host = s.Host()
	port = s.Port()
	viper.Set("REDIS_PORT", port)
	viper.Set("REDIS_HOST", host)
	redisFactoryService := RedisFactoryServiceImpl{}
	classUnderTest = redisFactoryService.GetConnection(host, port, 0)
	classUnderTestSpecificRedis = NewSpecificRedisService()
	m.Run()
}

func TestGetHGet(t *testing.T) {
	_ = NewRedisFactory()
	err := classUnderTest.Set(context.Background(), "tst", "tst")
	assert.Nil(t, err)

	response, err := classUnderTest.Ping(context.Background())

	response, err = classUnderTest.Get(context.Background(), "tst")
	assert.Equal(t, "tst", response)
	assert.Nil(t, err)

	err = classUnderTest.Del(context.Background(), "tst")
	assert.Nil(t, err)

	err = classUnderTest.SetEx(context.Background(), "tst", "tst", 10)
	assert.Nil(t, err)

	_, err = classUnderTest.Ttl(context.Background(), "tst")
	assert.Nil(t, err)

	response, err = classUnderTest.Get(context.Background(), "tst")
	assert.Equal(t, "tst", response)
	assert.Nil(t, err)

	err = classUnderTest.Del(context.Background(), "tst")
	assert.Nil(t, err)

	err = classUnderTest.HSet(context.Background(), "tst", "tst", "tst")
	assert.Nil(t, err)

	response, err = classUnderTest.HGet(context.Background(), "tst", "tst")
	assert.Equal(t, "tst", response)
	assert.Nil(t, err)

	myMap := make(map[string]string)
	myMap["tst2"] = "tst2"

	err = classUnderTest.HSetAll(context.Background(), "tst2", myMap)
	assert.Nil(t, err)

	keys, e := classUnderTest.Keys(context.Background(), "tst*")
	assert.Equal(t, 2, len(keys))
	assert.Nil(t, e)

	err = classUnderTest.HDel(context.Background(), "tst2", "tst2")

	responseMap, err := classUnderTest.HGetAll(context.Background(), "tst")
	assert.Equal(t, 1, len(responseMap))
	assert.Nil(t, err)

	responseBool, err := classUnderTest.SetNX(context.Background(), "1", "1", 10)
	assert.True(t, responseBool)
	responseBool, err = classUnderTest.SetNX(context.Background(), "1", "1", 10)
	assert.False(t, responseBool)

	err = classUnderTest.Expire(context.Background(), "1", 10)

	responseExists, err := classUnderTest.Exists(context.Background(), "1")
	assert.True(t, responseExists > 0)
}
