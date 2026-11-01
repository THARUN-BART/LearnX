import apiClient from './api';
import { Problem } from '../types/learning.types';
import { SAMPLE_PROBLEMS } from '../utils/constants';

class LearningService {
  async getProblems(topicId?: string): Promise<Problem[]> {
    try {
      const response = await apiClient.get('/api/v1/problems', {
        params: topicId ? { topicId } : {}
      });
      return response.data.data || response.data;
    } catch {
      if (topicId) {
        return SAMPLE_PROBLEMS.filter((p) => p.topicId === topicId);
      }
      return SAMPLE_PROBLEMS;
    }
  }

  async getProblemById(problemId: string): Promise<Problem | null> {
    try {
      const response = await apiClient.get(`/api/v1/problems/${problemId}`);
      return response.data.data || response.data;
    } catch {
      const prob = SAMPLE_PROBLEMS.find((p) => p.id === problemId || p.slug === problemId);
      return prob || SAMPLE_PROBLEMS[0];
    }
  }
}

export const learningService = new LearningService();
export default learningService;

