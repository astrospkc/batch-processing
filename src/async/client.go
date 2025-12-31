package async

import (
	// "batch-processing/env"
	"crypto/tls"

	"github.com/hibiken/asynq"
)


func NewAsynqClient()*asynq.Client{
	// envs:=env.NewEnv()
	return asynq.NewClient(asynq.RedisClientOpt{
		Addr: "127.0.0.1:6379",
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	})
}