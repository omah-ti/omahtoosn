package auth

import "time"

// request
type RegisterRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	FullName string `json:"full_name" validate:"required,min=2,max=120"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// response
type UserResponse struct {
	ID          string     `json:"id"`
	Email       string     `json:"email"`
	FullName    string     `json:"full_name"`
	IsActive    bool       `json:"is_active"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type LoginResponse struct {
	AccessExpiresAt  time.Time    `json:"access_expires_at"`
	CookieStrategy   string       `json:"cookie_strategy"`
	RefreshExpiresAt time.Time    `json:"refresh_expires_at"`
	User             UserResponse `json:"user"`
}

type ForgotPasswordResponse struct{}
