package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/THARUN-BART/LearnX/backend/internal/otp"
	"github.com/THARUN-BART/LearnX/backend/internal/supabase"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Signup handles POST /api/v1/auth/signup.
func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	res, err := h.service.Signup(r.Context(), req)
	if err != nil {
		mapAuthError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, res)
}

// VerifySignupOTP handles POST /api/v1/auth/verify-signup-otp.
func (h *Handler) VerifySignupOTP(w http.ResponseWriter, r *http.Request) {
	var req VerifySignupOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	res, err := h.service.VerifySignupOTP(r.Context(), req)
	if err != nil {
		mapAuthError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, res)
}

// Login handles POST /api/v1/auth/login.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	res, err := h.service.Login(r.Context(), req)
	if err != nil {
		mapAuthError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, res)
}

// ForgotPassword handles POST /api/v1/auth/forgot-password.
func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	res, err := h.service.ForgotPassword(r.Context(), req)
	if err != nil {
		mapAuthError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, res)
}

// VerifyResetOTP handles POST /api/v1/auth/verify-reset-otp.
func (h *Handler) VerifyResetOTP(w http.ResponseWriter, r *http.Request) {
	var req VerifyResetOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	res, err := h.service.VerifyResetOTP(r.Context(), req)
	if err != nil {
		mapAuthError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, res)
}

// ResetPassword handles POST /api/v1/auth/reset-password.
func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	res, err := h.service.ResetPassword(r.Context(), req)
	if err != nil {
		mapAuthError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, res)
}

// Logout handles POST /api/v1/auth/logout.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		respondError(w, http.StatusUnauthorized, "missing or invalid authorization header")
		return
	}

	token := strings.TrimSpace(parts[1])
	if err := h.service.Logout(r.Context(), token); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to logout")
		return
	}

	respondJSON(w, http.StatusOK, GenericResponse{
		Success: true,
		Message: "Logged out successfully.",
	})
}

// GoogleOAuth handles GET /api/v1/auth/google.
func (h *Handler) GoogleOAuth(w http.ResponseWriter, r *http.Request) {
	redirectTo := r.URL.Query().Get("redirect_to")
	authURL := h.service.GetOAuthURL("google", redirectTo)

	if r.Header.Get("Accept") == "application/json" || r.URL.Query().Get("format") == "json" {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"success":  true,
			"auth_url": authURL,
		})
		return
	}

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// OAuthCallback handles GET /api/v1/auth/callback.
func (h *Handler) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	errParam := r.URL.Query().Get("error_description")

	if errParam != "" {
		respondError(w, http.StatusBadRequest, errParam)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "OAuth exchange pending or completed",
		"code":    code,
	})
}

func mapAuthError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrUserAlreadyExists), errors.Is(err, supabase.ErrUserAlreadyExists):
		respondError(w, http.StatusConflict, "an account with this email already exists")
	case errors.Is(err, ErrPasswordMismatch):
		respondError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrWeakPassword):
		respondError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrInvalidRequest):
		respondError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, otp.ErrCooldownActive):
		respondError(w, http.StatusTooManyRequests, "OTP resend cooldown is active, please wait before requesting again")
	case errors.Is(err, otp.ErrTooManyAttempts):
		respondError(w, http.StatusTooManyRequests, "too many failed verification attempts, please request a new OTP")
	case errors.Is(err, otp.ErrOTPVerificationFailed):
		msg := err.Error()
		if idx := strings.LastIndex(msg, ": "); idx != -1 {
			msg = msg[idx+2:]
		}
		respondError(w, http.StatusBadRequest, msg)
	case errors.Is(err, ErrPendingSignupNotFound):
		respondError(w, http.StatusBadRequest, "no pending signup found or OTP has expired, please sign up again")
	case errors.Is(err, ErrInvalidResetToken):
		respondError(w, http.StatusBadRequest, "invalid, expired, or already used password reset token")
	case errors.Is(err, supabase.ErrInvalidCredentials):
		respondError(w, http.StatusUnauthorized, "invalid email or password")
	case errors.Is(err, supabase.ErrEmailNotConfirmed):
		respondError(w, http.StatusForbidden, "email address is not confirmed")
	case errors.Is(err, otp.ErrOTPServiceUnavailable), errors.Is(err, supabase.ErrSupabaseUnavailable):
		respondError(w, http.StatusServiceUnavailable, "service temporarily unavailable, please try again shortly")
	default:
		respondError(w, http.StatusBadRequest, err.Error())
	}
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]interface{}{
		"success": false,
		"error":   message,
	})
}
