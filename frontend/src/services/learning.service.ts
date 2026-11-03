import apiClient, { unwrap } from './api';
import { Problem } from '../types/learning.types';
import { cacheKey, cached, TTL } from '../utils/apiCache';
import { enrichProblemData } from '../utils/problemCatalog';

export function mapRawProblemToProblem(raw: any, sampleCases?: any[]): Problem {
  return enrichProblemData(raw, sampleCases);
}

export interface ProblemListParams {
  topicId?: string;
  subtopicId?: string;
  difficulty?: string;
  status?: string;
  q?: string;
  page?: number;
  limit?: number;
}

class LearningService {
  async getProblems(topicIdOrParams?: string | ProblemListParams): Promise<Problem[]> {
    const params: ProblemListParams =
      typeof topicIdOrParams === 'string' ? { topicId: topicIdOrParams } : topicIdOrParams || {};
    const limit = params.limit ?? 40;
    const page = params.page ?? 1;
    const key = cacheKey([
      'problems',
      params.topicId,
      params.subtopicId,
      params.difficulty,
      params.status,
      params.q,
      page,
      limit,
    ]);

    return cached(key, TTL.short, async () => {
      try {
        const response = await apiClient.get('/api/v1/problems', {
          params: {
            topic_id: params.topicId,
            subtopic_id: params.subtopicId,
            difficulty: params.difficulty,
            status: params.status,
            q: params.q,
            page,
            limit,
            sort: 'problem_order',
          },
        });
        const rawProblems = unwrap<any[]>(response.data, 'problems');
        if (!Array.isArray(rawProblems)) return [];
        return rawProblems.map((p) => mapRawProblemToProblem(p));
      } catch {
        return [];
      }
    });
  }

  async searchProblems(q: string, limit = 20): Promise<Problem[]> {
    if (!q.trim()) return [];
    try {
      const response = await apiClient.get('/api/v1/problems/search', {
        params: { q: q.trim(), limit },
      });
      const rawProblems = unwrap<any[]>(response.data, 'problems');
      if (!Array.isArray(rawProblems)) return [];
      return rawProblems.map((p) => mapRawProblemToProblem(p));
    } catch {
      return [];
    }
  }

  async getProblemsBySubtopic(subtopicId: string): Promise<Problem[]> {
    return cached(cacheKey(['subtopic-problems', subtopicId]), TTL.medium, async () => {
      try {
        const response = await apiClient.get(`/api/v1/subtopics/${subtopicId}/problems`, {
          params: { limit: 80 },
        });
        const rawProblems = unwrap<any[]>(response.data, 'problems');
        if (!Array.isArray(rawProblems)) return [];
        return rawProblems.map((p) => mapRawProblemToProblem(p));
      } catch {
        return this.getProblems({ subtopicId, limit: 80 });
      }
    });
  }

  async getProblemById(problemId: string): Promise<Problem | null> {
    return cached(cacheKey(['problem', problemId]), TTL.short, async () => {
      try {
        const [response, tcRes] = await Promise.allSettled([
          apiClient.get(`/api/v1/problems/${problemId}`),
          apiClient.get(`/api/v1/problems/${problemId}/testcases/sample`),
        ]);

        if (response.status !== 'fulfilled') return null;
        const raw = unwrap<any>(response.value.data, 'problem');
        if (!raw || !raw.id) return null;

        const related = response.value.data?.related;
        const samples =
          tcRes.status === 'fulfilled' ? unwrap<any[]>(tcRes.value.data, 'data', 'test_cases') : [];

        const problem = mapRawProblemToProblem(
          { ...raw, related },
          Array.isArray(samples) ? samples : []
        );
        return problem;
      } catch {
        return null;
      }
    });
  }
}

export const learningService = new LearningService();
export default learningService;
