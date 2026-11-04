import React, { useState, useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { CodeReviewResult } from '../../types/codereview.types';
import codeReviewService from '../../services/codereview.service';
import ScoreCards from '../../components/codereview/ScoreCards';
import IssueList from '../../components/codereview/IssueList';
import CodeDiffViewer from '../../components/codereview/CodeDiffViewer';
import ReviewHistoryDrawer from '../../components/codereview/ReviewHistoryDrawer';
import {
  Sparkles,
  Bot,
  Play,
  RotateCcw,
  History,
  CheckCircle2,
  GitBranch,
  ArrowRight,
  Eye,
} from 'lucide-react';
import Button from '../../components/ui/Button';
import { toast } from 'sonner';

const SAMPLE_CODE = `def two_sum(nums, target):
    # Find two indices that sum up to target
    for i in range(len(nums)):
        for j in range(len(nums)):
            if i != j and nums[i] + nums[j] == target:
                return [i, j]
    return []
`;

export const CodeReviewPage: React.FC = () => {
  const location = useLocation();
  const navigate = useNavigate();

  const [language, setLanguage] = useState<string>(location.state?.language || 'python');
  const [sourceCode, setSourceCode] = useState<string>(location.state?.sourceCode || SAMPLE_CODE);
  const [problemId, setProblemId] = useState<string>(location.state?.problemId || '');
  const [problemTitle, setProblemTitle] = useState<string>(location.state?.problemTitle || 'Custom Code Review');

  const [loading, setLoading] = useState(false);
  const [reviewResult, setReviewResult] = useState<CodeReviewResult | null>(null);
  const [historyOpen, setHistoryOpen] = useState(false);

  useEffect(() => {
    // If navigated with existing review id
    if (location.state?.reviewId) {
      codeReviewService.getReviewById(location.state.reviewId).then((res) => {
        if (res) {
          setReviewResult(res);
          setSourceCode(res.source_code);
          setLanguage(res.language);
          if (res.problem_title) setProblemTitle(res.problem_title);
        }
      });
    }
  }, [location.state]);

  const handleReviewCode = async () => {
    if (!sourceCode.trim()) {
      toast.error('Please enter code to review');
      return;
    }
    setLoading(true);
    try {
      const res = await codeReviewService.reviewCode({
        problem_id: problemId || undefined,
        language,
        source_code: sourceCode,
      });
      if (res) {
        setReviewResult(res);
        toast.success('AI Code Review Complete!');
      } else {
        toast.error('Could not complete code review');
      }
    } catch {
      toast.error('Review service error');
    } finally {
      setLoading(false);
    }
  };

  const handleSelectPastReview = async (id: string) => {
    setLoading(true);
    try {
      const res = await codeReviewService.getReviewById(id);
      if (res) {
        setReviewResult(res);
        setSourceCode(res.source_code);
        setLanguage(res.language);
        if (res.problem_title) setProblemTitle(res.problem_title);
        toast.info('Loaded review from history');
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex flex-col gap-6 font-mono text-[#0F172A] pb-12">
      {/* Hero Banner */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 p-6 bg-gradient-to-r from-[#F1F5F9] via-[#FFFFFF] to-[#F1F5F9] border border-[#E2E8F0] rounded-2xl relative overflow-hidden">
        <div className="space-y-1 relative z-10">
          <div className="flex items-center gap-2">
            <span className="text-[10px] uppercase font-bold tracking-widest text-[#EF4444] bg-[#EF4444]/10 px-2 py-0.5 rounded border border-[#EF4444]/20">
              AI Code Review Assistant
            </span>
            <span className="text-[10px] text-emerald-400 font-bold">
              Multi-Language Static & LLM Analysis
            </span>
          </div>
          <h1 className="text-xl sm:text-2xl font-bold font-sans tracking-tight text-white">
            Automated Code Quality, Security & Optimization Audit
          </h1>
          <p className="text-xs text-[#64748B] font-sans max-w-xl">
            Audit time/space bottlenecks, detect anti-patterns, generate clean production refactorings, and connect to interactive algorithm visualizers.
          </p>
        </div>

        <div className="flex items-center gap-2 relative z-10">
          <Button
            variant="secondary"
            size="sm"
            onClick={() => setHistoryOpen(true)}
            leftIcon={<History className="w-3.5 h-3.5" />}
            className="text-xs"
          >
            Past Reviews
          </Button>
        </div>
      </div>

      {/* Editor & Review Trigger Section */}
      <div className="bg-white border border-[#E2E8F0] rounded-2xl p-5 flex flex-col gap-4">
        <div className="flex flex-wrap items-center justify-between gap-3 border-b border-[#E2E8F0] pb-3">
          <div className="flex items-center gap-3">
            <span className="text-xs font-bold text-white font-sans">
              Target Code for Review
            </span>
            <span className="text-[11px] text-[#DC2626] bg-[#F1F5F9] px-2.5 py-0.5 rounded-lg border border-[#E2E8F0]">
              {problemTitle}
            </span>
          </div>

          <div className="flex items-center gap-2">
            <select
              value={language}
              onChange={(e) => setLanguage(e.target.value)}
              className="bg-[#F1F5F9] text-[#0F172A] text-xs px-3 py-1.5 rounded-lg border border-[#E2E8F0] outline-none"
            >
              <option value="python">Python</option>
              <option value="cpp">C++</option>
              <option value="java">Java</option>
              <option value="c">C</option>
              <option value="go">Go</option>
            </select>

            <Button
              variant="ghost"
              size="sm"
              onClick={() => setSourceCode(SAMPLE_CODE)}
              leftIcon={<RotateCcw className="w-3.5 h-3.5" />}
              className="text-xs text-[#64748B]"
            >
              Load Sample
            </Button>

            <Button
              variant="primary"
              size="sm"
              onClick={handleReviewCode}
              isLoading={loading}
              leftIcon={<Sparkles className="w-3.5 h-3.5" />}
              className="text-xs shadow-md shadow-[#EF4444]/30"
            >
              Review My Code
            </Button>
          </div>
        </div>

        {/* Code Input Area */}
        <textarea
          rows={9}
          value={sourceCode}
          onChange={(e) => setSourceCode(e.target.value)}
          placeholder="Paste or write your solution code here for comprehensive AI review..."
          className="w-full bg-[#F8FAFC] text-[#CBD5E1] font-mono text-xs p-4 rounded-xl border border-[#E2E8F0] outline-none focus:border-[#EF4444]/60 leading-relaxed resize-y"
        />
      </div>

      {/* Structured Review Results Dashboard */}
      {reviewResult && (
        <div className="space-y-6">
          {/* Score Cards */}
          <ScoreCards scores={reviewResult.scores} />

          {/* Cross-Feature Recommendation to Algorithm Visualizer */}
          {reviewResult.recommended_visualizer_algo_id && (
            <div className="p-4 rounded-2xl bg-gradient-to-r from-[#EF4444]/15 via-[#DC2626]/10 to-[#B91C1C]/15 border border-[#EF4444]/30 flex flex-col sm:flex-row sm:items-center justify-between gap-3 shadow-lg shadow-[#EF4444]/10">
              <div className="flex items-start gap-3">
                <div className="w-8 h-8 rounded-xl bg-[#EF4444] text-white flex items-center justify-center shrink-0 shadow-md shadow-[#EF4444]/30">
                  <Eye className="w-4 h-4" />
                </div>
                <div>
                  <div className="flex items-center gap-2">
                    <span className="text-[10px] font-bold uppercase tracking-wider text-[#D8B4FE]">
                      Recommended Visualizer
                    </span>
                  </div>
                  <p className="text-xs text-[#E2E8F0] font-sans mt-0.5">
                    {reviewResult.recommended_visualizer_reason}
                  </p>
                </div>
              </div>

              <button
                onClick={() =>
                  navigate(`/visualizer?algo=${reviewResult.recommended_visualizer_algo_id}`)
                }
                className="px-4 py-2 bg-[#EF4444] hover:bg-[#EF4444] text-white rounded-xl text-xs font-bold font-sans flex items-center gap-2 shrink-0 shadow-md shadow-[#EF4444]/30 transition-all"
              >
                Explore Visualizer <ArrowRight className="w-3.5 h-3.5" />
              </button>
            </div>
          )}

          {/* Correctness Observations */}
          {reviewResult.correctness_observations && reviewResult.correctness_observations.length > 0 && (
            <div className="p-4 bg-white border border-[#E2E8F0] rounded-2xl space-y-2">
              <span className="text-[10px] uppercase font-bold text-[#94A3B8] tracking-wider">
                Correctness & Structural Observations
              </span>
              <ul className="space-y-1.5 text-xs text-[#CBD5E1] font-sans">
                {reviewResult.correctness_observations.map((obs, idx) => (
                  <li key={idx} className="flex items-start gap-2">
                    <CheckCircle2 className="w-3.5 h-3.5 text-emerald-400 shrink-0 mt-0.5" />
                    <span>{obs}</span>
                  </li>
                ))}
              </ul>
            </div>
          )}

          {/* Issues List */}
          <IssueList issues={reviewResult.issues} />

          {/* Before & After Diff Viewer */}
          <CodeDiffViewer
            originalCode={reviewResult.source_code}
            optimizedCode={reviewResult.optimized_code}
            explanation={reviewResult.explanation_of_improvements}
            complexity={reviewResult.complexity}
            problemId={reviewResult.problem_id}
            language={reviewResult.language}
          />

          {/* Alternative Approaches */}
          {reviewResult.alternative_approaches && reviewResult.alternative_approaches.length > 0 && (
            <div className="bg-white border border-[#E2E8F0] rounded-2xl p-5 flex flex-col gap-4">
              <div className="flex items-center gap-2 border-b border-[#E2E8F0] pb-3">
                <GitBranch className="w-4 h-4 text-[#EF4444]" />
                <span className="text-sm font-bold text-white font-sans">
                  Alternative Algorithmic Paradigms & Trade-offs
                </span>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                {reviewResult.alternative_approaches.map((alt, idx) => (
                  <div key={idx} className="p-4 bg-[#F1F5F9] border border-[#E2E8F0] rounded-xl flex flex-col gap-2">
                    <div className="flex items-center justify-between">
                      <span className="text-xs font-bold text-[#0F172A] font-sans">{alt.name}</span>
                      <div className="flex items-center gap-2 text-[10px]">
                        <span className="text-emerald-400 font-bold">{alt.time}</span>
                        <span>·</span>
                        <span className="text-[#DC2626]">{alt.space}</span>
                      </div>
                    </div>
                    <p className="text-xs text-[#CBD5E1] font-sans leading-relaxed">{alt.description}</p>
                    <div className="text-[11px] text-[#94A3B8] font-sans pt-1 border-t border-[#E2E8F0]/60">
                      <span className="text-amber-400 font-bold">Trade-off: </span>
                      {alt.trade_off}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      )}

      {/* Review History Drawer */}
      <ReviewHistoryDrawer
        isOpen={historyOpen}
        onClose={() => setHistoryOpen(false)}
        onSelectReview={handleSelectPastReview}
      />
    </div>
  );
};
export default CodeReviewPage;

