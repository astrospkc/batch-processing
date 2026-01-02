package async

import (
	"batch-processing/src/connect"
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"go.mongodb.org/mongo-driver/v2/bson"
)


const (
	TypePostLike = "post:like"
)

type PostLikePayload struct{
	PostId string
	LikeCount int
}

func NewPostLikeTask(postId string , likes int)(*asynq.Task,error){
	payload, err:=json.Marshal(PostLikePayload{PostId: postId, LikeCount: likes})
	if err!=nil{
		return nil, err
	}
	// update data in db
	filter:= bson.M{
		"id":postId,
	}

	update:=bson.M{
		"$set":bson.M{
			"like_count":likes,
		},
	}
	_, err=connect.PostsCollection.UpdateByID(context.TODO(),filter, update )
	if err!=nil{
		return nil, err
	}
	return asynq.NewTask(TypePostLike, payload), nil
}