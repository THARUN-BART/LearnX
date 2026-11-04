import { normalizePercentage } from '../../utils/formatters';
import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  X,
  CheckCircle2,
  AlertTriangle,
  PlayCircle,
  Clock,
  Sparkles,
  ArrowRight,
  BookOpen,
  Lock,
  Layers,
  Flame,
} from 'lucide-react';
import { PathNode, ProgressionState } from '../../types/learning-path.types';
import { learningMemoryService } from '../../services/learning-memory.service';
import { ConceptDetailResponse } from '../../types/learning-memory.types';

interface ConceptDetailModalProps {
  node: PathNode | null;
  onClose: () => void;
}

export const ConceptDetailModal: React.FC<ConceptDetailModalProps> = ({ node, onClose }) => {
  const navigate = useNavigate();
  const [detail, setDetail] = useState<ConceptDetailResponse | null>(null);
  const [loading, setLoading] = useState<boolean>(false);

  useEffect(() => {
    if (!node) return;
    let cancelled = false;
    setLoading(true);

    learningMemoryService
      .getConceptDetail(node.concept_key)
      .then((data) => {
        if (!cancelled) setDetail(data);
      })
      .catch((err) => {
        console.warn('Could not load memory detail for concept:', err);
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });

    return () => {
      cancelled = true;
    };
  }, [node]);

  if (!node) return null;

  const isLocked = node.state === 'locked_by_prerequisite';

  const handleStartPractice = () => {
    if (node.recommended_problem_id) {
      navigate(`/practice/${node.recommended_problem_id}`);
    } else if (node.topic_id) {
      navigate(`/roadmap/${node.topic_id}`);
    } else {
      navigate('/practice');
    }
    onClose();
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm animate-fadeIn">
      <div className="relative w-full max-w-2xl rounded-2xl border border-[#252538] bg-[#0E0E17] p-6 shadow-2xl overflow-y-auto max-h-[90vh]">
        {/* Header */}
        <div className="flex items-start justify-between gap-4 pb-4 border-b border-[#E2E8F0]">
          <div>
            <div className="flex items-center gap-2 mb-1 text-xs font-mono text-[#64748B]">
              <span>Topic: {node.topic_name}</span>
              <span>•</span>
              <span className="capitalize">{node.state.replace(/_/g, ' ')}</span>
            </div>
            <h2 className="text-xl font-bold text-[#0F172A]">{node.concept_name}</h2>
          </div>
          <button
            onClick={onClose}
            className="p-1.5 rounded-lg text-[#64748B] hover:text-[#0F172A] hover:bg-[#F1F5F9] transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Content Body */}
        <div className="py-5 space-y-5">
          {/* Readiness & Difficulty Summary */}
          <div className="grid grid-cols-2 sm:grid-cols-3 gap-3">
            <div className="p-3.5 rounded-xl bg-[#141420] border border-[#E2E8F0]">
              <span className="text-[11px] font-mono text-[#64748B] block mb-1">Readiness Score</span>
              <span className="text-xl font-bold text-[#EF4444] font-mono">
                {normalizePercentage(node.readiness_score)}%
              </span>
            </div>
            <div className="p-3.5 rounded-xl bg-[#141420] border border-[#E2E8F0]">
              <span className="text-[11px] font-mono text-[#64748B] block mb-1">Recommended Difficulty</span>
              <span className="text-xl font-bold text-[#10B981] font-mono">
                {node.recommended_difficulty}
              </span>
            </div>
            <div className="p-3.5 rounded-xl bg-[#141420] border border-[#E2E8F0] col-span-2 sm:col-span-1">
              <span className="text-[11px] font-mono text-[#64748B] block mb-1">Curriculum Sequence</span>
              <span className="text-xl font-bold text-[#0F172A] font-mono">
                #{node.node_order}
              </span>
            </div>
          </div>

          {/* Prerequisites Chain */}
          <div className="p-4 rounded-xl bg-[#12121D] border border-[#E2E8F0]">
            <h4 className="text-xs font-semibold text-[#64748B] uppercase tracking-wider font-mono mb-2">
              Prerequisite Dependencies
            </h4>
            {node.prerequisite_concept_keys && node.prerequisite_concept_keys.length > 0 ? (
              <div className="flex flex-wrap gap-2">
                {node.prerequisite_concept_keys.map((prereq) => (
                  <span
                    key={prereq}
                    className="inline-flex items-center gap-1.5 px-3 py-1 rounded-lg bg-[#F1F5F9] text-xs font-mono text-[#C4C4D4] border border-[#252538]"
                  >
                    <BookOpen className="w-3.5 h-3.5 text-[#EF4444]" />
                    {prereq}
                  </span>
                ))}
              </div>
            ) : (
              <p className="text-xs text-emerald-400 font-mono">
                ✓ Foundational concept — No prior prerequisites required.
              </p>
            )}
          </div>

          {/* Learning Memory Diagnostics if available */}
          {detail?.concept ? (
            <div className="space-y-3">
              <h4 className="text-xs font-semibold text-[#64748B] uppercase tracking-wider font-mono">
                Diagnostic Learning Memory
              </h4>
              <div className="grid grid-cols-2 sm:grid-cols-4 gap-2.5 text-xs font-mono">
                <div className="p-2.5 rounded-lg bg-[#141420] border border-[#E2E8F0]">
                  <span className="text-[#64748B] block">Mastery Score</span>
                  <span className="font-bold text-[#0F172A] text-sm">
                    {normalizePercentage(detail.concept.estimated_mastery)}%
                  </span>
                </div>
                <div className="p-2.5 rounded-lg bg-[#141420] border border-[#E2E8F0]">
                  <span className="text-[#64748B] block">Confidence</span>
                  <span className="font-bold text-[#0F172A] text-sm capitalize">
                    {detail.concept.confidence_level || 'Provisional'}
                  </span>
                </div>
                <div className="p-2.5 rounded-lg bg-[#141420] border border-[#E2E8F0]">
                  <span className="text-[#64748B] block">Total Solves</span>
                  <span className="font-bold text-[#0F172A] text-sm">
                    {detail.concept.independent_solves || 0} Ind / {detail.concept.assisted_solves || 0} Asst
                  </span>
                </div>
                <div className="p-2.5 rounded-lg bg-[#141420] border border-[#E2E8F0]">
                  <span className="text-[#64748B] block">Mistake Count</span>
                  <span className="font-bold text-amber-400 text-sm">
                    {detail.concept.mistake_count || 0}
                  </span>
                </div>
              </div>

              {/* Active mistakes list */}
              {detail.mistakes && detail.mistakes.length > 0 && (
                <div className="p-3.5 rounded-xl bg-amber-500/10 border border-amber-500/20 text-xs">
                  <div className="flex items-center gap-2 text-amber-300 font-semibold mb-1.5">
                    <AlertTriangle className="w-4 h-4 text-amber-400" />
                    <span>Identified Error Patterns in this Concept:</span>
                  </div>
                  <ul className="list-disc list-inside space-y-1 text-amber-200/90 pl-1">
                    {detail.mistakes.map((m) => (
                      <li key={m.id}>
                        {m.subtype.replace(/_/g, ' ')} ({m.occurrence_count} occurrences)
                      </li>
                    ))}
                  </ul>
                </div>
              )}
            </div>
          ) : loading ? (
            <div className="h-16 bg-[#141420] rounded-xl animate-pulse" />
          ) : null}

          {/* Recommended Problem */}
          {node.recommended_problem_name && (
            <div className="p-4 rounded-xl bg-[#12121E] border border-[#EF4444]/30">
              <span className="text-[10px] font-mono uppercase tracking-wider text-[#EF4444] font-semibold block mb-1">
                Suggested Next Target
              </span>
              <p className="text-sm font-semibold text-[#0F172A]">
                {node.recommended_problem_name}
              </p>
            </div>
          )}
        </div>

        {/* Footer Actions */}
        <div className="pt-4 border-t border-[#E2E8F0] flex items-center justify-between gap-3">
          <button
            onClick={onClose}
            className="px-4 py-2 rounded-xl text-xs font-medium text-[#64748B] hover:text-[#0F172A] transition-colors"
          >
            Close
          </button>

          <button
            disabled={isLocked}
            onClick={handleStartPractice}
            className={`inline-flex items-center gap-2 px-5 py-2.5 rounded-xl text-xs font-semibold shadow-lg transition-all ${
              isLocked
                ? 'bg-[#F1F5F9] text-zinc-500 cursor-not-allowed'
                : 'bg-gradient-to-r from-[#EF4444] to-[#7C3AED] hover:from-[#7C3AED] hover:to-[#6D28D9] text-white shadow-[#8B5CF6]/25'
            }`}
          >
            {isLocked ? (
              <>
                <Lock className="w-3.5 h-3.5" />
                <span>Complete Prerequisites First</span>
              </>
            ) : (
              <>
                <span>Practice {node.concept_name}</span>
                <ArrowRight className="w-3.5 h-3.5" />
              </>
            )}
          </button>
        </div>
      </div>
    </div>
  );
};
