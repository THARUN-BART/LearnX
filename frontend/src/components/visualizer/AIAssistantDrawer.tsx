import React, { useState } from 'react';
import { ExplainStepResponse, ExecutionStep } from '../../types/visualizer.types';
import visualizerService from '../../services/visualizer.service';
import { Bot, Sparkles, X, Send, Lightbulb, Compass, Key } from 'lucide-react';
import Button from '../ui/Button';
import { toast } from 'sonner';

interface AIAssistantDrawerProps {
  isOpen: boolean;
  onClose: () => void;
  algorithmId: string;
  algorithmName: string;
  currentStepData: ExecutionStep | null;
  totalSteps: number;
}

export const AIAssistantDrawer: React.FC<AIAssistantDrawerProps> = ({
  isOpen,
  onClose,
  algorithmId,
  algorithmName,
  currentStepData,
  totalSteps,
}) => {
  const [loading, setLoading] = useState(false);
  const [explanation, setExplanation] = useState<ExplainStepResponse | null>(null);
  const [question, setQuestion] = useState('');

  if (!isOpen) return null;

  const handleRequestExplanation = async (userQ?: string) => {
    if (!currentStepData) return;
    setLoading(true);
    try {
      const res = await visualizerService.explainStep({
        algorithmId,
        stepNumber: currentStepData.step_number,
        totalSteps,
        operation: currentStepData.operation,
        explanation: currentStepData.explanation,
        currentState: currentStepData.state,
        userQuestion: userQ,
      });
      if (res) {
        setExplanation(res);
      } else {
        toast.error('AI Tutor unavailable at the moment');
      }
    } catch {
      toast.error('Could not retrieve AI explanation');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex justify-end bg-black/60 backdrop-blur-sm">
      <div className="w-full max-w-md bg-[#0D0D14] border-l border-[#1C1C28] h-full flex flex-col p-5 font-mono overflow-y-auto">
        {/* Drawer Header */}
        <div className="flex items-center justify-between border-b border-[#1C1C28] pb-3 mb-4">
          <div className="flex items-center gap-2">
            <div className="w-7 h-7 rounded-lg bg-[#8B5CF6]/20 border border-[#8B5CF6]/30 flex items-center justify-center text-[#A78BFA]">
              <Bot className="w-4 h-4" />
            </div>
            <div>
              <h4 className="text-sm font-bold text-[#F4F4F6] font-sans">AI Algorithm Tutor</h4>
              <span className="text-[10px] text-[#8E8EA0]">{algorithmName}</span>
            </div>
          </div>
          <button onClick={onClose} className="p-1.5 text-[#8E8EA0] hover:text-[#F4F4F6] rounded-lg">
            <X className="w-4 h-4" />
          </button>
        </div>

        {/* Current Step Summary */}
        {currentStepData && (
          <div className="p-3 bg-[#12121B] rounded-xl border border-[#1C1C28] mb-4 space-y-1">
            <div className="flex items-center justify-between text-[10px] text-[#A78BFA] font-bold">
              <span>STEP {currentStepData.step_number} OF {totalSteps}</span>
              <span className="uppercase">{currentStepData.operation}</span>
            </div>
            <p className="text-xs text-[#E2E8F0] font-sans">{currentStepData.explanation}</p>
          </div>
        )}

        {/* Action Trigger Button */}
        {!explanation && (
          <Button
            variant="primary"
            onClick={() => handleRequestExplanation()}
            isLoading={loading}
            leftIcon={<Sparkles className="w-3.5 h-3.5" />}
            className="w-full text-xs py-2.5 mb-4"
          >
            Explain This Step With AI
          </Button>
        )}

        {/* Structured Explanation Result */}
        {explanation && (
          <div className="space-y-4 mb-4 flex-1">
            <div className="p-3.5 bg-[#12121B] rounded-xl border border-[#8B5CF6]/30 space-y-2">
              <div className="flex items-center gap-1.5 text-[#E9D5FF] text-xs font-bold font-sans">
                <Lightbulb className="w-4 h-4 text-amber-400" /> Deep Step Breakdown
              </div>
              <p className="text-xs text-[#CBD5E1] font-sans leading-relaxed">
                {explanation.explanation}
              </p>
            </div>

            <div className="p-3.5 bg-[#12121B] rounded-xl border border-[#1C1C28] space-y-1.5">
              <div className="flex items-center gap-1.5 text-[#A78BFA] text-xs font-bold font-sans">
                <Compass className="w-4 h-4 text-[#8B5CF6]" /> Algorithmic Invariant
              </div>
              <p className="text-xs text-[#CBD5E1] font-sans leading-relaxed">
                {explanation.intuition}
              </p>
            </div>

            <div className="p-3.5 bg-[#12121B] rounded-xl border border-[#1C1C28] space-y-1.5">
              <div className="flex items-center gap-1.5 text-emerald-400 text-xs font-bold font-sans">
                <Compass className="w-4 h-4 text-emerald-400" /> Next Probable Action
              </div>
              <p className="text-xs text-[#CBD5E1] font-sans leading-relaxed">
                {explanation.next_expected_move}
              </p>
            </div>

            <div className="p-3.5 bg-gradient-to-r from-[#8B5CF6]/10 to-[#6366F1]/10 rounded-xl border border-[#8B5CF6]/20 space-y-1">
              <div className="flex items-center gap-1.5 text-[#D8B4FE] text-xs font-bold font-sans">
                <Key className="w-4 h-4 text-[#A78BFA]" /> Interview Takeaway
              </div>
              <p className="text-xs text-[#F4F4F6] font-sans font-medium">
                {explanation.key_takeaway}
              </p>
            </div>
          </div>
        )}

        {/* Ask Question Input */}
        <div className="pt-3 border-t border-[#1C1C28] mt-auto">
          <form
            onSubmit={(e) => {
              e.preventDefault();
              if (question.trim()) {
                handleRequestExplanation(question.trim());
                setQuestion('');
              }
            }}
            className="flex items-center gap-2"
          >
            <input
              value={question}
              onChange={(e) => setQuestion(e.target.value)}
              placeholder="Ask AI why this step occurred..."
              className="flex-1 bg-[#12121B] text-[#F4F4F6] text-xs px-3 py-2 rounded-lg border border-[#1C1C28] outline-none font-sans"
            />
            <Button size="sm" type="submit" disabled={loading || !question.trim()}>
              <Send className="w-3.5 h-3.5" />
            </Button>
          </form>
        </div>
      </div>
    </div>
  );
};
export default AIAssistantDrawer;

