package mistakes

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/THARUN-BART/LearnX/backend/internal/execution"
	"github.com/THARUN-BART/LearnX/backend/internal/submissions"
)

// Detector orchestrates the 4-layer mistake detection pipeline.
type Detector struct{}

func NewDetector() *Detector {
	return &Detector{}
}

// DetectionInput bundles all data needed for analysis.
type DetectionInput struct {
	Submission     *submissions.Submission
	TestResults    []submissions.SubmissionTestResult
	ActivePatterns []StudentMistakePattern
	ProblemTitle   string
	ProblemTopic   string
}

// Analyze runs Layer 1 -> Layer 2 -> Layer 3 -> Layer 4 detection.
func (d *Detector) Analyze(input DetectionInput) *Diagnosis {
	sub := input.Submission
	if sub == nil || sub.Status == execution.StatusAccepted {
		return nil
	}

	diag := &Diagnosis{
		StudentID:       sub.StudentID,
		SubmissionID:    sub.ID,
		ProblemID:       sub.ProblemID,
		Evidence:        make(map[string]interface{}),
		ConfidenceScore: 0.0,
	}

	// Layer 1: Deterministic Detection (Ground Truth Execution Facts)
	d.layer1Deterministic(input, diag)

	// Layer 2: Static Code Analysis (Syntactic & Structural Rules)
	d.layer2StaticAnalysis(input, diag)

	// Layer 3: Historical Pattern Analysis (Student Weakness Corroboration)
	d.layer3HistoricalPatterns(input, diag)

	// Layer 4: Root Cause & Conceptual Misconception Identification
	d.layer4RootCause(input, diag)

	// Determine overall confidence level based on score & evidence
	if diag.ConfidenceScore >= 0.85 {
		diag.ConfidenceLevel = ConfidenceHigh
	} else if diag.ConfidenceScore >= 0.60 {
		diag.ConfidenceLevel = ConfidenceMedium
	} else if diag.ConfidenceScore >= 0.35 {
		diag.ConfidenceLevel = ConfidenceLow
	} else {
		diag.ConfidenceLevel = ConfidenceInsufficientEvidence
	}

	return diag
}

func (d *Detector) layer1Deterministic(input DetectionInput, diag *Diagnosis) {
	sub := input.Submission
	errLower := strings.ToLower(sub.ErrorMessage)

	diag.Evidence["execution_status"] = sub.Status
	diag.Evidence["passed_test_cases"] = sub.PassedTestCases
	diag.Evidence["total_test_cases"] = sub.TotalTestCases
	diag.Evidence["execution_time_ms"] = sub.ExecutionTimeMs

	switch sub.Status {
	case execution.StatusCompilationError:
		diag.Category = CategoryCompilation
		diag.DetectionMethod = MethodDeterministic
		diag.ConfidenceScore = 0.95

		if strings.Contains(errLower, "syntax") || strings.Contains(errLower, "expected") || strings.Contains(errLower, "unexpected") {
			diag.Subtype = "syntax_error"
			diag.Description = "Your code failed to compile due to a syntax violation."
			diag.RelatedConcept = "Language Syntax & Grammatical Structure"
			diag.SuggestedCorrection = "Check matching parentheses, brackets, colons, or missing semicolons around the reported line."
		} else if strings.Contains(errLower, "mismatched types") || strings.Contains(errLower, "cannot convert") || strings.Contains(errLower, "type") {
			diag.Subtype = "type_mismatch"
			diag.Description = "A type mismatch occurred during compilation."
			diag.RelatedConcept = "Static Typing & Type Casting"
			diag.SuggestedCorrection = "Verify variable types and perform explicit type conversion where necessary."
		} else if strings.Contains(errLower, "undefined") || strings.Contains(errLower, "undeclared") || strings.Contains(errLower, "not found") {
			diag.Subtype = "undefined_identifier"
			diag.Description = "Reference to an undeclared or out-of-scope variable or function."
			diag.RelatedConcept = "Variable Scoping & Declarations"
			diag.SuggestedCorrection = "Ensure all identifiers are declared and accessible within the current scope before use."
		} else if strings.Contains(errLower, "missing return") || strings.Contains(errLower, "control reaches end of non-void") {
			diag.Subtype = "missing_return"
			diag.Description = "Function is missing a return statement along one or more execution branches."
			diag.RelatedConcept = "Function Return Contracts"
			diag.SuggestedCorrection = "Ensure every conditional path returns a value conforming to the function signature."
		} else {
			diag.Subtype = "compilation_failure"
			diag.Description = "Source code could not be compiled successfully."
			diag.RelatedConcept = "Compiler Build Pipeline"
			diag.SuggestedCorrection = "Inspect the compiler error output and correct invalid statements."
		}

	case execution.StatusRuntimeError:
		diag.Category = CategoryRuntime
		diag.DetectionMethod = MethodDeterministic
		diag.ConfidenceScore = 0.90

		if strings.Contains(errLower, "null pointer") || strings.Contains(errLower, "nil pointer") || strings.Contains(errLower, "nullpointerexception") {
			diag.Subtype = "null_pointer_dereference"
			diag.Description = "Attempted to dereference a null or nil pointer/reference."
			diag.RelatedConcept = "Pointer Safety & Reference Verification"
			diag.SuggestedCorrection = "Add defensive null/nil checks before accessing node or pointer properties (e.g., node != nil)."
		} else if strings.Contains(errLower, "index out of range") || strings.Contains(errLower, "out of bounds") || strings.Contains(errLower, "sigsegv") {
			diag.Subtype = "out_of_bounds"
			diag.Description = "Array or vector index accessed beyond its valid bounds."
			diag.RelatedConcept = "Array Boundaries & Safe Indexing"
			diag.SuggestedCorrection = "Inspect array indexing. Remember standard indices range from 0 to N-1."
		} else if strings.Contains(errLower, "divide by zero") || strings.Contains(errLower, "sigfpe") || strings.Contains(errLower, "division by zero") {
			diag.Subtype = "division_by_zero"
			diag.Description = "Arithmetic operation attempted division or modulo by zero."
			diag.RelatedConcept = "Arithmetic Safety & Non-Zero Invariants"
			diag.SuggestedCorrection = "Guard division and modulo operations with checks ensuring the denominator is non-zero."
		} else if strings.Contains(errLower, "stack overflow") || strings.Contains(errLower, "maximum recursion") {
			diag.Subtype = "stack_overflow"
			diag.Description = "Call stack depth exceeded due to infinite or excessively deep recursion."
			diag.RelatedConcept = "Recursion Base Cases & Call Stack Limits"
			diag.SuggestedCorrection = "Verify your recursive base case condition and ensure recursive parameters strictly converge toward it."
		} else {
			diag.Subtype = "runtime_exception"
			diag.Description = "The program aborted abnormally during execution."
			diag.RelatedConcept = "Runtime Safety & Exception Handling"
			diag.SuggestedCorrection = "Inspect variable states and operations prior to termination."
		}

	case execution.StatusTimeLimitExceeded:
		diag.Category = CategoryComplexity
		diag.Subtype = "time_limit_exceeded"
		diag.DetectionMethod = MethodDeterministic
		diag.RelatedConcept = "Time Complexity & Algorithmic Efficiency"
		diag.Description = "Your solution exceeded the allowable time limit for this problem."
		diag.SuggestedCorrection = "Review the constraints. If N >= 10^5, an O(N^2) brute force will TLE; aim for O(N) or O(N log N)."
		diag.ConfidenceScore = 0.90

	case execution.StatusMemoryLimitExceeded:
		diag.Category = CategoryComplexity
		diag.Subtype = "memory_limit_exceeded"
		diag.DetectionMethod = MethodDeterministic
		diag.RelatedConcept = "Space Complexity & Allocation"
		diag.Description = "Your solution exceeded the maximum allowable memory allocation."
		diag.SuggestedCorrection = "Avoid storing auxiliary tables or matrices when an in-place or streaming approach is possible."
		diag.ConfidenceScore = 0.90

	case execution.StatusWrongAnswer:
		diag.Category = CategoryLogical
		diag.Subtype = "incorrect_state_or_condition"
		diag.DetectionMethod = MethodDeterministic
		diag.RelatedConcept = "Logical Correctness & Invariants"
		diag.Description = "Your program produced output that differed from the expected answer on one or more test cases."
		diag.SuggestedCorrection = "Trace your code step-by-step with the failing input to verify intermediate states."
		diag.ConfidenceScore = 0.65
	}
}

func (d *Detector) layer2StaticAnalysis(input DetectionInput, diag *Diagnosis) {
	code := input.Submission.SourceCode
	if code == "" {
		return
	}

	// 1. Off-by-one loop boundary regexes: e.g., i <= arr.length or i <= len(...)
	offByOnePattern := regexp.MustCompile(`for\s*\(?.*?[<]=(\s*(?:len\(|arr\.length|v\.size\(\)|n|N|\.count))`)
	if offByOnePattern.MatchString(code) {
		diag.Evidence["static_suspicious_loop_bound"] = true
		if diag.Category == CategoryLogical || diag.Category == CategoryRuntime {
			diag.Category = CategoryLogical
			diag.Subtype = "off_by_one"
			diag.Description = "Suspicious loop boundary detected: using '<=' with array length frequently causes out-of-bounds or off-by-one errors."
			diag.RelatedConcept = "Loop Invariants & Boundary Off-by-One"
			diag.SuggestedCorrection = "Change the loop condition to '< length' instead of '<= length'."
			diag.DetectionMethod = MethodStatic
			diag.ConfidenceScore = 0.85
		}
	}

	// 2. Midpoint calculation integer overflow: (low + high) / 2
	midOverflowPattern := regexp.MustCompile(`\(\s*(?:low|left|l)\s*\+\s*(?:high|right|r)\s*\)\s*/\s*2`)
	if midOverflowPattern.MatchString(code) {
		diag.Evidence["static_midpoint_overflow_risk"] = true
		if diag.Subtype == "incorrect_state_or_condition" || diag.Category == CategoryLogical {
			diag.Category = CategoryEdgeCase
			diag.Subtype = "integer_overflow"
			diag.Description = "Calculating midpoint as (low + high) / 2 can trigger integer overflow when low and high are large."
			diag.RelatedConcept = "Integer Overflow & Safe Midpoint"
			diag.SuggestedCorrection = "Use 'low + (high - low) / 2' or 'left + (right - left) / 2' to compute midpoint safely."
			diag.DetectionMethod = MethodStatic
			diag.ConfidenceScore = 0.80
		}
	}

	// 3. Unchecked chained pointer access: e.g. .next.next or ->next->next
	chainedPointerPattern := regexp.MustCompile(`(?:->|\.)next(?:->|\.)next`)
	if chainedPointerPattern.MatchString(code) {
		diag.Evidence["static_chained_pointer_access"] = true
		if diag.Category == CategoryRuntime {
			diag.Category = CategoryDataStructure
			diag.Subtype = "null_pointer_traversal"
			diag.Description = "Chained node pointer traversal (.next.next) without checking that the intermediate node is non-nil."
			diag.RelatedConcept = "Linked List Pointer Invariants"
			diag.SuggestedCorrection = "Check that both `node != nil` and `node.next != nil` before advancing `node.next.next`."
			diag.DetectionMethod = MethodStatic
			diag.ConfidenceScore = 0.90
		}
	}

	// 4. Nested loops when problem constraints suggest O(N) or O(N log N)
	nestedLoopPattern := regexp.MustCompile(`for\b[^{]*\{[^}]*for\b`)
	if nestedLoopPattern.MatchString(code) {
		diag.Evidence["static_nested_loops_detected"] = true
		if diag.Category == CategoryComplexity && diag.Subtype == "time_limit_exceeded" {
			diag.Subtype = "quadratic_complexity"
			diag.Description = "Nested loops detected resulting in O(N^2) or worse time complexity, causing Time Limit Exceeded."
			diag.RelatedConcept = "Quadratic Time Complexity vs Linear/Log-linear"
			diag.SuggestedCorrection = "Consider using a Hash Map, Two Pointers, or Binary Search to reduce nested iterations to O(N) or O(N log N)."
			diag.DetectionMethod = MethodStatic
			diag.ConfidenceScore = 0.92
		}
	}

	// 5. Kadane's algorithm negative sum trap: initializing max_sum = 0
	kadaneZeroInitPattern := regexp.MustCompile(`(?:max_so_far|max_sum|maxSum|ans)\s*=\s*0`)
	if kadaneZeroInitPattern.MatchString(code) && strings.Contains(strings.ToLower(input.ProblemTitle), "subarray") {
		diag.Evidence["static_kadane_zero_init"] = true
		if diag.Category == CategoryLogical {
			diag.Category = CategoryEdgeCase
			diag.Subtype = "all_negative_elements"
			diag.Description = "Initialising maximum subarray sum to 0 fails when all elements in the input array are negative."
			diag.RelatedConcept = "Kadane's Algorithm & Extreme Values"
			diag.SuggestedCorrection = "Initialise your maximum accumulator to the first array element (arr[0]) or minimum possible integer instead of 0."
			diag.DetectionMethod = MethodStatic
			diag.ConfidenceScore = 0.88
		}
	}
}

func (d *Detector) layer3HistoricalPatterns(input DetectionInput, diag *Diagnosis) {
	if len(input.ActivePatterns) == 0 {
		return
	}

	for _, pattern := range input.ActivePatterns {
		// If the current diagnosis matches a recurring pattern for this student
		if pattern.Category == diag.Category && pattern.Subtype == diag.Subtype {
			diag.Evidence["recurring_mistake_detected"] = true
			diag.Evidence["previous_mistake_count"] = pattern.OccurrenceCount
			diag.Evidence["first_seen_at"] = pattern.FirstOccurrenceAt

			// Boost confidence since this student has made this exact mistake multiple times
			if diag.ConfidenceScore < 0.95 {
				diag.ConfidenceScore += 0.10
				if diag.ConfidenceScore > 1.0 {
					diag.ConfidenceScore = 1.0
				}
			}

			diag.Description = fmt.Sprintf("%s (Note: You have encountered this '%s' issue %d time(s) previously).",
				diag.Description, pattern.Subtype, pattern.OccurrenceCount)
			break
		}
	}
}

func (d *Detector) layer4RootCause(input DetectionInput, diag *Diagnosis) {
	code := input.Submission.SourceCode
	codeLower := strings.ToLower(code)

	// Binary search boundary misconception
	if (strings.Contains(codeLower, "binary") || strings.Contains(strings.ToLower(input.ProblemTopic), "binary search") ||
		strings.Contains(codeLower, "low <= high") || strings.Contains(codeLower, "left <= right")) {

		if strings.Contains(code, "low = mid") && !strings.Contains(code, "low = mid + 1") {
			diag.Evidence["root_cause_identified"] = "binary_search_infinite_loop_boundary"
			diag.Category = CategoryAlgorithmic
			diag.Subtype = "binary_search_boundary"
			diag.RelatedConcept = "Binary Search Invariant Maintenance"
			diag.Description = "Binary search boundary updated to 'mid' instead of 'mid + 1', causing an infinite loop when low and high differ by 1."
			diag.SuggestedCorrection = "Ensure boundaries strictly shrink each iteration: use 'low = mid + 1' or 'high = mid - 1'."
			diag.DetectionMethod = MethodRootCause
			diag.ConfidenceScore = 0.90
			return
		}
	}

	// Sliding window shrinking misconception
	if strings.Contains(codeLower, "window") || strings.Contains(strings.ToLower(input.ProblemTopic), "sliding window") {
		if strings.Contains(codeLower, "while") && strings.Contains(codeLower, "right") && !strings.Contains(codeLower, "left++") && !strings.Contains(codeLower, "left += 1") {
			diag.Evidence["root_cause_identified"] = "sliding_window_stagnant_left_pointer"
			diag.Category = CategoryAlgorithmic
			diag.Subtype = "sliding_window_condition"
			diag.RelatedConcept = "Sliding Window Window-Shrinking Invariant"
			diag.Description = "Sliding window expands with right pointer but does not correctly advance left pointer to maintain window constraints."
			diag.SuggestedCorrection = "Shrink the window from the left while window condition is violated: update frequencies and advance left pointer."
			diag.DetectionMethod = MethodRootCause
			diag.ConfidenceScore = 0.85
			return
		}
	}

	// Insufficient evidence fallback if Category or Subtype are still unknown
	if diag.Category == "" || diag.Category == CategoryUnknown {
		diag.Category = CategoryUnknown
		diag.Subtype = "insufficient_evidence"
		diag.RelatedConcept = "General Problem Logic"
		diag.Description = "The test cases did not pass, but static and runtime evidence are insufficient to isolate a single root cause."
		diag.SuggestedCorrection = "Review the problem description, print intermediate variables with a sample test case, or consult the AI Tutor."
		diag.DetectionMethod = MethodDeterministic
		diag.ConfidenceScore = 0.20
	}
}
