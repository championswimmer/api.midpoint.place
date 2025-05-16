package routes

import (
	"github.com/championswimmer/api.midpoint.place/src/controllers"
	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/gofiber/fiber/v2"
)

var waitlistController *controllers.WaitlistController

func WaitlistRoute() func(router fiber.Router) {
	waitlistController = controllers.CreateWaitlistController()

	return func(router fiber.Router) {
		router.Post("/signup", addToWaitlist)
	}
}

// @Summary Add a user to the waitlist
// @Description Add a user to the waitlist
// @Tags waitlist
// @ID add-to-waitlist
// @Accept json
// @Produce json
// @Param user body dto.WaitlistSignupRequest true "User"
// @Success 201 {object} dto.WaitlistSignupResponse "User added to waitlist successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request"
// @Failure 409 {object} dto.ErrorResponse "Email already signed up"
// @Failure 500 {object} dto.ErrorResponse "Failed to add to waitlist"
// @Router /waitlist/signup [post]
func addToWaitlist(ctx *fiber.Ctx) error {
	var req dto.WaitlistSignupRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	res, err := waitlistController.AddToWaitlist(&req)
	if err != nil {
		return ctx.Status(err.(*fiber.Error).Code).JSON(dto.CreateErrorResponse(err.(*fiber.Error).Code, err.Error()))
	}
	return ctx.Status(fiber.StatusCreated).JSON(res)
}
