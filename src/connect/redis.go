package connect

import (
	"batch-processing/env"
	"context"
	"fmt"

	"github.com/valkey-io/valkey-go"
)

var (
	RCtx = context.Background()
	RedisClient valkey.Client
   	
)
func InitRedisConnect(){
	envs:=env.NewEnv()
	redisURI:=envs.REDIS_SERVICE_URI
	
	if redisURI == "" {
		panic("AIVEN_SERVICE_URI not set")
	}

	// Parse the URI directly (handles TLS, password, username, host, port)
	client, err := valkey.NewClient(valkey.MustParseURL(redisURI))
	if err != nil {
		panic("Failed to initialize Redis client: " + err.Error())
	}

	RedisClient = client

	// Ping test to confirm successful connection
	err = RedisClient.Do(RCtx, RedisClient.B().Ping().Build()).Error()
	if err != nil {
		panic(fmt.Sprintf("Redis connection ping failed: %v", err))
	}

	fmt.Println("✅ Connected to Redis")

}