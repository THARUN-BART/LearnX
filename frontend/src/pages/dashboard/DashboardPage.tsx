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
import { AnalyticsSummary, Problem, Topic, MistakeRecord } from '../../types/learning.types';
import { LearningRecommendation } from '../../types/learning-path.types';
import { RecommendedActionCard } from '../../components/learning-path/RecommendedActionCard';
import Button from '../../components/ui/Button';
import { PageSkeleton } from '../../components/ui/Skeleton';

export const DashboardPage: React.FC = () => {
  const { user } = useAuth();
  const { profile } = useProfile();
  const navigate = useNavigate();

  const [analytics, setAnalytics] = useState<AnalyticsSummary | null>(null);
  const [recommendedProblems, setRecommendedProblems] = useState<Problem[]>([]);
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
      const [problems, mistakes] = await Promise.all([
        learningService.getProblems({ status: 'Not Started', limit: 8 }),
        analyticsService.getMistakes(),
      ]);
      if (cancelled) return;
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
      <div className="p-6 sm:p-7 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] relative overflow-hidden">
        <div className="absolute inset-0 bg-tech-grid opacity-30 pointer-events-none" />
        <div className="relative z-10 flex flex-col lg:flex-row lg:items-center justify-between gap-6">
          <div className="space-y-2">
            <div className="flex flex-wrap items-center gap-2 font-mono text-[11px]">
              <span className="px-2 py-0.5 rounded bg-[#8B5CF6]/15 text-[#A78BFA] border border-[#8B5CF6]/30">
                HOST: LEARNX-ENGINE
              </span>
              <span className="text-[#5A5A6E]">//</span>
              <span className="text-emerald-400 flex items-center gap-1">
                <span className="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse" />
                ONLINE
              </span>
              <span className="text-[#5A5A6E]">//</span>
              <span className="text-[#8E8EA0]">
                GOAL: {(profile?.learning_goal || 'INTERVIEW PREPARATION').toUpperCase()}
              </span>
            </div>
            <h1 className="text-2xl sm:text-3xl font-bold tracking-tight text-[#F4F4F6] font-mono">
              Welcome back, {studentName}
            </h1>
            <p className="text-xs text-[#8E8EA0] max-w-2xl leading-relaxed">
              Current focus: <strong className="text-white">{activeTopic?.title || 'Striver A2Z roadmap'}</strong>.
              {analytics?.needsRevision
                ? ` ${analytics.needsRevision} problems are marked for revision.`
                : ' Pick up the next unsolved problem from your active pillar.'}
            </p>
          </div>
          <div className="flex items-center gap-3 shrink-0">
            <Button variant="outline" size="sm" onClick={() => navigate('/roadmap')} leftIcon={<Layers className="w-4 h-4" />}>
              Roadmap
            </Button>
            <Button
              variant="primary"
              size="sm"
              onClick={() => navigate(recommendedProblems[0] ? `/practice/${recommendedProblems[0].id}` : '/practice')}
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
          <div className="flex items-center gap-2 font-mono text-xs text-[#8E8EA0] uppercase tracking-wider font-semibold">
            <Compass className="w-4 h-4 text-[#8B5CF6]" />
            <span>AI Adaptive Path Engine</span>
          </div>
          <button
            onClick={() => navigate('/dashboard/learning-path')}
            className="text-xs text-[#8B5CF6] hover:text-[#A78BFA] flex items-center gap-1 font-mono"
          >
            Explore Full Roadmap <ChevronRight className="w-3.5 h-3.5" />
          </button>
        </div>
        <RecommendedActionCard recommendation={adaptiveNextAction} loading={loading} />
      </div>

      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="p-5 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-1 font-mono">
          <div className="flex items-center justify-between text-xs text-[#5A5A6E]">
            <span>VERIFIED_SOLVED</span>
            <CheckCircle2 className="w-4 h-4 text-emerald-400" />
          </div>
          <div className="text-2xl sm:text-3xl font-bold text-[#F4F4F6]">{solved}</div>
          <span className="text-[11px] text-emerald-400 block">
            {analytics?.attemptedProblems ?? 0} attempted · {analytics?.pendingProblems ?? totalProblems} pending
          </span>
        </div>
        <div className="p-5 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-1 font-mono">
          <div className="flex items-center justify-between text-xs text-[#5A5A6E]">
            <span>STREAK_STABILITY</span>
            <Flame className="w-4 h-4 text-amber-400" />
          </div>
          <div className="text-2xl sm:text-3xl font-bold text-amber-400">{analytics?.streakDays ?? 0}d</div>
          <span className="text-[11px] text-[#8E8EA0] block">
            {(analytics?.streakDays ?? 0) > 0 ? 'Consecutive active days' : 'Solve today to start a streak'}
          </span>
        </div>
        <div className="p-5 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-1 font-mono">
          <div className="flex items-center justify-between text-xs text-[#5A5A6E]">
            <span>ACTIVE_RUNTIME</span>
            <Clock className="w-4 h-4 text-[#818CF8]" />
          </div>
          <div className="text-2xl sm:text-3xl font-bold text-[#818CF8]">{analytics?.hoursSpent ?? 0}h</div>
          <span className="text-[11px] text-[#8E8EA0] block">Estimated practice time</span>
        </div>
        <div className="p-5 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-1 font-mono">
          <div className="flex items-center justify-between text-xs text-[#5A5A6E]">
            <span>MASTERY_INDEX</span>
            <TrendingUp className="w-4 h-4 text-[#A78BFA]" />
          </div>
          <div className="text-2xl sm:text-3xl font-bold text-[#A78BFA]">
            {analytics?.completionPercentage ?? (totalProblems ? Math.round((solved / totalProblems) * 100) : 0)}%
          </div>
          <span className="text-[11px] text-emerald-400 block">
            {solved} / {totalProblems} Solved
          </span>
        </div>
      </div>

      <div className="p-5 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2 font-mono text-xs text-[#F4F4F6]">
            <Activity className="w-4 h-4 text-[#8B5CF6]" />
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
                        ? 'bg-[#8B5CF6]'
                        : cell.level === 1
                        ? 'bg-[#8B5CF6]/30'
                        : 'bg-[#181824]'
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
          <div className="p-6 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-4">
            <div>
              <h3 className="text-sm font-bold text-[#F4F4F6] font-mono">Weekly Solving Volume</h3>
              <p className="text-xs text-[#8E8EA0]">Verified solves aggregated by weekday</p>
            </div>
            <div className="h-56 w-full pt-2">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={analytics?.weeklyActivity || []}>
                  <XAxis dataKey="day" stroke="#5A5A6E" fontSize={11} tickLine={false} axisLine={{ stroke: '#1C1C28' }} />
                  <YAxis stroke="#5A5A6E" fontSize={11} tickLine={false} axisLine={{ stroke: '#1C1C28' }} />
                  <Tooltip
                    contentStyle={{
                      backgroundColor: '#0D0D14',
                      border: '1px solid #2B2B3C',
                      borderRadius: '0.75rem',
                      color: '#F4F4F6',
                      fontSize: '11px',
                      fontFamily: 'JetBrains Mono, monospace',
                    }}
                    cursor={{ fill: 'rgba(139, 92, 246, 0.05)' }}
                  />
                  <Bar dataKey="solved" fill="#8B5CF6" radius={[4, 4, 0, 0]} />
                </BarChart>
              </ResponsiveContainer>
            </div>
          </div>

          <div className="p-6 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-4">
            <div className="flex items-center justify-between">
              <span className="text-[10px] font-mono uppercase tracking-widest text-[#A78BFA] bg-[#8B5CF6]/15 px-2.5 py-0.5 rounded-full border border-[#8B5CF6]/30">
                Active Pillar Module
              </span>
              <button
                onClick={() => navigate(`/roadmap/${activeTopic?.slug || activeTopic?.id || ''}`)}
                className="text-xs text-[#8B5CF6] hover:text-[#A78BFA] flex items-center gap-1 font-mono"
              >
                Inspect Curriculum <ChevronRight className="w-3.5 h-3.5" />
              </button>
            </div>
            <div>
              <h3 className="text-xl font-bold text-[#F4F4F6] font-mono">{activeTopic?.title || 'Loading topic...'}</h3>
              <p className="text-xs text-[#8E8EA0] mt-1">{activeTopic?.description}</p>
            </div>
            <div className="space-y-1.5 font-mono text-xs">
              <div className="flex justify-between">
                <span className="text-[#8E8EA0]">Pillar Completion Rate</span>
                <span className="text-[#A78BFA] font-bold">
                  {activeTopic?.completedProblems || 0} / {activeTopic?.totalProblems || 0} (
                  {Math.round(((activeTopic?.completedProblems || 0) / Math.max(1, activeTopic?.totalProblems || 1)) * 100)}%)
                </span>
              </div>
              <div className="w-full bg-[#181824] h-2 rounded-full overflow-hidden border border-[#1C1C28]">
                <div
                  className="bg-[#8B5CF6] h-full rounded-full"
                  style={{
                    width: `${Math.min(100, Math.round(((activeTopic?.completedProblems || 0) / Math.max(1, activeTopic?.totalProblems || 1)) * 100))}%`,
                  }}
                />
              </div>
            </div>
          </div>
        </div>

        <div className="lg:col-span-4 space-y-6">
          <div className="p-5 rounded-2xl bg-[#12121E] border border-[#8B5CF6]/30 space-y-3 glow-purple">
            <div className="flex items-center justify-between text-xs font-mono text-[#A78BFA]">
              <span className="flex items-center gap-1.5">
                <Sparkles className="w-3.5 h-3.5 text-[#8B5CF6]" />
                Adaptive Diagnostic
              </span>
              <span className="text-[10px] text-[#5A5A6E]">LIVE</span>
            </div>
            <h4 className="text-sm font-bold text-[#F4F4F6] font-mono">
              {latestMistake?.categoryLabel || 'Diagnostic telemetry ready'}
            </h4>
            <p className="text-xs text-[#8E8EA0] leading-relaxed">
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
              <h3 className="text-xs font-bold text-[#F4F4F6] uppercase tracking-wider">Recommended Queued</h3>
              <button onClick={() => navigate('/practice')} className="text-xs text-[#8B5CF6] hover:text-[#A78BFA]">
                View all
              </button>
            </div>
            <div className="space-y-2">
              {recommendedProblems.length === 0 && (
                <p className="text-xs text-[#5A5A6E] font-mono">Loading next problems…</p>
              )}
              {recommendedProblems.map((prob) => (
                <div
                  key={prob.id}
                  onClick={() => navigate(`/practice/${prob.id}`)}
                  className="p-3.5 rounded-xl bg-[#0D0D14] hover:bg-[#14141E] border border-[#1C1C28] hover:border-[#2D2D3E] cursor-pointer transition-all flex items-center justify-between group"
                >
                  <div className="space-y-1 min-w-0">
                    <div className="text-xs font-bold text-[#F4F4F6] group-hover:text-[#A78BFA] truncate font-mono">
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
                      <span className="text-[10px] text-[#5A5A6E] font-mono">{prob.progress_status}</span>
                    </div>
                  </div>
                  <ChevronRight className="w-4 h-4 text-[#5A5A6E] group-hover:text-white" />
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default DashboardPage;
