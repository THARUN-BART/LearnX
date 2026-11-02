package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OpenAIProvider struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
}

func NewOpenAIProvider(apiKey, baseURL, model string, timeout time.Duration) *OpenAIProvider {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	if model == "" {
		model = "gpt-4o-mini"
	}
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &OpenAIProvider{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   model,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (p *OpenAIProvider) Name() string {
	return "openai"
}

type openAIChatRequest struct {
	Model          string                  `json:"model"`
	Messages       []ChatMessage           `json:"messages"`
	Temperature    float64                 `json:"temperature,omitempty"`
	MaxTokens      int                     `json:"max_tokens,omitempty"`
	ResponseFormat *openAIResponseFormat   `json:"response_format,omitempty"`
}

type openAIResponseFormat struct {
	Type string `json:"type"`
}

type openAIChatResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

func (p *OpenAIProvider) Generate(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	if p.apiKey == "" {
		return nil, errors.New("openai api key is not configured")
	}

	modelToUse := req.Model
	if modelToUse == "" {
		modelToUse = p.model
	}

	bodyReq := openAIChatRequest{
		Model:       modelToUse,
		Messages:    req.Messages,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
	}
	if req.ResponseFormat == "json_object" {
		bodyReq.ResponseFormat = &openAIResponseFormat{Type: "json_object"}
	}

	payloadBytes, err := json.Marshal(bodyReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal openai request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/chat/completions", bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to build openai http request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("openai request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read openai response: %w", err)
	}

	var parsed openAIChatResponse
	if err := json.Unmarshal(bodyBytes, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse openai json response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if parsed.Error != nil && parsed.Error.Message != "" {
			return nil, fmt.Errorf("openai error (status %d): %s", resp.StatusCode, parsed.Error.Message)
		}
		return nil, fmt.Errorf("openai returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if len(parsed.Choices) == 0 {
		return nil, errors.New("openai returned no choices")
	}

	return &ChatCompletionResponse{
		Content:      parsed.Choices[0].Message.Content,
		FinishReason: parsed.Choices[0].FinishReason,
		Usage: TokenUsage{
			PromptTokens:     parsed.Usage.PromptTokens,
			CompletionTokens: parsed.Usage.CompletionTokens,
			TotalTokens:      parsed.Usage.TotalTokens,
		},
	}, nil
}
