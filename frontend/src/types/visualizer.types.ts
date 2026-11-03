export type AlgorithmCategory =
  | 'Sorting'
  | 'Searching'
  | 'Arrays and Patterns'
  | 'Linked Lists'
  | 'Stack and Queue'
  | 'Recursion and Backtracking'
  | 'Trees'
  | 'Graphs'
  | 'Dynamic Programming';

export interface ComplexityInfo {
  best: string;
  average: string;
  worst: string;
  space: string;
}

export interface AlgorithmMetadata {
  id: string;
  name: string;
  category: AlgorithmCategory;
  description: string;
  problem_solved: string;
  complexity: ComplexityInfo;
  use_cases: string[];
  tags: string[];
  input_type: string;
  default_input: Record<string, any>;
}

export interface StepState {
  array?: number[];
  pointers?: Record<string, number>;
  comparing?: number[];
  swapping?: number[];
  highlighted?: number[];
  secondary?: number[];
  matrix?: number[][];
  active_cell?: [number, number];
  status_text?: string;
  nodes?: any[];
  edges?: any[];
}

export interface ExecutionStep {
  step_number: number;
  operation: string;
  description: string;
  explanation: string;
  state: StepState;
  is_complete: boolean;
}

export interface ExecuteResponse {
  algorithm_id: string;
  name: string;
  category: AlgorithmCategory;
  complexity: ComplexityInfo;
  total_steps: number;
  steps: ExecutionStep[];
}

export interface ExplainStepResponse {
  explanation: string;
  intuition: string;
  next_expected_move: string;
  key_takeaway: string;
}

export interface PracticeQuestion {
  id: string;
  algorithm_id: string;
  question: string;
  options: string[];
  correct_option_index: number;
  explanation: string;
  difficulty: 'Easy' | 'Medium' | 'Hard';
}

export interface VisualizerSession {
  id: string;
  student_id: string;
  algorithm_id: string;
  category: AlgorithmCategory;
  input_config: Record<string, any>;
  steps_completed: number;
  total_steps: number;
  time_spent_seconds: number;
  completed: boolean;
  created_at: string;
  updated_at: string;
}

