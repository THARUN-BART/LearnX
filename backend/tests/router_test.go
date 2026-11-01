package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/auth"
	"github.com/THARUN-BART/LearnX/backend/internal/config"
	"github.com/THARUN-BART/LearnX/backend/internal/crypto"
	"github.com/THARUN-BART/LearnX/backend/internal/middleware"
	"github.com/THARUN-BART/LearnX/backend/internal/otp"
	"github.com/THARUN-BART/LearnX/backend/internal/profile"
	"github.com/THARUN-BART/LearnX/backend/internal/server"
	"github.com/THARUN-BART/LearnX/backend/internal/supabase"
)

type routerMockSupabase struct {
	users map[string]*supabase.SupabaseUser
}

func newRouterMockSupabase() *routerMockSupabase {
	return &routerMockSupabase{
		users: make(map[string]*supabase.SupabaseUser),
	}
}

func (m *routerMockSupabase) AdminCreateUser(ctx context.Context, req supabase.AdminCreateUserRequest) (*supabase.SupabaseUser, error) {
	if _, ok := m.users[req.Email]; ok {
		return nil, supabase.ErrUserAlreadyExists
	}
	user := &supabase.SupabaseUser{
		ID:           "test-user-id-" + req.Email,
		Email:        req.Email,
		UserMetadata: req.UserMetadata,
		CreatedAt:    time.Now(),
	}
	m.users[req.Email] = user
	return user, nil
}

func (m *routerMockSupabase) AdminGetUserByEmail(ctx context.Context, email string) (*supabase.SupabaseUser, error) {
	if u, ok := m.users[email]; ok {
		return u, nil
	}
	return nil, supabase.ErrUserNotFound
}

func (m *routerMockSupabase) AdminUpdateUserPassword(ctx context.Context, userID, newPassword string) error {
	return nil
}

func (m *routerMockSupabase) SignInWithPassword(ctx context.Context, email, password string) (*supabase.TokenResponse, error) {
	user, ok := m.users[email]
	if !ok {
		return nil, supabase.ErrInvalidCredentials
	}
	return &supabase.TokenResponse{
		AccessToken:  "mock-jwt-token-for-" + user.ID,
		RefreshToken: "mock-refresh-token",
		TokenType:    "bearer",
		ExpiresIn:    3600,
		User:         *user,
	}, nil
}

func (m *routerMockSupabase) SignOut(ctx context.Context, accessToken string) error {
	return nil
}

func (m *routerMockSupabase) GetUser(ctx context.Context, accessToken string) (*supabase.SupabaseUser, error) {
	if len(accessToken) > 19 && accessToken[:19] == "mock-jwt-token-for-" {
		userID := accessToken[19:]
		for _, u := range m.users {
			if u.ID == userID {
				return u, nil
			}
		}
	}
	return nil, supabase.ErrInvalidCredentials
}

func (m *routerMockSupabase) GetOAuthURL(provider, redirectTo string) string {
	return "https://example.supabase.co/auth/v1/authorize?provider=" + provider
}

type routerMockOTPClient struct{}

func (m *routerMockOTPClient) SendOTP(ctx context.Context, email string) (*otp.SendOTPResponse, error) {
	return &otp.SendOTPResponse{Success: true, Message: "OTP sent successfully"}, nil
}

func (m *routerMockOTPClient) VerifyOTP(ctx context.Context, email, code string) (*otp.VerifyOTPResponse, error) {
	if code == "123456" {
		return &otp.VerifyOTPResponse{Success: true, Message: "OTP verified successfully"}, nil
	}
	return nil, otp.ErrOTPVerificationFailed
}

func setupFullServer() (http.Handler, *auth.MemoryRepository, *profile.MemoryRepository) {
	cfg := &config.Config{
		OTPAppName:               "LearnX",
		PendingSignupExpiry:      10 * time.Minute,
		PasswordResetTokenExpiry: 15 * time.Minute,
		EncryptionKey:            crypto.DeriveKey("test-aes-key-for-router-test!!"),
		ClientURL:                "*",
	}

	authRepo := auth.NewMemoryRepository()
	profileRepo := profile.NewMemoryRepository()
	sbClient := newRouterMockSupabase()

	otpService := otp.NewService(&routerMockOTPClient{}, 1*time.Second, 5)
	authService := auth.NewService(authRepo, otpService, sbClient, cfg.EncryptionKey, cfg.PendingSignupExpiry, cfg.PasswordResetTokenExpiry)
	profileService := profile.NewService(profileRepo)

	authHandler := auth.NewHandler(authService)
	profileHandler := profile.NewHandler(profileService)
	authMiddleware := middleware.NewAuthMiddleware(sbClient, "")
	rateLimiter := middleware.NewRateLimiter(100, 1*time.Minute)

	router := server.BuildRouter(cfg, authHandler, profileHandler, authMiddleware, rateLimiter)
	return router, authRepo, profileRepo
}

func TestRouter_CompleteFlowWithOnboarding(t *testing.T) {
	router, _, profileRepo := setupFullServer()

	// 1. Health Checks
	for _, path := range []string{"/health", "/api/v1/health"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("path %s expected 200, got %d", path, rec.Code)
		}
	}

	// 2. Signup with Confirm Password mismatch
	mismatchBody, _ := json.Marshal(auth.SignupRequest{
		FullName:        "Bob Student",
		Email:           "bob@learnx.com",
		Password:        "Password123!",
		ConfirmPassword: "WrongPassword123!",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", bytes.NewBuffer(mismatchBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for password mismatch, got %d", rec.Code)
	}

	// 3. Step 1: Valid Signup
	signupBody, _ := json.Marshal(auth.SignupRequest{
		FullName:        "Bob Student",
		Email:           "bob@learnx.com",
		Password:        "Password123!",
		ConfirmPassword: "Password123!",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/signup", bytes.NewBuffer(signupBody))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("signup failed with status %d: %s", rec.Code, rec.Body.String())
	}

	// 4. Step 2: Verify Signup OTP (6-digit verification)
	verifyBody, _ := json.Marshal(auth.VerifySignupOTPRequest{
		Email: "bob@learnx.com",
		Code:  "123456",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/verify-signup-otp", bytes.NewBuffer(verifyBody))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("verify signup OTP failed with status %d: %s", rec.Code, rec.Body.String())
	}

	var authResp auth.AuthResponse
	_ = json.NewDecoder(rec.Body).Decode(&authResp)
	token := authResp.Session.AccessToken
	userID := authResp.User.ID

	// Seed profile in memory store for profile queries
	profileRepo.Save(&profile.Profile{
		ID:                  userID,
		FullName:            "Bob Student",
		Email:               "bob@learnx.com",
		OnboardingCompleted: false,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	})

	// 5. Student Onboarding: POST /api/v1/profile/onboarding
	onboardingBody, _ := json.Marshal(profile.OnboardingRequest{
		College:           "Stanford University",
		Course:            "B.Tech",
		Department:        "Computer Science",
		CurrentYear:       "2nd Year",
		PreferredLanguage: "Go",
		LearningGoal:      "DSA",
		DSALevel:          "Intermediate",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/profile/onboarding", bytes.NewBuffer(onboardingBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("onboarding failed with status %d: %s", rec.Code, rec.Body.String())
	}

	// 6. Verify Full Student Profile fields (Supabase) via GET /api/v1/profile
	req = httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("profile fetch failed with status %d: %s", rec.Code, rec.Body.String())
	}

	var profileResp struct {
		Success bool            `json:"success"`
		Profile profile.Profile `json:"profile"`
	}
	_ = json.NewDecoder(rec.Body).Decode(&profileResp)
	p := profileResp.Profile

	// Verify all 14 fields required by the user:
	// id, full_name, email, college, course, department, current_year,
	// preferred_language, learning_goal, dsa_level, avatar_url, onboarding_completed, created_at, updated_at
	if p.ID == "" || p.FullName != "Bob Student" || p.Email != "bob@learnx.com" {
		t.Errorf("basic profile fields mismatch: %+v", p)
	}
	if p.College != "Stanford University" || p.Course != "B.Tech" || p.Department != "Computer Science" {
		t.Errorf("academic fields mismatch: %+v", p)
	}
	if p.CurrentYear != "2nd Year" || p.PreferredLanguage != "Go" || p.LearningGoal != "DSA" || p.DSALevel != "Intermediate" {
		t.Errorf("onboarding fields mismatch: %+v", p)
	}
	if !p.OnboardingCompleted {
		t.Errorf("expected onboarding_completed = true")
	}
	if p.CreatedAt.IsZero() || p.UpdatedAt.IsZero() {
		t.Errorf("timestamps should not be zero")
	}

	// 7. Step 4: Forgot Password
	forgotBody, _ := json.Marshal(auth.ForgotPasswordRequest{
		Email: "bob@learnx.com",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/forgot-password", bytes.NewBuffer(forgotBody))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("forgot password failed: %s", rec.Body.String())
	}

	// 8. Step 5: Verify Reset OTP (6-digit verification)
	verifyResetBody, _ := json.Marshal(auth.VerifyResetOTPRequest{
		Email: "bob@learnx.com",
		Code:  "123456",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/verify-reset-otp", bytes.NewBuffer(verifyResetBody))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("verify reset OTP failed: %s", rec.Body.String())
	}

	var resetOtpResp auth.VerifyResetResponse
	_ = json.NewDecoder(rec.Body).Decode(&resetOtpResp)
	resetToken := resetOtpResp.ResetToken

	// 9. Step 6: Reset Password with Confirm New Password
	resetPassBody, _ := json.Marshal(auth.ResetPasswordRequest{
		ResetToken:         resetToken,
		NewPassword:        "BrandNewPassword123!",
		ConfirmNewPassword: "BrandNewPassword123!",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/reset-password", bytes.NewBuffer(resetPassBody))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("reset password failed: %s", rec.Body.String())
	}

	// 10. Replay attack prevention on reset token
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/reset-password", bytes.NewBuffer(resetPassBody))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request for reused reset token, got %d", rec.Code)
	}
}
