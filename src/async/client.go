package async

import (
	// "batch-processing/env"
	"batch-processing/env"
	"log"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)


func NewAsynqClient()*asynq.Client{
	envs:=env.NewEnv()
	redisUri:= envs.REDIS_SERVICE_URI
	opt,err:= redis.ParseURL(redisUri)
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println(opt)

	return asynq.NewClient(
	asynq.RedisClientOpt{
		Addr:      opt.Addr,
		Username:  opt.Username,
		Password:  opt.Password,
		DB:        1,
		TLSConfig: opt.TLSConfig,
	},
	)

}