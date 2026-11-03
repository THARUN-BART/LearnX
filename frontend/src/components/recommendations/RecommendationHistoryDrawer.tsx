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
      <div className="relative w-full max-w-md bg-[#0D0D14] border-l border-[#1C1C28] h-full flex flex-col z-10 animate-in slide-in-from-right duration-200">
        {/* Header */}
        <div className="h-16 px-6 border-b border-[#1C1C28] flex items-center justify-between shrink-0">
          <div className="flex items-center gap-2 font-mono">
            <History className="w-4 h-4 text-[#8B5CF6]" />
            <h3 className="text-sm font-bold text-[#F4F4F6]">Recommendation Log</h3>
            <span className="text-[10px] text-[#5A5A6E] px-2 py-0.5 rounded bg-[#161622] border border-[#232336]">
              {total} Total
            </span>
          </div>

          <button
            onClick={onClose}
            className="p-1.5 rounded-lg text-[#8E8EA0] hover:text-[#F4F4F6] hover:bg-[#14141E] transition-colors"
          >
            <X className="w-4 h-4" />
          </button>
        </div>

        {/* List of items */}
        <div className="flex-1 overflow-y-auto p-6 space-y-4 font-mono">
          {loading ? (
            <div className="py-20 text-center text-xs text-[#5A5A6E]">
              Loading historical log...
            </div>
          ) : history.length === 0 ? (
            <div className="py-20 text-center space-y-2">
              <History className="w-8 h-8 text-[#2B2B3C] mx-auto" />
              <p className="text-xs text-[#5A5A6E]">No recommendation history logged yet.</p>
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
                  className="p-4 rounded-xl bg-[#12121A] border border-[#1C1C28] space-y-2 hover:border-[#2B2B3C] transition-all"
                >
                  <div className="flex items-center justify-between text-xs">
                    <span className="text-[10px] uppercase tracking-wider text-[#A78BFA] bg-[#8B5CF6]/15 px-2 py-0.5 rounded border border-[#8B5CF6]/20">
                      {rec.category.replace(/_/g, ' ')}
                    </span>

                    <span className="text-[10px] text-[#5A5A6E]">
                      {new Date(rec.created_at).toLocaleDateString()}
                    </span>
                  </div>

                  <div>
                    <h4
                      onClick={() => {
                        onClose();
                        navigate(`/practice/${rec.problem_id || rec.problem_slug}`);
                      }}
                      className="text-sm font-bold text-[#F4F4F6] hover:text-[#A78BFA] cursor-pointer transition-colors"
                    >
                      {rec.problem_title || 'DSA Problem'}
                    </h4>
                    <p className="text-[11px] text-[#8E8EA0] line-clamp-2 mt-1">
                      {rec.headline}
                    </p>
                  </div>

                  <div className="flex items-center justify-between pt-2 border-t border-[#1C1C28] text-xs">
                    <div className="flex items-center gap-2">
                      <span
                        className={`text-[10px] px-2 py-0.5 rounded flex items-center gap-1 ${
                          isCompleted
                            ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20'
                            : 'bg-[#181824] text-[#5A5A6E]'
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
                      className="text-xs text-[#8B5CF6] hover:text-[#A78BFA] flex items-center gap-1"
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

