package redis_factory

import (
	"context"
	"github.com/stretchr/testify/mock"
	"time"
)

type RedisFactoryService interface {
	GetConnection(host string, port string, db int) RedisFactoryInterface
}

type RedisFactoryInterface interface {
	Ping(ctx context.Context) (string, error)
	Set(ctx context.Context, key string, val string) error
	SetEx(ctx context.Context, key string, val string, ttlSeconds int) error
	SetNX(ctx context.Context, key string, val string, ttlSeconds int) (bool, error)
	Exists(ctx context.Context, key string) (int64, error)
	Get(ctx context.Context, key string) (string, error)
	Ttl(ctx context.Context, key string) (time.Duration, error)
	Expire(ctx context.Context, key string, seconds int) error
	Del(ctx context.Context, key string) error
	HGet(ctx context.Context, key string, innerKey string) (string, error)
	HSet(ctx context.Context, key string, innerKey string, val string) error
	HDel(ctx context.Context, key string, innerKey string) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HSetAll(ctx context.Context, key string, values map[string]string) error
	Keys(ctx context.Context, pattern string) ([]string, error)
}

type MockRedisFactory struct {
	mock.Mock
}

type MockRedisFactoryService struct {
	mock.Mock
}

func (redisFactoryMock MockRedisFactoryService) GetConnection(host string, port string, db int) RedisFactoryInterface {
	return new(MockRedisFactory)
}

func (factory MockRedisFactory) Ping(ctx context.Context) (string, error) {
	args := factory.Called(mock.Anything)
	return args.String(0), args.Error(1)
}

func (factory MockRedisFactory) Set(ctx context.Context, key string, val string) error {
	args := factory.Called(mock.Anything, mock.Anything, mock.Anything)
	return args.Error(0)
}

func (factory MockRedisFactory) SetEx(ctx context.Context, key string, val string, ttlSeconds int) error {
	args := factory.Called(mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	return args.Error(0)
}

func (factory MockRedisFactory) SetNX(ctx context.Context, key string, val string, ttlSeconds int) (bool, error) {
	args := factory.Called(mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	return args.Bool(0), args.Error(1)
}

func (factory MockRedisFactory) Exists(ctx context.Context, key string) (int64, error) {
	args := factory.Called(mock.Anything, mock.Anything)
	return int64(args.Int(0)), args.Error(1)
}

func (factory MockRedisFactory) Get(ctx context.Context, key string) (string, error) {
	args := factory.Called(mock.Anything, mock.Anything, mock.Anything)
	return args.String(0), args.Error(1)
}

func (factory MockRedisFactory) Ttl(ctx context.Context, key string) (time.Duration, error) {
	args := factory.Called(mock.Anything, mock.Anything, mock.Anything)
	return args.Get(0).(time.Duration), args.Error(1)
}

func (factory MockRedisFactory) Expire(ctx context.Context, key string, seconds int) error {
	args := factory.Called(mock.Anything, mock.Anything, mock.Anything)
	return args.Error(0)
}

func (factory MockRedisFactory) Del(ctx context.Context, key string) error {
	args := factory.Called(mock.Anything, mock.Anything, mock.Anything)
	return args.Error(0)
}

func (factory MockRedisFactory) HGet(ctx context.Context, key string, innerKey string) (string, error) {
	args := factory.Called(mock.Anything, mock.Anything, mock.Anything)
	return args.String(0), args.Error(1)
}

func (factory MockRedisFactory) HSet(ctx context.Context, key string, innerKey string, val string) error {
	args := factory.Called(mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	return args.Error(0)
}

func (factory MockRedisFactory) HDel(ctx context.Context, key string, innerKey string) error {
	args := factory.Called(mock.Anything, mock.Anything, mock.Anything)
	return args.Error(0)
}

func (factory MockRedisFactory) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	args := factory.Called(mock.Anything, mock.Anything)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]string), args.Error(1)
}

func (factory MockRedisFactory) HSetAll(ctx context.Context, key string, values map[string]string) error {
	args := factory.Called(mock.Anything, mock.Anything, mock.Anything)
	return args.Error(0)
}

func (factory MockRedisFactory) Keys(ctx context.Context, pattern string) ([]string, error) {
	args := factory.Called(mock.Anything, mock.Anything)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}