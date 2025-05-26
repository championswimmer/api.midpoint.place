package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/championswimmer/api.midpoint.place/tests"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestGroupsRoute_MemberRejoin(t *testing.T) {
	// 1. Create two users user1 and user2
	user1 := tests.TestUtil_CreateUser(t, "testuser_rejoin1@test.com", "testpassword_rejoin1")
	user2 := tests.TestUtil_CreateUser(t, "testuser_rejoin2@test.com", "testpassword_rejoin2")

	// 2. Create a group by user1
	group := tests.TestUtil_CreateGroup(t, user1.Token, "Test Rejoin Group")

	joinRequestBody := []byte(`{
		"latitude": 12.971645,
		"longitude": 77.594562
	}`)

	// Helper function to get group details and check member count
	checkGroupResp := func(assertFunc func(t *testing.T, groupResp dto.GroupResponse)) {
		time.Sleep(20 * time.Millisecond) // Allow time for eventual consistency
		req := httptest.NewRequest("GET", "/v1/groups/"+group.ID+"?includeUsers=true", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+user1.Token) // User1 is the group owner

		resp := lo.Must(tests.App.Test(req, -1))
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var groupResp dto.GroupResponse
		body := lo.Must(io.ReadAll(resp.Body))
		err := json.Unmarshal(body, &groupResp)
		assert.NoError(t, err)

		assertFunc(t, groupResp)
	}

	// 3. user2 joins group > check member_count = 2 (owner is not counted in Members array)
	t.Run("user2 joins group", func(t *testing.T) {
		reqJoin := httptest.NewRequest("PUT", "/v1/groups/"+group.ID+"/join", bytes.NewBuffer(joinRequestBody))
		reqJoin.Header.Set("Content-Type", "application/json")
		reqJoin.Header.Set("Authorization", "Bearer "+user2.Token)

		respJoin := lo.Must(tests.App.Test(reqJoin, -1))
		assert.Equal(t, fiber.StatusAccepted, respJoin.StatusCode)
		checkGroupResp(func(t *testing.T, groupResp dto.GroupResponse) {
			assert.Len(t, groupResp.Members, 1)
		}) // user2 is now a member, owner is not in Members array
	})

	// 4. user2 leaves group > check member_count = 1
	t.Run("user2 leaves group", func(t *testing.T) {
		reqLeave := httptest.NewRequest("DELETE", "/v1/groups/"+group.ID+"/join", nil)
		reqLeave.Header.Set("Content-Type", "application/json")
		reqLeave.Header.Set("Authorization", "Bearer "+user2.Token)

		respLeave := lo.Must(tests.App.Test(reqLeave, -1))
		assert.Equal(t, fiber.StatusAccepted, respLeave.StatusCode)
		checkGroupResp(func(t *testing.T, groupResp dto.GroupResponse) {
			assert.Len(t, groupResp.Members, 0)
		}) // user2 left
	})

	// 5. user2 joins group again > check member_count = 2
	t.Run("user2 rejoins group", func(t *testing.T) {
		reqRejoin := httptest.NewRequest("PUT", "/v1/groups/"+group.ID+"/join", bytes.NewBuffer(joinRequestBody))
		reqRejoin.Header.Set("Content-Type", "application/json")
		reqRejoin.Header.Set("Authorization", "Bearer "+user2.Token)

		respRejoin := lo.Must(tests.App.Test(reqRejoin, -1))
		assert.Equal(t, fiber.StatusAccepted, respRejoin.StatusCode)
		checkGroupResp(func(t *testing.T, groupResp dto.GroupResponse) {
			assert.Len(t, groupResp.Members, 1)
		}) // user2 rejoined
	})
}
