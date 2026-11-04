import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  User,
  GraduationCap,
  Target,
  Code2,
  Mail,
  Edit2,
  CheckCircle2,
  Save,
  Clock,
  Flame,
  Activity,
  Layers,
  ExternalLink,
  Filter,
  RefreshCw,
  Play
} from 'lucide-react';
import { useAuth } from '../../hooks/useAuth';
import { useProfile } from '../../hooks/useProfile';
import {
  COLLEGES,
  COURSES,
  DEPARTMENTS,
  ACADEMIC_YEARS,
  PROGRAMMING_LANGUAGES,
  LEARNING_GOALS,
  DSA_LEVELS
} from '../../utils/constants';
import { ProgrammingLanguage, DSALevel } from '../../types/profile.types';
import { SubmissionHistoryItem, AnalyticsSummary } from '../../types/learning.types';
import codingService from '../../services/coding.service';
import analyticsService from '../../services/analytics.service';
import Card from '../../components/ui/Card';
import Badge from '../../components/ui/Badge';
import Button from '../../components/ui/Button';
import { toast } from 'sonner';

export const ProfilePage: React.FC = () => {
  const navigate = useNavigate();
  const { user, updateLocalUser } = useAuth();
  const { profile, updateProfile } = useProfile();

  const [activeTab, setActiveTab] = useState<'profile' | 'submissions' | 'performance'>('profile');
  const [isEditing, setIsEditing] = useState<boolean>(false);
  const [isSaving, setIsSaving] = useState<boolean>(false);

  // Submissions state
  const [submissions, setSubmissions] = useState<SubmissionHistoryItem[]>([]);
  const [submissionTotal, setSubmissionTotal] = useState(0);
  const [submissionStatus, setSubmissionStatus] = useState<string>('');
  const [submissionsLoading, setSubmissionsLoading] = useState(false);

  // Performance telemetry state
  const [analytics, setAnalytics] = useState<AnalyticsSummary | null>(null);
  const [analyticsLoading, setAnalyticsLoading] = useState(false);

  // Edit form state
  const [formData, setFormData] = useState({
    full_name: profile?.full_name || user?.fullName || 'Student Name',
    college: profile?.college || COLLEGES[0],
    course: profile?.course || 'B.Tech',
    department: profile?.department || 'Computer Science and Engineering',
    current_year: profile?.current_year || '3rd Year',
    preferred_language: (profile?.preferred_language || 'Python') as ProgrammingLanguage,
    learning_goal: profile?.learning_goal || 'Interview Preparation',
    dsa_level: (profile?.dsa_level || 'Intermediate') as DSALevel,
    bio: profile?.bio || 'Passionate software engineering student mastering algorithms.'
  });

  useEffect(() => {
    if (profile) {
      setFormData({
        full_name: profile.full_name || user?.fullName || 'Student Name',
        college: profile.college || COLLEGES[0],
        course: profile.course || 'B.Tech',
        department: profile.department || 'Computer Science and Engineering',
        current_year: profile.current_year || '3rd Year',
        preferred_language: profile.preferred_language || 'Python',
        learning_goal: profile.learning_goal || 'Interview Preparation',
        dsa_level: profile.dsa_level || 'Intermediate',
        bio: profile.bio || 'Passionate software engineering student mastering algorithms.'
      });
    }
  }, [profile, user]);

  useEffect(() => {
    if (activeTab === 'submissions') {
      loadSubmissions();
    } else if (activeTab === 'performance') {
      loadPerformance();
    }
  }, [activeTab, submissionStatus]);

  const loadSubmissions = async () => {
    setSubmissionsLoading(true);
    try {
      const res = await codingService.getStudentSubmissions({
        limit: 30,
        status: submissionStatus || undefined,
      });
      setSubmissions(res.submissions);
      setSubmissionTotal(res.total);
    } catch {
      toast.error('Failed to load submission history');
    } finally {
      setSubmissionsLoading(false);
    }
  };

  const loadPerformance = async () => {
    setAnalyticsLoading(true);
    try {
      const data = await analyticsService.getAnalytics();
      setAnalytics(data);
    } catch {
      toast.error('Failed to load performance metrics');
    } finally {
      setAnalyticsLoading(false);
    }
  };

  const handleSave = async () => {
    setIsSaving(true);
    try {
      await updateProfile(formData);
      updateLocalUser({ fullName: formData.full_name, full_name: formData.full_name });
      setIsEditing(false);
      toast.success('Profile updated successfully!');
    } catch (err: any) {
      toast.error(err?.message || 'Failed to update profile');
    } finally {
      setIsSaving(false);
    }
  };

  const getStatusBadge = (status: string) => {
    const s = (status || '').toLowerCase();
    if (s.includes('accepted')) {
      return <span className="px-2 py-0.5 rounded bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 text-[10px] font-mono">Accepted</span>;
    }
    if (s.includes('wrong')) {
      return <span className="px-2 py-0.5 rounded bg-rose-500/10 text-rose-400 border border-rose-500/20 text-[10px] font-mono">Wrong Answer</span>;
    }
    if (s.includes('time')) {
      return <span className="px-2 py-0.5 rounded bg-amber-500/10 text-amber-400 border border-amber-500/20 text-[10px] font-mono">Time Limit Exceeded</span>;
    }
    if (s.includes('compile') || s.includes('syntax')) {
      return <span className="px-2 py-0.5 rounded bg-orange-500/10 text-orange-400 border border-orange-500/20 text-[10px] font-mono">Compile Error</span>;
    }
    return <span className="px-2 py-0.5 rounded bg-rose-500/10 text-rose-400 border border-rose-500/20 text-[10px] font-mono">{status}</span>;
  };

  return (
    <div className="space-y-8 max-w-5xl mx-auto pb-16">
      {/* Header Banner */}
      <div className="p-6 sm:p-8 rounded-2xl bg-white border border-[#E2E8F0] shadow-sm flex flex-col sm:flex-row items-center justify-between gap-6">
        <div className="flex items-center gap-5">
          <div className="w-16 h-16 rounded-2xl bg-gradient-to-tr from-[#EF4444] to-[#B91C1C] flex items-center justify-center text-white font-bold text-2xl shadow-xl shadow-[#EF4444]/20">
            {formData.full_name?.charAt(0) || 'S'}
          </div>
          <div>
            <h1 className="text-2xl font-bold text-[#0F172A] font-mono">
              {formData.full_name}
            </h1>
            <p className="text-xs text-[#64748B] flex items-center gap-1.5 mt-0.5">
              <Mail className="w-3.5 h-3.5" />
              <span>{user?.email || 'student@gmail.com'}</span>
            </p>
            <div className="flex items-center gap-2 mt-2">
              <Badge variant="purple" size="sm">
                {formData.dsa_level.toUpperCase()} CODER
              </Badge>
              <Badge variant="easy" size="sm">
                VERIFIED STUDENT
              </Badge>
            </div>
          </div>
        </div>

        {activeTab === 'profile' && (
          <Button
            variant={isEditing ? 'secondary' : 'primary'}
            size="sm"
            onClick={() => {
              if (isEditing) handleSave();
              else setIsEditing(true);
            }}
            isLoading={isSaving}
            leftIcon={isEditing ? <Save className="w-4 h-4" /> : <Edit2 className="w-4 h-4" />}
          >
            {isEditing ? 'Save Profile' : 'Edit Profile'}
          </Button>
        )}
      </div>

      {/* Navigation Tabs */}
      <div className="flex items-center gap-2 border-b border-[#E2E8F0] pb-3 text-xs font-mono">
        <button
          onClick={() => setActiveTab('profile')}
          className={`px-4 py-2 rounded-xl transition-all flex items-center gap-2 ${
            activeTab === 'profile'
              ? 'bg-[#F1F5F9] text-[#0F172A] font-semibold border border-[#E2E8F0]'
              : 'text-[#64748B] hover:text-[#0F172A]'
          }`}
        >
          <User className="w-4 h-4" />
          <span>Profile Info</span>
        </button>

        <button
          onClick={() => setActiveTab('submissions')}
          className={`px-4 py-2 rounded-xl transition-all flex items-center gap-2 ${
            activeTab === 'submissions'
              ? 'bg-[#F1F5F9] text-[#0F172A] font-semibold border border-[#E2E8F0]'
              : 'text-[#64748B] hover:text-[#0F172A]'
          }`}
        >
          <Code2 className="w-4 h-4" />
          <span>Submission History</span>
        </button>

        <button
          onClick={() => setActiveTab('performance')}
          className={`px-4 py-2 rounded-xl transition-all flex items-center gap-2 ${
            activeTab === 'performance'
              ? 'bg-[#F1F5F9] text-[#0F172A] font-semibold border border-[#E2E8F0]'
              : 'text-[#64748B] hover:text-[#0F172A]'
          }`}
        >
          <Activity className="w-4 h-4" />
          <span>Performance Telemetry</span>
        </button>
      </div>

      {/* Tab 1: Profile Information */}
      {activeTab === 'profile' && (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <Card className="p-6 space-y-4">
            <div className="flex items-center gap-2 pb-2 border-b border-[#E2E8F0] text-[#DC2626] font-mono text-xs font-bold">
              <GraduationCap className="w-4 h-4 text-[#EF4444]" />
              <span>Academic Credentials</span>
            </div>

            <div className="space-y-4 text-xs">
              <div className="space-y-1">
                <label className="text-[#94A3B8] font-mono text-[10px] uppercase">College / University</label>
                {isEditing ? (
                  <select
                    value={formData.college}
                    onChange={(e) => setFormData({ ...formData, college: e.target.value })}
                    className="w-full bg-white text-[#0F172A] p-2.5 rounded-xl border border-[#E2E8F0] outline-none"
                  >
                    {COLLEGES.map((c) => (
                      <option key={c} value={c}>{c}</option>
                    ))}
                  </select>
                ) : (
                  <p className="text-[#0F172A] font-medium p-1">{formData.college}</p>
                )}
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div className="space-y-1">
                  <label className="text-[#94A3B8] font-mono text-[10px] uppercase">Course</label>
                  {isEditing ? (
                    <select
                      value={formData.course}
                      onChange={(e) => setFormData({ ...formData, course: e.target.value })}
                      className="w-full bg-white text-[#0F172A] p-2.5 rounded-xl border border-[#E2E8F0] outline-none"
                    >
                      {COURSES.map((c) => (
                        <option key={c} value={c}>{c}</option>
                      ))}
                    </select>
                  ) : (
                    <p className="text-[#0F172A] font-medium p-1">{formData.course}</p>
                  )}
                </div>

                <div className="space-y-1">
                  <label className="text-[#94A3B8] font-mono text-[10px] uppercase">Year</label>
                  {isEditing ? (
                    <select
                      value={formData.current_year}
                      onChange={(e) => setFormData({ ...formData, current_year: e.target.value })}
                      className="w-full bg-white text-[#0F172A] p-2.5 rounded-xl border border-[#E2E8F0] outline-none"
                    >
                      {ACADEMIC_YEARS.map((y) => (
                        <option key={y} value={y}>{y}</option>
                      ))}
                    </select>
                  ) : (
                    <p className="text-[#0F172A] font-medium p-1">{formData.current_year}</p>
                  )}
                </div>
              </div>

              <div className="space-y-1">
                <label className="text-[#94A3B8] font-mono text-[10px] uppercase">Department</label>
                {isEditing ? (
                  <select
                    value={formData.department}
                    onChange={(e) => setFormData({ ...formData, department: e.target.value })}
                    className="w-full bg-white text-[#0F172A] p-2.5 rounded-xl border border-[#E2E8F0] outline-none"
                  >
                    {DEPARTMENTS.map((d) => (
                      <option key={d} value={d}>{d}</option>
                    ))}
                  </select>
                ) : (
                  <p className="text-[#0F172A] font-medium p-1">{formData.department}</p>
                )}
              </div>
            </div>
          </Card>

          <Card className="p-6 space-y-4">
            <div className="flex items-center gap-2 pb-2 border-b border-[#E2E8F0] text-[#DC2626] font-mono text-xs font-bold">
              <Code2 className="w-4 h-4 text-[#EF4444]" />
              <span>DSA & Programming Focus</span>
            </div>

            <div className="space-y-4 text-xs">
              <div className="space-y-1">
                <label className="text-[#94A3B8] font-mono text-[10px] uppercase">Preferred Language</label>
                {isEditing ? (
                  <select
                    value={formData.preferred_language}
                    onChange={(e) => setFormData({ ...formData, preferred_language: e.target.value as ProgrammingLanguage })}
                    className="w-full bg-white text-[#0F172A] p-2.5 rounded-xl border border-[#E2E8F0] outline-none"
                  >
                    {PROGRAMMING_LANGUAGES.map((l: any) => (
                      <option key={l.id || l} value={l.name || l}>{l.name || l}</option>
                    ))}
                  </select>
                ) : (
                  <p className="text-[#0F172A] font-medium p-1">{formData.preferred_language}</p>
                )}
              </div>

              <div className="space-y-1">
                <label className="text-[#94A3B8] font-mono text-[10px] uppercase">Target Learning Goal</label>
                {isEditing ? (
                  <select
                    value={formData.learning_goal}
                    onChange={(e) => setFormData({ ...formData, learning_goal: e.target.value })}
                    className="w-full bg-white text-[#0F172A] p-2.5 rounded-xl border border-[#E2E8F0] outline-none"
                  >
                    {LEARNING_GOALS.map((g: any) => (
                      <option key={g.id || g} value={g.title || g}>{g.title || g}</option>
                    ))}
                  </select>
                ) : (
                  <p className="text-[#0F172A] font-medium p-1">{formData.learning_goal}</p>
                )}
              </div>

              <div className="space-y-1">
                <label className="text-[#94A3B8] font-mono text-[10px] uppercase">Bio / Purpose</label>
                {isEditing ? (
                  <textarea
                    rows={2}
                    value={formData.bio}
                    onChange={(e) => setFormData({ ...formData, bio: e.target.value })}
                    className="w-full bg-white text-[#0F172A] p-2.5 rounded-xl border border-[#E2E8F0] outline-none resize-none"
                  />
                ) : (
                  <p className="text-[#64748B] p-1 italic leading-relaxed">{formData.bio}</p>
                )}
              </div>
            </div>
          </Card>
        </div>
      )}

      {/* Tab 2: Global Submission History */}
      {activeTab === 'submissions' && (
        <div className="space-y-4">
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 p-4 rounded-xl bg-white border border-[#E2E8F0]">
            <div className="flex items-center gap-2">
              <Code2 className="w-4 h-4 text-[#EF4444]" />
              <span className="text-xs font-mono font-bold text-[#0F172A]">
                All Submissions ({submissionTotal})
              </span>
            </div>

            <div className="flex items-center gap-2">
              <Filter className="w-3.5 h-3.5 text-[#94A3B8]" />
              <select
                value={submissionStatus}
                onChange={(e) => setSubmissionStatus(e.target.value)}
                className="bg-[#F1F5F9] text-xs text-[#64748B] px-3 py-1.5 rounded-lg border border-[#E2E8F0] outline-none font-mono"
              >
                <option value="">All Statuses</option>
                <option value="Accepted">Accepted</option>
                <option value="Wrong Answer">Wrong Answer</option>
                <option value="Runtime Error">Runtime Error</option>
                <option value="Time Limit Exceeded">Time Limit Exceeded</option>
                <option value="Compilation Error">Compilation Error</option>
              </select>

              <button
                onClick={loadSubmissions}
                disabled={submissionsLoading}
                className="p-1.5 rounded-lg bg-[#F1F5F9] border border-[#E2E8F0] text-[#64748B] hover:text-[#0F172A]"
              >
                <RefreshCw className={`w-3.5 h-3.5 ${submissionsLoading ? 'animate-spin' : ''}`} />
              </button>
            </div>
          </div>

          {submissionsLoading ? (
            <div className="space-y-3">
              {[1, 2, 3, 4, 5].map((i) => (
                <div key={i} className="h-16 rounded-xl bg-[#F1F5F9] border border-[#E2E8F0] animate-pulse" />
              ))}
            </div>
          ) : submissions.length > 0 ? (
            <div className="rounded-xl border border-[#E2E8F0] overflow-hidden bg-white">
              <div className="overflow-x-auto">
                <table className="w-full text-left text-xs font-mono">
                  <thead className="bg-[#F1F5F9] border-b border-[#E2E8F0] text-[#94A3B8] uppercase text-[10px]">
                    <tr>
                      <th className="px-4 py-3">Problem</th>
                      <th className="px-4 py-3">Status</th>
                      <th className="px-4 py-3">Language</th>
                      <th className="px-4 py-3">Runtime</th>
                      <th className="px-4 py-3">Memory</th>
                      <th className="px-4 py-3">Submitted At</th>
                      <th className="px-4 py-3 text-right">Action</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-[#E2E8F0]">
                    {submissions.map((sub) => (
                      <tr key={sub.id || sub.submission_id} className="hover:bg-[#F1F5F9]/60 transition-colors">
                        <td className="px-4 py-3 text-[#0F172A] font-semibold">
                          {sub.problem_title || sub.problem_id || 'Problem'}
                        </td>
                        <td className="px-4 py-3">
                          {getStatusBadge(sub.status)}
                        </td>
                        <td className="px-4 py-3 text-[#64748B] capitalize">
                          {sub.language || 'python'}
                        </td>
                        <td className="px-4 py-3 text-[#64748B]">
                          {sub.execution_time_ms ? `${sub.execution_time_ms}ms` : '—'}
                        </td>
                        <td className="px-4 py-3 text-[#64748B]">
                          {sub.memory_used_kb ? `${Math.round(sub.memory_used_kb / 1024 * 10) / 10}MB` : '—'}
                        </td>
                        <td className="px-4 py-3 text-[#94A3B8] text-[11px]">
                          {sub.submitted_at ? new Date(sub.submitted_at).toLocaleString() : '—'}
                        </td>
                        <td className="px-4 py-3 text-right">
                          <button
                            onClick={() => navigate(`/practice/${sub.problem_id}`)}
                            className="inline-flex items-center gap-1 text-[11px] text-[#EF4444] hover:text-[#DC2626]"
                          >
                            <span>Open</span>
                            <ExternalLink className="w-3 h-3" />
                          </button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          ) : (
            <div className="p-12 text-center rounded-xl bg-white border border-dashed border-[#E2E8F0] text-xs text-[#64748B] space-y-2">
              <Code2 className="w-8 h-8 text-[#94A3B8] mx-auto mb-2" />
              <p>No submissions recorded for this filter.</p>
              <Button size="sm" variant="outline" onClick={() => navigate('/practice')}>
                Go to Practice Arena
              </Button>
            </div>
          )}
        </div>
      )}

      {/* Tab 3: Performance Telemetry */}
      {activeTab === 'performance' && (
        <div className="space-y-6">
          {analyticsLoading || !analytics ? (
            <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
              {[1, 2, 3, 4].map((i) => (
                <div key={i} className="h-28 rounded-2xl bg-[#F1F5F9] border border-[#E2E8F0] animate-pulse" />
              ))}
            </div>
          ) : (
            <>
              <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 font-mono">
                <div className="p-5 rounded-2xl bg-white border border-[#E2E8F0] space-y-1">
                  <div className="flex items-center justify-between text-xs text-[#94A3B8]">
                    <span>TOTAL_SOLVED</span>
                    <CheckCircle2 className="w-4 h-4 text-emerald-400" />
                  </div>
                  <div className="text-3xl font-bold text-[#0F172A]">{analytics.totalSolved}</div>
                  <span className="text-[11px] text-emerald-400">
                    Accuracy: {analytics.accuracyRate}%
                  </span>
                </div>

                <div className="p-5 rounded-2xl bg-white border border-[#E2E8F0] space-y-1">
                  <div className="flex items-center justify-between text-xs text-[#94A3B8]">
                    <span>ACTIVE_STREAK</span>
                    <Flame className="w-4 h-4 text-amber-400" />
                  </div>
                  <div className="text-3xl font-bold text-amber-400">{analytics.streakDays}d</div>
                  <span className="text-[11px] text-[#64748B]">Consistent practice</span>
                </div>

                <div className="p-5 rounded-2xl bg-white border border-[#E2E8F0] space-y-1">
                  <div className="flex items-center justify-between text-xs text-[#94A3B8]">
                    <span>HOURS_INVESTED</span>
                    <Clock className="w-4 h-4 text-[#DC2626]" />
                  </div>
                  <div className="text-3xl font-bold text-[#DC2626]">{analytics.hoursSpent}h</div>
                  <span className="text-[11px] text-[#64748B]">{analytics.attemptedProblems} total attempts</span>
                </div>

                <div className="p-5 rounded-2xl bg-white border border-[#E2E8F0] space-y-1">
                  <div className="flex items-center justify-between text-xs text-[#94A3B8]">
                    <span>REVISION_QUEUE</span>
                    <Activity className="w-4 h-4 text-[#DC2626]" />
                  </div>
                  <div className="text-3xl font-bold text-[#DC2626]">{analytics.needsRevision}</div>
                  <span className="text-[11px] text-[#64748B]">Marked for review</span>
                </div>
              </div>

              {/* Difficulty Breakdown */}
              <div className="p-6 rounded-2xl bg-white border border-[#E2E8F0] space-y-4 font-mono">
                <h3 className="text-sm font-bold text-[#0F172A]">Difficulty Distribution</h3>
                <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                  {analytics.difficultyDistribution.map((d) => (
                    <div
                      key={d.name}
                      className="p-4 rounded-xl bg-[#F1F5F9] border border-[#E2E8F0] space-y-2"
                    >
                      <div className="flex items-center justify-between">
                        <span className="text-xs font-bold text-[#0F172A]">{d.name}</span>
                        <span className="text-xs text-[#64748B]">{d.count} solved</span>
                      </div>
                      <div className="w-full bg-[#F1F5F9] h-1.5 rounded-full overflow-hidden">
                        <div
                          className={`h-full ${
                            d.name === 'Easy'
                              ? 'bg-emerald-400'
                              : d.name === 'Medium'
                              ? 'bg-amber-400'
                              : 'bg-rose-400'
                          }`}
                          style={{
                            width: `${Math.min(100, Math.round((d.count / Math.max(1, analytics.totalSolved)) * 100))}%`,
                          }}
                        />
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </>
          )}
        </div>
      )}
    </div>
  );
};

export default ProfilePage;
