export type IssueSeverity = 'critical' | 'warning' | 'info';

export interface CodeIssue {
  id: string;
  severity: IssueSeverity;
  line_number: number;
  title: string;
  description: string;
  recommendation: string;
  code_snippet?: string;
  category: string;
  source_type: 'static_analysis' | 'execution_observation' | 'ai_synthesis';
}

export interface ReviewScores {
  overall: number;
  readability: number;
  efficiency: number;
  maintainability: number;
}

export interface ComplexityComparison {
  original_time: string;
  optimized_time: string;
  original_space: string;
  optimized_space: string;
  explanation: string;
}

export interface AlternativeApproach {
  name: string;
  description: string;
  time: string;
  space: string;
  trade_off: string;
}

export interface CodeReviewResult {
  id: string;
  student_id: string;
  problem_id?: string;
  problem_title?: string;
  submission_id?: string;
  language: string;
  source_code: string;
  scores: ReviewScores;
  correctness_observations: string[];
  issues: CodeIssue[];
  optimized_code: string;
  explanation_of_improvements: string;
  complexity: ComplexityComparison;
  alternative_approaches: AlternativeApproach[];
  recommended_visualizer_algo_id?: string;
  recommended_visualizer_reason?: string;
  created_at: string;
}

export interface ReviewHistoryItem {
  id: string;
  problem_id?: string;
  problem_title?: string;
  language: string;
  overall_score: number;
  created_at: string;
}

export interface CodeComparisonData {
  original_code: string;
  optimized_code: string;
  summary: string;
  complexity: ComplexityComparison;
}

export interface CodeReviewRequest {
  problem_id?: string;
  submission_id?: string;
  language: string;
  source_code: string;
  mode?: string;
}
