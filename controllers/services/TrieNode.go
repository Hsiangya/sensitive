package services

import (
	"context"
	"errors"
	"fmt"
	"sensitive/controllers/factory"
)

func InsertWord(ctx context.Context, word string) error {
	currentNode := "root"
	for _, char := range word {
		charStr := string(char)
		endExists, err := factory.RedisInstance.HExists(ctx, currentNode, "end")
		if err != nil {
			return errors.New("error checking for 'end': " + err.Error())
		}
		if endExists {
			return nil
		}
		err = factory.RedisInstance.HSet(ctx, currentNode, charStr, charStr)
		if err != nil {
			return errors.New("Error setting next hash: " + err.Error())
		}
		currentNode = charStr
	}
	err := factory.RedisInstance.HSet(ctx, currentNode, "end", true)
	if err != nil {
		fmt.Println("Error setting 'end':", err)
	}
	return nil
}

func IsSensitive(ctx context.Context, word string) bool {
	node := "root"
	for _, char := range word {
		isValue := factory.RedisInstance.HGet(ctx, node, string(char))
		fmt.Println("isValue %s", isValue)
		node = isValue

		isEnd := factory.RedisInstance.HGet(ctx, node, "end")
		if string(isEnd) == "1" {
			return true
		}
	}
	return false
}
