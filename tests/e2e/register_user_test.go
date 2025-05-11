package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/championswimmer/api.midpoint.place/tests"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestUsersRoute_RegisterUser(t *testing.T) {

	reqBody := []byte(`{"username": "testuser125", "password": "testpassword125"}`)

	req := httptest.NewRequest("POST", "/v1/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp := lo.Must(tests.App.Test(req, -1))

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var response dto.UserResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	assert.Equal(t, "testuser125", response.Username)
	assert.NotEmpty(t, response.Token)
	assert.NotEmpty(t, response.Id)
}
