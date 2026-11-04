import React, { useState, useEffect, useRef } from 'react';
import { useSearchParams } from 'react-router-dom';
import {
  AlgorithmMetadata,
  AlgorithmCategory,
  ExecuteResponse,
  ExecutionStep,
} from '../../types/visualizer.types';
import visualizerService from '../../services/visualizer.service';
import PlaybackControls from '../../components/visualizer/PlaybackControls';
import DataStructureCanvas from '../../components/visualizer/DataStructureCanvas';
import ExplanationPanel from '../../components/visualizer/ExplanationPanel';
import AIAssistantDrawer from '../../components/visualizer/AIAssistantDrawer';
import PracticeQuizModal from '../../components/visualizer/PracticeQuizModal';
import CustomInputModal from '../../components/visualizer/CustomInputModal';
import {
  Search,
  SlidersHorizontal,
  RotateCcw,
  Sparkles,
  Layers,
  ChevronRight,
} from 'lucide-react';
import Button from '../../components/ui/Button';
import { toast } from 'sonner';

const CATEGORIES: (AlgorithmCategory | 'All')[] = [
  'All',
  'Sorting',
  'Searching',
  'Arrays and Patterns',
  'Linked Lists',
  'Stack and Queue',
  'Recursion and Backtracking',
  'Trees',
  'Graphs',
  'Dynamic Programming',
];

export const VisualizerPage: React.FC = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  const [algorithms, setAlgorithms] = useState<AlgorithmMetadata[]>([]);
  const [selectedCategory, setSelectedCategory] = useState<string>('All');
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedAlgo, setSelectedAlgo] = useState<AlgorithmMetadata | null>(null);

  // Execution state
  const [executionData, setExecutionData] = useState<ExecuteResponse | null>(null);
  const [currentStepIndex, setCurrentStepIndex] = useState(0);
  const [isPlaying, setIsPlaying] = useState(false);
  const [speed, setSpeed] = useState(1);
  const [loading, setLoading] = useState(false);

  // Modals & Drawers
  const [customInputOpen, setCustomInputOpen] = useState(false);
  const [aiDrawerOpen, setAiDrawerOpen] = useState(false);
  const [quizModalOpen, setQuizModalOpen] = useState(false);

  // Active Session tracking
  const sessionIdRef = useRef<string | null>(null);
  const timerRef = useRef<any>(null);

  // Load all algorithms
  useEffect(() => {
    visualizerService.getAlgorithms().then((list) => {
      if (list && list.length > 0) {
        setAlgorithms(list);
        const queryAlgoId = searchParams.get('algo');
        const found = queryAlgoId ? list.find((a) => a.id === queryAlgoId) : null;
        const initial = found || list[0];
        setSelectedAlgo(initial);
      }
    });
  }, [searchParams]);

  // Execute selected algorithm
  const runAlgorithm = async (algo: AlgorithmMetadata, customInput?: Record<string, any>) => {
    setLoading(true);
    setIsPlaying(false);
    try {
      const resp = await visualizerService.executeAlgorithm(algo.id, customInput);
      if (resp) {
        setExecutionData(resp);
        setCurrentStepIndex(0);

        // Record backend session
        visualizerService
          .createSession({
            algorithmId: algo.id,
            category: algo.category,
            inputConfig: customInput || algo.default_input,
            totalSteps: resp.total_steps,
          })
          .then((sess) => {
            if (sess) sessionIdRef.current = sess.id;
          });
      } else {
        toast.error('Failed to generate algorithm trace');
      }
    } catch {
      toast.error('Visualizer execution error');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (selectedAlgo) {
      runAlgorithm(selectedAlgo);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedAlgo]);

  // Playback loop
  useEffect(() => {
    if (isPlaying && executionData) {
      const intervalMs = Math.max(250, 1000 / speed);
      timerRef.current = setInterval(() => {
        setCurrentStepIndex((prev) => {
          if (prev >= executionData.steps.length - 1) {
            setIsPlaying(false);
            if (sessionIdRef.current) {
              visualizerService.updateSession(sessionIdRef.current, {
                stepsCompleted: executionData.steps.length,
                timeSpentSeconds: 10,
                completed: true,
              });
            }
            return prev;
          }
          return prev + 1;
        });
      }, intervalMs);
    } else {
      if (timerRef.current) clearInterval(timerRef.current);
    }
    return () => {
      if (timerRef.current) clearInterval(timerRef.current);
    };
  }, [isPlaying, executionData, speed]);

  const handlePlayPause = () => {
    if (executionData && currentStepIndex >= executionData.steps.length - 1) {
      setCurrentStepIndex(0);
    }
    setIsPlaying(!isPlaying);
  };

  const handleNext = () => {
    if (executionData && currentStepIndex < executionData.steps.length - 1) {
      setCurrentStepIndex((prev) => prev + 1);
    }
  };

  const handlePrev = () => {
    if (currentStepIndex > 0) {
      setCurrentStepIndex((prev) => prev - 1);
    }
  };

  const handleRestart = () => {
    setIsPlaying(false);
    setCurrentStepIndex(0);
  };

  const handleSelectAlgo = (algo: AlgorithmMetadata) => {
    setSelectedAlgo(algo);
    setSearchParams({ algo: algo.id });
  };

  // Filter algorithms
  const filteredAlgorithms = algorithms.filter((a) => {
    const matchesCat = selectedCategory === 'All' || a.category === selectedCategory;
    const matchesQuery =
      a.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      a.description.toLowerCase().includes(searchQuery.toLowerCase()) ||
      a.tags.some((t) => t.toLowerCase().includes(searchQuery.toLowerCase()));
    return matchesCat && matchesQuery;
  });

  const currentStepData: ExecutionStep | null =
    executionData?.steps?.[currentStepIndex] || null;

  return (
    <div className="flex flex-col gap-6 font-mono text-[#0F172A] pb-12">
      {/* Top Banner / Hero */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 p-6 bg-gradient-to-r from-[#12121B] via-[#0D0D14] to-[#12121B] border border-[#E2E8F0] rounded-2xl relative overflow-hidden">
        <div className="space-y-1 relative z-10">
          <div className="flex items-center gap-2">
            <span className="text-[10px] uppercase font-bold tracking-widest text-[#EF4444] bg-[#EF4444]/10 px-2 py-0.5 rounded border border-[#EF4444]/20">
              Interactive Engine
            </span>
            <span className="text-[10px] text-emerald-400 font-bold">
              40+ Supported Algorithms
            </span>
          </div>
          <h1 className="text-xl sm:text-2xl font-bold font-sans tracking-tight text-white">
            Algorithm Visualizer & Execution Trace
          </h1>
          <p className="text-xs text-[#64748B] font-sans max-w-xl">
            Step through real algorithm executions, inspect active pointers and memory states in real-time, and get Socratic AI assistance.
          </p>
        </div>
      </div>

      {/* Category Pills & Search */}
      <div className="flex flex-col gap-3">
        <div className="flex items-center gap-3">
          <div className="relative flex-1 max-w-sm">
            <Search className="w-4 h-4 text-[#94A3B8] absolute left-3 top-1/2 -translate-y-1/2" />
            <input
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="Search algorithms, data structures, tags..."
              className="w-full bg-white text-[#0F172A] text-xs pl-9 pr-3 py-2 rounded-xl border border-[#E2E8F0] outline-none font-sans"
            />
          </div>
        </div>

        {/* Categories Bar */}
        <div className="flex items-center gap-1.5 overflow-x-auto pb-1 scrollbar-none">
          {CATEGORIES.map((cat) => (
            <button
              key={cat}
              onClick={() => setSelectedCategory(cat)}
              className={`text-xs px-3 py-1.5 rounded-xl whitespace-nowrap transition-all font-sans font-medium ${
                selectedCategory === cat
                  ? 'bg-[#EF4444] text-white shadow-md shadow-[#8B5CF6]/30 font-bold'
                  : 'bg-white text-[#64748B] hover:text-[#0F172A] border border-[#E2E8F0]'
              }`}
            >
              {cat}
            </button>
          ))}
        </div>
      </div>

      {/* Main Grid: Left Sidebar Selector (col-span-3) + Right Arena (col-span-9) */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-5 items-start">
        {/* Left: Algorithm Catalog List */}
        <div className="lg:col-span-4 bg-white border border-[#E2E8F0] rounded-2xl p-4 flex flex-col gap-3 max-h-[780px] overflow-y-auto">
          <div className="flex items-center justify-between text-xs text-[#64748B] pb-2 border-b border-[#E2E8F0]">
            <span className="font-bold uppercase tracking-wider text-[10px]">
              Algorithms ({filteredAlgorithms.length})
            </span>
          </div>

          <div className="space-y-2">
            {filteredAlgorithms.map((algo) => {
              const isSelected = selectedAlgo?.id === algo.id;
              return (
                <button
                  key={algo.id}
                  onClick={() => handleSelectAlgo(algo)}
                  className={`w-full text-left p-3 rounded-xl border transition-all flex flex-col gap-1.5 ${
                    isSelected
                      ? 'bg-[#EF4444]/15 border-[#EF4444] text-white shadow-md shadow-[#8B5CF6]/20'
                      : 'bg-[#F1F5F9] border-[#E2E8F0] text-[#64748B] hover:border-[#EF4444]/40 hover:text-[#0F172A]'
                  }`}
                >
                  <div className="flex items-center justify-between">
                    <span className="text-xs font-bold font-sans text-[#0F172A]">
                      {algo.name}
                    </span>
                    <span className="text-[9px] uppercase font-bold text-[#DC2626] bg-[#EF4444]/10 px-1.5 py-0.5 rounded">
                      {algo.category}
                    </span>
                  </div>
                  <p className="text-[11px] text-[#64748B] font-sans line-clamp-2">
                    {algo.description}
                  </p>
                  <div className="flex items-center gap-2 pt-1 text-[10px] text-[#94A3B8]">
                    <span>Avg: {algo.complexity.average}</span>
                    <span>·</span>
                    <span>Space: {algo.complexity.space}</span>
                  </div>
                </button>
              );
            })}
          </div>
        </div>

        {/* Right: Interactive Visualizer Workspace */}
        <div className="lg:col-span-8 flex flex-col gap-4">
          {selectedAlgo && (
            <>
              {/* Workspace Action Toolbar */}
              <div className="flex flex-wrap items-center justify-between gap-3 bg-white border border-[#E2E8F0] px-4 py-2.5 rounded-2xl">
                <div className="flex items-center gap-2">
                  <Layers className="w-4 h-4 text-[#EF4444]" />
                  <span className="text-sm font-bold text-white font-sans">
                    {selectedAlgo.name}
                  </span>
                  <span className="text-[10px] text-[#64748B] hidden sm:inline">
                    ({executionData?.total_steps ?? 0} Steps)
                  </span>
                </div>

                <div className="flex items-center gap-2">
                  <Button
                    variant="secondary"
                    size="sm"
                    onClick={() => setCustomInputOpen(true)}
                    leftIcon={<SlidersHorizontal className="w-3.5 h-3.5" />}
                    className="text-xs h-8"
                  >
                    Custom Input
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => runAlgorithm(selectedAlgo)}
                    leftIcon={<RotateCcw className="w-3.5 h-3.5" />}
                    className="text-xs h-8 text-[#64748B]"
                  >
                    Reset
                  </Button>
                </div>
              </div>

              {/* Data Structure Canvas */}
              <DataStructureCanvas
                currentStepData={currentStepData}
                category={selectedAlgo.category}
              />

              {/* Playback Controls */}
              <PlaybackControls
                isPlaying={isPlaying}
                onPlayPause={handlePlayPause}
                onNext={handleNext}
                onPrev={handlePrev}
                onRestart={handleRestart}
                currentStep={currentStepIndex}
                totalSteps={executionData?.total_steps ?? 0}
                onSeek={(step) => setCurrentStepIndex(step)}
                speed={speed}
                onSpeedChange={(s) => setSpeed(s)}
              />

              {/* Algorithm Explanation Panel */}
              <ExplanationPanel
                algorithm={selectedAlgo}
                currentStepData={currentStepData}
                onOpenAI={() => setAiDrawerOpen(true)}
                onOpenQuiz={() => setQuizModalOpen(true)}
              />
            </>
          )}
        </div>
      </div>

      {/* AI Assistant Drawer */}
      {selectedAlgo && (
        <AIAssistantDrawer
          isOpen={aiDrawerOpen}
          onClose={() => setAiDrawerOpen(false)}
          algorithmId={selectedAlgo.id}
          algorithmName={selectedAlgo.name}
          currentStepData={currentStepData}
          totalSteps={executionData?.total_steps ?? 0}
        />
      )}

      {/* Practice Quiz Modal */}
      {selectedAlgo && (
        <PracticeQuizModal
          isOpen={quizModalOpen}
          onClose={() => setQuizModalOpen(false)}
          algorithmId={selectedAlgo.id}
          algorithmName={selectedAlgo.name}
        />
      )}

      {/* Custom Input Modal */}
      {selectedAlgo && (
        <CustomInputModal
          isOpen={customInputOpen}
          onClose={() => setCustomInputOpen(false)}
          inputType={selectedAlgo.input_type}
          defaultInput={selectedAlgo.default_input}
          onApply={(newInput) => runAlgorithm(selectedAlgo, newInput)}
        />
      )}
    </div>
  );
};
export default VisualizerPage;

