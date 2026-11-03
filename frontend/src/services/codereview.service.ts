import apiClient, { unwrap } from './api';
import {
  CodeReviewRequest,
  CodeReviewResult,
  ReviewHistoryItem,
  CodeComparisonData,
} from '../types/codereview.types';

class CodeReviewService {
  async reviewCode(req: CodeReviewRequest): Promise<CodeReviewResult | null> {
    try {
      const res = await apiClient.post('/api/v1/ai/code-review', req);
      return unwrap<CodeReviewResult>(res.data, 'review') || unwrap<CodeReviewResult>(res.data, 'data') || null;
    } catch {
      return null;
    }
  }

  async getReviewById(id: string): Promise<CodeReviewResult | null> {
    try {
      const res = await apiClient.get(`/api/v1/ai/code-review/${id}`);
      return unwrap<CodeReviewResult>(res.data, 'review') || unwrap<CodeReviewResult>(res.data, 'data') || null;
    } catch {
      return null;
    }
  }

  async getComparison(id: string): Promise<CodeComparisonData | null> {
    try {
      const res = await apiClient.get(`/api/v1/ai/code-review/${id}/comparison`);
      return unwrap<CodeComparisonData>(res.data, 'comparison') || null;
    } catch {
      return null;
    }
  }

  async getStudentReviews(limit = 20, offset = 0): Promise<ReviewHistoryItem[]> {
    try {
      const res = await apiClient.get('/api/v1/student/code-reviews', {
        params: { limit, offset },
      });
      return unwrap<ReviewHistoryItem[]>(res.data, 'reviews') || [];
    } catch {
      return [];
    }
  }
}

export const codeReviewService = new CodeReviewService();
export default codeReviewService;

