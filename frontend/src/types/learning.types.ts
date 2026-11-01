export type Difficulty = 'Easy' | 'Medium' | 'Hard';

export interface TestCase {
  id: string;
  input: string;
  expectedOutput: string;
  actualOutput?: string;
  passed?: boolean;
  isHidden?: boolean;
}

export interface Example {
  input: string;
  output: string;
  explanation?: string;
}

export interface Problem {
  id: string;
  title: string;
  slug: string;
  topicId: string;
  difficulty: Difficulty;
  description: string;
  constraints: string[];
  examples: Example[];
  hints: string[];
  starterCode: Record<string, string>;
  testCases: TestCase[];
  accuracyRate?: number;
  solved?: boolean;
}

export interface Topic {
  id: string;
  title: string;
  slug: string;
  description: string;
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

export interface CodeExecutionResult {
  status: 'Accepted' | 'Wrong Answer' | 'Time Limit Exceeded' | 'Compilation Error' | 'Runtime Error';
  runtimeMs: number;
  memoryMb: number;
  stdout?: string;
  stderr?: string;
  passedTests: number;
  totalTests: number;
  details?: {
    testCaseId: string;
    input: string;
    expected: string;
    received: string;
  }[];
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

