package e2e

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/championswimmer/api.midpoint.place/tests"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestUsersRoute_RegisterUser(t *testing.T) {

	reqBody := []byte(`{"username": "testuser", "password": "testpassword"}`)

	req := httptest.NewRequest("POST", "/v1/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp := lo.Must(tests.App.Test(req, -1))

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}
