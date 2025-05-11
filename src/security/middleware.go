package security

import (
	"github.com/gofiber/fiber/v2"
)

// MandatoryJwtAuthMiddleware makes authentication mandatory
// will return 401 if no Authorization header is provided or if the JWT is invalid
// saves the user in the context locals as "user"
func MandatoryJwtAuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized: No Authorization header provided",
		})
	}
	// Splice out the "Bearer " prefix, if it exists
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		authHeader = authHeader[7:]
	}
	user, err := ValidateAndParseJWT(authHeader)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized: Invalid JWT token" + err.Error(),
		})
	}
	c.Locals("user", user)
	return c.Next()
}

// OptionalJwtAuthMiddleware makes authentication optional
// will run only if there is an Authorization header present
// will return 401 if the JWT is invalid inside the header
// will save the user in the context locals as "user" if the JWT is valid
func OptionalJwtAuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Next()
	}
	user, err := ValidateAndParseJWT(authHeader)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized: Invalid JWT token",
		})
	}
	c.Locals("user", user)
	return c.Next()
}
