package visualizer

import (
	"time"
)

type AlgorithmCategory string

const (
	CategorySorting               AlgorithmCategory = "Sorting"
	CategorySearching             AlgorithmCategory = "Searching"
	CategoryArraysPatterns        AlgorithmCategory = "Arrays and Patterns"
	CategoryLinkedLists           AlgorithmCategory = "Linked Lists"
	CategoryStackQueue            AlgorithmCategory = "Stack and Queue"
	CategoryRecursionBacktracking AlgorithmCategory = "Recursion and Backtracking"
	CategoryTrees                 AlgorithmCategory = "Trees"
	CategoryGraphs                AlgorithmCategory = "Graphs"
	CategoryDP                    AlgorithmCategory = "Dynamic Programming"
)

type ComplexityInfo struct {
	Best    string `json:"best"`
	Average string `json:"average"`
	Worst   string `json:"worst"`
	Space   string `json:"space"`
}

type AlgorithmMetadata struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Category      AlgorithmCategory `json:"category"`
	Description   string            `json:"description"`
	ProblemSolved string            `json:"problem_solved"`
	Complexity    ComplexityInfo    `json:"complexity"`
	UseCases      []string          `json:"use_cases"`
	Tags          []string          `json:"tags"`
	InputType     string            `json:"input_type"` // array, array_target, linked_list, stack_queue, tree, graph, number, n_queens, dp_string
	DefaultInput  map[string]any    `json:"default_input"`
}

type StepState struct {
	Array       []int          `json:"array,omitempty"`
	Pointers    map[string]int `json:"pointers,omitempty"`
	Comparing   []int          `json:"comparing,omitempty"`
	Swapping    []int          `json:"swapping,omitempty"`
	Highlighted []int          `json:"highlighted,omitempty"`
	Secondary   []int          `json:"secondary,omitempty"`
	PointersMap map[string]any `json:"pointers_map,omitempty"`
	Nodes       []any          `json:"nodes,omitempty"`
	Edges       []any          `json:"edges,omitempty"`
	Matrix      [][]int        `json:"matrix,omitempty"`
	ActiveCell  []int          `json:"active_cell,omitempty"` // [row, col]
	Auxiliary   any            `json:"auxiliary,omitempty"`
	StatusText  string         `json:"status_text,omitempty"`
}

type ExecutionStep struct {
	StepNumber  int       `json:"step_number"`
	Operation   string    `json:"operation"` // compare, swap, move_pointer, visit, insert, delete, backtrack, update_dp, result, complete
	Description string    `json:"description"`
	Explanation string    `json:"explanation"`
	State       StepState `json:"state"`
	IsComplete  bool      `json:"is_complete"`
}

type ExecuteRequest struct {
	AlgorithmID string         `json:"algorithm_id"`
	CustomInput map[string]any `json:"custom_input"`
}

type ExecuteResponse struct {
	AlgorithmID string            `json:"algorithm_id"`
	Name        string            `json:"name"`
	Category    AlgorithmCategory `json:"category"`
	Complexity  ComplexityInfo    `json:"complexity"`
	TotalSteps  int               `json:"total_steps"`
	Steps       []ExecutionStep   `json:"steps"`
}

type VisualizerSession struct {
	ID               string            `json:"id"`
	StudentID        string            `json:"student_id"`
	AlgorithmID      string            `json:"algorithm_id"`
	Category         AlgorithmCategory `json:"category"`
	InputConfig      map[string]any    `json:"input_config"`
	StepsCompleted   int               `json:"steps_completed"`
	TotalSteps       int               `json:"total_steps"`
	TimeSpentSeconds int               `json:"time_spent_seconds"`
	Completed        bool              `json:"completed"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
}

type CreateSessionRequest struct {
	AlgorithmID string         `json:"algorithm_id"`
	Category    string         `json:"category"`
	InputConfig map[string]any `json:"input_config"`
	TotalSteps  int            `json:"total_steps"`
}

type UpdateSessionRequest struct {
	StepsCompleted   int  `json:"steps_completed"`
	TimeSpentSeconds int  `json:"time_spent_seconds"`
	Completed        bool `json:"completed"`
}

type ExplainStepRequest struct {
	AlgorithmID  string         `json:"algorithm_id"`
	StepNumber   int            `json:"step_number"`
	TotalSteps   int            `json:"total_steps"`
	Operation    string         `json:"operation"`
	Explanation  string         `json:"explanation"`
	CurrentState map[string]any `json:"current_state"`
	UserQuestion string         `json:"user_question,omitempty"`
}

type ExplainStepResponse struct {
	Explanation      string `json:"explanation"`
	Intuition        string `json:"intuition"`
	NextExpectedMove string `json:"next_expected_move"`
	KeyTakeaway      string `json:"key_takeaway"`
}

type PracticeQuestion struct {
	ID                 string   `json:"id"`
	AlgorithmID        string   `json:"algorithm_id"`
	Question           string   `json:"question"`
	Options            []string `json:"options"`
	CorrectOptionIndex int      `json:"correct_option_index"`
	Explanation        string   `json:"explanation"`
	Difficulty         string   `json:"difficulty"`
}

