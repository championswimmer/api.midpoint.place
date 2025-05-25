package routes

import (
	"strconv"

	"github.com/championswimmer/api.midpoint.place/src/config"
	"github.com/championswimmer/api.midpoint.place/src/controllers"
	"github.com/championswimmer/api.midpoint.place/src/db/models"
	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/championswimmer/api.midpoint.place/src/security"
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
		router.Post("/:userid", security.MandatoryJwtAuthMiddleware, updateUserData)
	}

}

// @Summary Register a new user
// @Description Register a new user
// @Tags users
// @ID register-user
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "User"
// @Success 201 {object} dto.UserResponse "User created successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request"
// @Failure 409 {object} dto.ErrorResponse "Email already exists"
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

	return ctx.Status(fiber.StatusCreated).JSON(user)
}

// @Summary Login a user
// @Description Login a user
// @Tags users
// @ID login-user
// @Accept json
// @Produce json
// @Param user body dto.LoginUserRequest true "User"
// @Success 200 {object} dto.UserResponse "User logged in successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request"
// @Failure 401 {object} dto.ErrorResponse "Invalid credentials"
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
// @Param user body dto.UserUpdateRequest true "User"
// @Success 200 {object} dto.UserResponse "User updated successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request"
// @Failure 403 {object} dto.ErrorResponse "You are not allowed to update this user's data"
// @Router /users/{userid} [post]
// @Security BearerAuth
func updateUserData(ctx *fiber.Ctx) error {
	userID, err := strconv.ParseUint(ctx.Params("userid"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.CreateErrorResponse(fiber.StatusBadRequest, "Invalid user ID"))
	}
	// allow only for self
	if userID != uint64(ctx.Locals(config.LOCALS_USER).(*models.User).ID) {
		return ctx.Status(fiber.StatusForbidden).JSON(dto.CreateErrorResponse(fiber.StatusForbidden, "You are not allowed to update this user's data"))
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
