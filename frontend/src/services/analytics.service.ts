import apiClient, { unwrap } from './api';
import {
  AnalyticsSummary,
  MistakeRecord,
  MistakeCategory,
  MistakeSummary,
  StudentMistakePattern,
} from '../types/learning.types';
import { cacheKey, cached, TTL } from '../utils/apiCache';

const DAYS_OF_WEEK = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];

function mapCategory(raw: string | undefined): MistakeCategory {
  const catLower = (raw ?? '').toLowerCase();
  if (catLower === 'compilation' || catLower === 'syntax') return 'syntax';
  if (catLower === 'complexity' || catLower.includes('time')) return 'time_complexity';
  if (catLower === 'edge_case' || catLower === 'edge_cases') return 'edge_cases';
  if (catLower === 'algorithmic' || catLower === 'algorithm_selection') return 'algorithm_selection';
  if (catLower.includes('space') || catLower === 'data_structure') return 'space_complexity';
  return 'logical';
}

function mapDiagnosis(d: any, idx: number): MistakeRecord {
  const category = mapCategory(d.category ?? d.primary_category);
  return {
    id: d.id ?? `mistake-${idx}`,
    problemId: d.problem_id ?? '',
    problemTitle: d.evidence?.problem_title ?? d.problem_title ?? 'DSA Problem',
    topic: d.related_concept ?? d.dsa_concept ?? 'DSA Concept',
    category,
    categoryLabel: `${(d.category ?? d.primary_category ?? 'logical').replace(/_/g, ' ')}: ${(d.subtype ?? d.subcategory ?? '').replace(/_/g, ' ')}`.trim(),
    frequency: d.historical_occurrence_count ?? (d.evidence?.previous_mistake_count ? d.evidence.previous_mistake_count + 1 : 1),
    date: d.created_at ? new Date(d.created_at).toISOString().split('T')[0] : new Date().toISOString().split('T')[0],
    failedCodeSnippet: d.evidence?.source_code ?? '',
    explanation: d.description ?? d.what_went_wrong ?? 'Intelligent diagnosis generated.',
    whyFailed: d.probable_cause ?? `Detected via ${(d.detection_method ?? '').replace(/_/g, ' ')} (Confidence: ${d.confidence_level ?? d.confidence_score ?? 'Unknown'})`,
    suggestedCorrection: d.suggested_correction ?? d.actionable_fix ?? 'Review logic and invariants.',
    recommendedTopicId: '',
    recommendedTopicTitle: d.related_concept ?? d.dsa_concept ?? 'DSA Curriculum',
  };
}

class AnalyticsService {
  async getAnalytics(): Promise<AnalyticsSummary> {
    return cached(cacheKey(['analytics']), TTL.short, async () => {
      try {
        const [summaryRes, activityRes] = await Promise.allSettled([
          apiClient.get('/api/v1/student/progress/summary'),
          apiClient.get('/api/v1/student/analytics/activity?days=30'),
        ]);

        const summaryPayload =
          summaryRes.status === 'fulfilled' ? unwrap<any>(summaryRes.value.data, 'summary') : null;
        const overview: Record<string, any> = summaryPayload?.overview ?? {};
        const diffStats: any[] = summaryPayload?.difficulty ?? [];
        const topicsPayload: any[] = summaryPayload?.topics ?? [];

        const totalSolved: number = overview.solved_problems ?? 0;
        const attempted: number = overview.attempted_problems ?? 0;
        const totalProblems: number = overview.total_problems ?? 402;

        const easyStats = diffStats.find((d: any) => d.difficulty === 'Easy') ?? {};
        const medStats = diffStats.find((d: any) => d.difficulty === 'Medium') ?? {};
        const hardStats = diffStats.find((d: any) => d.difficulty === 'Hard') ?? {};

        const accuracyRate =
          (easyStats.accuracy ?? medStats.accuracy ?? hardStats.accuracy) != null
            ? Math.round(
                ((easyStats.accuracy || 0) + (medStats.accuracy || 0) + (hardStats.accuracy || 0)) /
                  Math.max(1, diffStats.length) *
                  1000
              ) / 10
            : attempted > 0
            ? Math.min(100, Math.round((totalSolved / attempted) * 1000) / 10)
            : totalSolved > 0
            ? 100
            : 0;

        const activityPayload =
          activityRes.status === 'fulfilled' ? unwrap<any>(activityRes.value.data, 'activity') : null;
        const rawPoints: any[] = Array.isArray(activityPayload?.activity)
          ? activityPayload.activity
          : Array.isArray(activityPayload)
          ? activityPayload
          : [];

        const weeklyActivity = DAYS_OF_WEEK.map((dayName, index) => {
          const matching = rawPoints.filter((p: any) => p.date && new Date(p.date).getDay() === index);
          const solved = matching.reduce((s, p) => s + (p.solved ?? 0), 0);
          const attempts = matching.reduce((s, p) => s + (p.attempts ?? 0), 0);
          return { day: dayName, solved, hours: Math.round(attempts * 0.3 * 10) / 10 };
        });

        let streakDays = 0;
        const sortedPts = [...rawPoints].sort((a: any, b: any) => (b.date ?? '').localeCompare(a.date ?? ''));
        for (const pt of sortedPts) {
          if ((pt.solved ?? 0) > 0 || (pt.attempts ?? 0) > 0) streakDays++;
          else break;
        }

        const topicProficiency = (Array.isArray(topicsPayload) ? topicsPayload : [])
          .map((t: any) => ({
            topic: t.topic_name ?? t.name ?? 'Topic',
            topicId: t.topic_id,
            score: Math.min(
              100,
              Math.round(t.completion_percentage ?? ((t.solved_problems ?? 0) / Math.max(1, t.total_problems ?? 1)) * 100)
            ),
            total: t.total_problems ?? 0,
            solved: t.solved_problems ?? 0,
          }))
          .sort((a, b) => b.score - a.score);

        const ranked = topicProficiency.filter((t) => t.total > 0);
        const strongestTopic = ranked[0];
        const weakestTopic = ranked.length ? ranked[ranked.length - 1] : undefined;

        const hoursSpent = Math.max(0, Math.round(attempted * 0.35 * 10) / 10);

        return {
          totalSolved,
          totalProblems,
          attemptedProblems: attempted,
          pendingProblems: overview.pending_problems ?? Math.max(0, totalProblems - totalSolved),
          needsRevision: overview.needs_revision ?? 0,
          completionPercentage: Math.round(overview.completion_percentage ?? 0),
          easySolved: easyStats.solved ?? 0,
          mediumSolved: medStats.solved ?? 0,
          hardSolved: hardStats.solved ?? 0,
          streakDays,
          hoursSpent,
          accuracyRate,
          weeklyActivity,
          dailyActivity: rawPoints.map((p: any) => ({
            date: p.date,
            solved: p.solved ?? 0,
            attempts: p.attempts ?? 0,
          })),
          topicProficiency: topicProficiency.slice(0, 12),
          difficultyDistribution: [
            { name: 'Easy', count: easyStats.total ?? easyStats.solved ?? 0, color: '#22C55E' },
            { name: 'Medium', count: medStats.total ?? medStats.solved ?? 0, color: '#F59E0B' },
            { name: 'Hard', count: hardStats.total ?? hardStats.solved ?? 0, color: '#EF4444' },
          ],
          monthlyProgress: [{ month: 'Current', solved: totalSolved }],
          strongestTopic,
          weakestTopic,
        };
      } catch {
        return this._emptyAnalytics();
      }
    });
  }

  async getMistakes(): Promise<MistakeRecord[]> {
    return cached(cacheKey(['mistakes']), TTL.short, async () => {
      try {
        const diagRes = await apiClient.get('/api/v1/student/mistakes', { params: { limit: 50 } });
        const diagnoses = unwrap<any[]>(diagRes.data, 'data');
        if (Array.isArray(diagnoses) && diagnoses.length > 0) {
          return diagnoses.map(mapDiagnosis);
        }

        const res = await apiClient.get('/api/v1/student/submissions', { params: { limit: 50 } });
        const payload = unwrap<any>(res.data, 'data', 'submissions');
        const submissions = payload?.submissions ?? payload;
        if (!Array.isArray(submissions) || submissions.length === 0) return [];

        return submissions
          .filter(
            (s: any) =>
              s.status && s.status !== 'Accepted' && s.status !== 'Queued' && s.status !== 'Processing'
          )
          .map((s: any, idx: number) => {
            const category = mapCategory(s.status);
            return {
              id: s.id ?? s.submission_id ?? `mistake-${idx}`,
              problemId: s.problem_id ?? '',
              problemTitle: s.problem_title ?? 'Problem Solution',
              topic: s.topic_name ?? 'DSA Curriculum',
              category,
              categoryLabel: s.status,
              frequency: 1,
              date: s.submitted_at ? new Date(s.submitted_at).toISOString().split('T')[0] : new Date().toISOString().split('T')[0],
              failedCodeSnippet: s.source_code ?? '',
              explanation:
                s.error_message ??
                `Execution terminated with ${s.status}. Passed ${s.passed_test_cases ?? 0} of ${s.total_test_cases ?? 0} test cases.`,
              whyFailed: s.error_message ?? `Test evaluation failed: status ${s.status}`,
              suggestedCorrection:
                'Review boundary conditions and verify the solution against sample and edge test cases.',
              recommendedTopicId: s.topic_id ?? '',
              recommendedTopicTitle: s.topic_name ?? 'Data Structures & Algorithms',
            };
          });
      } catch {
        return [];
      }
    });
  }

  async getMistakeSummary(): Promise<MistakeSummary | null> {
    try {
      const res = await apiClient.get('/api/v1/student/mistakes/summary');
      return unwrap<MistakeSummary>(res.data, 'data');
    } catch {
      return null;
    }
  }

  async getMistakePatterns(): Promise<StudentMistakePattern[]> {
    try {
      const res = await apiClient.get('/api/v1/student/mistakes/patterns');
      const list = unwrap<any[]>(res.data, 'data');
      return Array.isArray(list) ? list : [];
    } catch {
      return [];
    }
  }

  private _emptyAnalytics(): AnalyticsSummary {
    return {
      totalSolved: 0,
      totalProblems: 402,
      attemptedProblems: 0,
      pendingProblems: 402,
      needsRevision: 0,
      completionPercentage: 0,
      easySolved: 0,
      mediumSolved: 0,
      hardSolved: 0,
      streakDays: 0,
      hoursSpent: 0,
      accuracyRate: 0,
      weeklyActivity: DAYS_OF_WEEK.map((day) => ({ day, solved: 0, hours: 0 })),
      dailyActivity: [],
      topicProficiency: [],
      difficultyDistribution: [
        { name: 'Easy', count: 0, color: '#22C55E' },
        { name: 'Medium', count: 0, color: '#F59E0B' },
        { name: 'Hard', count: 0, color: '#EF4444' },
      ],
      monthlyProgress: [{ month: 'Current', solved: 0 }],
    };
  }
}

export const analyticsService = new AnalyticsService();
export default analyticsService;
