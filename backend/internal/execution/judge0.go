package execution

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	ErrExecutionServiceUnavailable = errors.New("execution service is unavailable")
	ErrExecutionTimeout            = errors.New("code execution exceeded timeout limit")
)

type Judge0Config struct {
	BaseURL    string
	APIKey     string
	APIHost    string
	Timeout    time.Duration
	HTTPClient *http.Client
	Logger     *slog.Logger
}

type Judge0Client struct {
	cfg        Judge0Config
	httpClient *http.Client
	logger     *slog.Logger
}

func NewJudge0Client(cfg Judge0Config) *Judge0Client {
	client := cfg.HTTPClient
	if client == nil {
		client = &http.Client{
			Timeout: cfg.Timeout,
		}
	}
	if client.Timeout == 0 {
		client.Timeout = 20 * time.Second
	}
	logger := cfg.Logger
	if logger == nil {
		logger = slog.Default()
	}

	return &Judge0Client{
		cfg:        cfg,
		httpClient: client,
		logger:     logger,
	}
}

func (c *Judge0Client) ProviderName() string {
	return "judge0"
}

type judge0SubmissionPayload struct {
	SourceCode     string   `json:"source_code"`
	LanguageID     int      `json:"language_id"`
	Stdin          string   `json:"stdin,omitempty"`
	ExpectedOutput string   `json:"expected_output,omitempty"`
	CPUTimeLimit   *float64 `json:"cpu_time_limit,omitempty"`
	MemoryLimitKB  *int     `json:"memory_limit,omitempty"`
}

type judge0SubmissionResponse struct {
	Token         string  `json:"token"`
	Stdout        *string `json:"stdout"`
	Stderr        *string `json:"stderr"`
	CompileOutput *string `json:"compile_output"`
	Message       *string `json:"message"`
	Time          *string `json:"time"`
	Memory        *int    `json:"memory"`
	Status        struct {
		ID          int    `json:"id"`
		Description string `json:"description"`
	} `json:"status"`
}

func (c *Judge0Client) Execute(ctx context.Context, req ExecutionRequest) (*ExecutionResponse, error) {
	baseURL := strings.TrimRight(c.cfg.BaseURL, "/")
	if baseURL == "" {
		baseURL = "http://localhost:2358"
	}
	url := fmt.Sprintf("%s/submissions/?base64_encoded=false&wait=true", baseURL)

	var cpuLimit *float64
	if req.CPUTimeLimit > 0 {
		sec := req.CPUTimeLimit.Seconds()
		cpuLimit = &sec
	}
	var memLimit *int
	if req.MemoryLimitKB > 0 {
		memLimit = &req.MemoryLimitKB
	}

	payload := judge0SubmissionPayload{
		SourceCode:     req.SourceCode,
		LanguageID:     req.LanguageID,
		Stdin:          req.Stdin,
		ExpectedOutput: req.ExpectedOutput,
		CPUTimeLimit:   cpuLimit,
		MemoryLimitKB:  memLimit,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal judge0 payload: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("create judge0 request: %w", err)
	}
	c.setHeaders(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		c.logger.Error("failed to call Judge0 execution service", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%w: %v", ErrExecutionServiceUnavailable, err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read judge0 response: %w", err)
	}

	if resp.StatusCode >= 500 {
		return nil, fmt.Errorf("%w: judge0 returned status %d", ErrExecutionServiceUnavailable, resp.StatusCode)
	}

	var jResp judge0SubmissionResponse
	if err := json.Unmarshal(respBytes, &jResp); err != nil {
		return nil, fmt.Errorf("unmarshal judge0 response: %w", err)
	}

	// If submission is still queued or processing (e.g. wait=true timed out), poll result
	if jResp.Status.ID == 1 || jResp.Status.ID == 2 {
		return c.pollResult(ctx, jResp.Token)
	}

	return c.mapResponse(&jResp), nil
}

func (c *Judge0Client) pollResult(ctx context.Context, token string) (*ExecutionResponse, error) {
	baseURL := strings.TrimRight(c.cfg.BaseURL, "/")
	url := fmt.Sprintf("%s/submissions/%s?base64_encoded=false", baseURL, token)

	ticker := time.NewTicker(400 * time.Millisecond)
	defer ticker.Stop()

	timeout := time.After(25 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-timeout:
			return nil, ErrExecutionTimeout
		case <-ticker.C:
			httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				return nil, err
			}
			c.setHeaders(httpReq)

			resp, err := c.httpClient.Do(httpReq)
			if err != nil {
				continue
			}
			bodyBytes, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			var jResp judge0SubmissionResponse
			if err := json.Unmarshal(bodyBytes, &jResp); err != nil {
				continue
			}

			if jResp.Status.ID != 1 && jResp.Status.ID != 2 {
				return c.mapResponse(&jResp), nil
			}
		}
	}
}

func (c *Judge0Client) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if c.cfg.APIKey != "" {
		if c.cfg.APIHost != "" {
			req.Header.Set("X-RapidAPI-Key", c.cfg.APIKey)
			req.Header.Set("X-RapidAPI-Host", c.cfg.APIHost)
		} else {
			req.Header.Set("X-Auth-Token", c.cfg.APIKey)
		}
	}
}

func (c *Judge0Client) mapResponse(j *judge0SubmissionResponse) *ExecutionResponse {
	out := &ExecutionResponse{
		Token:          j.Token,
		StatusID:       j.Status.ID,
		ProviderStatus: j.Status.Description,
	}

	if j.Stdout != nil {
		out.Stdout = *j.Stdout
	}
	if j.Stderr != nil {
		out.Stderr = *j.Stderr
	}
	if j.CompileOutput != nil {
		out.CompileOutput = *j.CompileOutput
	}
	if j.Message != nil {
		out.Message = *j.Message
	}
	if j.Time != nil {
		if sec, err := strconv.ParseFloat(*j.Time, 64); err == nil {
			out.TimeMs = int(sec * 1000)
		}
	}
	if j.Memory != nil {
		out.MemoryKB = *j.Memory
	}

	// Map Judge0 Status IDs to LearnX Statuses
	switch j.Status.ID {
	case 1:
		out.Status = StatusQueued
	case 2:
		out.Status = StatusProcessing
	case 3:
		out.Status = StatusAccepted
	case 4:
		out.Status = StatusWrongAnswer
	case 5:
		out.Status = StatusTimeLimitExceeded
	case 6:
		out.Status = StatusCompilationError
	case 7, 8, 9, 10, 11, 12:
		out.Status = StatusRuntimeError
	case 13:
		out.Status = StatusInternalError
	case 14:
		out.Status = StatusCompilationError
	default:
		out.Status = StatusInternalError
	}

	return out
}
