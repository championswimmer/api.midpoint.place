package controllers

import (
	"github.com/championswimmer/api.midpoint.place/src/config"
	"github.com/championswimmer/api.midpoint.place/src/db"
	"github.com/championswimmer/api.midpoint.place/src/db/models"
	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/championswimmer/api.midpoint.place/src/utils/applogger"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GroupPlacesController struct {
	db *gorm.DB
}

func CreateGroupPlacesController() *GroupPlacesController {
	appDb := db.GetAppDB()
	return &GroupPlacesController{
		db: appDb,
	}
}

// AddPlacesToGroup adds an array of places to a group
// If a place with the same PlaceId already exists for this group, it will be ignored
func (c *GroupPlacesController) AddPlacesToGroup(groupID string, req *dto.GroupPlacesAddRequest) ([]dto.GroupPlaceResponse, error) {
	// Check if group exists
	var group models.Group
	if err := c.db.First(&group, "id = ?", groupID).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Group not found")
	}

	// Create transaction to ensure atomicity
	var responses []dto.GroupPlaceResponse
	err := c.db.Transaction(func(tx *gorm.DB) error {
		// First, check for existing places to avoid duplicates
		var existingPlaceIDs []string
		if err := tx.Model(&models.GroupPlace{}).
			Where("group_id = ?", groupID).
			Pluck("place_id", &existingPlaceIDs).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch existing places")
		}

		// Create a map for O(1) lookups
		existingPlaceMap := make(map[string]bool)
		for _, id := range existingPlaceIDs {
			existingPlaceMap[id] = true
		}

		// Filter out places that already exist
		var newPlaces []models.GroupPlace
		for _, place := range req.Places {
			if existingPlaceMap[place.Id] {
				applogger.Info("Place", place.Id, "already exists for group", groupID, "- skipping")
				continue
			}

			// Add place to batch
			groupPlace := models.GroupPlace{
				ID:        uuid.New(),
				GroupId:   groupID,
				PlaceId:   place.Id,
				Name:      place.Name,
				Address:   place.Address,
				Type:      string(place.Type),
				Rating:    place.Rating,
				MapURI:    place.MapURI,
				Latitude:  place.Latitude,
				Longitude: place.Longitude,
			}
			newPlaces = append(newPlaces, groupPlace)
		}

		// Bulk insert if there are any new places
		if len(newPlaces) > 0 {
			if err := tx.Create(&newPlaces).Error; err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "Failed to add places to group")
			}

			// Prepare responses
			for _, groupPlace := range newPlaces {
				responses = append(responses, dto.GroupPlaceResponse{
					ID:        groupPlace.ID.String(),
					GroupID:   groupPlace.GroupId,
					PlaceID:   groupPlace.PlaceId,
					Name:      groupPlace.Name,
					Address:   groupPlace.Address,
					Type:      config.PlaceType(groupPlace.Type),
					Rating:    groupPlace.Rating,
					MapURI:    groupPlace.MapURI,
					Latitude:  groupPlace.Latitude,
					Longitude: groupPlace.Longitude,
				})
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return responses, nil
}

// RemoveAllPlacesFromGroup removes all places from a group
func (c *GroupPlacesController) RemoveAllPlacesFromGroup(groupID string) error {
	// Check if group exists
	var group models.Group
	if err := c.db.First(&group, "id = ?", groupID).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Group not found")
	}

	// Delete all places associated with the group
	if err := c.db.Where("group_id = ?", groupID).Delete(&models.GroupPlace{}).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to remove places from group")
	}

	return nil
}

// GetGroupPlaces retrieves all places for a group
func (c *GroupPlacesController) GetGroupPlaces(groupID string) ([]dto.GroupPlaceResponse, error) {
	// Check if group exists
	var group models.Group
	if err := c.db.First(&group, "id = ?", groupID).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Group not found")
	}

	var groupPlaces []models.GroupPlace
	if err := c.db.Where("group_id = ?", groupID).Find(&groupPlaces).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch group places")
	}

	responses := make([]dto.GroupPlaceResponse, len(groupPlaces))
	for i, place := range groupPlaces {
		responses[i] = dto.GroupPlaceResponse{
			ID:        place.ID.String(),
			GroupID:   place.GroupId,
			PlaceID:   place.PlaceId,
			Name:      place.Name,
			Address:   place.Address,
			Type:      config.PlaceType(place.Type),
			Rating:    place.Rating,
			MapURI:    place.MapURI,
			Latitude:  place.Latitude,
			Longitude: place.Longitude,
		}
	}

	return responses, nil
}
