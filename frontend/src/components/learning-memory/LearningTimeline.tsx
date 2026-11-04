import React from 'react';
import {
  CheckCircle2,
  XCircle,
  HelpCircle,
  AlertTriangle,
  RotateCcw,
  BookOpen,
} from 'lucide-react';
import { LearningEventItem } from '../../types/learning-memory.types';

interface Props {
  events: LearningEventItem[];
  loading?: boolean;
}

export const LearningTimeline: React.FC<Props> = ({ events, loading }) => {
  if (loading) {
    return (
      <div className="space-y-3">
        {[1, 2, 3].map((n) => (
          <div key={n} className="h-12 bg-[#F1F5F9] rounded-xl animate-pulse" />
        ))}
      </div>
    );
  }

  if (!events || events.length === 0) {
    return (
      <p className="text-xs text-[#64748B] italic py-4 text-center">
        No learning activity events logged yet.
      </p>
    );
  }

  const getEventIcon = (type: string) => {
    switch (type) {
      case 'submission_accepted':
        return <CheckCircle2 className="w-4 h-4 text-emerald-400" />;
      case 'submission_failed':
        return <XCircle className="w-4 h-4 text-rose-400" />;
      case 'hint_requested':
        return <HelpCircle className="w-4 h-4 text-blue-400" />;
      case 'mistake_detected':
        return <AlertTriangle className="w-4 h-4 text-amber-400" />;
      case 'mistake_corrected':
        return <RotateCcw className="w-4 h-4 text-teal-400" />;
      default:
        return <BookOpen className="w-4 h-4 text-purple-400" />;
    }
  };

  const getEventLabel = (e: LearningEventItem) => {
    switch (e.event_type) {
      case 'submission_accepted':
        return e.event_data?.is_independent
          ? 'Solved problem independently'
          : 'Solved problem with AI hint assistance';
      case 'submission_failed':
        return 'Code execution or assertion failed';
      case 'hint_requested':
        return `Requested Level ${e.event_data?.hint_level || 1} progressive hint`;
      case 'mistake_detected':
        return `Detected mistake pattern: ${e.event_data?.category || 'Algorithm'} (${e.event_data?.subtype || 'Edge case'})`;
      case 'mistake_corrected':
        return 'Successfully corrected previous mistake pattern';
      case 'walkthrough_requested':
        return 'Reviewed code algorithmic walkthrough';
      default:
        return e.event_type.replace('_', ' ');
    }
  };

  return (
    <div className="relative pl-6 space-y-4 before:content-[''] before:absolute before:left-2.5 before:top-2 before:bottom-2 before:w-0.5 before:bg-[#F1F5F9]">
      {events.map((e) => (
        <div key={e.id} className="relative flex items-start gap-3 text-xs">
          <div className="absolute -left-6 mt-0.5 w-5 h-5 rounded-full bg-[#F1F5F9] border border-[#E2E8F0] flex items-center justify-center">
            {getEventIcon(e.event_type)}
          </div>
          <div className="flex-1 bg-[#F1F5F9] p-3 rounded-xl border border-[#E2E8F0]">
            <div className="flex items-center justify-between">
              <span className="font-medium text-[#0F172A]">{getEventLabel(e)}</span>
              <span className="text-[10px] text-[#94A3B8] font-mono">
                {new Date(e.occurred_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
              </span>
            </div>
            {e.concept_key && (
              <span className="inline-block mt-1 px-2 py-0.5 rounded text-[10px] font-mono bg-[#161622] text-[#EF4444]">
                {e.concept_key}
              </span>
            )}
          </div>
        </div>
      ))}
    </div>
  );
};
