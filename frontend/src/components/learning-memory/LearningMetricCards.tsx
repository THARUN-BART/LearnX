import React from 'react';
import { Brain, CheckCircle, HelpCircle, AlertOctagon, RefreshCw, Calendar } from 'lucide-react';
import { LearningMemorySummary } from '../../types/learning-memory.types';

interface Props {
  summary: LearningMemorySummary;
}

export const LearningMetricCards: React.FC<Props> = ({ summary }) => {
  const cards = [
    {
      label: 'Concepts Practiced',
      value: summary.total_concepts_practised,
      subtext: `${summary.mastered_concepts} mastered, ${summary.struggling_concepts} focusing`,
      icon: Brain,
      gradient: 'from-[#8B5CF6]/20 to-[#6366F1]/10',
      border: 'border-[#8B5CF6]/30',
      textColor: 'text-[#8B5CF6]',
    },
    {
      label: 'Independent Solves',
      value: summary.independent_solves,
      subtext: 'Solved unassisted from start to finish',
      icon: CheckCircle,
      gradient: 'from-emerald-500/20 to-teal-500/10',
      border: 'border-emerald-500/30',
      textColor: 'text-emerald-400',
    },
    {
      label: 'Assisted Solves',
      value: summary.assisted_solves,
      subtext: 'Completed with progressive AI hints',
      icon: HelpCircle,
      gradient: 'from-blue-500/20 to-indigo-500/10',
      border: 'border-blue-500/30',
      textColor: 'text-blue-400',
    },
    {
      label: 'Recurring Mistakes',
      value: summary.total_mistakes,
      subtext: `${summary.active_weaknesses.length} active patterns tracked`,
      icon: AlertOctagon,
      gradient: 'from-rose-500/20 to-orange-500/10',
      border: 'border-rose-500/30',
      textColor: 'text-rose-400',
    },
    {
      label: 'Corrected Mistakes',
      value: summary.corrected_mistakes,
      subtext: 'Verified by subsequent successful solves',
      icon: RefreshCw,
      gradient: 'from-amber-500/20 to-yellow-500/10',
      border: 'border-amber-500/30',
      textColor: 'text-amber-400',
    },
    {
      label: 'Learning Days',
      value: summary.active_learning_days,
      subtext: 'Days with cognitive activity recorded',
      icon: Calendar,
      gradient: 'from-purple-500/20 to-pink-500/10',
      border: 'border-purple-500/30',
      textColor: 'text-purple-400',
    },
  ];

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      {cards.map((c, i) => {
        const Icon = c.icon;
        return (
          <div
            key={i}
            className={`p-5 rounded-2xl bg-[#0D0D14] border ${c.border} bg-gradient-to-br ${c.gradient} flex flex-col justify-between transition-all duration-200 hover:scale-[1.01]`}
          >
            <div className="flex items-center justify-between mb-3">
              <span className="text-xs font-mono uppercase tracking-wider text-[#8E8EA0]">
                {c.label}
              </span>
              <div className={`p-2 rounded-xl bg-[#161622] ${c.textColor}`}>
                <Icon className="w-5 h-5" />
              </div>
            </div>
            <div>
              <div className="text-3xl font-extrabold text-[#F4F4F6] font-mono tracking-tight">
                {c.value}
              </div>
              <p className="text-xs text-[#8E8EA0] mt-1">{c.subtext}</p>
            </div>
          </div>
        );
      })}
    </div>
  );
};
