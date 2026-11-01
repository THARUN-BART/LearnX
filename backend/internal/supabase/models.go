package supabase

import (
	"errors"
	"time"
)

var (
	ErrUserAlreadyExists  = errors.New("user already registered")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailNotConfirmed  = errors.New("email address has not been confirmed")
	ErrSupabaseUnavailable = errors.New("supabase auth service unavailable")
)

type SupabaseUser struct {
	ID           string                 `json:"id"`
	Email        string                 `json:"email"`
	UserMetadata map[string]interface{} `json:"user_metadata"`
	CreatedAt    time.Time              `json:"created_at"`
	ConfirmedAt  *time.Time             `json:"confirmed_at,omitempty"`
}

type AdminCreateUserRequest struct {
	Email        string                 `json:"email"`
	Password     string                 `json:"password"`
	EmailConfirm bool                   `json:"email_confirm"`
	UserMetadata map[string]interface{} `json:"user_metadata"`
}

type TokenResponse struct {
	AccessToken  string       `json:"access_token"`
	TokenType    string       `json:"token_type"`
	ExpiresIn    int          `json:"expires_in"`
	ExpiresAt    int64        `json:"expires_at"`
	RefreshToken string       `json:"refresh_token"`
	User         SupabaseUser `json:"user"`
}

type ErrorResponse struct {
	Code             int    `json:"code"`
	ErrorCode        string `json:"error_code"`
	Msg              string `json:"msg"`
	ErrorDescription string `json:"error_description"`
	Message          string `json:"message"`
}

func (e *ErrorResponse) GetMessage() string {
	if e.Msg != "" {
		return e.Msg
	}
	if e.Message != "" {
		return e.Message
	}
	if e.ErrorDescription != "" {
		return e.ErrorDescription
	}
	return "unknown supabase auth error"
}
