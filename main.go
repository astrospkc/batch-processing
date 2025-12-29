package main

import (
	"batch-processing/src/connect"
	"batch-processing/src/routes"

	"github.com/gofiber/fiber/v2"
)


func main(){
	app := fiber.New()
    connect.MongoConnect()
    app.Get("/", func (c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })

    routes.NormalRoutes(app)
    app.Listen(":3000")
}