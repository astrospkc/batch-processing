package async

import (
	"encoding/json"

	"github.com/hibiken/asynq"
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
	return asynq.NewTask(TypePostLike, payload), nil
}