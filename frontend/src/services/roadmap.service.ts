import apiClient from './api';
import { Topic } from '../types/learning.types';
import { DSA_TOPICS } from '../utils/constants';

class RoadmapService {
  async getRoadmap(): Promise<Topic[]> {
    try {
      const response = await apiClient.get('/api/v1/roadmap');
      return response.data.data || response.data;
    } catch {
      return DSA_TOPICS;
    }
  }

  async getTopicById(topicId: string): Promise<Topic | null> {
    try {
      const response = await apiClient.get(`/api/v1/roadmap/${topicId}`);
      return response.data.data || response.data;
    } catch {
      const topic = DSA_TOPICS.find((t) => t.id === topicId || t.slug === topicId);
      return topic || null;
    }
  }
}

export const roadmapService = new RoadmapService();
export default roadmapService;

