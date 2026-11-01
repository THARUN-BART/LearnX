import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Flame,
  CheckCircle2,
  Clock,
  TrendingUp,
  ArrowRight,
  Sparkles,
  Play,
  Bot,
  BrainCircuit,
  AlertTriangle,
  ChevronRight,
  Terminal,
  Activity,
  Cpu,
  Layers
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
import { AnalyticsSummary, Problem } from '../../types/learning.types';
import Card from '../../components/ui/Card';
import Badge from '../../components/ui/Badge';
import Button from '../../components/ui/Button';

export const DashboardPage: React.FC = () => {
  const { user } = useAuth();
  const { profile } = useProfile();
  const navigate = useNavigate();

  const [analytics, setAnalytics] = useState<AnalyticsSummary | null>(null);
  const [recommendedProblems, setRecommendedProblems] = useState<Problem[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const loadDashboardData = async () => {
      try {
        const [stats, problems] = await Promise.all([
          analyticsService.getAnalytics(),
          learningService.getProblems()
        ]);
        setAnalytics(stats);
        setRecommendedProblems(problems.slice(0, 3));
      } catch (e) {
        console.error(e);
      } finally {
        setLoading(false);
      }
    };
    loadDashboardData();
  }, []);

  const studentName = user?.fullName || user?.full_name || 'Student';

  // Generate simulated 20-week activity heatmap blocks
  const activityWeeks = Array.from({ length: 20 }, (_, w) => {
    return Array.from({ length: 7 }, (_, d) => {
      const activeProb = Math.random();
      const level = activeProb > 0.7 ? 3 : activeProb > 0.4 ? 2 : activeProb > 0.2 ? 1 : 0;
      return { week: w, day: d, level };
    });
  });

  return (
    <div className="space-y-8">
      {/* Top Cockpit Telemetry Banner */}
      <div className="p-6 sm:p-7 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] relative overflow-hidden">
        {/* Subtle grid texture overlay */}
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
              Curriculum focus is locked on <strong className="text-white">Arrays & Two Pointers</strong>. You have 2 high-priority edge-case validations queued for verification.
            </p>
          </div>

          <div className="flex items-center gap-3 shrink-0">
            <Button
              variant="outline"
              size="sm"
              onClick={() => navigate('/roadmap')}
              leftIcon={<Layers className="w-4 h-4" />}
            >
              Roadmap
            </Button>
            <Button
              variant="primary"
              size="sm"
              onClick={() => navigate('/practice')}
              rightIcon={<Play className="w-4 h-4" />}
            >
              Resume Coding
            </Button>
          </div>
        </div>
      </div>

      {/* KPI METRIC TILES */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        {/* Metric 1 */}
        <div className="p-5 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-1 font-mono">
          <div className="flex items-center justify-between text-xs text-[#5A5A6E]">
            <span>VERIFIED_SOLVED</span>
            <CheckCircle2 className="w-4 h-4 text-emerald-400" />
          </div>
          <div className="text-2xl sm:text-3xl font-bold text-[#F4F4F6]">
            {analytics?.totalSolved || profile?.problemsSolved || 47}
          </div>
          <span className="text-[11px] text-emerald-400 block">+6 this week</span>
        </div>

        {/* Metric 2 */}
        <div className="p-5 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-1 font-mono">
          <div className="flex items-center justify-between text-xs text-[#5A5A6E]">
            <span>STREAK_STABILITY</span>
            <Flame className="w-4 h-4 text-amber-400" />
          </div>
          <div className="text-2xl sm:text-3xl font-bold text-amber-400">
            {analytics?.streakDays || profile?.streakDays || 14}d
          </div>
          <span className="text-[11px] text-[#8E8EA0] block">Top 5% consistency</span>
        </div>

        {/* Metric 3 */}
        <div className="p-5 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-1 font-mono">
          <div className="flex items-center justify-between text-xs text-[#5A5A6E]">
            <span>ACTIVE_RUNTIME</span>
            <Clock className="w-4 h-4 text-[#818CF8]" />
          </div>
          <div className="text-2xl sm:text-3xl font-bold text-[#818CF8]">
            {analytics?.hoursSpent || profile?.practiceHours || 38}h
          </div>
          <span className="text-[11px] text-[#8E8EA0] block">1.8h / daily avg</span>
        </div>

        {/* Metric 4 */}
        <div className="p-5 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-1 font-mono">
          <div className="flex items-center justify-between text-xs text-[#5A5A6E]">
            <span>MASTERY_INDEX</span>
            <TrendingUp className="w-4 h-4 text-[#A78BFA]" />
          </div>
          <div className="text-2xl sm:text-3xl font-bold text-[#A78BFA]">
            {profile?.overallProgress || 42}%
          </div>
          <span className="text-[11px] text-emerald-400 block">3 / 16 Pillars Mastered</span>
        </div>
      </div>

      {/* HEATMAP CONSISTENCY MATRIX */}
      <div className="p-5 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2 font-mono text-xs text-[#F4F4F6]">
            <Activity className="w-4 h-4 text-[#8B5CF6]" />
            <span>Problem-Solving Frequency Heatmap</span>
          </div>
          <div className="flex items-center gap-1.5 text-[10px] font-mono text-[#5A5A6E]">
            <span>Less</span>
            <span className="w-2.5 h-2.5 rounded-xs bg-[#1C1C28]" />
            <span className="w-2.5 h-2.5 rounded-xs bg-[#8B5CF6]/40" />
            <span className="w-2.5 h-2.5 rounded-xs bg-[#8B5CF6]" />
            <span className="w-2.5 h-2.5 rounded-xs bg-emerald-400" />
            <span>More</span>
          </div>
        </div>

        {/* Matrix grid */}
        <div className="overflow-x-auto pb-1">
          <div className="flex gap-1 min-w-[500px]">
            {activityWeeks.map((week, wIdx) => (
              <div key={wIdx} className="flex flex-col gap-1">
                {week.map((cell, dIdx) => (
                  <div
                    key={dIdx}
                    title={`Week ${cell.week}, Day ${cell.day}`}
                    className={`w-3 h-3 rounded-xs transition-colors ${
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

      {/* TWO COLUMN CONTENT: Left (Charts & Topics), Right (AI Insights & Problems) */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        {/* Left Column (8 cols) */}
        <div className="lg:col-span-8 space-y-6">
          {/* Weekly Activity Bar Chart */}
          <div className="p-6 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <h3 className="text-sm font-bold text-[#F4F4F6] font-mono">
                  Weekly Solving Volume
                </h3>
                <p className="text-xs text-[#8E8EA0]">
                  Number of verified submissions per 24-hour cycle
                </p>
              </div>
              <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-[#181824] text-[#8E8EA0] border border-[#2B2B3C]">
                Telemetry Window: 7d
              </span>
            </div>

            <div className="h-56 w-full pt-2">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={analytics?.weeklyActivity || []}>
                  <XAxis
                    dataKey="day"
                    stroke="#5A5A6E"
                    fontSize={11}
                    tickLine={false}
                    axisLine={{ stroke: '#1C1C28' }}
                  />
                  <YAxis
                    stroke="#5A5A6E"
                    fontSize={11}
                    tickLine={false}
                    axisLine={{ stroke: '#1C1C28' }}
                  />
                  <Tooltip
                    contentStyle={{
                      backgroundColor: '#0D0D14',
                      border: '1px solid #2B2B3C',
                      borderRadius: '0.75rem',
                      color: '#F4F4F6',
                      fontSize: '11px',
                      fontFamily: 'JetBrains Mono, monospace'
                    }}
                    cursor={{ fill: 'rgba(139, 92, 246, 0.05)' }}
                  />
                  <Bar
                    dataKey="solved"
                    fill="#8B5CF6"
                    radius={[4, 4, 0, 0]}
                  />
                </BarChart>
              </ResponsiveContainer>
            </div>
          </div>

          {/* Active Learning Topic Card */}
          <div className="p-6 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-4 relative overflow-hidden">
            <div className="flex items-center justify-between">
              <span className="text-[10px] font-mono uppercase tracking-widest text-[#A78BFA] bg-[#8B5CF6]/15 px-2.5 py-0.5 rounded-full border border-[#8B5CF6]/30">
                Active Pillar Module
              </span>
              <button
                onClick={() => navigate('/roadmap/arrays')}
                className="text-xs text-[#8B5CF6] hover:text-[#A78BFA] flex items-center gap-1 font-mono transition-colors"
              >
                <span>Inspect Curriculum</span>
                <ChevronRight className="w-3.5 h-3.5" />
              </button>
            </div>

            <div>
              <h3 className="text-xl font-bold text-[#F4F4F6] font-mono">
                Arrays, Sliders & Two-Pointers
              </h3>
              <p className="text-xs text-[#8E8EA0] mt-1">
                Contiguous memory layout, sliding window bounds, Kadane’s maximum subarray routine, and fast-slow pointer logic.
              </p>
            </div>

            <div className="space-y-1.5 font-mono text-xs">
              <div className="flex justify-between">
                <span className="text-[#8E8EA0]">Pillar Completion Rate</span>
                <span className="text-[#A78BFA] font-bold">18 / 24 Problems (75%)</span>
              </div>
              <div className="w-full bg-[#181824] h-2 rounded-full overflow-hidden border border-[#1C1C28]">
                <div className="bg-[#8B5CF6] h-full w-[75%] rounded-full" />
              </div>
            </div>
          </div>
        </div>

        {/* Right Column (4 cols): AI Insights & Recommended Problems */}
        <div className="lg:col-span-4 space-y-6">
          {/* AI Diagnostic Alert Card */}
          <div className="p-5 rounded-2xl bg-[#12121E] border border-[#8B5CF6]/30 space-y-3 relative glow-purple">
            <div className="flex items-center justify-between text-xs font-mono text-[#A78BFA]">
              <span className="flex items-center gap-1.5">
                <Sparkles className="w-3.5 h-3.5 text-[#8B5CF6]" />
                Adaptive Diagnostic
              </span>
              <span className="text-[10px] text-[#5A5A6E]">LIVE</span>
            </div>

            <h4 className="text-sm font-bold text-[#F4F4F6] font-mono">
              Avoid Off-by-one Edge Slips
            </h4>

            <p className="text-xs text-[#8E8EA0] leading-relaxed">
              In your last 5 submissions, 40% of panics were caused by loop inclusive bounds (`high = mid` instead of `high = mid - 1`).
            </p>

            <Button
              variant="outline"
              size="sm"
              className="w-full text-xs font-mono"
              onClick={() => navigate('/analysis')}
            >
              Inspect Diagnostics
            </Button>
          </div>

          {/* Recommended Problems */}
          <div className="space-y-3">
            <div className="flex items-center justify-between font-mono">
              <h3 className="text-xs font-bold text-[#F4F4F6] uppercase tracking-wider">
                Recommended Queued
              </h3>
              <button
                onClick={() => navigate('/practice')}
                className="text-xs text-[#8B5CF6] hover:text-[#A78BFA]"
              >
                View all
              </button>
            </div>

            <div className="space-y-2">
              {recommendedProblems.map((prob) => (
                <div
                  key={prob.id}
                  onClick={() => navigate(`/practice/${prob.id}`)}
                  className="p-3.5 rounded-xl bg-[#0D0D14] hover:bg-[#14141E] border border-[#1C1C28] hover:border-[#2D2D3E] cursor-pointer transition-all flex items-center justify-between group"
                >
                  <div className="space-y-1 min-w-0">
                    <div className="text-xs font-bold text-[#F4F4F6] group-hover:text-[#A78BFA] transition-colors truncate font-mono">
                      {prob.title}
                    </div>
                    <div className="flex items-center gap-2">
                      <span className={`text-[10px] font-mono px-2 py-0.5 rounded ${
                        prob.difficulty === 'Easy'
                          ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20'
                          : 'bg-amber-500/10 text-amber-400 border border-amber-500/20'
                      }`}>
                        {prob.difficulty}
                      </span>
                      <span className="text-[10px] text-[#5A5A6E] font-mono">
                        {prob.accuracyRate}% pass
                      </span>
                    </div>
                  </div>
                  <ChevronRight className="w-4 h-4 text-[#5A5A6E] group-hover:text-white group-hover:translate-x-1 transition-all" />
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
