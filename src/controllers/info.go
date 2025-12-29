package controllers

import (
	"github.com/gofiber/fiber/v2"
)


func FetchUserId(c *fiber.Ctx)(string, error){
	var user_id string
	userIdInterface:=c.Locals("user_id")


	user_id, ok := userIdInterface.(string)
	if !ok{
		return user_id, nil
	}
	return user_id, nil

}