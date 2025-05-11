package routes

import (
	"strconv"

	"github.com/championswimmer/api.midpoint.place/src/controllers"
	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/championswimmer/api.midpoint.place/src/server/parsers"
	"github.com/championswimmer/api.midpoint.place/src/server/validators"
	"github.com/gofiber/fiber/v2"
)

var usersController *controllers.UsersController

func UsersRoute() func(router fiber.Router) {

	usersController = controllers.CreateUsersController()

	return func(router fiber.Router) {
		router.Post("/", registerUser)
		router.Post("/login", loginUser)
		router.Post("/:userid", updateUserData)
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

	user, err := usersController.CreateUser(u)
	if err != nil {
		return ctx.Status(err.(*fiber.Error).Code).JSON(dto.CreateErrorResponse(err.(*fiber.Error).Code, err.Error()))
	}

	// TODO: Generate JWT token here when implementing authentication
	// TODO: Do not return password hash
	return ctx.Status(fiber.StatusCreated).JSON(user)
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
	u, parseError := parsers.ParseBody[dto.LoginUserRequest](ctx)
	if parseError != nil {
		return parsers.SendParsingError(ctx, parseError)
	}

	validateErr := validators.ValidateLoginUserRequest(u)
	if validateErr != nil {
		return validators.SendValidationError(ctx, validateErr)
	}

	user, err := usersController.LoginUser(u)
	if err != nil {
		return ctx.Status(err.(*fiber.Error).Code).JSON(dto.CreateErrorResponse(err.(*fiber.Error).Code, err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(user)
}

// @Summary Update user location
// @Description Update location details for a user
// @Tags users
// @ID update-user-location
// @Accept json
// @Produce json
// @Param userid path string true "User ID"
// @Param location body dto.UserUpdateLocationRequest true "Location"
// @Router /users/{userid} [post]
func updateUserData(ctx *fiber.Ctx) error {
	userID, err := strconv.ParseUint(ctx.Params("userid"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.CreateErrorResponse(fiber.StatusBadRequest, "Invalid user ID"))
	}

	req, parseError := parsers.ParseBody[dto.UserUpdateRequest](ctx)
	if parseError != nil {
		return parsers.SendParsingError(ctx, parseError)
	}

	validateErr := validators.ValidateLocation(req.Location)
	if validateErr != nil {
		return validators.SendValidationError(ctx, validateErr)
	}

	user, err := usersController.UpdateUserLocation(uint(userID), req)
	if err != nil {
		return ctx.Status(err.(*fiber.Error).Code).JSON(dto.CreateErrorResponse(err.(*fiber.Error).Code, err.Error()))
	}

	return ctx.Status(fiber.StatusAccepted).JSON(user)
}
