import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  AlertTriangle,
  Play,
  CheckCircle2,
  Bug,
  Clock,
  Layers,
  ChevronRight,
  Terminal,
  Activity,
  FileCode,
  ShieldCheck
} from 'lucide-react';
import {
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell
} from 'recharts';
import { MistakeRecord } from '../../types/learning.types';
import analyticsService from '../../services/analytics.service';
import Card from '../../components/ui/Card';
import Badge from '../../components/ui/Badge';
import Button from '../../components/ui/Button';

export const MistakeAnalysisPage: React.FC = () => {
  const navigate = useNavigate();
  const [mistakes, setMistakes] = useState<MistakeRecord[]>([]);
  const [selectedCategory, setSelectedCategory] = useState<string>('all');
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const fetchMistakes = async () => {
      const data = await analyticsService.getMistakes();
      setMistakes(data);
      setLoading(false);
    };
    fetchMistakes();
  }, []);

  const categories = [
    { id: 'all', label: 'All Failure Modes' },
    { id: 'edge_cases', label: 'Edge Cases' },
    { id: 'time_complexity', label: 'Time Complexity' },
    { id: 'logical', label: 'Logical Errors' },
    { id: 'algorithm_selection', label: 'Algorithm Selection' }
  ];

  const filteredMistakes = mistakes.filter((m) => {
    if (selectedCategory === 'all') return true;
    return m.category === selectedCategory;
  });

  const categoryChartData = [
    { name: 'Edge Cases', value: 40, color: '#F59E0B' },
    { name: 'Time Complexity', value: 25, color: '#EF4444' },
    { name: 'Logical Flow', value: 20, color: '#8B5CF6' },
    { name: 'Algorithm Choice', value: 15, color: '#6366F1' }
  ];

  return (
    <div className="space-y-8 font-mono">
      {/* Page Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl sm:text-3xl font-bold tracking-tight text-[#F4F4F6] flex items-center gap-2">
            <AlertTriangle className="w-6 h-6 text-amber-400" />
            <span>Mistake Diagnostics & Anti-Patterns</span>
          </h1>
          <p className="text-xs text-[#8E8EA0] mt-1 font-sans">
            Automated code diagnostics classifying runtime faults, complexity regressions, and memory underflows.
          </p>
        </div>

        <div className="flex items-center gap-2 px-3 py-1 rounded-full bg-[#12121B] border border-[#1C1C28] text-xs text-[#A78BFA]">
          <Activity className="w-3.5 h-3.5 text-[#8B5CF6]" />
          <span>{mistakes.length} Telemetry Records Logged</span>
        </div>
      </div>

      {/* METRIC SUMMARY STRIP & PIE CHART */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <div className="lg:col-span-8 grid grid-cols-1 sm:grid-cols-3 gap-4">
          <div className="p-5 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-1">
            <span className="text-xs text-[#5A5A6E]">FAILED_RUNS</span>
            <div className="text-3xl font-bold text-red-400">18</div>
            <span className="text-[11px] text-[#8E8EA0] font-sans">Across 12 unique test suites</span>
          </div>

          <div className="p-5 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-1">
            <span className="text-xs text-[#5A5A6E]">TOP_ANTI_PATTERN</span>
            <div className="text-lg font-bold text-amber-400 truncate">Edge Case Underflow</div>
            <span className="text-[11px] text-[#8E8EA0] font-sans">Off-by-one index boundaries</span>
          </div>

          <div className="p-5 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-1">
            <span className="text-xs text-[#5A5A6E]">POST_HINT_RECOVERY</span>
            <div className="text-3xl font-bold text-emerald-400">84%</div>
            <span className="text-[11px] text-emerald-400/80 font-sans">Resolved in &lt;2 attempts</span>
          </div>
        </div>

        {/* Small Distribution Chart */}
        <div className="lg:col-span-4 p-5 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] flex items-center justify-between">
          <div className="space-y-1">
            <h4 className="text-xs font-bold text-[#F4F4F6]">Failure Category Ratio</h4>
            <div className="space-y-1 pt-1 text-[11px]">
              {categoryChartData.map((d) => (
                <div key={d.name} className="flex items-center gap-1.5 text-[#8E8EA0]">
                  <span className="w-2 h-2 rounded-full" style={{ backgroundColor: d.color }} />
                  <span>{d.name}</span>
                </div>
              ))}
            </div>
          </div>
          <div className="w-24 h-24">
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie
                  data={categoryChartData}
                  innerRadius={25}
                  outerRadius={40}
                  paddingAngle={4}
                  dataKey="value"
                >
                  {categoryChartData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={entry.color} />
                  ))}
                </Pie>
              </PieChart>
            </ResponsiveContainer>
          </div>
        </div>
      </div>

      {/* FILTER TABS */}
      <div className="flex flex-wrap items-center gap-2 bg-[#0D0D14] p-2 rounded-xl border border-[#1C1C28]">
        {categories.map((cat) => (
          <button
            key={cat.id}
            onClick={() => setSelectedCategory(cat.id)}
            className={`px-3 py-1.5 rounded-lg text-xs font-mono transition-all cursor-pointer ${
              selectedCategory === cat.id
                ? 'bg-[#8B5CF6] text-white shadow-sm'
                : 'bg-[#12121B] text-[#8E8EA0] hover:text-white border border-[#1C1C28]'
            }`}
          >
            {cat.label}
          </button>
        ))}
      </div>

      {/* DETAILED MISTAKE FORENSIC CARDS */}
      <div className="space-y-4">
        {filteredMistakes.map((record) => (
          <div
            key={record.id}
            className="p-6 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-4 hover:border-[#2D2D3E] transition-all"
          >
            {/* Header info */}
            <div className="flex flex-wrap items-center justify-between gap-3">
              <div className="flex items-center gap-3">
                <span className="text-[10px] uppercase font-mono px-2 py-0.5 rounded bg-red-500/10 text-red-400 border border-red-500/20">
                  {record.categoryLabel}
                </span>
                <span className="font-bold text-sm text-[#F4F4F6]">
                  {record.problemTitle}
                </span>
                <span className="text-xs text-[#5A5A6E]">
                  // {record.topic}
                </span>
              </div>

              <div className="flex items-center gap-2 text-xs">
                <span className="text-amber-400 bg-amber-500/10 px-2 py-0.5 rounded border border-amber-500/20 text-[11px]">
                  Failed {record.frequency}x
                </span>
                <span className="text-[#5A5A6E] text-[11px]">{record.date}</span>
              </div>
            </div>

            {/* Code Diff Simulation */}
            <div className="bg-[#09090E] rounded-xl border border-[#1C1C28] overflow-hidden text-xs">
              <div className="px-3.5 py-1.5 bg-[#0D0D14] border-b border-[#1C1C28] flex items-center justify-between text-[11px] text-[#5A5A6E]">
                <div className="flex items-center gap-1.5">
                  <FileCode className="w-3.5 h-3.5 text-red-400" />
                  <span>Sub-optimal Submission Excerpt</span>
                </div>
                <span>Flagged by LearnX Static AST Scanner</span>
              </div>

              <div className="p-3 text-red-300 overflow-x-auto whitespace-pre font-mono text-xs diff-del">
                {record.failedCodeSnippet}
              </div>
            </div>

            {/* Why Failed & Suggested Correction Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-xs font-sans">
              <div className="p-3.5 bg-[#12121B] rounded-xl border border-[#1C1C28] space-y-1">
                <span className="font-bold text-red-400 block font-mono">
                  Fault Mechanism
                </span>
                <p className="text-[#8E8EA0] leading-relaxed">{record.whyFailed}</p>
                <p className="text-[#5A5A6E] text-[11px] pt-1">{record.explanation}</p>
              </div>

              <div className="p-3.5 bg-[#12121B] rounded-xl border border-emerald-500/20 space-y-1">
                <span className="font-bold text-emerald-400 block font-mono">
                  Prescribed Fix
                </span>
                <p className="text-[#D1D1DE] leading-relaxed">
                  {record.suggestedCorrection}
                </p>
                <div className="pt-2 flex items-center gap-2 text-[11px] font-mono text-[#A78BFA]">
                  <span>Pillar:</span>
                  <strong className="underline">{record.recommendedTopicTitle}</strong>
                </div>
              </div>
            </div>

            {/* Action button */}
            <div className="pt-1 flex items-center justify-end">
              <Button
                variant="primary"
                size="sm"
                onClick={() => navigate(`/practice/${record.problemId}`)}
                leftIcon={<Play className="w-3 h-3" />}
                className="text-xs font-mono h-8"
              >
                Launch Drill
              </Button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default MistakeAnalysisPage;
