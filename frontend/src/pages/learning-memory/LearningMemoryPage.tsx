import React, { useEffect, useState } from 'react';
import { Sparkles, Brain, Search, Filter, RefreshCcw } from 'lucide-react';
import { learningMemoryService } from '../../services/learning-memory.service';
import {
  LearningMemorySummary,
  ConceptMemory,
  ConceptStrength,
  ConceptWeakness,
  TimelinePoint,
  StudentFacingNarrative,
  LearningSnapshot,
} from '../../types/learning-memory.types';
import { LearningMetricCards } from '../../components/learning-memory/LearningMetricCards';
import { ConceptMemoryList } from '../../components/learning-memory/ConceptMemoryList';
import { StrengthsWeaknessesSection } from '../../components/learning-memory/StrengthsWeaknessesSection';
import { LearningActivityChart } from '../../components/learning-memory/LearningActivityChart';
import { MistakeCorrectionChart } from '../../components/learning-memory/MistakeCorrectionChart';
import { LearningSnapshotsHistory } from '../../components/learning-memory/LearningSnapshotsHistory';

export const LearningMemoryPage: React.FC = () => {
  const [summary, setSummary] = useState<LearningMemorySummary | null>(null);
  const [narrative, setNarrative] = useState<StudentFacingNarrative | null>(null);
  const [concepts, setConcepts] = useState<ConceptMemory[]>([]);
  const [strengths, setStrengths] = useState<ConceptStrength[]>([]);
  const [weaknesses, setWeaknesses] = useState<ConceptWeakness[]>([]);
  const [timeline, setTimeline] = useState<TimelinePoint[]>([]);
  const [snapshots, setSnapshots] = useState<LearningSnapshot[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  // Filters
  const [search, setSearch] = useState<string>('');
  const [stateFilter, setStateFilter] = useState<string>('');

  const loadData = async () => {
    setLoading(true);
    try {
      const [sum, narr, concRes, strList, weakList, timeData, snapData] = await Promise.all([
        learningMemoryService.getSummary(),
        learningMemoryService.getNarrative(),
        learningMemoryService.listConcepts({ search, learning_state: stateFilter }),
        learningMemoryService.getStrengths(5),
        learningMemoryService.getWeaknesses(5),
        learningMemoryService.getTimeline(14),
        learningMemoryService.getSnapshots(30),
      ]);
      setSummary(sum);
      setNarrative(narr);
      setConcepts(concRes.concepts);
      setStrengths(strList);
      setWeaknesses(weakList);
      setTimeline(timeData);
      setSnapshots(snapData);
    } catch (err) {
      console.error('Failed to load learning memory profile', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, [stateFilter]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    loadData();
  };

  return (
    <div className="space-y-8 pb-12">
      {/* Top Narrative Banner */}
      <div className="p-6 rounded-3xl bg-gradient-to-r from-[#8B5CF6]/15 via-[#6366F1]/10 to-transparent border border-[#8B5CF6]/30 relative overflow-hidden">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div className="space-y-1.5 max-w-2xl">
            <div className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-[11px] font-mono uppercase bg-[#8B5CF6]/20 text-[#8B5CF6] border border-[#8B5CF6]/40">
              <Sparkles className="w-3.5 h-3.5" />
              Cognitive Learning Profile
            </div>
            <h2 className="text-xl md:text-2xl font-bold text-[#F4F4F6] tracking-tight">
              {narrative?.headline || 'Your Learning Memory Engine'}
            </h2>
            <p className="text-xs text-[#8E8EA0] leading-relaxed">
              {narrative?.narrative || 'LearnX continuously observes your solves, hint requests, and mistake corrections across all DSA topics.'}
            </p>
          </div>

          <button
            onClick={loadData}
            className="self-start md:self-auto px-4 py-2 rounded-xl bg-[#12121B] hover:bg-[#161622] text-[#8E8EA0] hover:text-[#F4F4F6] border border-[#1C1C28] text-xs font-mono flex items-center gap-2 transition-all"
          >
            <RefreshCcw className="w-3.5 h-3.5" />
            Refresh Telemetry
          </button>
        </div>
      </div>

      {/* Summary Metric Cards */}
      {summary && <LearningMetricCards summary={summary} />}

      {/* Charts Section */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <LearningActivityChart data={timeline} />
        <MistakeCorrectionChart data={timeline} />
      </div>

      {/* Mastery Snapshots & Long-Term Milestones */}
      <LearningSnapshotsHistory snapshots={snapshots} loading={loading} />

      {/* Strengths & Weaknesses */}
      <StrengthsWeaknessesSection strengths={strengths} weaknesses={weaknesses} />

      {/* Concept Learning Overview List with Filters */}
      <div className="space-y-4">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-[#1C1C28] pb-4">
          <div className="flex items-center gap-2">
            <div className="p-2 rounded-xl bg-[#161622] text-[#8B5CF6]">
              <Brain className="w-5 h-5" />
            </div>
            <div>
              <h3 className="text-base font-bold text-[#F4F4F6]">Concept Memory Telemetry</h3>
              <p className="text-xs text-[#8E8EA0]">Persistent mastery breakdown and attempt evidence</p>
            </div>
          </div>

          <div className="flex flex-wrap items-center gap-3">
            <form onSubmit={handleSearch} className="relative">
              <Search className="w-3.5 h-3.5 text-[#5A5A6E] absolute left-3 top-1/2 -translate-y-1/2" />
              <input
                type="text"
                placeholder="Search concepts..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="pl-9 pr-3 py-1.5 rounded-xl bg-[#0D0D14] border border-[#1C1C28] text-xs text-[#F4F4F6] placeholder-[#5A5A6E] focus:outline-none focus:border-[#8B5CF6]"
              />
            </form>

            <div className="flex items-center gap-1.5 bg-[#0D0D14] border border-[#1C1C28] p-1 rounded-xl">
              <Filter className="w-3.5 h-3.5 text-[#5A5A6E] ml-2" />
              <select
                value={stateFilter}
                onChange={(e) => setStateFilter(e.target.value)}
                className="bg-transparent text-xs text-[#8E8EA0] px-2 py-1 focus:outline-none cursor-pointer"
              >
                <option value="" className="bg-[#0D0D14]">All States</option>
                <option value="discovered" className="bg-[#0D0D14]">Discovered</option>
                <option value="practicing" className="bg-[#0D0D14]">Practicing</option>
                <option value="reinforcing" className="bg-[#0D0D14]">Reinforcing</option>
                <option value="mastered" className="bg-[#0D0D14]">Mastered</option>
                <option value="struggling" className="bg-[#0D0D14]">Needs Focus</option>
              </select>
            </div>
          </div>
        </div>

        <ConceptMemoryList concepts={concepts} loading={loading} />
      </div>
    </div>
  );
};

export default LearningMemoryPage;
