export type Difficulty = 'Easy' | 'Medium' | 'Hard' | 'Unspecified';
export type ProgressStatus = 'Not Started' | 'Attempted' | 'Solved' | 'Needs Revision';

export interface Roadmap {
  id: string;
  name: string;
  slug: string;
  source_name: string;
  source_url: string;
  topic_count: number;
  problem_count: number;
  is_active: boolean;
}

export interface RoadmapTopic {
  id: string;
  roadmap_id: string;
  title: string;
  slug: string;
  description?: string;
  topic_order: number;
  subtopic_count?: number;
  problem_count: number;
  solved_count?: number;
  completion_percentage?: number;
}

export interface Subtopic {
  id: string;
  topic_id: string;
  title: string;
  slug: string;
  description?: string;
  subtopic_order: number;
  problem_count: number;
}

export interface A2ZProblem {
  id: string;
  title: string;
  slug: string;
  summary?: string;
  difficulty: Difficulty;
  external_platform?: string;
  external_problem_url: string;
  source_sheet?: string;
  source_problem_id?: string;
  problem_order: number;
  tags: string[];
  progress_status?: ProgressStatus;
  subtopic_id?: string;
  topic_id?: string;
  related?: A2ZProblem[];
}

export interface StudentProblemProgress {
  student_id: string;
  problem_id: string;
  status: ProgressStatus;
  attempt_count: number;
  successful_attempts: number;
  last_attempted_at?: string;
  solved_at?: string;
}

export interface ProgressSummary {
  roadmap_id: string;
  roadmap_name: string;
  total_problems: number;
  solved_problems: number;
  attempted_problems: number;
  pending_problems: number;
  needs_revision: number;
  completion_percentage: number;
}

export interface TopicProgressItem {
  topic_id: string;
  topic_name: string;
  total_problems: number;
  solved_problems: number;
  completion_percentage: number;
}

export interface AnalyticsDifficultyItem {
  difficulty: Difficulty;
  total_problems: number;
  solved_problems: number;
  attempted_problems: number;
  accuracy_rate: number;
}

export interface AnalyticsActivityDay {
  date: string;
  attempt_count: number;
  solved_count: number;
}

// Legacy UI compatibility types
export interface TestCase {
  id: string;
  input: string;
  expectedOutput: string;
  actualOutput?: string;
  passed?: boolean;
}

export interface Problem extends A2ZProblem {
  topicId?: string;
  description: string;
  constraints: string[];
  examples: { input: string; output: string; explanation?: string }[];
  hints: string[];
  starterCode: Record<string, string>;
  testCases: TestCase[];
  accuracyRate?: number;
}

export interface Topic extends RoadmapTopic {
  difficulty: Difficulty;
  estimatedHours: number;
  totalProblems: number;
  completedProblems: number;
  status: 'locked' | 'in_progress' | 'completed';
  prerequisites: string[];
  subtopics: string[];
  explanation?: string;
}

export interface CodeExecutionRequest {
  problemId: string;
  language: string;
  code: string;
  customInput?: string;
}

export interface SubmissionTestResult {
  test_case_order: number;
  status: string;
  execution_time_ms: number;
  memory_used_kb: number;
  is_hidden: boolean;
  actual_output?: string;
  error_message?: string;
}

export interface CodeExecutionResult {
  status: string;
  runtimeMs: number;
  memoryMb: number;
  stdout?: string;
  stderr?: string;
  passedTests: number;
  totalTests: number;
  submissionId?: string;
  diagnosis?: MistakeDiagnosis;
  results?: SubmissionTestResult[];
}

export interface PedagogyAnalysis {
  runner_result?: {
    passed?: boolean;
    error_type?: string;
    output?: string;
    traceback?: string;
  };
  diagnosis?: {
    skill?: string;
    misconception_id?: string;
    misconception?: string;
    confidence?: number;
    evidence?: string;
  };
  hint?: {
    hint_level?: string;
    hint?: string;
  };
}

export interface SubmissionHistoryItem {
  id?: string;
  submission_id?: string;
  problem_id?: string;
  problem_title?: string;
  status: string;
  language?: string;
  execution_time_ms?: number;
  memory_used_kb?: number;
  submitted_at?: string;
  source_code?: string;
}

export interface AIConversation {
  id: string;
  title: string;
  problem_id?: string;
  conversation_type?: string;
  updated_at?: string;
  created_at?: string;
}

export interface StudentMistakePattern {
  id: string;
  category: string;
  subtype: string;
  related_concept: string;
  occurrence_count: number;
  correction_status: 'active' | 'corrected';
  last_occurrence_at: string;
}

export interface MistakeSummary {
  total_mistakes: number;
  category_counts: Record<string, number>;
  recurring_patterns: StudentMistakePattern[];
  needs_attention_concepts: string[];
}

export interface MistakeDiagnosis {
  id: string;
  student_id: string;
  submission_id: string;
  problem_id: string;
  category: string;
  subtype: string;
  description: string;
  related_concept: string;
  detection_method: string;
  confidence_level: 'High' | 'Medium' | 'Low' | 'Insufficient Evidence';
  confidence_score: number;
  evidence: Record<string, any>;
  suggested_correction: string;
  created_at: string;
}

export interface ProgressiveHint {
  problem_id: string;
  level: number;
  title: string;
  content: string;
  guiding_question?: string;
  next_level_available: boolean;
  created_at?: string;
}

export interface DryRunStep {
  step_number: number;
  line_number: number;
  action: string;
  variables: Record<string, any>;
  explanation: string;
}

export interface DryRunResult {
  test_case_input: string;
  steps: DryRunStep[];
  final_output: string;
}

export interface ComplexityResult {
  time_complexity: string;
  time_explanation: string;
  space_complexity: string;
  space_explanation: string;
  is_optimal: boolean;
  optimisation_notes?: string;
}

export interface ExplainMistakeResult {
  submission_id: string;
  category: string;
  subtype: string;
  what_went_wrong: string;
  why_it_happened: string;
  how_to_fix: string;
  related_dsa_concepts: string[];
}

export type AssistantMode = 'hint' | 'explain' | 'debug';

export interface AssistantMessage {
  id: string;
  role: 'user' | 'assistant' | 'system';
  content: string;
  codeSnippet?: string;
  timestamp: string;
  mode?: AssistantMode;
  conversationId?: string;
}

export type MistakeCategory =
  | 'logical'
  | 'syntax'
  | 'time_complexity'
  | 'space_complexity'
  | 'edge_cases'
  | 'algorithm_selection';

export interface MistakeRecord {
  id: string;
  problemId: string;
  problemTitle: string;
  topic: string;
  category: MistakeCategory;
  categoryLabel: string;
  frequency: number;
  date: string;
  failedCodeSnippet: string;
  explanation: string;
  whyFailed: string;
  suggestedCorrection: string;
  recommendedTopicId: string;
  recommendedTopicTitle: string;
}

export interface AnalyticsSummary {
  totalSolved: number;
  totalProblems: number;
  attemptedProblems: number;
  pendingProblems: number;
  needsRevision: number;
  completionPercentage: number;
  easySolved: number;
  mediumSolved: number;
  hardSolved: number;
  streakDays: number;
  hoursSpent: number;
  accuracyRate: number;
  weeklyActivity: { day: string; solved: number; hours: number }[];
  dailyActivity: { date: string; solved: number; attempts: number }[];
  topicProficiency: { topic: string; topicId?: string; score: number; total: number; solved: number }[];
  difficultyDistribution: { name: string; count: number; color: string }[];
  monthlyProgress: { month: string; solved: number }[];
  strongestTopic?: { topic: string; score: number };
  weakestTopic?: { topic: string; score: number };
}
