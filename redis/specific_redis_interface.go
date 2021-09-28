package redis_factory

import "github.com/stretchr/testify/mock"

type SpecificRedisService interface {
	GetSpecificRedis(host ...string) RedisFactoryInterface
	GetSpecificRedisWithPort(host string, port ...string) RedisFactoryInterface
	GetSpecificRedisWithParams(host string, port string, database ...int) RedisFactoryInterface
}

type MockSpecificRedisFactory struct {
	mock.Mock
}

func (factory MockSpecificRedisFactory) GetSpecificRedis(host ...string) RedisFactoryInterface {
	args := factory.Called(mock.Anything)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(RedisFactoryInterface)
}
func (factory MockSpecificRedisFactory) GetSpecificRedisWithPort(host string, port ...string) RedisFactoryInterface {
	args := factory.Called(mock.Anything, mock.Anything)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(RedisFactoryInterface)
}
func (factory MockSpecificRedisFactory) GetSpecificRedisWithParams(host string, port string, database ...int) RedisFactoryInterface {
	args := factory.Called(mock.Anything, mock.Anything, mock.Anything)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(RedisFactoryInterface)
}
