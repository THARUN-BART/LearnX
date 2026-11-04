import { normalizePercentage } from '../../utils/formatters';
import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Flame,
  CheckCircle2,
  Clock,
  TrendingUp,
  Play,
  Sparkles,
  ChevronRight,
  Activity,
  Layers,
  Compass,
} from 'lucide-react';
import {
  ResponsiveContainer,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip
} from 'recharts';
import { useAuth } from '../../hooks/useAuth';
import { useProfile } from '../../hooks/useProfile';
import { analyticsService } from '../../services/analytics.service';
import { learningService } from '../../services/learning.service';
import { roadmapService } from '../../services/roadmap.service';
import { learningPathService } from '../../services/learning-path.service';
import { recommendationService } from '../../services/recommendation.service';
import { AnalyticsSummary, Problem, Topic, MistakeRecord } from '../../types/learning.types';
import { LearningRecommendation } from '../../types/learning-path.types';
import { ProblemRecommendation } from '../../types/recommendations.types';
import { RecommendedActionCard } from '../../components/learning-path/RecommendedActionCard';
import Button from '../../components/ui/Button';
import { PageSkeleton } from '../../components/ui/Skeleton';

export const DashboardPage: React.FC = () => {
  const { user } = useAuth();
  const { profile } = useProfile();
  const navigate = useNavigate();

  const [analytics, setAnalytics] = useState<AnalyticsSummary | null>(null);
  const [recommendedProblems, setRecommendedProblems] = useState<Problem[]>([]);
  const [smartRecommendations, setSmartRecommendations] = useState<ProblemRecommendation[]>([]);
  const [activeTopic, setActiveTopic] = useState<Topic | null>(null);
  const [latestMistake, setLatestMistake] = useState<MistakeRecord | null>(null);
  const [adaptiveNextAction, setAdaptiveNextAction] = useState<LearningRecommendation | null>(null);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    let cancelled = false;

    const loadCore = async () => {
      try {
        const [stats, topics, nextAction] = await Promise.all([
          analyticsService.getAnalytics(),
          roadmapService.getRoadmap(),
          learningPathService.getNextAction().catch(() => null),
        ]);
        if (cancelled) return;
        setAnalytics(stats);
        setAdaptiveNextAction(nextAction);
        if (Array.isArray(topics) && topics.length > 0) {
          const current =
            topics.find((t) => t.status === 'in_progress' && (t.completedProblems || 0) > 0) ||
            topics.find((t) => t.status === 'in_progress') ||
            topics[0];
          setActiveTopic(current);
        }
      } finally {
        if (!cancelled) setLoading(false);
      }
    };

    const loadSecondary = async () => {
      const [problems, mistakes, recRes] = await Promise.all([
        learningService.getProblems({ status: 'Not Started', limit: 8 }),
        analyticsService.getMistakes(),
        recommendationService.getRecommendations().catch(() => null),
      ]);
      if (cancelled) return;

      if (recRes) {
        const top = recRes.top_pick ? [recRes.top_pick] : [];
        const general = recRes.categories?.recommended_for_you || [];
        const combined = [...top, ...general.filter((r) => r.id !== recRes.top_pick?.id)].slice(0, 3);
        setSmartRecommendations(combined);
      }

      let list = Array.isArray(problems) ? problems.slice(0, 3) : [];
      if (!list.length) {
        const fallback = await learningService.getProblems({ limit: 8 });
        if (cancelled) return;
        list = fallback.slice(0, 3);
      }
      setRecommendedProblems(list);
      if (Array.isArray(mistakes) && mistakes.length > 0) setLatestMistake(mistakes[0]);
    };

    loadCore().then(loadSecondary);
    return () => {
      cancelled = true;
    };
  }, []);

  const studentName = user?.fullName || user?.full_name || 'Student';
  const totalProblems = analytics?.totalProblems || 402;
  const solved = analytics?.totalSolved ?? 0;

  const activityCells = (() => {
    const points = analytics?.dailyActivity || [];
    const byDate = new Map(points.map((p) => [p.date.slice(0, 10), p]));
    const cells: { date: string; level: number }[] = [];
    for (let i = 139; i >= 0; i--) {
      const d = new Date();
      d.setDate(d.getDate() - i);
      const key = d.toISOString().slice(0, 10);
      const pt = byDate.get(key);
      const solvedDay = pt?.solved ?? 0;
      const attempts = pt?.attempts ?? 0;
      const level = solvedDay >= 3 ? 3 : solvedDay >= 1 ? 2 : attempts > 0 ? 1 : 0;
      cells.push({ date: key, level });
    }
    const weeks: { date: string; level: number }[][] = [];
    for (let i = 0; i < cells.length; i += 7) weeks.push(cells.slice(i, i + 7));
    return weeks;
  })();

  if (loading && !analytics) {
    return <PageSkeleton rows={3} />;
  }

  return (
    <div className="space-y-8">
      <div className="p-6 sm:p-7 rounded-2xl bg-white border border-[#E2E8F0] relative overflow-hidden">
        <div className="absolute inset-0 bg-tech-grid opacity-30 pointer-events-none" />
        <div className="relative z-10 flex flex-col lg:flex-row lg:items-center justify-between gap-6">
          <div className="space-y-2">
            <div className="flex flex-wrap items-center gap-2 font-mono text-[11px]">
              <span className="px-2 py-0.5 rounded bg-[#EF4444]/15 text-[#DC2626] border border-[#EF4444]/30">
                HOST: LEARNX-ENGINE
              </span>
              <span className="text-[#94A3B8]">//</span>
              <span className="text-emerald-400 flex items-center gap-1">
                <span className="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse" />
                ONLINE
              </span>
              <span className="text-[#94A3B8]">//</span>
              <span className="text-[#64748B]">
                GOAL: {(profile?.learning_goal || 'INTERVIEW PREPARATION').toUpperCase()}
              </span>
            </div>
            <h1 className="text-2xl sm:text-3xl font-bold tracking-tight text-[#0F172A] font-mono">
              Welcome back, {studentName}
            </h1>
            <p className="text-xs text-[#64748B] max-w-2xl leading-relaxed">
              Current focus: <strong className="text-white">{activeTopic?.title || 'Striver A2Z roadmap'}</strong>.
              {analytics?.needsRevision ? (
                <span>
                  {' '}
                  <button
                    onClick={() => navigate('/practice', { state: { filterStatus: 'Needs Revision' } })}
                    className="text-amber-400 underline hover:text-amber-300 font-semibold"
                  >
                    {analytics.needsRevision} problems are marked for revision
                  </button>
                  .
                </span>
              ) : (
                ' Pick up the next unsolved problem from your active pillar.'
              )}
            </p>
          </div>
          <div className="flex items-center gap-3 shrink-0">
            <Button variant="outline" size="sm" onClick={() => navigate('/roadmap')} leftIcon={<Layers className="w-4 h-4" />}>
              Roadmap
            </Button>
            <Button
              variant="primary"
              size="sm"
              onClick={() => {
                if (smartRecommendations[0]) {
                  navigate(`/practice/${smartRecommendations[0].problem_id}`, {
                    state: {
                      fromRecommendation: true,
                      recommendationId: smartRecommendations[0].id,
                      category: smartRecommendations[0].category,
                      primaryReason: smartRecommendations[0].explanation,
                    }
                  });
                } else if (recommendedProblems[0]) {
                  navigate(`/practice/${recommendedProblems[0].id}`);
                } else {
                  navigate('/practice');
                }
              }}
              rightIcon={<Play className="w-4 h-4" />}
            >
              Resume Coding
            </Button>
          </div>
        </div>
      </div>

      {/* Adaptive Learning Path: Recommended Next Action */}
      <div className="space-y-2">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2 font-mono text-xs text-[#64748B] uppercase tracking-wider font-semibold">
            <Compass className="w-4 h-4 text-[#EF4444]" />
            <span>AI Adaptive Path Engine</span>
          </div>
          <button
            onClick={() => navigate('/dashboard/learning-path')}
            className="text-xs text-[#EF4444] hover:text-[#DC2626] flex items-center gap-1 font-mono"
          >
            Explore Full Roadmap <ChevronRight className="w-3.5 h-3.5" />
          </button>
        </div>
        <RecommendedActionCard recommendation={adaptiveNextAction} loading={loading} />
      </div>

      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="p-5 rounded-2xl bg-white border border-[#E2E8F0] space-y-1 font-mono">
          <div className="flex items-center justify-between text-xs text-[#94A3B8]">
            <span>VERIFIED_SOLVED</span>
            <CheckCircle2 className="w-4 h-4 text-emerald-400" />
          </div>
          <div className="text-2xl sm:text-3xl font-bold text-[#0F172A]">{solved}</div>
          <span className="text-[11px] text-emerald-400 block">
            {analytics?.attemptedProblems ?? 0} attempted · {analytics?.pendingProblems ?? totalProblems} pending
          </span>
        </div>
        <div className="p-5 rounded-2xl bg-white border border-[#E2E8F0] space-y-1 font-mono">
          <div className="flex items-center justify-between text-xs text-[#94A3B8]">
            <span>STREAK_STABILITY</span>
            <Flame className="w-4 h-4 text-amber-400" />
          </div>
          <div className="text-2xl sm:text-3xl font-bold text-amber-400">{analytics?.streakDays ?? 0}d</div>
          <span className="text-[11px] text-[#64748B] block">
            {(analytics?.streakDays ?? 0) > 0 ? 'Consecutive active days' : 'Solve today to start a streak'}
          </span>
        </div>
        <div className="p-5 rounded-2xl bg-white border border-[#E2E8F0] space-y-1 font-mono">
          <div className="flex items-center justify-between text-xs text-[#94A3B8]">
            <span>ACTIVE_RUNTIME</span>
            <Clock className="w-4 h-4 text-[#DC2626]" />
          </div>
          <div className="text-2xl sm:text-3xl font-bold text-[#DC2626]">{analytics?.hoursSpent ?? 0}h</div>
          <span className="text-[11px] text-[#64748B] block">Estimated practice time</span>
        </div>
        <div className="p-5 rounded-2xl bg-white border border-[#E2E8F0] space-y-1 font-mono">
          <div className="flex items-center justify-between text-xs text-[#94A3B8]">
            <span>MASTERY_INDEX</span>
            <TrendingUp className="w-4 h-4 text-[#DC2626]" />
          </div>
          <div className="text-2xl sm:text-3xl font-bold text-[#DC2626]">
            {analytics?.completionPercentage ?? (totalProblems ? Math.round((solved / totalProblems) * 100) : 0)}%
          </div>
          <span className="text-[11px] text-emerald-400 block">
            {solved} / {totalProblems} Solved
          </span>
        </div>
      </div>

      <div className="p-5 rounded-2xl bg-white border border-[#E2E8F0] space-y-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2 font-mono text-xs text-[#0F172A]">
            <Activity className="w-4 h-4 text-[#EF4444]" />
            <span>Problem-Solving Frequency (last 20 weeks)</span>
          </div>
        </div>
        <div className="overflow-x-auto pb-1">
          <div className="flex gap-1 min-w-[500px]">
            {activityCells.map((week, wIdx) => (
              <div key={wIdx} className="flex flex-col gap-1">
                {week.map((cell) => (
                  <div
                    key={cell.date}
                    title={`${cell.date}`}
                    className={`w-3 h-3 rounded-xs ${
                      cell.level === 3
                        ? 'bg-emerald-400'
                        : cell.level === 2
                        ? 'bg-[#EF4444]'
                        : cell.level === 1
                        ? 'bg-[#EF4444]/30'
                        : 'bg-[#E2E8F0]'
                    }`}
                  />
                ))}
              </div>
            ))}
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <div className="lg:col-span-8 space-y-6">
          <div className="p-6 rounded-2xl bg-white border border-[#E2E8F0] space-y-4">
            <div>
              <h3 className="text-sm font-bold text-[#0F172A] font-mono">Weekly Solving Volume</h3>
              <p className="text-xs text-[#64748B]">Verified solves aggregated by weekday</p>
            </div>
            <div className="h-56 w-full pt-2">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={analytics?.weeklyActivity || []}>
                  <XAxis dataKey="day" stroke="#5A5A6E" fontSize={11} tickLine={false} axisLine={{ stroke: '#E2E8F0' }} />
                  <YAxis stroke="#5A5A6E" fontSize={11} tickLine={false} axisLine={{ stroke: '#E2E8F0' }} />
                  <Tooltip
                    contentStyle={{
                      backgroundColor: '#FFFFFF',
                      border: '1px solid #CBD5E1',
                      borderRadius: '0.75rem',
                      color: '#F4F4F6',
                      fontSize: '11px',
                      fontFamily: 'JetBrains Mono, monospace',
                    }}
                    cursor={{ fill: 'rgba(139, 92, 246, 0.05)' }}
                  />
                  <Bar dataKey="solved" fill="#EF4444" radius={[4, 4, 0, 0]} />
                </BarChart>
              </ResponsiveContainer>
            </div>
          </div>

          <div className="p-6 rounded-2xl bg-white border border-[#E2E8F0] space-y-4">
            <div className="flex items-center justify-between">
              <span className="text-[10px] font-mono uppercase tracking-widest text-[#DC2626] bg-[#EF4444]/15 px-2.5 py-0.5 rounded-full border border-[#EF4444]/30">
                Active Pillar Module
              </span>
              <button
                onClick={() => navigate(`/roadmap/${activeTopic?.slug || activeTopic?.id || ''}`)}
                className="text-xs text-[#EF4444] hover:text-[#DC2626] flex items-center gap-1 font-mono"
              >
                Inspect Curriculum <ChevronRight className="w-3.5 h-3.5" />
              </button>
            </div>
            <div>
              <h3 className="text-xl font-bold text-[#0F172A] font-mono">{activeTopic?.title || 'Loading topic...'}</h3>
              <p className="text-xs text-[#64748B] mt-1">{activeTopic?.description}</p>
            </div>
            <div className="space-y-1.5 font-mono text-xs">
              <div className="flex justify-between">
                <span className="text-[#64748B]">Pillar Completion Rate</span>
                <span className="text-[#DC2626] font-bold">
                  {activeTopic?.completedProblems || 0} / {activeTopic?.totalProblems || 0} (
                  {Math.round(((activeTopic?.completedProblems || 0) / Math.max(1, activeTopic?.totalProblems || 1)) * 100)}%)
                </span>
              </div>
              <div className="w-full bg-[#E2E8F0] h-2 rounded-full overflow-hidden border border-[#E2E8F0]">
                <div
                  className="bg-[#EF4444] h-full rounded-full"
                  style={{
                    width: `${Math.min(100, Math.round(((activeTopic?.completedProblems || 0) / Math.max(1, activeTopic?.totalProblems || 1)) * 100))}%`,
                  }}
                />
              </div>
            </div>
          </div>
        </div>

        <div className="lg:col-span-4 space-y-6">
          <div className="p-5 rounded-2xl bg-[#F1F5F9] border border-[#EF4444]/30 space-y-3 glow-purple">
            <div className="flex items-center justify-between text-xs font-mono text-[#DC2626]">
              <span className="flex items-center gap-1.5">
                <Sparkles className="w-3.5 h-3.5 text-[#EF4444]" />
                Adaptive Diagnostic
              </span>
              <span className="text-[10px] text-[#94A3B8]">LIVE</span>
            </div>
            <h4 className="text-sm font-bold text-[#0F172A] font-mono">
              {latestMistake?.categoryLabel || 'Diagnostic telemetry ready'}
            </h4>
            <p className="text-xs text-[#64748B] leading-relaxed">
              {latestMistake?.explanation ||
                'Submit solutions in Practice to trigger test diagnostics, complexity notes, and edge-case analysis.'}
            </p>
            <Button
              variant="outline"
              size="sm"
              className="w-full text-xs font-mono"
              onClick={() => navigate(latestMistake ? '/analysis' : '/practice')}
            >
              {latestMistake ? 'Inspect Diagnostics' : 'Start Practice Session'}
            </Button>
          </div>

          <div className="space-y-3">
            <div className="flex items-center justify-between font-mono">
              <h3 className="text-xs font-bold text-[#0F172A] uppercase tracking-wider flex items-center gap-1.5">
                <Sparkles className="w-3.5 h-3.5 text-[#EF4444]" />
                <span>Recommended Queued</span>
              </h3>
              <button
                onClick={() => navigate('/recommendations')}
                className="text-xs text-[#EF4444] hover:text-[#DC2626] flex items-center gap-1 font-mono"
              >
                Explore Picks <ChevronRight className="w-3.5 h-3.5" />
              </button>
            </div>

            <div className="space-y-2">
              {smartRecommendations.length > 0 ? (
                smartRecommendations.map((rec) => (
                  <div
                    key={rec.id}
                    onClick={() =>
                      navigate(`/practice/${rec.problem_id}`, {
                        state: {
                          fromRecommendation: true,
                          recommendationId: rec.id,
                          category: rec.category,
                          primaryReason: rec.explanation,
                        }
                      })
                    }
                    className="p-3.5 rounded-xl bg-white hover:bg-[#F8FAFC] border border-[#E2E8F0] hover:border-[#EF4444]/40 cursor-pointer transition-all flex items-center justify-between group"
                  >
                    <div className="space-y-1 min-w-0 pr-2">
                      <div className="text-xs font-bold text-[#0F172A] group-hover:text-[#DC2626] truncate font-mono">
                        {rec.problem_title || rec.headline}
                      </div>
                      <p className="text-[11px] text-[#64748B] truncate">
                        {rec.explanation}
                      </p>
                      <div className="flex items-center gap-2 pt-0.5">
                        <span
                          className={`text-[9px] font-mono px-1.5 py-0.5 rounded ${
                            (rec.problem_difficulty || rec.recommended_difficulty) === 'Easy'
                              ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20'
                              : (rec.problem_difficulty || rec.recommended_difficulty) === 'Medium'
                              ? 'bg-amber-500/10 text-amber-400 border border-amber-500/20'
                              : 'bg-rose-500/10 text-rose-400 border border-rose-500/20'
                          }`}
                        >
                          {(rec.problem_difficulty || rec.recommended_difficulty)}
                        </span>
                        <span className="text-[9px] font-mono text-[#DC2626] bg-[#EF4444]/10 px-1.5 py-0.5 rounded">
                          {normalizePercentage(rec.match_score)}% Match
                        </span>
                      </div>
                    </div>
                    <ChevronRight className="w-4 h-4 text-[#94A3B8] group-hover:text-[#0F172A] shrink-0" />
                  </div>
                ))
              ) : recommendedProblems.length > 0 ? (
                recommendedProblems.map((prob) => (
                  <div
                    key={prob.id}
                    onClick={() => navigate(`/practice/${prob.id}`)}
                    className="p-3.5 rounded-xl bg-white hover:bg-[#F8FAFC] border border-[#E2E8F0] hover:border-[#2D2D3E] cursor-pointer transition-all flex items-center justify-between group"
                  >
                    <div className="space-y-1 min-w-0">
                      <div className="text-xs font-bold text-[#0F172A] group-hover:text-[#DC2626] truncate font-mono">
                        {prob.title}
                      </div>
                      <div className="flex items-center gap-2">
                        <span
                          className={`text-[10px] font-mono px-2 py-0.5 rounded ${
                            prob.difficulty === 'Easy'
                              ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20'
                              : 'bg-amber-500/10 text-amber-400 border border-amber-500/20'
                          }`}
                        >
                          {prob.difficulty}
                        </span>
                        <span className="text-[10px] text-[#94A3B8] font-mono">{prob.progress_status}</span>
                      </div>
                    </div>
                    <ChevronRight className="w-4 h-4 text-[#94A3B8] group-hover:text-[#0F172A] shrink-0" />
                  </div>
                ))
              ) : (
                <p className="text-xs text-[#94A3B8] font-mono p-4 text-center">Loading queued problems…</p>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default DashboardPage;
