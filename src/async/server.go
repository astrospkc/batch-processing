package async

import (
	// "batch-processing/env"
	"batch-processing/env"
	"batch-processing/src/connect"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type PostLikePayload struct{
	PostId string
}

const (
	TypePostLike = "post:like"
)
func NewServer() *asynq.Server{
	envs:= env.NewEnv()
	redis_uri := envs.REDIS_SERVICE_URI
	opt, err:= redis.ParseURL(redis_uri)
	if err!=nil{
		log.Fatal(err)
	}
	client := asynq.RedisClientOpt{
		Addr:      opt.Addr,
		Username:  opt.Username,
		Password:  opt.Password,
		DB:        opt.DB,
		TLSConfig: opt.TLSConfig,
	}
	return asynq.NewServer(client, 
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
	filter:= bson.M{
		"id":payload.PostId,
	}
	key := fmt.Sprintf("PostId_%s", payload.PostId)
	data, err :=connect.RedisClient.Get(connect.RCtx, key).Int()
	if err == redis.Nil || data == 0 {
		return nil // nothing to flush
	}
	if err != nil {
		return err
	}
	update:=bson.M{
		"$inc":bson.M{
			"like_count":data,
		},
	}
	_, err=connect.PostsCollection.UpdateByID(context.TODO(),filter, update )
	if err!=nil{
		return err
	}
	connect.RedisClient.Del(connect.RCtx, key)

	log.Println("here comes the payload: ",  payload.PostId)

	return nil
}