package dto

type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AuthResponse struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token,omitempty"`
	RefreshToken string        `json:"refresh_token,omitempty"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
