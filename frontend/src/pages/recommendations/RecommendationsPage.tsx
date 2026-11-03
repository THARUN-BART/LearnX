import React, { useEffect, useState, useCallback, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Sparkles,
  RefreshCw,
  History,
  Target,
  Zap,
  TrendingUp,
  BrainCircuit,
  HelpCircle,
  FolderOpen
} from 'lucide-react';
import { recommendationService } from '../../services/recommendation.service';
import {
  CategorizedRecommendations,
  ProblemRecommendation,
  RecommendationCategory,
  FeedbackRating
} from '../../types/recommendations.types';
import { TopPickHero } from '../../components/recommendations/TopPickHero';
import { RecommendationCategoryTabs } from '../../components/recommendations/RecommendationCategoryTabs';
import { ProblemRecommendationCard } from '../../components/recommendations/ProblemRecommendationCard';
import { RecommendationHistoryDrawer } from '../../components/recommendations/RecommendationHistoryDrawer';

export const RecommendationsPage: React.FC = () => {
  const navigate = useNavigate();

  const [data, setData] = useState<CategorizedRecommendations | null>(null);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [activeTab, setActiveTab] = useState<string>('recommended_for_you');
  const [historyOpen, setHistoryOpen] = useState(false);

  const loadData = useCallback(async () => {
    try {
      setLoading(true);
      const res = await recommendationService.getRecommendations();
      setData(res);
    } catch (err) {
      console.error('[RecommendationsPage] Error loading recommendations:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadData();
  }, [loadData]);

  const handleRefresh = async () => {
    try {
      setRefreshing(true);
      const res = await recommendationService.refreshRecommendations();
      setData(res);
    } catch (err) {
      console.error('[RecommendationsPage] Refresh failed:', err);
    } finally {
      setRefreshing(false);
    }
  };

  const handleFeedback = async (recId: string, rating: FeedbackRating, notes?: string) => {
    try {
      await recommendationService.submitFeedback(recId, { rating, notes });
      setData((prev) => {
        if (!prev) return prev;
        const updateRec = (rec: ProblemRecommendation): ProblemRecommendation => {
          if (rec.id === recId) {
            return { ...rec, feedback_rating: rating, feedback_notes: notes };
          }
          return rec;
        };

        const updatedCategories: Record<string, ProblemRecommendation[]> = {};
        for (const [catKey, list] of Object.entries(prev.categories || {})) {
          updatedCategories[catKey] = (list || []).map(updateRec);
        }

        return {
          ...prev,
          top_pick: prev.top_pick ? updateRec(prev.top_pick) : undefined,
          categories: updatedCategories,
        };
      });
    } catch (err) {
      console.error('[RecommendationsPage] Submit feedback failed:', err);
    }
  };

  const categoryCounts = useMemo<Record<string, number>>(() => {
    const counts: Record<string, number> = {};
    if (data?.categories) {
      for (const [k, v] of Object.entries(data.categories)) {
        counts[k] = v?.length || 0;
      }
    }
    return counts;
  }, [data]);

  const currentCategoryProblems: ProblemRecommendation[] = useMemo(() => {
    if (!data?.categories) return [];
    if (activeTab === 'all') {
      const all: ProblemRecommendation[] = [];
      const seen = new Set<string>();
      if (data.top_pick) {
        all.push(data.top_pick);
        seen.add(data.top_pick.id);
      }
      for (const list of Object.values(data.categories)) {
        for (const item of list) {
          if (!seen.has(item.id)) {
            seen.add(item.id);
            all.push(item);
          }
        }
      }
      return all;
    }
    return data.categories[activeTab] || [];
  }, [data, activeTab]);

  return (
    <div className="space-y-8 pb-16">
      {/* Top Banner Header */}
      <div className="relative overflow-hidden rounded-2xl border border-[#1C1C28] bg-gradient-to-r from-[#12121B] via-[#101019] to-[#12121B] p-6 lg:p-8">
        <div className="absolute top-0 right-0 w-96 h-96 bg-[#8B5CF6]/10 rounded-full blur-3xl pointer-events-none -mr-20 -mt-20" />
        <div className="relative z-10 flex flex-col md:flex-row md:items-center justify-between gap-6">
          <div className="space-y-2 max-w-2xl">
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-[#8B5CF6]/10 border border-[#8B5CF6]/30 text-[#A78BFA] text-xs font-mono font-medium">
              <Sparkles className="w-3.5 h-3.5" />
              <span>Deterministic Learning Vector Engine</span>
            </div>
            <h1 className="text-2xl lg:text-3xl font-bold tracking-tight text-[#F4F4F6]">
              Personalised Recommendations
            </h1>
            <p className="text-sm text-[#8E8EA0] leading-relaxed">
              Every suggested problem is scored in real-time against your concept mastery, recurring error patterns,
              and adaptive learning path to maximize retention.
            </p>
          </div>

          <div className="flex items-center gap-3">
            <button
              onClick={() => setHistoryOpen(true)}
              className="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl border border-[#2A2A3C] bg-[#161622] text-[#F4F4F6] text-xs font-medium hover:bg-[#1E1E30] transition-colors"
            >
              <History className="w-4 h-4 text-[#8E8EA0]" />
              <span>History</span>
            </button>

            <button
              onClick={handleRefresh}
              disabled={refreshing || loading}
              className="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl bg-gradient-to-r from-[#8B5CF6] to-[#6366F1] text-white text-xs font-medium font-mono hover:opacity-95 transition-opacity disabled:opacity-50 shadow-lg shadow-[#8B5CF6]/20"
            >
              <RefreshCw className={`w-4 h-4 ${refreshing ? 'animate-spin' : ''}`} />
              <span>{refreshing ? 'Re-scoring...' : 'Recalculate'}</span>
            </button>
          </div>
        </div>

        {/* Telemetry quick bar */}
        {data && (
          <div className="relative z-10 mt-6 grid grid-cols-2 sm:grid-cols-4 gap-3 pt-6 border-t border-[#1C1C28]/80 text-xs font-mono">
            <div className="flex items-center gap-2.5 p-2.5 rounded-lg bg-[#0D0D14]/70 border border-[#1C1C28]">
              <Target className="w-4 h-4 text-purple-400" />
              <div>
                <span className="text-[#5A5A6E] block text-[10px]">FOCUS CONCEPT</span>
                <span className="text-[#F4F4F6] font-semibold truncate block max-w-[130px]">
                  {data.primary_focus_concept || 'Arrays & Hashing'}
                </span>
              </div>
            </div>

            <div className="flex items-center gap-2.5 p-2.5 rounded-lg bg-[#0D0D14]/70 border border-[#1C1C28]">
              <TrendingUp className="w-4 h-4 text-emerald-400" />
              <div>
                <span className="text-[#5A5A6E] block text-[10px]">ADAPTIVE DIFFICULTY</span>
                <span className="text-emerald-400 font-semibold capitalize">
                  {data.recommended_difficulty || 'Medium'}
                </span>
              </div>
            </div>

            <div className="flex items-center gap-2.5 p-2.5 rounded-lg bg-[#0D0D14]/70 border border-[#1C1C28]">
              <BrainCircuit className="w-4 h-4 text-amber-400" />
              <div>
                <span className="text-[#5A5A6E] block text-[10px]">TARGET WEAKNESSES</span>
                <span className="text-amber-400 font-semibold">
                  {data.categories?.improve_weaknesses?.length || 0} Queued
                </span>
              </div>
            </div>

            <div className="flex items-center gap-2.5 p-2.5 rounded-lg bg-[#0D0D14]/70 border border-[#1C1C28]">
              <Zap className="w-4 h-4 text-blue-400" />
              <div>
                <span className="text-[#5A5A6E] block text-[10px]">CURATED PICKS</span>
                <span className="text-[#F4F4F6] font-semibold">
                  {data.total_recommendations || currentCategoryProblems.length} Available
                </span>
              </div>
            </div>
          </div>
        )}
      </div>

      {/* Loading Skeleton */}
      {loading ? (
        <div className="space-y-6">
          <div className="h-56 rounded-2xl bg-[#12121B] border border-[#1C1C28] animate-pulse" />
          <div className="h-12 w-full max-w-xl rounded-xl bg-[#12121B] animate-pulse" />
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
            {[1, 2, 3, 4, 5, 6].map((i) => (
              <div key={i} className="h-72 rounded-2xl bg-[#12121B] border border-[#1C1C28] animate-pulse" />
            ))}
          </div>
        </div>
      ) : data ? (
        <>
          {/* Top Pick Hero Section */}
          {data.top_pick && (
            <div>
              <div className="flex items-center justify-between mb-3 px-1">
                <div className="flex items-center gap-2">
                  <span className="inline-block w-2 h-2 rounded-full bg-[#8B5CF6] animate-ping" />
                  <span className="text-xs font-mono uppercase tracking-wider text-[#A78BFA] font-bold">
                    System Priority #1
                  </span>
                </div>
                <span className="text-[11px] text-[#5A5A6E] font-mono">
                  Calculated via composite weight algorithm
                </span>
              </div>

              <TopPickHero
                recommendation={data.top_pick}
                onFeedbackClick={(rec) => handleFeedback(rec.id, 'helpful')}
              />
            </div>
          )}

          {/* Categorized Recommendations Section */}
          <div className="space-y-6">
            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pt-4 border-t border-[#1C1C28]">
              <div>
                <h2 className="text-lg font-bold text-[#F4F4F6]">Curated Workstreams</h2>
                <p className="text-xs text-[#8E8EA0]">
                  Filtered algorithms targeting your specific learning vector
                </p>
              </div>
            </div>

            {/* Filter Tabs */}
            <RecommendationCategoryTabs
              activeTab={activeTab}
              onTabChange={setActiveTab}
              counts={categoryCounts}
            />

            {/* Problem Cards Grid */}
            {currentCategoryProblems.length > 0 ? (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
                {currentCategoryProblems.map((rec) => (
                  <ProblemRecommendationCard
                    key={rec.id}
                    recommendation={rec}
                    onFeedbackSubmit={(recId, rating) => handleFeedback(recId, rating)}
                  />
                ))}
              </div>
            ) : (
              <div className="rounded-2xl border border-dashed border-[#2A2A3C] bg-[#12121B]/50 p-12 text-center">
                <FolderOpen className="w-10 h-10 text-[#5A5A6E] mx-auto mb-3" />
                <h3 className="text-sm font-semibold text-[#F4F4F6] mb-1">
                  No problems queued in this workstream
                </h3>
                <p className="text-xs text-[#8E8EA0] max-w-sm mx-auto mb-4">
                  All current problems in this category have either been mastered or superseded by higher-priority objectives.
                </p>
                <button
                  onClick={() => setActiveTab('recommended_for_you')}
                  className="px-4 py-2 rounded-xl bg-[#1C1C28] text-xs text-[#F4F4F6] hover:bg-[#2A2A3C] transition-colors"
                >
                  View All Recommended
                </button>
              </div>
            )}
          </div>
        </>
      ) : (
        <div className="rounded-2xl border border-[#1C1C28] bg-[#12121B] p-12 text-center">
          <HelpCircle className="w-12 h-12 text-[#5A5A6E] mx-auto mb-4" />
          <h2 className="text-base font-semibold text-[#F4F4F6] mb-2">Unable to load recommendations</h2>
          <p className="text-xs text-[#8E8EA0] max-w-md mx-auto mb-6">
            We could not reach the recommendation engine. Please verify your connection or recalculate.
          </p>
          <button
            onClick={loadData}
            className="px-4 py-2 bg-[#8B5CF6] text-white rounded-xl text-xs font-medium hover:bg-[#7C3AED]"
          >
            Retry Connection
          </button>
        </div>
      )}

      {/* History Drawer */}
      <RecommendationHistoryDrawer
        isOpen={historyOpen}
        onClose={() => setHistoryOpen(false)}
      />
    </div>
  );
};

export default RecommendationsPage;
