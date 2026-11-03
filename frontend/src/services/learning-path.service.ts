import { apiClient } from './api';
import {
  AdaptivePath,
  LearningRecommendation,
  PathProgressResponse,
  RecommendationHistoryResponse,
} from '../types/learning-path.types';

export const learningPathService = {
  getLearningPath: async (): Promise<AdaptivePath> => {
    const res = await apiClient.get<AdaptivePath>('/api/v1/student/learning-path');
    return res.data;
  },

  getNextAction: async (): Promise<LearningRecommendation> => {
    const res = await apiClient.get<LearningRecommendation>('/api/v1/student/learning-path/next-action');
    return res.data;
  },

  getProgress: async (): Promise<PathProgressResponse> => {
    const res = await apiClient.get<PathProgressResponse>('/api/v1/student/learning-path/progress');
    return res.data;
  },

  getHistory: async (limit = 10, offset = 0): Promise<RecommendationHistoryResponse> => {
    const res = await apiClient.get<RecommendationHistoryResponse>('/api/v1/student/learning-path/history', {
      params: { limit, offset },
    });
    return res.data;
  },

  recalculate: async (): Promise<AdaptivePath> => {
    const res = await apiClient.post<AdaptivePath>('/api/v1/student/learning-path/recalculate', {});
    return res.data;
  },
};
