package tests

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/otp"
)

func TestOTPClient_SendOTP_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/send-otp" {
			t.Errorf("expected path /api/send-otp, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var req otp.SendOTPRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.Email != "user@example.com" || req.AppName != "LearnX" {
			t.Errorf("unexpected request data: %+v", req)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(otp.SendOTPResponse{
			Success: true,
			Message: "OTP sent successfully",
			AppName: "LearnX",
		})
	}))
	defer server.Close()

	client, err := otp.NewClient(server.URL, "LearnX", 2*time.Second)
	if err != nil {
		t.Fatalf("unexpected NewClient error: %v", err)
	}

	resp, err := client.SendOTP(context.Background(), "user@example.com")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !resp.Success {
		t.Fatalf("expected success = true, got false")
	}
	if resp.Message != "OTP sent successfully" {
		t.Errorf("unexpected message: %s", resp.Message)
	}
}

func TestOTPClient_SendOTP_ServiceUnavailable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"success":false,"error":"Internal service crash"}`))
	}))
	defer server.Close()

	client, err := otp.NewClient(server.URL, "LearnX", 2*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = client.SendOTP(context.Background(), "user@example.com")
	if err == nil {
		t.Fatalf("expected error for 503 response, got nil")
	}
	if !errors.Is(err, otp.ErrOTPServiceUnavailable) {
		t.Errorf("expected ErrOTPServiceUnavailable, got %v", err)
	}
}

func TestOTPClient_SendOTP_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := otp.NewClient(server.URL, "LearnX", 10*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = client.SendOTP(context.Background(), "user@example.com")
	if err == nil {
		t.Fatalf("expected timeout error, got nil")
	}
}

func TestOTPClient_VerifyOTP_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/verify-otp" {
			t.Errorf("expected /api/verify-otp, got %s", r.URL.Path)
		}

		var req otp.VerifyOTPRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode: %v", err)
		}

		if req.Email != "user@example.com" || req.Code != "123456" {
			t.Errorf("unexpected request data: %+v", req)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(otp.VerifyOTPResponse{
			Success: true,
			Message: "OTP verified successfully",
			AppName: "LearnX",
		})
	}))
	defer server.Close()

	client, err := otp.NewClient(server.URL, "LearnX", 2*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resp, err := client.VerifyOTP(context.Background(), "user@example.com", "123456")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if !resp.Success {
		t.Fatalf("expected success true")
	}
}

func TestOTPClient_VerifyOTP_InvalidCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Invalid OTP code",
		})
	}))
	defer server.Close()

	client, err := otp.NewClient(server.URL, "LearnX", 2*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = client.VerifyOTP(context.Background(), "user@example.com", "000000")
	if err == nil {
		t.Fatalf("expected error for invalid OTP, got nil")
	}
	if !errors.Is(err, otp.ErrOTPVerificationFailed) {
		t.Errorf("expected ErrOTPVerificationFailed, got: %v", err)
	}
}

func TestOTPClient_VerifyOTP_ExpiredOrNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "OTP not found for LearnX. Please request a new one",
		})
	}))
	defer server.Close()

	client, err := otp.NewClient(server.URL, "LearnX", 2*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = client.VerifyOTP(context.Background(), "user@example.com", "999999")
	if err == nil {
		t.Fatalf("expected error for expired/not found OTP, got nil")
	}
	if !errors.Is(err, otp.ErrOTPVerificationFailed) {
		t.Errorf("expected ErrOTPVerificationFailed, got: %v", err)
	}
}

type mockSimpleClient struct {
	sendFunc   func(ctx context.Context, email string) (*otp.SendOTPResponse, error)
	verifyFunc func(ctx context.Context, email string, code string) (*otp.VerifyOTPResponse, error)
}

func (m *mockSimpleClient) SendOTP(ctx context.Context, email string) (*otp.SendOTPResponse, error) {
	if m.sendFunc != nil {
		return m.sendFunc(ctx, email)
	}
	return &otp.SendOTPResponse{Success: true, Message: "OTP sent successfully"}, nil
}

func (m *mockSimpleClient) VerifyOTP(ctx context.Context, email string, code string) (*otp.VerifyOTPResponse, error) {
	if m.verifyFunc != nil {
		return m.verifyFunc(ctx, email, code)
	}
	return &otp.VerifyOTPResponse{Success: true, Message: "OTP verified successfully"}, nil
}

func TestOTPService_CooldownAndAttempts(t *testing.T) {
	mock := &mockSimpleClient{}
	svc := otp.NewService(mock, 1*time.Second, 3)

	// First send should succeed
	resp, err := svc.SendOTP(context.Background(), "student@learnx.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Success {
		t.Fatalf("expected success")
	}

	// Immediate second send fails with ErrCooldownActive
	_, err = svc.SendOTP(context.Background(), "student@learnx.com")
	if !errors.Is(err, otp.ErrCooldownActive) {
		t.Fatalf("expected ErrCooldownActive, got: %v", err)
	}

	// Wait for cooldown
	time.Sleep(1050 * time.Millisecond)

	// Now send should succeed again
	_, err = svc.SendOTP(context.Background(), "student@learnx.com")
	if err != nil {
		t.Fatalf("expected success after cooldown, got: %v", err)
	}

	// Max attempts verification
	mock.verifyFunc = func(ctx context.Context, email, code string) (*otp.VerifyOTPResponse, error) {
		return nil, otp.ErrOTPVerificationFailed
	}

	for i := 1; i <= 3; i++ {
		_, err := svc.VerifyOTP(context.Background(), "student@learnx.com", "123456")
		if !errors.Is(err, otp.ErrOTPVerificationFailed) {
			t.Fatalf("attempt %d: expected ErrOTPVerificationFailed, got %v", i, err)
		}
	}

	// 4th attempt blocked by max attempts
	_, err = svc.VerifyOTP(context.Background(), "student@learnx.com", "123456")
	if !errors.Is(err, otp.ErrTooManyAttempts) {
		t.Fatalf("expected ErrTooManyAttempts, got: %v", err)
	}
}
