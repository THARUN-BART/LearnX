package tests

import (
	"testing"

	"github.com/THARUN-BART/LearnX/backend/internal/codereview"
)

func TestStaticCodeAnalysis(t *testing.T) {
	badPythonCode := `def search(items, targets):
    for x in targets:
        for y in items:
            for z in items:
                for w in items:
                    if x in list(y):
                        return 1000000007
`
	issues := codereview.RunStaticAnalysis(badPythonCode, "python")
	if len(issues) == 0 {
		t.Fatalf("expected static analysis to flag deep nesting and quadratic lookup")
	}

	foundNesting := false
	foundMagicNum := false
	for _, iss := range issues {
		if iss.Category == "maintainability" {
			foundNesting = true
		}
		if iss.Category == "readability" {
			foundMagicNum = true
		}
	}

	if !foundNesting {
		t.Errorf("expected deep nesting issue to be detected")
	}
	if !foundMagicNum {
		t.Errorf("expected magic constant issue to be detected")
	}
}

func TestVisualizerRecommendationInference(t *testing.T) {
	codeWithPointers := `def maxArea(height):
    left = 0
    right = len(height) - 1
    return 0`

	algoID, reason := codereview.InferVisualizerRecommendation(codeWithPointers, "python", nil)
	if algoID != "two-pointers" {
		t.Errorf("expected 'two-pointers', got '%s'", algoID)
	}
	if reason == "" {
		t.Errorf("expected non-empty reason")
	}
}

