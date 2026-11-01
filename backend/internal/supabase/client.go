package supabase

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

type Client interface {
	AdminCreateUser(ctx context.Context, req AdminCreateUserRequest) (*SupabaseUser, error)
	AdminGetUserByEmail(ctx context.Context, email string) (*SupabaseUser, error)
	AdminUpdateUserPassword(ctx context.Context, userID string, newPassword string) error
	SignInWithPassword(ctx context.Context, email, password string) (*TokenResponse, error)
	SignOut(ctx context.Context, accessToken string) error
	GetUser(ctx context.Context, accessToken string) (*SupabaseUser, error)
	GetOAuthURL(provider, redirectTo string) string
}

type HTTPClient struct {
	baseURL        string
	anonKey        string
	serviceRoleKey string
	httpClient     *http.Client
	logger         *slog.Logger
}

type ClientOption func(*HTTPClient)

func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *HTTPClient) {
		c.httpClient = client
	}
}

func WithLogger(logger *slog.Logger) ClientOption {
	return func(c *HTTPClient) {
		c.logger = logger
	}
}

func NewClient(baseURL, anonKey, serviceRoleKey string, opts ...ClientOption) *HTTPClient {
	c := &HTTPClient{
		baseURL:        strings.TrimRight(baseURL, "/"),
		anonKey:        anonKey,
		serviceRoleKey: serviceRoleKey,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		logger: slog.Default(),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// AdminCreateUser creates a user using the Supabase GoTrue Admin API.
// Email verification is performed securely server-side.
func (c *HTTPClient) AdminCreateUser(ctx context.Context, req AdminCreateUserRequest) (*SupabaseUser, error) {
	endpoint := c.baseURL + "/auth/v1/admin/users"

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal admin create user payload: %w", err)
	}

	c.logger.Info("admin creating user in Supabase Auth", slog.String("email", req.Email))

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("apikey", c.serviceRoleKey)
	httpReq.Header.Set("Authorization", "Bearer "+c.serviceRoleKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		c.logger.Error("supabase admin create user request failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%w: %v", ErrSupabaseUnavailable, err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode == http.StatusConflict || resp.StatusCode == http.StatusUnprocessableEntity {
		return nil, ErrUserAlreadyExists
	}

	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		_ = json.Unmarshal(bodyBytes, &errResp)
		msg := errResp.GetMessage()
		if strings.Contains(strings.ToLower(msg), "already") || strings.Contains(strings.ToLower(msg), "exists") {
			return nil, ErrUserAlreadyExists
		}
		return nil, fmt.Errorf("supabase admin create user failed (status %d): %s", resp.StatusCode, msg)
	}

	var user SupabaseUser
	if err := json.Unmarshal(bodyBytes, &user); err != nil {
		return nil, fmt.Errorf("failed to decode user response: %w", err)
	}

	return &user, nil
}

// AdminGetUserByEmail checks if a user exists in Supabase Auth.
func (c *HTTPClient) AdminGetUserByEmail(ctx context.Context, email string) (*SupabaseUser, error) {
	endpoint := fmt.Sprintf("%s/auth/v1/admin/users?email=%s", c.baseURL, url.QueryEscape(email))

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("apikey", c.serviceRoleKey)
	httpReq.Header.Set("Authorization", "Bearer "+c.serviceRoleKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSupabaseUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("supabase admin get users failed with status %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	type listResponse struct {
		Users []SupabaseUser `json:"users"`
	}
	var res listResponse
	if err := json.Unmarshal(bodyBytes, &res); err != nil {
		// Try unmarshaling as direct array
		var users []SupabaseUser
		if errArray := json.Unmarshal(bodyBytes, &users); errArray == nil && len(users) > 0 {
			for _, u := range users {
				if strings.EqualFold(u.Email, email) {
					return &u, nil
				}
			}
		}
		return nil, ErrUserNotFound
	}

	for _, u := range res.Users {
		if strings.EqualFold(u.Email, email) {
			return &u, nil
		}
	}

	return nil, ErrUserNotFound
}

// AdminUpdateUserPassword updates a user's password using the admin credentials.
func (c *HTTPClient) AdminUpdateUserPassword(ctx context.Context, userID, newPassword string) error {
	endpoint := fmt.Sprintf("%s/auth/v1/admin/users/%s", c.baseURL, url.PathEscape(userID))

	payload := map[string]interface{}{
		"password": newPassword,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	c.logger.Info("admin updating password for user", slog.String("user_id", userID))

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPut, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("apikey", c.serviceRoleKey)
	httpReq.Header.Set("Authorization", "Bearer "+c.serviceRoleKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSupabaseUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("supabase admin update password failed (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// SignInWithPassword authenticates the user using email and password.
func (c *HTTPClient) SignInWithPassword(ctx context.Context, email, password string) (*TokenResponse, error) {
	endpoint := c.baseURL + "/auth/v1/token?grant_type=password"

	payload := map[string]string{
		"email":    email,
		"password": password,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	c.logger.Info("authenticating user credentials with Supabase", slog.String("email", email))

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("apikey", c.anonKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSupabaseUnavailable, err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		_ = json.Unmarshal(bodyBytes, &errResp)
		msg := strings.ToLower(errResp.GetMessage())
		if strings.Contains(msg, "invalid") || strings.Contains(msg, "credential") {
			return nil, ErrInvalidCredentials
		}
		if strings.Contains(msg, "confirm") || strings.Contains(msg, "unverified") {
			return nil, ErrEmailNotConfirmed
		}
		return nil, fmt.Errorf("%w: %s", ErrInvalidCredentials, errResp.GetMessage())
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(bodyBytes, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &tokenResp, nil
}

// SignOut revokes the user's session token.
func (c *HTTPClient) SignOut(ctx context.Context, accessToken string) error {
	endpoint := c.baseURL + "/auth/v1/logout"

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, nil)
	if err != nil {
		return err
	}
	httpReq.Header.Set("apikey", c.anonKey)
	httpReq.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSupabaseUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 && resp.StatusCode != http.StatusNotFound {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("supabase sign out failed: %s", string(body))
	}

	return nil
}

// GetUser retrieves the user information associated with an access token.
func (c *HTTPClient) GetUser(ctx context.Context, accessToken string) (*SupabaseUser, error) {
	endpoint := c.baseURL + "/auth/v1/user"

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("apikey", c.anonKey)
	httpReq.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrSupabaseUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrInvalidCredentials
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("supabase get user failed with status %d", resp.StatusCode)
	}

	var user SupabaseUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// GetOAuthURL builds the Google OAuth authorization URL.
func (c *HTTPClient) GetOAuthURL(provider, redirectTo string) string {
	if provider == "" {
		provider = "google"
	}
	u := fmt.Sprintf("%s/auth/v1/authorize?provider=%s", c.baseURL, url.QueryEscape(provider))
	if redirectTo != "" {
		u += fmt.Sprintf("&redirect_to=%s", url.QueryEscape(redirectTo))
	}
	return u
}
