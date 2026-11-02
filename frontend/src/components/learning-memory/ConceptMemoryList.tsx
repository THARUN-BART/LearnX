import React from 'react';
import { Link } from 'react-router-dom';
import { ArrowRight, Sparkles, CheckCircle2, AlertCircle } from 'lucide-react';
import { ConceptMemory } from '../../types/learning-memory.types';

interface Props {
  concepts: ConceptMemory[];
  loading?: boolean;
}

export const ConceptMemoryList: React.FC<Props> = ({ concepts, loading }) => {
  if (loading) {
    return (
      <div className="space-y-3">
        {[1, 2, 3].map((n) => (
          <div key={n} className="h-16 rounded-xl bg-[#12121B] animate-pulse border border-[#1C1C28]" />
        ))}
      </div>
    );
  }

  if (concepts.length === 0) {
    return (
      <div className="p-8 text-center bg-[#0D0D14] border border-[#1C1C28] rounded-2xl">
        <Sparkles className="w-8 h-8 text-[#8B5CF6] mx-auto mb-2 opacity-50" />
        <h4 className="text-sm font-semibold text-[#F4F4F6]">No concepts recorded yet</h4>
        <p className="text-xs text-[#8E8EA0] max-w-sm mx-auto mt-1">
          As you submit code in the Practice Arena, your concept profile and cognitive telemetry will appear here automatically.
        </p>
      </div>
    );
  }

  const getStateBadge = (state: string) => {
    switch (state) {
      case 'mastered':
        return <span className="px-2.5 py-0.5 rounded-full text-[10px] font-mono uppercase bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">Mastered</span>;
      case 'struggling':
        return <span className="px-2.5 py-0.5 rounded-full text-[10px] font-mono uppercase bg-rose-500/10 text-rose-400 border border-rose-500/20">Needs Focus</span>;
      case 'reinforcing':
        return <span className="px-2.5 py-0.5 rounded-full text-[10px] font-mono uppercase bg-blue-500/10 text-blue-400 border border-blue-500/20">Reinforcing</span>;
      default:
        return <span className="px-2.5 py-0.5 rounded-full text-[10px] font-mono uppercase bg-purple-500/10 text-purple-400 border border-purple-500/20">Practicing</span>;
    }
  };

  return (
    <div className="space-y-3">
      {concepts.map((c) => (
        <div
          key={c.id || c.concept_key}
          className="p-4 rounded-xl bg-[#0D0D14] border border-[#1C1C28] hover:border-[#8B5CF6]/40 transition-all duration-200 flex flex-col md:flex-row md:items-center justify-between gap-4"
        >
          <div className="flex items-start gap-3">
            <div className="mt-0.5 p-2 rounded-lg bg-[#161622] text-[#8B5CF6]">
              {c.learning_state === 'mastered' ? (
                <CheckCircle2 className="w-4 h-4 text-emerald-400" />
              ) : c.learning_state === 'struggling' ? (
                <AlertCircle className="w-4 h-4 text-rose-400" />
              ) : (
                <Sparkles className="w-4 h-4" />
              )}
            </div>
            <div>
              <div className="flex items-center gap-2">
                <h4 className="font-semibold text-sm text-[#F4F4F6]">{c.concept_name}</h4>
                {getStateBadge(c.learning_state)}
              </div>
              <p className="text-xs text-[#8E8EA0] mt-0.5">
                Topic: {c.topic_name || 'General DSA'} • Last practiced: {new Date(c.last_practiced_at).toLocaleDateString()}
              </p>
            </div>
          </div>

          <div className="flex items-center justify-between md:justify-end gap-6 border-t md:border-t-0 pt-3 md:pt-0 border-[#1C1C28]">
            <div className="flex items-center gap-4 text-xs font-mono">
              <div className="text-center">
                <span className="text-[#8E8EA0] block text-[10px] uppercase">Solves</span>
                <span className="text-emerald-400 font-bold">{c.independent_solves} ind</span> /{' '}
                <span className="text-blue-400 font-bold">{c.assisted_solves} ass</span>
              </div>
              <div className="text-center">
                <span className="text-[#8E8EA0] block text-[10px] uppercase">Mistakes</span>
                <span className="text-rose-400 font-bold">{c.mistake_count}</span> ({c.corrected_mistakes} fixed)
              </div>
              <div className="text-center">
                <span className="text-[#8E8EA0] block text-[10px] uppercase">Est. Mastery</span>
                <span className="text-[#F4F4F6] font-bold">{c.estimated_mastery}%</span>
              </div>
            </div>

            <Link
              to={`/dashboard/learning-memory/concepts/${encodeURIComponent(c.concept_key)}`}
              className="p-2 rounded-lg bg-[#161622] text-[#8E8EA0] hover:text-[#F4F4F6] hover:bg-[#8B5CF6]/20 transition-colors"
              title="View full concept cognitive memory"
            >
              <ArrowRight className="w-4 h-4" />
            </Link>
          </div>
        </div>
      ))}
    </div>
  );
};
