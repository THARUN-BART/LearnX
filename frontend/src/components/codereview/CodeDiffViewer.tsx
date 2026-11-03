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
    <div className="bg-[#0D0D14] border border-[#1C1C28] rounded-2xl p-5 flex flex-col gap-5 font-mono">
      {/* Header */}
      <div className="flex flex-wrap items-center justify-between gap-3 border-b border-[#1C1C28] pb-3">
        <div className="flex items-center gap-2">
          <Code2 className="w-4 h-4 text-[#8B5CF6]" />
          <span className="text-sm font-bold text-white font-sans">
            Before & After Code Optimization
          </span>
        </div>

        <div className="flex items-center gap-2">
          <button
            onClick={handleCopyOptimized}
            className="text-xs text-[#8E8EA0] hover:text-[#F4F4F6] flex items-center gap-1.5 bg-[#12121B] hover:bg-[#181824] border border-[#1C1C28] px-3 py-1.5 rounded-lg transition-colors"
          >
            <Copy className="w-3.5 h-3.5" /> Copy Optimized
          </button>
          <button
            onClick={handleLoadIntoEditor}
            className="text-xs text-white font-bold flex items-center gap-1.5 bg-[#8B5CF6] hover:bg-[#7C3AED] px-3.5 py-1.5 rounded-lg shadow-md shadow-[#8B5CF6]/30 transition-all font-sans"
          >
            Load in Practice Arena <ArrowRight className="w-3.5 h-3.5" />
          </button>
        </div>
      </div>

      {/* Complexity Comparison Table */}
      <div className="p-4 bg-[#12121B] border border-[#1C1C28] rounded-xl flex flex-col gap-3">
        <div className="flex items-center gap-1.5 text-xs font-bold text-[#E9D5FF] font-sans">
          <Zap className="w-4 h-4 text-amber-400" /> Asymptotic Complexity Comparison
        </div>
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-2 text-xs">
          <div className="p-2.5 bg-[#0D0D14] rounded-lg border border-[#1C1C28]">
            <span className="text-[10px] text-[#5A5A6E] block font-bold">ORIGINAL TIME</span>
            <span className="text-rose-400 font-bold">{complexity.original_time}</span>
          </div>
          <div className="p-2.5 bg-[#0D0D14] rounded-lg border border-[#1C1C28]">
            <span className="text-[10px] text-[#5A5A6E] block font-bold">OPTIMIZED TIME</span>
            <span className="text-emerald-400 font-bold">{complexity.optimized_time}</span>
          </div>
          <div className="p-2.5 bg-[#0D0D14] rounded-lg border border-[#1C1C28]">
            <span className="text-[10px] text-[#5A5A6E] block font-bold">ORIGINAL SPACE</span>
            <span className="text-[#A78BFA] font-bold">{complexity.original_space}</span>
          </div>
          <div className="p-2.5 bg-[#0D0D14] rounded-lg border border-[#1C1C28]">
            <span className="text-[10px] text-[#5A5A6E] block font-bold">OPTIMIZED SPACE</span>
            <span className="text-emerald-400 font-bold">{complexity.optimized_space}</span>
          </div>
        </div>
        {complexity.explanation && (
          <p className="text-xs text-[#CBD5E1] font-sans pt-1 border-t border-[#1C1C28]/60">
            {complexity.explanation}
          </p>
        )}
      </div>

      {/* Explanation of Improvements */}
      {explanation && (
        <div className="p-3.5 bg-[#8B5CF6]/5 border border-[#8B5CF6]/20 rounded-xl text-xs font-sans text-[#E2E8F0]">
          <span className="font-bold text-[#A78BFA]">Why this is better: </span>
          {explanation}
        </div>
      )}

      {/* Side-by-Side Code Boxes */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {/* Original Code */}
        <div className="flex flex-col gap-2">
          <div className="flex items-center justify-between text-xs text-[#8E8EA0] px-1">
            <span className="font-bold uppercase text-[10px] text-rose-400">Original Submission</span>
            <span>{language}</span>
          </div>
          <pre className="p-4 bg-[#07070B] border border-[#1C1C28] rounded-xl text-xs text-[#CBD5E1] overflow-x-auto max-h-[380px] overflow-y-auto whitespace-pre font-mono leading-relaxed">
            {originalCode}
          </pre>
        </div>

        {/* Optimized Code */}
        <div className="flex flex-col gap-2">
          <div className="flex items-center justify-between text-xs text-[#8E8EA0] px-1">
            <span className="font-bold uppercase text-[10px] text-emerald-400">Optimized Recommendation</span>
            <span>{language}</span>
          </div>
          <pre className="p-4 bg-[#07070B] border border-emerald-500/30 rounded-xl text-xs text-emerald-300 overflow-x-auto max-h-[380px] overflow-y-auto whitespace-pre font-mono leading-relaxed shadow-inner shadow-emerald-500/5">
            {optimizedCode}
          </pre>
        </div>
      </div>
    </div>
  );
};
export default CodeDiffViewer;

