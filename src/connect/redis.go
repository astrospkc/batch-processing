package connect

import (
	"batch-processing/env"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var (
	RCtx = context.Background()
	RedisClient redis.Client
   	
)
func InitRedisConnect(){
	envs:=env.NewEnv()
	redisURI:=envs.REDIS_SERVICE_URI
	
	if redisURI == "" {
		panic("AIVEN_SERVICE_URI not set")
	}

	// Parse the URI directly (handles TLS, password, username, host, port)
	client, err := redis.ParseURL(redisURI)
	if err != nil {
		panic("Failed to initialize Redis client: " + err.Error())
	}

	RedisClient = *redis.NewClient(client)

	// Ping test to confirm successful connection
	err = RedisClient.Ping(RCtx).Err()
	if err != nil {
		panic(fmt.Sprintf("Redis connection ping failed: %v", err))
	}

	fmt.Println("✅ Connected to Redis")

}