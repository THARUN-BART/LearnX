package tests

import (
	"testing"

	"github.com/THARUN-BART/LearnX/backend/internal/visualizer"
)

func TestVisualizerRegistry(t *testing.T) {
	algos := visualizer.AlgorithmsRegistry
	if len(algos) < 30 {
		t.Fatalf("expected at least 30 algorithms in registry, got %d", len(algos))
	}

	for _, a := range algos {
		if a.ID == "" || a.Name == "" || a.Category == "" {
			t.Errorf("algorithm metadata missing required fields: %+v", a)
		}
	}
}

func TestSortingVisualizations(t *testing.T) {
	sortingAlgos := []string{"bubble-sort", "selection-sort", "insertion-sort", "merge-sort", "quick-sort", "heap-sort"}
	sample := map[string]any{"array": []int{45, 12, 85, 32, 89, 39, 69, 44}}

	for _, algoID := range sortingAlgos {
		steps, err := visualizer.GenerateExecutionSteps(algoID, sample)
		if err != nil {
			t.Fatalf("failed generating steps for %s: %v", algoID, err)
		}
		if len(steps) < 3 {
			t.Fatalf("%s produced too few steps: %d", algoID, len(steps))
		}
		lastStep := steps[len(steps)-1]
		if !lastStep.IsComplete {
			t.Errorf("last step of %s should be marked complete", algoID)
		}
	}
}

func TestSearchingVisualizations(t *testing.T) {
	input := map[string]any{"array": []int{2, 5, 8, 12, 16, 23, 38, 56, 72, 91}, "target": 23}

	for _, algoID := range []string{"linear-search", "binary-search"} {
		steps, err := visualizer.GenerateExecutionSteps(algoID, input)
		if err != nil {
			t.Fatalf("failed generating steps for %s: %v", algoID, err)
		}
		if len(steps) < 2 {
			t.Fatalf("%s produced too few steps: %d", algoID, len(steps))
		}
		lastStep := steps[len(steps)-1]
		if !lastStep.IsComplete {
			t.Errorf("last step of %s should be marked complete", algoID)
		}
	}
}

func TestDPVisualizations(t *testing.T) {
	dpAlgos := []string{"dp-fibonacci", "climbing-stairs", "knapsack-01", "longest-common-subsequence", "coin-change"}

	for _, algoID := range dpAlgos {
		steps, err := visualizer.GenerateExecutionSteps(algoID, nil)
		if err != nil {
			t.Fatalf("failed generating steps for %s: %v", algoID, err)
		}
		if len(steps) < 2 {
			t.Fatalf("%s produced too few steps: %d", algoID, len(steps))
		}
	}
}

func TestBacktrackingNQueens(t *testing.T) {
	steps, err := visualizer.GenerateExecutionSteps("backtracking-n-queens", map[string]any{"n": 4})
	if err != nil {
		t.Fatalf("N-Queens error: %v", err)
	}
	if len(steps) < 5 {
		t.Fatalf("N-Queens produced too few steps: %d", len(steps))
	}
}

