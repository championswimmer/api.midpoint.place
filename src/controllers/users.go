package controllers

import (
	"github.com/championswimmer/api.midpoint.place/src/db"
	"github.com/championswimmer/api.midpoint.place/src/db/models"
	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/championswimmer/api.midpoint.place/src/security"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UsersController struct {
	db *gorm.DB
}

func CreateUsersController() *UsersController {
	appDb := db.GetAppDB()
	return &UsersController{
		db: appDb,
	}
}

func (c *UsersController) CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Check if username already exists
	var existingUser models.User
	if err := c.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, fiber.NewError(fiber.StatusConflict, "Username already exists")
	}

	// Create new user
	user := models.User{
		Username: req.Username,
		Password: security.HashPassword(req.Password),
	}

	if err := c.db.Create(&user).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to create user")
	}

	// TODO: Generate JWT token here when implementing authentication
	return &dto.UserResponse{
		Id:       string(user.ID),
		Username: user.Username,
		Token:    "", // Will be implemented with JWT
	}, nil
}
