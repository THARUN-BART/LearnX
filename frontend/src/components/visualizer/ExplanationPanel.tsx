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
    <div className="bg-[#0D0D14] border border-[#1C1C28] rounded-2xl p-5 flex flex-col gap-5 font-mono text-xs text-[#8E8EA0] leading-relaxed">
      {/* Header Info */}
      <div className="flex flex-col gap-2">
        <div className="flex items-center justify-between gap-2">
          <span className="text-[10px] font-bold uppercase tracking-wider text-[#A78BFA] bg-[#8B5CF6]/10 px-2 py-0.5 rounded border border-[#8B5CF6]/20">
            {algorithm.category}
          </span>
          <div className="flex items-center gap-1.5">
            <button
              onClick={onOpenQuiz}
              className="text-[11px] text-[#A78BFA] hover:text-white flex items-center gap-1 bg-[#12121B] hover:bg-[#181824] border border-[#1C1C28] px-2.5 py-1 rounded-lg transition-colors"
            >
              <BrainCircuit className="w-3.5 h-3.5" /> Practice Quiz
            </button>
            <button
              onClick={onOpenAI}
              className="text-[11px] text-emerald-400 hover:text-white flex items-center gap-1 bg-emerald-500/10 hover:bg-emerald-500/20 border border-emerald-500/30 px-2.5 py-1 rounded-lg transition-colors"
            >
              <Bot className="w-3.5 h-3.5" /> Explain Step
            </button>
          </div>
        </div>
        <h3 className="text-lg font-bold text-[#F4F4F6] font-sans">
          {algorithm.name}
        </h3>
        <p className="text-xs text-[#CBD5E1] font-sans leading-normal">
          {algorithm.description}
        </p>
      </div>

      {/* Dynamic Step Explanation Box */}
      {currentStepData && (
        <div className="p-3.5 bg-[#12121B] border border-[#8B5CF6]/30 rounded-xl space-y-1.5">
          <div className="flex items-center gap-2">
            <Zap className="w-3.5 h-3.5 text-[#A78BFA]" />
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
        <span className="text-[10px] font-bold uppercase tracking-wider text-[#5A5A6E] flex items-center gap-1">
          <Zap className="w-3 h-3 text-amber-400" /> Complexity Analysis
        </span>
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-2">
          <div className="p-2.5 bg-[#12121B] rounded-xl border border-[#1C1C28] flex flex-col gap-0.5">
            <span className="text-[9px] uppercase text-[#5A5A6E] font-bold">Best Time</span>
            <span className="text-xs font-bold text-emerald-400">{algorithm.complexity.best}</span>
          </div>
          <div className="p-2.5 bg-[#12121B] rounded-xl border border-[#1C1C28] flex flex-col gap-0.5">
            <span className="text-[9px] uppercase text-[#5A5A6E] font-bold">Average</span>
            <span className="text-xs font-bold text-amber-400">{algorithm.complexity.average}</span>
          </div>
          <div className="p-2.5 bg-[#12121B] rounded-xl border border-[#1C1C28] flex flex-col gap-0.5">
            <span className="text-[9px] uppercase text-[#5A5A6E] font-bold">Worst Time</span>
            <span className="text-xs font-bold text-rose-400">{algorithm.complexity.worst}</span>
          </div>
          <div className="p-2.5 bg-[#12121B] rounded-xl border border-[#1C1C28] flex flex-col gap-0.5">
            <span className="text-[9px] uppercase text-[#5A5A6E] font-bold">Space</span>
            <span className="text-xs font-bold text-[#A78BFA]">{algorithm.complexity.space}</span>
          </div>
        </div>
      </div>

      {/* Problem Solved & Use Cases */}
      <div className="space-y-2">
        <span className="text-[10px] font-bold uppercase tracking-wider text-[#5A5A6E] flex items-center gap-1">
          <Target className="w-3 h-3 text-[#8B5CF6]" /> Problem It Solves
        </span>
        <p className="text-xs text-[#CBD5E1] font-sans">
          {algorithm.problem_solved}
        </p>
      </div>

      {/* Use Cases */}
      {algorithm.use_cases && algorithm.use_cases.length > 0 && (
        <div className="space-y-2">
          <span className="text-[10px] font-bold uppercase tracking-wider text-[#5A5A6E]">
            Common Use Cases
          </span>
          <ul className="space-y-1 text-[11px] text-[#A78BFA] font-sans">
            {algorithm.use_cases.map((uc, idx) => (
              <li key={idx} className="flex items-start gap-1.5">
                <span className="text-[#8B5CF6]">•</span>
                <span>{uc}</span>
              </li>
            ))}
          </ul>
        </div>
      )}

      {/* Cross-Feature Links: Solve Problems in Practice Arena */}
      <div className="pt-2 border-t border-[#1C1C28] flex items-center justify-between">
        <button
          onClick={() => navigate('/practice')}
          className="text-xs text-[#A78BFA] hover:text-white flex items-center gap-1 font-sans font-medium"
        >
          <ExternalLink className="w-3.5 h-3.5" /> Practice related problems in Practice Arena
        </button>
        <button
          onClick={() => navigate('/code-review')}
          className="text-xs text-[#8E8EA0] hover:text-[#F4F4F6] flex items-center gap-1 font-sans"
        >
          Review Code with AI →
        </button>
      </div>
    </div>
  );
};
export default ExplanationPanel;

