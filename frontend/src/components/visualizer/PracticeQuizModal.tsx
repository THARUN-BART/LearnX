import React, { useState, useEffect } from 'react';
import { PracticeQuestion } from '../../types/visualizer.types';
import visualizerService from '../../services/visualizer.service';
import { X, CheckCircle2, XCircle, BrainCircuit } from 'lucide-react';
import Button from '../ui/Button';

interface PracticeQuizModalProps {
  isOpen: boolean;
  onClose: () => void;
  algorithmId: string;
  algorithmName: string;
}

export const PracticeQuizModal: React.FC<PracticeQuizModalProps> = ({
  isOpen,
  onClose,
  algorithmId,
  algorithmName,
}) => {
  const [questions, setQuestions] = useState<PracticeQuestion[]>([]);
  const [loading, setLoading] = useState(false);
  const [selectedAnswers, setSelectedAnswers] = useState<Record<number, number>>({});
  const [submitted, setSubmitted] = useState<Record<number, boolean>>({});

  useEffect(() => {
    if (isOpen && algorithmId) {
      setLoading(true);
      setSelectedAnswers({});
      setSubmitted({});
      visualizerService.getPracticeQuestions(algorithmId).then((res) => {
        setQuestions(res);
        setLoading(false);
      });
    }
  }, [isOpen, algorithmId]);

  if (!isOpen) return null;

  const handleSelectOption = (qIdx: number, optIdx: number) => {
    if (submitted[qIdx]) return;
    setSelectedAnswers((prev) => ({ ...prev, [qIdx]: optIdx }));
  };

  const handleSubmitQuestion = (qIdx: number) => {
    setSubmitted((prev) => ({ ...prev, [qIdx]: true }));
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm">
      <div className="w-full max-w-xl bg-white border border-[#E2E8F0] rounded-2xl p-6 font-mono max-h-[85vh] flex flex-col">
        {/* Header */}
        <div className="flex items-center justify-between border-b border-[#E2E8F0] pb-3 mb-4">
          <div className="flex items-center gap-2">
            <BrainCircuit className="w-5 h-5 text-[#EF4444]" />
            <div>
              <h3 className="text-base font-bold text-[#0F172A] font-sans">
                {algorithmName} Practice Quiz
              </h3>
              <span className="text-[10px] text-[#64748B]">
                Conceptual & Complexity Calibration
              </span>
            </div>
          </div>
          <button onClick={onClose} className="p-1.5 text-[#64748B] hover:text-[#0F172A] rounded-lg">
            <X className="w-4 h-4" />
          </button>
        </div>

        {/* Content */}
        <div className="flex-1 overflow-y-auto space-y-6 pr-1">
          {loading ? (
            <div className="py-12 text-center text-xs text-[#94A3B8]">
              Loading algorithm practice questions…
            </div>
          ) : questions.length === 0 ? (
            <div className="py-12 text-center text-xs text-[#64748B]">
              No practice questions registered for this algorithm yet.
            </div>
          ) : (
            questions.map((q, qIdx) => {
              const userPick = selectedAnswers[qIdx];
              const isSubmitted = submitted[qIdx];
              const isCorrect = userPick === q.correct_option_index;

              return (
                <div key={q.id || qIdx} className="p-4 bg-[#F1F5F9] rounded-xl border border-[#E2E8F0] space-y-3">
                  <div className="flex items-center justify-between gap-2">
                    <span className="text-[10px] uppercase font-bold text-[#DC2626]">
                      Question {qIdx + 1} · {q.difficulty}
                    </span>
                    {isSubmitted && (
                      <span className={`text-[10px] font-bold flex items-center gap-1 ${
                        isCorrect ? 'text-emerald-400' : 'text-rose-400'
                      }`}>
                        {isCorrect ? <CheckCircle2 className="w-3.5 h-3.5" /> : <XCircle className="w-3.5 h-3.5" />}
                        {isCorrect ? 'Correct!' : 'Incorrect'}
                      </span>
                    )}
                  </div>

                  <p className="text-xs text-[#0F172A] font-sans font-medium">{q.question}</p>

                  <div className="space-y-1.5">
                    {q.options.map((opt, optIdx) => {
                      const isSelected = userPick === optIdx;
                      let btnStyle = 'bg-[#E2E8F0] border-[#E2E8F0] text-[#CBD5E1] hover:border-[#EF4444]/50';

                      if (isSubmitted) {
                        if (optIdx === q.correct_option_index) {
                          btnStyle = 'bg-emerald-500/20 border-emerald-500 text-emerald-300 font-bold';
                        } else if (isSelected) {
                          btnStyle = 'bg-rose-500/20 border-rose-500 text-rose-300';
                        } else {
                          btnStyle = 'opacity-40 border-[#E2E8F0] text-[#94A3B8]';
                        }
                      } else if (isSelected) {
                        btnStyle = 'bg-[#EF4444]/20 border-[#EF4444] text-white font-bold';
                      }

                      return (
                        <button
                          key={optIdx}
                          disabled={isSubmitted}
                          onClick={() => handleSelectOption(qIdx, optIdx)}
                          className={`w-full text-left p-2.5 rounded-lg border text-xs font-sans transition-all flex items-center justify-between ${btnStyle}`}
                        >
                          <span>{opt}</span>
                        </button>
                      );
                    })}
                  </div>

                  {!isSubmitted ? (
                    <Button
                      size="sm"
                      disabled={userPick === undefined}
                      onClick={() => handleSubmitQuestion(qIdx)}
                      className="text-xs"
                    >
                      Check Answer
                    </Button>
                  ) : (
                    <div className="p-3 bg-white rounded-lg border border-[#E2E8F0] text-[11px] text-[#DC2626] font-sans">
                      <span className="font-bold text-white">Explanation: </span>
                      {q.explanation}
                    </div>
                  )}
                </div>
              );
            })
          )}
        </div>
      </div>
    </div>
  );
};
export default PracticeQuizModal;

