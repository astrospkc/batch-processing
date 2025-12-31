package routes

import (
	"batch-processing/src/controllers"
	"batch-processing/src/middleware"

	"github.com/gofiber/fiber/v2"
)

func NormalRoutes(app *fiber.App){
	
	auth := app.Group("/auth")
	auth.Post("/", controllers.CreateUser())

	post:=app.Group("post", middleware.ValidateAPIKey())
	post.Post("/", controllers.CreatePost())
	
	postLike := app.Group("/like", middleware.ValidateAPIKey())
	postLike.Post("/:post_id", controllers.LikePost())
}