import React, { useState, useEffect } from 'react';
import {
  TrendingUp,
  Clock,
  Flame,
  CheckCircle2,
  Calendar,
  Layers,
  Award,
  Zap,
  BarChart2,
  Activity,
  Cpu
} from 'lucide-react';
import {
  ResponsiveContainer,
  AreaChart,
  Area,
  BarChart,
  Bar,
  PieChart,
  Pie,
  Cell,
  XAxis,
  YAxis,
  Tooltip
} from 'recharts';
import { AnalyticsSummary } from '../../types/learning.types';
import analyticsService from '../../services/analytics.service';
import Card from '../../components/ui/Card';
import Badge from '../../components/ui/Badge';

export const AnalyticsPage: React.FC = () => {
  const [stats, setStats] = useState<AnalyticsSummary | null>(null);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const loadStats = async () => {
      const data = await analyticsService.getAnalytics();
      setStats(data);
      setLoading(false);
    };
    loadStats();
  }, []);

  if (loading || !stats) {
    return (
      <div className="py-20 text-center text-xs font-mono text-[#5A5A6E]">
        Aggregating performance telemetry...
      </div>
    );
  }

  return (
    <div className="space-y-8 font-mono">
      {/* Page Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl sm:text-3xl font-bold tracking-tight text-[#F4F4F6] flex items-center gap-2">
            <BarChart2 className="w-6 h-6 text-[#8B5CF6]" />
            <span>Telemetry & Performance Metrics</span>
          </h1>
          <p className="text-xs text-[#8E8EA0] mt-1 font-sans">
            Asymptotic runtime statistics, testsuite pass accuracy, and topic competency matrices.
          </p>
        </div>

        <div className="flex items-center gap-2 px-3 py-1 rounded-full bg-[#12121B] border border-[#1C1C28] text-xs text-[#8E8EA0]">
          <span className="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse" />
          <span>Real-time Sampling Engine</span>
        </div>
      </div>

      {/* High-level KPI Strip */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="p-5 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-1">
          <span className="text-xs text-[#5A5A6E]">TOTAL_VERIFIED</span>
          <div className="text-3xl font-bold text-[#F4F4F6]">
            {stats.totalSolved}
          </div>
          <span className="text-[11px] text-emerald-400">
            Accuracy: {stats.accuracyRate}%
          </span>
        </div>

        <div className="p-5 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-1">
          <span className="text-xs text-[#5A5A6E]">AVG_SOLVE_DURATION</span>
          <div className="text-3xl font-bold text-[#818CF8]">
            {stats.totalSolved > 0 ? `${Math.round((stats.hoursSpent * 60) / stats.totalSolved)}m` : '0m'}
          </div>
          <span className="text-[11px] text-[#8E8EA0]">
            {stats.totalSolved > 0 ? 'Average solve duration' : 'No recorded sessions'}
          </span>
        </div>

        <div className="p-5 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-1">
          <span className="text-xs text-[#5A5A6E]">ACTIVE_STREAK</span>
          <div className="text-3xl font-bold text-amber-400">
            {stats.streakDays} Days
          </div>
          <span className="text-[11px] text-amber-400/80">Unbroken consistency</span>
        </div>

        <div className="p-5 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-1">
          <span className="text-xs text-[#5A5A6E]">TOTAL_HOURS</span>
          <div className="text-3xl font-bold text-[#A78BFA]">
            {stats.hoursSpent}h
          </div>
          <span className="text-[11px] text-[#8E8EA0]">Deep IDE runtime</span>
        </div>
      </div>

      {/* TOPIC ACCURACY & DIFFICULTY SPLIT */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        {/* Topic-Wise Proficiency Bar Chart */}
        <div className="lg:col-span-8 p-6 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-4">
          <div className="flex items-center justify-between">
            <div>
              <h3 className="text-sm font-bold text-[#F4F4F6]">
                Topic-Wise Accuracy & Proficiency Index
              </h3>
              <p className="text-xs text-[#8E8EA0] font-sans">
                Computed from first-attempt clean passes without hints
              </p>
            </div>
            <span className="text-[10px] px-2 py-0.5 rounded bg-[#181824] text-[#8E8EA0] border border-[#2B2B3C]">
              Scaled 0-100
            </span>
          </div>

          <div className="h-64 w-full pt-2 flex items-center justify-center">
            {stats.topicProficiency.length === 0 ? (
              <span className="text-xs text-[#5A5A6E] font-mono">
                No topic activity recorded yet. Start solving problems to populate competency scores.
              </span>
            ) : (
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={stats.topicProficiency} layout="vertical">
                  <XAxis type="number" domain={[0, 100]} stroke="#5A5A6E" fontSize={11} />
                  <YAxis dataKey="topic" type="category" stroke="#8E8EA0" fontSize={11} width={80} />
                  <Tooltip
                    contentStyle={{
                      backgroundColor: '#0D0D14',
                      border: '1px solid #2B2B3C',
                      borderRadius: '0.75rem',
                      color: '#F4F4F6',
                      fontSize: '11px',
                      fontFamily: 'JetBrains Mono, monospace'
                    }}
                  />
                  <Bar dataKey="score" fill="#8B5CF6" radius={[0, 4, 4, 0]} />
                </BarChart>
              </ResponsiveContainer>
            )}
          </div>
        </div>

        {/* Difficulty Distribution Pie */}
        <div className="lg:col-span-4 p-6 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-4 flex flex-col justify-between">
          <div>
            <h3 className="text-sm font-bold text-[#F4F4F6]">
              Problem Difficulty Split
            </h3>
            <p className="text-xs text-[#8E8EA0] font-sans">Curriculum tier breakdown</p>
          </div>

          <div className="h-44 w-full">
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie
                  data={stats.difficultyDistribution}
                  innerRadius={35}
                  outerRadius={65}
                  paddingAngle={5}
                  dataKey="count"
                >
                  {stats.difficultyDistribution.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={entry.color} />
                  ))}
                </Pie>
                <Tooltip
                  contentStyle={{
                    backgroundColor: '#0D0D14',
                    border: '1px solid #2B2B3C',
                    borderRadius: '0.75rem',
                    color: '#F4F4F6',
                    fontSize: '11px',
                    fontFamily: 'JetBrains Mono, monospace'
                  }}
                />
              </PieChart>
            </ResponsiveContainer>
          </div>

          <div className="grid grid-cols-3 gap-2 text-center text-xs pt-3 border-t border-[#1C1C28]">
            <div>
              <span className="text-emerald-400 font-bold block">{stats.easySolved}</span>
              <span className="text-[10px] text-[#5A5A6E]">EASY</span>
            </div>
            <div>
              <span className="text-amber-400 font-bold block">{stats.mediumSolved}</span>
              <span className="text-[10px] text-[#5A5A6E]">MEDIUM</span>
            </div>
            <div>
              <span className="text-red-400 font-bold block">{stats.hardSolved}</span>
              <span className="text-[10px] text-[#5A5A6E]">HARD</span>
            </div>
          </div>
        </div>
      </div>

      {/* MONTHLY GROWTH TIMELINE & STRENGTH/WEAKNESS */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        {/* Monthly Cumulative Progress Area Chart */}
        <div className="lg:col-span-8 p-6 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-4">
          <div>
            <h3 className="text-sm font-bold text-[#F4F4F6]">
              Monthly Cumulative Velocity
            </h3>
            <p className="text-xs text-[#8E8EA0] font-sans">Volume of verified problem completions over time</p>
          </div>

          <div className="h-56 w-full pt-2">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={stats.monthlyProgress}>
                <defs>
                  <linearGradient id="colorSolved" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#8B5CF6" stopOpacity={0.3} />
                    <stop offset="95%" stopColor="#8B5CF6" stopOpacity={0.0} />
                  </linearGradient>
                </defs>
                <XAxis dataKey="month" stroke="#5A5A6E" fontSize={11} tickLine={false} />
                <YAxis stroke="#5A5A6E" fontSize={11} tickLine={false} />
                <Tooltip
                  contentStyle={{
                    backgroundColor: '#0D0D14',
                    border: '1px solid #2B2B3C',
                    borderRadius: '0.75rem',
                    color: '#F4F4F6',
                    fontSize: '11px',
                    fontFamily: 'JetBrains Mono, monospace'
                  }}
                />
                <Area
                  type="monotone"
                  dataKey="solved"
                  stroke="#8B5CF6"
                  strokeWidth={2}
                  fillOpacity={1}
                  fill="url(#colorSolved)"
                />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Strongest vs Weakest Topics */}
        <div className="lg:col-span-4 p-6 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-4">
          <h3 className="text-sm font-bold text-[#F4F4F6]">
            Competency Assessment
          </h3>

          <div className="space-y-4 text-xs font-sans">
            {/* Strongest */}
            <div className="p-3.5 bg-emerald-500/10 border border-emerald-500/20 rounded-xl space-y-1">
              <span className="text-emerald-400 font-bold block flex items-center gap-1.5 font-mono">
                <CheckCircle2 className="w-3.5 h-3.5" /> High Mastery Index
              </span>
              <div className="text-sm font-semibold text-[#F4F4F6] font-mono">
                Hashing & Frequency Maps (90%)
              </div>
              <p className="text-[#8E8EA0] text-[11px]">
                Consistently selects O(1) expected time structures with minimal collision overhead.
              </p>
            </div>

            {/* Weakest */}
            <div className="p-3.5 bg-amber-500/10 border border-amber-500/20 rounded-xl space-y-1">
              <span className="text-amber-400 font-bold block flex items-center gap-1.5 font-mono">
                <Flame className="w-3.5 h-3.5" /> Queued for Drill
              </span>
              <div className="text-sm font-semibold text-[#F4F4F6] font-mono">
                Dynamic Programming (35%)
              </div>
              <p className="text-[#8E8EA0] text-[11px]">
                Frequent recursion depth exceedances; reinforce 2D state memoization tables.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AnalyticsPage;
