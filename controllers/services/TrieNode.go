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

//func IsSensitive(ctx context.Context, word string) bool {
//	node := "root"
//	for _, char := range word {
//		isValue := factory.RedisInstance.HGet(ctx, node, string(char))
//		fmt.Println("isValue %s", isValue)
//		node = isValue
//
//		isEnd := factory.RedisInstance.HGet(ctx, node, "end")
//		if string(isEnd) == "1" {
//			return true
//		}
//	}
//	return false
//}

func IsSensitive(ctx context.Context, word string) bool {
	for i := 0; i < len(word); i++ {
		node := "root"
		for j := i; j < len(word); j++ {
			char := string(word[j])
			isValue := factory.RedisInstance.HGet(ctx, node, char)
			if isValue == "" {
				break
			}
			node = isValue

			isEnd := factory.RedisInstance.HGet(ctx, node, "end")
			if string(isEnd) == "1" {
				return true
			}
		}
	}
	return false
}

type ACNode struct {
	children map[rune]*ACNode
	failure  *ACNode
	isEnd    bool
}

func NewACNode() *ACNode {
	return &ACNode{
		children: make(map[rune]*ACNode),
		failure:  nil,
		isEnd:    false,
	}
}

type ACSensitiveMatcher struct {
	root *ACNode
}

func NewACSensitiveMatcher() *ACSensitiveMatcher {
	return &ACSensitiveMatcher{
		root: NewACNode(),
	}
}

func (matcher *ACSensitiveMatcher) InsertWord(ctx context.Context, word string) error {
	current := matcher.root
	for _, char := range word {
		if _, ok := current.children[char]; !ok {
			current.children[char] = NewACNode()
		}
		current = current.children[char]
	}
	current.isEnd = true
	return nil
}

func (matcher *ACSensitiveMatcher) buildFailurePointers() {
	queue := make([]*ACNode, 0)
	queue = append(queue, matcher.root)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for char, child := range current.children {
			if current == matcher.root {
				child.failure = matcher.root
			} else {
				failure := current.failure
				for failure != nil {
					if _, ok := failure.children[char]; ok {
						child.failure = failure.children[char]
						break
					}
					failure = failure.failure
				}
				if failure == nil {
					child.failure = matcher.root
				}
			}
			queue = append(queue, child)
		}
	}
}

func (matcher *ACSensitiveMatcher) IsSensitive(ctx context.Context, text string) bool {
	current := matcher.root
	for _, char := range text {
		for current != matcher.root && current.children[char] == nil {
			current = current.failure
		}
		if current.children[char] != nil {
			current = current.children[char]
		}
		if current.isEnd {
			return true
		}
	}
	return false
}
