package request

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=4,max=100"`
	Email    string `json:"email" validate:"required,email,max=200"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=200"`
	Password string `json:"password" validate:"required,min=8"`
}
