package services

import (
	"context"
	"errors"
	"fmt"
	"sensitive/controllers/factory"
)

//func InsertWord(ctx context.Context, word string) error {
//	currentHash := "root"
//	for _, char := range word {
//		charStr := string(char)
//		endExists, err := factory.RedisInstance.HExists(ctx, currentHash, "end")
//		if err != nil {
//			return errors.New("error checking for 'end': " + err.Error())
//		}
//		if endExists {
//			return nil
//		}
//		nextHash, err := factory.RedisInstance.HGet(ctx, currentHash, charStr)
//		if errors.Is(redis.Nil, err) {
//			nextHash = currentHash + ":" + charStr
//			err := factory.RedisInstance.HSet(ctx, currentHash, charStr, nextHash)
//			if err != nil {
//				return errors.New("Error setting next hash: " + err.Error())
//			}
//		} else if err != nil {
//			return errors.New("Error getting next hash: " + err.Error())
//		}
//		currentHash = nextHash
//	}
//	err := factory.RedisInstance.HSet(ctx, currentHash, "end", true)
//	if err != nil {
//		fmt.Println("Error setting 'end':", err)
//	}
//	return nil
//}

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
