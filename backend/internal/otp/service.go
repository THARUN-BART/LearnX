package otp

import (
	"context"
	"net/mail"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	otpCodeRegex = regexp.MustCompile(`^\d{4,8}$`)
)

type Service interface {
	SendOTP(ctx context.Context, email string) (*SendOTPResponse, error)
	VerifyOTP(ctx context.Context, email string, code string) (*VerifyOTPResponse, error)
	GetRemainingCooldown(email string) time.Duration
}

type otpState struct {
	lastSentAt   time.Time
	attemptsLeft int
}

type otpService struct {
	client       Client
	cooldown     time.Duration
	maxAttempts  int
	mu           sync.RWMutex
	stateByEmail map[string]*otpState
}

// NewService creates a new OTP domain service with rate limiting and cooldown management.
func NewService(client Client, cooldown time.Duration, maxAttempts int) Service {
	if cooldown <= 0 {
		cooldown = 60 * time.Second
	}
	if maxAttempts <= 0 {
		maxAttempts = 5
	}
	return &otpService{
		client:       client,
		cooldown:     cooldown,
		maxAttempts:  maxAttempts,
		stateByEmail: make(map[string]*otpState),
	}
}

// CleanEmail normalizes and validates email addresses.
func CleanEmail(email string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return "", ErrInvalidEmail
	}
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return "", ErrInvalidEmail
	}
	return email, nil
}

// CleanCode validates the format of an OTP code without exposing it.
func CleanCode(code string) (string, error) {
	code = strings.TrimSpace(code)
	if !otpCodeRegex.MatchString(code) {
		return "", ErrInvalidCode
	}
	return code, nil
}

// GetRemainingCooldown returns the remaining cooldown duration for an email.
func (s *otpService) GetRemainingCooldown(email string) time.Duration {
	clean, err := CleanEmail(email)
	if err != nil {
		return 0
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	state, exists := s.stateByEmail[clean]
	if !exists {
		return 0
	}

	elapsed := time.Since(state.lastSentAt)
	if elapsed < s.cooldown {
		return s.cooldown - elapsed
	}
	return 0
}

// SendOTP checks cooldown, resets attempts, and triggers external dispatch.
func (s *otpService) SendOTP(ctx context.Context, email string) (*SendOTPResponse, error) {
	cleanEmail, err := CleanEmail(email)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	state, exists := s.stateByEmail[cleanEmail]
	now := time.Now()
	if exists {
		elapsed := now.Sub(state.lastSentAt)
		if elapsed < s.cooldown {
			s.mu.Unlock()
			return nil, ErrCooldownActive
		}
	}

	// Update cooldown and reset verification attempts
	s.stateByEmail[cleanEmail] = &otpState{
		lastSentAt:   now,
		attemptsLeft: s.maxAttempts,
	}
	s.mu.Unlock()

	resp, err := s.client.SendOTP(ctx, cleanEmail)
	if err != nil {
		// If sending failed, release cooldown so user can retry
		s.mu.Lock()
		delete(s.stateByEmail, cleanEmail)
		s.mu.Unlock()
		return nil, err
	}

	return resp, nil
}

// VerifyOTP verifies the OTP code while guarding against brute-force attempts.
func (s *otpService) VerifyOTP(ctx context.Context, email string, code string) (*VerifyOTPResponse, error) {
	cleanEmail, err := CleanEmail(email)
	if err != nil {
		return nil, err
	}

	cleanCode, err := CleanCode(code)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	state, exists := s.stateByEmail[cleanEmail]
	if exists {
		if state.attemptsLeft <= 0 {
			s.mu.Unlock()
			return nil, ErrTooManyAttempts
		}
		state.attemptsLeft--
	}
	s.mu.Unlock()

	resp, err := s.client.VerifyOTP(ctx, cleanEmail, cleanCode)
	if err != nil {
		return nil, err
	}

	// On successful verification, clear state for this email
	s.mu.Lock()
	delete(s.stateByEmail, cleanEmail)
	s.mu.Unlock()

	return resp, nil
}
