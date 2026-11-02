import React, { useState, useEffect } from 'react';
import {
  CheckCircle2,
  Flame,
  BarChart2,
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
import { PageSkeleton } from '../../components/ui/Skeleton';

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
    return <PageSkeleton rows={3} />;
  }

  const pieData = stats.difficultyDistribution.filter((d) => d.count > 0);
  const daily = stats.dailyActivity.slice(-30);

  return (
    <div className="space-y-8 font-mono">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl sm:text-3xl font-bold tracking-tight text-[#F4F4F6] flex items-center gap-2">
            <BarChart2 className="w-6 h-6 text-[#8B5CF6]" />
            Telemetry & Performance Metrics
          </h1>
          <p className="text-xs text-[#8E8EA0] mt-1 font-sans">
            Live progress from your A2Z sheet — {stats.totalSolved} / {stats.totalProblems} solved
            {stats.needsRevision ? ` · ${stats.needsRevision} need revision` : ''}.
          </p>
        </div>
      </div>

      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="p-5 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-1">
          <span className="text-xs text-[#5A5A6E]">TOTAL_VERIFIED</span>
          <div className="text-3xl font-bold text-[#F4F4F6]">{stats.totalSolved}</div>
          <span className="text-[11px] text-emerald-400">Accuracy: {stats.accuracyRate}%</span>
        </div>
        <div className="p-5 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-1">
          <span className="text-xs text-[#5A5A6E]">COMPLETION</span>
          <div className="text-3xl font-bold text-[#818CF8]">{stats.completionPercentage}%</div>
          <span className="text-[11px] text-[#8E8EA0]">{stats.pendingProblems} still pending</span>
        </div>
        <div className="p-5 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-1">
          <span className="text-xs text-[#5A5A6E]">ACTIVE_STREAK</span>
          <div className="text-3xl font-bold text-amber-400">{stats.streakDays} Days</div>
          <span className="text-[11px] text-amber-400/80">Consecutive active days</span>
        </div>
        <div className="p-5 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-1">
          <span className="text-xs text-[#5A5A6E]">TOTAL_HOURS</span>
          <div className="text-3xl font-bold text-[#A78BFA]">{stats.hoursSpent}h</div>
          <span className="text-[11px] text-[#8E8EA0]">{stats.attemptedProblems} attempts logged</span>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <div className="lg:col-span-8 p-6 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-4">
          <h3 className="text-sm font-bold text-[#F4F4F6]">Topic-Wise Completion</h3>
          <div className="h-64 w-full pt-2 flex items-center justify-center">
            {stats.topicProficiency.length === 0 ? (
              <span className="text-xs text-[#5A5A6E]">No topic activity yet. Solve problems to populate scores.</span>
            ) : (
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={stats.topicProficiency} layout="vertical">
                  <XAxis type="number" domain={[0, 100]} stroke="#5A5A6E" fontSize={11} />
                  <YAxis dataKey="topic" type="category" stroke="#8E8EA0" fontSize={11} width={110} />
                  <Tooltip
                    contentStyle={{
                      backgroundColor: '#0D0D14',
                      border: '1px solid #2B2B3C',
                      borderRadius: '0.75rem',
                      color: '#F4F4F6',
                      fontSize: '11px',
                    }}
                  />
                  <Bar dataKey="score" fill="#8B5CF6" radius={[0, 4, 4, 0]} />
                </BarChart>
              </ResponsiveContainer>
            )}
          </div>
        </div>

        <div className="lg:col-span-4 p-6 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-4 flex flex-col justify-between">
          <div>
            <h3 className="text-sm font-bold text-[#F4F4F6]">Difficulty Split</h3>
            <p className="text-xs text-[#8E8EA0] font-sans">Problems in your catalogue by difficulty</p>
          </div>
          <div className="h-44 w-full">
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie
                  data={pieData.length ? pieData : stats.difficultyDistribution}
                  innerRadius={35}
                  outerRadius={65}
                  paddingAngle={5}
                  dataKey="count"
                >
                  {(pieData.length ? pieData : stats.difficultyDistribution).map((entry, index) => (
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

      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <div className="lg:col-span-8 p-6 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-4">
          <h3 className="text-sm font-bold text-[#F4F4F6]">30-day activity</h3>
          <div className="h-56 w-full pt-2">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={daily.length ? daily : stats.monthlyProgress.map((m) => ({ date: m.month, solved: m.solved, attempts: 0 }))}>
                <defs>
                  <linearGradient id="colorSolved" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#8B5CF6" stopOpacity={0.3} />
                    <stop offset="95%" stopColor="#8B5CF6" stopOpacity={0.0} />
                  </linearGradient>
                </defs>
                <XAxis dataKey="date" stroke="#5A5A6E" fontSize={10} tickLine={false} hide={daily.length > 14} />
                <YAxis stroke="#5A5A6E" fontSize={11} tickLine={false} />
                <Tooltip
                  contentStyle={{
                    backgroundColor: '#0D0D14',
                    border: '1px solid #2B2B3C',
                    borderRadius: '0.75rem',
                    color: '#F4F4F6',
                    fontSize: '11px',
                  }}
                />
                <Area type="monotone" dataKey="solved" stroke="#8B5CF6" strokeWidth={2} fillOpacity={1} fill="url(#colorSolved)" />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </div>

        <div className="lg:col-span-4 p-6 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-4">
          <h3 className="text-sm font-bold text-[#F4F4F6]">Competency Assessment</h3>
          <div className="space-y-4 text-xs font-sans">
            <div className="p-3.5 bg-emerald-500/10 border border-emerald-500/20 rounded-xl space-y-1">
              <span className="text-emerald-400 font-bold flex items-center gap-1.5 font-mono">
                <CheckCircle2 className="w-3.5 h-3.5" /> High Mastery Index
              </span>
              <div className="text-sm font-semibold text-[#F4F4F6] font-mono">
                {stats.strongestTopic
                  ? `${stats.strongestTopic.topic} (${stats.strongestTopic.score}%)`
                  : 'Not enough data yet'}
              </div>
              <p className="text-[#8E8EA0] text-[11px]">
                Highest completion among topics you have started.
              </p>
            </div>
            <div className="p-3.5 bg-amber-500/10 border border-amber-500/20 rounded-xl space-y-1">
              <span className="text-amber-400 font-bold flex items-center gap-1.5 font-mono">
                <Flame className="w-3.5 h-3.5" /> Queued for Drill
              </span>
              <div className="text-sm font-semibold text-[#F4F4F6] font-mono">
                {stats.weakestTopic
                  ? `${stats.weakestTopic.topic} (${stats.weakestTopic.score}%)`
                  : 'Start more topics to rank weaknesses'}
              </div>
              <p className="text-[#8E8EA0] text-[11px]">
                Lowest completion — a good next practice target.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AnalyticsPage;
