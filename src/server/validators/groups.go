package validators

import (
	"regexp"

	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/gofiber/fiber/v2"
)

func validateName(name string) *ValidationError {
	if len(name) > 100 || len(name) < 1 {
		return &ValidationError{
			status:  fiber.StatusUnprocessableEntity,
			message: "Name must be between 1 and 100 characters",
		}
	}
	return nil
}

func validateSecret(secret string) *ValidationError {
	if secret == "" {
		return nil
	}
	if matched, _ := regexp.MatchString(`^\d{6}$`, secret); !matched {
		return &ValidationError{
			status:  fiber.StatusUnprocessableEntity,
			message: "Secret must be a 6-digit string or empty",
		}
	}
	return nil
}

func validateRadius(radius int) *ValidationError {
	if radius < 0 {
		return &ValidationError{
			status:  fiber.StatusUnprocessableEntity,
			message: "Radius must be a positive integer",
		}
	}
	return nil
}

func ValidateCreateGroupRequest(req *dto.CreateGroupRequest) *ValidationError {
	if err := validateName(req.Name); err != nil {
		return err
	}
	if err := validateSecret(req.Secret); err != nil {
		return err
	}
	if err := validateRadius(req.Radius); err != nil {
		return err
	}
	// Add any other specific validations for CreateGroupRequest
	return nil
}

func ValidateUpdateGroupRequest(req *dto.UpdateGroupRequest) *ValidationError {
	if req.Name != "" { // Only validate if provided for update
		if err := validateName(req.Name); err != nil {
			return err
		}
	}
	if req.Secret != "" { // Only validate if provided for update, or if explicitly set to be cleared (empty string)
		if err := validateSecret(req.Secret); err != nil {
			return err
		}
	}
	if req.Radius > 0 { // Only validate if provided for update and positive
		if err := validateRadius(req.Radius); err != nil {
			return err
		}
	}
	// Add any other specific validations for UpdateGroupRequest
	return nil
}
