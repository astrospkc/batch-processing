package routes

import (
	"batch-processing/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func NormalRoutes(app *fiber.App){
	
	auth := app.Group("/auth")
	auth.Post("/", controllers.CreateUser())
}