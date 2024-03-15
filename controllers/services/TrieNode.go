package services

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"sensitive/controllers/factory"
)

func InsertWord(ctx context.Context, word string) error {
	node := "root"
	for _, char := range word {
		key := string(char)
		isEnd, err := factory.RedisInstance.Client.HGet(ctx, key, "isEnd").Int64()
		if err != nil && !errors.Is(redis.Nil, err) {
			return err
		}
		if isEnd == 1 {
			return nil
		}
		_, err = factory.RedisInstance.Client.HSetNX(ctx, key, "isEnd", 0).Result()
		if err != nil {
			return err
		}
		node = key
	}
	_, err := factory.RedisInstance.Client.HSet(ctx, node, "isEnd", 1).Result()
	return err
}

func IsSensitive(ctx context.Context, word string) (bool, error) {
	for _, char := range word {
		key := string(char)
		isEnd, err := factory.RedisInstance.Client.HGet(ctx, key, "isEnd").Int64()
		if errors.Is(redis.Nil, err) {
			return false, nil
		}
		if err != nil {
			return false, err
		}
		if isEnd == 1 {
			return true, nil
		}
	}
	return false, nil
}
