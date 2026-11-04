import React from 'react';
import { useNavigate } from 'react-router-dom';
import {
  History,
  CheckCircle2,
  Clock,
  ArrowRight,
  Sparkles,
  AlertCircle,
  RefreshCw,
  BookOpen,
} from 'lucide-react';
import { LearningRecommendation, ActionType } from '../../types/learning-path.types';

interface RecommendationHistoryListProps {
  history: LearningRecommendation[];
  loading?: boolean;
}

const actionBadges: Record<ActionType, { label: string; text: string; bg: string }> = {
  learn_concept: { label: 'New Concept', text: 'text-emerald-400', bg: 'bg-emerald-500/10' },
  revise_concept: { label: 'Revision', text: 'text-amber-400', bg: 'bg-amber-500/10' },
  solve_easy: { label: 'Easy Problem', text: 'text-blue-400', bg: 'bg-blue-500/10' },
  practice_medium: { label: 'Medium Problem', text: 'text-purple-400', bg: 'bg-purple-500/10' },
  attempt_hard: { label: 'Hard Problem', text: 'text-rose-400', bg: 'bg-rose-500/10' },
  review_mistake: { label: 'Mistake Review', text: 'text-amber-400', bg: 'bg-amber-500/10' },
  complete_prerequisite: { label: 'Prerequisite', text: 'text-indigo-400', bg: 'bg-indigo-500/10' },
};

export const RecommendationHistoryList: React.FC<RecommendationHistoryListProps> = ({
  history,
  loading = false,
}) => {
  const navigate = useNavigate();

  if (loading) {
    return (
      <div className="rounded-2xl border border-[#E2E8F0] bg-white p-5 space-y-3 animate-pulse">
        {[...Array(3)].map((_, i) => (
          <div key={i} className="h-16 bg-[#161622] rounded-xl" />
        ))}
      </div>
    );
  }

  if (history.length === 0) {
    return (
      <div className="rounded-2xl border border-[#E2E8F0] bg-white p-8 text-center text-xs text-[#64748B] font-mono">
        <History className="w-8 h-8 text-[#94A3B8] mx-auto mb-2" />
        No previous recommendation history yet. Complete problems to see adaptive transitions over time.
      </div>
    );
  }

  return (
    <div className="rounded-2xl border border-[#E2E8F0] bg-white p-5">
      <div className="flex items-center justify-between mb-4 pb-3 border-b border-[#E2E8F0]">
        <h3 className="text-sm font-semibold text-[#0F172A] flex items-center gap-2">
          <History className="w-4 h-4 text-[#EF4444]" />
          <span>Adaptive Recommendation Timeline</span>
        </h3>
        <span className="text-xs font-mono text-[#64748B]">{history.length} records</span>
      </div>

      <div className="space-y-3">
        {history.map((item) => {
          const badge = actionBadges[item.action_type] || actionBadges.practice_medium;
          const date = new Date(item.created_at).toLocaleDateString(undefined, {
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
          });

          return (
            <div
              key={item.id}
              className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 p-3.5 rounded-xl bg-[#141420] border border-[#E2E8F0] hover:border-[#252538] transition-colors"
            >
              <div className="space-y-1">
                <div className="flex items-center gap-2">
                  <span
                    className={`px-2 py-0.5 rounded text-[10px] font-semibold font-mono uppercase ${badge.bg} ${badge.text}`}
                  >
                    {badge.label}
                  </span>
                  <span className="text-xs font-semibold text-[#0F172A]">
                    {item.headline || item.concept_name}
                  </span>
                  {item.is_active && (
                    <span className="px-1.5 py-0.2 rounded bg-emerald-500/10 text-emerald-400 text-[10px] font-mono">
                      Current
                    </span>
                  )}
                </div>
                <p className="text-xs text-[#64748B] line-clamp-1">{item.explanation}</p>
                <div className="flex items-center gap-3 text-[10px] font-mono text-[#94A3B8]">
                  <span>{date}</span>
                  {item.problem_title && <span>Problem: {item.problem_title}</span>}
                  <span>Diff: {item.recommended_difficulty}</span>
                </div>
              </div>

              {item.problem_id && (
                <button
                  onClick={() => navigate(`/practice/${item.problem_id}`)}
                  className="shrink-0 inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-[#F1F5F9] hover:bg-[#252538] text-xs font-mono text-[#C4C4D4] hover:text-[#0F172A] transition-colors"
                >
                  <span>Solve</span>
                  <ArrowRight className="w-3 h-3" />
                </button>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
};
