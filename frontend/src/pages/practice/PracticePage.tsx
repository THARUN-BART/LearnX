import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import Editor from '@monaco-editor/react';
import {
  Play,
  Send,
  RotateCcw,
  Copy,
  Maximize2,
  Minimize2,
  Terminal,
  CheckCircle2,
  AlertCircle,
  HelpCircle,
  FileCode,
  Sparkles,
  Bot,
  Layers,
  ChevronDown,
  Cpu,
  X
} from 'lucide-react';
import { SAMPLE_PROBLEMS, PROGRAMMING_LANGUAGES } from '../../utils/constants';
import {
  Problem,
  CodeExecutionResult,
  ProgressiveHint,
  DryRunResult,
  ComplexityResult,
  ExplainMistakeResult
} from '../../types/learning.types';
import codingService from '../../services/coding.service';
import learningService from '../../services/learning.service';
import aiService from '../../services/ai.service';
import Button from '../../components/ui/Button';
import Badge from '../../components/ui/Badge';
import { toast } from 'sonner';

export const PracticePage: React.FC = () => {
  const { problemId } = useParams<{ problemId: string }>();
  const navigate = useNavigate();

  // Selected problem
  const [problem, setProblem] = useState<Problem>(() => {
    return SAMPLE_PROBLEMS.find((p) => p.id === problemId || p.slug === problemId) || SAMPLE_PROBLEMS[0];
  });

  // Editor states
  const [selectedLanguage, setSelectedLanguage] = useState<string>('python');
  const [code, setCode] = useState<string>(() => problem.starterCode?.['python'] || '');
  const [activeTab, setActiveTab] = useState<'description' | 'hints'>('description');
  const [activeBottomTab, setActiveBottomTab] = useState<'testcases' | 'console'>('testcases');
  const [selectedTestCaseIndex, setSelectedTestCaseIndex] = useState<number>(0);
  const [isFullscreen, setIsFullscreen] = useState<boolean>(false);

  // Execution states
  const [isRunning, setIsRunning] = useState<boolean>(false);
  const [executionResult, setExecutionResult] = useState<CodeExecutionResult | null>(null);
  const [allProblems, setAllProblems] = useState<Problem[]>([]);

  // AI Tutor states
  const [unlockedHints, setUnlockedHints] = useState<ProgressiveHint[]>([]);
  const [isLoadingHint, setIsLoadingHint] = useState<boolean>(false);
  const [aiModal, setAiModal] = useState<{
    isOpen: boolean;
    title: string;
    type: 'explain' | 'dryrun' | 'complexity';
    loading: boolean;
    explain?: ExplainMistakeResult;
    dryRun?: DryRunResult;
    complexity?: ComplexityResult;
  }>({
    isOpen: false,
    title: '',
    type: 'explain',
    loading: false
  });

  // Fetch real problems list from API for selector
  useEffect(() => {
    const fetchAllProblems = async () => {
      const list = await learningService.getProblems();
      if (Array.isArray(list) && list.length > 0) {
        setAllProblems(list);
        if (!problemId) {
          setProblem(list[0]);
          setCode(list[0].starterCode?.[selectedLanguage.toLowerCase()] || list[0].starterCode?.['python'] || '');
        }
      }
    };
    fetchAllProblems();
  }, [problemId]);

  // Fetch problem from API when problemId changes
  useEffect(() => {
    let isMounted = true;
    const fetchProblem = async () => {
      if (!problemId) return;
      const prob = await learningService.getProblemById(problemId);
      if (prob && isMounted) {
        setProblem(prob);
        const starter =
          prob.starterCode?.[selectedLanguage.toLowerCase()] ||
          prob.starterCode?.['python'] ||
          '';
        setCode(starter);
        setExecutionResult(null);

        // Fetch unlocked progressive hints
        try {
          const hints = await aiService.getUnlockedHints(prob.id);
          if (Array.isArray(hints) && hints.length > 0) {
            setUnlockedHints(hints);
          } else {
            setUnlockedHints([]);
          }
        } catch {
          setUnlockedHints([]);
        }
      }
    };
    fetchProblem();
    return () => {
      isMounted = false;
    };
  }, [problemId]);

  const handleLanguageChange = (lang: string) => {
    setSelectedLanguage(lang);
    const starter =
      problem.starterCode?.[lang.toLowerCase()] ||
      problem.starterCode?.['python'] ||
      '';
    setCode(starter);
  };

  const handleResetCode = () => {
    const starter =
      problem.starterCode?.[selectedLanguage.toLowerCase()] ||
      problem.starterCode?.['python'] ||
      '';
    setCode(starter);
    toast.info('Starter code reset');
  };

  const handleCopyCode = () => {
    navigator.clipboard.writeText(code);
    toast.success('Code copied to clipboard');
  };

  const handleUnlockNextHint = async () => {
    const nextLevel = unlockedHints.length + 1;
    if (nextLevel > 5) {
      toast.info('All 5 progressive hint levels are already unlocked.');
      return;
    }
    setIsLoadingHint(true);
    try {
      const hint = await aiService.getProgressiveHint(
        problem.id,
        nextLevel,
        executionResult?.submissionId
      );
      setUnlockedHints((prev) => [...prev, hint]);
      toast.success(`Unlocked Hint Level ${nextLevel}: ${hint.title}`);
    } catch {
      toast.error('Failed to unlock hint from AI Tutor');
    } finally {
      setIsLoadingHint(false);
    }
  };

  const handleExplainMistake = async (submissionId: string) => {
    setAiModal({
      isOpen: true,
      title: 'AI Root Cause Diagnosis',
      type: 'explain',
      loading: true
    });
    try {
      const data = await aiService.explainMistake(submissionId);
      setAiModal((prev) => ({ ...prev, loading: false, explain: data }));
    } catch {
      toast.error('Could not retrieve root cause explanation');
      setAiModal((prev) => ({ ...prev, loading: false }));
    }
  };

  const handleDryRun = async () => {
    setAiModal({
      isOpen: true,
      title: 'Execution Trace & Dry Run',
      type: 'dryrun',
      loading: true
    });
    try {
      const data = await aiService.getDryRun(problem.id, code, selectedLanguage);
      setAiModal((prev) => ({ ...prev, loading: false, dryRun: data }));
    } catch {
      toast.error('Could not generate dry run trace');
      setAiModal((prev) => ({ ...prev, loading: false }));
    }
  };

  const handleComplexityCheck = async () => {
    setAiModal({
      isOpen: true,
      title: 'Time & Space Complexity Audit',
      type: 'complexity',
      loading: true
    });
    try {
      const data = await aiService.getComplexity(problem.id, code, selectedLanguage);
      setAiModal((prev) => ({ ...prev, loading: false, complexity: data }));
    } catch {
      toast.error('Could not generate complexity audit');
      setAiModal((prev) => ({ ...prev, loading: false }));
    }
  };

  const handleRunCode = async () => {
    setIsRunning(true);
    setActiveBottomTab('console');
    try {
      const result = await codingService.runCode({
        problemId: problem.id,
        language: selectedLanguage,
        code
      });
      setExecutionResult(result);
      if (result.status === 'Accepted') {
        toast.success('Test run passed!');
      } else {
        toast.error(`Execution status: ${result.status}`);
      }
    } catch {
      toast.error('Execution failed');
    } finally {
      setIsRunning(false);
    }
  };

  const handleSubmitCode = async () => {
    setIsRunning(true);
    setActiveBottomTab('console');
    try {
      const result = await codingService.submitCode({
        problemId: problem.id,
        language: selectedLanguage,
        code
      });
      setExecutionResult(result);
      if (result.status === 'Accepted') {
        toast.success('Solution Accepted! Great job!');
      } else {
        toast.error(`Submission status: ${result.status}`);
      }
    } catch {
      toast.error('Submission failed');
    } finally {
      setIsRunning(false);
    }
  };

  const monacoLanguageMap: Record<string, string> = {
    python: 'python',
    cpp: 'cpp',
    java: 'java',
    c: 'c',
    go: 'go'
  };

  return (
    <div
      className={`flex flex-col gap-3 font-mono ${
        isFullscreen
          ? 'fixed inset-0 z-50 bg-[#07070B] p-4 h-screen'
          : 'h-[calc(100vh-6.5rem)]'
      }`}
    >
      {/* Top Workspace Header & Toolstrip */}
      <div className="flex flex-wrap items-center justify-between gap-3 bg-[#0D0D14] border border-[#1C1C28] px-4 py-2.5 rounded-xl shrink-0">
        <div className="flex items-center gap-3">
          <div className="flex items-center gap-1.5 text-xs text-[#5A5A6E]">
            <Cpu className="w-3.5 h-3.5 text-[#8B5CF6]" />
            <span className="hidden sm:inline">LOC:</span>
          </div>

          <select
            value={problem.id}
            onChange={(e) => navigate(`/practice/${e.target.value}`)}
            className="bg-[#12121B] text-[#F4F4F6] text-xs font-mono font-semibold px-3 py-1.5 rounded-lg border border-[#1C1C28] outline-none cursor-pointer max-w-[220px] truncate"
          >
            {(allProblems.length > 0 ? allProblems : [problem]).map((p) => (
              <option key={p.id} value={p.id}>
                {p.title} ({p.difficulty})
              </option>
            ))}
          </select>

          <span className={`text-[10px] px-2 py-0.5 rounded font-mono ${
            problem.difficulty === 'Easy'
              ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20'
              : 'bg-amber-500/10 text-amber-400 border border-amber-500/20'
          }`}>
            {problem.difficulty}
          </span>
        </div>

        {/* Action Controls */}
        <div className="flex items-center gap-2">
          {/* Language selector */}
          <select
            value={selectedLanguage}
            onChange={(e) => handleLanguageChange(e.target.value)}
            className="bg-[#12121B] text-[#F4F4F6] text-xs font-mono px-3 py-1.5 rounded-lg border border-[#1C1C28] outline-none cursor-pointer"
          >
            {PROGRAMMING_LANGUAGES.map((l) => (
              <option key={l.id} value={l.id.toLowerCase()}>
                {l.name}
              </option>
            ))}
          </select>

          <button
            onClick={handleResetCode}
            title="Reset code boilerplate"
            className="p-1.5 text-[#8E8EA0] hover:text-[#F4F4F6] hover:bg-[#14141E] rounded-lg transition-colors border border-[#1C1C28] cursor-pointer"
          >
            <RotateCcw className="w-3.5 h-3.5" />
          </button>

          <button
            onClick={handleCopyCode}
            title="Copy code"
            className="p-1.5 text-[#8E8EA0] hover:text-[#F4F4F6] hover:bg-[#14141E] rounded-lg transition-colors border border-[#1C1C28] cursor-pointer"
          >
            <Copy className="w-3.5 h-3.5" />
          </button>

          <button
            onClick={() => setIsFullscreen(!isFullscreen)}
            title="Toggle fullscreen"
            className="p-1.5 text-[#8E8EA0] hover:text-[#F4F4F6] hover:bg-[#14141E] rounded-lg transition-colors border border-[#1C1C28] cursor-pointer"
          >
            {isFullscreen ? <Minimize2 className="w-3.5 h-3.5" /> : <Maximize2 className="w-3.5 h-3.5" />}
          </button>

          <Button
            variant="secondary"
            size="sm"
            onClick={handleRunCode}
            isLoading={isRunning}
            leftIcon={<Play className="w-3 h-3" />}
            className="text-xs font-mono h-8"
          >
            Run
          </Button>

          <Button
            variant="primary"
            size="sm"
            onClick={handleSubmitCode}
            isLoading={isRunning}
            leftIcon={<Send className="w-3 h-3" />}
            className="text-xs font-mono h-8"
          >
            Submit
          </Button>
        </div>
      </div>

      {/* Main Multi-Panel Workspace: LEFT (Problem Details), RIGHT (Editor + Console) */}
      <div className="flex-1 grid grid-cols-1 lg:grid-cols-12 gap-3 min-h-0">
        {/* LEFT PANEL: Problem Details, Examples, Constraints, Hints */}
        <div className="lg:col-span-5 bg-[#0D0D14] border border-[#1C1C28] rounded-xl flex flex-col overflow-hidden">
          <div className="flex items-center justify-between border-b border-[#1C1C28] px-4 py-2 bg-[#0A0A10]">
            <div className="flex items-center gap-2">
              <button
                onClick={() => setActiveTab('description')}
                className={`px-3 py-1 rounded-lg text-xs font-mono transition-colors cursor-pointer ${
                  activeTab === 'description'
                    ? 'bg-[#181824] text-[#A78BFA] border border-[#8B5CF6]/30'
                    : 'text-[#8E8EA0] hover:text-white'
                }`}
              >
                Problem Specs
              </button>
              <button
                onClick={() => setActiveTab('hints')}
                className={`px-3 py-1 rounded-lg text-xs font-mono transition-colors cursor-pointer ${
                  activeTab === 'hints'
                    ? 'bg-[#181824] text-[#A78BFA] border border-[#8B5CF6]/30'
                    : 'text-[#8E8EA0] hover:text-white'
                }`}
              >
                Hints ({problem.hints.length})
              </button>
            </div>

            <button
              onClick={() => navigate('/assistant')}
              className="text-[11px] text-[#8B5CF6] hover:text-[#A78BFA] flex items-center gap-1 font-mono"
            >
              <Bot className="w-3 h-3" />
              <span>Ask AI</span>
            </button>
          </div>

          <div className="flex-1 p-5 overflow-y-auto space-y-6 text-xs text-[#8E8EA0] leading-relaxed">
            {activeTab === 'description' ? (
              <>
                <div>
                  <h2 className="text-base font-bold text-[#F4F4F6] font-mono mb-2">
                    {problem.title}
                  </h2>
                  <div className="whitespace-pre-line text-[#D1D1DE] font-sans">
                    {problem.description}
                  </div>
                </div>

                {/* Examples */}
                <div className="space-y-3">
                  <span className="text-[10px] font-bold text-[#5A5A6E] uppercase tracking-wider font-mono">
                    Test Examples
                  </span>
                  {problem.examples.map((ex, i) => (
                    <div
                      key={i}
                      className="p-3 bg-[#12121B] rounded-xl border border-[#1C1C28] space-y-1 font-mono text-xs"
                    >
                      <div className="text-[#8E8EA0]">
                        <span className="text-[#5A5A6E]">Input:</span> {ex.input}
                      </div>
                      <div className="text-emerald-400">
                        <span className="text-[#5A5A6E]">Output:</span> {ex.output}
                      </div>
                      {ex.explanation && (
                        <div className="text-[#5A5A6E] text-[11px] pt-0.5">
                          Note: {ex.explanation}
                        </div>
                      )}
                    </div>
                  ))}
                </div>

                {/* Constraints */}
                <div className="space-y-2">
                  <span className="text-[10px] font-bold text-[#5A5A6E] uppercase tracking-wider font-mono">
                    Complexity Constraints
                  </span>
                  <ul className="list-disc list-inside space-y-1 font-mono text-[11px] text-[#A78BFA]">
                    {problem.constraints.map((c, i) => (
                      <li key={i}>{c}</li>
                    ))}
                  </ul>
                </div>
              </>
            ) : (
              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <div className="text-xs text-[#A78BFA] font-mono flex items-center gap-1.5">
                    <Sparkles className="w-4 h-4 text-[#8B5CF6]" />
                    <span>Progressive AI Scaffolding</span>
                  </div>
                  <span className="text-[10px] text-[#5A5A6E] font-mono">
                    {unlockedHints.length} / 5 Levels
                  </span>
                </div>

                {/* Unlocked progressive hints */}
                {unlockedHints.length > 0 ? (
                  unlockedHints.map((h, i) => (
                    <div
                      key={i}
                      className="p-3.5 bg-[#12121B] rounded-xl border border-[#8B5CF6]/30 space-y-1.5"
                    >
                      <div className="flex items-center justify-between">
                        <span className="text-xs font-bold text-[#E9D5FF] font-mono flex items-center gap-1.5">
                          <span className="w-2 h-2 rounded-full bg-[#8B5CF6]" />
                          Level {h.level}: {h.title}
                        </span>
                        <span className="text-[10px] text-[#A78BFA] bg-[#8B5CF6]/10 px-2 py-0.5 rounded font-mono">
                          {h.level === 1 && 'Concept'}
                          {h.level === 2 && 'Approach'}
                          {h.level === 3 && 'Edge Cases'}
                          {h.level === 4 && 'Blueprint'}
                          {h.level === 5 && 'Full Code Fix'}
                        </span>
                      </div>
                      <p className="text-xs text-[#D1D1DE] pt-1 font-sans whitespace-pre-wrap leading-relaxed">
                        {h.content}
                      </p>
                      {h.guiding_question && (
                        <div className="text-[11px] text-[#A78BFA] italic pt-1 border-t border-[#1C1C28]">
                          🤔 {h.guiding_question}
                        </div>
                      )}
                    </div>
                  ))
                ) : (
                  <div className="p-3 bg-[#12121B] rounded-xl border border-[#1C1C28] text-xs text-[#8E8EA0] text-center py-4">
                    No hints unlocked yet. LearnX uses Socratic progressive hints so you learn to derive the algorithm yourself without spoilers!
                  </div>
                )}

                {/* Unlock next level button */}
                {unlockedHints.length < 5 && (
                  <Button
                    onClick={handleUnlockNextHint}
                    disabled={isLoadingHint}
                    className="w-full text-xs font-mono py-2 bg-[#8B5CF6]/10 hover:bg-[#8B5CF6]/20 border border-[#8B5CF6]/40 text-[#E9D5FF]"
                  >
                    <Sparkles className="w-3.5 h-3.5 mr-1.5 text-[#A78BFA]" />
                    {isLoadingHint
                      ? 'Generating Hint...'
                      : `Unlock Next Hint (Level ${unlockedHints.length + 1})`}
                  </Button>
                )}
              </div>
            )}
          </div>
        </div>

        {/* RIGHT PANEL: Center Monaco Editor + Bottom Console/Test Cases */}
        <div className="lg:col-span-7 flex flex-col gap-3 min-h-0">
          {/* Monaco Editor Container */}
          <div className="flex-1 bg-[#0D0D14] border border-[#1C1C28] rounded-xl overflow-hidden flex flex-col min-h-[300px]">
            {/* Editor tab strip */}
            <div className="h-8 px-4 bg-[#0A0A10] border-b border-[#1C1C28] flex items-center justify-between text-xs font-mono text-[#5A5A6E]">
              <div className="flex items-center gap-2">
                <FileCode className="w-3.5 h-3.5 text-[#8B5CF6]" />
                <span className="text-[#F4F4F6] text-xs">solution.{selectedLanguage.toLowerCase()}</span>
              </div>
              <span className="text-[10px]">UTF-8 • Tabs: 4</span>
            </div>

            <div className="flex-1 w-full min-h-0">
              <Editor
                height="100%"
                language={monacoLanguageMap[selectedLanguage.toLowerCase()] || 'python'}
                theme="vs-dark"
                value={code}
                onChange={(newVal) => setCode(newVal || '')}
                options={{
                  minimap: { enabled: false },
                  fontSize: 13,
                  lineNumbers: 'on',
                  scrollBeyondLastLine: false,
                  automaticLayout: true,
                  tabSize: 4,
                  fontFamily: 'JetBrains Mono, monospace'
                }}
              />
            </div>
          </div>

          {/* BOTTOM PANEL: Test Cases & Console Output */}
          <div className="h-52 bg-[#0D0D14] border border-[#1C1C28] rounded-xl flex flex-col overflow-hidden shrink-0">
            <div className="h-8 px-4 bg-[#0A0A10] border-b border-[#1C1C28] flex items-center justify-between">
              <div className="flex items-center gap-2">
                <button
                  onClick={() => setActiveBottomTab('testcases')}
                  className={`px-2.5 py-1 rounded text-xs font-mono transition-colors cursor-pointer ${
                    activeBottomTab === 'testcases'
                      ? 'bg-[#181824] text-[#F4F4F6]'
                      : 'text-[#8E8EA0] hover:text-white'
                  }`}
                >
                  Test Cases
                </button>
                <button
                  onClick={() => setActiveBottomTab('console')}
                  className={`px-2.5 py-1 rounded text-xs font-mono transition-colors flex items-center gap-1.5 cursor-pointer ${
                    activeBottomTab === 'console'
                      ? 'bg-[#181824] text-[#F4F4F6]'
                      : 'text-[#8E8EA0] hover:text-white'
                  }`}
                >
                  <Terminal className="w-3 h-3 text-[#8B5CF6]" />
                  <span>Execution Output</span>
                </button>
              </div>

              {executionResult && (
                <div className="flex items-center gap-2 text-xs font-mono">
                  <span
                    className={
                      executionResult.status === 'Accepted'
                        ? 'text-emerald-400 font-semibold'
                        : 'text-amber-400 font-semibold'
                    }
                  >
                    {executionResult.status}
                  </span>
                </div>
              )}
            </div>

            <div className="flex-1 p-3.5 overflow-y-auto text-xs font-mono">
              {activeBottomTab === 'testcases' ? (
                <div className="space-y-3">
                  <div className="flex items-center gap-2">
                    {problem.testCases.map((tc, idx) => (
                      <button
                        key={tc.id}
                        onClick={() => setSelectedTestCaseIndex(idx)}
                        className={`px-2.5 py-1 rounded-md text-xs font-mono border transition-all cursor-pointer ${
                          selectedTestCaseIndex === idx
                            ? 'bg-[#8B5CF6]/20 border-[#8B5CF6] text-white'
                            : 'bg-[#12121B] border-[#1C1C28] text-[#8E8EA0] hover:border-[#2D2D3E]'
                        }`}
                      >
                        Case {idx + 1}
                      </button>
                    ))}
                  </div>

                  {problem.testCases[selectedTestCaseIndex] && (
                    <div className="p-3 bg-[#12121B] rounded-xl border border-[#1C1C28] space-y-2">
                      <div className="text-[#8E8EA0]">
                        <span className="text-[#5A5A6E]">Input:</span>{' '}
                        {problem.testCases[selectedTestCaseIndex].input}
                      </div>
                      <div className="text-emerald-400">
                        <span className="text-[#5A5A6E]">Expected Output:</span>{' '}
                        {problem.testCases[selectedTestCaseIndex].expectedOutput}
                      </div>
                    </div>
                  )}
                </div>
              ) : (
                <div className="space-y-2">
                  {executionResult ? (
                    <div className="space-y-2">
                      <div className="flex items-center gap-4 text-[#8E8EA0] pb-1 border-b border-[#1C1C28]">
                        <span>Runtime: {executionResult.runtimeMs}ms</span>
                        <span>Memory: {executionResult.memoryMb}MB</span>
                        <span>
                          Passed: {executionResult.passedTests} / {executionResult.totalTests}
                        </span>
                      </div>

                      {executionResult.stderr && (
                        <div className="text-amber-400 whitespace-pre-wrap">
                          {executionResult.stderr}
                        </div>
                      )}

                      {executionResult.stdout && (
                        <div className="text-[#F4F4F6] whitespace-pre-wrap">
                          {executionResult.stdout}
                        </div>
                      )}

                      {/* Intelligent Mistake Diagnosis from Module A */}
                      {executionResult.diagnosis && (
                        <div className="mt-3 p-3.5 bg-[#17122B] border border-[#8B5CF6]/40 rounded-xl space-y-2">
                          <div className="flex items-center justify-between">
                            <div className="flex items-center gap-2">
                              <Sparkles className="w-4 h-4 text-[#A78BFA]" />
                              <span className="font-bold text-xs text-[#E9D5FF] uppercase tracking-wider">
                                Intelligent Diagnosis: {executionResult.diagnosis.category} / {executionResult.diagnosis.subtype}
                              </span>
                            </div>
                            <span className="text-[10px] px-2 py-0.5 rounded-full bg-[#8B5CF6]/20 text-[#A78BFA] border border-[#8B5CF6]/30">
                              {executionResult.diagnosis.confidence_level} Confidence
                            </span>
                          </div>
                          <p className="text-xs text-[#D1D1DE] leading-relaxed">
                            {executionResult.diagnosis.description}
                          </p>
                          <div className="text-[11px] text-[#A78BFA]">
                            <span className="text-[#8E8EA0]">DSA Concept:</span> {executionResult.diagnosis.related_concept}
                          </div>
                          <div className="text-[11px] text-[#34D399]">
                            <span className="text-[#8E8EA0]">Suggested Fix:</span> {executionResult.diagnosis.suggested_correction}
                          </div>

                          {/* Quick AI Tutor Tools */}
                          <div className="flex flex-wrap gap-2 pt-2 border-t border-[#8B5CF6]/20">
                            <button
                              onClick={() => handleExplainMistake(executionResult.submissionId || executionResult.diagnosis!.submission_id)}
                              className="px-2.5 py-1 text-[11px] rounded-lg bg-[#8B5CF6]/20 hover:bg-[#8B5CF6]/30 text-[#E9D5FF] border border-[#8B5CF6]/40 flex items-center gap-1.5 cursor-pointer transition-all"
                            >
                              <Bot className="w-3.5 h-3.5 text-[#A78BFA]" /> Explain Root Cause
                            </button>
                            <button
                              onClick={handleDryRun}
                              className="px-2.5 py-1 text-[11px] rounded-lg bg-[#8B5CF6]/20 hover:bg-[#8B5CF6]/30 text-[#E9D5FF] border border-[#8B5CF6]/40 flex items-center gap-1.5 cursor-pointer transition-all"
                            >
                              <Play className="w-3.5 h-3.5 text-[#A78BFA]" /> Simulate Dry Run
                            </button>
                            <button
                              onClick={handleComplexityCheck}
                              className="px-2.5 py-1 text-[11px] rounded-lg bg-[#8B5CF6]/20 hover:bg-[#8B5CF6]/30 text-[#E9D5FF] border border-[#8B5CF6]/40 flex items-center gap-1.5 cursor-pointer transition-all"
                            >
                              <Cpu className="w-3.5 h-3.5 text-[#A78BFA]" /> Complexity Audit
                            </button>
                          </div>
                        </div>
                      )}
                    </div>
                  ) : (
                    <div className="text-[#5A5A6E] italic">
                      Click "Run" or "Submit" to compile and execute your solution against test suites.
                    </div>
                  )}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* AI Tutor Intelligence Modal */}
      {aiModal.isOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-sm p-4">
          <div className="bg-[#0D0D14] border border-[#8B5CF6]/30 rounded-2xl max-w-2xl w-full max-h-[85vh] flex flex-col overflow-hidden shadow-2xl font-mono">
            <div className="flex items-center justify-between px-5 py-3.5 border-b border-[#1C1C28] bg-[#12121B]">
              <div className="flex items-center gap-2">
                <Sparkles className="w-4 h-4 text-[#8B5CF6]" />
                <h3 className="text-sm font-bold text-[#F4F4F6] font-mono">{aiModal.title}</h3>
              </div>
              <button
                onClick={() => setAiModal((prev) => ({ ...prev, isOpen: false }))}
                className="text-[#8E8EA0] hover:text-white p-1 rounded-lg cursor-pointer"
              >
                <X className="w-4 h-4" />
              </button>
            </div>

            <div className="p-5 overflow-y-auto space-y-4 text-xs font-mono leading-relaxed">
              {aiModal.loading ? (
                <div className="py-12 flex flex-col items-center justify-center gap-3 text-[#A78BFA]">
                  <div className="w-6 h-6 border-2 border-[#8B5CF6] border-t-transparent rounded-full animate-spin" />
                  <span>Synthesizing pedagogical breakdown...</span>
                </div>
              ) : (
                <>
                  {aiModal.type === 'explain' && aiModal.explain && (
                    <div className="space-y-3 font-sans">
                      <div className="p-3 bg-[#16122B] border border-[#8B5CF6]/30 rounded-xl space-y-1">
                        <span className="text-[11px] font-bold text-[#A78BFA] font-mono uppercase">What Went Wrong</span>
                        <p className="text-xs text-[#E9D5FF] leading-relaxed">{aiModal.explain.what_went_wrong}</p>
                      </div>
                      <div className="p-3 bg-[#12121B] border border-[#1C1C28] rounded-xl space-y-1">
                        <span className="text-[11px] font-bold text-amber-400 font-mono uppercase">Why It Happened (Conceptual Cause)</span>
                        <p className="text-xs text-[#D1D1DE] leading-relaxed">{aiModal.explain.why_it_happened}</p>
                      </div>
                      <div className="p-3 bg-[#0E1F1A] border border-emerald-500/30 rounded-xl space-y-1">
                        <span className="text-[11px] font-bold text-emerald-400 font-mono uppercase">How to Resolve</span>
                        <p className="text-xs text-[#A7F3D0] leading-relaxed">{aiModal.explain.how_to_fix}</p>
                      </div>
                      {aiModal.explain.related_dsa_concepts && aiModal.explain.related_dsa_concepts.length > 0 && (
                        <div className="flex items-center gap-2 pt-2">
                          <span className="text-[#5A5A6E] text-[11px] font-mono">Related Concepts:</span>
                          {aiModal.explain.related_dsa_concepts.map((c, i) => (
                            <span key={i} className="px-2 py-0.5 rounded bg-[#8B5CF6]/10 text-[#A78BFA] text-[10px] font-mono border border-[#8B5CF6]/20">
                              {c}
                            </span>
                          ))}
                        </div>
                      )}
                    </div>
                  )}

                  {aiModal.type === 'dryrun' && aiModal.dryRun && (
                    <div className="space-y-3">
                      <div className="p-2.5 bg-[#12121B] border border-[#1C1C28] rounded-xl text-[#8E8EA0]">
                        <span className="text-[#5A5A6E]">Input Trace:</span> {aiModal.dryRun.test_case_input}
                      </div>
                      <div className="space-y-2">
                        {aiModal.dryRun.steps.map((st) => (
                          <div key={st.step_number} className="p-3 bg-[#12121B] border border-[#1C1C28] rounded-xl space-y-1">
                            <div className="flex items-center justify-between text-[11px] text-[#A78BFA]">
                              <span>Step #{st.step_number} (Line {st.line_number})</span>
                              <span>{st.action}</span>
                            </div>
                            <div className="text-[#8E8EA0] text-[11px]">{st.explanation}</div>
                            {st.variables && Object.keys(st.variables).length > 0 && (
                              <div className="text-[#34D399] text-[11px] bg-[#0A0A10] p-1.5 rounded border border-[#1C1C28]">
                                State: {JSON.stringify(st.variables)}
                              </div>
                            )}
                          </div>
                        ))}
                      </div>
                      <div className="p-2.5 bg-[#0E1F1A] border border-emerald-500/30 rounded-xl text-emerald-300">
                        Final Result: {aiModal.dryRun.final_output}
                      </div>
                    </div>
                  )}

                  {aiModal.type === 'complexity' && aiModal.complexity && (
                    <div className="space-y-3 font-sans">
                      <div className="grid grid-cols-2 gap-3">
                        <div className="p-3 bg-[#12121B] border border-[#1C1C28] rounded-xl space-y-1">
                          <span className="text-[10px] font-mono text-[#5A5A6E] uppercase">Time Complexity</span>
                          <div className="text-base font-bold text-[#A78BFA] font-mono">{aiModal.complexity.time_complexity}</div>
                          <p className="text-[11px] text-[#8E8EA0]">{aiModal.complexity.time_explanation}</p>
                        </div>
                        <div className="p-3 bg-[#12121B] border border-[#1C1C28] rounded-xl space-y-1">
                          <span className="text-[10px] font-mono text-[#5A5A6E] uppercase">Space Complexity</span>
                          <div className="text-base font-bold text-[#34D399] font-mono">{aiModal.complexity.space_complexity}</div>
                          <p className="text-[11px] text-[#8E8EA0]">{aiModal.complexity.space_explanation}</p>
                        </div>
                      </div>
                      {aiModal.complexity.optimisation_notes && (
                        <div className="p-3 bg-[#16122B] border border-[#8B5CF6]/30 rounded-xl space-y-1">
                          <span className="text-[11px] font-bold text-[#E9D5FF] font-mono">Optimisation Notes</span>
                          <p className="text-xs text-[#D1D1DE]">{aiModal.complexity.optimisation_notes}</p>
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
