package controllers

import (
	"github.com/championswimmer/api.midpoint.place/src/db"
	"github.com/championswimmer/api.midpoint.place/src/db/models"
	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/championswimmer/api.midpoint.place/src/utils/applogger"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type GroupUsersController struct {
	db *gorm.DB
}

func CreateGroupUsersController() *GroupUsersController {
	appDb := db.GetAppDB()
	return &GroupUsersController{
		db: appDb,
	}
}

// JoinGroup adds a user to a group with specified location
// This operation is idempotent - if the user is already in the group,
// their location will be updated and a warning will be logged
func (c *GroupUsersController) JoinGroup(req *dto.GroupUserJoinRequest) (*dto.GroupUserResponse, error) {
	// Check if user exists
	var user models.User
	if err := c.db.First(&user, req.UserID).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	// Check if group exists
	var group models.Group
	if err := c.db.First(&group, "id = ?", req.GroupID).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Group not found")
	}

	// Check if the user is already in the group
	var existingGroupUser models.GroupUser
	result := c.db.Where("user_id = ? AND group_id = ?", req.UserID, req.GroupID).First(&existingGroupUser)

	if result.Error == nil {
		// User is already in the group, update their location
		applogger.Warn("User", req.UserID, "is already in group", req.GroupID, "- updating location")

		existingGroupUser.Latitude = req.Latitude
		existingGroupUser.Longitude = req.Longitude

		if err := c.db.Save(&existingGroupUser).Error; err != nil {
			return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to update user location in group")
		}

		return &dto.GroupUserResponse{
			UserID:    existingGroupUser.UserID,
			GroupID:   existingGroupUser.GroupID,
			Latitude:  existingGroupUser.Latitude,
			Longitude: existingGroupUser.Longitude,
		}, nil
	}

	// Add user to group
	groupUser := models.GroupUser{
		UserID:    req.UserID,
		GroupID:   req.GroupID,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}

	if err := c.db.Create(&groupUser).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to add user to group")
	}

	return &dto.GroupUserResponse{
		UserID:    groupUser.UserID,
		GroupID:   groupUser.GroupID,
		Latitude:  groupUser.Latitude,
		Longitude: groupUser.Longitude,
	}, nil
}

// LeaveGroup removes a user from a group
// This operation is idempotent - if the user is not in the group,
// a warning will be logged but no error will be returned
func (c *GroupUsersController) LeaveGroup(req *dto.GroupUserLeaveRequest) error {
	// Check if the mapping exists
	var groupUser models.GroupUser
	result := c.db.Where("user_id = ? AND group_id = ?", req.UserID, req.GroupID).First(&groupUser)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// User is not in the group, log a warning but don't return an error
			applogger.Warn("User", req.UserID, "is not in group", req.GroupID, "- no action needed")
			return nil
		}
		// Return other errors
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to query user in group")
	}

	// Remove the user from the group
	if err := c.db.Delete(&groupUser).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to remove user from group")
	}

	return nil
}
