import { apiClient } from './api';
import {
  LearningMemorySummary,
  ConceptMemory,
  ConceptDetailResponse,
  ConceptStrength,
  ConceptWeakness,
  LearningEventItem,
  TimelinePoint,
  StudentFacingNarrative,
  LearningSnapshot,
} from '../types/learning-memory.types';

export const learningMemoryService = {
  getSummary: async (): Promise<LearningMemorySummary> => {
    const res = await apiClient.get<LearningMemorySummary>('/api/v1/student/learning-memory');
    return res.data;
  },

  getNarrative: async (): Promise<StudentFacingNarrative> => {
    const res = await apiClient.get<StudentFacingNarrative>('/api/v1/student/learning-memory/summary');
    return res.data;
  },

  listConcepts: async (params?: {
    topic_id?: string;
    learning_state?: string;
    search?: string;
    limit?: number;
    offset?: number;
  }): Promise<{ concepts: ConceptMemory[]; total: number }> => {
    const res = await apiClient.get<{ concepts: ConceptMemory[]; total: number }>(
      '/api/v1/student/learning-memory/concepts',
      { params }
    );
    return res.data;
  },

  getConceptDetail: async (conceptKey: string): Promise<ConceptDetailResponse> => {
    const res = await apiClient.get<ConceptDetailResponse>(
      `/api/v1/student/learning-memory/concepts/${encodeURIComponent(conceptKey)}`
    );
    return res.data;
  },

  getStrengths: async (limit = 5): Promise<ConceptStrength[]> => {
    const res = await apiClient.get<{ strengths: ConceptStrength[] }>(
      '/api/v1/student/learning-memory/strengths',
      { params: { limit } }
    );
    return res.data.strengths || [];
  },

  getWeaknesses: async (limit = 5): Promise<ConceptWeakness[]> => {
    const res = await apiClient.get<{ weaknesses: ConceptWeakness[] }>(
      '/api/v1/student/learning-memory/weaknesses',
      { params: { limit } }
    );
    return res.data.weaknesses || [];
  },

  getHistory: async (params?: {
    event_type?: string;
    concept_key?: string;
    limit?: number;
    offset?: number;
  }): Promise<{ events: LearningEventItem[]; total: number }> => {
    const res = await apiClient.get<{ events: LearningEventItem[]; total: number }>(
      '/api/v1/student/learning-memory/history',
      { params }
    );
    return res.data;
  },

  getTimeline: async (days = 14): Promise<TimelinePoint[]> => {
    const res = await apiClient.get<{ timeline: TimelinePoint[] }>(
      '/api/v1/student/learning-memory/timeline',
      { params: { days } }
    );
    return res.data.timeline || [];
  },

  getSnapshots: async (days = 30): Promise<LearningSnapshot[]> => {
    try {
      const res = await apiClient.get<{ snapshots: LearningSnapshot[] }>(
        '/api/v1/student/learning-memory/snapshots',
        { params: { days } }
      );
      return res.data?.snapshots || [];
    } catch {
      return [];
    }
  },
};

export default learningMemoryService;
