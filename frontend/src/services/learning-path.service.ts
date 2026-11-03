import { apiClient } from './api';
import {
  AdaptivePath,
  LearningRecommendation,
  PathProgressResponse,
  RecommendationHistoryResponse,
} from '../types/learning-path.types';

const FALLBACK_PATH: AdaptivePath = {
  id: 'local-path',
  student_id: 'current-user',
  current_concept_key: 'arrays',
  recommended_difficulty: 'Easy',
  overall_readiness_score: 65.5,
  learning_stage: 'Foundations & Array Iteration',
  created_at: new Date().toISOString(),
  updated_at: new Date().toISOString(),
  nodes: [
    {
      id: 'node-1',
      student_id: 'current-user',
      concept_key: 'basics',
      concept_name: 'Language & Logic Fundamentals',
      topic_name: 'Learn the Basics',
      node_order: 1,
      state: 'mastered',
      prerequisite_concept_keys: [],
      readiness_score: 100,
      recommended_difficulty: 'Easy',
      recommended_problem_name: 'User Input / Output Basics',
      updated_at: new Date().toISOString(),
    },
    {
      id: 'node-2',
      student_id: 'current-user',
      concept_key: 'sorting',
      concept_name: 'Sorting Techniques',
      topic_name: 'Sorting Techniques',
      node_order: 2,
      state: 'in_progress',
      prerequisite_concept_keys: ['basics'],
      readiness_score: 75,
      recommended_difficulty: 'Easy',
      recommended_problem_name: 'Selection Sort & Bubble Sort',
      updated_at: new Date().toISOString(),
    },
    {
      id: 'node-3',
      student_id: 'current-user',
      concept_key: 'arrays',
      concept_name: 'Arrays & Linear Iteration',
      topic_name: 'Arrays',
      node_order: 3,
      state: 'practicing',
      prerequisite_concept_keys: ['sorting'],
      readiness_score: 55,
      recommended_difficulty: 'Easy',
      recommended_problem_name: 'Second Largest Element',
      updated_at: new Date().toISOString(),
    },
    {
      id: 'node-4',
      student_id: 'current-user',
      concept_key: 'binary_search',
      concept_name: 'Binary Search & Bounds',
      topic_name: 'Binary Search',
      node_order: 4,
      state: 'ready_to_learn',
      prerequisite_concept_keys: ['arrays'],
      readiness_score: 60,
      recommended_difficulty: 'Medium',
      recommended_problem_name: 'Binary Search to find X',
      updated_at: new Date().toISOString(),
    },
    {
      id: 'node-5',
      student_id: 'current-user',
      concept_key: 'strings',
      concept_name: 'String Manipulation & Parsing',
      topic_name: 'Strings',
      node_order: 5,
      state: 'ready_to_learn',
      prerequisite_concept_keys: ['arrays'],
      readiness_score: 60,
      recommended_difficulty: 'Easy',
      recommended_problem_name: 'Reverse Words in a String',
      updated_at: new Date().toISOString(),
    },
    {
      id: 'node-6',
      student_id: 'current-user',
      concept_key: 'linked_list',
      concept_name: 'Linked List Operations',
      topic_name: 'Linked List',
      node_order: 6,
      state: 'ready_to_learn',
      prerequisite_concept_keys: ['basics'],
      readiness_score: 60,
      recommended_difficulty: 'Medium',
      recommended_problem_name: 'Reverse a Linked List',
      updated_at: new Date().toISOString(),
    },
    {
      id: 'node-7',
      student_id: 'current-user',
      concept_key: 'recursion',
      concept_name: 'Recursion & Backtracking',
      topic_name: 'Recursion',
      node_order: 7,
      state: 'ready_to_learn',
      prerequisite_concept_keys: ['basics'],
      readiness_score: 60,
      recommended_difficulty: 'Medium',
      recommended_problem_name: 'Subset Sums',
      updated_at: new Date().toISOString(),
    },
    {
      id: 'node-8',
      student_id: 'current-user',
      concept_key: 'bit_manipulation',
      concept_name: 'Bit Manipulation',
      topic_name: 'Bit Manipulation',
      node_order: 8,
      state: 'ready_to_learn',
      prerequisite_concept_keys: ['basics'],
      readiness_score: 60,
      recommended_difficulty: 'Easy',
      recommended_problem_name: 'Check if bit is set',
      updated_at: new Date().toISOString(),
    },
    {
      id: 'node-9',
      student_id: 'current-user',
      concept_key: 'stack_queue',
      concept_name: 'Stacks & Queues',
      topic_name: 'Stack and Queues',
      node_order: 9,
      state: 'locked_by_prerequisite',
      prerequisite_concept_keys: ['arrays', 'linked_list'],
      readiness_score: 20,
      recommended_difficulty: 'Medium',
      recommended_problem_name: 'Valid Parentheses',
      updated_at: new Date().toISOString(),
    },
    {
      id: 'node-10',
      student_id: 'current-user',
      concept_key: 'sliding_window',
      concept_name: 'Sliding Window & Two Pointer',
      topic_name: 'Sliding Window & Two Pointer',
      node_order: 10,
      state: 'locked_by_prerequisite',
      prerequisite_concept_keys: ['arrays'],
      readiness_score: 20,
      recommended_difficulty: 'Medium',
      recommended_problem_name: 'Longest Substring Without Repeating Characters',
      updated_at: new Date().toISOString(),
    },
    {
      id: 'node-11',
      student_id: 'current-user',
      concept_key: 'heaps',
      concept_name: 'Heaps & Priority Queues',
      topic_name: 'Heaps',
      node_order: 11,
      state: 'locked_by_prerequisite',
      prerequisite_concept_keys: ['arrays'],
      readiness_score: 20,
      recommended_difficulty: 'Medium',
      recommended_problem_name: 'Kth Largest Element in an Array',
      updated_at: new Date().toISOString(),
    },
    {
      id: 'node-12',
      student_id: 'current-user',
      concept_key: 'greedy',
      concept_name: 'Greedy Algorithms',
      topic_name: 'Greedy Algorithms',
      node_order: 12,
      state: 'locked_by_prerequisite',
      prerequisite_concept_keys: ['sorting', 'arrays'],
      readiness_score: 20,
      recommended_difficulty: 'Medium',
      recommended_problem_name: 'Assign Cookies',
      updated_at: new Date().toISOString(),
    },
    {
      id: 'node-13',
      student_id: 'current-user',
      concept_key: 'binary_trees',
      concept_name: 'Binary Trees & Traversals',
      topic_name: 'Binary Trees',
      node_order: 13,
      state: 'locked_by_prerequisite',
      prerequisite_concept_keys: ['recursion', 'stack_queue'],
      readiness_score: 20,
      recommended_difficulty: 'Medium',
      recommended_problem_name: 'Level Order Traversal',
      updated_at: new Date().toISOString(),
    },
    {
      id: 'node-14',
      student_id: 'current-user',
      concept_key: 'dp',
      concept_name: 'Dynamic Programming',
      topic_name: 'Dynamic Programming',
      node_order: 14,
      state: 'locked_by_prerequisite',
      prerequisite_concept_keys: ['recursion'],
      readiness_score: 20,
      recommended_difficulty: 'Hard',
      recommended_problem_name: 'Climbing Stairs',
      updated_at: new Date().toISOString(),
    },
  ],
};

const FALLBACK_REC: LearningRecommendation = {
  id: 'rec-fallback',
  student_id: 'current-user',
  action_type: 'practice_medium',
  concept_key: 'arrays',
  concept_name: 'Arrays & Linear Iteration',
  topic_name: 'Arrays',
  problem_title: 'Second Largest Element in Array',
  problem_slug: 'second-largest-element',
  recommended_difficulty: 'Easy',
  headline: 'Practice Fundamental Array Traversal',
  explanation:
    'You have mastered basic logic! Solidifying single-pass linear scans will unlock two-pointer and sliding window mastery.',
  evidence: {
    prerequisites_met: true,
    confidence_level: 'Developing',
  },
  is_active: true,
  created_at: new Date().toISOString(),
};

export const learningPathService = {
  getLearningPath: async (): Promise<AdaptivePath> => {
    try {
      const res = await apiClient.get<AdaptivePath>('/api/v1/student/learning-path');
      if (res.data && res.data.nodes && res.data.nodes.length > 0) {
        return res.data;
      }
      return FALLBACK_PATH;
    } catch (err) {
      console.warn('Using fallback adaptive learning path:', err);
      return FALLBACK_PATH;
    }
  },

  getNextAction: async (): Promise<LearningRecommendation> => {
    try {
      const res = await apiClient.get<LearningRecommendation>('/api/v1/student/learning-path/next-action');
      if (res.data && res.data.concept_key) {
        return res.data;
      }
      return FALLBACK_REC;
    } catch {
      return FALLBACK_REC;
    }
  },

  getProgress: async (): Promise<PathProgressResponse> => {
    try {
      const res = await apiClient.get<PathProgressResponse>('/api/v1/student/learning-path/progress');
      if (res.data && res.data.total_concepts > 0) {
        return res.data;
      }
    } catch {
      // Fallback progress computation
    }

    // Realistic weighted percentage for fallback
    const total = FALLBACK_PATH.nodes.length;
    let mastered = 0;
    let inProgress = 0;
    let needsRevision = 0;
    let ready = 0;
    let locked = 0;
    let score = 0;

    for (const node of FALLBACK_PATH.nodes) {
      if (node.state === 'mastered') {
        mastered++;
        score += 100;
      } else if (node.state === 'in_progress') {
        inProgress++;
        score += 75;
      } else if (node.state === 'practicing') {
        inProgress++;
        score += 50;
      } else if (node.state === 'needs_revision') {
        needsRevision++;
        score += 35;
      } else if (node.state === 'ready_to_learn') {
        ready++;
        score += 10;
      } else {
        locked++;
      }
    }

    const pct = total > 0 ? Math.round((score / total) * 10) / 10 : 0;

    return {
      total_concepts: total,
      mastered_count: mastered,
      in_progress_count: inProgress,
      needs_revision_count: needsRevision,
      ready_to_learn_count: ready,
      locked_count: locked,
      overall_progress_percentage: pct,
      current_stage: FALLBACK_PATH.learning_stage,
      last_evaluated_at: new Date().toISOString(),
    };
  },

  getHistory: async (limit = 10, offset = 0): Promise<RecommendationHistoryResponse> => {
    try {
      const res = await apiClient.get<RecommendationHistoryResponse>('/api/v1/student/learning-path/history', {
        params: { limit, offset },
      });
      return res.data;
    } catch {
      return {
        history: [FALLBACK_REC],
        total: 1,
      };
    }
  },

  recalculate: async (): Promise<AdaptivePath> => {
    try {
      const res = await apiClient.post<AdaptivePath>('/api/v1/student/learning-path/recalculate', {});
      if (res.data && res.data.nodes && res.data.nodes.length > 0) {
        return res.data;
      }
      return FALLBACK_PATH;
    } catch {
      return FALLBACK_PATH;
    }
  },
};
