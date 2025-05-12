package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/championswimmer/api.midpoint.place/src/server"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

var App *fiber.App

func init() {
	App = server.CreateServer()
}

func TestUtil_CreateUser(t *testing.T, username string, password string) (userResponse *dto.UserResponse) {
	user := dto.CreateUserRequest{
		Username: username,
		Password: password,
	}

	body := lo.Must(json.Marshal(user))

	req := httptest.NewRequest("POST", "/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := lo.Must(App.Test(req, -1))

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var response dto.UserResponse
	body = lo.Must(io.ReadAll(resp.Body))
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err)

	assert.NotNil(t, response.Token)

	return &response
}
