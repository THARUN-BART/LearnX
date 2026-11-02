package ai

import (
	"context"
	"encoding/json"
	"strings"
)

type MockProvider struct {
	CustomResponse *ChatCompletionResponse
	CustomError    error
	CapturedReqs   []ChatCompletionRequest
}

func NewMockProvider() *MockProvider {
	return &MockProvider{
		CapturedReqs: make([]ChatCompletionRequest, 0),
	}
}

func (m *MockProvider) Name() string {
	return "mock"
}

func (m *MockProvider) Generate(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	m.CapturedReqs = append(m.CapturedReqs, req)
	if m.CustomError != nil {
		return nil, m.CustomError
	}
	if m.CustomResponse != nil {
		return m.CustomResponse, nil
	}

	// Detect request type from messages to provide realistic mock replies
	lastMsg := ""
	if len(req.Messages) > 0 {
		lastMsg = strings.ToLower(req.Messages[len(req.Messages)-1].Content)
	}
	systemMsg := ""
	if len(req.Messages) > 0 && req.Messages[0].Role == "system" {
		systemMsg = strings.ToLower(req.Messages[0].Content)
	}

	content := "This is a simulated AI Tutor response for your DSA query."

	if req.ResponseFormat == "json_object" || strings.Contains(systemMsg, "json") || strings.Contains(lastMsg, "json") {
		switch {
		case strings.Contains(systemMsg, "dry-run simulation") || strings.Contains(lastMsg, "dry_run") || strings.Contains(lastMsg, "dry-run"):
			contentObj := map[string]interface{}{
				"test_case_input": "[2, 7, 11, 15], target = 9",
				"steps": []map[string]interface{}{
					{
						"step_number": 1,
						"line_number": 3,
						"action":      "Initialize hash map for complements",
						"variables":   map[string]interface{}{"lookup": "{}"},
						"explanation": "Map is empty at the start of iteration.",
					},
					{
						"step_number": 2,
						"line_number": 4,
						"action":      "Check complement 9 - 2 = 7 in map",
						"variables":   map[string]interface{}{"num": 2, "complement": 7},
						"explanation": "7 not found in map, storing 2 at index 0.",
					},
					{
						"step_number": 3,
						"line_number": 4,
						"action":      "Found complement 9 - 7 = 2 in map",
						"variables":   map[string]interface{}{"num": 7, "complement": 2, "index": 1},
						"explanation": "Found complement 2 at index 0. Returning [0, 1].",
					},
				},
				"final_output": "[0, 1]",
			}
			b, _ := json.Marshal(contentObj)
			content = string(b)

		case strings.Contains(systemMsg, "time and space complexity") || strings.Contains(lastMsg, "complexity"):
			contentObj := map[string]interface{}{
				"time_complexity":    "O(N)",
				"time_explanation":   "We iterate through the array of N elements exactly once, doing O(1) hash map lookups.",
				"space_complexity":   "O(N)",
				"space_explanation":  "In the worst case, we store up to N elements in the hash map.",
				"is_optimal":         true,
				"optimisation_notes": "This matches the theoretical lower bound for this problem.",
			}
			b, _ := json.Marshal(contentObj)
			content = string(b)

		case strings.Contains(systemMsg, "code walkthrough") || strings.Contains(lastMsg, "walkthrough"):
			contentObj := map[string]interface{}{
				"overall_summary": "Clean iterative approach using hash map.",
				"sections": []map[string]interface{}{
					{
						"title":       "Initialization",
						"line_range":  "Lines 1-3",
						"observation": "Allocates lookup table",
						"evaluation":  "positive",
					},
				},
				"key_strengths":   []string{"Optimal time complexity O(N)"},
				"areas_to_improve": []string{"Handle empty arrays early"},
			}
			b, _ := json.Marshal(contentObj)
			content = string(b)

		case strings.Contains(systemMsg, "level 1") || strings.Contains(systemMsg, "level 2") ||
			strings.Contains(systemMsg, "level 3") || strings.Contains(systemMsg, "level 4") ||
			strings.Contains(systemMsg, "level 5") || strings.Contains(lastMsg, "hint"):
			contentObj := map[string]interface{}{
				"hint_level":       1,
				"title":            "Conceptual Strategy",
				"content":          "Think about how you can avoid checking every pair by remembering numbers you have already visited using a hash table.",
				"guiding_question": "If you are at element x, what number y would satisfy x + y = target, and how quickly can you check if y was seen?",
			}
			b, _ := json.Marshal(contentObj)
			content = string(b)

		default:
			contentObj := map[string]interface{}{
				"summary":     "Analysis completed successfully.",
				"explanation": "Your code approach is solid, but remember to verify edge cases such as empty input arrays.",
			}
			b, _ := json.Marshal(contentObj)
			content = string(b)
		}
	}

	return &ChatCompletionResponse{
		Content:      content,
		FinishReason: "stop",
		Usage: TokenUsage{
			PromptTokens:     85,
			CompletionTokens: 120,
			TotalTokens:      205,
		},
	}, nil
}
