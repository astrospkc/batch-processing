package connect

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var (
	RCtx = context.Background()
	RedisClient *redis.Client
   	
)
func InitRedisConnect(){
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
	},
	)

	if err := rdb.Ping(RCtx).Err(); err!=nil{
		log.Fatalf("Redis not connected %v", err)
	}

	RedisClient = rdb
	fmt.Print("redis connected ")
}