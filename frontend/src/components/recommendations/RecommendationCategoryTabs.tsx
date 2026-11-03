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
                ? 'bg-[#1C1C2B] text-[#A78BFA] border border-[#8B5CF6]/40 shadow-sm shadow-[#8B5CF6]/10'
                : 'bg-[#0D0D14] text-[#8E8EA0] hover:text-[#F4F4F6] hover:bg-[#14141E] border border-[#1C1C28]'
            }`}
          >
            <span className={isActive ? 'text-[#8B5CF6]' : 'text-[#5A5A6E]'}>{tab.icon}</span>
            <span>{tab.label}</span>
            {tab.count !== undefined && tab.count > 0 && (
              <span
                className={`text-[10px] px-1.5 py-0.2 rounded-full font-bold ${
                  isActive ? 'bg-[#8B5CF6]/20 text-[#A78BFA]' : 'bg-[#181824] text-[#5A5A6E]'
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

