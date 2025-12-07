package request

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=4,max=100"`
	Email    string `json:"email" validate:"required,email,min=6,max=255"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,min=6,max=255"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}
