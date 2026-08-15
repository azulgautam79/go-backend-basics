package user

import "time"

type RegisterRequest struct {
	Name     string `json:"name" minLength:"2"`
	Email    string `json:"email" format:"email"`
	Password string `json:"password" minLength:"8"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User         User
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}
