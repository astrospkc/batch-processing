package middleware

import (
	"batch-processing/src/connect"
	"batch-processing/src/models"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
)


func ValidateAPIKey() fiber.Handler {
	return func(c *fiber.Ctx) error{

	apikey := c.Get("X-API-Key")

	filter := bson.M{
		"key":apikey,
	}
	var u models.User
	err := connect.UsersCollection.FindOne(context.TODO(), filter).Decode(&u)
	if err!=nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":"no user found with this apikey",
		})
	}

	c.Locals("user_id", u.Id)
	return c.Next()
}

}