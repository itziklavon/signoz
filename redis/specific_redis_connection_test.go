package redis_factory

import (
	"testing"
)

func TestGetSpecificRedis(t *testing.T) {
	classUnderTestSpecificRedis.GetSpecificRedis()
	classUnderTestSpecificRedis.GetSpecificRedis()
}

func TestGetSpecificRedisHost(t *testing.T) {
	classUnderTestSpecificRedis.GetSpecificRedis(host)
	classUnderTestSpecificRedis.GetSpecificRedis(host)
}

func TestGetSpecificRedisWithPort(t *testing.T) {
	classUnderTestSpecificRedis.GetSpecificRedisWithPort(host)
	classUnderTestSpecificRedis.GetSpecificRedisWithPort(host)
}

func TestGetSpecificRedisWithPortAndDatabase(t *testing.T) {
	classUnderTestSpecificRedis.GetSpecificRedisWithParams(host, port)
	classUnderTestSpecificRedis.GetSpecificRedisWithParams(host, port)
}
