import React from 'react';
import { AlgorithmMetadata, ExecutionStep } from '../../types/visualizer.types';
import {
  BookOpen,
  Zap,
  Target,
  ExternalLink,
  Bot,
  BrainCircuit,
} from 'lucide-react';
import { useNavigate } from 'react-router-dom';

interface ExplanationPanelProps {
  algorithm: AlgorithmMetadata;
  currentStepData: ExecutionStep | null;
  onOpenAI: () => void;
  onOpenQuiz: () => void;
}

export const ExplanationPanel: React.FC<ExplanationPanelProps> = ({
  algorithm,
  currentStepData,
  onOpenAI,
  onOpenQuiz,
}) => {
  const navigate = useNavigate();

  return (
    <div className="bg-white border border-[#E2E8F0] rounded-2xl p-5 flex flex-col gap-5 font-mono text-xs text-[#64748B] leading-relaxed">
      {/* Header Info */}
      <div className="flex flex-col gap-2">
        <div className="flex items-center justify-between gap-2">
          <span className="text-[10px] font-bold uppercase tracking-wider text-[#DC2626] bg-[#EF4444]/10 px-2 py-0.5 rounded border border-[#EF4444]/20">
            {algorithm.category}
          </span>
          <div className="flex items-center gap-1.5">
            <button
              onClick={onOpenQuiz}
              className="text-[11px] text-[#DC2626] hover:text-[#0F172A] flex items-center gap-1 bg-[#F1F5F9] hover:bg-[#E2E8F0] border border-[#E2E8F0] px-2.5 py-1 rounded-lg transition-colors"
            >
              <BrainCircuit className="w-3.5 h-3.5" /> Practice Quiz
            </button>
            <button
              onClick={onOpenAI}
              className="text-[11px] text-emerald-400 hover:text-[#0F172A] flex items-center gap-1 bg-emerald-500/10 hover:bg-emerald-500/20 border border-emerald-500/30 px-2.5 py-1 rounded-lg transition-colors"
            >
              <Bot className="w-3.5 h-3.5" /> Explain Step
            </button>
          </div>
        </div>
        <h3 className="text-lg font-bold text-[#0F172A] font-sans">
          {algorithm.name}
        </h3>
        <p className="text-xs text-[#CBD5E1] font-sans leading-normal">
          {algorithm.description}
        </p>
      </div>

      {/* Dynamic Step Explanation Box */}
      {currentStepData && (
        <div className="p-3.5 bg-[#F1F5F9] border border-[#EF4444]/30 rounded-xl space-y-1.5">
          <div className="flex items-center gap-2">
            <Zap className="w-3.5 h-3.5 text-[#DC2626]" />
            <span className="text-[11px] font-bold text-[#E9D5FF] uppercase tracking-wide">
              Step {currentStepData.step_number}: {currentStepData.operation.replace(/_/g, ' ')}
            </span>
          </div>
          <p className="text-xs text-[#E2E8F0] font-sans leading-relaxed">
            {currentStepData.explanation}
          </p>
        </div>
      )}

      {/* Complexity Breakdown Badges */}
      <div className="space-y-2">
        <span className="text-[10px] font-bold uppercase tracking-wider text-[#94A3B8] flex items-center gap-1">
          <Zap className="w-3 h-3 text-amber-400" /> Complexity Analysis
        </span>
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-2">
          <div className="p-2.5 bg-[#F1F5F9] rounded-xl border border-[#E2E8F0] flex flex-col gap-0.5">
            <span className="text-[9px] uppercase text-[#94A3B8] font-bold">Best Time</span>
            <span className="text-xs font-bold text-emerald-400">{algorithm.complexity.best}</span>
          </div>
          <div className="p-2.5 bg-[#F1F5F9] rounded-xl border border-[#E2E8F0] flex flex-col gap-0.5">
            <span className="text-[9px] uppercase text-[#94A3B8] font-bold">Average</span>
            <span className="text-xs font-bold text-amber-400">{algorithm.complexity.average}</span>
          </div>
          <div className="p-2.5 bg-[#F1F5F9] rounded-xl border border-[#E2E8F0] flex flex-col gap-0.5">
            <span className="text-[9px] uppercase text-[#94A3B8] font-bold">Worst Time</span>
            <span className="text-xs font-bold text-rose-400">{algorithm.complexity.worst}</span>
          </div>
          <div className="p-2.5 bg-[#F1F5F9] rounded-xl border border-[#E2E8F0] flex flex-col gap-0.5">
            <span className="text-[9px] uppercase text-[#94A3B8] font-bold">Space</span>
            <span className="text-xs font-bold text-[#DC2626]">{algorithm.complexity.space}</span>
          </div>
        </div>
      </div>

      {/* Problem Solved & Use Cases */}
      <div className="space-y-2">
        <span className="text-[10px] font-bold uppercase tracking-wider text-[#94A3B8] flex items-center gap-1">
          <Target className="w-3 h-3 text-[#EF4444]" /> Problem It Solves
        </span>
        <p className="text-xs text-[#CBD5E1] font-sans">
          {algorithm.problem_solved}
        </p>
      </div>

      {/* Use Cases */}
      {algorithm.use_cases && algorithm.use_cases.length > 0 && (
        <div className="space-y-2">
          <span className="text-[10px] font-bold uppercase tracking-wider text-[#94A3B8]">
            Common Use Cases
          </span>
          <ul className="space-y-1 text-[11px] text-[#DC2626] font-sans">
            {algorithm.use_cases.map((uc, idx) => (
              <li key={idx} className="flex items-start gap-1.5">
                <span className="text-[#EF4444]">•</span>
                <span>{uc}</span>
              </li>
            ))}
          </ul>
        </div>
      )}

      {/* Cross-Feature Links: Solve Problems in Practice Arena */}
      <div className="pt-2 border-t border-[#E2E8F0] flex items-center justify-between">
        <button
          onClick={() => navigate('/practice')}
          className="text-xs text-[#DC2626] hover:text-[#0F172A] flex items-center gap-1 font-sans font-medium"
        >
          <ExternalLink className="w-3.5 h-3.5" /> Practice related problems in Practice Arena
        </button>
        <button
          onClick={() => navigate('/code-review')}
          className="text-xs text-[#64748B] hover:text-[#0F172A] flex items-center gap-1 font-sans"
        >
          Review Code with AI →
        </button>
      </div>
    </div>
  );
};
export default ExplanationPanel;

