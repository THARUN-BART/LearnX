package otp

import "errors"

var (
	ErrInvalidEmail          = errors.New("invalid email address")
	ErrInvalidCode           = errors.New("invalid or malformed OTP code")
	ErrCooldownActive        = errors.New("OTP resend cooldown active, please wait before requesting again")
	ErrTooManyAttempts       = errors.New("maximum verification attempts exceeded, please request a new OTP")
	ErrOTPVerificationFailed = errors.New("invalid or expired OTP code")
	ErrOTPServiceUnavailable = errors.New("external OTP service is currently unavailable")
	ErrUnexpectedResponse    = errors.New("unexpected response received from OTP service")
)

// SendOTPRequest represents the JSON payload to send an OTP via external service.
type SendOTPRequest struct {
	Email   string `json:"email"`
	AppName string `json:"app_name"`
}

// SendOTPResponse represents the JSON response from the external OTP service.
type SendOTPResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	AppName string `json:"app_name,omitempty"`
	Error   string `json:"error,omitempty"`
}

// VerifyOTPRequest represents the JSON payload to verify an OTP.
type VerifyOTPRequest struct {
	Email   string `json:"email"`
	Code    string `json:"code"`
	AppName string `json:"app_name"`
}

// VerifyOTPResponse represents the verification response from the external service.
type VerifyOTPResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	AppName string `json:"app_name,omitempty"`
	Error   string `json:"error,omitempty"`
}
