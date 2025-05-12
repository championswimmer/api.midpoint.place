package controllers

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/championswimmer/api.midpoint.place/src/config"
	"github.com/championswimmer/api.midpoint.place/src/db"
	"github.com/championswimmer/api.midpoint.place/src/db/models"
	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/gofiber/fiber/v2"
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
	for i := range code {
		n := lo.Must(rand.Int(rand.Reader, big.NewInt(int64(len(config.GROUP_CODE_CHARSET)))))

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

func (c *GroupsController) CreateGroup(creatorID uint, req *dto.CreateGroupRequest) (*dto.GroupResponse, error) {
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

	return &dto.GroupResponse{
		ID:           group.ID,
		Name:         group.Name,
		Type:         group.Type,
		Code:         group.Code,
		CreatorID:    group.CreatorID,
		MidpointLat:  group.MidpointLatitude,
		MidpointLong: group.MidpointLongitude,
		Radius:       group.Radius,
	}, nil
}

func (c *GroupsController) UpdateGroup(groupID string, req *dto.UpdateGroupRequest) (*dto.GroupResponse, error) {
	var group models.Group
	if err := c.db.Where("id = ?", groupID).First(&group).Error; err != nil {
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
		ID:           group.ID,
		Name:         group.Name,
		Type:         group.Type,
		Code:         group.Code,
		CreatorID:    group.CreatorID,
		MidpointLat:  group.MidpointLatitude,
		MidpointLong: group.MidpointLongitude,
		Radius:       group.Radius,
	}, nil
}
