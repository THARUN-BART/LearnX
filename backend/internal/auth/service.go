package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"
	"unicode"

	"github.com/THARUN-BART/LearnX/backend/internal/crypto"
	"github.com/THARUN-BART/LearnX/backend/internal/otp"
	"github.com/THARUN-BART/LearnX/backend/internal/profile"
	"github.com/THARUN-BART/LearnX/backend/internal/supabase"
)

type Service interface {
	Signup(ctx context.Context, req SignupRequest) (*SignupPendingResponse, error)
	VerifySignupOTP(ctx context.Context, req VerifySignupOTPRequest) (*AuthResponse, error)
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
	ForgotPassword(ctx context.Context, req ForgotPasswordRequest) (*GenericResponse, error)
	VerifyResetOTP(ctx context.Context, req VerifyResetOTPRequest) (*VerifyResetResponse, error)
	ResetPassword(ctx context.Context, req ResetPasswordRequest) (*GenericResponse, error)
	Logout(ctx context.Context, accessToken string) error
	GetOAuthURL(provider, redirectTo string) string
}

type authService struct {
	repo           Repository
	otpService     otp.Service
	supabaseClient supabase.Client
	encryptionKey  []byte
	signupExpiry   time.Duration
	resetExpiry    time.Duration
	logger         *slog.Logger
}

func NewService(
	repo Repository,
	otpService otp.Service,
	supabaseClient supabase.Client,
	encryptionKey []byte,
	signupExpiry time.Duration,
	resetExpiry time.Duration,
) Service {
	if signupExpiry <= 0 {
		signupExpiry = 10 * time.Minute
	}
	if resetExpiry <= 0 {
		resetExpiry = 15 * time.Minute
	}
	return &authService{
		repo:           repo,
		otpService:     otpService,
		supabaseClient: supabaseClient,
		encryptionKey:  encryptionKey,
		signupExpiry:   signupExpiry,
		resetExpiry:    resetExpiry,
		logger:         slog.Default(),
	}
}

// ValidatePassword ensures password meets strength policy: >=8 chars, uppercase, lowercase, digit.
func ValidatePassword(pass string) error {
	if len(pass) < 8 {
		return ErrWeakPassword
	}
	var hasUpper, hasLower, hasDigit bool
	for _, ch := range pass {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		}
	}
	if !hasUpper || !hasLower || !hasDigit {
		return ErrWeakPassword
	}
	return nil
}

// Signup handles Step 1 of the authentication flow.
// Validates Full Name, Email, Password, Confirm Password, and uniqueness.
func (s *authService) Signup(ctx context.Context, req SignupRequest) (*SignupPendingResponse, error) {
	fullName := strings.TrimSpace(req.FullName)
	if len(fullName) < 2 {
		return nil, fmt.Errorf("%w: full name must be at least 2 characters", ErrInvalidRequest)
	}

	cleanEmail, err := otp.CleanEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid email address", ErrInvalidRequest)
	}

	// Validate Confirm Password if supplied
	if req.ConfirmPassword != "" && req.Password != req.ConfirmPassword {
		return nil, ErrPasswordMismatch
	}

	if err := ValidatePassword(req.Password); err != nil {
		return nil, err
	}

	// Check if account already exists in database
	exists, err := s.repo.UserExists(ctx, cleanEmail)
	if err != nil {
		s.logger.Error("failed checking existing user", slog.String("error", err.Error()))
	}
	if exists {
		return nil, ErrUserAlreadyExists
	}

	// Encrypt sensitive password at rest using AES-256-GCM.
	encryptedPassword, err := crypto.Encrypt([]byte(req.Password), s.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to secure pending credentials: %w", err)
	}
	passwordHash := crypto.HashString(req.Password)

	now := time.Now()
	pending := &PendingSignup{
		Email:             cleanEmail,
		FullName:          fullName,
		PasswordEncrypted: encryptedPassword,
		PasswordHash:      passwordHash,
		CreatedAt:         now,
		ExpiresAt:         now.Add(s.signupExpiry),
	}

	if err := s.repo.SavePendingSignup(ctx, pending); err != nil {
		return nil, fmt.Errorf("failed to record pending registration: %w", err)
	}

	// Dispatch OTP using external OTP service
	_, err = s.otpService.SendOTP(ctx, cleanEmail)
	if err != nil {
		_ = s.repo.DeletePendingSignup(ctx, cleanEmail)
		return nil, err
	}

	return &SignupPendingResponse{
		Success:          true,
		Message:          "OTP sent to your email. Please verify to complete signup.",
		Email:            cleanEmail,
		ExpiresInSeconds: int(s.signupExpiry.Seconds()),
	}, nil
}

// VerifySignupOTP handles Step 2 of the authentication flow.
func (s *authService) VerifySignupOTP(ctx context.Context, req VerifySignupOTPRequest) (*AuthResponse, error) {
	cleanEmail, err := otp.CleanEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid email address", ErrInvalidRequest)
	}

	cleanCode, err := otp.CleanCode(req.Code)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid OTP code format", ErrInvalidRequest)
	}

	// 1. Verify OTP with external OTP service
	_, err = s.otpService.VerifyOTP(ctx, cleanEmail, cleanCode)
	if err != nil {
		return nil, err
	}

	// 2. Fetch pending signup record
	pending, err := s.repo.GetPendingSignup(ctx, cleanEmail)
	if err != nil {
		return nil, ErrPendingSignupNotFound
	}

	// 3. Decrypt password and verify integrity
	decryptedPasswordBytes, err := crypto.Decrypt(pending.PasswordEncrypted, s.encryptionKey)
	if err != nil {
		return nil, errors.New("failed to retrieve registration credentials")
	}
	plainPassword := string(decryptedPasswordBytes)

	if crypto.HashString(plainPassword) != pending.PasswordHash {
		return nil, errors.New("registration credentials integrity violation")
	}

	// 4. Create user in Supabase Auth using GoTrue Admin API
	adminReq := supabase.AdminCreateUserRequest{
		Email:        cleanEmail,
		Password:     plainPassword,
		EmailConfirm: true,
		UserMetadata: map[string]interface{}{
			"full_name": pending.FullName,
		},
	}

	sbUser, err := s.supabaseClient.AdminCreateUser(ctx, adminReq)
	if err != nil {
		if errors.Is(err, supabase.ErrUserAlreadyExists) {
			s.logger.Info("user already registered in Supabase Auth, syncing existing account credentials", slog.String("email", cleanEmail))
			existingUser, getErr := s.supabaseClient.AdminGetUserByEmail(ctx, cleanEmail)
			if getErr != nil || existingUser == nil {
				_ = s.repo.DeletePendingSignup(ctx, cleanEmail)
				return nil, ErrUserAlreadyExists
			}
			sbUser = existingUser
			_ = s.supabaseClient.AdminUpdateUserPassword(ctx, existingUser.ID, plainPassword)
		} else {
			s.logger.Error("supabase user creation failed", slog.String("error", err.Error()))
			return nil, fmt.Errorf("failed to create user account: %w", err)
		}
	}

	// 5. Create initial student profile
	userProfile := &profile.Profile{
		ID:                  sbUser.ID,
		Email:               cleanEmail,
		FullName:            pending.FullName,
		OnboardingCompleted: false,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
	if err := s.repo.CreateProfile(ctx, userProfile); err != nil {
		s.logger.Error("failed to create student profile in database", slog.String("error", err.Error()))
	}

	// 6. Invalidate pending signup
	_ = s.repo.DeletePendingSignup(ctx, cleanEmail)

	// 7. Obtain session tokens for authenticated user
	tokenResp, err := s.supabaseClient.SignInWithPassword(ctx, cleanEmail, plainPassword)
	if err != nil {
		s.logger.Warn("user created but initial session sign-in failed", slog.String("error", err.Error()))
		return &AuthResponse{
			Success: true,
			Message: "User account created and verified successfully. Please log in.",
			User: &UserSummary{
				ID:                  sbUser.ID,
				Email:               cleanEmail,
				FullName:            pending.FullName,
				Role:                "student",
				OnboardingCompleted: false,
			},
		}, nil
	}

	return &AuthResponse{
		Success: true,
		Message: "Signup and OTP verification completed successfully.",
		User: &UserSummary{
			ID:                  sbUser.ID,
			Email:               cleanEmail,
			FullName:            pending.FullName,
			Role:                "student",
			OnboardingCompleted: false,
		},
		Session: &SessionInfo{
			AccessToken:  tokenResp.AccessToken,
			RefreshToken: tokenResp.RefreshToken,
			TokenType:    tokenResp.TokenType,
			ExpiresIn:    tokenResp.ExpiresIn,
		},
	}, nil
}

// Login handles Step 3 of the authentication flow.
func (s *authService) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	cleanEmail, err := otp.CleanEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid email address", ErrInvalidRequest)
	}
	if strings.TrimSpace(req.Password) == "" {
		return nil, fmt.Errorf("%w: password cannot be blank", ErrInvalidRequest)
	}

	tokenResp, err := s.supabaseClient.SignInWithPassword(ctx, cleanEmail, req.Password)
	if err != nil {
		return nil, err
	}

	// Fetch user profile from database
	prof, err := s.repo.GetProfileByID(ctx, tokenResp.User.ID)
	fullName := cleanEmail
	onboardingCompleted := false
	if err == nil && prof != nil {
		fullName = prof.FullName
		onboardingCompleted = prof.OnboardingCompleted
	} else {
		if metaName, ok := tokenResp.User.UserMetadata["full_name"].(string); ok && metaName != "" {
			fullName = metaName
		}
		_ = s.repo.CreateProfile(ctx, &profile.Profile{
			ID:                  tokenResp.User.ID,
			Email:               cleanEmail,
			FullName:            fullName,
			OnboardingCompleted: false,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		})
	}

	return &AuthResponse{
		Success: true,
		Message: "Login successful.",
		User: &UserSummary{
			ID:                  tokenResp.User.ID,
			Email:               cleanEmail,
			FullName:            fullName,
			Role:                "student",
			OnboardingCompleted: onboardingCompleted,
		},
		Session: &SessionInfo{
			AccessToken:  tokenResp.AccessToken,
			RefreshToken: tokenResp.RefreshToken,
			TokenType:    tokenResp.TokenType,
			ExpiresIn:    tokenResp.ExpiresIn,
		},
	}, nil
}

// ForgotPassword handles Step 4 of the authentication flow.
func (s *authService) ForgotPassword(ctx context.Context, req ForgotPasswordRequest) (*GenericResponse, error) {
	cleanEmail, err := otp.CleanEmail(req.Email)
	if err != nil {
		return &GenericResponse{
			Success: true,
			Message: "If an account with that email exists, a password reset OTP has been sent.",
		}, nil
	}

	exists, err := s.repo.UserExists(ctx, cleanEmail)
	if err != nil {
		s.logger.Error("error checking user existence for password reset", slog.String("error", err.Error()))
	}

	if exists {
		_, err = s.otpService.SendOTP(ctx, cleanEmail)
		if err != nil {
			s.logger.Error("failed to dispatch password reset OTP", slog.String("error", err.Error()))
		}
	} else {
		s.logger.Info("password reset requested for non-existent email (enumeration guarded)", slog.String("email", cleanEmail))
	}

	return &GenericResponse{
		Success: true,
		Message: "If an account with that email exists, a password reset OTP has been sent.",
	}, nil
}

// VerifyResetOTP handles Step 5 of the authentication flow.
func (s *authService) VerifyResetOTP(ctx context.Context, req VerifyResetOTPRequest) (*VerifyResetResponse, error) {
	cleanEmail, err := otp.CleanEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid email address", ErrInvalidRequest)
	}

	cleanCode, err := otp.CleanCode(req.Code)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid OTP code format", ErrInvalidRequest)
	}

	_, err = s.otpService.VerifyOTP(ctx, cleanEmail, cleanCode)
	if err != nil {
		return nil, err
	}

	prof, err := s.repo.GetProfileByEmail(ctx, cleanEmail)
	if err != nil {
		return nil, ErrUserNotFound
	}

	rawToken, err := crypto.GenerateRandomToken(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate reset token: %w", err)
	}
	tokenHash := crypto.HashString(rawToken)

	now := time.Now()
	pr := &PasswordReset{
		TokenHash: tokenHash,
		UserID:    prof.ID,
		Email:     cleanEmail,
		CreatedAt: now,
		ExpiresAt: now.Add(s.resetExpiry),
	}

	if err := s.repo.SavePasswordReset(ctx, pr); err != nil {
		return nil, fmt.Errorf("failed to store reset authorization: %w", err)
	}

	return &VerifyResetResponse{
		Success:          true,
		Message:          "OTP verified successfully. You may now reset your password.",
		ResetToken:       rawToken,
		ExpiresInSeconds: int(s.resetExpiry.Seconds()),
	}, nil
}

// ResetPassword handles Step 6 of the authentication flow.
// Validates new_password and confirm_new_password.
func (s *authService) ResetPassword(ctx context.Context, req ResetPasswordRequest) (*GenericResponse, error) {
	rawToken := strings.TrimSpace(req.ResetToken)
	if rawToken == "" {
		return nil, ErrInvalidResetToken
	}

	// Validate Confirm New Password if supplied
	if req.ConfirmNewPassword != "" && req.NewPassword != req.ConfirmNewPassword {
		return nil, ErrPasswordMismatch
	}

	if err := ValidatePassword(req.NewPassword); err != nil {
		return nil, err
	}

	tokenHash := crypto.HashString(rawToken)

	pr, err := s.repo.GetPasswordResetByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, ErrInvalidResetToken
	}

	if pr.UsedAt != nil || time.Now().After(pr.ExpiresAt) {
		return nil, ErrInvalidResetToken
	}

	if err := s.repo.MarkPasswordResetUsed(ctx, tokenHash); err != nil {
		return nil, ErrInvalidResetToken
	}

	if err := s.supabaseClient.AdminUpdateUserPassword(ctx, pr.UserID, req.NewPassword); err != nil {
		s.logger.Error("failed to update user password in Supabase", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to update password: %w", err)
	}

	return &GenericResponse{
		Success: true,
		Message: "Password has been reset successfully. Please log in with your new password.",
	}, nil
}

func (s *authService) Logout(ctx context.Context, accessToken string) error {
	if accessToken == "" {
		return errors.New("missing access token")
	}
	return s.supabaseClient.SignOut(ctx, accessToken)
}

func (s *authService) GetOAuthURL(provider, redirectTo string) string {
	return s.supabaseClient.GetOAuthURL(provider, redirectTo)
}
