import { Difficulty, ProgressStatus } from './learning.types';

export type RecommendationCategory =
  | 'recommended_for_you'
  | 'improve_weaknesses'
  | 'strengthen_fundamentals'
  | 'challenge_yourself'
  | 'similar_problems'
  | 'revision_practice';

export type FeedbackRating = 'helpful' | 'not_relevant' | 'too_easy' | 'too_difficult';

export interface ProblemRecommendation {
  id: string;
  student_id: string;
  problem_id: string;
  category: RecommendationCategory;
  recommended_difficulty: Difficulty;
  primary_concept: string;
  related_concepts: string[];
  match_score: number;
  headline: string;
  explanation: string;
  learning_objective: string;
  status: 'active' | 'completed' | 'skipped' | 'expired';
  is_top_pick: boolean;
  feedback_rating?: FeedbackRating;
  feedback_notes?: string;
  completed_at?: string;
  created_at: string;
  updated_at: string;

  // Hydrated problem metadata
  problem_title?: string;
  problem_slug?: string;
  problem_difficulty?: Difficulty;
  problem_summary?: string;
  external_problem_url?: string;
  topic_name?: string;
  subtopic_name?: string;
  tags?: string[];
  attempt_status?: ProgressStatus;
}

export interface CategorizedRecommendations {
  top_pick?: ProblemRecommendation;
  categories: Record<string, ProblemRecommendation[]>;
  total_recommendations: number;
  recommended_difficulty: Difficulty;
  primary_focus_concept: string;
  generated_at: string;
}

export interface RecommendationFeedbackPayload {
  rating: FeedbackRating;
  notes?: string;
}

export interface RecommendationHistoryResponse {
  history: ProblemRecommendation[];
  total: number;
  limit: number;
  offset: number;
}

