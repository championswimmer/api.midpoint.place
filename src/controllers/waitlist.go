package controllers

import (
	"regexp"

	"github.com/championswimmer/api.midpoint.place/src/db"
	"github.com/championswimmer/api.midpoint.place/src/db/models"
	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type WaitlistController struct {
	db *gorm.DB
}

func CreateWaitlistController() *WaitlistController {
	appDb := db.GetAppDB()
	return &WaitlistController{
		db: appDb,
	}
}

func (c *WaitlistController) AddToWaitlist(req *dto.WaitlistSignupRequest) (*dto.WaitlistSignupResponse, error) {

	// Validate email format
	if !isValidEmail(req.Email) {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid email format")
	}

	// Check for duplicates
	var existingSignup models.WaitlistSignup
	if err := c.db.Where("email = ?", req.Email).First(&existingSignup).Error; err == nil {
		return nil, fiber.NewError(fiber.StatusConflict, "Email already signed up")
	}

	// Add to waitlist
	waitlistSignup := models.WaitlistSignup{
		Email: req.Email,
	}
	if err := c.db.Create(&waitlistSignup).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to add to waitlist")
	}

	return &dto.WaitlistSignupResponse{
		Message: "Successfully added to waitlist",
	}, nil
}

func isValidEmail(email string) bool {
	// Simple regex for email validation
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
