package ratelimit

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// GlobalRateLimiter creates a new rate limiter for global requests.
func GlobalRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": "Too many requests (5 req/sec)",
			})
		},
	})
}

// GlobalRateLimiterMinute creates a new rate limiter for global requests per minute.
func GlobalRateLimiterMinute() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        50,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": "Too many requests (50 req/min)",
			})
		},
	})
}

// UserCreateRateLimiter creates a new rate limiter for user creation.
func UserCreateRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        2,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": "Too many user create requests (2 req/min)",
			})
		},
	})
}

// GroupCreateRateLimiter creates a new rate limiter for group creation.
func GroupCreateRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        2,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": "Too many group create requests (2 req/min)",
			})
		},
	})
}
