import React from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Sparkles,
  Play,
  ExternalLink,
  Target,
  CheckCircle2,
  BrainCircuit,
  TrendingUp,
  Flame
} from 'lucide-react';
import { ProblemRecommendation } from '../../types/recommendations.types';
import Badge from '../ui/Badge';
import Button from '../ui/Button';

interface TopPickHeroProps {
  recommendation: ProblemRecommendation;
  onFeedbackClick?: (rec: ProblemRecommendation) => void;
}

export const TopPickHero: React.FC<TopPickHeroProps> = ({
  recommendation,
  onFeedbackClick
}) => {
  const navigate = useNavigate();

  const handleStartPractice = () => {
    navigate(`/practice/${recommendation.problem_id || recommendation.problem_slug}`);
  };

  return (
    <div className="relative overflow-hidden rounded-2xl bg-gradient-to-r from-[#12121D] via-[#151522] to-[#0D0D14] border border-[#EF4444]/40 p-6 sm:p-8 shadow-xl shadow-[#8B5CF6]/5 glow-purple">
      {/* Background glow & subtle tech grid */}
      <div className="absolute top-0 right-0 w-96 h-96 bg-[#EF4444]/10 rounded-full blur-3xl pointer-events-none -mr-20 -mt-20" />
      <div className="absolute inset-0 bg-tech-grid opacity-20 pointer-events-none" />

      <div className="relative z-10 flex flex-col lg:flex-row lg:items-center justify-between gap-6">
        <div className="space-y-4 max-w-3xl">
          {/* Top badges bar */}
          <div className="flex flex-wrap items-center gap-2.5 font-mono text-xs">
            <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-[#EF4444]/20 text-[#DC2626] border border-[#EF4444]/40 font-bold uppercase tracking-wider text-[10px]">
              <Sparkles className="w-3.5 h-3.5 text-[#EF4444]" />
              Top Priority Pick
            </span>

            <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 text-[11px]">
              <Target className="w-3 h-3" />
              {Math.round(recommendation.match_score)}% Concept Alignment
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

            <span className="text-[#94A3B8] hidden sm:inline">//</span>
            <span className="text-xs text-[#64748B]">{recommendation.primary_concept}</span>
          </div>

          {/* Title & Headline */}
          <div>
            <h2 className="text-2xl sm:text-3xl font-bold tracking-tight text-[#0F172A] font-mono group cursor-pointer hover:text-[#DC2626] transition-colors"
                onClick={handleStartPractice}>
              {recommendation.problem_title || 'Target DSA Problem'}
            </h2>
            <p className="text-sm font-semibold text-[#DC2626] mt-1 font-mono">
              {recommendation.headline}
            </p>
          </div>

          {/* Explanation rationale */}
          <div className="p-4 rounded-xl bg-[#09090F]/70 border border-[#E2E8F0] space-y-2">
            <div className="flex items-center gap-2 text-xs font-mono text-[#64748B]">
              <BrainCircuit className="w-3.5 h-3.5 text-[#EF4444]" />
              <span className="uppercase tracking-wider font-semibold">Pedagogical Rationale</span>
            </div>
            <p className="text-xs sm:text-sm text-[#D4D4D8] leading-relaxed font-sans">
              {recommendation.explanation}
            </p>

            {recommendation.learning_objective && (
              <div className="pt-2 border-t border-[#E2E8F0]/80 flex items-start gap-2 text-xs text-[#64748B]">
                <CheckCircle2 className="w-3.5 h-3.5 text-emerald-400 shrink-0 mt-0.5" />
                <span>
                  <strong className="text-white font-mono">Learning Goal:</strong>{' '}
                  {recommendation.learning_objective}
                </span>
              </div>
            )}
          </div>
        </div>

        {/* Action Controls */}
        <div className="flex flex-col sm:flex-row lg:flex-col gap-3 shrink-0 lg:w-56">
          <Button
            variant="primary"
            size="lg"
            className="w-full font-mono text-sm shadow-lg shadow-[#8B5CF6]/20"
            rightIcon={<Play className="w-4 h-4" />}
            onClick={handleStartPractice}
          >
            Solve Problem
          </Button>

          {recommendation.external_problem_url && (
            <Button
              variant="outline"
              size="sm"
              className="w-full font-mono text-xs"
              rightIcon={<ExternalLink className="w-3.5 h-3.5" />}
              onClick={() => window.open(recommendation.external_problem_url, '_blank')}
            >
              TakeUForward Sheet
            </Button>
          )}

          {onFeedbackClick && (
            <button
              onClick={() => onFeedbackClick(recommendation)}
              className="text-[11px] font-mono text-[#64748B] hover:text-[#0F172A] text-center pt-1 transition-colors underline underline-offset-4 decoration-[#2D2D3E]"
            >
              Rate this suggestion
            </button>
          )}
        </div>
      </div>
    </div>
  );
};

export default TopPickHero;

