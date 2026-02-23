package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestGlobalRateLimiter(t *testing.T) {
	t.Run("enabled (non-test env)", func(t *testing.T) {
		t.Setenv("ENV", "production")

		app := fiber.New(fiber.Config{ProxyHeader: fiber.HeaderXForwardedFor})
		app.Use(GlobalRateLimiter())
		app.Get("/", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

		// First 5 requests should succeed
		for range 5 {
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set(fiber.HeaderXForwardedFor, "10.0.0.1")
			resp := assertRequest(t, app, req)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		}

		// 6th request should be rate limited
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set(fiber.HeaderXForwardedFor, "10.0.0.1")
		resp := assertRequest(t, app, req)
		assert.Equal(t, fiber.StatusTooManyRequests, resp.StatusCode)
	})

	t.Run("disabled (test env)", func(t *testing.T) {
		t.Setenv("ENV", "test")

		app := fiber.New(fiber.Config{ProxyHeader: fiber.HeaderXForwardedFor})
		app.Use(GlobalRateLimiter())
		app.Get("/", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

		// All requests should succeed when rate limiting is disabled
		for range 10 {
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set(fiber.HeaderXForwardedFor, "10.0.0.1")
			resp := assertRequest(t, app, req)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		}
	})
}

func TestGlobalRateLimiterMinute(t *testing.T) {
	t.Run("enabled (non-test env)", func(t *testing.T) {
		t.Setenv("ENV", "production")

		app := fiber.New(fiber.Config{ProxyHeader: fiber.HeaderXForwardedFor})
		app.Use(GlobalRateLimiterMinute())
		app.Get("/", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

		// First 50 requests should succeed
		for range 50 {
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set(fiber.HeaderXForwardedFor, "10.0.0.2")
			resp := assertRequest(t, app, req)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		}

		// 51st request should be rate limited
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set(fiber.HeaderXForwardedFor, "10.0.0.2")
		resp := assertRequest(t, app, req)
		assert.Equal(t, fiber.StatusTooManyRequests, resp.StatusCode)
	})

	t.Run("disabled (test env)", func(t *testing.T) {
		t.Setenv("ENV", "test")

		app := fiber.New(fiber.Config{ProxyHeader: fiber.HeaderXForwardedFor})
		app.Use(GlobalRateLimiterMinute())
		app.Get("/", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

		// All requests should succeed when rate limiting is disabled
		for range 100 {
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set(fiber.HeaderXForwardedFor, "10.0.0.2")
			resp := assertRequest(t, app, req)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		}
	})
}

func TestUserCreateRateLimiter(t *testing.T) {
	t.Run("enabled (non-test env)", func(t *testing.T) {
		t.Setenv("ENV", "production")

		app := fiber.New(fiber.Config{ProxyHeader: fiber.HeaderXForwardedFor})
		app.Post("/users", UserCreateRateLimiter(), func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusCreated) })

		// First 2 requests should succeed
		for range 2 {
			req := httptest.NewRequest("POST", "/users", nil)
			req.Header.Set(fiber.HeaderXForwardedFor, "10.0.0.3")
			resp := assertRequest(t, app, req)
			assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		}

		// 3rd request should be rate limited
		req := httptest.NewRequest("POST", "/users", nil)
		req.Header.Set(fiber.HeaderXForwardedFor, "10.0.0.3")
		resp := assertRequest(t, app, req)
		assert.Equal(t, fiber.StatusTooManyRequests, resp.StatusCode)
	})

	t.Run("disabled (test env)", func(t *testing.T) {
		t.Setenv("ENV", "test")

		app := fiber.New(fiber.Config{ProxyHeader: fiber.HeaderXForwardedFor})
		app.Post("/users", UserCreateRateLimiter(), func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusCreated) })

		// All requests should succeed when rate limiting is disabled
		for range 10 {
			req := httptest.NewRequest("POST", "/users", nil)
			req.Header.Set(fiber.HeaderXForwardedFor, "10.0.0.3")
			resp := assertRequest(t, app, req)
			assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		}
	})
}

func TestGroupCreateRateLimiter(t *testing.T) {
	t.Run("enabled (non-test env)", func(t *testing.T) {
		t.Setenv("ENV", "production")

		app := fiber.New(fiber.Config{ProxyHeader: fiber.HeaderXForwardedFor})
		app.Post("/groups", GroupCreateRateLimiter(), func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusCreated) })

		// First 2 requests should succeed
		for range 2 {
			req := httptest.NewRequest("POST", "/groups", nil)
			req.Header.Set(fiber.HeaderXForwardedFor, "10.0.0.4")
			resp := assertRequest(t, app, req)
			assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		}

		// 3rd request should be rate limited
		req := httptest.NewRequest("POST", "/groups", nil)
		req.Header.Set(fiber.HeaderXForwardedFor, "10.0.0.4")
		resp := assertRequest(t, app, req)
		assert.Equal(t, fiber.StatusTooManyRequests, resp.StatusCode)
	})

	t.Run("disabled (test env)", func(t *testing.T) {
		t.Setenv("ENV", "test")

		app := fiber.New(fiber.Config{ProxyHeader: fiber.HeaderXForwardedFor})
		app.Post("/groups", GroupCreateRateLimiter(), func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusCreated) })

		// All requests should succeed when rate limiting is disabled
		for range 10 {
			req := httptest.NewRequest("POST", "/groups", nil)
			req.Header.Set(fiber.HeaderXForwardedFor, "10.0.0.4")
			resp := assertRequest(t, app, req)
			assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		}
	})
}

func assertRequest(t *testing.T, app *fiber.App, req *http.Request) *http.Response {
	t.Helper()
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	return resp
}
