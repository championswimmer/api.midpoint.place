package dto

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserRequest = CreateUserRequest

type UserResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}
