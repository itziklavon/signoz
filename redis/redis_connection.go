package redis_factory

import (
	"context"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

var redisFactoryOnce sync.Once
var redisFactoryService *RedisFactoryServiceImpl

type RedisFactory struct {
	*redis.Client
}

type RedisFactoryServiceImpl struct {
}

// NewRedisFactory init new redis factory singleton instance,
//factory will contain all available redis methods
func NewRedisFactory() *RedisFactoryServiceImpl {
	redisFactoryOnce.Do(func() {
		redisFactoryService = &RedisFactoryServiceImpl{}
	})
	return redisFactoryService
}

func (RedisFactoryService *RedisFactoryServiceImpl) GetConnection(host string, port string, db int) RedisFactoryInterface {
	return RedisFactory{redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: "",
		DB:       db,
	}),
	}
}

// Ping connection and check for liveness
func (factory RedisFactory) Ping(ctx context.Context) (string, error) {
	return factory.Client.Ping(ctx).Result()
}

// Set set value in redis - expireation 0 means indefinite
func (factory RedisFactory) Set(ctx context.Context, key string, val string) error {
	return factory.Client.Set(ctx, key, val, 0).Err()
}

// SetEx set value in redis - with expiration, expireation 0 means indefinite
func (factory RedisFactory) SetEx(ctx context.Context, key string, val string, ttlSeconds int) error {
	return factory.Client.SetEX(ctx, key, val, time.Duration(ttlSeconds)*time.Second).Err()
}

// SetNx set value in redis if not exists - with expiration, expireation 0 means indefinite
func (factory RedisFactory) SetNX(ctx context.Context, key string, val string, ttlSeconds int) (bool, error) {
	return factory.Client.SetNX(ctx, key, val, time.Duration(ttlSeconds)*time.Second).Result()
}

// Exists check valie exists in redis (response > 0)
func (factory RedisFactory) Exists(ctx context.Context, key string) (int64, error) {
	return factory.Client.Exists(ctx, key).Result()
}

// Get get value from redis, if not exists, redis.nil
func (factory RedisFactory) Get(ctx context.Context, key string) (string, error) {
	return factory.Client.Get(ctx, key).Result()
}

// Ttl check time to live of specific key
func (factory RedisFactory) Ttl(ctx context.Context, key string) (time.Duration, error) {
	return factory.Client.TTL(ctx, key).Result()
}

// Expire change time to leave of specific key
func (factory RedisFactory) Expire(ctx context.Context, key string, seconds int) error {
	return factory.Client.Expire(ctx, key, time.Duration(seconds)*time.Second).Err()
}

// Del Delete key from redis
func (factory RedisFactory) Del(ctx context.Context, key string) error {
	return factory.Client.Del(ctx, key).Err()
}

// HGet get specific value from inner map
func (factory RedisFactory) HGet(ctx context.Context, key string, innerKey string) (string, error) {
	return factory.Client.HGet(ctx, key, innerKey).Result()
}

// HGet set value on inner map
func (factory RedisFactory) HSet(ctx context.Context, key string, innerKey string, val string) error {
	return factory.Client.HSet(ctx, key, innerKey, val).Err()
}

// HDel deletes key from inner map
func (factory RedisFactory) HDel(ctx context.Context, key string, innerKey string) error {
	return factory.Client.HDel(ctx, key, innerKey).Err()
}

// HGetAll get all values from inner key
func (factory RedisFactory) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return factory.Client.HGetAll(ctx, key).Result()
}

// HSetAll set values on inner map
func (factory RedisFactory) HSetAll(ctx context.Context, key string, values map[string]string) error {
	var args []interface{}
	for k, v := range values {
		args = append(args, k, v)
	}
	return factory.Client.HMSet(ctx, key, args).Err()
}

// Keys get keys by pattern
func (factory RedisFactory) Keys(ctx context.Context, pattern string) ([]string, error) {
	return factory.Client.Keys(ctx, pattern).Result()
}
