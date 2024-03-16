package dependences

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisClient struct {
	Client *redis.Client
}

func (r *RedisClient) Connect() error {
	return r.Client.Ping(context.Background()).Err()
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *RedisClient) HGet(ctx context.Context, key string, field string) string {
	result, _ := r.Client.HGet(ctx, key, field).Result()
	return result
}
func (r *RedisClient) HSet(ctx context.Context, key string, field string, value interface{}) error {
	return r.Client.HSet(ctx, key, field, value).Err()
}

func (r *RedisClient) HSetNX(ctx context.Context, key string, field string, value interface{}) (bool, error) {
	return r.Client.HSetNX(ctx, key, field, value).Result()
}

func (r *RedisClient) HExists(ctx context.Context, key string, field string) (bool, error) {
	return r.Client.HExists(ctx, key, field).Result()
}

func (r *RedisClient) Disconnect(ctx context.Context) error {
	return r.Client.Close()
}
