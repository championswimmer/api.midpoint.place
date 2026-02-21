package validators

import (
	"testing"

	"github.com/championswimmer/api.midpoint.place/src/config"
	"github.com/championswimmer/api.midpoint.place/src/dto"
	"github.com/stretchr/testify/assert"
)

func TestValidateCreateGroupRequest_PlaceTypes(t *testing.T) {
	valid := ValidateCreateGroupRequest(&dto.CreateGroupRequest{
		Name:       "Test Group",
		PlaceTypes: []config.PlaceType{config.PlaceTypeCafe, config.PlaceTypeMuseum},
	})
	assert.Nil(t, valid)

	invalid := ValidateCreateGroupRequest(&dto.CreateGroupRequest{
		Name:       "Test Group",
		PlaceTypes: []config.PlaceType{"invalid_type"},
	})
	assert.NotNil(t, invalid)
}

func TestValidateUpdateGroupRequest_PlaceTypes(t *testing.T) {
	empty := []config.PlaceType{}
	invalid := []config.PlaceType{"invalid_type"}
	valid := []config.PlaceType{config.PlaceTypeBookstore}

	assert.NotNil(t, ValidateUpdateGroupRequest(&dto.UpdateGroupRequest{PlaceTypes: &empty}))
	assert.NotNil(t, ValidateUpdateGroupRequest(&dto.UpdateGroupRequest{PlaceTypes: &invalid}))
	assert.Nil(t, ValidateUpdateGroupRequest(&dto.UpdateGroupRequest{PlaceTypes: &valid}))
}
