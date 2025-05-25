package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/championswimmer/api.midpoint.place/src/utils/applogger"
	"github.com/championswimmer/api.midpoint.place/tests"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestGroupsRoute_ListGroups(t *testing.T) {
	// 1. Create user testuser501
	createdUser := tests.TestUtil_CreateUser(t, "testuser5501@test.com", "testpassword5501")

	// 2. Make a new public group with that user's token
	groupName := "Test Public Group for Listing"
	createGroupReqBody := []byte(fmt.Sprintf(`{
		"name": "%s",
		"type": "public",
		"radius": 1500
	}`, groupName))

	createReq := httptest.NewRequest("POST", "/v1/groups", bytes.NewBuffer(createGroupReqBody))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+createdUser.Token)

	createResp := lo.Must(tests.App.Test(createReq, -1))
	assert.Equal(t, fiber.StatusCreated, createResp.StatusCode)

	var createdGroup dto.GroupResponse
	createBody := lo.Must(io.ReadAll(createResp.Body))
	err := json.Unmarshal(createBody, &createdGroup)
	assert.NoError(t, err)
	assert.Equal(t, groupName, createdGroup.Name)

	// 3. List all public groups and verify the group exists in the list
	t.Run("list all public groups", func(t *testing.T) {
		listReq := httptest.NewRequest("GET", "/v1/groups", nil)
		listReq.Header.Set("Authorization", "Bearer "+createdUser.Token)

		listResp := lo.Must(tests.App.Test(listReq, -1))
		assert.Equal(t, fiber.StatusOK, listResp.StatusCode)

		var groups []dto.GroupResponse
		body := lo.Must(io.ReadAll(listResp.Body))
		err := json.Unmarshal(body, &groups)
		assert.NoError(t, err)
		assert.Greater(t, len(groups), 0, "Expected at least one group in the list")
		applogger.Info("groups", groups)

		foundGroup := lo.ContainsBy(groups, func(g dto.GroupResponse) bool {
			return g.ID == createdGroup.ID
		})
		assert.True(t, foundGroup, "Created group not found in public list")
		// Also check member count is 0 for the new group
		for _, g := range groups {
			if g.ID == createdGroup.ID {
				assert.Equal(t, 0, g.MemberCount, "Newly created group should have 0 members")
			}
		}
	})

	// 4. List all groups where self=creator and verify the group exists in the list
	t.Run("list groups by creator", func(t *testing.T) {
		listReq := httptest.NewRequest("GET", "/v1/groups?self=creator", nil)
		listReq.Header.Set("Authorization", "Bearer "+createdUser.Token)

		listResp := lo.Must(tests.App.Test(listReq, -1))
		assert.Equal(t, fiber.StatusOK, listResp.StatusCode)

		var groups []dto.GroupResponse
		body := lo.Must(io.ReadAll(listResp.Body))
		err := json.Unmarshal(body, &groups)
		assert.NoError(t, err)

		foundGroup := lo.ContainsBy(groups, func(g dto.GroupResponse) bool {
			return g.ID == createdGroup.ID
		})
		assert.True(t, foundGroup, "Created group not found in self=creator list")
		// Also check member count is 0 for the new group
		for _, g := range groups {
			if g.ID == createdGroup.ID {
				assert.Equal(t, 0, g.MemberCount, "Newly created group should have 0 members")
			}
		}
	})
}
