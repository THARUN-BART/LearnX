import React, { useEffect, useState, useCallback } from 'react';
import {
  Compass,
  RefreshCw,
  AlertCircle,
  Sparkles,
  Layers,
  ArrowRight,
  TrendingUp,
} from 'lucide-react';
import { learningPathService } from '../../services/learning-path.service';
import {
  AdaptivePath,
  LearningRecommendation,
  PathProgressResponse,
  PathNode,
} from '../../types/learning-path.types';
import { RecommendedActionCard } from '../../components/learning-path/RecommendedActionCard';
import { AdaptiveRoadmapGraph } from '../../components/learning-path/AdaptiveRoadmapGraph';
import { ConceptDetailModal } from '../../components/learning-path/ConceptDetailModal';
import { PathProgressStats } from '../../components/learning-path/PathProgressStats';
import { RecommendationHistoryList } from '../../components/learning-path/RecommendationHistoryList';

export const LearningPathPage: React.FC = () => {
  const [path, setPath] = useState<AdaptivePath | null>(null);
  const [nextAction, setNextAction] = useState<LearningRecommendation | null>(null);
  const [progress, setProgress] = useState<PathProgressResponse | null>(null);
  const [history, setHistory] = useState<LearningRecommendation[]>([]);
  const [selectedNode, setSelectedNode] = useState<PathNode | null>(null);

  const [loading, setLoading] = useState<boolean>(true);
  const [recalculating, setRecalculating] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const loadData = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);

      const [pathData, nextActionData, progressData, historyData] = await Promise.all([
        learningPathService.getLearningPath(),
        learningPathService.getNextAction().catch(() => null),
        learningPathService.getProgress().catch(() => null),
        learningPathService.getHistory(15, 0).catch(() => ({ history: [], total: 0 })),
      ]);

      setPath(pathData);
      setNextAction(nextActionData);
      setProgress(progressData);
      setHistory(historyData.history || []);
    } catch (err: any) {
      console.error('Failed to load adaptive learning path:', err);
      setError(err?.response?.data?.error || 'Failed to load your adaptive learning path.');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadData();
  }, [loadData]);

  const handleRecalculate = async () => {
    try {
      setRecalculating(true);
      const updatedPath = await learningPathService.recalculate();
      setPath(updatedPath);

      // Refresh next action and progress stats
      const [nextActionData, progressData, historyData] = await Promise.all([
        learningPathService.getNextAction().catch(() => null),
        learningPathService.getProgress().catch(() => null),
        learningPathService.getHistory(15, 0).catch(() => ({ history: [], total: 0 })),
      ]);
      setNextAction(nextActionData);
      setProgress(progressData);
      setHistory(historyData.history || []);
    } catch (err: any) {
      console.error('Recalculation error:', err);
    } finally {
      setRecalculating(false);
    }
  };

  return (
    <div className="min-h-screen bg-[#F8FAFC] text-[#0F172A] p-4 sm:p-6 lg:p-8 space-y-8 max-w-7xl mx-auto">
      {/* Top Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-2 border-b border-[#E2E8F0]">
        <div>
          <div className="flex items-center gap-2 mb-1">
            <div className="w-8 h-8 rounded-xl bg-gradient-to-tr from-[#EF4444] to-[#B91C1C] flex items-center justify-center shadow-md shadow-[#EF4444]/30">
              <Compass className="w-4 h-4 text-white" />
            </div>
            <h1 className="text-2xl font-bold tracking-tight text-[#0F172A]">
              Adaptive Learning Path
            </h1>
          </div>
          <p className="text-xs text-[#64748B] max-w-xl">
            Real-time personalized DSA progression powered by Learning Memory, prerequisite graphs, and error-pattern analysis.
          </p>
        </div>

        <button
          onClick={handleRecalculate}
          disabled={recalculating}
          className="inline-flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl bg-[#F1F5F9] hover:bg-[#E2E8F0] border border-[#E2E8F0] text-xs font-semibold text-[#0F172A] transition-all disabled:opacity-50"
        >
          <RefreshCw className={`w-3.5 h-3.5 ${recalculating ? 'animate-spin text-[#EF4444]' : ''}`} />
          <span>{recalculating ? 'Evaluating Evidence...' : 'Re-evaluate Path'}</span>
        </button>
      </div>

      {/* Error state if any */}
      {error && (
        <div className="p-4 rounded-xl bg-rose-500/10 border border-rose-500/30 text-xs text-rose-300 flex items-center justify-between gap-3">
          <div className="flex items-center gap-2">
            <AlertCircle className="w-4 h-4 text-rose-400 shrink-0" />
            <span>{error}</span>
          </div>
          <button
            onClick={loadData}
            className="px-3 py-1 bg-rose-500/20 hover:bg-rose-500/30 rounded text-rose-200 font-semibold"
          >
            Retry
          </button>
        </div>
      )}

      {/* Section 1: Recommended Next Action Hero */}
      <div>
        <div className="flex items-center gap-2 mb-3">
          <Sparkles className="w-4 h-4 text-[#EF4444]" />
          <h2 className="text-sm font-semibold uppercase tracking-wider font-mono text-[#64748B]">
            Recommended Next Action
          </h2>
        </div>
        <RecommendedActionCard
          recommendation={nextAction}
          loading={loading}
          onRefresh={handleRecalculate}
        />
      </div>

      {/* Section 2: Progress & Mastery Breakdown */}
      <PathProgressStats progress={progress} loading={loading} />

      {/* Section 3: Interactive Adaptive Roadmap Graph */}
      <div>
        <AdaptiveRoadmapGraph
          nodes={path?.nodes || []}
          selectedConceptKey={selectedNode?.concept_key || null}
          onSelectConcept={(node) => setSelectedNode(node)}
        />
      </div>

      {/* Section 4: Recommendation History & Progression Timeline */}
      <RecommendationHistoryList history={history} loading={loading} />

      {/* Concept Details Modal */}
      {selectedNode && (
        <ConceptDetailModal node={selectedNode} onClose={() => setSelectedNode(null)} />
      )}
    </div>
  );
};

export default LearningPathPage;
