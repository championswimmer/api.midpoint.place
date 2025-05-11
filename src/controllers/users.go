package controllers

import (
	"strconv"

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
		Id:       strconv.FormatUint(uint64(user.ID), 10),
		Username: user.Username,
		Token:    "", // Will be implemented with JWT
	}, nil
}

func (c *UsersController) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := c.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	return &user, nil
}

func (c *UsersController) LoginUser(req *dto.LoginUserRequest) (*dto.UserResponse, error) {
	var user models.User
	if err := c.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	if !security.CheckPasswordHash(user.Password, req.Password) {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid password")
	}

	token := security.CreateJWTFromUser(&user)

	return &dto.UserResponse{
		Id:       strconv.FormatUint(uint64(user.ID), 10),
		Username: user.Username,
		Token:    token,
	}, nil
}
