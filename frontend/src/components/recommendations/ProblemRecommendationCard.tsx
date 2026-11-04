import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Play,
  ExternalLink,
  Target,
  BrainCircuit,
  CheckCircle2,
  ThumbsUp,
  Zap,
  Flame,
  ChevronRight,
  Sparkles
} from 'lucide-react';
import { ProblemRecommendation, FeedbackRating } from '../../types/recommendations.types';
import Badge from '../ui/Badge';
import Button from '../ui/Button';
import { toast } from 'sonner';

interface ProblemRecommendationCardProps {
  recommendation: ProblemRecommendation;
  onFeedbackSubmit?: (recId: string, rating: FeedbackRating) => void;
}

export const ProblemRecommendationCard: React.FC<ProblemRecommendationCardProps> = ({
  recommendation,
  onFeedbackSubmit
}) => {
  const navigate = useNavigate();
  const [rated, setRated] = useState<string | null>(recommendation.feedback_rating || null);

  const handleStartPractice = () => {
    navigate(`/practice/${recommendation.problem_id || recommendation.problem_slug}`);
  };

  const handleRating = (rating: FeedbackRating, label: string) => {
    setRated(rating);
    if (onFeedbackSubmit) {
      onFeedbackSubmit(recommendation.id, rating);
    }
    toast.success(`Feedback recorded: ${label}`);
  };

  const getCategoryMeta = (cat: string) => {
    switch (cat) {
      case 'improve_weaknesses':
        return {
          label: 'Weakness Target',
          color: 'bg-rose-500/10 text-rose-400 border-rose-500/20'
        };
      case 'strengthen_fundamentals':
        return {
          label: 'Foundations',
          color: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20'
        };
      case 'challenge_yourself':
        return {
          label: 'High Challenge',
          color: 'bg-amber-500/10 text-amber-400 border-amber-500/20'
        };
      case 'revision_practice':
        return {
          label: 'Revision Queue',
          color: 'bg-sky-500/10 text-sky-400 border-sky-500/20'
        };
      case 'similar_problems':
        return {
          label: 'Pattern Variant',
          color: 'bg-purple-500/10 text-purple-400 border-purple-500/20'
        };
      default:
        return {
          label: 'Recommended',
          color: 'bg-[#EF4444]/10 text-[#DC2626] border-[#EF4444]/20'
        };
    }
  };

  const meta = getCategoryMeta(recommendation.category);

  return (
    <div className="p-5 sm:p-6 rounded-2xl bg-white hover:bg-[#F1F5F9] border border-[#E2E8F0] hover:border-[#EF4444]/40 transition-all duration-200 flex flex-col justify-between group relative overflow-hidden">
      <div className="space-y-4">
        {/* Header tags */}
        <div className="flex items-center justify-between gap-2 flex-wrap">
          <div className="flex items-center gap-2">
            <span
              className={`text-[10px] font-mono font-semibold px-2.5 py-0.5 rounded-full border ${meta.color}`}
            >
              {meta.label}
            </span>

            <Badge
              variant={
                recommendation.recommended_difficulty === 'Easy'
                  ? 'easy'
                  : recommendation.recommended_difficulty === 'Medium'
                  ? 'medium'
                  : 'hard'
              }
              size="sm"
            >
              {recommendation.recommended_difficulty}
            </Badge>
          </div>

          <span className="text-[11px] font-mono text-[#94A3B8] flex items-center gap-1">
            <Target className="w-3 h-3 text-[#EF4444]" />
            {Math.round(recommendation.match_score)}% match
          </span>
        </div>

        {/* Title and Headline */}
        <div>
          <h3
            onClick={handleStartPractice}
            className="text-base sm:text-lg font-bold text-[#0F172A] group-hover:text-[#DC2626] font-mono cursor-pointer transition-colors"
          >
            {recommendation.problem_title || 'DSA Problem'}
          </h3>
          <p className="text-xs font-semibold text-[#64748B] mt-1 font-mono line-clamp-1">
            {recommendation.headline}
          </p>
        </div>

        {/* Pedagogical Explanation */}
        <div className="p-3.5 rounded-xl bg-[#F8FAFC] border border-[#E2E8F0] space-y-2">
          <div className="flex items-center gap-1.5 text-[10px] font-mono text-[#94A3B8] uppercase tracking-wider font-semibold">
            <BrainCircuit className="w-3 h-3 text-[#EF4444]" />
            <span>Why this problem</span>
          </div>
          <p className="text-xs text-[#C4C4D0] leading-relaxed line-clamp-3 font-sans">
            {recommendation.explanation}
          </p>

          {recommendation.learning_objective && (
            <div className="pt-2 border-t border-[#E2E8F0] flex items-start gap-1.5 text-[11px] text-[#64748B]">
              <CheckCircle2 className="w-3 h-3 text-emerald-400 shrink-0 mt-0.5" />
              <span className="line-clamp-2">{recommendation.learning_objective}</span>
            </div>
          )}
        </div>

        {/* Concept pills */}
        <div className="flex flex-wrap gap-1.5 pt-1">
          <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-[#F1F5F9] text-[#DC2626] border border-[#EF4444]/20">
            {recommendation.primary_concept}
          </span>
          {(recommendation.related_concepts || []).slice(0, 2).map((c, i) => (
            <span
              key={i}
              className="text-[10px] font-mono px-2 py-0.5 rounded bg-[#F1F5F9] text-[#64748B] border border-[#E2E8F0]"
            >
              {c}
            </span>
          ))}
        </div>
      </div>

      {/* Footer controls & Feedback */}
      <div className="pt-4 mt-4 border-t border-[#E2E8F0] flex items-center justify-between gap-3">
        {/* Quick Feedback Pills */}
        <div className="flex items-center gap-1 text-xs">
          <button
            title="Helpful suggestion"
            onClick={() => handleRating('helpful', 'Helpful')}
            className={`p-1.5 rounded-lg border transition-all text-xs ${
              rated === 'helpful'
                ? 'bg-emerald-500/20 text-emerald-400 border-emerald-500/40'
                : 'bg-[#F1F5F9] text-[#94A3B8] hover:text-[#0F172A] border-[#E2E8F0]'
            }`}
          >
            <ThumbsUp className="w-3 h-3" />
          </button>
          <button
            title="Too easy"
            onClick={() => handleRating('too_easy', 'Marked as Too Easy')}
            className={`p-1.5 rounded-lg border transition-all text-xs ${
              rated === 'too_easy'
                ? 'bg-amber-500/20 text-amber-400 border-amber-500/40'
                : 'bg-[#F1F5F9] text-[#94A3B8] hover:text-[#0F172A] border-[#E2E8F0]'
            }`}
          >
            <Zap className="w-3 h-3" />
          </button>
          <button
            title="Too difficult"
            onClick={() => handleRating('too_difficult', 'Marked as Too Difficult')}
            className={`p-1.5 rounded-lg border transition-all text-xs ${
              rated === 'too_difficult'
                ? 'bg-rose-500/20 text-rose-400 border-rose-500/40'
                : 'bg-[#F1F5F9] text-[#94A3B8] hover:text-[#0F172A] border-[#E2E8F0]'
            }`}
          >
            <Flame className="w-3 h-3" />
          </button>
        </div>

        {/* CTA Solve */}
        <Button
          variant="secondary"
          size="sm"
          className="font-mono text-xs"
          rightIcon={<Play className="w-3 h-3" />}
          onClick={handleStartPractice}
        >
          Solve
        </Button>
      </div>
    </div>
  );
};

export default ProblemRecommendationCard;

