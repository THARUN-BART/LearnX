import apiClient, { unwrap } from './api';
import { Subtopic, Topic } from '../types/learning.types';
import { DSA_TOPICS } from '../utils/constants';
import { cacheKey, cached, TTL } from '../utils/apiCache';

function mapTopic(t: any, activeRoadmapId: string, progressMap: Map<string, { solved: number; total: number }>): Topic {
  const prog = progressMap.get(t.id);
  const meta = DSA_TOPICS.find(
    (m) =>
      m.slug === t.slug ||
      m.topic_order === t.topic_order ||
      m.title.toLowerCase() === (t.name ?? '').toLowerCase()
  );

  const totalProbs = t.problem_count ?? prog?.total ?? meta?.totalProblems ?? 0;
  const solvedProbs = prog?.solved ?? 0;
  const isCompleted = totalProbs > 0 && solvedProbs >= totalProbs;

  return {
    id: t.id,
    roadmap_id: t.roadmap_id ?? activeRoadmapId,
    title: t.name ?? meta?.title ?? 'DSA Topic',
    slug: t.slug ?? meta?.slug ?? '',
    description: t.description ?? meta?.description ?? '',
    topic_order: t.topic_order ?? meta?.topic_order ?? 1,
    subtopic_count: t.subtopic_count ?? meta?.subtopics?.length ?? 1,
    problem_count: totalProbs,
    totalProblems: totalProbs,
    completedProblems: solvedProbs,
    solved_count: solvedProbs,
    completion_percentage: totalProbs > 0 ? Math.round((solvedProbs / totalProbs) * 100) : 0,
    status: isCompleted ? 'completed' : 'in_progress',
    difficulty:
      meta?.difficulty ??
      (t.topic_order <= 3 ? 'Easy' : t.topic_order <= 14 ? 'Medium' : 'Hard'),
    estimatedHours: meta?.estimatedHours ?? Math.max(6, Math.round(totalProbs * 0.4)),
    prerequisites: meta?.prerequisites ?? [],
    subtopics: meta?.subtopics ?? [],
    explanation: meta?.explanation ?? t.description ?? '',
  };
}

class RoadmapService {
  async listRoadmaps(): Promise<any[]> {
    return cached(cacheKey(['roadmaps']), TTL.long, async () => {
      const r = await apiClient.get('/api/v1/roadmap');
      const list = unwrap<any[]>(r.data, 'roadmaps');
      return Array.isArray(list) ? list : [];
    });
  }

  async getRoadmap(): Promise<Topic[]> {
    return cached(cacheKey(['roadmap', 'topics']), TTL.medium, async () => {
      try {
        const roadmaps = await this.listRoadmaps();
        if (!roadmaps.length) return DSA_TOPICS;

        const activeRoadmap =
          roadmaps.find((r: any) => r.slug === 'strivers-a2z-dsa-sheet') ??
          roadmaps.find((r: any) => r.is_active) ??
          roadmaps[0];

        const [topicsRes, progressRes] = await Promise.allSettled([
          apiClient.get(`/api/v1/roadmap/${activeRoadmap.id}/topics`),
          apiClient.get('/api/v1/student/progress/topics'),
        ]);

        const backendTopics =
          topicsRes.status === 'fulfilled'
            ? unwrap<any[]>(topicsRes.value.data, 'topics')
            : [];

        if (!Array.isArray(backendTopics) || backendTopics.length === 0) {
          return DSA_TOPICS;
        }

        const progressMap = new Map<string, { solved: number; total: number }>();
        if (progressRes.status === 'fulfilled') {
          const list = unwrap<any[]>(progressRes.value.data, 'topics');
          (Array.isArray(list) ? list : []).forEach((p: any) => {
            progressMap.set(p.topic_id, {
              solved: p.solved_problems ?? 0,
              total: p.total_problems ?? 0,
            });
          });
        }

        return backendTopics.map((t: any) => mapTopic(t, activeRoadmap.id, progressMap));
      } catch {
        return DSA_TOPICS;
      }
    });
  }

  async getTopicById(topicIdOrSlug: string): Promise<Topic | null> {
    return cached(cacheKey(['topic', topicIdOrSlug]), TTL.medium, async () => {
      try {
        const r = await apiClient.get(`/api/v1/topics/${topicIdOrSlug}`);
        const rawTopic = unwrap<any>(r.data, 'topic');
        if (!rawTopic) return this._fallbackTopic(topicIdOrSlug);

        const [subs, progress] = await Promise.allSettled([
          this.getSubtopics(rawTopic.id),
          apiClient.get('/api/v1/student/progress/topics'),
        ]);

        const subtopics = subs.status === 'fulfilled' ? subs.value : [];
        let solvedCount = 0;
        if (progress.status === 'fulfilled') {
          const list = unwrap<any[]>(progress.value.data, 'topics');
          const item = (Array.isArray(list) ? list : []).find((p: any) => p.topic_id === rawTopic.id);
          if (item) solvedCount = item.solved_problems ?? 0;
        }

        const meta = DSA_TOPICS.find(
          (m) =>
            m.id === rawTopic.id ||
            m.slug === rawTopic.slug ||
            m.topic_order === rawTopic.topic_order
        );

        const total = rawTopic.problem_count ?? meta?.totalProblems ?? 0;
        return {
          id: rawTopic.id,
          roadmap_id: rawTopic.roadmap_id ?? '',
          title: rawTopic.name ?? rawTopic.title ?? meta?.title ?? 'Topic',
          slug: rawTopic.slug ?? meta?.slug ?? '',
          description: rawTopic.description ?? meta?.description ?? '',
          topic_order: rawTopic.topic_order ?? meta?.topic_order ?? 1,
          subtopic_count: subtopics.length || rawTopic.subtopic_count || meta?.subtopics?.length || 1,
          problem_count: total,
          totalProblems: total,
          completedProblems: solvedCount,
          solved_count: solvedCount,
          completion_percentage: total > 0 ? Math.round((solvedCount / total) * 100) : 0,
          status: total > 0 && solvedCount >= total ? 'completed' : 'in_progress',
          difficulty:
            meta?.difficulty ??
            (rawTopic.topic_order <= 3 ? 'Easy' : rawTopic.topic_order <= 14 ? 'Medium' : 'Hard'),
          estimatedHours: meta?.estimatedHours ?? Math.max(6, Math.round(total * 0.4)),
          prerequisites: meta?.prerequisites ?? [],
          subtopics: subtopics.length ? subtopics.map((s) => s.title) : meta?.subtopics ?? [],
          explanation: meta?.explanation ?? rawTopic.description ?? 'Master core patterns and problem-solving techniques.',
        };
      } catch {
        return this._fallbackTopic(topicIdOrSlug);
      }
    });
  }

  async getSubtopics(topicId: string): Promise<Subtopic[]> {
    return cached(cacheKey(['subtopics', topicId]), TTL.long, async () => {
      try {
        const sr = await apiClient.get(`/api/v1/topics/${topicId}/subtopics`);
        const subs = unwrap<any[]>(sr.data, 'subtopics');
        return (Array.isArray(subs) ? subs : []).map((s: any) => ({
          id: s.id,
          topic_id: s.topic_id ?? topicId,
          title: s.name ?? s.title ?? s.slug,
          slug: s.slug ?? '',
          description: s.description ?? '',
          subtopic_order: s.subtopic_order ?? 0,
          problem_count: s.problem_count ?? 0,
        }));
      } catch {
        return [];
      }
    });
  }

  private _fallbackTopic(topicIdOrSlug: string): Topic | null {
    return (
      DSA_TOPICS.find((t) => t.id === topicIdOrSlug || t.slug === topicIdOrSlug) ??
      DSA_TOPICS[0] ??
      null
    );
  }
}

export const roadmapService = new RoadmapService();
export default roadmapService;
