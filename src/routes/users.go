package routes

import (
	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/championswimmer/api.midpoint.place/src/server/parsers"
	"github.com/championswimmer/api.midpoint.place/src/server/validators"
	"github.com/gofiber/fiber/v2"
)

func UsersRoute() func(router fiber.Router) {

	return func(router fiber.Router) {
		router.Post("/", registerUser)
		router.Post("/login", loginUser)
	}

}

// @Summary Register a new user
// @Description Register a new user
// @Tags users
// @ID register-user
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "User"
// @Router /users [post]
func registerUser(ctx *fiber.Ctx) error {

	u, parseError := parsers.ParseBody[dto.CreateUserRequest](ctx)
	if parseError != nil {
		return parsers.SendParsingError(ctx, parseError)
	}

	validateErr := validators.ValidateCreateUserRequest(u)
	if validateErr != nil {
		return validators.SendValidationError(ctx, validateErr)
	}

	return ctx.Status(fiber.StatusCreated).JSON(u)
}

// @Summary Login a user
// @Description Login a user
// @Tags users
// @ID login-user
// @Accept json
// @Produce json
// @Param user body dto.LoginUserRequest true "User"
// @Router /users/login [post]
func loginUser(ctx *fiber.Ctx) error {
	return ctx.SendString("PLACEHOLDER: login user")
}
