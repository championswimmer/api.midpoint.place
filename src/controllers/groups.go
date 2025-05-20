package controllers

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/championswimmer/api.midpoint.place/src/config"
	"github.com/championswimmer/api.midpoint.place/src/db"
	"github.com/championswimmer/api.midpoint.place/src/db/models"
	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/championswimmer/api.midpoint.place/src/server/validators"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type GroupsController struct {
	db *gorm.DB
}

func CreateGroupsController() *GroupsController {
	appDb := db.GetAppDB()
	return &GroupsController{
		db: appDb,
	}
}

// generateGroupCode generates a random 10-character alphanumeric string
func generateGroupCode() string {
	// TODO: move this to consts
	code := make([]byte, 10)
	charsetLength := big.NewInt(int64(len(config.GROUP_CODE_CHARSET)))
	for i := range code {
		n := lo.Must(rand.Int(rand.Reader, charsetLength))

		code[i] = config.GROUP_CODE_CHARSET[n.Int64()]
	}
	return string(code)
}

// generateRandomSecret generates a random 6-digit numeric secret
func generateRandomSecret() string {
	// Generate a random number between 100000 and 999999
	n := lo.Must(rand.Int(rand.Reader, big.NewInt(900000)))

	// Add 100000 to ensure 6 digits
	return fmt.Sprintf("%06d", n.Int64()+100000)
}

func (c *GroupsController) GetGroupByIDorCode(groupIDorCode string, includeUsers bool, includePlaces bool) (*dto.GroupResponse, error) {
	var group models.Group

	// Check if input is valid UUID or 10-char alphanumeric code
	isValidUUID := uuid.Validate(groupIDorCode) == nil
	isValidCode := len(groupIDorCode) == 10 && func() bool {
		for _, c := range groupIDorCode {
			if !strings.ContainsRune(config.GROUP_CODE_CHARSET, c) {
				return false
			}
		}
		return true
	}()

	if !isValidUUID && !isValidCode {
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity, "Invalid group ID or code")
	}

	// Always preload creator
	query := c.db.Preload("Creator").Joins("LEFT JOIN users ON groups.creator_id = users.id")

	if includeUsers {
		query = query.Preload("Members").Joins("LEFT JOIN group_users ON groups.id = group_users.group_id")
	}
	if includePlaces {
		query = query.Preload("Places").Joins("LEFT JOIN group_places ON groups.id = group_places.group_id")
	}

	if isValidUUID {
		query = query.Where("groups.id = ?", groupIDorCode)
	} else {
		query = query.Where("groups.code = ?", groupIDorCode)
	}

	if err := query.First(&group).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Group not found")
	}

	groupResponse := &dto.GroupResponse{
		ID:   group.ID,
		Name: group.Name,
		Type: group.Type,
		Code: group.Code,
		Creator: dto.GroupCreator{
			ID:       group.Creator.ID,
			Username: group.Creator.Username,
		},
		MidpointLatitude:  group.MidpointLatitude,
		MidpointLongitude: group.MidpointLongitude,
		Radius:            group.Radius,
	}

	if includeUsers {
		groupResponse.Members = lo.Map(group.Members, func(member models.GroupUser, _ int) dto.GroupUserResponse {
			return dto.GroupUserResponse{
				UserID:    member.UserID,
				GroupID:   member.GroupID,
				Latitude:  member.Latitude,
				Longitude: member.Longitude,
				Role:      member.Role,
			}
		})
	}
	if includePlaces {
		groupResponse.Places = lo.Map(group.Places, func(place models.GroupPlace, _ int) dto.GroupPlaceResponse {
			return dto.GroupPlaceResponse{
				PlaceID:   place.PlaceID,
				GroupID:   place.GroupID,
				Name:      place.Name,
				Address:   place.Address,
				Type:      place.Type,
				Rating:    place.Rating,
				MapURI:    place.MapURI,
				Latitude:  place.Latitude,
				Longitude: place.Longitude,
			}
		})
	}
	return groupResponse, nil
}

func (c *GroupsController) CreateGroup(creatorID uint, req *dto.CreateGroupRequest) (*dto.GroupResponse, error) {
	// Validate request
	if err := validators.ValidateCreateGroupRequest(req); err != nil {
		return nil, err
	}

	// Generate random code and secret if not provided
	code := generateGroupCode()

	secret := req.Secret
	if secret == "" {
		secret = generateRandomSecret()
	}

	// Set default type if not provided
	groupType := req.Type
	if groupType == "" {
		groupType = config.GroupTypePublic
	}

	// Create new group
	group := models.Group{
		ID:        uuid.New().String(),
		CreatorID: creatorID,
		Name:      req.Name,
		Type:      groupType,
		Code:      code,
		Secret:    secret,
		Radius:    req.Radius,
	}

	if err := c.db.Create(&group).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to create group")
	}

	// Fetch the creator information
	var creator models.User
	if err := c.db.First(&creator, creatorID).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch creator information")
	}

	return &dto.GroupResponse{
		ID:   group.ID,
		Name: group.Name,
		Type: group.Type,
		Code: group.Code,
		Creator: dto.GroupCreator{
			ID:       creator.ID,
			Username: creator.Username,
		},
		MidpointLatitude:  group.MidpointLatitude,
		MidpointLongitude: group.MidpointLongitude,
		Radius:            group.Radius,
	}, nil
}

func (c *GroupsController) UpdateGroup(groupID string, req *dto.UpdateGroupRequest) (*dto.GroupResponse, error) {
	var group models.Group
	if err := c.db.Preload("Creator").Where("id = ?", groupID).First(&group).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Group not found")
	}

	// Update only allowed fields
	if req.Name != "" {
		group.Name = req.Name
	}
	if req.Type != "" {
		group.Type = req.Type
	}
	if req.Secret != "" {
		group.Secret = req.Secret
	}
	if req.Radius > 0 {
		group.Radius = req.Radius
	}

	if err := c.db.Save(&group).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to update group")
	}

	return &dto.GroupResponse{
		ID:   group.ID,
		Name: group.Name,
		Type: group.Type,
		Code: group.Code,
		Creator: dto.GroupCreator{
			ID:       group.Creator.ID,
			Username: group.Creator.Username,
		},
		MidpointLatitude:  group.MidpointLatitude,
		MidpointLongitude: group.MidpointLongitude,
		Radius:            group.Radius,
	}, nil
}

func (c *GroupsController) UpdateGroupMidpoint(groupID string, req *dto.UpdateGroupMidpointRequest) (*dto.GroupResponse, error) {
	var group models.Group
	if err := c.db.Preload("Creator").Where("id = ?", groupID).First(&group).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Group not found")
	}

	group.MidpointLatitude = req.Latitude
	group.MidpointLongitude = req.Longitude

	if err := c.db.Save(&group).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to update group location")
	}
	return &dto.GroupResponse{
		ID:   group.ID,
		Name: group.Name,
		Type: group.Type,
		Code: group.Code,
		Creator: dto.GroupCreator{
			ID:       group.Creator.ID,
			Username: group.Creator.Username,
		},
		MidpointLatitude:  group.MidpointLatitude,
		MidpointLongitude: group.MidpointLongitude,
		Radius:            group.Radius,
	}, nil
}

func (c *GroupsController) GetGroupsByCreator(creatorID uint) ([]dto.GroupResponse, error) {
	var groups []models.Group

	// Fetch all groups created by the user, preload creator info
	if err := c.db.Preload("Creator").Where("creator_id = ?", creatorID).Find(&groups).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch groups")
	}

	// Convert to response DTOs
	groupResponses := make([]dto.GroupResponse, len(groups))
	for i, group := range groups {
		groupResponses[i] = dto.GroupResponse{
			ID:   group.ID,
			Name: group.Name,
			Type: group.Type,
			Code: group.Code,
			Creator: dto.GroupCreator{
				ID:       group.Creator.ID,
				Username: group.Creator.Username,
			},
			MidpointLatitude:  group.MidpointLatitude,
			MidpointLongitude: group.MidpointLongitude,
			Radius:            group.Radius,
		}
	}

	return groupResponses, nil
}

func (c *GroupsController) GetPublicGroups(limit int) ([]dto.GroupResponse, error) {
	var groups []models.Group

	if err := c.db.Preload("Creator").
		Where("type = ?", config.GroupTypePublic).
		Order("created_at desc").
		Limit(limit).
		Find(&groups).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch public groups")
	}

	// Convert to response DTOs
	groupResponses := lo.Map(groups, func(group models.Group, _ int) dto.GroupResponse {
		return dto.GroupResponse{
			ID:   group.ID,
			Name: group.Name,
			Type: group.Type,
			Code: group.Code,
			Creator: dto.GroupCreator{
				ID:       group.Creator.ID,
				Username: group.Creator.Username,
			},
			MidpointLatitude:  group.MidpointLatitude,
			MidpointLongitude: group.MidpointLongitude,
			Radius:            group.Radius,
		}
	})

	return groupResponses, nil
}

func (c *GroupsController) GetPublicGroupsByCreator(creatorID uint) ([]dto.GroupResponse, error) {
	var groups []models.Group

	// Fetch all public groups created by the user, preload creator info
	if err := c.db.Preload("Creator").
		Where("groups.creator_id = ? AND groups.type = ?", creatorID, config.GroupTypePublic).
		Find(&groups).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch public groups by creator")
	}

	// Convert to response DTOs
	groupResponses := make([]dto.GroupResponse, len(groups))
	for i, group := range groups {
		groupResponses[i] = dto.GroupResponse{
			ID:   group.ID,
			Name: group.Name,
			Type: group.Type,
			Code: group.Code,
			Creator: dto.GroupCreator{
				ID:       group.Creator.ID,
				Username: group.Creator.Username,
			},
			MidpointLatitude:  group.MidpointLatitude,
			MidpointLongitude: group.MidpointLongitude,
			Radius:            group.Radius,
		}
	}

	return groupResponses, nil
}
