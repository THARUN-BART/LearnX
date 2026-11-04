import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  AlertTriangle,
  Play,
  Activity,
  FileCode,
  ShieldCheck,
  ThumbsUp,
  ThumbsDown,
} from 'lucide-react';
import { ResponsiveContainer, PieChart, Pie, Cell } from 'recharts';
import { MistakeRecord, MistakeSummary, StudentMistakePattern } from '../../types/learning.types';
import analyticsService from '../../services/analytics.service';
import aiService from '../../services/ai.service';
import Button from '../../components/ui/Button';
import { PageSkeleton } from '../../components/ui/Skeleton';
import { toast } from 'sonner';

export const MistakeAnalysisPage: React.FC = () => {
  const navigate = useNavigate();
  const [mistakes, setMistakes] = useState<MistakeRecord[]>([]);
  const [summary, setSummary] = useState<MistakeSummary | null>(null);
  const [patterns, setPatterns] = useState<StudentMistakePattern[]>([]);
  const [selectedCategory, setSelectedCategory] = useState<string>('all');
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const fetchMistakes = async () => {
      const [data, sum, pats] = await Promise.all([
        analyticsService.getMistakes(),
        analyticsService.getMistakeSummary(),
        analyticsService.getMistakePatterns(),
      ]);
      setMistakes(data);
      setSummary(sum);
      setPatterns(pats);
      setLoading(false);
    };
    fetchMistakes();
  }, []);

  const categories = [
    { id: 'all', label: 'All Failure Modes' },
    { id: 'edge_cases', label: 'Edge Cases' },
    { id: 'time_complexity', label: 'Time Complexity' },
    { id: 'logical', label: 'Logical Errors' },
    { id: 'algorithm_selection', label: 'Algorithm Selection' },
    { id: 'syntax', label: 'Compilation' },
  ];

  const filteredMistakes = mistakes.filter((m) => {
    if (selectedCategory === 'all') return true;
    return m.category === selectedCategory;
  });

  const countsFromSummary = summary?.category_counts || {};
  const edgeCount =
    countsFromSummary.edge_case ??
    countsFromSummary.edge_cases ??
    mistakes.filter((m) => m.category === 'edge_cases').length;
  const timeCount =
    countsFromSummary.complexity ??
    mistakes.filter((m) => m.category === 'time_complexity').length;
  const logicCount =
    (countsFromSummary.logical ?? 0) + (countsFromSummary.runtime ?? 0) ||
    mistakes.filter((m) => m.category === 'logical' || m.category === 'syntax').length;
  const algoCount =
    (countsFromSummary.algorithmic ?? 0) + (countsFromSummary.data_structure ?? 0) ||
    mistakes.filter((m) => m.category === 'algorithm_selection' || m.category === 'space_complexity').length;

  const rawChartData = [
    { name: 'Edge Cases', value: edgeCount, color: '#F59E0B' },
    { name: 'Time Complexity', value: timeCount, color: '#EF4444' },
    { name: 'Logical Flow', value: logicCount, color: '#8B5CF6' },
    { name: 'Algorithm Choice', value: algoCount, color: '#6366F1' },
  ].filter((item) => item.value > 0);

  const categoryChartData =
    rawChartData.length > 0 ? rawChartData : [{ name: 'Clean Runs', value: 1, color: '#22C55E' }];

  const topPattern = patterns[0] || summary?.recurring_patterns?.[0];
  const totalFailed = summary?.total_mistakes ?? mistakes.length;
  const topLabel = topPattern
    ? `${(topPattern.category || '').replace(/_/g, ' ')} / ${(topPattern.subtype || '').replace(/_/g, ' ')}`
    : 'None yet';

  const handleFeedback = async (id: string, type: 'helpful' | 'unhelpful') => {
    try {
      await aiService.submitMistakeFeedback(id, type);
      toast.success('Feedback recorded');
    } catch {
      toast.error('Could not record feedback');
    }
  };

  if (loading) return <PageSkeleton rows={3} />;

  return (
    <div className="space-y-8 font-mono">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl sm:text-3xl font-bold tracking-tight text-[#0F172A] flex items-center gap-2">
            <AlertTriangle className="w-6 h-6 text-amber-400" />
            Mistake Diagnostics & Anti-Patterns
          </h1>
          <p className="text-xs text-[#64748B] mt-1 font-sans">
            Diagnoses from failed submissions, recurring patterns, and suggested drills.
          </p>
        </div>
        <div className="flex items-center gap-2 px-3 py-1 rounded-full bg-[#F1F5F9] border border-[#E2E8F0] text-xs text-[#DC2626]">
          <Activity className="w-3.5 h-3.5 text-[#EF4444]" />
          {totalFailed} Telemetry Records Logged
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <div className="lg:col-span-8 grid grid-cols-1 sm:grid-cols-3 gap-4">
          <div className="p-5 rounded-2xl bg-white border border-[#E2E8F0] space-y-1">
            <span className="text-xs text-[#94A3B8]">FAILED_RUNS</span>
            <div className="text-3xl font-bold text-red-400">{totalFailed}</div>
            <span className="text-[11px] text-[#64748B] font-sans">Diagnosed submissions</span>
          </div>
          <div className="p-5 rounded-2xl bg-white border border-[#E2E8F0] space-y-1">
            <span className="text-xs text-[#94A3B8]">TOP_ANTI_PATTERN</span>
            <div className="text-lg font-bold text-amber-400 truncate">{topLabel}</div>
            <span className="text-[11px] text-[#64748B] font-sans">
              {topPattern ? `${topPattern.occurrence_count} occurrences` : 'No recurring pattern yet'}
            </span>
          </div>
          <div className="p-5 rounded-2xl bg-white border border-[#E2E8F0] space-y-1">
            <span className="text-xs text-[#94A3B8]">FOCUS_CONCEPTS</span>
            <div className="text-sm font-bold text-emerald-400 truncate">
              {(summary?.needs_attention_concepts || []).slice(0, 2).join(', ') || 'Keep solving'}
            </div>
            <span className="text-[11px] text-emerald-400/80 font-sans">Concepts to drill next</span>
          </div>
        </div>

        <div className="lg:col-span-4 p-5 rounded-2xl bg-white border border-[#E2E8F0] flex items-center justify-between">
          <div className="space-y-1">
            <h4 className="text-xs font-bold text-[#0F172A]">Failure Category Ratio</h4>
            <div className="space-y-1 pt-1 text-[11px]">
              {categoryChartData.map((d) => (
                <div key={d.name} className="flex items-center gap-1.5 text-[#64748B]">
                  <span className="w-2 h-2 rounded-full" style={{ backgroundColor: d.color }} />
                  <span>
                    {d.name}
                    {d.name !== 'Clean Runs' ? ` (${d.value})` : ''}
                  </span>
                </div>
              ))}
            </div>
          </div>
          <div className="w-24 h-24">
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie data={categoryChartData} innerRadius={25} outerRadius={40} paddingAngle={4} dataKey="value">
                  {categoryChartData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={entry.color} />
                  ))}
                </Pie>
              </PieChart>
            </ResponsiveContainer>
          </div>
        </div>
      </div>

      {patterns.length > 0 && (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
          {patterns.slice(0, 3).map((p) => (
            <div key={p.id} className="p-4 rounded-xl bg-white border border-[#E2E8F0] text-xs space-y-1">
              <div className="text-[#DC2626] font-bold">
                {(p.related_concept || p.subtype || 'Pattern').replace(/_/g, ' ')}
              </div>
              <div className="text-[#64748B]">
                {p.category.replace(/_/g, ' ')} · {p.occurrence_count}x · {p.correction_status}
              </div>
            </div>
          ))}
        </div>
      )}

      <div className="flex flex-wrap items-center gap-2 bg-white p-2 rounded-xl border border-[#E2E8F0]">
        {categories.map((cat) => (
          <button
            key={cat.id}
            onClick={() => setSelectedCategory(cat.id)}
            className={`px-3 py-1.5 rounded-lg text-xs font-mono transition-all cursor-pointer ${
              selectedCategory === cat.id
                ? 'bg-[#EF4444] text-white shadow-sm'
                : 'bg-[#F1F5F9] text-[#64748B] hover:text-[#0F172A] border border-[#E2E8F0]'
            }`}
          >
            {cat.label}
          </button>
        ))}
      </div>

      <div className="space-y-4">
        {filteredMistakes.length === 0 ? (
          <div className="p-12 text-center rounded-2xl bg-white border border-[#E2E8F0] space-y-3">
            <ShieldCheck className="w-8 h-8 text-emerald-400 mx-auto" />
            <h3 className="text-sm font-bold text-[#0F172A]">No Failure Modes Detected</h3>
            <p className="text-xs text-[#64748B] max-w-md mx-auto">
              Submit solutions in the Practice Arena to track debugging telemetry.
            </p>
            <Button size="sm" onClick={() => navigate('/practice')} className="text-xs">
              Open Practice Arena
            </Button>
          </div>
        ) : (
          filteredMistakes.map((record) => (
            <div
              key={record.id}
              className="p-6 rounded-2xl bg-white border border-[#E2E8F0] space-y-4 hover:border-[#2D2D3E] transition-all"
            >
              <div className="flex flex-wrap items-center justify-between gap-3">
                <div className="flex items-center gap-3">
                  <span className="text-[10px] uppercase font-mono px-2 py-0.5 rounded bg-red-500/10 text-red-400 border border-red-500/20">
                    {record.categoryLabel}
                  </span>
                  <span className="font-bold text-sm text-[#0F172A]">{record.problemTitle}</span>
                  <span className="text-xs text-[#94A3B8]">// {record.topic}</span>
                </div>
                <div className="flex items-center gap-2 text-xs">
                  <span className="text-amber-400 bg-amber-500/10 px-2 py-0.5 rounded border border-amber-500/20 text-[11px]">
                    Failed {record.frequency}x
                  </span>
                  <span className="text-[#94A3B8] text-[11px]">{record.date}</span>
                </div>
              </div>

              {record.failedCodeSnippet && (
                <div className="bg-[#09090E] rounded-xl border border-[#E2E8F0] overflow-hidden text-xs">
                  <div className="px-3.5 py-1.5 bg-white border-b border-[#E2E8F0] flex items-center gap-1.5 text-[11px] text-[#94A3B8]">
                    <FileCode className="w-3.5 h-3.5 text-red-400" />
                    Submission excerpt
                  </div>
                  <div className="p-3 text-red-300 overflow-x-auto whitespace-pre font-mono text-xs">
                    {record.failedCodeSnippet}
                  </div>
                </div>
              )}

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-xs font-sans">
                <div className="p-3.5 bg-[#F1F5F9] rounded-xl border border-[#E2E8F0] space-y-1">
                  <span className="font-bold text-red-400 block font-mono">Fault Mechanism</span>
                  <p className="text-[#64748B] leading-relaxed">{record.whyFailed}</p>
                  <p className="text-[#94A3B8] text-[11px] pt-1">{record.explanation}</p>
                </div>
                <div className="p-3.5 bg-[#F1F5F9] rounded-xl border border-emerald-500/20 space-y-1">
                  <span className="font-bold text-emerald-400 block font-mono">Prescribed Fix</span>
                  <p className="text-[#D1D1DE] leading-relaxed">{record.suggestedCorrection}</p>
                  <div className="pt-2 text-[11px] font-mono text-[#DC2626]">
                    Pillar: <strong>{record.recommendedTopicTitle}</strong>
                  </div>
                </div>
              </div>

              <div className="pt-1 flex items-center justify-end gap-2">
                <button
                  onClick={() => handleFeedback(record.id, 'helpful')}
                  className="p-1.5 rounded-lg border border-[#E2E8F0] text-[#64748B] hover:text-emerald-400"
                  title="Helpful"
                >
                  <ThumbsUp className="w-3.5 h-3.5" />
                </button>
                <button
                  onClick={() => handleFeedback(record.id, 'unhelpful')}
                  className="p-1.5 rounded-lg border border-[#E2E8F0] text-[#64748B] hover:text-red-400"
                  title="Unhelpful"
                >
                  <ThumbsDown className="w-3.5 h-3.5" />
                </button>
                {record.problemId && (
                  <Button
                    variant="primary"
                    size="sm"
                    onClick={() => navigate(`/practice/${record.problemId}`)}
                    leftIcon={<Play className="w-3 h-3" />}
                    className="text-xs font-mono h-8"
                  >
                    Launch Drill
                  </Button>
                )}
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
};

export default MistakeAnalysisPage;
