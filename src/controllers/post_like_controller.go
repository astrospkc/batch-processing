package controllers

import (
	"batch-processing/src/async"
	"batch-processing/src/connect"
	"batch-processing/src/models"
	"context"
	"fmt"
	"strconv"
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


func MiscLikePost()fiber.Handler{
	return func(c *fiber.Ctx) error{
		AsynqClient := async.NewAsynqClient()

		fmt.Println("why")
		postId:= c.Params("post_id")
		_, err:= FetchUserId(c)
		if err!=nil{
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":"Failed to fetch user_id",
			})
		}
		var postLike models.Post
		filter:=bson.M{
			"id":postId,
		}
		if err = connect.PostsCollection.FindOne(context.TODO(),filter ).Decode(&postLike); err!=nil{
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":"Failed to fetch post",
			})
		}
		// add task in redis and later update it to db after some time
		// 1st check if any data related with like is present or not.
		// setnX -> set if not exists
		key := fmt.Sprintf("PostId_%s", postId)
		ok, err := connect.RedisClient.SetNX(connect.RCtx, key,postLike.LikeCount, time.Duration(12*time.Now().Hour())*60*time.Second).Result()
		if err!=nil{
			return err
		}				
		if !ok {
			// key already exists
			// update the value of the key
			var val string
			val, err = connect.RedisClient.Get(context.TODO(), key).Result()
			if err!=nil{
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":"Failed while getting the value of like from redis",
				})
			}
			value,err:= strconv.Atoi(val)
			if err!=nil{
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":"Failed while type conversion",
				})
			}

			ttl:= 12*time.Minute
			result, err := connect.RedisClient.Set(context.TODO(), key, value+1, time.Duration(ttl) ).Result()
			if err!=nil{
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":"Failed to set value",
				})
			}

			task := asynq.NewTask(
				"post:like",
				[]byte(fmt.Sprintf(`{"post_id":"%s"}`, postId)),

			)

			_, _ = AsynqClient.Enqueue(
				task,
				asynq.ProcessIn(10*time.Second), // batch window
				asynq.Unique(10*time.Second), 
			)
			// assign the task 

			return c.JSON(fiber.Map{
				"message":"successfully done",
				"data":result,
				"like count":value+1,
			})
		}

		return c.JSON(fiber.Map{
			"message":"successfully done",
			"data":"",
			"like count":postLike.LikeCount,
		})
		
		
	}
}

