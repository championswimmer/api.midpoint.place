package parsers

import (
	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/gofiber/fiber/v2"
)

type ParsingError struct {
	status  int
	message string
}

func (e *ParsingError) Error() string {
	return e.message
}

func (e *ParsingError) ErrorDetails() (int, string) {
	return e.status, e.message
}

var bodyParsingError = &ParsingError{
	status:  fiber.StatusBadRequest,
	message: "The request body is not valid",
}

func SendParsingError(ctx *fiber.Ctx, err *ParsingError) error {
	return ctx.Status(err.status).JSON(dto.CreateErrorResponse(err.ErrorDetails()))

}

func ParseBody[T any](ctx *fiber.Ctx) (*T, *ParsingError) {
	var body = new(T)
	if err := ctx.BodyParser(body); err != nil {
		return nil, bodyParsingError
	}
	return body, nil
}
