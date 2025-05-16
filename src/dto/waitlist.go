package dto

type WaitlistSignupRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type WaitlistSignupResponse struct {
	Message string `json:"message"`
}
