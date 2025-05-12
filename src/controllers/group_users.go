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
func (c *GroupUsersController) JoinGroup(groupID string, userID uint, req *dto.GroupUserJoinRequest) (*dto.GroupUserResponse, error) {
	// Check if user exists
	var user models.User
	if err := c.db.First(&user, userID).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	// Check if group exists
	var group models.Group
	if err := c.db.First(&group, "id = ?", groupID).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Group not found")
	}

	// Check if the user is already in the group
	var existingGroupUser models.GroupUser
	result := c.db.Where("user_id = ? AND group_id = ?", userID, groupID).First(&existingGroupUser)

	if result.Error == nil {
		// User is already in the group, update their location
		applogger.Warn("User", userID, "is already in group", groupID, "- updating location")

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
		UserID:    userID,
		GroupID:   groupID,
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
func (c *GroupUsersController) LeaveGroup(groupID string, userID uint) error {
	// Check if the mapping exists
	var groupUser models.GroupUser
	result := c.db.Where("user_id = ? AND group_id = ?", userID, groupID).First(&groupUser)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// User is not in the group, log a warning but don't return an error
			applogger.Warn("User", userID, "is not in group", groupID, "- no action needed")
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

func (c *GroupUsersController) CalculateGroupCentroid(groupID string) (latitude float64, longitude float64, err error) {
	groupMembers := []models.GroupUser{}

	if err := c.db.Where("group_id = ?", groupID).Find(&groupMembers).Error; err != nil {
		return 0, 0, fiber.NewError(fiber.StatusInternalServerError, "Failed to get group members")
	}

	totalMembers := len(groupMembers)
	totalLatitude := 0.0
	totalLongitude := 0.0

	for _, member := range groupMembers {
		totalLatitude += member.Latitude
		totalLongitude += member.Longitude
	}

	centroidLatitude := totalLatitude / float64(totalMembers)
	centroidLongitude := totalLongitude / float64(totalMembers)

	return centroidLatitude, centroidLongitude, nil
}
