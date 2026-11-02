import React from 'react';
import {
  CheckCircle2,
  Clock,
  AlertTriangle,
  PlayCircle,
  TrendingUp,
  Layers,
  Sparkles,
} from 'lucide-react';
import { PathProgressResponse } from '../../types/learning-path.types';

interface PathProgressStatsProps {
  progress: PathProgressResponse | null;
  loading?: boolean;
}

export const PathProgressStats: React.FC<PathProgressStatsProps> = ({ progress, loading = false }) => {
  if (loading || !progress) {
    return (
      <div className="grid grid-cols-2 md:grid-cols-4 gap-3.5 animate-pulse">
        {[...Array(4)].map((_, i) => (
          <div key={i} className="h-24 rounded-2xl bg-[#0D0D14] border border-[#1C1C28]" />
        ))}
      </div>
    );
  }

  const stats = [
    {
      label: 'Concepts Mastered',
      value: progress.mastered_count,
      total: progress.total_concepts,
      icon: CheckCircle2,
      accent: 'text-emerald-400',
      bg: 'bg-emerald-500/10',
      border: 'border-emerald-500/20',
    },
    {
      label: 'In Progress / Active',
      value: progress.in_progress_count,
      total: progress.total_concepts,
      icon: Clock,
      accent: 'text-blue-400',
      bg: 'bg-blue-500/10',
      border: 'border-blue-500/20',
    },
    {
      label: 'Needs Revision',
      value: progress.needs_revision_count,
      total: progress.total_concepts,
      icon: AlertTriangle,
      accent: 'text-amber-400',
      bg: 'bg-amber-500/10',
      border: 'border-amber-500/20',
    },
    {
      label: 'Ready to Learn',
      value: progress.ready_to_learn_count,
      total: progress.total_concepts,
      icon: PlayCircle,
      accent: 'text-cyan-400',
      bg: 'bg-cyan-500/10',
      border: 'border-cyan-500/20',
    },
  ];

  return (
    <div className="space-y-4">
      {/* Top Level Progress Meter */}
      <div className="rounded-2xl border border-[#1C1C28] bg-[#0D0D14] p-5">
        <div className="flex flex-wrap items-center justify-between gap-3 mb-3">
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 rounded-xl bg-[#8B5CF6]/15 flex items-center justify-center text-[#8B5CF6]">
              <Layers className="w-4 h-4" />
            </div>
            <div>
              <h3 className="text-sm font-semibold text-[#F4F4F6]">Overall Curriculum Mastery</h3>
              <p className="text-xs text-[#8E8EA0] font-mono">Stage: {progress.current_stage}</p>
            </div>
          </div>
          <div className="text-right">
            <span className="text-2xl font-bold font-mono text-[#8B5CF6]">
              {Math.round(progress.overall_progress_percentage)}%
            </span>
          </div>
        </div>

        <div className="w-full h-2.5 rounded-full bg-[#1C1C28] overflow-hidden">
          <div
            className="h-full rounded-full bg-gradient-to-r from-[#8B5CF6] to-[#06B6D4] transition-all duration-700"
            style={{ width: `${Math.min(100, Math.max(2, progress.overall_progress_percentage))}%` }}
          />
        </div>
      </div>

      {/* Breakdown Cards */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-3.5">
        {stats.map((s) => {
          const Icon = s.icon;
          return (
            <div
              key={s.label}
              className={`rounded-2xl border ${s.border} ${s.bg} p-4 flex flex-col justify-between`}
            >
              <div className="flex items-center justify-between mb-2">
                <span className="text-xs font-mono text-[#8E8EA0]">{s.label}</span>
                <Icon className={`w-4 h-4 ${s.accent}`} />
              </div>
              <div className="flex items-baseline gap-1.5">
                <span className={`text-2xl font-bold font-mono ${s.accent}`}>{s.value}</span>
                <span className="text-xs font-mono text-[#5A5A6E]">/ {s.total}</span>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};
