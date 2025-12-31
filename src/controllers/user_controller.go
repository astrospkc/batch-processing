package controllers

import (
	"batch-processing/src/connect"
	"batch-processing/src/models"
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserCreatedResponse struct{
	Message 	string		`bson:"message" json:"message"`
	Email		string 		`bson:"email" json:"email"`
	Code        int         `bson:"code" json:"code"`
	Data        models.User `bson:"data" json:"data"`
}

func GenerateApiKey() (string, error){
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err!=nil{
		return "",nil
	}
	return hex.EncodeToString(bytes), nil
}




func CreateUser()fiber.Handler{
	return func(c *fiber.Ctx)error{
		var body models.User
		if err:=c.BodyParser(&body);err!=nil{
			return c.Status(fiber.StatusBadRequest).JSON(
				fiber.Map{
					"error":"Invalid request body ",
				},
			)
		}
		// check if existing user
		var existingUser models.User
		err := connect.UsersCollection.FindOne(context.TODO(), bson.M{"email":body.Email}).Decode(&existingUser)
		if err== nil {
			return c.Status(fiber.StatusBadRequest).JSON(UserCreatedResponse{
				Message: "email already in use",
				Email:   body.Email,
				Code:    fiber.StatusBadRequest,
				
			})
		}

		apikey, err := GenerateApiKey()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(UserCreatedResponse{
				Message: "Failed to generate api key",
				Email:    body.Email,
				Code:    fiber.StatusInternalServerError,
			})
		}

		newuser:=models.User{
			Id:primitive.NewObjectID().Hex(),
			Name:body.Name,
			Email:body.Email,
			API_Key: apikey,
			CreatedAt: time.Now().UTC(),
		}
		_,err = connect.UsersCollection.InsertOne(context.TODO(),newuser)
		if err!=nil{
			return c.Status(fiber.StatusBadRequest).JSON(UserCreatedResponse{
				Message: "Failed to create user",
				Email:   body.Email,
				Code:    fiber.StatusBadRequest,

			})
		}
		
		
		return c.JSON(UserCreatedResponse{
			Message: "User Created Successfully",
			Email: body.Email,
			Code: fiber.StatusAccepted,
			Data: newuser,
		})


	}
}
