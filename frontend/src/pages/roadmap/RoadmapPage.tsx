import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  CheckCircle2,
  Play,
  Clock,
  Code2,
  ChevronRight,
  Filter,
  Sparkles
} from 'lucide-react';
import { Topic, Difficulty } from '../../types/learning.types';
import roadmapService from '../../services/roadmap.service';
import Card from '../../components/ui/Card';
import Badge from '../../components/ui/Badge';
import { PageSkeleton } from '../../components/ui/Skeleton';

export const RoadmapPage: React.FC = () => {
  const navigate = useNavigate();
  const [topics, setTopics] = useState<Topic[]>([]);
  const [filterDifficulty, setFilterDifficulty] = useState<string>('all');
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const fetchTopics = async () => {
      const data = await roadmapService.getRoadmap();
      setTopics(data);
      setLoading(false);
    };
    fetchTopics();
  }, []);

  const filteredTopics = topics.filter((t) => {
    if (filterDifficulty === 'all') return true;
    return t.difficulty.toLowerCase() === filterDifficulty.toLowerCase();
  });

  const completedCount = topics.filter((t) => t.status === 'completed').length;
  const inProgressCount = topics.filter((t) => t.status === 'in_progress').length;

  if (loading) {
    return <PageSkeleton rows={3} />;
  }

  return (
    <div className="space-y-8">
      {/* Header & Stats Banner */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl sm:text-3xl font-bold tracking-tight text-[#F5F5F7] font-mono">
            DSA Mastery Roadmap
          </h1>
          <p className="text-xs text-[#9292A3] mt-1">
            Official Striver A2Z sheet — {topics.length} pillars with live progress from your account.
          </p>
        </div>

        {/* Milestone Summary Pill */}
        <div className="flex items-center gap-3 bg-[#101019] p-2 rounded-xl border border-[#232336] text-xs font-mono">
          <div className="px-3 py-1 bg-emerald-500/10 text-emerald-400 rounded-lg border border-emerald-500/20">
            {completedCount} Completed
          </div>
          <div className="px-3 py-1 bg-[#8B5CF6]/10 text-[#A78BFA] rounded-lg border border-[#8B5CF6]/20">
            {inProgressCount} In Progress
          </div>
          <button
            onClick={() => navigate('/practice', { state: { filterStatus: 'Needs Revision' } })}
            className="px-3 py-1 bg-amber-500/10 text-amber-400 hover:bg-amber-500/20 rounded-lg border border-amber-500/20 transition-colors"
          >
            Revision Queue
          </button>
        </div>
      </div>

      {/* Filter toolbar */}
      <div className="flex items-center justify-between bg-[#171722] p-3 rounded-xl border border-[#232336]">
        <div className="flex items-center gap-2">
          <Filter className="w-4 h-4 text-[#9292A3]" />
          <span className="text-xs text-[#9292A3] font-medium hidden sm:inline">Filter Difficulty:</span>
          {['all', 'easy', 'medium', 'hard'].map((diff) => (
            <button
              key={diff}
              onClick={() => setFilterDifficulty(diff)}
              className={`px-3 py-1 rounded-lg text-xs font-mono capitalize transition-all ${
                filterDifficulty === diff
                  ? 'bg-[#8B5CF6] text-white shadow-sm'
                  : 'bg-[#101019] text-[#9292A3] hover:text-white border border-[#232336]'
              }`}
            >
              {diff}
            </button>
          ))}
        </div>
        <span className="text-xs font-mono text-[#65657B]">
          Showing {filteredTopics.length} of {topics.length} Pillars
        </span>
      </div>

      {/* Topics Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {filteredTopics.map((topic, index) => {
          const isCompleted = topic.status === 'completed';
          const isInProgress = topic.status === 'in_progress' || !isCompleted;
          const totalProbs = topic.totalProblems || topic.problem_count || 1;
          const completedProbs = topic.completedProblems || topic.solved_count || 0;
          const percent = Math.min(100, Math.round((completedProbs / totalProbs) * 100));

          return (
            <Card
              key={topic.id}
              onClick={() => {
                navigate(`/roadmap/${topic.slug || topic.id}`);
              }}
              className={`p-6 flex flex-col justify-between transition-all duration-200 hover:border-[#8B5CF6]/50 hover:bg-[#1A1A28] cursor-pointer ${
                isInProgress ? 'border-[#8B5CF6]/40 glow-purple' : ''
              }`}
            >
              <div className="space-y-3">
                {/* Status pill & Difficulty */}
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <span className="text-xs font-mono text-[#65657B] font-bold">
                      #{String(topic.topic_order || index + 1).padStart(2, '0')}
                    </span>
                    <Badge
                      variant={
                        topic.difficulty === 'Easy'
                          ? 'easy'
                          : topic.difficulty === 'Medium'
                          ? 'medium'
                          : 'hard'
                      }
                      size="sm"
                    >
                      {topic.difficulty}
                    </Badge>
                  </div>

                  {isCompleted && (
                    <Badge variant="easy" size="sm" icon={<CheckCircle2 className="w-3 h-3" />}>
                      Done
                    </Badge>
                  )}
                  {!isCompleted && (
                    <Badge variant="purple" size="sm" icon={<Sparkles className="w-3 h-3" />}>
                      In Progress
                    </Badge>
                  )}
                </div>

                {/* Title & Description */}
                <div>
                  <h3 className="text-base font-bold text-[#F5F5F7] font-mono group-hover:text-[#A78BFA] transition-colors">
                    {topic.title}
                  </h3>
                  <p className="text-xs text-[#9292A3] mt-1.5 leading-relaxed line-clamp-2">
                    {topic.description}
                  </p>
                </div>

                {/* Subtopic tags */}
                <div className="flex flex-wrap gap-1.5 pt-1">
                  {(topic.subtopics || []).slice(0, 2).map((sub, i) => (
                    <span
                      key={i}
                      className="text-[10px] font-mono px-2 py-0.5 rounded bg-[#101019] text-[#9292A3] border border-[#232336]"
                    >
                      {sub}
                    </span>
                  ))}
                  {(topic.subtopics || []).length > 2 && (
                    <span className="text-[10px] font-mono text-[#65657B] px-1 py-0.5">
                      +{(topic.subtopics || []).length - 2}
                    </span>
                  )}
                </div>
              </div>

              {/* Progress Bar & Footer */}
              <div className="pt-5 mt-4 border-t border-[#232336] space-y-2">
                <div className="flex justify-between text-xs font-mono">
                  <span className="text-[#9292A3]">
                    {completedProbs} / {totalProbs} Solved
                  </span>
                  <span className="text-[#F5F5F7] font-semibold">{percent}%</span>
                </div>
                <div className="w-full bg-[#101019] h-1.5 rounded-full overflow-hidden">
                  <div
                    className={`h-full transition-all duration-300 ${
                      isCompleted ? 'bg-emerald-500' : 'bg-[#8B5CF6]'
                    }`}
                    style={{ width: `${percent}%` }}
                  />
                </div>
              </div>
            </Card>
          );
        })}
      </div>
    </div>
  );
};

export default RoadmapPage;

