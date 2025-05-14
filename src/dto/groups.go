package dto

import "github.com/championswimmer/api.midpoint.place/src/config"

type CreateGroupRequest struct {
	Name   string           `json:"name" validate:"required"`
	Type   config.GroupType `json:"type" validate:"omitempty,oneof=public protected private"`
	Secret string           `json:"secret" validate:"omitempty"`
	Radius int              `json:"radius" validate:"omitempty,min=0"`
}

type UpdateGroupRequest struct {
	Name   string           `json:"name" validate:"omitempty"`
	Type   config.GroupType `json:"type" validate:"omitempty,oneof=public protected private"`
	Secret string           `json:"secret" validate:"omitempty"`
	Radius int              `json:"radius" validate:"omitempty,min=0"`
}

type UpdateGroupMidpointRequest struct {
	Location
}

type GroupResponse struct {
	ID                string           `json:"id"`
	Name              string           `json:"name"`
	Type              config.GroupType `json:"type"`
	Code              string           `json:"code"`
	CreatorID         string           `json:"creator_id"`
	MidpointLatitude  float64          `json:"midpoint_latitude"`
	MidpointLongitude float64          `json:"midpoint_longitude"`
	Radius            int              `json:"radius"`
}
