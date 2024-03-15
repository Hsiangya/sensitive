package factory

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
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
)

func ConnectMongo() *dependences.MongoDBClient {
	mongoURL := viper.GetString("mongo.main_uri")
	mongoOnce.Do(func() {
		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURL))
		fmt.Printf("connect %s \n", mongoURL)
		if err != nil {
			panic(err)
		}
		MongoInstance = &dependences.MongoDBClient{Client: client}
	})
	return MongoInstance
}

func ConnectRedis() *dependences.RedisClient {
	redisURL := viper.GetString("redis.main_uri")
	redisDB := viper.GetInt("redis.main_db")
	redisPassword := viper.GetString("redis.main_pass")

	redisOnce.Do(func() {
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			panic(err)
		}
		opt.DB = redisDB
		opt.Password = redisPassword
		client := redis.NewClient(opt)

		err = client.Ping(context.Background()).Err()
		if err != nil {
			panic(err)
		}

		fmt.Printf("connected to Redis: %s, DB: %d\n", redisURL, redisDB)
		RedisInstance = &dependences.RedisClient{Client: client}
	})

	return RedisInstance
}
