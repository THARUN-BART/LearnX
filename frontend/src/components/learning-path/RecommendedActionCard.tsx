import React from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Sparkles,
  ArrowRight,
  AlertCircle,
  BookOpen,
  Target,
  RefreshCw,
  Code2,
  CheckCircle,
} from 'lucide-react';
import { LearningRecommendation, ActionType } from '../../types/learning-path.types';

interface RecommendedActionCardProps {
  recommendation: LearningRecommendation | null;
  loading?: boolean;
  onRefresh?: () => void;
  compact?: boolean;
}

const actionMeta: Record<
  ActionType,
  { label: string; icon: React.FC<{ className?: string }>; bg: string; text: string; border: string }
> = {
  learn_concept: {
    label: 'Learn New Concept',
    icon: BookOpen,
    bg: 'bg-emerald-500/10',
    text: 'text-emerald-400',
    border: 'border-emerald-500/30',
  },
  revise_concept: {
    label: 'Revise Concept',
    icon: RefreshCw,
    bg: 'bg-amber-500/10',
    text: 'text-amber-400',
    border: 'border-amber-500/30',
  },
  solve_easy: {
    label: 'Warm Up Practice',
    icon: Code2,
    bg: 'bg-blue-500/10',
    text: 'text-blue-400',
    border: 'border-blue-500/30',
  },
  practice_medium: {
    label: 'Core Practice',
    icon: Target,
    bg: 'bg-purple-500/10',
    text: 'text-purple-400',
    border: 'border-purple-500/30',
  },
  attempt_hard: {
    label: 'Challenge Problem',
    icon: Sparkles,
    bg: 'bg-rose-500/10',
    text: 'text-rose-400',
    border: 'border-rose-500/30',
  },
  review_mistake: {
    label: 'Review Weakness',
    icon: AlertCircle,
    bg: 'bg-amber-500/10',
    text: 'text-amber-400',
    border: 'border-amber-500/30',
  },
  complete_prerequisite: {
    label: 'Prerequisite Required',
    icon: Target,
    bg: 'bg-indigo-500/10',
    text: 'text-indigo-400',
    border: 'border-indigo-500/30',
  },
};

const diffMeta = {
  Easy: { bg: 'bg-emerald-500/10', text: 'text-emerald-400', border: 'border-emerald-500/20' },
  Medium: { bg: 'bg-amber-500/10', text: 'text-amber-400', border: 'border-amber-500/20' },
  Hard: { bg: 'bg-rose-500/10', text: 'text-rose-400', border: 'border-rose-500/20' },
};

export const RecommendedActionCard: React.FC<RecommendedActionCardProps> = ({
  recommendation,
  loading = false,
  onRefresh,
  compact = false,
}) => {
  const navigate = useNavigate();

  if (loading) {
    return (
      <div className="rounded-2xl border border-[#1C1C28] bg-gradient-to-br from-[#0D0D14] to-[#12121B] p-6 animate-pulse">
        <div className="h-6 w-48 bg-[#1C1C28] rounded mb-4" />
        <div className="h-8 w-3/4 bg-[#1C1C28] rounded mb-3" />
        <div className="h-4 w-full bg-[#1C1C28] rounded mb-6" />
        <div className="h-10 w-36 bg-[#1C1C28] rounded" />
      </div>
    );
  }

  if (!recommendation) {
    return (
      <div className="rounded-2xl border border-[#1C1C28] bg-[#0D0D14] p-6 flex flex-col items-center text-center">
        <CheckCircle className="w-10 h-10 text-emerald-400 mb-2" />
        <h3 className="text-base font-semibold text-[#F4F4F6]">All Set!</h3>
        <p className="text-sm text-[#8E8EA0] max-w-sm mt-1">
          Your adaptive learning path is up to date. Solve any problem or recalculate your recommendations.
        </p>
        {onRefresh && (
          <button
            onClick={onRefresh}
            className="mt-4 px-4 py-2 rounded-xl bg-[#1C1C28] hover:bg-[#252538] text-xs font-medium text-[#F4F4F6] transition-colors"
          >
            Recalculate Path
          </button>
        )}
      </div>
    );
  }

  const meta = actionMeta[recommendation.action_type] || actionMeta.practice_medium;
  const ActionIcon = meta.icon;
  const diff = diffMeta[recommendation.recommended_difficulty] || diffMeta.Easy;

  const handleAction = () => {
    if (recommendation.problem_id) {
      navigate(`/practice/${recommendation.problem_id}`);
    } else if (recommendation.topic_id) {
      navigate(`/roadmap/${recommendation.topic_id}`);
    } else {
      navigate('/practice');
    }
  };

  return (
    <div className="relative overflow-hidden rounded-2xl border border-[#8B5CF6]/30 bg-gradient-to-br from-[#0D0D14] via-[#12121E] to-[#161226] p-6 shadow-xl shadow-[#8B5CF6]/5 group">
      {/* Glow accent */}
      <div className="absolute -top-24 -right-24 w-60 h-60 bg-[#8B5CF6]/15 rounded-full blur-3xl pointer-events-none group-hover:bg-[#8B5CF6]/20 transition-all duration-500" />

      {/* Top Header */}
      <div className="flex flex-wrap items-center justify-between gap-3 relative z-10 mb-4">
        <div className="flex items-center gap-2">
          <span
            className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold uppercase tracking-wider border ${meta.bg} ${meta.text} ${meta.border}`}
          >
            <ActionIcon className="w-3.5 h-3.5" />
            {meta.label}
          </span>

          <span
            className={`px-2.5 py-0.5 rounded-full text-xs font-semibold border ${diff.bg} ${diff.text} ${diff.border}`}
          >
            {recommendation.recommended_difficulty}
          </span>
        </div>

        {onRefresh && (
          <button
            onClick={onRefresh}
            title="Re-evaluate learning path"
            className="flex items-center gap-1.5 text-xs text-[#8E8EA0] hover:text-[#F4F4F6] transition-colors bg-[#1C1C28]/60 hover:bg-[#1C1C28] px-2.5 py-1 rounded-lg border border-[#252538]"
          >
            <RefreshCw className="w-3 h-3" />
            <span>Recalculate</span>
          </button>
        )}
      </div>

      {/* Headline & Context */}
      <div className="relative z-10 mb-3">
        <h2 className="text-xl md:text-2xl font-bold text-[#F4F4F6] tracking-tight group-hover:text-white transition-colors">
          {recommendation.headline}
        </h2>
        <div className="flex items-center gap-2 text-xs text-[#A0A0B2] mt-1 font-mono">
          <span>Topic: {recommendation.topic_name || recommendation.concept_name}</span>
          {recommendation.problem_title && (
            <>
              <span>•</span>
              <span className="text-[#8B5CF6]">Problem: {recommendation.problem_title}</span>
            </>
          )}
        </div>
      </div>

      {/* Pedagogical Explanation */}
      <div className="relative z-10 bg-[#09090F]/70 rounded-xl p-3.5 border border-[#1C1C28] mb-4 text-sm text-[#C4C4D4] leading-relaxed">
        <p className="text-xs text-[#8E8EA0] uppercase tracking-wider font-semibold font-mono mb-1">
          Why this recommendation?
        </p>
        <p>{recommendation.explanation}</p>
      </div>

      {/* Weakness Alert if applicable */}
      {recommendation.relevant_weakness && (
        <div className="relative z-10 flex items-center gap-2.5 p-3 rounded-xl bg-amber-500/10 border border-amber-500/20 text-xs text-amber-300 mb-5">
          <AlertCircle className="w-4 h-4 shrink-0 text-amber-400" />
          <span>
            <strong>Identified Weakness:</strong> {recommendation.relevant_weakness}
          </span>
        </div>
      )}

      {/* CTA Button */}
      <div className="relative z-10 flex items-center justify-between pt-2">
        <button
          onClick={handleAction}
          className="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-gradient-to-r from-[#8B5CF6] to-[#7C3AED] hover:from-[#7C3AED] hover:to-[#6D28D9] text-white font-medium text-sm shadow-lg shadow-[#8B5CF6]/25 hover:shadow-[#8B5CF6]/40 transition-all active:scale-[0.98]"
        >
          <span>
            {recommendation.problem_id
              ? `Solve ${recommendation.problem_title || 'Problem'}`
              : `Start ${recommendation.concept_name}`}
          </span>
          <ArrowRight className="w-4 h-4" />
        </button>

        {!compact && (
          <span className="text-xs text-[#5A5A6E] font-mono hidden sm:inline">
            Evaluated deterministically from learning memory
          </span>
        )}
      </div>
    </div>
  );
};
