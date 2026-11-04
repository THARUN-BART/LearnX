import React, { useEffect, useState } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import {
  ArrowLeft,
  BookOpen,
  CheckCircle2,
  Clock,
  Code2,
  Play,
} from 'lucide-react';
import roadmapService from '../../services/roadmap.service';
import learningService from '../../services/learning.service';
import { Topic, Problem, Subtopic } from '../../types/learning.types';
import Card from '../../components/ui/Card';
import Badge from '../../components/ui/Badge';
import Button from '../../components/ui/Button';
import { PageSkeleton } from '../../components/ui/Skeleton';

export const TopicDetailPage: React.FC = () => {
  const { topicId } = useParams<{ topicId: string }>();
  const navigate = useNavigate();

  const [topic, setTopic] = useState<Topic | null>(null);
  const [subtopics, setSubtopics] = useState<Subtopic[]>([]);
  const [problems, setProblems] = useState<Problem[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const loadTopic = async () => {
      if (!topicId) return;
      setLoading(true);
      const data = await roadmapService.getTopicById(topicId);
      setTopic(data);
      if (!data) {
        setLoading(false);
        return;
      }
      const [probList, subList] = await Promise.all([
        learningService.getProblems({ topicId: data.id, limit: 80 }),
        roadmapService.getSubtopics(data.id),
      ]);
      setProblems(Array.isArray(probList) ? probList : []);
      setSubtopics(subList);
      setLoading(false);
    };
    loadTopic();
  }, [topicId]);

  if (loading || !topic) {
    return <PageSkeleton rows={3} />;
  }

  const totalProbs = topic.totalProblems || topic.problem_count || 1;
  const completedProbs = topic.completedProblems || topic.solved_count || 0;
  const percent = Math.min(100, Math.round((completedProbs / totalProbs) * 100));
  const prereqs = topic.prerequisites || [];

  const grouped = subtopics.length
    ? subtopics.map((sub) => ({
        sub,
        items: problems.filter((p) => p.subtopic_id === sub.id),
      }))
    : [{ sub: null, items: problems }];

  const ungrouped = problems.filter(
    (p) => subtopics.length && !subtopics.some((s) => s.id === p.subtopic_id)
  );

  const renderProblem = (prob: Problem) => (
    <div
      key={prob.id}
      onClick={() => navigate(`/practice/${prob.id}`)}
      className="p-4 rounded-xl bg-[#F8FAFC] hover:bg-[#1C1C2A] border border-[#E2E8F0] hover:border-[#EF4444]/50 transition-all flex items-center justify-between cursor-pointer group"
    >
      <div className="space-y-1 min-w-0">
        <div className="text-sm font-bold text-[#0F172A] group-hover:text-[#DC2626] truncate">
          {prob.title}
        </div>
        <div className="flex items-center gap-2.5 text-xs">
          <Badge
            variant={prob.difficulty === 'Easy' ? 'easy' : prob.difficulty === 'Hard' ? 'hard' : 'medium'}
            size="sm"
          >
            {prob.difficulty}
          </Badge>
          <span className="text-[#94A3B8] font-mono text-[11px]">{prob.progress_status || 'Not Started'}</span>
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
  );

  return (
    <div className="space-y-8 max-w-5xl mx-auto">
      <Link
        to="/roadmap"
        className="inline-flex items-center gap-1.5 text-xs text-[#64748B] hover:text-[#0F172A] font-mono"
      >
        <ArrowLeft className="w-3.5 h-3.5" />
        Back to Roadmap
      </Link>

      <div className="p-6 sm:p-8 rounded-2xl bg-gradient-to-r from-[#171722] via-[#1A1A28] to-[#101019] border border-[#E2E8F0] space-y-4">
        <div className="flex flex-wrap items-center justify-between gap-3">
          <div className="flex items-center gap-2">
            <Badge variant={topic.difficulty === 'Easy' ? 'easy' : topic.difficulty === 'Hard' ? 'hard' : 'medium'}>
              {topic.difficulty}
            </Badge>
            <Badge variant="purple" icon={<Clock className="w-3 h-3" />}>
              ~{topic.estimatedHours} Hours
            </Badge>
          </div>
          <span className="text-xs font-mono text-[#DC2626] bg-[#EF4444]/10 px-3 py-1 rounded-full border border-[#EF4444]/20">
            {completedProbs} / {totalProbs} Solved ({percent}%)
          </span>
        </div>
        <h1 className="text-2xl sm:text-4xl font-bold tracking-tight text-[#0F172A] font-mono">{topic.title}</h1>
        <p className="text-sm text-[#64748B] leading-relaxed max-w-3xl">{topic.description}</p>
        {prereqs.length > 0 && (
          <div className="flex items-center gap-2 pt-2 text-xs text-[#94A3B8]">
            <span className="font-semibold text-[#64748B]">Prerequisites:</span>
            {prereqs.map((p, i) => (
              <span key={i} className="px-2 py-0.5 rounded bg-white text-[#64748B] border border-[#E2E8F0] font-mono text-[11px]">
                {p}
              </span>
            ))}
          </div>
        )}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
        <div className="lg:col-span-8 space-y-6">
          <Card className="p-6 space-y-4">
            <div className="flex items-center gap-2 text-[#DC2626] font-mono text-xs font-semibold">
              <BookOpen className="w-4 h-4 text-[#EF4444]" />
              Core Conceptual Mechanics
            </div>
            <p className="text-sm text-[#0F172A] leading-relaxed">{topic.explanation}</p>
          </Card>

          <div className="space-y-6">
            <h3 className="text-base font-bold text-[#0F172A] font-mono flex items-center gap-2">
              <Code2 className="w-4 h-4 text-[#EF4444]" />
              Practice Problems ({problems.length})
            </h3>
            {grouped.map(({ sub, items }, idx) => (
              <div key={sub?.id || `group-${idx}`} className="space-y-3">
                {sub && (
                  <h4 className="text-xs font-mono text-[#DC2626] uppercase tracking-wider">
                    {sub.title} · {items.length || sub.problem_count} problems
                  </h4>
                )}
                {items.length === 0 && sub && (
                  <p className="text-xs text-[#94A3B8] font-mono">No problems loaded for this subtopic yet.</p>
                )}
                {items.map(renderProblem)}
              </div>
            ))}
            {ungrouped.length > 0 && (
              <div className="space-y-3">
                <h4 className="text-xs font-mono text-[#DC2626] uppercase tracking-wider">Other</h4>
                {ungrouped.map(renderProblem)}
              </div>
            )}
          </div>
        </div>

        <div className="lg:col-span-4 space-y-6">
          <Card className="p-5 space-y-4">
            <h4 className="text-xs font-bold text-[#0F172A] uppercase tracking-wider font-mono">
              Key Subtopics Checklist
            </h4>
            <div className="space-y-2.5">
              {(subtopics.length ? subtopics.map((s) => s.title) : topic.subtopics || []).map((sub, i) => (
                <div key={i} className="flex items-center gap-2.5 text-xs text-[#64748B]">
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
