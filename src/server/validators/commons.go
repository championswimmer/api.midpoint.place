package validators

import (
	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/gofiber/fiber/v2"
)

type ValidationError struct {
	status  int
	message string
}

func (e *ValidationError) Error() string {
	return e.message
}

func (e *ValidationError) ErrorDetails() (int, string) {
	return e.status, e.message
}

func SendValidationError(ctx *fiber.Ctx, err *ValidationError) error {
	return ctx.Status(err.status).JSON(dto.CreateErrorResponse(err.ErrorDetails()))
}
