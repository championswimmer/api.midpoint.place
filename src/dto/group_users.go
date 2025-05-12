package dto

// GroupUserJoinRequest represents the request to add a user to a group
type GroupUserJoinRequest struct {
	Location
}

// GroupUserResponse represents the response for group user operations
type GroupUserResponse struct {
	UserID    uint    `json:"user_id"`
	GroupID   string  `json:"group_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Role      string  `json:"role"`
}
