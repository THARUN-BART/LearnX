package auth

import (
	"errors"
	"time"
)

var (
	ErrInvalidRequest        = errors.New("invalid request parameters")
	ErrWeakPassword          = errors.New("password must be at least 8 characters and include uppercase, lowercase, and numbers")
	ErrPasswordMismatch      = errors.New("password and confirmation password do not match")
	ErrUserAlreadyExists     = errors.New("user already registered with this email")
	ErrPendingSignupNotFound = errors.New("no pending signup found or OTP has expired, please sign up again")
	ErrInvalidResetToken     = errors.New("invalid, expired, or already used password reset token")
	ErrUserNotFound          = errors.New("account not found")
	ErrUnauthorized          = errors.New("unauthorized")
)

type SignupRequest struct {
	FullName        string `json:"full_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password,omitempty"`
}

type VerifySignupOTPRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type VerifyResetOTPRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type ResetPasswordRequest struct {
	ResetToken         string `json:"reset_token"`
	NewPassword        string `json:"new_password"`
	ConfirmNewPassword string `json:"confirm_new_password,omitempty"`
}

type PendingSignup struct {
	ID                string
	Email             string
	FullName          string
	PasswordEncrypted string
	PasswordHash      string
	CreatedAt         time.Time
	ExpiresAt         time.Time
}

type PasswordReset struct {
	ID        string
	TokenHash string
	UserID    string
	Email     string
	CreatedAt time.Time
	ExpiresAt time.Time
	UsedAt    *time.Time
}

type UserSummary struct {
	ID                  string `json:"id"`
	Email               string `json:"email"`
	FullName            string `json:"full_name"`
	Role                string `json:"role"`
	OnboardingCompleted bool   `json:"onboarding_completed"`
}

type SessionInfo struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type SignupPendingResponse struct {
	Success          bool   `json:"success"`
	Message          string `json:"message"`
	Email            string `json:"email"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type AuthResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	User    *UserSummary `json:"user,omitempty"`
	Session *SessionInfo `json:"session,omitempty"`
}

type VerifyResetResponse struct {
	Success          bool   `json:"success"`
	Message          string `json:"message"`
	ResetToken       string `json:"reset_token"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type GenericResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
