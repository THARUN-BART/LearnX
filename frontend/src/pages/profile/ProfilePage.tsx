import React, { useState, useEffect } from 'react';
import {
  User,
  GraduationCap,
  Target,
  Code2,
  Mail,
  Edit2,
  CheckCircle2,
  Save
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
import Card from '../../components/ui/Card';
import Badge from '../../components/ui/Badge';
import Button from '../../components/ui/Button';
import { toast } from 'sonner';

export const ProfilePage: React.FC = () => {
  const { user, updateLocalUser } = useAuth();
  const { profile, updateProfile } = useProfile();

  const [isEditing, setIsEditing] = useState<boolean>(false);
  const [isSaving, setIsSaving] = useState<boolean>(false);

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

  return (
    <div className="space-y-8 max-w-4xl mx-auto">
      {/* Header Banner */}
      <div className="p-6 sm:p-8 rounded-2xl bg-gradient-to-r from-[#171722] via-[#1A1A28] to-[#101019] border border-[#232336] flex flex-col sm:flex-row items-center justify-between gap-6">
        <div className="flex items-center gap-5">
          <div className="w-16 h-16 rounded-2xl bg-gradient-to-tr from-[#6366F1] to-[#8B5CF6] flex items-center justify-center text-white font-bold text-2xl shadow-xl shadow-[#8B5CF6]/20">
            {formData.full_name?.charAt(0) || 'S'}
          </div>
          <div>
            <h1 className="text-2xl font-bold text-[#F5F5F7] font-mono">
              {formData.full_name}
            </h1>
            <p className="text-xs text-[#9292A3] flex items-center gap-1.5 mt-0.5">
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
      </div>

      {/* Main Profile Info Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Academic Card */}
        <Card className="p-6 space-y-4">
          <div className="flex items-center gap-2 pb-2 border-b border-[#232336] text-[#A78BFA] font-mono text-xs font-bold">
            <GraduationCap className="w-4 h-4 text-[#8B5CF6]" />
            <span>Academic Credentials (PostgreSQL)</span>
          </div>

          <div className="space-y-4 text-xs">
            <div className="space-y-1">
              <label className="text-[#65657B]">Full Name</label>
              {isEditing ? (
                <input
                  type="text"
                  value={formData.full_name}
                  onChange={(e) => setFormData({ ...formData, full_name: e.target.value })}
                  className="w-full bg-[#101019] text-[#F5F5F7] rounded-xl px-3 py-2 border border-[#232336]"
                />
              ) : (
                <p className="font-semibold text-[#F5F5F7]">{formData.full_name}</p>
              )}
            </div>

            <div className="space-y-1">
              <label className="text-[#65657B]">University / Institution</label>
              {isEditing ? (
                <input
                  type="text"
                  value={formData.college}
                  onChange={(e) => setFormData({ ...formData, college: e.target.value })}
                  className="w-full bg-[#101019] text-[#F5F5F7] rounded-xl px-3 py-2 border border-[#232336]"
                />
              ) : (
                <p className="font-semibold text-[#F5F5F7]">{formData.college}</p>
              )}
            </div>

            <div className="space-y-1">
              <label className="text-[#65657B]">Degree / Course</label>
              {isEditing ? (
                <input
                  type="text"
                  value={formData.course}
                  onChange={(e) => setFormData({ ...formData, course: e.target.value })}
                  className="w-full bg-[#101019] text-[#F5F5F7] rounded-xl px-3 py-2 border border-[#232336]"
                />
              ) : (
                <p className="font-semibold text-[#F5F5F7]">{formData.course}</p>
              )}
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-1">
                <label className="text-[#65657B]">Department</label>
                {isEditing ? (
                  <input
                    type="text"
                    value={formData.department}
                    onChange={(e) => setFormData({ ...formData, department: e.target.value })}
                    className="w-full bg-[#101019] text-[#F5F5F7] rounded-xl px-3 py-2 border border-[#232336]"
                  />
                ) : (
                  <p className="font-semibold text-[#F5F5F7]">{formData.department}</p>
                )}
              </div>

              <div className="space-y-1">
                <label className="text-[#65657B]">Current Year</label>
                {isEditing ? (
                  <select
                    value={formData.current_year}
                    onChange={(e) => setFormData({ ...formData, current_year: e.target.value })}
                    className="w-full bg-[#101019] text-[#F5F5F7] rounded-xl px-3 py-2 border border-[#232336]"
                  >
                    {ACADEMIC_YEARS.map((y) => (
                      <option key={y} value={y}>{y}</option>
                    ))}
                  </select>
                ) : (
                  <p className="font-semibold text-[#F5F5F7]">{formData.current_year}</p>
                )}
              </div>
            </div>
          </div>
        </Card>

        {/* Learning Preferences Card */}
        <Card className="p-6 space-y-4">
          <div className="flex items-center gap-2 pb-2 border-b border-[#232336] text-[#A78BFA] font-mono text-xs font-bold">
            <Target className="w-4 h-4 text-[#8B5CF6]" />
            <span>Target Calibration</span>
          </div>

          <div className="space-y-4 text-xs">
            <div className="space-y-1">
              <label className="text-[#65657B]">Preferred Language (C, C++, Java, Python, Go)</label>
              {isEditing ? (
                <select
                  value={formData.preferred_language}
                  onChange={(e) =>
                    setFormData({ ...formData, preferred_language: e.target.value as ProgrammingLanguage })
                  }
                  className="w-full bg-[#101019] text-[#F5F5F7] rounded-xl px-3 py-2 border border-[#232336]"
                >
                  {PROGRAMMING_LANGUAGES.map((l) => (
                    <option key={l.id} value={l.id}>{l.name}</option>
                  ))}
                </select>
              ) : (
                <p className="font-semibold text-[#F5F5F7] font-mono">{formData.preferred_language}</p>
              )}
            </div>

            <div className="space-y-1">
              <label className="text-[#65657B]">Learning Goal</label>
              {isEditing ? (
                <select
                  value={formData.learning_goal}
                  onChange={(e) => setFormData({ ...formData, learning_goal: e.target.value })}
                  className="w-full bg-[#101019] text-[#F5F5F7] rounded-xl px-3 py-2 border border-[#232336]"
                >
                  {LEARNING_GOALS.map((g) => (
                    <option key={g.id} value={g.id}>{g.title}</option>
                  ))}
                </select>
              ) : (
                <p className="font-semibold text-[#F5F5F7]">{formData.learning_goal}</p>
              )}
            </div>

            <div className="space-y-1">
              <label className="text-[#65657B]">DSA Proficiency Tier (Beginner, Intermediate, Advanced)</label>
              {isEditing ? (
                <select
                  value={formData.dsa_level}
                  onChange={(e) => setFormData({ ...formData, dsa_level: e.target.value as DSALevel })}
                  className="w-full bg-[#101019] text-[#F5F5F7] rounded-xl px-3 py-2 border border-[#232336]"
                >
                  {DSA_LEVELS.map((lvl) => (
                    <option key={lvl.id} value={lvl.id}>{lvl.title}</option>
                  ))}
                </select>
              ) : (
                <p className="font-semibold text-[#F5F5F7]">{formData.dsa_level}</p>
              )}
            </div>
          </div>
        </Card>
      </div>

      {/* Bio / Summary Statement */}
      <Card className="p-6 space-y-3">
        <h4 className="text-xs font-bold text-[#F5F5F7] font-mono">Bio & Focus Statement</h4>
        {isEditing ? (
          <textarea
            rows={3}
            value={formData.bio}
            onChange={(e) => setFormData({ ...formData, bio: e.target.value })}
            className="w-full bg-[#101019] text-[#F5F5F7] rounded-xl p-3 border border-[#232336] text-xs outline-none"
          />
        ) : (
          <p className="text-xs text-[#9292A3] leading-relaxed">
            {formData.bio}
          </p>
        )}
      </Card>
    </div>
  );
};

export default ProfilePage;
