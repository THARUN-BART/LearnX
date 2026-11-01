package testcases

import "time"

type TestCase struct {
	ID             string    `json:"id"`
	ProblemID      string    `json:"problem_id"`
	InputData      string    `json:"input_data"`
	ExpectedOutput string    `json:"expected_output"`
	IsHidden       bool      `json:"is_hidden"`
	TestOrder      int       `json:"test_order"`
	CreatedAt      time.Time `json:"created_at"`
}

// PublicTestCase is safe for students and never exposes hidden inputs/outputs.
type PublicTestCase struct {
	ID             string `json:"id"`
	ProblemID      string `json:"problem_id"`
	InputData      string `json:"input_data,omitempty"`
	ExpectedOutput string `json:"expected_output,omitempty"`
	IsHidden       bool   `json:"is_hidden"`
	TestOrder      int    `json:"test_order"`
}

type CreateTestCaseRequest struct {
	InputData      string `json:"input_data"`
	ExpectedOutput string `json:"expected_output"`
	IsHidden       bool   `json:"is_hidden"`
	TestOrder      int    `json:"test_order"`
}
