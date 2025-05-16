package routes

import (
	"github.com/championswimmer/api.midpoint.place/src/controllers"
	"github.com/gofiber/fiber/v2"
)

var waitlistController *controllers.WaitlistController

func WaitlistRoute() func(router fiber.Router) {
	waitlistController = controllers.CreateWaitlistController()

	return func(router fiber.Router) {
		router.Post("/signup", addToWaitlist)
	}
}

func addToWaitlist(ctx *fiber.Ctx) error {
	return waitlistController.AddToWaitlist(ctx)
}
