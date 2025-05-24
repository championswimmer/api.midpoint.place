package e2e

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"strconv"

	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/championswimmer/api.midpoint.place/tests"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestFetchUserGroups(t *testing.T) {
	// Create users user1 and user2
	user1 := tests.TestUtil_CreateUser(t, "testuser4001", "testpassword4001")
	user2 := tests.TestUtil_CreateUser(t, "testuser4002", "testpassword4002")

	// Create a group by user1
	group := tests.TestUtil_CreateGroup(t, user1.Token, "Test Group")

	// Fetch all groups of user1 via token of user2
	req := httptest.NewRequest("GET", "/v1/users/"+strconv.FormatUint(uint64(user1.ID), 10)+"/groups", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+user2.Token)

	resp := lo.Must(tests.App.Test(req, -1))
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Verify that the newly created group is there in the response
	var groups []dto.GroupResponse
	body := lo.Must(io.ReadAll(resp.Body))
	err := json.Unmarshal(body, &groups)
	assert.NoError(t, err)

	found := false
	for _, g := range groups {
		if g.ID == group.ID {
			found = true
			break
		}
	}
	assert.True(t, found, "Newly created group not found in the response")
}
