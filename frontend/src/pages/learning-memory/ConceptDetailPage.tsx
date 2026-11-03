import { normalizePercentage } from '../../utils/formatters';
import React, { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import {
  ArrowLeft,
  Brain,
  CheckCircle2,
  AlertTriangle,
  Lightbulb,
  Clock,
  Play,
  Calendar,
} from 'lucide-react';
import { learningMemoryService } from '../../services/learning-memory.service';
import { ConceptDetailResponse } from '../../types/learning-memory.types';
import { LearningTimeline } from '../../components/learning-memory/LearningTimeline';

export const ConceptDetailPage: React.FC = () => {
  const { conceptKey } = useParams<{ conceptKey: string }>();
  const [detail, setDetail] = useState<ConceptDetailResponse | null>(null);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    if (conceptKey) {
      setLoading(true);
      learningMemoryService
        .getConceptDetail(conceptKey)
        .then((res) => setDetail(res))
        .catch((err) => console.error('Failed to load concept details', err))
        .finally(() => setLoading(false));
    }
  }, [conceptKey]);

  if (loading) {
    return (
      <div className="space-y-4 animate-pulse">
        <div className="h-8 w-48 bg-[#12121B] rounded-lg" />
        <div className="h-48 bg-[#12121B] rounded-2xl" />
      </div>
    );
  }

  if (!detail) {
    return (
      <div className="text-center py-12">
        <p className="text-xs text-[#8E8EA0]">Concept not found.</p>
        <Link to="/dashboard/learning-memory" className="text-xs text-[#8B5CF6] hover:underline mt-2 inline-block">
          Return to Learning Memory
        </Link>
      </div>
    );
  }

  const { concept, mistakes, timeline, next_steps } = detail;

  return (
    <div className="space-y-6 pb-12">
      {/* Back Link */}
      <Link
        to="/dashboard/learning-memory"
        className="inline-flex items-center gap-1.5 text-xs text-[#8E8EA0] hover:text-[#F4F4F6] transition-colors"
      >
        <ArrowLeft className="w-3.5 h-3.5" />
        Back to Learning Memory
      </Link>

      {/* Header Profile */}
      <div className="p-6 rounded-3xl bg-[#0D0D14] border border-[#1C1C28] flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div className="flex items-start gap-4">
          <div className="p-3 rounded-2xl bg-[#161622] text-[#8B5CF6]">
            <Brain className="w-6 h-6" />
          </div>
          <div>
            <div className="flex items-center gap-3">
              <h2 className="text-xl font-bold text-[#F4F4F6]">{concept.concept_name}</h2>
              <span className="px-2.5 py-0.5 rounded-full text-[10px] font-mono uppercase bg-[#8B5CF6]/10 text-[#8B5CF6] border border-[#8B5CF6]/20">
                {concept.learning_state}
              </span>
            </div>
            <p className="text-xs text-[#8E8EA0] mt-1">
              Topic: {concept.topic_name || 'DSA Core'} • Concept Identifier:{' '}
              <span className="font-mono text-[#5A5A6E]">{concept.concept_key}</span>
            </p>
          </div>
        </div>

        <Link
          to={`/practice`}
          className="px-4 py-2.5 rounded-xl bg-gradient-to-r from-[#8B5CF6] to-[#6366F1] text-white text-xs font-semibold flex items-center justify-center gap-2 hover:opacity-95 transition-opacity"
        >
          <Play className="w-4 h-4 fill-white" />
          Practice Related Problems
        </Link>
      </div>

      {/* Metric Breakdown Grid */}
      <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
        <div className="p-4 rounded-xl bg-[#0D0D14] border border-[#1C1C28]">
          <span className="text-[10px] uppercase font-mono text-[#8E8EA0]">Estimated Mastery</span>
          <div className="text-2xl font-bold text-[#F4F4F6] font-mono mt-1">
            {normalizePercentage(concept.estimated_mastery)}%
          </div>
          <span className="text-[10px] text-[#5A5A6E]">Provisional calculation</span>
        </div>

        <div className="p-4 rounded-xl bg-[#0D0D14] border border-[#1C1C28]">
          <span className="text-[10px] uppercase font-mono text-[#8E8EA0]">Independent Solves</span>
          <div className="text-2xl font-bold text-emerald-400 font-mono mt-1">
            {concept.independent_solves}
          </div>
          <span className="text-[10px] text-[#5A5A6E]">Unassisted verification</span>
        </div>

        <div className="p-4 rounded-xl bg-[#0D0D14] border border-[#1C1C28]">
          <span className="text-[10px] uppercase font-mono text-[#8E8EA0]">Assisted Solves</span>
          <div className="text-2xl font-bold text-blue-400 font-mono mt-1">
            {concept.assisted_solves}
          </div>
          <span className="text-[10px] text-[#5A5A6E]">{concept.hint_requests} hints requested</span>
        </div>

        <div className="p-4 rounded-xl bg-[#0D0D14] border border-[#1C1C28]">
          <span className="text-[10px] uppercase font-mono text-[#8E8EA0]">Mistakes</span>
          <div className="text-2xl font-bold text-rose-400 font-mono mt-1">
            {concept.mistake_count}
          </div>
          <span className="text-[10px] text-teal-400 font-mono">{concept.corrected_mistakes} verified fixed</span>
        </div>
      </div>

      {/* Recommended Next Steps */}
      {next_steps && next_steps.length > 0 && (
        <div className="p-5 rounded-2xl bg-[#0D0D14] border border-[#1C1C28] space-y-2">
          <div className="flex items-center gap-2 text-xs font-semibold text-[#F4F4F6]">
            <Lightbulb className="w-4 h-4 text-amber-400" />
            AI Tutor Pedagogical Next Steps
          </div>
          <ul className="space-y-1.5 pl-6 list-disc text-xs text-[#8E8EA0]">
            {next_steps.map((step, idx) => (
              <li key={idx}>{step}</li>
            ))}
          </ul>
        </div>
      )}

      {/* Mistake Pattern History */}
      <div className="space-y-3">
        <h3 className="text-sm font-bold text-[#F4F4F6] flex items-center gap-2">
          <AlertTriangle className="w-4 h-4 text-rose-400" />
          Recurring Mistake Memory for this Concept
        </h3>
        {mistakes.length === 0 ? (
          <p className="text-xs text-[#8E8EA0] italic bg-[#0D0D14] p-4 rounded-xl border border-[#1C1C28]">
            No recurring mistake patterns logged for this concept.
          </p>
        ) : (
          <div className="space-y-2">
            {mistakes.map((m) => (
              <div key={m.id} className="p-3.5 rounded-xl bg-[#0D0D14] border border-[#1C1C28] flex items-center justify-between">
                <div>
                  <div className="flex items-center gap-2">
                    <span className="text-xs font-semibold text-[#F4F4F6]">{m.category}</span>
                    <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-[#161622] text-[#8E8EA0]">
                      {m.subtype}
                    </span>
                  </div>
                  <p className="text-[11px] text-[#5A5A6E] mt-0.5">
                    Occurrences: {m.occurrence_count} • Corrected: {m.corrected_count}
                  </p>
                </div>
                <span
                  className={`px-2.5 py-0.5 rounded-full text-[10px] font-mono uppercase ${
                    m.status === 'resolved'
                      ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20'
                      : m.status === 'improving'
                      ? 'bg-blue-500/10 text-blue-400 border border-blue-500/20'
                      : 'bg-rose-500/10 text-rose-400 border border-rose-500/20'
                  }`}
                >
                  {m.status}
                </span>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Activity Timeline */}
      <div className="space-y-3">
        <h3 className="text-sm font-bold text-[#F4F4F6] flex items-center gap-2">
          <Clock className="w-4 h-4 text-[#8B5CF6]" />
          Recent Cognitive Activity Timeline
        </h3>
        <LearningTimeline events={timeline} />
      </div>
    </div>
  );
};

export default ConceptDetailPage;
