import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  X,
  History,
  CheckCircle2,
  Clock,
  ThumbsUp,
  Zap,
  Flame,
  ArrowRight,
  ExternalLink
} from 'lucide-react';
import { ProblemRecommendation } from '../../types/recommendations.types';
import recommendationService from '../../services/recommendation.service';
import Badge from '../ui/Badge';
import Button from '../ui/Button';

interface RecommendationHistoryDrawerProps {
  isOpen: boolean;
  onClose: () => void;
}

export const RecommendationHistoryDrawer: React.FC<RecommendationHistoryDrawerProps> = ({
  isOpen,
  onClose
}) => {
  const navigate = useNavigate();
  const [history, setHistory] = useState<ProblemRecommendation[]>([]);
  const [total, setTotal] = useState<number>(0);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    if (isOpen) {
      const loadHistory = async () => {
        setLoading(true);
        const res = await recommendationService.getHistory(30, 0);
        setHistory(res.history || []);
        setTotal(res.total || 0);
        setLoading(false);
      };
      loadHistory();
    }
  }, [isOpen]);

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex justify-end">
      {/* Backdrop */}
      <div className="fixed inset-0 bg-black/70 backdrop-blur-xs transition-opacity" onClick={onClose} />

      {/* Drawer content */}
      <div className="relative w-full max-w-md bg-white border-l border-[#E2E8F0] h-full flex flex-col z-10 animate-in slide-in-from-right duration-200">
        {/* Header */}
        <div className="h-16 px-6 border-b border-[#E2E8F0] flex items-center justify-between shrink-0">
          <div className="flex items-center gap-2 font-mono">
            <History className="w-4 h-4 text-[#EF4444]" />
            <h3 className="text-sm font-bold text-[#0F172A]">Recommendation Log</h3>
            <span className="text-[10px] text-[#94A3B8] px-2 py-0.5 rounded bg-[#F1F5F9] border border-[#E2E8F0]">
              {total} Total
            </span>
          </div>

          <button
            onClick={onClose}
            className="p-1.5 rounded-lg text-[#64748B] hover:text-[#0F172A] hover:bg-[#F8FAFC] transition-colors"
          >
            <X className="w-4 h-4" />
          </button>
        </div>

        {/* List of items */}
        <div className="flex-1 overflow-y-auto p-6 space-y-4 font-mono">
          {loading ? (
            <div className="py-20 text-center text-xs text-[#94A3B8]">
              Loading historical log...
            </div>
          ) : history.length === 0 ? (
            <div className="py-20 text-center space-y-2">
              <History className="w-8 h-8 text-[#CBD5E1] mx-auto" />
              <p className="text-xs text-[#94A3B8]">No recommendation history logged yet.</p>
              <p className="text-[11px] text-[#4A4A5E]">
                As you receive and practice problem picks, your history will populate here.
              </p>
            </div>
          ) : (
            history.map((rec) => {
              const isCompleted = rec.status === 'completed';
              return (
                <div
                  key={rec.id}
                  className="p-4 rounded-xl bg-[#F1F5F9] border border-[#E2E8F0] space-y-2 hover:border-[#CBD5E1] transition-all"
                >
                  <div className="flex items-center justify-between text-xs">
                    <span className="text-[10px] uppercase tracking-wider text-[#DC2626] bg-[#EF4444]/15 px-2 py-0.5 rounded border border-[#EF4444]/20">
                      {rec.category.replace(/_/g, ' ')}
                    </span>

                    <span className="text-[10px] text-[#94A3B8]">
                      {new Date(rec.created_at).toLocaleDateString()}
                    </span>
                  </div>

                  <div>
                    <h4
                      onClick={() => {
                        onClose();
                        navigate(`/practice/${rec.problem_id || rec.problem_slug}`);
                      }}
                      className="text-sm font-bold text-[#0F172A] hover:text-[#DC2626] cursor-pointer transition-colors"
                    >
                      {rec.problem_title || 'DSA Problem'}
                    </h4>
                    <p className="text-[11px] text-[#64748B] line-clamp-2 mt-1">
                      {rec.headline}
                    </p>
                  </div>

                  <div className="flex items-center justify-between pt-2 border-t border-[#E2E8F0] text-xs">
                    <div className="flex items-center gap-2">
                      <span
                        className={`text-[10px] px-2 py-0.5 rounded flex items-center gap-1 ${
                          isCompleted
                            ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20'
                            : 'bg-[#E2E8F0] text-[#94A3B8]'
                        }`}
                      >
                        {isCompleted ? (
                          <>
                            <CheckCircle2 className="w-3 h-3 text-emerald-400" />
                            Completed
                          </>
                        ) : (
                          'Suggested'
                        )}
                      </span>

                      {rec.feedback_rating && (
                        <span className="text-[10px] text-amber-400 font-mono">
                          Rated: {rec.feedback_rating}
                        </span>
                      )}
                    </div>

                    <button
                      onClick={() => {
                        onClose();
                        navigate(`/practice/${rec.problem_id || rec.problem_slug}`);
                      }}
                      className="text-xs text-[#EF4444] hover:text-[#DC2626] flex items-center gap-1"
                    >
                      <span>Open</span>
                      <ArrowRight className="w-3 h-3" />
                    </button>
                  </div>
                </div>
              );
            })
          )}
        </div>
      </div>
    </div>
  );
};

export default RecommendationHistoryDrawer;

