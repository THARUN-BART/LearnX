import { normalizePercentage } from '../../utils/formatters';
import React from 'react';
import {
  Lock,
  CheckCircle2,
  AlertTriangle,
  PlayCircle,
  Clock,
  Sparkles,
  ArrowRight,
  GitCommit,
} from 'lucide-react';
import { PathNode, ProgressionState } from '../../types/learning-path.types';

interface AdaptiveRoadmapGraphProps {
  nodes: PathNode[];
  selectedConceptKey: string | null;
  onSelectConcept: (node: PathNode) => void;
}

const stateConfig: Record<
  ProgressionState,
  {
    label: string;
    badgeBg: string;
    badgeText: string;
    border: string;
    bg: string;
    icon: React.FC<{ className?: string }>;
    accentColor: string;
  }
> = {
  mastered: {
    label: 'Mastered',
    badgeBg: 'bg-emerald-500/15',
    badgeText: 'text-emerald-400',
    border: 'border-emerald-500/40',
    bg: 'bg-emerald-950/20',
    icon: CheckCircle2,
    accentColor: '#10B981',
  },
  needs_revision: {
    label: 'Needs Revision',
    badgeBg: 'bg-amber-500/15',
    badgeText: 'text-amber-400',
    border: 'border-amber-500/40',
    bg: 'bg-amber-950/20',
    icon: AlertTriangle,
    accentColor: '#F59E0B',
  },
  practicing: {
    label: 'Practicing',
    badgeBg: 'bg-purple-500/15',
    badgeText: 'text-purple-400',
    border: 'border-purple-500/40',
    bg: 'bg-purple-950/20',
    icon: Sparkles,
    accentColor: '#8B5CF6',
  },
  in_progress: {
    label: 'In Progress',
    badgeBg: 'bg-blue-500/15',
    badgeText: 'text-blue-400',
    border: 'border-blue-500/40',
    bg: 'bg-blue-950/20',
    icon: Clock,
    accentColor: '#3B82F6',
  },
  ready_to_learn: {
    label: 'Ready to Learn',
    badgeBg: 'bg-cyan-500/15',
    badgeText: 'text-cyan-400',
    border: 'border-cyan-500/40',
    bg: 'bg-cyan-950/20',
    icon: PlayCircle,
    accentColor: '#06B6D4',
  },
  not_started: {
    label: 'Not Started',
    badgeBg: 'bg-zinc-800',
    badgeText: 'text-zinc-400',
    border: 'border-zinc-700/50',
    bg: 'bg-zinc-900/30',
    icon: GitCommit,
    accentColor: '#71717A',
  },
  locked_by_prerequisite: {
    label: 'Locked',
    badgeBg: 'bg-red-500/10',
    badgeText: 'text-zinc-500',
    border: 'border-zinc-800',
    bg: 'bg-zinc-950/40 opacity-70',
    icon: Lock,
    accentColor: '#4B5563',
  },
};

export const AdaptiveRoadmapGraph: React.FC<AdaptiveRoadmapGraphProps> = ({
  nodes,
  selectedConceptKey,
  onSelectConcept,
}) => {
  return (
    <div className="rounded-2xl border border-[#E2E8F0] bg-white p-5 md:p-6">
      <div className="flex flex-wrap items-center justify-between gap-4 mb-6 pb-4 border-b border-[#E2E8F0]">
        <div>
          <h3 className="text-base font-semibold text-[#0F172A] flex items-center gap-2">
            <span>Adaptive DSA Curriculum Graph</span>
            <span className="text-xs font-mono font-normal text-[#64748B]">
              ({nodes.length} Concepts)
            </span>
          </h3>
          <p className="text-xs text-[#64748B] mt-0.5">
            Nodes dynamically adapt their progression state based on your real performance and prerequisite mastery.
          </p>
        </div>

        {/* Legend */}
        <div className="flex flex-wrap items-center gap-2 text-[11px] font-mono">
          <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
            <span className="w-1.5 h-1.5 rounded-full bg-emerald-400" /> Mastered
          </span>
          <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-cyan-500/10 text-cyan-400 border border-cyan-500/20">
            <span className="w-1.5 h-1.5 rounded-full bg-cyan-400" /> Ready
          </span>
          <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-purple-500/10 text-purple-400 border border-purple-500/20">
            <span className="w-1.5 h-1.5 rounded-full bg-purple-400" /> Practicing
          </span>
          <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-amber-500/10 text-amber-400 border border-amber-500/20">
            <span className="w-1.5 h-1.5 rounded-full bg-amber-400" /> Needs Revision
          </span>
          <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-zinc-800 text-zinc-400 border border-zinc-700">
            <span className="w-1.5 h-1.5 rounded-full bg-zinc-400" /> Locked
          </span>
        </div>
      </div>

      {/* Grid of Nodes with Prerequisite Connectors */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3.5">
        {nodes.map((node, index) => {
          const cfg = stateConfig[node.state] || stateConfig.not_started;
          const StateIcon = cfg.icon;
          const isSelected = selectedConceptKey === node.concept_key;
          const isLocked = node.state === 'locked_by_prerequisite';

          return (
            <div
              key={node.concept_key}
              onClick={() => onSelectConcept(node)}
              className={`relative cursor-pointer rounded-xl border p-4 transition-all duration-200 flex flex-col justify-between ${
                cfg.bg
              } ${isSelected ? 'ring-2 ring-[#8B5CF6] border-[#EF4444]' : cfg.border} ${
                isLocked
                  ? 'hover:border-zinc-700'
                  : 'hover:scale-[1.01] hover:border-[#EF4444]/50 shadow-lg shadow-black/20'
              }`}
            >
              {/* Top Row: Order & Status */}
              <div className="flex items-center justify-between gap-2 mb-2">
                <div className="flex items-center gap-2">
                  <span className="w-5 h-5 rounded-md bg-[#F1F5F9] text-[10px] font-mono text-[#64748B] flex items-center justify-center font-bold">
                    {index + 1}
                  </span>
                  <span className="text-[10px] uppercase font-mono tracking-wider text-[#717182] truncate max-w-[120px]">
                    {node.topic_name}
                  </span>
                </div>

                <span
                  className={`inline-flex items-center gap-1 px-2 py-0.5 rounded-md text-[10px] font-semibold border ${cfg.badgeBg} ${cfg.badgeText} ${cfg.border}`}
                >
                  <StateIcon className="w-3 h-3" />
                  <span>{cfg.label}</span>
                </span>
              </div>

              {/* Center: Concept Name & Readiness */}
              <div className="my-1.5">
                <h4 className="text-sm font-semibold text-[#0F172A] tracking-tight line-clamp-1">
                  {node.concept_name}
                </h4>

                {/* Readiness score bar */}
                <div className="mt-2.5">
                  <div className="flex items-center justify-between text-[11px] font-mono text-[#64748B] mb-1">
                    <span>Readiness</span>
                    <span className="font-semibold" style={{ color: cfg.accentColor }}>
                      {normalizePercentage(node.readiness_score)}%
                    </span>
                  </div>
                  <div className="w-full h-1.5 rounded-full bg-[#F1F5F9] overflow-hidden">
                    <div
                      className="h-full rounded-full transition-all duration-500"
                      style={{
                        width: `${Math.min(100, Math.max(5, normalizePercentage(node.readiness_score)))}%`,
                        backgroundColor: cfg.accentColor,
                      }}
                    />
                  </div>
                </div>
              </div>

              {/* Bottom: Prerequisites & Difficulty */}
              <div className="pt-2.5 mt-2 border-t border-[#E2E8F0]/80 flex items-center justify-between text-[10px] font-mono text-[#64748B]">
                <div className="truncate max-w-[160px]">
                  {node.prerequisite_concept_keys && node.prerequisite_concept_keys.length > 0 ? (
                    <span className="text-[#64748B]">
                      Req: {node.prerequisite_concept_keys.join(', ')}
                    </span>
                  ) : (
                    <span className="text-emerald-500/80">No Prereq</span>
                  )}
                </div>

                <span className="text-[#A0A0B2] font-medium bg-[#F1F5F9] px-1.5 py-0.5 rounded">
                  {node.recommended_difficulty}
                </span>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};
