import apiClient from './api';
import {
  CategorizedRecommendations,
  ProblemRecommendation,
  RecommendationFeedbackPayload,
  RecommendationHistoryResponse
} from '../types/recommendations.types';
import { SAMPLE_PROBLEMS } from '../utils/constants';

class RecommendationService {
  /**
   * Fetches all categorized recommendations for the current student.
   */
  async getRecommendations(): Promise<CategorizedRecommendations> {
    try {
      const res = await apiClient.get('/api/v1/student/recommendations');
      return res.data?.data || res.data;
    } catch {
      return this.generateFallbackRecommendations();
    }
  }

  /**
   * Retrieves the #1 highest priority next recommendation.
   */
  async getNextRecommendation(): Promise<ProblemRecommendation | null> {
    try {
      const res = await apiClient.get('/api/v1/student/recommendations/next');
      return res.data?.data || res.data;
    } catch {
      const fallback = this.generateFallbackRecommendations();
      return fallback.top_pick || null;
    }
  }

  /**
   * Retrieves recommendations specifically addressing active weaknesses and mistake patterns.
   */
  async getWeaknessRecommendations(): Promise<ProblemRecommendation[]> {
    try {
      const res = await apiClient.get('/api/v1/student/recommendations/weaknesses');
      const data = res.data?.data || res.data;
      return data?.weakness_recommendations || (Array.isArray(data) ? data : []);
    } catch {
      const fallback = this.generateFallbackRecommendations();
      return fallback.categories['improve_weaknesses'] || [];
    }
  }

  /**
   * Retrieves similar problems that reinforce the same conceptual pattern.
   */
  async getSimilarProblems(problemId: string): Promise<ProblemRecommendation[]> {
    try {
      const res = await apiClient.get(`/api/v1/student/recommendations/similar/${problemId}`);
      const data = res.data?.data || res.data;
      return data?.similar_problems || (Array.isArray(data) ? data : []);
    } catch {
      return this.generateFallbackSimilar(problemId);
    }
  }

  /**
   * Forces dynamic re-evaluation and returns refreshed recommendations.
   */
  async refreshRecommendations(): Promise<CategorizedRecommendations> {
    try {
      const res = await apiClient.post('/api/v1/student/recommendations/refresh');
      return res.data?.data || res.data;
    } catch {
      return this.generateFallbackRecommendations();
    }
  }

  /**
   * Records student feedback for a recommended problem.
   */
  async submitFeedback(
    recommendationId: string,
    payload: RecommendationFeedbackPayload
  ): Promise<void> {
    try {
      await apiClient.post(`/api/v1/student/recommendations/${recommendationId}/feedback`, payload);
    } catch (e) {
      console.warn('[RecommendationService] Could not submit feedback to backend:', e);
    }
  }

  /**
   * Fetches historical log of recommendations.
   */
  async getHistory(limit = 20, offset = 0): Promise<RecommendationHistoryResponse> {
    try {
      const res = await apiClient.get('/api/v1/student/recommendations/history', {
        params: { limit, offset }
      });
      return res.data?.data || res.data;
    } catch {
      return {
        history: [],
        total: 0,
        limit,
        offset
      };
    }
  }

  // -------------------------------------------------------------------------
  // Offline / Dev Fallback Generator
  // -------------------------------------------------------------------------
  private generateFallbackRecommendations(): CategorizedRecommendations {
    const p1 = SAMPLE_PROBLEMS[0];
    const p2 = SAMPLE_PROBLEMS[1] || p1;
    const p3 = SAMPLE_PROBLEMS[2] || p1;

    const topPick: ProblemRecommendation = {
      id: 'rec-fallback-top',
      student_id: 'usr_current',
      problem_id: p2.id,
      category: 'recommended_for_you',
      recommended_difficulty: p2.difficulty,
      primary_concept: 'Arrays',
      related_concepts: p2.tags,
      match_score: 96.0,
      headline: `Top Priority: ${p2.title}`,
      explanation:
        'Your recent submissions showed strong grasp of linear traversal. Solving this problem will reinforce two-pointer optimization and boundary handling.',
      learning_objective: 'Master two-pointer movement and boundary invariant verification.',
      status: 'active',
      is_top_pick: true,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      problem_title: p2.title,
      problem_slug: p2.slug,
      problem_difficulty: p2.difficulty,
      problem_summary: p2.summary,
      external_problem_url: p2.external_problem_url,
      topic_name: 'Arrays',
      tags: p2.tags,
      attempt_status: 'Not Started'
    };

    const weaknessRec: ProblemRecommendation = {
      id: 'rec-fallback-weak',
      student_id: 'usr_current',
      problem_id: p1.id,
      category: 'improve_weaknesses',
      recommended_difficulty: p1.difficulty,
      primary_concept: 'Arrays & Two Pointers',
      related_concepts: p1.tags,
      match_score: 92.0,
      headline: `Overcome Weakness: ${p1.title}`,
      explanation:
        'Targeted practice to address recurring edge-case errors with index bounds and duplicate values.',
      learning_objective: 'Eliminate off-by-one errors and handle identical element edge cases.',
      status: 'active',
      is_top_pick: false,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      problem_title: p1.title,
      problem_slug: p1.slug,
      problem_difficulty: p1.difficulty,
      problem_summary: p1.summary,
      external_problem_url: p1.external_problem_url,
      topic_name: 'Arrays',
      tags: p1.tags,
      attempt_status: 'Not Started'
    };

    const fundRec: ProblemRecommendation = {
      id: 'rec-fallback-fund',
      student_id: 'usr_current',
      problem_id: p1.id,
      category: 'strengthen_fundamentals',
      recommended_difficulty: 'Easy',
      primary_concept: 'Linear Iteration',
      related_concepts: ['Arrays'],
      match_score: 85.0,
      headline: `Strengthen Foundations: ${p1.title}`,
      explanation:
        'Solidify standard iteration logic to ensure clean, bug-free implementations under time pressure.',
      learning_objective: 'Achieve zero-compile-error speed and clean memory bounds.',
      status: 'active',
      is_top_pick: false,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      problem_title: p1.title,
      problem_slug: p1.slug,
      problem_difficulty: 'Easy',
      problem_summary: p1.summary,
      external_problem_url: p1.external_problem_url,
      topic_name: 'Arrays',
      tags: ['Arrays'],
      attempt_status: 'Not Started'
    };

    const challengeRec: ProblemRecommendation = {
      id: 'rec-fallback-chal',
      student_id: 'usr_current',
      problem_id: p3.id,
      category: 'challenge_yourself',
      recommended_difficulty: 'Hard',
      primary_concept: 'Two Pointers & Monotonic Stack',
      related_concepts: p3.tags,
      match_score: 88.0,
      headline: `Push Limits: ${p3.title}`,
      explanation:
        'Test your problem-solving depth with strict O(N) runtime and O(1) space constraints.',
      learning_objective: 'Develop rigorous edge-case analysis on multi-pointer problems.',
      status: 'active',
      is_top_pick: false,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      problem_title: p3.title,
      problem_slug: p3.slug,
      problem_difficulty: 'Hard',
      problem_summary: p3.summary,
      external_problem_url: p3.external_problem_url,
      topic_name: 'Arrays',
      tags: p3.tags,
      attempt_status: 'Not Started'
    };

    return {
      top_pick: topPick,
      categories: {
        recommended_for_you: [topPick],
        improve_weaknesses: [weaknessRec],
        strengthen_fundamentals: [fundRec],
        challenge_yourself: [challengeRec],
        revision_practice: [weaknessRec]
      },
      total_recommendations: 4,
      recommended_difficulty: 'Medium',
      primary_focus_concept: 'Arrays & Two Pointers',
      generated_at: new Date().toISOString()
    };
  }

  private generateFallbackSimilar(problemId: string): ProblemRecommendation[] {
    return SAMPLE_PROBLEMS.filter((p) => p.id !== problemId).map((p, idx) => ({
      id: `sim-rec-${idx}`,
      student_id: 'usr_current',
      problem_id: p.id,
      category: 'similar_problems',
      recommended_difficulty: p.difficulty,
      primary_concept: 'DSA Curriculum',
      related_concepts: p.tags,
      match_score: 80.0,
      headline: `Similar Pattern: ${p.title}`,
      explanation:
        'Applies comparable algorithmic techniques with slight variations in constraint bounds.',
      learning_objective: 'Reinforce the algorithmic pattern across diverse constraints.',
      status: 'active',
      is_top_pick: false,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      problem_title: p.title,
      problem_slug: p.slug,
      problem_difficulty: p.difficulty,
      problem_summary: p.summary,
      external_problem_url: p.external_problem_url,
      topic_name: 'Arrays',
      tags: p.tags,
      attempt_status: 'Not Started'
    }));
  }
}

export const recommendationService = new RecommendationService();
export default recommendationService;

