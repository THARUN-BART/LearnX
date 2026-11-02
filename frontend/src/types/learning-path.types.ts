export type ProgressionState =
  | 'not_started'
  | 'ready_to_learn'
  | 'in_progress'
  | 'practicing'
  | 'needs_revision'
  | 'mastered'
  | 'locked_by_prerequisite';

export type ActionType =
  | 'learn_concept'
  | 'revise_concept'
  | 'solve_easy'
  | 'practice_medium'
  | 'attempt_hard'
  | 'review_mistake'
  | 'complete_prerequisite';

export interface PathNode {
  id: string;
  student_id: string;
  topic_id?: string;
  topic_name: string;
  concept_key: string;
  concept_name: string;
  node_order: number;
  state: ProgressionState;
  prerequisite_concept_keys: string[];
  readiness_score: number;
  recommended_difficulty: 'Easy' | 'Medium' | 'Hard';
  recommended_problem_id?: string;
  recommended_problem_name?: string;
  updated_at: string;
}

export interface AdaptivePath {
  id: string;
  student_id: string;
  current_topic_id?: string;
  current_concept_key: string;
  recommended_difficulty: 'Easy' | 'Medium' | 'Hard';
  overall_readiness_score: number;
  learning_stage: string;
  nodes: PathNode[];
  created_at: string;
  updated_at: string;
}

export interface LearningRecommendation {
  id: string;
  student_id: string;
  action_type: ActionType;
  concept_key: string;
  concept_name: string;
  topic_id?: string;
  topic_name?: string;
  problem_id?: string;
  problem_title?: string;
  problem_slug?: string;
  recommended_difficulty: 'Easy' | 'Medium' | 'Hard';
  headline: string;
  explanation: string;
  relevant_weakness?: string;
  evidence: Record<string, any>;
  is_active: boolean;
  completed_at?: string;
  created_at: string;
}

export interface PathProgressResponse {
  total_concepts: number;
  mastered_count: number;
  in_progress_count: number;
  needs_revision_count: number;
  ready_to_learn_count: number;
  locked_count: number;
  overall_progress_percentage: number;
  current_stage: string;
  last_evaluated_at: string;
}

export interface RecommendationHistoryResponse {
  history: LearningRecommendation[];
  total: number;
}
