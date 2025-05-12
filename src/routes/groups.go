package routes

import (
	"github.com/championswimmer/api.midpoint.place/src/config"
	"github.com/championswimmer/api.midpoint.place/src/controllers"
	"github.com/championswimmer/api.midpoint.place/src/db/models"
	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/championswimmer/api.midpoint.place/src/security"
	"github.com/championswimmer/api.midpoint.place/src/server/parsers"
	"github.com/championswimmer/api.midpoint.place/src/server/validators"
	"github.com/gofiber/fiber/v2"
)

var groupsController *controllers.GroupsController

func GroupsRoute() func(router fiber.Router) {
	groupsController = controllers.CreateGroupsController()

	return func(router fiber.Router) {
		router.Post("/", security.MandatoryJwtAuthMiddleware, createGroup)
		router.Patch("/:groupid", security.MandatoryJwtAuthMiddleware, updateGroup) // Assuming PATCH for partial updates
	}
}

// @Summary Create a new group
// @Description Create a new group
// @Tags groups
// @ID create-group
// @Accept json
// @Produce json
// @Param group body dto.CreateGroupRequest true "Group"
// @Success 201 {object} dto.GroupResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 422 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /groups [post]
func createGroup(ctx *fiber.Ctx) error {
	user := ctx.Locals(config.LOCALS_USER).(*models.User)

	req, parseError := parsers.ParseBody[dto.CreateGroupRequest](ctx)
	if parseError != nil {
		return parsers.SendParsingError(ctx, parseError)
	}

	validateErr := validators.ValidateCreateGroupRequest(req)
	if validateErr != nil {
		return validators.SendValidationError(ctx, validateErr)
	}

	group, err := groupsController.CreateGroup(user.ID, req)
	if err != nil {
		return ctx.Status(err.(*fiber.Error).Code).JSON(dto.CreateErrorResponse(err.(*fiber.Error).Code, err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).JSON(group)
}

// @Summary Update an existing group
// @Description Update an existing group's details
// @Tags groups
// @ID update-group
// @Accept json
// @Produce json
// @Param groupid path string true "Group ID"
// @Param group body dto.UpdateGroupRequest true "Group Update Data"
// @Success 200 {object} dto.GroupResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 422 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /groups/{groupid} [patch]
func updateGroup(ctx *fiber.Ctx) error {
	groupID := ctx.Params("groupid")

	group, err := groupsController.GetGroupByID(groupID)
	if err != nil {
		return ctx.Status(err.(*fiber.Error).Code).JSON(dto.CreateErrorResponse(err.(*fiber.Error).Code, err.Error()))
	}

	req, parseError := parsers.ParseBody[dto.UpdateGroupRequest](ctx)
	if parseError != nil {
		return parsers.SendParsingError(ctx, parseError)
	}

	validateErr := validators.ValidateUpdateGroupRequest(req)
	if validateErr != nil {
		return validators.SendValidationError(ctx, validateErr)
	}

	group, err = groupsController.UpdateGroup(group.ID, req)
	if err != nil {
		return ctx.Status(err.(*fiber.Error).Code).JSON(dto.CreateErrorResponse(err.(*fiber.Error).Code, err.Error()))
	}

	return ctx.Status(fiber.StatusAccepted).JSON(group)
}
