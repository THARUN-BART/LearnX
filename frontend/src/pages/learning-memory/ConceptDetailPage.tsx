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
        <div className="h-8 w-48 bg-[#F1F5F9] rounded-lg" />
        <div className="h-48 bg-[#F1F5F9] rounded-2xl" />
      </div>
    );
  }

  if (!detail) {
    return (
      <div className="text-center py-12">
        <p className="text-xs text-[#64748B]">Concept not found.</p>
        <Link to="/dashboard/learning-memory" className="text-xs text-[#EF4444] hover:underline mt-2 inline-block">
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
        className="inline-flex items-center gap-1.5 text-xs text-[#64748B] hover:text-[#0F172A] transition-colors"
      >
        <ArrowLeft className="w-3.5 h-3.5" />
        Back to Learning Memory
      </Link>

      {/* Header Profile */}
      <div className="p-6 rounded-3xl bg-white border border-[#E2E8F0] flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div className="flex items-start gap-4">
          <div className="p-3 rounded-2xl bg-[#F1F5F9] text-[#EF4444]">
            <Brain className="w-6 h-6" />
          </div>
          <div>
            <div className="flex items-center gap-3">
              <h2 className="text-xl font-bold text-[#0F172A]">{concept.concept_name}</h2>
              <span className="px-2.5 py-0.5 rounded-full text-[10px] font-mono uppercase bg-[#EF4444]/10 text-[#EF4444] border border-[#EF4444]/20">
                {concept.learning_state}
              </span>
            </div>
            <p className="text-xs text-[#64748B] mt-1">
              Topic: {concept.topic_name || 'DSA Core'} • Concept Identifier:{' '}
              <span className="font-mono text-[#94A3B8]">{concept.concept_key}</span>
            </p>
          </div>
        </div>

        <Link
          to={`/practice`}
          className="px-4 py-2.5 rounded-xl bg-gradient-to-r from-[#EF4444] to-[#B91C1C] text-white text-xs font-semibold flex items-center justify-center gap-2 hover:opacity-95 transition-opacity"
        >
          <Play className="w-4 h-4 fill-white" />
          Practice Related Problems
        </Link>
      </div>

      {/* Metric Breakdown Grid */}
      <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
        <div className="p-4 rounded-xl bg-white border border-[#E2E8F0]">
          <span className="text-[10px] uppercase font-mono text-[#64748B]">Estimated Mastery</span>
          <div className="text-2xl font-bold text-[#0F172A] font-mono mt-1">
            {normalizePercentage(concept.estimated_mastery)}%
          </div>
          <span className="text-[10px] text-[#94A3B8]">Provisional calculation</span>
        </div>

        <div className="p-4 rounded-xl bg-white border border-[#E2E8F0]">
          <span className="text-[10px] uppercase font-mono text-[#64748B]">Independent Solves</span>
          <div className="text-2xl font-bold text-emerald-400 font-mono mt-1">
            {concept.independent_solves}
          </div>
          <span className="text-[10px] text-[#94A3B8]">Unassisted verification</span>
        </div>

        <div className="p-4 rounded-xl bg-white border border-[#E2E8F0]">
          <span className="text-[10px] uppercase font-mono text-[#64748B]">Assisted Solves</span>
          <div className="text-2xl font-bold text-blue-400 font-mono mt-1">
            {concept.assisted_solves}
          </div>
          <span className="text-[10px] text-[#94A3B8]">{concept.hint_requests} hints requested</span>
        </div>

        <div className="p-4 rounded-xl bg-white border border-[#E2E8F0]">
          <span className="text-[10px] uppercase font-mono text-[#64748B]">Mistakes</span>
          <div className="text-2xl font-bold text-rose-400 font-mono mt-1">
            {concept.mistake_count}
          </div>
          <span className="text-[10px] text-teal-400 font-mono">{concept.corrected_mistakes} verified fixed</span>
        </div>
      </div>

      {/* Recommended Next Steps */}
      {next_steps && next_steps.length > 0 && (
        <div className="p-5 rounded-2xl bg-white border border-[#E2E8F0] space-y-2">
          <div className="flex items-center gap-2 text-xs font-semibold text-[#0F172A]">
            <Lightbulb className="w-4 h-4 text-amber-400" />
            AI Tutor Pedagogical Next Steps
          </div>
          <ul className="space-y-1.5 pl-6 list-disc text-xs text-[#64748B]">
            {next_steps.map((step, idx) => (
              <li key={idx}>{step}</li>
            ))}
          </ul>
        </div>
      )}

      {/* Mistake Pattern History */}
      <div className="space-y-3">
        <h3 className="text-sm font-bold text-[#0F172A] flex items-center gap-2">
          <AlertTriangle className="w-4 h-4 text-rose-400" />
          Recurring Mistake Memory for this Concept
        </h3>
        {mistakes.length === 0 ? (
          <p className="text-xs text-[#64748B] italic bg-white p-4 rounded-xl border border-[#E2E8F0]">
            No recurring mistake patterns logged for this concept.
          </p>
        ) : (
          <div className="space-y-2">
            {mistakes.map((m) => (
              <div key={m.id} className="p-3.5 rounded-xl bg-white border border-[#E2E8F0] flex items-center justify-between">
                <div>
                  <div className="flex items-center gap-2">
                    <span className="text-xs font-semibold text-[#0F172A]">{m.category}</span>
                    <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-[#F1F5F9] text-[#64748B]">
                      {m.subtype}
                    </span>
                  </div>
                  <p className="text-[11px] text-[#94A3B8] mt-0.5">
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
        <h3 className="text-sm font-bold text-[#0F172A] flex items-center gap-2">
          <Clock className="w-4 h-4 text-[#EF4444]" />
          Recent Cognitive Activity Timeline
        </h3>
        <LearningTimeline events={timeline} />
      </div>
    </div>
  );
};

export default ConceptDetailPage;
