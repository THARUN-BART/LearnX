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
  description: string;
  constraints: string[];
  examples: { input: string; output: string; explanation?: string }[];
  hints: string[];
  starterCode: Record<string, string>;
  testCases: TestCase[];
  accuracyRate?: number;
}

export interface Topic extends RoadmapTopic {
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

export interface CodeExecutionResult {
  status: 'Accepted' | 'Wrong Answer' | 'Time Limit Exceeded' | 'Compilation Error' | 'Runtime Error';
  runtimeMs: number;
  memoryMb: number;
  stdout?: string;
  stderr?: string;
  passedTests: number;
  totalTests: number;
}

export type AssistantMode = 'hint' | 'explain' | 'debug';

export interface AssistantMessage {
  id: string;
  role: 'user' | 'assistant' | 'system';
  content: string;
  codeSnippet?: string;
  timestamp: string;
  mode?: AssistantMode;
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
  easySolved: number;
  mediumSolved: number;
  hardSolved: number;
  streakDays: number;
  hoursSpent: number;
  accuracyRate: number;
  weeklyActivity: { day: string; solved: number; hours: number }[];
  topicProficiency: { topic: string; score: number; total: number }[];
  difficultyDistribution: { name: string; count: number; color: string }[];
  monthlyProgress: { month: string; solved: number }[];
}
