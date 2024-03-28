package factory

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sensitive/controllers/dependences"
	"sync"
)

var (
	MongoInstance *dependences.MongoDBClient
	mongoOnce     sync.Once
	RedisInstance *dependences.RedisClient
	redisOnce     sync.Once
	DfaInstance   *dependences.DFATree
	dfaOnce       sync.Once
)

func CreateMongoApp(Url string) *dependences.MongoDBClient {
	mongoOnce.Do(func() {
		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(Url))
		fmt.Printf("connect %s \n", Url)
		if err != nil {
			panic(err)
		}
		MongoInstance = &dependences.MongoDBClient{Client: client}
	})
	return MongoInstance
}

func CreateRedisApp(url string, db int, password string) *dependences.RedisClient {
	redisOnce.Do(func() {
		opt, err := redis.ParseURL(url)
		if err != nil {
			panic(err)
		}
		opt.DB = db
		opt.Password = password
		client := redis.NewClient(opt)

		err = client.Ping(context.Background()).Err()
		if err != nil {
			panic(err)
		}

		fmt.Printf("connected to Redis: %s, DB: %d\n", url, db)
		RedisInstance = &dependences.RedisClient{Client: client}
	})
	return RedisInstance
}

func CreateDFATree() *dependences.DFATree {
	dfaOnce.Do(func() {
		DfaInstance = &dependences.DFATree{Root: &dependences.Node{Children: make(map[rune]*dependences.Node)}}

		DfaInstance.LoadSensitiveWord(MongoInstance)
	})
	return DfaInstance
}
