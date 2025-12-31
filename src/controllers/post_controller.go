package controllers

import (
	"batch-processing/src/connect"
	"batch-processing/src/models"
	"context"

	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreatePost()fiber.Handler{
	return func(c *fiber.Ctx) error{
		userId,err := FetchUserId(c)
	
		if err!=nil{
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":"Failed to fetch user id",
			})
		}

		var body models.Post
		if err = c.BodyParser(&body); err!=nil{
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":"Invalid request body",
			})
		}

		newPost:= models.Post{
			Id:primitive.NewObjectID().Hex(),
			UserId:userId,
			Title:body.Title,
			Content:body.Content,
			LikeCount: 0,
			CreatedAt: time.Now().UTC(),
		}

		
		
		_,err = connect.PostsCollection.InsertOne(context.TODO(),newPost)
		if err!=nil{
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":"Failed to create post",
			})
		}

		return c.JSON(fiber.Map{
			"message":"created  successfully",
			"data":newPost,
		})
	}
}