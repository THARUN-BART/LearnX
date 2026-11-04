import React from 'react';
import { ComplexityComparison } from '../../types/codereview.types';
import { Copy, ArrowRight, Zap, Code2 } from 'lucide-react';
import { toast } from 'sonner';
import { useNavigate } from 'react-router-dom';

interface CodeDiffViewerProps {
  originalCode: string;
  optimizedCode: string;
  explanation: string;
  complexity: ComplexityComparison;
  problemId?: string;
  language: string;
}

export const CodeDiffViewer: React.FC<CodeDiffViewerProps> = ({
  originalCode,
  optimizedCode,
  explanation,
  complexity,
  problemId,
  language,
}) => {
  const navigate = useNavigate();

  const handleCopyOptimized = () => {
    navigator.clipboard.writeText(optimizedCode);
    toast.success('Optimized code copied to clipboard!');
  };

  const handleLoadIntoEditor = () => {
    if (problemId) {
      navigate(`/practice/${problemId}`, { state: { initialCode: optimizedCode, language } });
    } else {
      navigate('/practice', { state: { initialCode: optimizedCode, language } });
    }
    toast.success('Loaded optimized code into Practice Arena!');
  };

  return (
    <div className="bg-white border border-[#E2E8F0] rounded-2xl p-5 flex flex-col gap-5 font-mono">
      {/* Header */}
      <div className="flex flex-wrap items-center justify-between gap-3 border-b border-[#E2E8F0] pb-3">
        <div className="flex items-center gap-2">
          <Code2 className="w-4 h-4 text-[#EF4444]" />
          <span className="text-sm font-bold text-white font-sans">
            Before & After Code Optimization
          </span>
        </div>

        <div className="flex items-center gap-2">
          <button
            onClick={handleCopyOptimized}
            className="text-xs text-[#64748B] hover:text-[#0F172A] flex items-center gap-1.5 bg-[#F1F5F9] hover:bg-[#E2E8F0] border border-[#E2E8F0] px-3 py-1.5 rounded-lg transition-colors"
          >
            <Copy className="w-3.5 h-3.5" /> Copy Optimized
          </button>
          <button
            onClick={handleLoadIntoEditor}
            className="text-xs text-white font-bold flex items-center gap-1.5 bg-[#EF4444] hover:bg-[#EF4444] px-3.5 py-1.5 rounded-lg shadow-md shadow-[#EF4444]/30 transition-all font-sans"
          >
            Load in Practice Arena <ArrowRight className="w-3.5 h-3.5" />
          </button>
        </div>
      </div>

      {/* Complexity Comparison Table */}
      <div className="p-4 bg-[#F1F5F9] border border-[#E2E8F0] rounded-xl flex flex-col gap-3">
        <div className="flex items-center gap-1.5 text-xs font-bold text-[#E9D5FF] font-sans">
          <Zap className="w-4 h-4 text-amber-400" /> Asymptotic Complexity Comparison
        </div>
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-2 text-xs">
          <div className="p-2.5 bg-white rounded-lg border border-[#E2E8F0]">
            <span className="text-[10px] text-[#94A3B8] block font-bold">ORIGINAL TIME</span>
            <span className="text-rose-400 font-bold">{complexity.original_time}</span>
          </div>
          <div className="p-2.5 bg-white rounded-lg border border-[#E2E8F0]">
            <span className="text-[10px] text-[#94A3B8] block font-bold">OPTIMIZED TIME</span>
            <span className="text-emerald-400 font-bold">{complexity.optimized_time}</span>
          </div>
          <div className="p-2.5 bg-white rounded-lg border border-[#E2E8F0]">
            <span className="text-[10px] text-[#94A3B8] block font-bold">ORIGINAL SPACE</span>
            <span className="text-[#DC2626] font-bold">{complexity.original_space}</span>
          </div>
          <div className="p-2.5 bg-white rounded-lg border border-[#E2E8F0]">
            <span className="text-[10px] text-[#94A3B8] block font-bold">OPTIMIZED SPACE</span>
            <span className="text-emerald-400 font-bold">{complexity.optimized_space}</span>
          </div>
        </div>
        {complexity.explanation && (
          <p className="text-xs text-[#CBD5E1] font-sans pt-1 border-t border-[#E2E8F0]/60">
            {complexity.explanation}
          </p>
        )}
      </div>

      {/* Explanation of Improvements */}
      {explanation && (
        <div className="p-3.5 bg-[#EF4444]/5 border border-[#EF4444]/20 rounded-xl text-xs font-sans text-[#E2E8F0]">
          <span className="font-bold text-[#DC2626]">Why this is better: </span>
          {explanation}
        </div>
      )}

      {/* Side-by-Side Code Boxes */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {/* Original Code */}
        <div className="flex flex-col gap-2">
          <div className="flex items-center justify-between text-xs text-[#64748B] px-1">
            <span className="font-bold uppercase text-[10px] text-rose-400">Original Submission</span>
            <span>{language}</span>
          </div>
          <pre className="p-4 bg-[#F8FAFC] border border-[#E2E8F0] rounded-xl text-xs text-[#CBD5E1] overflow-x-auto max-h-[380px] overflow-y-auto whitespace-pre font-mono leading-relaxed">
            {originalCode}
          </pre>
        </div>

        {/* Optimized Code */}
        <div className="flex flex-col gap-2">
          <div className="flex items-center justify-between text-xs text-[#64748B] px-1">
            <span className="font-bold uppercase text-[10px] text-emerald-400">Optimized Recommendation</span>
            <span>{language}</span>
          </div>
          <pre className="p-4 bg-[#F8FAFC] border border-emerald-500/30 rounded-xl text-xs text-emerald-300 overflow-x-auto max-h-[380px] overflow-y-auto whitespace-pre font-mono leading-relaxed shadow-inner shadow-emerald-500/5">
            {optimizedCode}
          </pre>
        </div>
      </div>
    </div>
  );
};
export default CodeDiffViewer;

