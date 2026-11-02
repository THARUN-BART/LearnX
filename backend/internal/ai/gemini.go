package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type GeminiProvider struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
}

func NewGeminiProvider(apiKey, baseURL, model string, timeout time.Duration) *GeminiProvider {
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com/v1beta"
	}
	if model == "" {
		model = "gemini-1.5-flash"
	}
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &GeminiProvider{
		apiKey:  apiKey,
		baseURL: strings.TrimRight(baseURL, "/"),
		model:   model,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (p *GeminiProvider) Name() string {
	return "gemini"
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"` // "user" or "model"
	Parts []geminiPart `json:"parts"`
}

type geminiGenerationConfig struct {
	Temperature      float64 `json:"temperature,omitempty"`
	MaxOutputTokens  int     `json:"maxOutputTokens,omitempty"`
	ResponseMimeType string  `json:"responseMimeType,omitempty"`
}

type geminiGenerateRequest struct {
	Contents          []geminiContent         `json:"contents"`
	SystemInstruction *geminiContent          `json:"systemInstruction,omitempty"`
	GenerationConfig  *geminiGenerationConfig `json:"generationConfig,omitempty"`
}

type geminiGenerateResponse struct {
	Candidates []struct {
		Content struct {
			Parts []geminiPart `json:"parts"`
			Role  string       `json:"role"`
		} `json:"content"`
		FinishReason string `json:"finishReason"`
	} `json:"candidates"`
	UsageMetadata struct {
		PromptTokenCount     int `json:"promptTokenCount"`
		CandidatesTokenCount int `json:"candidatesTokenCount"`
		TotalTokenCount      int `json:"totalTokenCount"`
	} `json:"usageMetadata"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error,omitempty"`
}

func (p *GeminiProvider) Generate(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	if p.apiKey == "" {
		return nil, errors.New("gemini api key is not configured")
	}

	modelToUse := req.Model
	if modelToUse == "" {
		modelToUse = p.model
	}

	bodyReq := geminiGenerateRequest{
		Contents: make([]geminiContent, 0),
	}

	// Separate system instructions from conversational dialog
	for _, msg := range req.Messages {
		if msg.Role == "system" {
			bodyReq.SystemInstruction = &geminiContent{
				Parts: []geminiPart{{Text: msg.Content}},
			}
		} else {
			role := "user"
			if msg.Role == "assistant" {
				role = "model"
			}
			bodyReq.Contents = append(bodyReq.Contents, geminiContent{
				Role:  role,
				Parts: []geminiPart{{Text: msg.Content}},
			})
		}
	}

	// Fallback: If no user messages were passed, create one from system
	if len(bodyReq.Contents) == 0 && bodyReq.SystemInstruction != nil {
		bodyReq.Contents = append(bodyReq.Contents, geminiContent{
			Role:  "user",
			Parts: bodyReq.SystemInstruction.Parts,
		})
		bodyReq.SystemInstruction = nil
	}

	genConfig := &geminiGenerationConfig{
		Temperature: req.Temperature,
	}
	if req.MaxTokens > 0 {
		genConfig.MaxOutputTokens = req.MaxTokens
	}
	if req.ResponseFormat == "json_object" {
		genConfig.ResponseMimeType = "application/json"
	}
	bodyReq.GenerationConfig = genConfig

	payloadBytes, err := json.Marshal(bodyReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal gemini request: %w", err)
	}

	targetURL := fmt.Sprintf("%s/models/%s:generateContent?key=%s", p.baseURL, modelToUse, p.apiKey)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to build gemini http request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("gemini request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read gemini response: %w", err)
	}

	var parsed geminiGenerateResponse
	if err := json.Unmarshal(bodyBytes, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse gemini json response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if parsed.Error != nil && parsed.Error.Message != "" {
			return nil, fmt.Errorf("gemini error (status %d): %s", resp.StatusCode, parsed.Error.Message)
		}
		return nil, fmt.Errorf("gemini returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if len(parsed.Candidates) == 0 || len(parsed.Candidates[0].Content.Parts) == 0 {
		return nil, errors.New("gemini returned no content candidates")
	}

	var contentBuilder strings.Builder
	for _, part := range parsed.Candidates[0].Content.Parts {
		contentBuilder.WriteString(part.Text)
	}

	return &ChatCompletionResponse{
		Content:      contentBuilder.String(),
		FinishReason: parsed.Candidates[0].FinishReason,
		Usage: TokenUsage{
			PromptTokens:     parsed.UsageMetadata.PromptTokenCount,
			CompletionTokens: parsed.UsageMetadata.CandidatesTokenCount,
			TotalTokens:      parsed.UsageMetadata.TotalTokenCount,
		},
	}, nil
}
