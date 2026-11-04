import React from 'react';
import {
  Sparkles,
  Target,
  AlertTriangle,
  Layers,
  Flame,
  RotateCcw,
  Zap
} from 'lucide-react';
import { RecommendationCategory } from '../../types/recommendations.types';

interface TabItem {
  id: string;
  label: string;
  icon: React.ReactNode;
  count?: number;
}

interface RecommendationCategoryTabsProps {
  activeTab: string;
  onTabChange: (tabId: string) => void;
  counts: Record<string, number>;
}

export const RecommendationCategoryTabs: React.FC<RecommendationCategoryTabsProps> = ({
  activeTab,
  onTabChange,
  counts
}) => {
  const tabs: TabItem[] = [
    {
      id: 'all',
      label: 'All Picks',
      icon: <Sparkles className="w-3.5 h-3.5" />,
      count: Object.values(counts).reduce((a, b) => a + b, 0)
    },
    {
      id: 'recommended_for_you',
      label: 'Next for You',
      icon: <Target className="w-3.5 h-3.5" />,
      count: counts['recommended_for_you'] || 0
    },
    {
      id: 'improve_weaknesses',
      label: 'Weakness Targets',
      icon: <AlertTriangle className="w-3.5 h-3.5" />,
      count: counts['improve_weaknesses'] || 0
    },
    {
      id: 'strengthen_fundamentals',
      label: 'Fundamentals',
      icon: <Layers className="w-3.5 h-3.5" />,
      count: counts['strengthen_fundamentals'] || 0
    },
    {
      id: 'challenge_yourself',
      label: 'Challenges',
      icon: <Flame className="w-3.5 h-3.5" />,
      count: counts['challenge_yourself'] || 0
    },
    {
      id: 'revision_practice',
      label: 'Revision',
      icon: <RotateCcw className="w-3.5 h-3.5" />,
      count: counts['revision_practice'] || 0
    }
  ];

  return (
    <div className="flex items-center gap-1.5 overflow-x-auto pb-1 select-none no-scrollbar">
      {tabs.map((tab) => {
        const isActive = activeTab === tab.id;
        return (
          <button
            key={tab.id}
            onClick={() => onTabChange(tab.id)}
            className={`flex items-center gap-2 px-3.5 py-1.5 rounded-xl font-mono text-xs whitespace-nowrap transition-all duration-150 cursor-pointer ${
              isActive
                ? 'bg-[#F1F5F9] text-[#DC2626] border border-[#EF4444]/40 shadow-sm shadow-[#EF4444]/10'
                : 'bg-white text-[#64748B] hover:text-[#0F172A] hover:bg-[#F8FAFC] border border-[#E2E8F0]'
            }`}
          >
            <span className={isActive ? 'text-[#EF4444]' : 'text-[#94A3B8]'}>{tab.icon}</span>
            <span>{tab.label}</span>
            {tab.count !== undefined && tab.count > 0 && (
              <span
                className={`text-[10px] px-1.5 py-0.2 rounded-full font-bold ${
                  isActive ? 'bg-[#EF4444]/20 text-[#DC2626]' : 'bg-[#E2E8F0] text-[#94A3B8]'
                }`}
              >
                {tab.count}
              </span>
            )}
          </button>
        );
      })}
    </div>
  );
};

export default RecommendationCategoryTabs;

