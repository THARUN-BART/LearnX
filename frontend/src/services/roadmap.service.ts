import apiClient from './api';
import { Topic } from '../types/learning.types';
import { DSA_TOPICS, TOKEN_KEY } from '../utils/constants';

class RoadmapService {
  /**
   * Fetches the complete list of pillars/topics for Striver's A2Z DSA sheet from the live database.
   * Merges real-time student progress when authenticated.
   */
  async getRoadmap(): Promise<Topic[]> {
    try {
      // 1. Fetch available roadmaps (use authenticated route if logged in, otherwise public alias)
      let roadmaps: any[] = [];
      const hasToken = typeof window !== 'undefined' && !!localStorage.getItem(TOKEN_KEY);

      if (hasToken) {
        try {
          const roadmapsRes = await apiClient.get('/api/v1/roadmaps');
          roadmaps = roadmapsRes.data.roadmaps || roadmapsRes.data.data || [];
        } catch {
          // Fallback to public route if token expired
        }
      }

      if (!Array.isArray(roadmaps) || roadmaps.length === 0) {
        const fallbackRes = await apiClient.get('/api/v1/roadmap');
        roadmaps = fallbackRes.data.roadmaps || fallbackRes.data.data || [];
      }

      if (!Array.isArray(roadmaps) || roadmaps.length === 0) {
        return DSA_TOPICS;
      }

      // Active roadmap (Striver's A2Z or first roadmap)
      const activeRoadmap = roadmaps.find((r: any) => r.slug === 'strivers-a2z') || roadmaps[0];

      // 2. Fetch topics under this roadmap
      let backendTopics: any[] = [];
      if (hasToken) {
        try {
          const topicsRes = await apiClient.get(`/api/v1/roadmaps/${activeRoadmap.id}/topics`);
          backendTopics = topicsRes.data.topics || topicsRes.data.data || [];
        } catch {
          // Fallback to public route
        }
      }

      if (!Array.isArray(backendTopics) || backendTopics.length === 0) {
        const fallbackRes = await apiClient.get(`/api/v1/roadmap/${activeRoadmap.id}/topics`);
        backendTopics = fallbackRes.data.topics || fallbackRes.data.data || [];
      }

      if (!Array.isArray(backendTopics) || backendTopics.length === 0) {
        return DSA_TOPICS;
      }

      // 3. Try to fetch student's live topic progress from the database
      const progressMap = new Map<string, { solved: number; total: number }>();
      try {
        const progRes = await apiClient.get('/api/v1/student/progress/topics');
        const progList = progRes.data.topics || progRes.data.data || [];
        if (Array.isArray(progList)) {
          progList.forEach((p: any) => {
            progressMap.set(p.topic_id, {
              solved: p.solved_problems || 0,
              total: p.total_problems || 0,
            });
          });
        }
      } catch {
        // Unauthenticated - continue with zero progress
      }

      // 4. Map real backend topics to UI Topic model
      return backendTopics.map((t: any) => {
        const prog = progressMap.get(t.id);
        const meta = DSA_TOPICS.find(
          (m) => m.slug === t.slug || m.topic_order === t.topic_order || m.title.toLowerCase() === t.name?.toLowerCase()
        );

        const totalProbs = t.problem_count || meta?.totalProblems || 0;
        const solvedProbs = prog?.solved || 0;
        const isCompleted = totalProbs > 0 && solvedProbs >= totalProbs;

        return {
          id: t.id,
          roadmap_id: t.roadmap_id || activeRoadmap.id,
          title: t.name || t.title || meta?.title || 'DSA Topic',
          slug: t.slug || meta?.slug || '',
          description: t.description || meta?.description || '',
          topic_order: t.topic_order || meta?.topic_order || 1,
          subtopic_count: t.subtopic_count || meta?.subtopics?.length || 1,
          problem_count: totalProbs,
          totalProblems: totalProbs,
          completedProblems: solvedProbs,
          status: isCompleted ? 'completed' : 'in_progress',
          difficulty: meta?.difficulty || (t.topic_order <= 3 ? 'Easy' : t.topic_order <= 14 ? 'Medium' : 'Hard'),
          estimatedHours: meta?.estimatedHours || Math.max(6, Math.round(totalProbs * 0.4)),
          prerequisites: meta?.prerequisites || [],
          subtopics: meta?.subtopics || [],
          explanation: meta?.explanation || t.description || ''
        };
      });
    } catch {
      return DSA_TOPICS;
    }
  }

  /**
   * Fetches real topic details, subtopics, and current completion metrics from the backend.
   */
  async getTopicById(topicIdOrSlug: string): Promise<Topic | null> {
    try {
      // 1. Fetch real topic details from database
      const response = await apiClient.get(`/api/v1/topics/${topicIdOrSlug}`);
      const rawTopic = response.data.topic || response.data.data;
      if (!rawTopic) {
        return this.fallbackTopic(topicIdOrSlug);
      }

      // 2. Fetch real subtopics for this topic
      let subtopicNames: string[] = [];
      try {
        const subRes = await apiClient.get(`/api/v1/topics/${rawTopic.id}/subtopics`);
        const subs = subRes.data.subtopics || subRes.data.data || [];
        if (Array.isArray(subs)) {
          subtopicNames = subs.map((s: any) => s.name || s.title || s.slug);
        }
      } catch {
        // Fallback subtopics from constants
      }

      // 3. Match with metadata
      const meta = DSA_TOPICS.find(
        (m) => m.id === rawTopic.id || m.slug === rawTopic.slug || m.topic_order === rawTopic.topic_order
      );

      // 4. Fetch real student progress
      let solvedCount = 0;
      try {
        const progRes = await apiClient.get('/api/v1/student/progress/topics');
        const progList = progRes.data.topics || progRes.data.data || [];
        const item = progList.find((p: any) => p.topic_id === rawTopic.id);
        if (item) {
          solvedCount = item.solved_problems || 0;
        }
      } catch {
        // Default 0
      }

      const total = rawTopic.problem_count || meta?.totalProblems || 0;
      const isCompleted = total > 0 && solvedCount >= total;

      return {
        id: rawTopic.id,
        roadmap_id: rawTopic.roadmap_id || '',
        title: rawTopic.name || rawTopic.title || meta?.title || 'Topic',
        slug: rawTopic.slug || meta?.slug || '',
        description: rawTopic.description || meta?.description || '',
        topic_order: rawTopic.topic_order || meta?.topic_order || 1,
        subtopic_count: subtopicNames.length || rawTopic.subtopic_count || meta?.subtopics?.length || 1,
        problem_count: total,
        totalProblems: total,
        completedProblems: solvedCount,
        status: isCompleted ? 'completed' : 'in_progress',
        difficulty: meta?.difficulty || (rawTopic.topic_order <= 3 ? 'Easy' : rawTopic.topic_order <= 14 ? 'Medium' : 'Hard'),
        estimatedHours: meta?.estimatedHours || Math.max(6, Math.round(total * 0.4)),
        prerequisites: meta?.prerequisites || [],
        subtopics: subtopicNames.length > 0 ? subtopicNames : (meta?.subtopics || []),
        explanation: meta?.explanation || rawTopic.description || 'Master core patterns and problem-solving techniques.'
      };
    } catch {
      return this.fallbackTopic(topicIdOrSlug);
    }
  }

  private fallbackTopic(topicIdOrSlug: string): Topic | null {
    const topic = DSA_TOPICS.find((t) => t.id === topicIdOrSlug || t.slug === topicIdOrSlug);
    return topic || DSA_TOPICS[0] || null;
  }
}

export const roadmapService = new RoadmapService();
export default roadmapService;
