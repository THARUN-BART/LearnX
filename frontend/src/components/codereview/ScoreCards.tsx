import React from 'react';
import { ReviewScores } from '../../types/codereview.types';
import { CheckCircle2, Award, Zap, ShieldCheck } from 'lucide-react';

interface ScoreCardsProps {
  scores: ReviewScores;
}

export const ScoreCards: React.FC<ScoreCardsProps> = ({ scores }) => {
  const getScoreColor = (score: number) => {
    if (score >= 85) return 'text-emerald-400 border-emerald-500/30 bg-emerald-500/10';
    if (score >= 65) return 'text-amber-400 border-amber-500/30 bg-amber-500/10';
    return 'text-rose-400 border-rose-500/30 bg-rose-500/10';
  };

  return (
    <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 font-mono">
      {/* Overall Score */}
      <div className={`p-4 rounded-2xl border flex flex-col justify-between ${getScoreColor(scores.overall)}`}>
        <div className="flex items-center justify-between">
          <span className="text-[10px] uppercase font-bold tracking-wider">Overall Score</span>
          <Award className="w-4 h-4" />
        </div>
        <div className="mt-2 flex items-baseline gap-1">
          <span className="text-3xl font-bold font-sans">{scores.overall}</span>
          <span className="text-xs opacity-60">/ 100</span>
        </div>
      </div>

      {/* Readability */}
      <div className="p-4 rounded-2xl border border-[#E2E8F0] bg-[#F1F5F9] flex flex-col justify-between">
        <div className="flex items-center justify-between text-[#64748B]">
          <span className="text-[10px] uppercase font-bold tracking-wider">Readability</span>
          <CheckCircle2 className="w-4 h-4 text-blue-400" />
        </div>
        <div className="mt-2 flex items-baseline gap-1">
          <span className="text-2xl font-bold font-sans text-blue-400">{scores.readability}</span>
          <span className="text-xs text-[#94A3B8]">/ 100</span>
        </div>
      </div>

      {/* Efficiency */}
      <div className="p-4 rounded-2xl border border-[#E2E8F0] bg-[#F1F5F9] flex flex-col justify-between">
        <div className="flex items-center justify-between text-[#64748B]">
          <span className="text-[10px] uppercase font-bold tracking-wider">Efficiency</span>
          <Zap className="w-4 h-4 text-amber-400" />
        </div>
        <div className="mt-2 flex items-baseline gap-1">
          <span className="text-2xl font-bold font-sans text-amber-400">{scores.efficiency}</span>
          <span className="text-xs text-[#94A3B8]">/ 100</span>
        </div>
      </div>

      {/* Maintainability */}
      <div className="p-4 rounded-2xl border border-[#E2E8F0] bg-[#F1F5F9] flex flex-col justify-between">
        <div className="flex items-center justify-between text-[#64748B]">
          <span className="text-[10px] uppercase font-bold tracking-wider">Maintainability</span>
          <ShieldCheck className="w-4 h-4 text-[#DC2626]" />
        </div>
        <div className="mt-2 flex items-baseline gap-1">
          <span className="text-2xl font-bold font-sans text-[#DC2626]">{scores.maintainability}</span>
          <span className="text-xs text-[#94A3B8]">/ 100</span>
        </div>
      </div>
    </div>
  );
};
export default ScoreCards;

