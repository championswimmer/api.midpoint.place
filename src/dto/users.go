package dto

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserRequest = CreateUserRequest

type Location struct {
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

type UserUpdateRequest struct {
	Location Location `json:"location"`
}

type UserResponse struct {
	Id       string   `json:"id"`
	Username string   `json:"username"`
	Token    string   `json:"token"`
	Location Location `json:"location,omitempty"`
}
