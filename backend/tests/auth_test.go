package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/auth"
	"github.com/THARUN-BART/LearnX/backend/internal/crypto"
	"github.com/THARUN-BART/LearnX/backend/internal/otp"
	"github.com/THARUN-BART/LearnX/backend/internal/profile"
	"github.com/THARUN-BART/LearnX/backend/internal/supabase"
)

type mockSupabaseClient struct {
	adminCreateUserFunc         func(ctx context.Context, req supabase.AdminCreateUserRequest) (*supabase.SupabaseUser, error)
	adminGetUserByEmailFunc     func(ctx context.Context, email string) (*supabase.SupabaseUser, error)
	adminUpdateUserPasswordFunc func(ctx context.Context, userID string, newPassword string) error
	signInWithPasswordFunc      func(ctx context.Context, email, password string) (*supabase.TokenResponse, error)
	signOutFunc                 func(ctx context.Context, accessToken string) error
	getUserFunc                 func(ctx context.Context, accessToken string) (*supabase.SupabaseUser, error)
	getOAuthURLFunc             func(provider, redirectTo string) string
}

func (m *mockSupabaseClient) AdminCreateUser(ctx context.Context, req supabase.AdminCreateUserRequest) (*supabase.SupabaseUser, error) {
	if m.adminCreateUserFunc != nil {
		return m.adminCreateUserFunc(ctx, req)
	}
	return &supabase.SupabaseUser{
		ID:    "user-uuid-1234",
		Email: req.Email,
	}, nil
}

func (m *mockSupabaseClient) AdminGetUserByEmail(ctx context.Context, email string) (*supabase.SupabaseUser, error) {
	if m.adminGetUserByEmailFunc != nil {
		return m.adminGetUserByEmailFunc(ctx, email)
	}
	return nil, supabase.ErrUserNotFound
}

func (m *mockSupabaseClient) AdminUpdateUserPassword(ctx context.Context, userID string, newPassword string) error {
	if m.adminUpdateUserPasswordFunc != nil {
		return m.adminUpdateUserPasswordFunc(ctx, userID, newPassword)
	}
	return nil
}

func (m *mockSupabaseClient) SignInWithPassword(ctx context.Context, email, password string) (*supabase.TokenResponse, error) {
	if m.signInWithPasswordFunc != nil {
		return m.signInWithPasswordFunc(ctx, email, password)
	}
	return &supabase.TokenResponse{
		AccessToken:  "mock-access-token",
		RefreshToken: "mock-refresh-token",
		TokenType:    "bearer",
		ExpiresIn:    3600,
		User: supabase.SupabaseUser{
			ID:    "user-uuid-1234",
			Email: email,
		},
	}, nil
}

func (m *mockSupabaseClient) SignOut(ctx context.Context, accessToken string) error {
	if m.signOutFunc != nil {
		return m.signOutFunc(ctx, accessToken)
	}
	return nil
}

func (m *mockSupabaseClient) GetUser(ctx context.Context, accessToken string) (*supabase.SupabaseUser, error) {
	if m.getUserFunc != nil {
		return m.getUserFunc(ctx, accessToken)
	}
	return &supabase.SupabaseUser{
		ID:    "user-uuid-1234",
		Email: "student@learnx.com",
	}, nil
}

func (m *mockSupabaseClient) GetOAuthURL(provider, redirectTo string) string {
	if m.getOAuthURLFunc != nil {
		return m.getOAuthURLFunc(provider, redirectTo)
	}
	return "https://example.supabase.co/auth/v1/authorize?provider=google"
}

type mockAuthOTPService struct {
	sendOTPFunc   func(ctx context.Context, email string) (*otp.SendOTPResponse, error)
	verifyOTPFunc func(ctx context.Context, email string, code string) (*otp.VerifyOTPResponse, error)
}

func (m *mockAuthOTPService) SendOTP(ctx context.Context, email string) (*otp.SendOTPResponse, error) {
	if m.sendOTPFunc != nil {
		return m.sendOTPFunc(ctx, email)
	}
	return &otp.SendOTPResponse{Success: true, Message: "OTP sent successfully"}, nil
}

func (m *mockAuthOTPService) VerifyOTP(ctx context.Context, email string, code string) (*otp.VerifyOTPResponse, error) {
	if m.verifyOTPFunc != nil {
		return m.verifyOTPFunc(ctx, email, code)
	}
	return &otp.VerifyOTPResponse{Success: true, Message: "OTP verified successfully"}, nil
}

func (m *mockAuthOTPService) GetRemainingCooldown(email string) time.Duration {
	return 0
}

func setupTestAuthService() (auth.Service, *auth.MemoryRepository, *mockAuthOTPService, *mockSupabaseClient) {
	repo := auth.NewMemoryRepository()
	otpMock := &mockAuthOTPService{}
	sbMock := &mockSupabaseClient{}
	encKey := crypto.DeriveKey("test-secret-key-32-bytes-long!")

	svc := auth.NewService(repo, otpMock, sbMock, encKey, 10*time.Minute, 15*time.Minute)
	return svc, repo, otpMock, sbMock
}

func TestAuth_Signup_SuccessWithConfirmPassword(t *testing.T) {
	svc, repo, otpMock, _ := setupTestAuthService()

	otpSent := false
	otpMock.sendOTPFunc = func(ctx context.Context, email string) (*otp.SendOTPResponse, error) {
		if email == "student@learnx.com" {
			otpSent = true
		}
		return &otp.SendOTPResponse{Success: true, Message: "OTP sent successfully"}, nil
	}

	resp, err := svc.Signup(context.Background(), auth.SignupRequest{
		FullName:        "Test Student",
		Email:           "student@learnx.com",
		Password:        "StrongPassword123!",
		ConfirmPassword: "StrongPassword123!",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Success {
		t.Fatalf("expected success")
	}
	if !otpSent {
		t.Fatalf("expected OTP service to be called")
	}

	pending, err := repo.GetPendingSignup(context.Background(), "student@learnx.com")
	if err != nil {
		t.Fatalf("pending signup not found in repo: %v", err)
	}
	if pending.PasswordEncrypted == "StrongPassword123!" {
		t.Fatalf("password must not be stored in plain text!")
	}
}

func TestAuth_Signup_PasswordMismatch(t *testing.T) {
	svc, _, _, _ := setupTestAuthService()

	_, err := svc.Signup(context.Background(), auth.SignupRequest{
		FullName:        "Test Student",
		Email:           "student@learnx.com",
		Password:        "StrongPassword123!",
		ConfirmPassword: "DifferentPassword123!",
	})
	if !errors.Is(err, auth.ErrPasswordMismatch) {
		t.Fatalf("expected ErrPasswordMismatch, got: %v", err)
	}
}

func TestAuth_Signup_DuplicateSignup(t *testing.T) {
	svc, repo, _, _ := setupTestAuthService()

	_ = repo.CreateProfile(context.Background(), &profile.Profile{
		ID:       "existing-uuid",
		Email:    "existing@learnx.com",
		FullName: "Existing Student",
	})

	_, err := svc.Signup(context.Background(), auth.SignupRequest{
		FullName:        "Duplicate Attempt",
		Email:           "existing@learnx.com",
		Password:        "StrongPassword123!",
		ConfirmPassword: "StrongPassword123!",
	})
	if !errors.Is(err, auth.ErrUserAlreadyExists) {
		t.Fatalf("expected ErrUserAlreadyExists, got: %v", err)
	}
}

func TestAuth_SignupAfterOTPVerification_Success(t *testing.T) {
	svc, repo, _, sbMock := setupTestAuthService()

	_, err := svc.Signup(context.Background(), auth.SignupRequest{
		FullName:        "Student Verified",
		Email:           "verified@learnx.com",
		Password:        "StrongPassword123!",
		ConfirmPassword: "StrongPassword123!",
	})
	if err != nil {
		t.Fatalf("signup initiation failed: %v", err)
	}

	var createdUserReq supabase.AdminCreateUserRequest
	sbMock.adminCreateUserFunc = func(ctx context.Context, req supabase.AdminCreateUserRequest) (*supabase.SupabaseUser, error) {
		createdUserReq = req
		return &supabase.SupabaseUser{
			ID:    "verified-uuid-456",
			Email: req.Email,
		}, nil
	}

	authResp, err := svc.VerifySignupOTP(context.Background(), auth.VerifySignupOTPRequest{
		Email: "verified@learnx.com",
		Code:  "123456",
	})
	if err != nil {
		t.Fatalf("verify signup otp failed: %v", err)
	}

	if !authResp.Success {
		t.Fatalf("expected authResp success = true")
	}
	if authResp.User.ID != "verified-uuid-456" {
		t.Errorf("unexpected user ID: %s", authResp.User.ID)
	}
	if !createdUserReq.EmailConfirm {
		t.Errorf("expected EmailConfirm=true on Supabase Admin call")
	}

	prof, err := repo.GetProfileByID(context.Background(), "verified-uuid-456")
	if err != nil {
		t.Fatalf("profile was not created: %v", err)
	}
	if prof.FullName != "Student Verified" {
		t.Errorf("expected profile name 'Student Verified', got '%s'", prof.FullName)
	}
}

func TestAuth_Login_Success(t *testing.T) {
	svc, repo, _, sbMock := setupTestAuthService()

	_ = repo.CreateProfile(context.Background(), &profile.Profile{
		ID:                  "user-100",
		Email:               "login@learnx.com",
		FullName:            "Learner One",
		OnboardingCompleted: true,
	})

	sbMock.signInWithPasswordFunc = func(ctx context.Context, email, password string) (*supabase.TokenResponse, error) {
		return &supabase.TokenResponse{
			AccessToken:  "valid-jwt-token",
			RefreshToken: "valid-refresh-token",
			TokenType:    "bearer",
			ExpiresIn:    3600,
			User: supabase.SupabaseUser{
				ID:    "user-100",
				Email: email,
			},
		}, nil
	}

	resp, err := svc.Login(context.Background(), auth.LoginRequest{
		Email:    "login@learnx.com",
		Password: "StrongPassword123!",
	})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if resp.User.Email != "login@learnx.com" || resp.User.FullName != "Learner One" {
		t.Errorf("unexpected user summary: %+v", resp.User)
	}
	if !resp.User.OnboardingCompleted {
		t.Errorf("expected OnboardingCompleted = true")
	}
}

func TestAuth_ForgotPassword_GenericResponse(t *testing.T) {
	svc, repo, otpMock, _ := setupTestAuthService()

	otpCalled := false
	otpMock.sendOTPFunc = func(ctx context.Context, email string) (*otp.SendOTPResponse, error) {
		if email == "registered@learnx.com" {
			otpCalled = true
		}
		return &otp.SendOTPResponse{Success: true, Message: "OTP sent"}, nil
	}

	_ = repo.CreateProfile(context.Background(), &profile.Profile{
		ID:       "registered-id",
		Email:    "registered@learnx.com",
		FullName: "Registered User",
	})

	resp1, err := svc.ForgotPassword(context.Background(), auth.ForgotPasswordRequest{
		Email: "registered@learnx.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !otpCalled {
		t.Fatalf("expected OTP to be dispatched")
	}

	resp2, err := svc.ForgotPassword(context.Background(), auth.ForgotPasswordRequest{
		Email: "unregistered@learnx.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp1.Message != resp2.Message {
		t.Fatalf("messages must match to prevent enumeration: %q vs %q", resp1.Message, resp2.Message)
	}
}

func TestAuth_PasswordReset_FlowAndTokenReuse(t *testing.T) {
	svc, repo, _, sbMock := setupTestAuthService()

	_ = repo.CreateProfile(context.Background(), &profile.Profile{
		ID:       "user-reset-id",
		Email:    "reset@learnx.com",
		FullName: "Reset User",
	})

	verifyResp, err := svc.VerifyResetOTP(context.Background(), auth.VerifyResetOTPRequest{
		Email: "reset@learnx.com",
		Code:  "123456",
	})
	if err != nil {
		t.Fatalf("verify reset OTP failed: %v", err)
	}

	// 1. Password mismatch during reset
	_, err = svc.ResetPassword(context.Background(), auth.ResetPasswordRequest{
		ResetToken:         verifyResp.ResetToken,
		NewPassword:        "BrandNewPassword999!",
		ConfirmNewPassword: "DifferentPassword123!",
	})
	if !errors.Is(err, auth.ErrPasswordMismatch) {
		t.Fatalf("expected ErrPasswordMismatch, got: %v", err)
	}

	// 2. Successful reset with matching passwords
	sbMock.adminUpdateUserPasswordFunc = func(ctx context.Context, userID, newPassword string) error {
		return nil
	}

	resetResp, err := svc.ResetPassword(context.Background(), auth.ResetPasswordRequest{
		ResetToken:         verifyResp.ResetToken,
		NewPassword:        "BrandNewPassword999!",
		ConfirmNewPassword: "BrandNewPassword999!",
	})
	if err != nil {
		t.Fatalf("reset password failed: %v", err)
	}
	if !resetResp.Success {
		t.Fatalf("expected reset success = true")
	}

	// 3. Reusing the same reset token must fail
	_, err = svc.ResetPassword(context.Background(), auth.ResetPasswordRequest{
		ResetToken:         verifyResp.ResetToken,
		NewPassword:        "AnotherPassword123!",
		ConfirmNewPassword: "AnotherPassword123!",
	})
	if !errors.Is(err, auth.ErrInvalidResetToken) {
		t.Fatalf("expected ErrInvalidResetToken when reusing reset token, got: %v", err)
	}
}
