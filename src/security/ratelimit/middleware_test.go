package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/championswimmer/api.midpoint.place/src/config"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func withNonTestEnv(t *testing.T, fn func()) {
	t.Helper()
	origEnv := config.Env
	config.Env = "local"
	defer func() { config.Env = origEnv }()
	fn()
}

func TestGlobalRateLimiter(t *testing.T) {
	withNonTestEnv(t, func() {
		app := fiber.New(fiber.Config{ProxyHeader: fiber.HeaderXForwardedFor})
		app.Use(GlobalRateLimiter())
		app.Get("/", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

		for range 5 {
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set(fiber.HeaderXForwardedFor, "10.0.0.1")
			resp := assertRequest(t, app, req)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		}

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set(fiber.HeaderXForwardedFor, "10.0.0.1")
		resp := assertRequest(t, app, req)
		assert.Equal(t, fiber.StatusTooManyRequests, resp.StatusCode)
	})
}

func TestGlobalRateLimiterMinute(t *testing.T) {
	withNonTestEnv(t, func() {
		app := fiber.New(fiber.Config{ProxyHeader: fiber.HeaderXForwardedFor})
		app.Use(GlobalRateLimiterMinute())
		app.Get("/", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

		for range 50 {
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set(fiber.HeaderXForwardedFor, "10.0.0.2")
			resp := assertRequest(t, app, req)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		}

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set(fiber.HeaderXForwardedFor, "10.0.0.2")
		resp := assertRequest(t, app, req)
		assert.Equal(t, fiber.StatusTooManyRequests, resp.StatusCode)
	})
}

func TestUserCreateRateLimiter(t *testing.T) {
	withNonTestEnv(t, func() {
		app := fiber.New(fiber.Config{ProxyHeader: fiber.HeaderXForwardedFor})
		app.Post("/users", UserCreateRateLimiter(), func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusCreated) })

		for range 2 {
			req := httptest.NewRequest("POST", "/users", nil)
			req.Header.Set(fiber.HeaderXForwardedFor, "10.0.0.3")
			resp := assertRequest(t, app, req)
			assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		}

		req := httptest.NewRequest("POST", "/users", nil)
		req.Header.Set(fiber.HeaderXForwardedFor, "10.0.0.3")
		resp := assertRequest(t, app, req)
		assert.Equal(t, fiber.StatusTooManyRequests, resp.StatusCode)
	})
}

func TestGroupCreateRateLimiter(t *testing.T) {
	withNonTestEnv(t, func() {
		app := fiber.New(fiber.Config{ProxyHeader: fiber.HeaderXForwardedFor})
		app.Post("/groups", GroupCreateRateLimiter(), func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusCreated) })

		for range 2 {
			req := httptest.NewRequest("POST", "/groups", nil)
			req.Header.Set(fiber.HeaderXForwardedFor, "10.0.0.4")
			resp := assertRequest(t, app, req)
			assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		}

		req := httptest.NewRequest("POST", "/groups", nil)
		req.Header.Set(fiber.HeaderXForwardedFor, "10.0.0.4")
		resp := assertRequest(t, app, req)
		assert.Equal(t, fiber.StatusTooManyRequests, resp.StatusCode)
	})
}

func assertRequest(t *testing.T, app *fiber.App, req *http.Request) *http.Response {
	t.Helper()
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	return resp
}
