import apiClient from './api';
import {
  StudentProblemProgress,
  ProgressSummary,
  TopicProgressItem,
  ProgressStatus
} from '../types/learning.types';

class ProgressService {
  // GET /api/v1/student/progress
  async getAllProgress(): Promise<StudentProblemProgress[]> {
    try {
      const response = await apiClient.get('/api/v1/student/progress');
      return response.data.progress || response.data.data || [];
    } catch {
      return [];
    }
  }

  // GET /api/v1/student/progress/problems/{problem_id}
  async getProblemProgress(problemId: string): Promise<StudentProblemProgress | null> {
    try {
      const response = await apiClient.get(`/api/v1/student/progress/problems/${problemId}`);
      return response.data.progress || response.data.data || null;
    } catch {
      return null;
    }
  }

  // POST /api/v1/student/progress/problems/{problem_id}/attempt
  async recordAttempt(problemId: string): Promise<void> {
    try {
      await apiClient.post(`/api/v1/student/progress/problems/${problemId}/attempt`);
    } catch (e) {
      console.warn('[ProgressService] Error recording attempt:', e);
    }
  }

  // PATCH /api/v1/student/progress/problems/{problem_id}
  async updateProblemStatus(
    problemId: string,
    status: ProgressStatus,
    reset: boolean = false
  ): Promise<StudentProblemProgress> {
    try {
      const payload: { status: ProgressStatus; reset?: boolean } = { status };
      if (reset) payload.reset = true;
      const response = await apiClient.patch(
        `/api/v1/student/progress/problems/${problemId}`,
        payload
      );
      return response.data.progress || response.data;
    } catch (err: any) {
      // Local fallback for offline mode
      const mock: StudentProblemProgress = {
        student_id: 'usr_local',
        problem_id: problemId,
        status,
        attempt_count: 1,
        successful_attempts: status === 'Solved' ? 1 : 0,
        solved_at: status === 'Solved' ? new Date().toISOString() : undefined
      };
      return mock;
    }
  }

  // GET /api/v1/student/progress/summary
  async getProgressSummary(): Promise<ProgressSummary> {
    try {
      const response = await apiClient.get('/api/v1/student/progress/summary');
      return response.data.summary || response.data.data || response.data;
    } catch {
      return {
        roadmap_id: '',
        roadmap_name: "Striver's A2Z DSA Sheet",
        total_problems: 402,
        solved_problems: 0,
        attempted_problems: 0,
        pending_problems: 402,
        needs_revision: 0,
        completion_percentage: 0
      };
    }
  }

  // GET /api/v1/student/progress/topics
  async getTopicProgress(): Promise<TopicProgressItem[]> {
    try {
      const response = await apiClient.get('/api/v1/student/progress/topics');
      return response.data.topics || response.data.data || [];
    } catch {
      return [];
    }
  }
}

export const progressService = new ProgressService();
export default progressService;

