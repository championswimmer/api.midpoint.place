package tests

import (
	"github.com/championswimmer/api.midpoint.place/src/server"
	"github.com/gofiber/fiber/v2"
)

var App *fiber.App

func init() {
	App = server.CreateServer()
}
