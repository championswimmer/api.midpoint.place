package routes

import (
	"github.com/championswimmer/api.midpoint.place/src/config"
	"github.com/championswimmer/api.midpoint.place/src/controllers"
	"github.com/championswimmer/api.midpoint.place/src/db/models"
	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/championswimmer/api.midpoint.place/src/security"
	"github.com/championswimmer/api.midpoint.place/src/server/parsers"
	"github.com/championswimmer/api.midpoint.place/src/server/validators"
	"github.com/championswimmer/api.midpoint.place/src/utils/applogger"
	"github.com/gofiber/fiber/v2"
)

var groupsController *controllers.GroupsController
var groupUsersController *controllers.GroupUsersController

func GroupsRoute() func(router fiber.Router) {
	groupsController = controllers.CreateGroupsController()
	groupUsersController = controllers.CreateGroupUsersController()

	return func(router fiber.Router) {
		router.Post("/", security.MandatoryJwtAuthMiddleware, createGroup)
		router.Get("/:groupIdOrCode", security.MandatoryJwtAuthMiddleware, getGroup)
		router.Patch("/:groupIdOrCode", security.MandatoryJwtAuthMiddleware, updateGroup) // Assuming PATCH for partial updates
		router.Put("/:groupIdOrCode/join", security.MandatoryJwtAuthMiddleware, joinGroup)
		router.Delete("/:groupIdOrCode/join", security.MandatoryJwtAuthMiddleware, leaveGroup)
	}
}

// @Summary Create a new group
// @Description Create a new group
// @Tags groups
// @ID create-group
// @Accept json
// @Produce json
// @Param group body dto.CreateGroupRequest true "Group"
// @Success 201 {object} dto.GroupResponse "Group created successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request"
// @Failure 422 {object} dto.ErrorResponse "Group info validation failed"
// @Failure 500 {object} dto.ErrorResponse "Failed to create group"
// @Router /groups [post]
// @Security BearerAuth
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
// @Param groupIdOrCode path string true "Group ID or Code"
// @Param group body dto.UpdateGroupRequest true "Group Update Data"
// @Success 200 {object} dto.GroupResponse "Group updated successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request"
// @Failure 404 {object} dto.ErrorResponse "Group not found"
// @Failure 422 {object} dto.ErrorResponse "Group info validation failed"
// @Failure 500 {object} dto.ErrorResponse "Failed to update group"
// @Router /groups/{groupIdOrCode} [patch]
// @Security BearerAuth
func updateGroup(ctx *fiber.Ctx) error {
	groupID := ctx.Params("groupIdOrCode")

	group, err := groupsController.GetGroupByIDorCode(groupID, false)
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

// @Summary Join a group
// @Description Join an existing group
// @Tags groups
// @ID join-group
// @Produce json
// @Param groupIdOrCode path string true "Group ID or Code"
// @Param groupUser body dto.GroupUserJoinRequest true "Group User"
// @Success 200 {object} dto.GroupUserResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid request"
// @Failure 404 {object} dto.ErrorResponse "Group not found"
// @Failure 500 {object} dto.ErrorResponse "Failed to join group"
// @Router /groups/{groupIdOrCode}/join [put]
// @Security BearerAuth
func joinGroup(ctx *fiber.Ctx) error {
	user := ctx.Locals(config.LOCALS_USER).(*models.User)
	groupIDOrCode := ctx.Params("groupIdOrCode")

	group, err := groupsController.GetGroupByIDorCode(groupIDOrCode, false)
	if err != nil {
		return ctx.Status(err.(*fiber.Error).Code).JSON(dto.CreateErrorResponse(err.(*fiber.Error).Code, err.Error()))
	}

	groupUserReq, parseError := parsers.ParseBody[dto.GroupUserJoinRequest](ctx)
	if parseError != nil {
		return parsers.SendParsingError(ctx, parseError)
	}

	groupUserResp, err := groupUsersController.JoinGroup(group.ID, user.ID, groupUserReq)
	if err != nil {
		return ctx.Status(err.(*fiber.Error).Code).JSON(dto.CreateErrorResponse(err.(*fiber.Error).Code, err.Error()))
	}

	go func() {
		err := _recalculateGroupMidpoint(group.ID)
		if err != nil {
			applogger.Error("Error recalculating group location", err)
		}
	}()

	return ctx.Status(fiber.StatusAccepted).JSON(groupUserResp)
}

// @Summary Leave a group
// @Description Leave an existing group
// @Tags groups
// @ID leave-group
// @Produce json
// @Param groupIdOrCode path string true "Group ID or Code"
// @Success 200 {object} dto.GroupUserResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid request"
// @Failure 404 {object} dto.ErrorResponse "Group not found"
// @Failure 500 {object} dto.ErrorResponse "Failed to leave group"
// @Router /groups/{groupIdOrCode}/join [delete]
// @Security BearerAuth
func leaveGroup(ctx *fiber.Ctx) error {
	user := ctx.Locals(config.LOCALS_USER).(*models.User)
	groupIDOrCode := ctx.Params("groupIdOrCode")

	group, err := groupsController.GetGroupByIDorCode(groupIDOrCode, false)
	if err != nil {
		return ctx.Status(err.(*fiber.Error).Code).JSON(dto.CreateErrorResponse(err.(*fiber.Error).Code, err.Error()))
	}

	err = groupUsersController.LeaveGroup(group.ID, user.ID)
	if err != nil {
		return ctx.Status(err.(*fiber.Error).Code).JSON(dto.CreateErrorResponse(err.(*fiber.Error).Code, err.Error()))
	}
	return ctx.Status(fiber.StatusAccepted).JSON([]byte("{}"))
}

// @Summary Get group information
// @Description Get details of a group by ID or code
// @Tags groups
// @ID get-group
// @Produce json
// @Param groupIdOrCode path string true "Group ID or Code"
// @Param includeUsers query bool false "Include Users"
// @Success 200 {object} dto.GroupResponse
// @Failure 400 {object} dto.ErrorResponse "Invalid request"
// @Failure 404 {object} dto.ErrorResponse "Group not found"
// @Failure 500 {object} dto.ErrorResponse "Failed to get group"
// @Router /groups/{groupIdOrCode} [get]
// @Security BearerAuth
func getGroup(ctx *fiber.Ctx) error {
	groupIDOrCode := ctx.Params("groupIdOrCode")
	includeUsers := ctx.QueryBool("includeUsers", false)

	group, err := groupsController.GetGroupByIDorCode(groupIDOrCode)
	if err != nil {
		return ctx.Status(err.(*fiber.Error).Code).JSON(dto.CreateErrorResponse(err.(*fiber.Error).Code, err.Error()))
	}

	if includeUsers {
		groupUsers, err := groupUsersController.GetGroupMembers(group.ID)
		if err != nil {
			return ctx.Status(err.(*fiber.Error).Code).JSON(dto.CreateErrorResponse(err.(*fiber.Error).Code, err.Error()))
		}
		group.Members = groupUsers
	}

	return ctx.Status(fiber.StatusOK).JSON(group)
}

func _recalculateGroupMidpoint(groupID string) error {
	applogger.Info("Recalculating group midpoint for group", groupID)
	centroidLatitude, centroidLongitude, err := groupUsersController.CalculateGroupCentroid(groupID)
	if err != nil {
		return err
	}
	applogger.Info("Recalculated group midpoint for group", groupID, "to", centroidLatitude, centroidLongitude)

	groupMidpointUpdateRequest := &dto.UpdateGroupMidpointRequest{}
	groupMidpointUpdateRequest.Latitude = centroidLatitude
	groupMidpointUpdateRequest.Longitude = centroidLongitude

	_, err = groupsController.UpdateGroupMidpoint(groupID, groupMidpointUpdateRequest)
	if err != nil {
		return err
	}

	return nil
}
