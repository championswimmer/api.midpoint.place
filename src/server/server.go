package server

import (
	"github.com/championswimmer/api.midpoint.place/src/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// @title Midpoint Place API
// @version 1.0
// @description This is the API for the Midpoint Place project
// @host api.midpoint.place
// @BasePath /v1
// @schemes http https
func CreateServer() *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork: true,
	})

	// enable cors
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	apiV1 := app.Group("/v1")

	apiV1.Route("/users", routes.UsersRoute())
	apiV1.Route("/groups", routes.GroupsRoute())

	return app
}
