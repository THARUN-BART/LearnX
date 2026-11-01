package otp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client defines the interface for interacting with the external OTP service.
type Client interface {
	SendOTP(ctx context.Context, email string) (*SendOTPResponse, error)
	VerifyOTP(ctx context.Context, email string, code string) (*VerifyOTPResponse, error)
}

// HTTPClient implements the Client interface using Go's net/http package.
type HTTPClient struct {
	baseURL    string
	appName    string
	httpClient *http.Client
	logger     *slog.Logger
}

// ClientOption defines functional options for HTTPClient.
type ClientOption func(*HTTPClient)

// WithHTTPClient allows passing a custom net/http Client (e.g. for testing with httptest).
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *HTTPClient) {
		c.httpClient = client
	}
}

// WithLogger sets a structured logger.
func WithLogger(logger *slog.Logger) ClientOption {
	return func(c *HTTPClient) {
		c.logger = logger
	}
}

// NewClient creates a new HTTPClient instance for the external OTP service.
func NewClient(baseURL string, appName string, timeout time.Duration, opts ...ClientOption) (*HTTPClient, error) {
	if baseURL == "" {
		baseURL = "https://otp-latest-0k73.onrender.com"
	}
	parsedURL, err := url.Parse(baseURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, fmt.Errorf("invalid OTP service base URL: %s", baseURL)
	}

	if appName == "" {
		appName = "LearnX"
	}

	c := &HTTPClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		appName: appName,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		logger: slog.Default(),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// SendOTP calls POST /api/send-otp to dispatch an OTP to the given email.
// Never logs OTP codes.
func (c *HTTPClient) SendOTP(ctx context.Context, email string) (*SendOTPResponse, error) {
	endpoint := c.baseURL + "/api/send-otp"

	reqBody := SendOTPRequest{
		Email:   strings.TrimSpace(email),
		AppName: c.appName,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal send-otp request: %w", err)
	}

	c.logger.Info("sending OTP request to external service",
		slog.String("email", reqBody.Email),
		slog.String("app_name", c.appName),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("failed to connect to OTP service", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%w: %v", ErrOTPServiceUnavailable, err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var res SendOTPResponse
	if err := json.Unmarshal(bodyBytes, &res); err != nil {
		c.logger.Error("failed to unmarshal OTP response", slog.String("body", string(bodyBytes)))
		return nil, fmt.Errorf("%w: response is not valid JSON", ErrUnexpectedResponse)
	}

	if resp.StatusCode >= 500 {
		c.logger.Error("OTP service returned server error", slog.Int("status", resp.StatusCode), slog.String("error", res.Error))
		return nil, fmt.Errorf("%w: status %d", ErrOTPServiceUnavailable, resp.StatusCode)
	}

	if resp.StatusCode >= 400 || !res.Success {
		errMsg := res.Message
		if errMsg == "" {
			errMsg = res.Error
		}
		if errMsg == "" {
			errMsg = fmt.Sprintf("status %d", resp.StatusCode)
		}
		c.logger.Warn("OTP send request rejected by provider", slog.String("message", errMsg), slog.Int("status", resp.StatusCode))
		return nil, fmt.Errorf("failed to send OTP: %s", errMsg)
	}

	return &res, nil
}

// VerifyOTP calls POST /api/verify-otp to verify the provided code.
// Never logs OTP codes under any circumstances!
func (c *HTTPClient) VerifyOTP(ctx context.Context, email string, code string) (*VerifyOTPResponse, error) {
	endpoint := c.baseURL + "/api/verify-otp"

	reqBody := VerifyOTPRequest{
		Email:   strings.TrimSpace(email),
		Code:    strings.TrimSpace(code),
		AppName: c.appName,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal verify-otp request: %w", err)
	}

	// Security requirement: Never log OTP codes!
	c.logger.Info("verifying OTP with external service",
		slog.String("email", reqBody.Email),
		slog.String("app_name", c.appName),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("failed to connect to OTP service for verification", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%w: %v", ErrOTPServiceUnavailable, err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var res VerifyOTPResponse
	if err := json.Unmarshal(bodyBytes, &res); err != nil {
		c.logger.Error("failed to unmarshal OTP verify response", slog.String("body", string(bodyBytes)))
		return nil, fmt.Errorf("%w: response is not valid JSON", ErrUnexpectedResponse)
	}

	if resp.StatusCode >= 500 {
		c.logger.Error("OTP service returned server error during verification", slog.Int("status", resp.StatusCode))
		return nil, fmt.Errorf("%w: status %d", ErrOTPServiceUnavailable, resp.StatusCode)
	}

	if resp.StatusCode >= 400 || !res.Success {
		errMsg := res.Message
		if errMsg == "" {
			errMsg = res.Error
		}
		if errMsg == "" {
			errMsg = "Invalid or expired OTP code"
		}
		c.logger.Warn("OTP verification failed", slog.String("email", reqBody.Email), slog.String("reason", errMsg))
		return nil, fmt.Errorf("%w: %s", ErrOTPVerificationFailed, errMsg)
	}

	return &res, nil
}
