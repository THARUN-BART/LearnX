import React, { useEffect, useState } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import {
  ArrowLeft,
  BookOpen,
  CheckCircle2,
  Clock,
  Code2,
  Play,
  Sparkles,
  Layers,
  ChevronRight
} from 'lucide-react';
import roadmapService from '../../services/roadmap.service';
import learningService from '../../services/learning.service';
import { Topic, Problem } from '../../types/learning.types';
import Card from '../../components/ui/Card';
import Badge from '../../components/ui/Badge';
import Button from '../../components/ui/Button';

export const TopicDetailPage: React.FC = () => {
  const { topicId } = useParams<{ topicId: string }>();
  const navigate = useNavigate();

  const [topic, setTopic] = useState<Topic | null>(null);
  const [problems, setProblems] = useState<Problem[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const loadTopic = async () => {
      if (!topicId) return;
      const data = await roadmapService.getTopicById(topicId);
      setTopic(data);

      const probList = await learningService.getProblems(data?.id || topicId);
      setProblems(Array.isArray(probList) ? probList : []);
      setLoading(false);
    };
    loadTopic();
  }, [topicId]);

  if (loading || !topic) {
    return (
      <div className="py-20 text-center text-xs font-mono text-[#65657B]">
        Loading curriculum topic details...
      </div>
    );
  }

  const totalProbs = topic.totalProblems || topic.problem_count || 1;
  const completedProbs = topic.completedProblems || topic.solved_count || 0;
  const percent = Math.min(100, Math.round((completedProbs / totalProbs) * 100));
  const prereqs = topic.prerequisites || [];
  const subtopicsList = topic.subtopics || [];

  return (
    <div className="space-y-8 max-w-5xl mx-auto">
      {/* Top back navigation */}
      <div>
        <Link
          to="/roadmap"
          className="inline-flex items-center gap-1.5 text-xs text-[#9292A3] hover:text-[#F5F5F7] transition-colors font-mono"
        >
          <ArrowLeft className="w-3.5 h-3.5" />
          <span>Back to Roadmap</span>
        </Link>
      </div>

      {/* Hero Header */}
      <div className="p-6 sm:p-8 rounded-2xl bg-gradient-to-r from-[#171722] via-[#1A1A28] to-[#101019] border border-[#232336] space-y-4">
        <div className="flex flex-wrap items-center justify-between gap-3">
          <div className="flex items-center gap-2">
            <Badge
              variant={
                topic.difficulty === 'Easy'
                  ? 'easy'
                  : topic.difficulty === 'Medium'
                  ? 'medium'
                  : 'hard'
              }
            >
              {topic.difficulty}
            </Badge>
            <Badge variant="purple" icon={<Clock className="w-3 h-3" />}>
              ~{topic.estimatedHours} Hours
            </Badge>
          </div>
          <span className="text-xs font-mono text-[#A78BFA] bg-[#8B5CF6]/10 px-3 py-1 rounded-full border border-[#8B5CF6]/20">
            {completedProbs} / {totalProbs} Solved ({percent}%)
          </span>
        </div>

        <h1 className="text-2xl sm:text-4xl font-bold tracking-tight text-[#F5F5F7] font-mono">
          {topic.title}
        </h1>

        <p className="text-sm text-[#9292A3] leading-relaxed max-w-3xl">
          {topic.description}
        </p>

        {/* Prerequisites */}
        {prereqs.length > 0 && (
          <div className="flex items-center gap-2 pt-2 text-xs text-[#65657B]">
            <span className="font-semibold text-[#9292A3]">Prerequisites:</span>
            {prereqs.map((p, i) => (
              <span key={i} className="px-2 py-0.5 rounded bg-[#101019] text-[#9292A3] border border-[#232336] font-mono text-[11px]">
                {p}
              </span>
            ))}
          </div>
        )}
      </div>

      {/* Algorithmic Mental Model & Subtopics */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
        <div className="lg:col-span-8 space-y-6">
          <Card className="p-6 space-y-4">
            <div className="flex items-center gap-2 text-[#A78BFA] font-mono text-xs font-semibold">
              <BookOpen className="w-4 h-4 text-[#8B5CF6]" />
              <span>Core Conceptual Mechanics</span>
            </div>
            <p className="text-sm text-[#F5F5F7] leading-relaxed font-sans">
              {topic.explanation}
            </p>
          </Card>

          {/* Practice Problems Section */}
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <h3 className="text-base font-bold text-[#F5F5F7] font-mono flex items-center gap-2">
                <Code2 className="w-4 h-4 text-[#8B5CF6]" />
                <span>Pillar Practice Problems ({problems.length})</span>
              </h3>
            </div>

            <div className="space-y-3">
              {problems.map((prob) => (
                <div
                  key={prob.id}
                  onClick={() => navigate(`/practice/${prob.id}`)}
                  className="p-4 rounded-xl bg-[#171722] hover:bg-[#1C1C2A] border border-[#232336] hover:border-[#8B5CF6]/50 transition-all flex items-center justify-between cursor-pointer group"
                >
                  <div className="space-y-1">
                    <div className="text-sm font-bold text-[#F5F5F7] group-hover:text-[#A78BFA] transition-colors">
                      {prob.title}
                    </div>
                    <div className="flex items-center gap-2.5 text-xs">
                      <Badge
                        variant={prob.difficulty === 'Easy' ? 'easy' : 'medium'}
                        size="sm"
                      >
                        {prob.difficulty}
                      </Badge>
                      <span className="text-[#65657B] font-mono text-[11px]">
                        Pass Rate: {prob.accuracyRate}%
                      </span>
                    </div>
                  </div>

                  <Button
                    variant="secondary"
                    size="sm"
                    rightIcon={<Play className="w-3.5 h-3.5" />}
                    onClick={(e) => {
                      e.stopPropagation();
                      navigate(`/practice/${prob.id}`);
                    }}
                  >
                    Solve
                  </Button>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Right: Subtopics Checklist */}
        <div className="lg:col-span-4 space-y-6">
          <Card className="p-5 space-y-4">
            <h4 className="text-xs font-bold text-[#F5F5F7] uppercase tracking-wider font-mono">
              Key Subtopics Checklist
            </h4>
            <div className="space-y-2.5">
              {subtopicsList.map((sub, i) => (
                <div key={i} className="flex items-center gap-2.5 text-xs text-[#9292A3]">
                  <CheckCircle2 className="w-4 h-4 text-emerald-400 shrink-0" />
                  <span>{sub}</span>
                </div>
              ))}
            </div>
          </Card>
        </div>
      </div>
    </div>
  );
};

export default TopicDetailPage;

