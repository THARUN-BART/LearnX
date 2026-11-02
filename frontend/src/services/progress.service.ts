import apiClient, { unwrap } from './api';
import {
  StudentProblemProgress,
  ProgressSummary,
  TopicProgressItem,
  ProgressStatus,
} from '../types/learning.types';
import { cacheInvalidate } from '../utils/apiCache';

function mapProgress(raw: any): StudentProblemProgress {
  return {
    student_id: raw.student_id,
    problem_id: raw.problem_id,
    status: raw.status,
    attempt_count: raw.total_attempts ?? raw.attempt_count ?? 0,
    successful_attempts: raw.successful_attempts ?? 0,
    last_attempted_at: raw.last_attempt_at ?? raw.last_attempted_at,
    solved_at: raw.solved_at,
  };
}

class ProgressService {
  async getAllProgress(): Promise<StudentProblemProgress[]> {
    try {
      const response = await apiClient.get('/api/v1/student/progress');
      const list = unwrap<any[]>(response.data, 'progress');
      return (Array.isArray(list) ? list : []).map(mapProgress);
    } catch {
      return [];
    }
  }

  async getProblemProgress(problemId: string): Promise<StudentProblemProgress | null> {
    try {
      const response = await apiClient.get(`/api/v1/student/progress/problems/${problemId}`);
      const raw = unwrap<any>(response.data, 'progress');
      return raw ? mapProgress(raw) : null;
    } catch {
      return null;
    }
  }

  async recordAttempt(problemId: string): Promise<void> {
    try {
      await apiClient.post(`/api/v1/student/progress/problems/${problemId}/attempt`);
      cacheInvalidate('problems');
      cacheInvalidate('roadmap');
    } catch {
      // non-blocking
    }
  }

  async updateProblemStatus(
    problemId: string,
    status: ProgressStatus,
    reset: boolean = false
  ): Promise<StudentProblemProgress> {
    const payload: { status: ProgressStatus; reset?: boolean } = { status };
    if (reset) payload.reset = true;
    const response = await apiClient.patch(
      `/api/v1/student/progress/problems/${problemId}`,
      payload
    );
    cacheInvalidate('problems');
    cacheInvalidate('problem|');
    cacheInvalidate('roadmap');
    cacheInvalidate('analytics');
    const raw = unwrap<any>(response.data, 'progress');
    return mapProgress(raw || { problem_id: problemId, status });
  }

  async getProgressSummary(): Promise<ProgressSummary> {
    try {
      const response = await apiClient.get('/api/v1/student/progress/summary');
      const summary = unwrap<any>(response.data, 'summary');
      const overview = summary?.overview ?? summary;
      return {
        roadmap_id: overview.roadmap_id ?? '',
        roadmap_name: overview.roadmap_name ?? "Striver's A2Z DSA Sheet",
        total_problems: overview.total_problems ?? 0,
        solved_problems: overview.solved_problems ?? 0,
        attempted_problems: overview.attempted_problems ?? 0,
        pending_problems: overview.pending_problems ?? 0,
        needs_revision: overview.needs_revision ?? 0,
        completion_percentage: overview.completion_percentage ?? 0,
      };
    } catch {
      return {
        roadmap_id: '',
        roadmap_name: "Striver's A2Z DSA Sheet",
        total_problems: 402,
        solved_problems: 0,
        attempted_problems: 0,
        pending_problems: 402,
        needs_revision: 0,
        completion_percentage: 0,
      };
    }
  }

  async getTopicProgress(): Promise<TopicProgressItem[]> {
    try {
      const response = await apiClient.get('/api/v1/student/progress/topics');
      const list = unwrap<any[]>(response.data, 'topics');
      return Array.isArray(list) ? list : [];
    } catch {
      return [];
    }
  }
}

export const progressService = new ProgressService();
export default progressService;
