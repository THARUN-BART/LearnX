package ai

import "context"

type ChatMessage struct {
	Role    string `json:"role"` // system, user, assistant
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model          string        `json:"model"`
	Messages       []ChatMessage `json:"messages"`
	Temperature    float64       `json:"temperature,omitempty"`
	MaxTokens      int           `json:"max_tokens,omitempty"`
	ResponseFormat string        `json:"response_format,omitempty"` // "json_object" or empty
}

type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatCompletionResponse struct {
	Content      string     `json:"content"`
	FinishReason string     `json:"finish_reason"`
	Usage        TokenUsage `json:"usage"`
}

// LLMProvider abstracts interactions with language models (OpenAI, Anthropic, Gemini, Mock).
type LLMProvider interface {
	Generate(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error)
	Name() string
}
