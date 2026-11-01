import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { motion, AnimatePresence } from 'framer-motion';
import {
  GraduationCap,
  Sparkles,
  ArrowRight,
  ArrowLeft,
  CheckCircle2,
  Target,
  Edit2
} from 'lucide-react';
import {
  COLLEGES,
  ACADEMIC_YEARS,
  PROGRAMMING_LANGUAGES,
  LEARNING_GOALS,
  DSA_LEVELS
} from '../../utils/constants';
import { OnboardingData, ProgrammingLanguage, DSALevel } from '../../types/profile.types';
import profileService from '../../services/profile.service';
import { useAuth } from '../../hooks/useAuth';
import Button from '../../components/ui/Button';
import Card from '../../components/ui/Card';
import Badge from '../../components/ui/Badge';
import { toast } from 'sonner';

export const OnboardingPage: React.FC = () => {
  const navigate = useNavigate();
  const { user, updateLocalUser } = useAuth();
  const [step, setStep] = useState<number>(1);
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);

  // If user is already onboarded, avoid stuck state and go straight to dashboard
  useEffect(() => {
    if (user?.onboarding_completed || user?.isOnboarded) {
      navigate('/dashboard', { replace: true });
    }
  }, [user, navigate]);

  // Form state aligned with Go backend schema
  const [formData, setFormData] = useState<OnboardingData>({
    college: COLLEGES[0],
    course: 'B.Tech',
    department: 'Computer Science and Engineering',
    current_year: '3rd Year',
    preferred_language: 'Python',
    learning_goal: 'Interview Preparation',
    dsa_level: 'Intermediate'
  });

  const handleNext = () => {
    if (step < 3) {
      setStep(step + 1);
    }
  };

  const handleBack = () => {
    if (step > 1) {
      setStep(step - 1);
    }
  };

  const handleComplete = async () => {
    setIsSubmitting(true);
    try {
      const profile = await profileService.submitOnboarding(formData);
      updateLocalUser({
        onboarding_completed: true,
        isOnboarded: true,
        full_name: profile.full_name || user?.fullName || 'Student Name'
      });
      toast.success('Onboarding complete! Welcome to your cockpit.');
      setTimeout(() => {
        navigate('/dashboard', { replace: true });
      }, 100);
    } catch (err: any) {
      toast.error(err?.message || 'Failed to complete profile');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="min-h-screen bg-[#07070B] flex flex-col items-center justify-center p-4 sm:p-6 text-[#F4F4F6] relative overflow-hidden font-mono">
      {/* Background glow */}
      <div className="absolute top-1/4 left-1/2 -translate-x-1/2 w-[600px] h-[400px] bg-[#8B5CF6]/10 blur-[130px] pointer-events-none" />

      {/* Brand & Progress Header */}
      <div className="w-full max-w-2xl mb-8 space-y-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="w-7 h-7 rounded-xl bg-gradient-to-tr from-[#8B5CF6] to-[#6366F1] flex items-center justify-center shadow-md shadow-[#8B5CF6]/30">
              <Sparkles className="w-3.5 h-3.5 text-white" />
            </div>
            <span className="font-bold text-lg text-[#F4F4F6] tracking-tight font-mono">
              Learn<span className="text-[#8B5CF6]">X</span>
            </span>
          </div>
          <span className="text-xs font-mono text-[#8E8EA0]">Stage {step} of 3</span>
        </div>

        {/* Progress Bar */}
        <div className="grid grid-cols-3 gap-2">
          {[1, 2, 3].map((s) => (
            <div
              key={s}
              className={`h-1.5 rounded-full transition-all duration-300 ${
                s <= step ? 'bg-[#8B5CF6]' : 'bg-[#1C1C28]'
              }`}
            />
          ))}
        </div>
      </div>

      {/* Form Card with Animated Transitions */}
      <Card className="w-full max-w-2xl bg-[#0D0D14] border border-[#1C1C28] p-6 sm:p-8 rounded-2xl shadow-2xl relative">
        <AnimatePresence mode="wait">
          {/* STAGE 1: ACADEMIC INFORMATION */}
          {step === 1 && (
            <motion.div
              key="step-1"
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: -20 }}
              transition={{ duration: 0.2 }}
              className="space-y-6"
            >
              <div className="space-y-1">
                <div className="inline-flex items-center gap-1.5 text-xs text-[#A78BFA] font-mono">
                  <GraduationCap className="w-4 h-4" />
                  <span>Academic Calibration</span>
                </div>
                <h2 className="text-xl sm:text-2xl font-bold font-mono text-[#F4F4F6]">
                  Tell us about your university
                </h2>
                <p className="text-xs text-[#8E8EA0] font-sans">
                  We use your academic year and department to tailor the recommended problem pacing.
                </p>
              </div>

              <div className="space-y-4 text-xs font-sans">
                <div className="space-y-1.5">
                  <label className="font-medium text-[#8E8EA0] font-mono">College / University Name</label>
                  <input
                    type="text"
                    value={formData.college}
                    onChange={(e) => setFormData({ ...formData, college: e.target.value })}
                    placeholder="e.g. Stanford University or IIT Bombay"
                    className="w-full bg-[#12121B] text-[#F4F4F6] text-xs font-mono rounded-xl px-4 py-2.5 border border-[#1C1C28] focus:border-[#8B5CF6] focus:outline-none"
                  />
                </div>

                <div className="space-y-1.5">
                  <label className="font-medium text-[#8E8EA0] font-mono">Course (e.g. B.Tech, B.E, B.Sc, BCA)</label>
                  <input
                    type="text"
                    value={formData.course}
                    onChange={(e) => setFormData({ ...formData, course: e.target.value })}
                    placeholder="e.g. B.Tech"
                    className="w-full bg-[#12121B] text-[#F4F4F6] text-xs font-mono rounded-xl px-4 py-2.5 border border-[#1C1C28] focus:border-[#8B5CF6] focus:outline-none"
                  />
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div className="space-y-1.5">
                    <label className="font-medium text-[#8E8EA0] font-mono">Department</label>
                    <input
                      type="text"
                      value={formData.department}
                      onChange={(e) => setFormData({ ...formData, department: e.target.value })}
                      placeholder="e.g. Computer Science and Engineering"
                      className="w-full bg-[#12121B] text-[#F4F4F6] text-xs font-mono rounded-xl px-4 py-2.5 border border-[#1C1C28] focus:border-[#8B5CF6] focus:outline-none"
                    />
                  </div>

                  <div className="space-y-1.5">
                    <label className="font-medium text-[#8E8EA0] font-mono">Current Academic Year</label>
                    <select
                      value={formData.current_year}
                      onChange={(e) => setFormData({ ...formData, current_year: e.target.value })}
                      className="w-full bg-[#12121B] text-[#F4F4F6] text-xs font-mono rounded-xl px-4 py-2.5 border border-[#1C1C28] focus:border-[#8B5CF6] focus:outline-none cursor-pointer"
                    >
                      {ACADEMIC_YEARS.map((y) => (
                        <option key={y} value={y} className="bg-[#12121B] text-[#F4F4F6]">
                          {y}
                        </option>
                      ))}
                    </select>
                  </div>
                </div>
              </div>
            </motion.div>
          )}

          {/* STAGE 2: LEARNING PREFERENCES */}
          {step === 2 && (
            <motion.div
              key="step-2"
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: -20 }}
              transition={{ duration: 0.2 }}
              className="space-y-6"
            >
              <div className="space-y-1">
                <div className="inline-flex items-center gap-1.5 text-xs text-[#A78BFA] font-mono">
                  <Target className="w-4 h-4" />
                  <span>Target Calibration</span>
                </div>
                <h2 className="text-xl sm:text-2xl font-bold font-mono text-[#F4F4F6]">
                  Your coding weapon & career goals
                </h2>
                <p className="text-xs text-[#8E8EA0] font-sans">
                  Select your primary language and technical benchmark.
                </p>
              </div>

              {/* Language selection */}
              <div className="space-y-2">
                <label className="text-xs font-mono text-[#8E8EA0]">
                  Preferred Language (Allowed: C, C++, Java, Python, Go)
                </label>
                <div className="grid grid-cols-2 sm:grid-cols-5 gap-2">
                  {PROGRAMMING_LANGUAGES.map((lang) => {
                    const isSelected = formData.preferred_language === lang.id;
                    return (
                      <button
                        type="button"
                        key={lang.id}
                        onClick={() =>
                          setFormData({ ...formData, preferred_language: lang.id as ProgrammingLanguage })
                        }
                        className={`p-3 rounded-xl border text-center transition-all cursor-pointer ${
                          isSelected
                            ? 'bg-[#8B5CF6]/20 border-[#8B5CF6] text-white shadow-md shadow-[#8B5CF6]/20'
                            : 'bg-[#12121B] border-[#1C1C28] text-[#8E8EA0] hover:border-[#2D2D3E]'
                        }`}
                      >
                        <div className="font-mono font-bold text-sm">{lang.name}</div>
                        <div className="text-[10px] text-[#5A5A6E] mt-0.5">{lang.ext}</div>
                      </button>
                    );
                  })}
                </div>
              </div>

              {/* Learning goal */}
              <div className="space-y-2">
                <label className="text-xs font-mono text-[#8E8EA0]">Learning Goal</label>
                <div className="grid grid-cols-1 sm:grid-cols-3 gap-2.5">
                  {LEARNING_GOALS.map((goal) => {
                    const isSelected = formData.learning_goal === goal.id;
                    return (
                      <button
                        type="button"
                        key={goal.id}
                        onClick={() => setFormData({ ...formData, learning_goal: goal.id })}
                        className={`p-3 rounded-xl border text-left transition-all cursor-pointer ${
                          isSelected
                            ? 'bg-[#8B5CF6]/20 border-[#8B5CF6] text-white shadow-sm'
                            : 'bg-[#12121B] border-[#1C1C28] text-[#8E8EA0] hover:border-[#2D2D3E]'
                        }`}
                      >
                        <div className="text-xs font-bold text-[#F4F4F6] font-mono">{goal.title}</div>
                        <div className="text-[11px] text-[#5A5A6E] mt-1 font-sans">{goal.desc}</div>
                      </button>
                    );
                  })}
                </div>
              </div>

              {/* DSA Level */}
              <div className="space-y-2">
                <label className="text-xs font-mono text-[#8E8EA0]">Current DSA Level</label>
                <div className="grid grid-cols-1 sm:grid-cols-3 gap-2.5">
                  {DSA_LEVELS.map((lvl) => {
                    const isSelected = formData.dsa_level === lvl.id;
                    return (
                      <button
                        type="button"
                        key={lvl.id}
                        onClick={() => setFormData({ ...formData, dsa_level: lvl.id as DSALevel })}
                        className={`p-3 rounded-xl border text-left transition-all cursor-pointer ${
                          isSelected
                            ? 'bg-[#8B5CF6]/20 border-[#8B5CF6] text-white'
                            : 'bg-[#12121B] border-[#1C1C28] text-[#8E8EA0] hover:border-[#2D2D3E]'
                        }`}
                      >
                        <div className="text-xs font-bold text-[#F4F4F6] font-mono">{lvl.title}</div>
                        <div className="text-[10px] text-[#5A5A6E] mt-1 font-sans">{lvl.desc}</div>
                      </button>
                    );
                  })}
                </div>
              </div>
            </motion.div>
          )}

          {/* STAGE 3: PROFILE SUMMARY */}
          {step === 3 && (
            <motion.div
              key="step-3"
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: -20 }}
              transition={{ duration: 0.2 }}
              className="space-y-6"
            >
              <div className="space-y-1">
                <div className="inline-flex items-center gap-1.5 text-xs text-emerald-400 font-mono">
                  <CheckCircle2 className="w-4 h-4" />
                  <span>Review & Confirmation</span>
                </div>
                <h2 className="text-xl sm:text-2xl font-bold font-mono text-[#F4F4F6]">
                  Verify Your Curriculum Configuration
                </h2>
                <p className="text-xs text-[#8E8EA0] font-sans">
                  Your data will be securely synced with the LearnX backend and Supabase PostgreSQL.
                </p>
              </div>

              <div className="space-y-4">
                {/* Academic summary */}
                <div className="p-4 bg-[#12121B] rounded-xl border border-[#1C1C28] space-y-2">
                  <div className="flex items-center justify-between pb-2 border-b border-[#1C1C28]">
                    <span className="text-xs font-bold text-[#F4F4F6] flex items-center gap-1.5 font-mono">
                      <GraduationCap className="w-3.5 h-3.5 text-[#8B5CF6]" /> Academic Background
                    </span>
                    <button
                      onClick={() => setStep(1)}
                      className="text-xs text-[#8B5CF6] hover:text-[#A78BFA] flex items-center gap-1 font-mono cursor-pointer"
                    >
                      <Edit2 className="w-3 h-3" /> Edit
                    </button>
                  </div>
                  <div className="grid grid-cols-2 gap-2 text-xs font-mono">
                    <div>
                      <span className="text-[#5A5A6E] block text-[11px]">College</span>
                      <span className="font-medium text-[#F4F4F6] truncate block">{formData.college}</span>
                    </div>
                    <div>
                      <span className="text-[#5A5A6E] block text-[11px]">Degree</span>
                      <span className="font-medium text-[#F4F4F6]">{formData.course}</span>
                    </div>
                    <div>
                      <span className="text-[#5A5A6E] block text-[11px]">Department</span>
                      <span className="font-medium text-[#F4F4F6] truncate block">{formData.department}</span>
                    </div>
                    <div>
                      <span className="text-[#5A5A6E] block text-[11px]">Current Year</span>
                      <span className="font-medium text-[#F4F4F6]">{formData.current_year}</span>
                    </div>
                  </div>
                </div>

                {/* Preferences summary */}
                <div className="p-4 bg-[#12121B] rounded-xl border border-[#1C1C28] space-y-2">
                  <div className="flex items-center justify-between pb-2 border-b border-[#1C1C28]">
                    <span className="text-xs font-bold text-[#F4F4F6] flex items-center gap-1.5 font-mono">
                      <Target className="w-3.5 h-3.5 text-[#8B5CF6]" /> Target Calibration
                    </span>
                    <button
                      onClick={() => setStep(2)}
                      className="text-xs text-[#8B5CF6] hover:text-[#A78BFA] flex items-center gap-1 font-mono cursor-pointer"
                    >
                      <Edit2 className="w-3 h-3" /> Edit
                    </button>
                  </div>
                  <div className="grid grid-cols-3 gap-2 text-xs font-mono">
                    <div>
                      <span className="text-[#5A5A6E] block text-[11px]">Preferred Language</span>
                      <span className="text-[#A78BFA] font-bold block mt-1">
                        {formData.preferred_language}
                      </span>
                    </div>
                    <div>
                      <span className="text-[#5A5A6E] block text-[11px]">Goal</span>
                      <span className="font-medium text-[#F4F4F6] block mt-1">
                        {formData.learning_goal}
                      </span>
                    </div>
                    <div>
                      <span className="text-[#5A5A6E] block text-[11px]">DSA Level</span>
                      <span className="font-medium text-[#F4F4F6] block mt-1">
                        {formData.dsa_level}
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            </motion.div>
          )}
        </AnimatePresence>

        {/* Footer controls: Back / Continue / Complete */}
        <div className="mt-8 pt-5 border-t border-[#1C1C28] flex items-center justify-between">
          {step > 1 ? (
            <Button
              variant="ghost"
              size="sm"
              onClick={handleBack}
              leftIcon={<ArrowLeft className="w-4 h-4" />}
            >
              Back
            </Button>
          ) : (
            <div />
          )}

          {step < 3 ? (
            <Button
              variant="primary"
              size="md"
              onClick={handleNext}
              rightIcon={<ArrowRight className="w-4 h-4" />}
            >
              Continue
            </Button>
          ) : (
            <Button
              variant="primary"
              size="md"
              onClick={handleComplete}
              isLoading={isSubmitting}
              rightIcon={<CheckCircle2 className="w-4 h-4" />}
            >
              Launch LearnX Cockpit
            </Button>
          )}
        </div>
      </Card>
    </div>
  );
};

export default OnboardingPage;
