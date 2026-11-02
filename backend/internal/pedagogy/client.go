package pedagogy

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

var (
	ErrEmptyStudentCode   = errors.New("student_code cannot be empty")
	ErrServiceUnavailable = errors.New("pedagogy service is unreachable or returned an error")
)

type Client interface {
	Analyze(ctx context.Context, req AnalyzeRequest) (*AnalyzeResponse, error)
	Diagnose(ctx context.Context, req AnalyzeRequest) (*DiagnoseResponse, error)
	Health(ctx context.Context) (*HealthResponse, error)
}

type clientImpl struct {
	baseURL    string
	httpClient *http.Client
	logger     *slog.Logger
}

func NewClient(baseURL string, timeout time.Duration, logger *slog.Logger) Client {
	if baseURL == "" {
		baseURL = "http://localhost:8001"
	}
	baseURL = strings.TrimRight(baseURL, "/")
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	if logger == nil {
		logger = slog.Default()
	}

	return &clientImpl{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		logger: logger,
	}
}

func (c *clientImpl) Analyze(ctx context.Context, req AnalyzeRequest) (*AnalyzeResponse, error) {
	if strings.TrimSpace(req.StudentCode) == "" {
		return nil, ErrEmptyStudentCode
	}
	if req.Skill == "" {
		req.Skill = "SKL-02"
	}
	if req.HintLevel == "" {
		req.HintLevel = "L1"
	}
	if req.Timeout <= 0 {
		req.Timeout = 3.0
	}

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal analyze request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/analyze", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		c.logger.Warn("pedagogy analyze request failed", slog.String("url", c.baseURL+"/analyze"), slog.String("error", err.Error()))
		return nil, fmt.Errorf("%w: %v", ErrServiceUnavailable, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading pedagogy response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		c.logger.Warn("pedagogy returned non-2xx status", slog.Int("status", resp.StatusCode), slog.String("body", string(respBody)))
		return nil, fmt.Errorf("%w: status %d", ErrServiceUnavailable, resp.StatusCode)
	}

	var result AnalyzeResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed parsing pedagogy response: %w", err)
	}

	return &result, nil
}

func (c *clientImpl) Diagnose(ctx context.Context, req AnalyzeRequest) (*DiagnoseResponse, error) {
	if strings.TrimSpace(req.StudentCode) == "" {
		return nil, ErrEmptyStudentCode
	}
	if req.Skill == "" {
		req.Skill = "SKL-02"
	}
	if req.Timeout <= 0 {
		req.Timeout = 3.0
	}

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal diagnose request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/diagnose", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		c.logger.Warn("pedagogy diagnose request failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%w: %v", ErrServiceUnavailable, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading pedagogy response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%w: status %d", ErrServiceUnavailable, resp.StatusCode)
	}

	var result DiagnoseResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed parsing pedagogy response: %w", err)
	}

	return &result, nil
}

func (c *clientImpl) Health(ctx context.Context) (*HealthResponse, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/health", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var health HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return nil, err
	}

	return &health, nil
}
