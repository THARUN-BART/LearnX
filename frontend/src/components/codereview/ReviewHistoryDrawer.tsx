import React, { useState, useEffect } from 'react';
import { ReviewHistoryItem } from '../../types/codereview.types';
import codeReviewService from '../../services/codereview.service';
import { X, History, ChevronRight } from 'lucide-react';

interface ReviewHistoryDrawerProps {
  isOpen: boolean;
  onClose: () => void;
  onSelectReview: (id: string) => void;
}

export const ReviewHistoryDrawer: React.FC<ReviewHistoryDrawerProps> = ({
  isOpen,
  onClose,
  onSelectReview,
}) => {
  const [reviews, setReviews] = useState<ReviewHistoryItem[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (isOpen) {
      setLoading(true);
      codeReviewService.getStudentReviews(30, 0).then((res) => {
        setReviews(res);
        setLoading(false);
      });
    }
  }, [isOpen]);

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex justify-end bg-black/60 backdrop-blur-sm">
      <div className="w-full max-w-md bg-[#0D0D14] border-l border-[#1C1C28] h-full flex flex-col p-5 font-mono overflow-y-auto">
        {/* Header */}
        <div className="flex items-center justify-between border-b border-[#1C1C28] pb-3 mb-4">
          <div className="flex items-center gap-2">
            <History className="w-4 h-4 text-[#8B5CF6]" />
            <h4 className="text-sm font-bold text-[#F4F4F6] font-sans">Code Review History</h4>
          </div>
          <button onClick={onClose} className="p-1.5 text-[#8E8EA0] hover:text-[#F4F4F6] rounded-lg">
            <X className="w-4 h-4" />
          </button>
        </div>

        {/* List */}
        <div className="space-y-2.5 flex-1 overflow-y-auto">
          {loading ? (
            <div className="py-12 text-center text-xs text-[#5A5A6E]">Loading reviews…</div>
          ) : reviews.length === 0 ? (
            <div className="py-12 text-center text-xs text-[#8E8EA0]">
              No past code reviews found.
            </div>
          ) : (
            reviews.map((r) => (
              <button
                key={r.id}
                onClick={() => {
                  onSelectReview(r.id);
                  onClose();
                }}
                className="w-full text-left p-3.5 bg-[#12121B] hover:bg-[#181824] border border-[#1C1C28] hover:border-[#8B5CF6]/50 rounded-xl transition-all flex items-center justify-between group"
              >
                <div>
                  <div className="text-xs font-bold text-[#F4F4F6] font-sans group-hover:text-[#A78BFA] transition-colors">
                    {r.problem_title || 'General Code Review'}
                  </div>
                  <div className="text-[10px] text-[#5A5A6E] mt-1">
                    {r.language} · {r.created_at ? new Date(r.created_at).toLocaleDateString() : 'Recent'}
                  </div>
                </div>

                <div className="flex items-center gap-2 shrink-0">
                  <span
                    className={`text-xs font-bold px-2 py-0.5 rounded ${
                      r.overall_score >= 80
                        ? 'bg-emerald-500/10 text-emerald-400'
                        : r.overall_score >= 60
                        ? 'bg-amber-500/10 text-amber-400'
                        : 'bg-rose-500/10 text-rose-400'
                    }`}
                  >
                    {r.overall_score}/100
                  </span>
                  <ChevronRight className="w-4 h-4 text-[#5A5A6E] group-hover:text-white" />
                </div>
              </button>
            ))
          )}
        </div>
      </div>
    </div>
  );
};
export default ReviewHistoryDrawer;

