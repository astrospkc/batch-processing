package async

import (
	// "batch-processing/env"
	"context"
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"
)

func NewServer() *asynq.Server{
	// envs:= env.NewEnv()
	return asynq.NewServer(asynq.RedisClientOpt{
		Addr:"127.0.0.1:6379",
	}, 
	asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			"critical":6, 
			"default":3,
			"low":1,
		},
	},
)
}


func NewMux()*asynq.ServeMux{
	mux:=asynq.NewServeMux()
	mux.HandleFunc(TypePostLike, HandlePostLike)
	return mux
}

func HandlePostLike(ctx context.Context, t *asynq.Task) error{
	var payload PostLikePayload
	if err:= json.Unmarshal(t.Payload(),&payload); err!=nil{
		return err
	}
	log.Println("her comes the payload: ", payload.LikeCount, payload.PostId)

	return nil
}