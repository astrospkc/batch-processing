package controllers

import (
	"batch-processing/src/async"
	"batch-processing/src/connect"
	"batch-processing/src/models"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
)


type LikePostResponse struct {
	Message string `bson:"message" json:"message"`
	Code    int  `bson:"code" json:"code"`
	Data    models.Post `bson:"data" json:"data"`
}



type PostLikePayload struct{
	PostId string
	LikeCount int
}

var AsynqClient *asynq.Client

func LikePost()fiber.Handler{
	return func(c *fiber.Ctx) error{
		postId:= c.Params("post_id")
		userId,err:= FetchUserId(c)
		
		if err!=nil{
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":"Failed to fetch user id",
			})
		}
		postlike:=models.Likes{
			Id: primitive.NewObjectID().Hex(),
			PostId: postId,
			UserId: userId,
			CreatedAt: time.Now().UTC(),
		}

		_,err =connect.PostLikesCollection.InsertOne(context.TODO(),postlike)
		if err!=nil{
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":"Failed to insert like or like",
			})
		}

		// now after insertion , increase count
		// update post likecount
		filter:=bson.M{
			"id":postId,
		}

		update:=bson.M{
			"$inc":bson.M{
				"like_count":1,
			},
		}

		
		var updatedPost models.Post
		if err = connect.PostsCollection.FindOneAndUpdate(context.TODO(),filter, update).Decode(&updatedPost); err!=nil{
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":"failed to update like count",
			})
		}

		

		return c.JSON(LikePostResponse{
			Message: "post like is successfull",
			Code:fiber.StatusOK,
			Data:updatedPost,
		})

	}
}


func MiscLikePost() fiber.Handler {
	return func(c *fiber.Ctx) error {
		AsynqClient := async.NewAsynqClient()
		// fmt.Println("AsynqClient",AsynqClient)
		AsynqClient.Enqueue(asynq.NewTask("debug:test", nil))

		// fmt.Println("why")
		postId := c.Params("post_id")
		_, err := FetchUserId(c)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to fetch user_id",
			})
		}
		var postLike models.Post
		filter := bson.M{
			"id": postId,
		}
		if err = connect.PostsCollection.FindOne(context.TODO(), filter).Decode(&postLike); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to fetch post",
			})
		}

		// Redis Key
		key := fmt.Sprintf("PostId_%s", postId)
		// Fixed TTL of 12 hours
		ttl := 12 * time.Hour

		// 1. Initialize key if not exists (SetNX) using current DB count
		// If key exists, this does nothing and returns false (which we ignore here)
		_, err = connect.RedisClient.SetNX(connect.RCtx, key, postLike.LikeCount, ttl).Result()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to access redis",
			})
		}

		// 2. Atomically Increment the count
		newCount, err := connect.RedisClient.Incr(connect.RCtx, key).Result()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to increment like count",
			})
		}

		// 3. Enqueue background task to sync to DB
		// Determine wait time or batching strategy. User had 10s unique lock.
		task := asynq.NewTask(
			"post:like",
			[]byte(fmt.Sprintf(`{"post_id":"%s"}`, postId)),
		)

		// Enqueue with uniqueness to avoid swamping the worker if likes come in fast bursts
		_, err = AsynqClient.Enqueue(
			task,
			asynq.ProcessIn(10*time.Second), // Debounce/Batch window
			asynq.Unique(10*time.Second),    // Dedup task for this post for 10s
		)
		if err != nil {
			// Log error but don't fail the request since Redis is updated
			log.Printf("Failed to enqueue task: %v", err)
		}

		return c.JSON(fiber.Map{
			"message":    "successfully done",
			"data":       "queued",
			"like count": newCount,
		})
	}
}




// docker run --rm \
//   --name asynqmon \
//   -p 8080:8080 \
//   -e REDIS_ADDR=scheduler-punampandit-5c65.b.aivencloud.com:20072 \
//   -e REDIS_USERNAME=default \
//   -e REDIS_PASSWORD=AVNrFn_QxT \
//   hibiken/asynqmon
// Asynq Monitoring WebUI server is listening on port 8080