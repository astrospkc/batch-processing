package controllers

import (
	"batch-processing/src/connect"
	"batch-processing/src/models"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
)


type LikePostResponse struct {
	Message string `bson:"message" json:"message"`
	Code    int  `bson:"code" json:"code"`
	Data    models.Post `bson:"data" json:"data"`
}

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