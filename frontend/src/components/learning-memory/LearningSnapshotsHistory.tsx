import React from 'react';
import {
  History,
  TrendingUp,
  Brain,
  Sparkles,
  CheckCircle2,
  AlertTriangle,
  Calendar,
} from 'lucide-react';
import {
  ResponsiveContainer,
  AreaChart,
  Area,
  XAxis,
  YAxis,
  Tooltip,
} from 'recharts';
import { LearningSnapshot } from '../../types/learning-memory.types';

interface LearningSnapshotsHistoryProps {
  snapshots: LearningSnapshot[];
  loading?: boolean;
}

export const LearningSnapshotsHistory: React.FC<LearningSnapshotsHistoryProps> = ({
  snapshots,
  loading = false,
}) => {
  if (loading) {
    return (
      <div className="h-64 rounded-2xl bg-[#12121B] border border-[#1C1C28] animate-pulse" />
    );
  }

  const chartData = snapshots.map((s) => ({
    date: s.snapshot_date ? s.snapshot_date.slice(5) : 'Day',
    independent: s.independent_solves,
    assisted: s.assisted_solves,
    mistakes: s.mistake_count,
    corrected: s.corrected_mistake_count,
  }));

  const totalIndependent = snapshots.reduce((acc, s) => acc + s.independent_solves, 0);
  const totalAssisted = snapshots.reduce((acc, s) => acc + s.assisted_solves, 0);
  const totalCorrected = snapshots.reduce((acc, s) => acc + s.corrected_mistake_count, 0);

  return (
    <div className="rounded-2xl bg-[#0D0D14] border border-[#1C1C28] p-6 space-y-6 font-mono">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2 text-xs font-semibold uppercase tracking-wider text-[#A78BFA]">
            <History className="w-4 h-4 text-[#8B5CF6]" />
            <span>Mastery Milestones & Snapshots</span>
          </div>
          <p className="text-xs text-[#8E8EA0] mt-1 font-sans">
            Longitudinal telemetry tracking independent vs AI-assisted problem solves and mistake resolution.
          </p>
        </div>

        <div className="flex items-center gap-2 text-xs">
          <span className="px-2.5 py-1 rounded-lg bg-emerald-500/10 border border-emerald-500/20 text-emerald-400">
            {totalIndependent} Independent Solves
          </span>
          <span className="px-2.5 py-1 rounded-lg bg-[#8B5CF6]/10 border border-[#8B5CF6]/20 text-[#A78BFA]">
            {totalAssisted} AI-Assisted
          </span>
          <span className="px-2.5 py-1 rounded-lg bg-amber-500/10 border border-amber-500/20 text-amber-400">
            {totalCorrected} Corrected
          </span>
        </div>
      </div>

      {chartData.length > 0 ? (
        <div className="h-56 w-full pt-2">
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={chartData}>
              <defs>
                <linearGradient id="colorIndep" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#10B981" stopOpacity={0.4} />
                  <stop offset="95%" stopColor="#10B981" stopOpacity={0.0} />
                </linearGradient>
                <linearGradient id="colorAssist" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#8B5CF6" stopOpacity={0.4} />
                  <stop offset="95%" stopColor="#8B5CF6" stopOpacity={0.0} />
                </linearGradient>
              </defs>
              <XAxis dataKey="date" stroke="#5A5A6E" fontSize={11} tickLine={false} axisLine={{ stroke: '#1C1C28' }} />
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
              />
              <Area
                type="monotone"
                dataKey="independent"
                name="Independent Solves"
                stroke="#10B981"
                fillOpacity={1}
                fill="url(#colorIndep)"
              />
              <Area
                type="monotone"
                dataKey="assisted"
                name="AI-Assisted Solves"
                stroke="#8B5CF6"
                fillOpacity={1}
                fill="url(#colorAssist)"
              />
            </AreaChart>
          </ResponsiveContainer>
        </div>
      ) : (
        <div className="p-8 text-center rounded-xl bg-[#12121B] border border-dashed border-[#2A2A3C] text-xs text-[#8E8EA0]">
          <Brain className="w-8 h-8 text-[#5A5A6E] mx-auto mb-2" />
          <p>No snapshots logged yet. Solving problems logs your mastery milestones automatically.</p>
        </div>
      )}

      {/* Snapshot Log Cards */}
      {snapshots.length > 0 && (
        <div className="space-y-2 pt-2">
          <div className="text-[11px] uppercase tracking-wider text-[#5A5A6E] font-semibold">
            Historical Snapshots ({snapshots.length})
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
            {snapshots.slice(0, 6).map((snap) => (
              <div
                key={snap.id || snap.snapshot_date}
                className="p-3.5 rounded-xl bg-[#12121B] border border-[#1C1C28] flex items-center justify-between text-xs"
              >
                <div className="space-y-1">
                  <div className="flex items-center gap-1.5 text-[#F4F4F6] font-semibold">
                    <Calendar className="w-3.5 h-3.5 text-[#8B5CF6]" />
                    <span>{snap.snapshot_date}</span>
                  </div>
                  <div className="text-[11px] text-[#8E8EA0]">
                    {snap.total_concepts_practised} concepts touched
                  </div>
                </div>
                <div className="text-right space-y-0.5">
                  <div className="text-emerald-400 font-bold">
                    +{snap.independent_solves} solo
                  </div>
                  <div className="text-[10px] text-amber-400">
                    {snap.corrected_mistake_count} fixes
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default LearningSnapshotsHistory;

