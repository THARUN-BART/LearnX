package codereview

import (
	"fmt"
	"strings"
)

// RunStaticAnalysis performs deterministic rule-based checks on student code.
func RunStaticAnalysis(code, lang string) []CodeIssue {
	issues := make([]CodeIssue, 0)
	lines := strings.Split(code, "\n")
	lang = strings.ToLower(lang)

	nestingLevel := 0
	for lineIdx, line := range lines {
		lineNum := lineIdx + 1
		trimmed := strings.TrimSpace(line)

		// 1. Deep nesting check
		if strings.Contains(line, "{") || (lang == "python" && strings.HasSuffix(trimmed, ":")) {
			nestingLevel++
			if nestingLevel >= 4 {
				issues = append(issues, CodeIssue{
					ID:             fmt.Sprintf("nesting-%d", lineNum),
					Severity:       SeverityWarning,
					LineNumber:     lineNum,
					Title:          "Excessive Nested Control Flow",
					Description:    fmt.Sprintf("Indentation nesting level reached %d. Deeply nested blocks increase cognitive complexity and bug risk.", nestingLevel),
					Recommendation: "Extract inner logic into dedicated helper functions or invert guard conditions to return early.",
					CodeSnippet:    trimmed,
					Category:       "maintainability",
					SourceType:     "static_analysis",
				})
			}
		}
		if strings.Contains(line, "}") {
			if nestingLevel > 0 {
				nestingLevel--
			}
		}

		// 2. Magic Numbers check
		if !strings.HasPrefix(trimmed, "//") && !strings.HasPrefix(trimmed, "#") {
			if strings.Contains(trimmed, "100000") || strings.Contains(trimmed, "999999") || strings.Contains(trimmed, "1000000007") {
				if !strings.Contains(trimmed, "MOD") && !strings.Contains(trimmed, "const") && !strings.Contains(trimmed, "final") {
					issues = append(issues, CodeIssue{
						ID:             fmt.Sprintf("magic-num-%d", lineNum),
						Severity:       SeverityInfo,
						LineNumber:     lineNum,
						Title:          "Unlabeled Magic Constant",
						Description:    "Literal numeric constant used directly in conditional or arithmetic statement.",
						Recommendation: "Declare meaningful named constants (e.g. `const MOD = 1e9 + 7`) at the top of the file.",
						CodeSnippet:    trimmed,
						Category:       "readability",
						SourceType:     "static_analysis",
					})
				}
			}
		}

		// 3. Busy loops / quadratic lookups
		if lang == "python" {
			if strings.Contains(trimmed, "for ") && strings.Contains(trimmed, " in ") {
				if strings.Contains(trimmed, ".index(") || strings.Contains(trimmed, " in list(") {
					issues = append(issues, CodeIssue{
						ID:             fmt.Sprintf("quadratic-search-%d", lineNum),
						Severity:       SeverityCritical,
						LineNumber:     lineNum,
						Title:          "O(N) Search Inside Loop (O(N²) Quadratic Hazard)",
						Description:    "Calling linear membership `.index()` or searching inside a list within a loop produces quadratic O(N²) running time.",
						Recommendation: "Convert the target collection to a `set` or `dict` beforehand to obtain O(1) hash table lookups.",
						CodeSnippet:    trimmed,
						Category:       "performance",
						SourceType:     "static_analysis",
					})
				}
			}
		}

		// 4. Infinite loop risk
		if strings.Contains(trimmed, "while true") || strings.Contains(trimmed, "while (true)") || strings.Contains(trimmed, "while 1:") {
			if !strings.Contains(code, "break") && !strings.Contains(code, "return") {
				issues = append(issues, CodeIssue{
					ID:             fmt.Sprintf("infinite-loop-%d", lineNum),
					Severity:       SeverityCritical,
					LineNumber:     lineNum,
					Title:          "Potential Infinite Loop",
					Description:    "Unconditional `while(true)` without apparent terminating `break` or `return` statement.",
					Recommendation: "Ensure all execution paths guarantee progress towards termination, or replace with bounded loop.",
					CodeSnippet:    trimmed,
					Category:       "safety",
					SourceType:     "static_analysis",
				})
			}
		}
	}

	return issues
}

// InferVisualizerRecommendation maps code patterns to relevant visualizer algorithms
func InferVisualizerRecommendation(code, language string, issues []CodeIssue) (string, string) {
	codeLower := strings.ToLower(code)

	// Check if two pointers is relevant
	if strings.Contains(codeLower, "left") && strings.Contains(codeLower, "right") {
		return "two-pointers", "Your code employs two pointers. Explore the Two Pointers interactive visualizer to see optimal convergence steps."
	}
	// Sorting inefficiency
	if strings.Contains(codeLower, "for ") && strings.Contains(codeLower, "range(len(") {
		for _, iss := range issues {
			if iss.Category == "performance" {
				return "merge-sort", "Detected quadratic iterative loop pattern. Explore Merge Sort or Quick Sort visualizers for optimal O(N log N) divide-and-conquer."
			}
		}
	}
	// Linked list
	if strings.Contains(codeLower, ".next") || strings.Contains(codeLower, "->next") {
		return "linked-list-reversal", "You are manipulating pointers on linked nodes. Inspect the Linked List Reversal visualizer for in-place pointer rewiring."
	}
	// Trees
	if strings.Contains(codeLower, ".left") || strings.Contains(codeLower, "->left") {
		return "binary-tree-traversals", "Explore the Binary Tree Traversals visualizer to master Inorder, Preorder, and Postorder recursion."
	}
	// Graph BFS/DFS
	if strings.Contains(codeLower, "queue") || strings.Contains(codeLower, "visited") {
		return "graph-bfs", "Explore the Graph BFS visualizer to inspect queue level-order expansions."
	}
	// Dynamic Programming
	if strings.Contains(codeLower, "dp[") {
		return "dp-fibonacci", "Inspect the Dynamic Programming Tabulation visualizer to master subproblem reuse."
	}

	return "binary-search", "Explore the Binary Search visualizer to master logarithmic boundary narrowing."
}

