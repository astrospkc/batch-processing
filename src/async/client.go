package async

import (
	// "batch-processing/env"
	"batch-processing/env"
	"log"

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

	return asynq.NewClient(
	asynq.RedisClientOpt{
		Addr:      opt.Addr,
		Username:  opt.Username,
		Password:  opt.Password,
		DB:        opt.DB,
		TLSConfig: opt.TLSConfig,
	},
	)

}