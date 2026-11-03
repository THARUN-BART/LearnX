import { normalizePercentage } from '../../utils/formatters';
import React, { Suspense, lazy, useEffect, useMemo, useState } from 'react';
import { useParams, useNavigate, useLocation } from 'react-router-dom';
import {
  Play,
  Send,
  RotateCcw,
  Copy,
  Maximize2,
  Minimize2,
  Terminal,
  FileCode,
  Sparkles,
  Bot,
  X,
  ExternalLink,
  Bookmark,
} from 'lucide-react';
import { PROGRAMMING_LANGUAGES } from '../../utils/constants';
import {
  Problem,
  CodeExecutionResult,
  ProgressiveHint,
  DryRunResult,
  ComplexityResult,
  ExplainMistakeResult,
  PedagogyAnalysis,
  SubmissionHistoryItem,
  ProgressStatus,
} from '../../types/learning.types';
import codingService from '../../services/coding.service';
import learningService from '../../services/learning.service';
import aiService from '../../services/ai.service';
import progressService from '../../services/progress.service';
import { recommendationService } from '../../services/recommendation.service';
import { ProblemRecommendation } from '../../types/recommendations.types';
import { starterCodeKey, toJudgeLanguage, toMonacoLanguage } from '../../utils/language';
import Button from '../../components/ui/Button';
import { toast } from 'sonner';

const Editor = lazy(() => import('@monaco-editor/react'));

function starterFor(problem: Problem, lang: string): string {
  const key = starterCodeKey(lang);
  return problem.starterCode?.[key] || problem.starterCode?.python || '';
}

function formatInlineCode(text: string) {
  const parts = text.split(/(`[^`]+`)/g);
  return parts.map((part, idx) => {
    if (part.startsWith('`') && part.endsWith('`')) {
      return (
        <code key={idx} className="px-1.5 py-0.5 rounded bg-[#181824] text-[#A78BFA] font-mono text-[11px] border border-[#1C1C28]">
          {part.slice(1, -1)}
        </code>
      );
    }
    return part;
  });
}

function RenderProblemText({ text }: { text: string }) {
  const lines = text.split('\n');
  return (
    <div className="space-y-2 font-sans leading-relaxed">
      {lines.map((line, idx) => {
        const trimmed = line.trim();
        if (!trimmed) return <div key={idx} className="h-1" />;
        if (trimmed.startsWith('### ')) {
          return (
            <h4 key={idx} className="text-xs font-bold text-[#E9D5FF] uppercase tracking-wider pt-2">
              {trimmed.replace(/^###\s+/, '')}
            </h4>
          );
        }
        if (trimmed.startsWith('## ')) {
          return (
            <h3 key={idx} className="text-sm font-bold text-[#F4F4F6] pt-2">
              {trimmed.replace(/^##\s+/, '')}
            </h3>
          );
        }
        if (trimmed.startsWith('- ')) {
          return (
            <div key={idx} className="flex items-start gap-2 pl-2 text-[11px] text-[#CBD5E1]">
              <span className="text-[#8B5CF6]">•</span>
              <span>{formatInlineCode(trimmed.slice(2))}</span>
            </div>
          );
        }
        return (
          <p key={idx} className="text-xs text-[#CBD5E1]">
            {formatInlineCode(trimmed)}
          </p>
        );
      })}
    </div>
  );
}

export const PracticePage: React.FC = () => {
  const { problemId } = useParams<{ problemId: string }>();
  const navigate = useNavigate();
  const location = useLocation();

  const [problem, setProblem] = useState<Problem | null>(null);
  const [problemLoading, setProblemLoading] = useState(true);
  const [selectedLanguage, setSelectedLanguage] = useState('python');
  const [code, setCode] = useState('');
  const [activeTab, setActiveTab] = useState<'description' | 'hints' | 'history' | 'related'>('description');
  const [activeBottomTab, setActiveBottomTab] = useState<'testcases' | 'console'>('testcases');
  const [selectedTestCaseIndex, setSelectedTestCaseIndex] = useState(0);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [isRunning, setIsRunning] = useState(false);
  const [executionResult, setExecutionResult] = useState<CodeExecutionResult | null>(null);
  const [allProblems, setAllProblems] = useState<Problem[]>([]);
  const [problemQuery, setProblemQuery] = useState('');
  const [unlockedHints, setUnlockedHints] = useState<ProgressiveHint[]>([]);
  const [isLoadingHint, setIsLoadingHint] = useState(false);
  const [history, setHistory] = useState<SubmissionHistoryItem[]>([]);
  const [pedagogy, setPedagogy] = useState<PedagogyAnalysis | null>(null);
  const [similarRecs, setSimilarRecs] = useState<ProblemRecommendation[]>([]);
  const [statusFilter, setStatusFilter] = useState<string>(location.state?.filterStatus || '');
  const [progressStatus, setProgressStatus] = useState<ProgressStatus>('Not Started');
  const [aiModal, setAiModal] = useState<{
    isOpen: boolean;
    title: string;
    type: 'explain' | 'dryrun' | 'complexity' | 'pedagogy';
    loading: boolean;
    explain?: ExplainMistakeResult;
    dryRun?: DryRunResult;
    complexity?: ComplexityResult;
  }>({ isOpen: false, title: '', type: 'explain', loading: false });

  useEffect(() => {
    if (statusFilter) {
      learningService.getProblems({ status: statusFilter, limit: 50 }).then((list) => {
        if (Array.isArray(list)) setAllProblems(list);
      });
    } else {
      learningService.getProblems({ limit: 40 }).then((list) => {
        if (Array.isArray(list)) setAllProblems(list);
      });
    }
  }, [statusFilter]);

  useEffect(() => {
    const q = problemQuery.trim();
    if (q.length < 2) return;
    const t = setTimeout(async () => {
      const found = await learningService.searchProblems(q, 20);
      if (found.length) setAllProblems(found);
    }, 250);
    return () => clearTimeout(t);
  }, [problemQuery]);

  useEffect(() => {
    let mounted = true;
    const load = async () => {
      setProblemLoading(true);
      let prob: Problem | null = null;
      if (problemId) {
        prob = await learningService.getProblemById(problemId);
      } else {
        const list = allProblems.length ? allProblems : await learningService.getProblems({ limit: 1 });
        prob = list[0] || null;
        if (prob && mounted) navigate(`/practice/${prob.id}`, { replace: true });
      }
      if (!mounted) return;
      if (!prob) {
        setProblem(null);
        setProblemLoading(false);
        return;
      }
      setProblem(prob);
      if (location.state?.initialCode) {
        setCode(location.state.initialCode);
      } else {
        setCode(starterFor(prob, selectedLanguage));
      }
      setExecutionResult(null);
      setPedagogy(null);
      setProgressStatus((prob.progress_status as ProgressStatus) || 'Not Started');
      setSelectedTestCaseIndex(0);
      setActiveTab('description');
      setProblemLoading(false);

      const [hints, subs, prog] = await Promise.allSettled([
        aiService.getUnlockedHints(prob.id),
        codingService.getProblemSubmissions(prob.id),
        progressService.getProblemProgress(prob.id),
      ]);
      if (!mounted) return;
      setUnlockedHints(hints.status === 'fulfilled' && Array.isArray(hints.value) ? hints.value : []);
      setHistory(subs.status === 'fulfilled' ? subs.value : []);
      if (prog.status === 'fulfilled' && prog.value?.status) {
        setProgressStatus(prog.value.status);
      }
      recommendationService.getSimilarProblems(prob.id).then((sim) => {
        if (mounted && Array.isArray(sim)) setSimilarRecs(sim);
      });
    };
    load();
    return () => {
      mounted = false;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [problemId]);

  const related = useMemo(() => (problem?.related || []).slice(0, 6), [problem]);

  const handleLanguageChange = (lang: string) => {
    setSelectedLanguage(lang);
    if (problem) setCode(starterFor(problem, lang));
  };

  const handleResetCode = () => {
    if (!problem) return;
    setCode(starterFor(problem, selectedLanguage));
    toast.info('Starter code reset');
  };

  const handleCopyCode = () => {
    navigator.clipboard.writeText(code);
    toast.success('Code copied to clipboard');
  };

  const handleUnlockNextHint = async () => {
    if (!problem) return;
    const nextLevel = unlockedHints.length + 1;
    if (nextLevel > 5) {
      toast.info('All 5 progressive hint levels are already unlocked.');
      return;
    }
    setIsLoadingHint(true);
    try {
      const hint = await aiService.getProgressiveHint(problem.id, nextLevel, executionResult?.submissionId);
      setUnlockedHints((prev) => [...prev, hint]);
      toast.success(`Unlocked Hint Level ${nextLevel}: ${hint.title}`);
    } catch {
      toast.error('Failed to unlock hint from AI Tutor');
    } finally {
      setIsLoadingHint(false);
    }
  };

  const handleExplainMistake = async (submissionId: string) => {
    setAiModal({ isOpen: true, title: 'AI Root Cause Diagnosis', type: 'explain', loading: true });
    try {
      const data = await aiService.explainMistake(submissionId);
      setAiModal((prev) => ({ ...prev, loading: false, explain: data }));
    } catch {
      toast.error('Could not retrieve root cause explanation');
      setAiModal((prev) => ({ ...prev, loading: false }));
    }
  };

  const handleDryRun = async () => {
    if (!problem) return;
    setAiModal({ isOpen: true, title: 'Execution Trace & Dry Run', type: 'dryrun', loading: true });
    try {
      const sample = problem.testCases[0]?.input;
      const data = await aiService.getDryRun(problem.id, code, selectedLanguage, sample);
      setAiModal((prev) => ({ ...prev, loading: false, dryRun: data }));
    } catch {
      toast.error('Could not generate dry run trace');
      setAiModal((prev) => ({ ...prev, loading: false }));
    }
  };

  const handleComplexityCheck = async () => {
    if (!problem) return;
    setAiModal({ isOpen: true, title: 'Time & Space Complexity Audit', type: 'complexity', loading: true });
    try {
      const data = await aiService.getComplexity(problem.id, code, selectedLanguage);
      setAiModal((prev) => ({ ...prev, loading: false, complexity: data }));
    } catch {
      toast.error('Could not generate complexity audit');
      setAiModal((prev) => ({ ...prev, loading: false }));
    }
  };

  const handlePedagogy = async () => {
    if (!problem) return;
    const sample = problem.testCases[0];
    const testCode = sample
      ? `# sample assertion based on visible case\nprint(${JSON.stringify(sample.expectedOutput)})`
      : 'assert True';
    setAiModal({ isOpen: true, title: 'Pedagogy Diagnosis', type: 'pedagogy', loading: true });
    try {
      const data = await aiService.analyzePedagogy(code, testCode, 'L1');
      setPedagogy(data);
      setAiModal((prev) => ({ ...prev, loading: false }));
    } catch {
      toast.error('Pedagogy service unavailable');
      setAiModal((prev) => ({ ...prev, loading: false }));
    }
  };

  const runOrSubmit = async (submit: boolean) => {
    if (!problem) return;
    setIsRunning(true);
    setActiveBottomTab('console');
    try {
      await progressService.recordAttempt(problem.id);
      const result = submit
        ? await codingService.submitCode({ problemId: problem.id, language: selectedLanguage, code })
        : await codingService.runCode({ problemId: problem.id, language: selectedLanguage, code });
      setExecutionResult(result);
      if (result.status === 'Accepted') {
        setProgressStatus('Solved');
        toast.success(submit ? 'Solution Accepted!' : 'Test run passed!');
      } else {
        if (progressStatus === 'Not Started') setProgressStatus('Attempted');
        toast.error(`Status: ${result.status}`);
      }
      const subs = await codingService.getProblemSubmissions(problem.id);
      setHistory(subs);
    } catch {
      toast.error(submit ? 'Submission failed' : 'Execution failed');
    } finally {
      setIsRunning(false);
    }
  };

  const markStatus = async (status: ProgressStatus) => {
    if (!problem) return;
    try {
      await progressService.updateProblemStatus(problem.id, status, status === 'Not Started');
      setProgressStatus(status);
      toast.success(`Marked ${status}`);
    } catch {
      toast.error('Could not update progress');
    }
  };

  if (problemLoading && !problem) {
    return <div className="py-20 text-center text-xs font-mono text-[#5A5A6E]">Loading problem workspace…</div>;
  }
  if (!problem) {
    return (
      <div className="py-20 text-center text-xs font-mono text-[#8E8EA0] space-y-3">
        <p>No problem found. Import the A2Z sheet on the backend, then refresh.</p>
        <Button size="sm" onClick={() => navigate('/roadmap')}>Open Roadmap</Button>
      </div>
    );
  }

  const tests = problem.testCases || [];

  return (
    <div className={`flex flex-col gap-3 font-mono ${isFullscreen ? 'fixed inset-0 z-50 bg-[#07070B] p-4 h-screen' : 'h-[calc(100vh-6.5rem)]'}`}>
      <div className="flex flex-wrap items-center justify-between gap-3 bg-[#0D0D14] border border-[#1C1C28] px-4 py-2.5 rounded-xl shrink-0">
        <div className="flex items-center gap-3 min-w-0">
          <input
            value={problemQuery}
            onChange={(e) => setProblemQuery(e.target.value)}
            placeholder="Search problems"
            className="w-28 sm:w-36 bg-[#12121B] text-[#F4F4F6] text-[11px] px-2 py-1.5 rounded-lg border border-[#1C1C28] outline-none"
          />
          <select
            value={statusFilter}
            onChange={(e) => setStatusFilter(e.target.value)}
            className="bg-[#12121B] text-[#A78BFA] text-[11px] px-2 py-1.5 rounded-lg border border-[#1C1C28] outline-none font-mono"
          >
            <option value="">All Statuses</option>
            <option value="Needs Revision">Needs Revision</option>
            <option value="Attempted">Attempted</option>
            <option value="Solved">Solved</option>
            <option value="Not Started">Not Started</option>
          </select>
          <select
            value={problem.id}
            onChange={(e) => navigate(`/practice/${e.target.value}`)}
            className="bg-[#12121B] text-[#F4F4F6] text-xs font-semibold px-3 py-1.5 rounded-lg border border-[#1C1C28] outline-none max-w-[220px] truncate"
          >
            {(allProblems.length ? allProblems : [problem]).map((p) => (
              <option key={p.id} value={p.id}>
                {p.title} ({p.difficulty})
              </option>
            ))}
          </select>
          <span className={`text-[10px] px-2 py-0.5 rounded ${
            problem.difficulty === 'Easy'
              ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20'
              : problem.difficulty === 'Hard'
              ? 'bg-red-500/10 text-red-400 border border-red-500/20'
              : 'bg-amber-500/10 text-amber-400 border border-amber-500/20'
          }`}>
            {problem.difficulty}
          </span>
          <span className="text-[10px] text-[#A78BFA] hidden sm:inline">{progressStatus}</span>
        </div>

        <div className="flex items-center gap-2">
          <select
            value={selectedLanguage}
            onChange={(e) => handleLanguageChange(e.target.value)}
            className="bg-[#12121B] text-[#F4F4F6] text-xs px-3 py-1.5 rounded-lg border border-[#1C1C28]"
          >
            {PROGRAMMING_LANGUAGES.map((l) => (
              <option key={l.id} value={toJudgeLanguage(l.id)}>
                {l.name}
              </option>
            ))}
          </select>
          <button onClick={handleResetCode} title="Reset" className="p-1.5 text-[#8E8EA0] hover:text-[#F4F4F6] rounded-lg border border-[#1C1C28]">
            <RotateCcw className="w-3.5 h-3.5" />
          </button>
          <button onClick={handleCopyCode} title="Copy" className="p-1.5 text-[#8E8EA0] hover:text-[#F4F4F6] rounded-lg border border-[#1C1C28]">
            <Copy className="w-3.5 h-3.5" />
          </button>
          <button onClick={() => setIsFullscreen(!isFullscreen)} title="Fullscreen" className="p-1.5 text-[#8E8EA0] hover:text-[#F4F4F6] rounded-lg border border-[#1C1C28]">
            {isFullscreen ? <Minimize2 className="w-3.5 h-3.5" /> : <Maximize2 className="w-3.5 h-3.5" />}
          </button>
          <Button
            variant="secondary"
            size="sm"
            onClick={() =>
              navigate('/code-review', {
                state: {
                  problemId: problem.id,
                  problemTitle: problem.title,
                  language: selectedLanguage,
                  sourceCode: code,
                },
              })
            }
            leftIcon={<Sparkles className="w-3 h-3 text-[#A78BFA]" />}
            className="text-xs h-8 text-[#A78BFA] border-[#8B5CF6]/30 hover:border-[#8B5CF6]"
          >
            Review Code
          </Button>
          <Button variant="secondary" size="sm" onClick={() => runOrSubmit(false)} isLoading={isRunning} leftIcon={<Play className="w-3 h-3" />} className="text-xs h-8">
            Run
          </Button>
          <Button variant="primary" size="sm" onClick={() => runOrSubmit(true)} isLoading={isRunning} leftIcon={<Send className="w-3 h-3" />} className="text-xs h-8">
            Submit
          </Button>
        </div>
      </div>

      <div className="flex-1 grid grid-cols-1 lg:grid-cols-12 gap-3 min-h-0">
        <div className="lg:col-span-5 bg-[#0D0D14] border border-[#1C1C28] rounded-xl flex flex-col overflow-hidden">
          <div className="flex items-center justify-between border-b border-[#1C1C28] px-3 py-2 bg-[#0A0A10] gap-2 overflow-x-auto">
            <div className="flex items-center gap-1">
              {(['description', 'hints', 'history', 'related'] as const).map((tab) => (
                <button
                  key={tab}
                  onClick={() => setActiveTab(tab)}
                  className={`px-2.5 py-1 rounded-lg text-[11px] capitalize ${
                    activeTab === tab ? 'bg-[#181824] text-[#A78BFA] border border-[#8B5CF6]/30' : 'text-[#8E8EA0]'
                  }`}
                >
                  {tab === 'hints' ? `Hints (${unlockedHints.length}/5)` : tab}
                </button>
              ))}
            </div>
            <button onClick={() => navigate('/assistant')} className="text-[11px] text-[#8B5CF6] flex items-center gap-1 shrink-0">
              <Bot className="w-3 h-3" /> Ask AI
            </button>
          </div>

          <div className="flex-1 p-5 overflow-y-auto space-y-5 text-xs text-[#8E8EA0] leading-relaxed">
            {activeTab === 'description' && (
              <>
                {location.state?.fromRecommendation && (
                  <div className="p-3.5 rounded-xl bg-[#8B5CF6]/10 border border-[#8B5CF6]/30 flex items-start gap-2.5">
                    <Sparkles className="w-4 h-4 text-[#A78BFA] shrink-0 mt-0.5" />
                    <div className="space-y-1 min-w-0">
                      <div className="flex items-center gap-2">
                        <span className="text-[10px] font-mono uppercase tracking-wider text-[#A78BFA] font-bold">
                          Recommended by AI Engine
                        </span>
                        {location.state?.category && (
                          <span className="text-[9px] font-mono px-1.5 py-0.5 rounded bg-[#8B5CF6]/20 text-[#D8B4FE]">
                            {location.state.category.replace(/_/g, ' ')}
                          </span>
                        )}
                      </div>
                      <p className="text-xs text-[#E2E8F0]">
                        {location.state.primaryReason || 'Selected based on your learning curve and practice vector.'}
                      </p>
                    </div>
                  </div>
                )}
                <div className="flex items-start justify-between gap-3">
                  <h2 className="text-base font-bold text-[#F4F4F6]">{problem.title}</h2>
                  {problem.external_problem_url && (
                    <a
                      href={problem.external_problem_url}
                      target="_blank"
                      rel="noreferrer"
                      className="text-[11px] text-[#A78BFA] flex items-center gap-1 shrink-0"
                    >
                      <ExternalLink className="w-3 h-3" /> Official statement
                    </a>
                  )}
                </div>
                <RenderProblemText text={problem.description} />
                <div className="flex flex-wrap gap-2">
                  {(['Solved', 'Needs Revision', 'Attempted'] as ProgressStatus[]).map((s) => (
                    <button
                      key={s}
                      onClick={() => markStatus(s)}
                      className={`px-2 py-1 rounded-lg text-[10px] border ${
                        progressStatus === s ? 'border-[#8B5CF6] text-[#A78BFA]' : 'border-[#1C1C28] text-[#8E8EA0]'
                      }`}
                    >
                      <Bookmark className="w-3 h-3 inline mr-1" />
                      {s}
                    </button>
                  ))}
                </div>
                <div className="space-y-3">
                  <span className="text-[10px] font-bold text-[#5A5A6E] uppercase">Test Examples</span>
                  {problem.examples.map((ex, i) => (
                    <div key={i} className="p-3 bg-[#12121B] rounded-xl border border-[#1C1C28] space-y-1">
                      <div><span className="text-[#5A5A6E]">Input:</span> {ex.input}</div>
                      <div className="text-emerald-400"><span className="text-[#5A5A6E]">Output:</span> {ex.output}</div>
                    </div>
                  ))}
                </div>
                <ul className="list-disc list-inside space-y-1 text-[11px] text-[#A78BFA]">
                  {problem.constraints.map((c, i) => (
                    <li key={i}>{c}</li>
                  ))}
                </ul>
              </>
            )}

            {activeTab === 'hints' && (
              <div className="space-y-4">
                {unlockedHints.length === 0 && (
                  <div className="p-4 bg-[#12121B] rounded-xl border border-[#1C1C28] text-center">
                    Unlock Socratic hints one level at a time. No spoilers until you ask.
                  </div>
                )}
                {unlockedHints.map((h, i) => (
                  <div key={i} className="p-3.5 bg-[#12121B] rounded-xl border border-[#8B5CF6]/30 space-y-1.5">
                    <span className="text-xs font-bold text-[#E9D5FF]">Level {h.level}: {h.title}</span>
                    <p className="text-[#D1D1DE] font-sans whitespace-pre-wrap">{h.content}</p>
                    {h.guiding_question && <div className="text-[11px] text-[#A78BFA] italic">{h.guiding_question}</div>}
                  </div>
                ))}
                {unlockedHints.length < 5 && (
                  <Button onClick={handleUnlockNextHint} disabled={isLoadingHint} className="w-full text-xs py-2">
                    {isLoadingHint ? 'Generating…' : `Unlock Next Hint (Level ${unlockedHints.length + 1})`}
                  </Button>
                )}
              </div>
            )}

            {activeTab === 'history' && (
              <div className="space-y-2">
                {history.length === 0 && <p>No submissions yet for this problem.</p>}
                {history.map((s, i) => (
                  <div key={s.id || s.submission_id || i} className="p-3 bg-[#12121B] rounded-xl border border-[#1C1C28] flex justify-between">
                    <div>
                      <div className="text-[#F4F4F6]">{s.status}</div>
                      <div className="text-[10px] text-[#5A5A6E]">
                        {s.language} · {s.submitted_at ? new Date(s.submitted_at).toLocaleString() : ''}
                      </div>
                    </div>
                    {s.source_code && (
                      <button className="text-[#A78BFA] text-[11px]" onClick={() => setCode(s.source_code || '')}>
                        Load code
                      </button>
                    )}
                  </div>
                ))}
              </div>
            )}

            {activeTab === 'related' && (
              <div className="space-y-4">
                {similarRecs.length > 0 && (
                  <div className="space-y-2">
                    <div className="text-[11px] font-mono uppercase tracking-wider text-[#A78BFA] font-semibold flex items-center gap-1.5">
                      <Sparkles className="w-3.5 h-3.5 text-[#8B5CF6]" />
                      <span>AI Curated Similar Problems</span>
                    </div>
                    {similarRecs.map((sim) => (
                      <button
                        key={sim.id}
                        onClick={() =>
                          navigate(`/practice/${sim.problem_id}`, {
                            state: {
                              fromRecommendation: true,
                              recommendationId: sim.id,
                              category: sim.category,
                              primaryReason: sim.explanation,
                            }
                          })
                        }
                        className="w-full text-left p-3 bg-[#12121B] rounded-xl border border-[#1C1C28] hover:border-[#8B5CF6]/40 transition-colors space-y-1 group"
                      >
                        <div className="flex items-center justify-between">
                          <span className="text-xs font-semibold text-[#F4F4F6] group-hover:text-[#A78BFA] truncate">
                            {sim.problem_title || sim.headline}
                          </span>
                          <span className="text-[10px] font-mono text-purple-400 shrink-0">
                            {normalizePercentage(sim.match_score)}% match
                          </span>
                        </div>
                        <p className="text-[11px] text-[#8E8EA0] line-clamp-2">
                          {sim.explanation}
                        </p>
                        <div className="text-[10px] font-mono text-[#5A5A6E] flex items-center gap-2 pt-0.5">
                          <span>{sim.problem_difficulty || sim.recommended_difficulty}</span>
                          <span>·</span>
                          <span>{sim.primary_concept}</span>
                        </div>
                      </button>
                    ))}
                  </div>
                )}
                <div className="text-[11px] font-mono uppercase tracking-wider text-[#8E8EA0] font-semibold">
                  Subtopic Problems
                </div>
              <div className="space-y-2">
                {related.length === 0 && <p>No related problems from this subtopic.</p>}
                {related.map((r) => (
                  <button
                    key={r.id}
                    onClick={() => navigate(`/practice/${r.id}`)}
                    className="w-full text-left p-3 bg-[#12121B] rounded-xl border border-[#1C1C28] hover:border-[#8B5CF6]/40"
                  >
                    <div className="text-[#F4F4F6]">{r.title}</div>
                    <div className="text-[10px] text-[#5A5A6E]">{r.difficulty} · {r.progress_status}</div>
                  </button>
                ))}
              </div>
            </div>
            )}
          </div>
        </div>

        <div className="lg:col-span-7 flex flex-col gap-3 min-h-0">
          <div className="flex-1 bg-[#0D0D14] border border-[#1C1C28] rounded-xl overflow-hidden flex flex-col min-h-[300px]">
            <div className="h-8 px-4 bg-[#0A0A10] border-b border-[#1C1C28] flex items-center justify-between text-xs text-[#5A5A6E]">
              <div className="flex items-center gap-2">
                <FileCode className="w-3.5 h-3.5 text-[#8B5CF6]" />
                <span className="text-[#F4F4F6]">solution.{toJudgeLanguage(selectedLanguage)}</span>
              </div>
            </div>
            <div className="flex-1 w-full min-h-0">
              <Suspense fallback={<div className="h-full flex items-center justify-center text-[11px] text-[#5A5A6E]">Loading editor…</div>}>
                <Editor
                  height="100%"
                  language={toMonacoLanguage(selectedLanguage)}
                  theme="vs-dark"
                  value={code}
                  onChange={(newVal) => setCode(newVal || '')}
                  options={{
                    minimap: { enabled: false },
                    fontSize: 13,
                    scrollBeyondLastLine: false,
                    automaticLayout: true,
                    tabSize: 4,
                    fontFamily: 'JetBrains Mono, monospace',
                  }}
                />
              </Suspense>
            </div>
          </div>

          <div className="h-56 bg-[#0D0D14] border border-[#1C1C28] rounded-xl flex flex-col overflow-hidden shrink-0">
            <div className="h-8 px-4 bg-[#0A0A10] border-b border-[#1C1C28] flex items-center justify-between">
              <div className="flex items-center gap-2">
                <button onClick={() => setActiveBottomTab('testcases')} className={`px-2.5 py-1 rounded text-xs ${activeBottomTab === 'testcases' ? 'bg-[#181824] text-[#F4F4F6]' : 'text-[#8E8EA0]'}`}>
                  Test Cases
                </button>
                <button onClick={() => setActiveBottomTab('console')} className={`px-2.5 py-1 rounded text-xs flex items-center gap-1 ${activeBottomTab === 'console' ? 'bg-[#181824] text-[#F4F4F6]' : 'text-[#8E8EA0]'}`}>
                  <Terminal className="w-3 h-3 text-[#8B5CF6]" /> Output
                </button>
              </div>
              {executionResult && (
                <span className={executionResult.status === 'Accepted' ? 'text-emerald-400 text-xs' : 'text-amber-400 text-xs'}>
                  {executionResult.status}
                </span>
              )}
            </div>
            <div className="flex-1 p-3.5 overflow-y-auto text-xs">
              {activeBottomTab === 'testcases' ? (
                <div className="space-y-3">
                  <div className="flex items-center gap-2 flex-wrap">
                    {tests.map((tc, idx) => (
                      <button
                        key={tc.id}
                        onClick={() => setSelectedTestCaseIndex(idx)}
                        className={`px-2.5 py-1 rounded-md text-xs border ${
                          selectedTestCaseIndex === idx
                            ? 'bg-[#8B5CF6]/20 border-[#8B5CF6] text-white'
                            : 'bg-[#12121B] border-[#1C1C28] text-[#8E8EA0]'
                        }`}
                      >
                        Case {idx + 1}
                      </button>
                    ))}
                    {tests.length === 0 && <span className="text-[#5A5A6E]">No sample tests registered yet.</span>}
                  </div>
                  {tests[selectedTestCaseIndex] && (
                    <div className="p-3 bg-[#12121B] rounded-xl border border-[#1C1C28] space-y-2">
                      <div><span className="text-[#5A5A6E]">Input:</span> {tests[selectedTestCaseIndex].input || '(empty)'}</div>
                      <div className="text-emerald-400"><span className="text-[#5A5A6E]">Expected:</span> {tests[selectedTestCaseIndex].expectedOutput || '(empty)'}</div>
                    </div>
                  )}
                </div>
              ) : (
                <div className="space-y-2">
                  {executionResult ? (
                    <>
                      <div className="flex items-center gap-4 text-[#8E8EA0] pb-1 border-b border-[#1C1C28]">
                        <span>Runtime: {executionResult.runtimeMs}ms</span>
                        <span>Memory: {executionResult.memoryMb}MB</span>
                        <span>Passed: {executionResult.passedTests} / {executionResult.totalTests}</span>
                      </div>
                      {executionResult.results?.map((r) => (
                        <div key={r.test_case_order} className="p-2 rounded-lg bg-[#12121B] border border-[#1C1C28]">
                          Case {r.test_case_order} · {r.status} {r.is_hidden ? '(hidden)' : ''}
                          {!r.is_hidden && r.actual_output && <div className="text-[#D1D1DE] whitespace-pre-wrap">{r.actual_output}</div>}
                          {r.error_message && <div className="text-amber-400">{r.error_message}</div>}
                        </div>
                      ))}
                      {executionResult.stderr && <div className="text-amber-400 whitespace-pre-wrap">{executionResult.stderr}</div>}
                      {executionResult.stdout && <div className="text-[#F4F4F6] whitespace-pre-wrap">{executionResult.stdout}</div>}
                      {executionResult.diagnosis && (
                        <div className="mt-3 p-3.5 bg-[#17122B] border border-[#8B5CF6]/40 rounded-xl space-y-2">
                          <div className="flex items-center gap-2 text-[#E9D5FF]">
                            <Sparkles className="w-4 h-4" />
                            {executionResult.diagnosis.category} / {executionResult.diagnosis.subtype}
                          </div>
                          <p className="text-[#D1D1DE]">{executionResult.diagnosis.description}</p>
                          <div className="text-[11px] text-[#34D399]">{executionResult.diagnosis.suggested_correction}</div>
                          <div className="flex flex-wrap gap-2 pt-2">
                            {executionResult.submissionId && (
                              <button onClick={() => handleExplainMistake(executionResult.submissionId!)} className="px-2.5 py-1 text-[11px] rounded-lg bg-[#8B5CF6]/20 text-[#E9D5FF] border border-[#8B5CF6]/40">
                                Explain Root Cause
                              </button>
                            )}
                            <button onClick={handleDryRun} className="px-2.5 py-1 text-[11px] rounded-lg bg-[#8B5CF6]/20 text-[#E9D5FF] border border-[#8B5CF6]/40">Dry Run</button>
                            <button onClick={handleComplexityCheck} className="px-2.5 py-1 text-[11px] rounded-lg bg-[#8B5CF6]/20 text-[#E9D5FF] border border-[#8B5CF6]/40">Complexity</button>
                            {toJudgeLanguage(selectedLanguage) === 'python' && (
                              <button onClick={handlePedagogy} className="px-2.5 py-1 text-[11px] rounded-lg bg-[#8B5CF6]/20 text-[#E9D5FF] border border-[#8B5CF6]/40">Pedagogy</button>
                            )}
                          </div>
                        </div>
                      )}
                      {!executionResult.diagnosis && (
                        <div className="flex flex-wrap gap-2 pt-2">
                          <button onClick={handleDryRun} className="px-2.5 py-1 text-[11px] rounded-lg bg-[#8B5CF6]/20 text-[#E9D5FF] border border-[#8B5CF6]/40">Dry Run</button>
                          <button onClick={handleComplexityCheck} className="px-2.5 py-1 text-[11px] rounded-lg bg-[#8B5CF6]/20 text-[#E9D5FF] border border-[#8B5CF6]/40">Complexity</button>
                          {toJudgeLanguage(selectedLanguage) === 'python' && (
                            <button onClick={handlePedagogy} className="px-2.5 py-1 text-[11px] rounded-lg bg-[#8B5CF6]/20 text-[#E9D5FF] border border-[#8B5CF6]/40">Pedagogy</button>
                          )}
                        </div>
                      )}
                    </>
                  ) : (
                    <div className="text-[#5A5A6E] italic">Click Run or Submit to evaluate against the test suite.</div>
                  )}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      {aiModal.isOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-sm p-4">
          <div className="bg-[#0D0D14] border border-[#8B5CF6]/30 rounded-2xl max-w-2xl w-full max-h-[85vh] flex flex-col overflow-hidden">
            <div className="flex items-center justify-between px-5 py-3.5 border-b border-[#1C1C28] bg-[#12121B]">
              <h3 className="text-sm font-bold text-[#F4F4F6]">{aiModal.title}</h3>
              <button onClick={() => setAiModal((prev) => ({ ...prev, isOpen: false }))} className="text-[#8E8EA0] p-1">
                <X className="w-4 h-4" />
              </button>
            </div>
            <div className="p-5 overflow-y-auto space-y-4 text-xs">
              {aiModal.loading ? (
                <div className="py-12 flex flex-col items-center gap-3 text-[#A78BFA]">
                  <div className="w-6 h-6 border-2 border-[#8B5CF6] border-t-transparent rounded-full animate-spin" />
                  Synthesizing…
                </div>
              ) : (
                <>
                  {aiModal.type === 'explain' && aiModal.explain && (
                    <div className="space-y-3 font-sans">
                      <div className="p-3 bg-[#16122B] border border-[#8B5CF6]/30 rounded-xl">
                        <span className="text-[11px] font-bold text-[#A78BFA] font-mono uppercase">What Went Wrong</span>
                        <p className="text-[#E9D5FF]">{aiModal.explain.what_went_wrong}</p>
                      </div>
                      <div className="p-3 bg-[#12121B] border border-[#1C1C28] rounded-xl">
                        <span className="text-[11px] font-bold text-amber-400 font-mono uppercase">Why</span>
                        <p className="text-[#D1D1DE]">{aiModal.explain.why_it_happened}</p>
                      </div>
                      <div className="p-3 bg-[#0E1F1A] border border-emerald-500/30 rounded-xl">
                        <span className="text-[11px] font-bold text-emerald-400 font-mono uppercase">How to Fix</span>
                        <p className="text-[#A7F3D0]">{aiModal.explain.how_to_fix}</p>
                      </div>
                    </div>
                  )}
                  {aiModal.type === 'dryrun' && aiModal.dryRun && (
                    <div className="space-y-2">
                      <div className="p-2.5 bg-[#12121B] border border-[#1C1C28] rounded-xl">Input: {aiModal.dryRun.test_case_input}</div>
                      {aiModal.dryRun.steps?.map((st) => (
                        <div key={st.step_number} className="p-3 bg-[#12121B] border border-[#1C1C28] rounded-xl">
                          <div className="text-[#A78BFA]">Step #{st.step_number} (Line {st.line_number}) · {st.action}</div>
                          <div className="text-[#8E8EA0]">{st.explanation}</div>
                        </div>
                      ))}
                      <div className="p-2.5 bg-[#0E1F1A] border border-emerald-500/30 rounded-xl text-emerald-300">Final: {aiModal.dryRun.final_output}</div>
                    </div>
                  )}
                  {aiModal.type === 'complexity' && aiModal.complexity && (
                    <div className="grid grid-cols-2 gap-3">
                      <div className="p-3 bg-[#12121B] border border-[#1C1C28] rounded-xl">
                        <div className="text-[#A78BFA] font-bold">{aiModal.complexity.time_complexity}</div>
                        <p className="text-[#8E8EA0]">{aiModal.complexity.time_explanation}</p>
                      </div>
                      <div className="p-3 bg-[#12121B] border border-[#1C1C28] rounded-xl">
                        <div className="text-[#34D399] font-bold">{aiModal.complexity.space_complexity}</div>
                        <p className="text-[#8E8EA0]">{aiModal.complexity.space_explanation}</p>
                      </div>
                    </div>
                  )}
                  {aiModal.type === 'pedagogy' && pedagogy && (
                    <div className="space-y-3">
                      <div className="p-3 bg-[#12121B] border border-[#1C1C28] rounded-xl">
                        Runner: {pedagogy.runner_result?.passed ? 'passed' : pedagogy.runner_result?.error_type || 'failed'}
                        {pedagogy.runner_result?.output && <pre className="mt-2 whitespace-pre-wrap">{pedagogy.runner_result.output}</pre>}
                      </div>
                      <div className="p-3 bg-[#16122B] border border-[#8B5CF6]/30 rounded-xl">
                        {pedagogy.diagnosis?.misconception} ({normalizePercentage(pedagogy.diagnosis?.confidence)}%)
                        <p className="text-[#8E8EA0] mt-1">{pedagogy.diagnosis?.evidence}</p>
                      </div>
                      {pedagogy.hint?.hint && (
                        <div className="p-3 bg-[#0E1F1A] border border-emerald-500/30 rounded-xl text-emerald-200">
                          {pedagogy.hint.hint_level}: {pedagogy.hint.hint}
                        </div>
                      )}
                    </div>
                  )}
                </>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default PracticePage;
