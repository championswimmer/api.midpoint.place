package controllers

import (
	"github.com/championswimmer/api.midpoint.place/src/db"
	"github.com/championswimmer/api.midpoint.place/src/db/models"
	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/championswimmer/api.midpoint.place/src/utils/applogger"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
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

	// Check if the user is already in the group and update their location if necessary
	var groupUser models.GroupUser
	err := c.db.Transaction(func(tx *gorm.DB) error {
		applogger.Info("Joining group", groupID, "for user", userID, "transaction started")
		if err := tx.Where("user_id = ? AND group_id = ?", userID, groupID).FirstOrCreate(&groupUser, models.GroupUser{
			UserID:    userID,
			GroupID:   groupID,
			Latitude:  req.Latitude,
			Longitude: req.Longitude,
		}).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to add/update user in group")
		}

		// If the user is already in the group and the location has changed, update the location
		if groupUser.Latitude != req.Latitude || groupUser.Longitude != req.Longitude {
			groupUser.Latitude = req.Latitude
			groupUser.Longitude = req.Longitude
			if err := tx.Save(&groupUser).Error; err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "Failed to update user location in group")
			}
			applogger.Warn("User", userID, "is already in group", groupID, "- updating location")
		}
		applogger.Info("Joining group", groupID, "for user", userID, "transaction completed")
		return nil
	})

	if err != nil {
		return nil, err
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
	var result struct {
		Latitude  float64
		Longitude float64
	}
	if err := c.db.Model(&models.GroupUser{}).
		Select("AVG(latitude) as Latitude, AVG(longitude) as Longitude").
		Where("group_id = ?", groupID).
		Scan(&result).Error; err != nil {
		return 0, 0, fiber.NewError(fiber.StatusInternalServerError, "Failed to calculate group centroid")
	}

	return result.Latitude, result.Longitude, nil
}

// GroupMembershipCheck checks if a user belongs to a specified group
// Returns true if the user is a member of the group, false otherwise
func (c *GroupUsersController) GroupMembershipCheck(groupID string, userID uint) (bool, error) {
	var groupUser models.GroupUser
	result := c.db.Where("user_id = ? AND group_id = ?", userID, groupID).First(&groupUser)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// User is not in the group
			return false, nil
		}
		// Return other errors
		return false, fiber.NewError(fiber.StatusInternalServerError, "Failed to check group membership")
	}

	// User is in the group
	return true, nil
}

// GetGroupMembers fetches all members of a group
// Returns a slice of GroupUserResponse for all users in the group
func (c *GroupUsersController) GetGroupMembers(groupID string) ([]dto.GroupUserResponse, error) {
	var groupUsers []models.GroupUser
	if err := c.db.Where("group_id = ?", groupID).Find(&groupUsers).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch group members")
	}

	response := make([]dto.GroupUserResponse, len(groupUsers))
	for i, groupUser := range groupUsers {
		response[i] = dto.GroupUserResponse{
			UserID:    groupUser.UserID,
			GroupID:   groupUser.GroupID,
			Latitude:  groupUser.Latitude,
			Longitude: groupUser.Longitude,
		}
	}

	return response, nil
}

// GetGroupsContainingMember returns all groups that contain the specified user as a member
func (c *GroupUsersController) GetGroupsContainingMember(userID uint) ([]dto.GroupResponse, error) {
	var groups []models.Group
	if err := c.db.Preload("Creator").
		Preload("Members").
		Preload("Members.User").
		Joins("JOIN group_users ON group_users.group_id = groups.id AND group_users.deleted_at IS NULL").
		Where("group_users.user_id = ?", userID).
		Find(&groups).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch groups containing member")
	}

	response := make([]dto.GroupResponse, len(groups))
	for i, group := range groups {
		members := lo.Map(group.Members, func(member models.GroupUser, _ int) dto.GroupUserResponse {
			return dto.GroupUserResponse{
				UserID:    member.UserID,
				GroupID:   member.GroupID,
				Latitude:  member.Latitude,
				Longitude: member.Longitude,
				Role:      member.Role,
			}
		})

		response[i] = dto.GroupResponse{
			ID:                group.ID,
			Name:              group.Name,
			Type:              group.Type,
			Code:              group.Code,
			Creator:           dto.GroupCreator{ID: group.Creator.ID, DisplayName: group.Creator.DisplayName},
			MidpointLatitude:  group.MidpointLatitude,
			MidpointLongitude: group.MidpointLongitude,
			Radius:            group.Radius,
			Members:           members,
		}
	}

	return response, nil
}
