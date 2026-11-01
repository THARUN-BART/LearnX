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
  Cpu
} from 'lucide-react';
import { SAMPLE_PROBLEMS, PROGRAMMING_LANGUAGES } from '../../utils/constants';
import { Problem, CodeExecutionResult } from '../../types/learning.types';
import codingService from '../../services/coding.service';
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
  const [code, setCode] = useState<string>(() => problem.starterCode['python'] || '');
  const [activeTab, setActiveTab] = useState<'description' | 'hints'>('description');
  const [activeBottomTab, setActiveBottomTab] = useState<'testcases' | 'console'>('testcases');
  const [selectedTestCaseIndex, setSelectedTestCaseIndex] = useState<number>(0);
  const [isFullscreen, setIsFullscreen] = useState<boolean>(false);

  // Execution states
  const [isRunning, setIsRunning] = useState<boolean>(false);
  const [executionResult, setExecutionResult] = useState<CodeExecutionResult | null>(null);

  // When problem or language changes, update starter code
  useEffect(() => {
    const currentProb =
      SAMPLE_PROBLEMS.find((p) => p.id === problemId || p.slug === problemId) || SAMPLE_PROBLEMS[0];
    setProblem(currentProb);
    setCode(currentProb.starterCode[selectedLanguage.toLowerCase()] || currentProb.starterCode['python'] || '');
    setExecutionResult(null);
  }, [problemId, selectedLanguage]);

  const handleLanguageChange = (lang: string) => {
    setSelectedLanguage(lang);
    setCode(problem.starterCode[lang.toLowerCase()] || '');
  };

  const handleResetCode = () => {
    setCode(problem.starterCode[selectedLanguage.toLowerCase()] || '');
    toast.info('Starter code reset');
  };

  const handleCopyCode = () => {
    navigator.clipboard.writeText(code);
    toast.success('Code copied to clipboard');
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
            className="bg-[#12121B] text-[#F4F4F6] text-xs font-mono font-semibold px-3 py-1.5 rounded-lg border border-[#1C1C28] outline-none cursor-pointer"
          >
            {SAMPLE_PROBLEMS.map((p) => (
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
                <div className="text-xs text-[#A78BFA] font-mono flex items-center gap-1.5">
                  <Sparkles className="w-4 h-4 text-[#8B5CF6]" />
                  <span>Cognitive Thinking Cues</span>
                </div>
                {problem.hints.map((hint, i) => (
                  <div
                    key={i}
                    className="p-3.5 bg-[#12121B] rounded-xl border border-[#1C1C28] space-y-1"
                  >
                    <span className="text-xs font-bold text-[#F4F4F6] font-mono">
                      Phase #{i + 1}
                    </span>
                    <p className="text-xs text-[#8E8EA0] pt-1 font-sans">{hint}</p>
                  </div>
                ))}
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
    </div>
  );
};

export default PracticePage;
