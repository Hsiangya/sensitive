package dependences

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"unicode"
)

type Node struct {
	Children map[rune]*Node
	IsEnd    bool
}

type DFATree struct {
	Root *Node
}

func (d *DFATree) AddWord(word string) {
	node := d.Root
	for _, char := range word {
		_, exist := node.Children[char]
		if !exist {
			node.Children[char] = &Node{Children: make(map[rune]*Node)}
		}
		node = node.Children[char]
	}
	node.IsEnd = true
}

func (d *DFATree) CheckChinese(words string) bool {
	for index, char := range words {
		if !unicode.Is(unicode.Han, char) {
			continue
		}
		node := d.Root
		for _, nextChar := range words[index:] {
			if !unicode.Is(unicode.Han, nextChar) {
				continue
			}
			fmt.Println(nextChar)
			nextNode, exists := node.Children[nextChar]
			if !exists {
				break
			}

			if nextNode.IsEnd {
				return true
			}
			node = nextNode
		}

	}
	return false
}

func (d *DFATree) CheckEnglish(words string) bool {
	for index, char := range words {
		if !isEnglishLetter(char) {
			continue
		}

		node := d.Root
		for _, nextChar := range words[index:] {
			if !isEnglishLetter(nextChar) {
				continue
			}
			nextNode, exists := node.Children[nextChar]
			if !exists {
				break
			}

			if nextNode.IsEnd {
				return true
			}
			node = nextNode
		}

	}
	return false
}

func (d *DFATree) LoadSensitiveWord(mongo *MongoDBClient) {
	fmt.Println("begin loading")
	filter := bson.M{}
	ctx := context.Background()
	sensitiveWords := mongo.FindMany(ctx, "public_info", "sensitive", filter)
	for _, wordMapping := range sensitiveWords {
		if text, ok := wordMapping["text"].(string); ok {
			d.AddWord(text)
		}
		fmt.Println(wordMapping["text"].(string))

	}

}

func isEnglishLetter(r rune) bool {
	return ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z')
}
