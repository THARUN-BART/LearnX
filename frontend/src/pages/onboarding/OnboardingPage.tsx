import React, { useState } from 'react';
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
  COURSES,
  DEPARTMENTS,
  ACADEMIC_YEARS,
  PROGRAMMING_LANGUAGES,
  LEARNING_GOALS,
  DSA_LEVELS
} from '../../utils/constants';
import { OnboardingData, ProgrammingLanguage, DSALevel } from '../../types/profile.types';
import profileService from '../../services/profile.service';
import Button from '../../components/ui/Button';
import Card from '../../components/ui/Card';
import Badge from '../../components/ui/Badge';
import { toast } from 'sonner';

export const OnboardingPage: React.FC = () => {
  const navigate = useNavigate();
  const [step, setStep] = useState<number>(1);
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);

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
      await profileService.submitOnboarding(formData);
      toast.success('Onboarding complete! Welcome to your cockpit.');
      navigate('/dashboard', { replace: true });
    } catch (err: any) {
      toast.error(err?.message || 'Failed to complete profile');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="min-h-screen bg-[#08080D] flex flex-col items-center justify-center p-4 sm:p-6 text-[#F5F5F7] relative overflow-hidden">
      {/* Background glow */}
      <div className="absolute top-1/4 left-1/2 -translate-x-1/2 w-[600px] h-[400px] bg-[#8B5CF6]/10 blur-[130px] pointer-events-none" />

      {/* Brand & Progress Header */}
      <div className="w-full max-w-2xl mb-8 space-y-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="w-7 h-7 rounded-xl bg-gradient-to-tr from-[#8B5CF6] to-[#6366F1] flex items-center justify-center shadow-md shadow-[#8B5CF6]/30">
              <Sparkles className="w-3.5 h-3.5 text-white" />
            </div>
            <span className="font-bold text-lg text-[#F5F5F7] tracking-tight font-mono">
              Learn<span className="text-[#8B5CF6]">X</span>
            </span>
          </div>
          <span className="text-xs font-mono text-[#9292A3]">Stage {step} of 3</span>
        </div>

        {/* Progress Bar */}
        <div className="grid grid-cols-3 gap-2">
          {[1, 2, 3].map((s) => (
            <div
              key={s}
              className={`h-1.5 rounded-full transition-all duration-300 ${
                s <= step ? 'bg-[#8B5CF6]' : 'bg-[#232336]'
              }`}
            />
          ))}
        </div>
      </div>

      {/* Form Card with Animated Transitions */}
      <Card className="w-full max-w-2xl bg-[#171722]/90 backdrop-blur-xl border border-[#232336] p-6 sm:p-8 rounded-2xl shadow-2xl relative">
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
                  <span>Academic Information</span>
                </div>
                <h2 className="text-xl sm:text-2xl font-bold font-mono text-[#F5F5F7]">
                  Tell us about your university
                </h2>
                <p className="text-xs text-[#9292A3]">
                  We use your academic year and department to tailor the recommended problem pacing.
                </p>
              </div>

              <div className="space-y-4">
                <div className="space-y-1.5">
                  <label className="text-xs font-medium text-[#9292A3]">College / University Name</label>
                  <input
                    type="text"
                    value={formData.college}
                    onChange={(e) => setFormData({ ...formData, college: e.target.value })}
                    placeholder="e.g. Stanford University or IIT Bombay"
                    className="w-full bg-[#101019] text-[#F5F5F7] text-sm rounded-xl px-4 py-2.5 border border-[#232336] focus:border-[#8B5CF6] focus:outline-none"
                  />
                </div>

                <div className="space-y-1.5">
                  <label className="text-xs font-medium text-[#9292A3]">Course (e.g. B.Tech, B.E, B.Sc, BCA)</label>
                  <input
                    type="text"
                    value={formData.course}
                    onChange={(e) => setFormData({ ...formData, course: e.target.value })}
                    placeholder="e.g. B.Tech"
                    className="w-full bg-[#101019] text-[#F5F5F7] text-sm rounded-xl px-4 py-2.5 border border-[#232336] focus:border-[#8B5CF6] focus:outline-none"
                  />
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div className="space-y-1.5">
                    <label className="text-xs font-medium text-[#9292A3]">Department</label>
                    <input
                      type="text"
                      value={formData.department}
                      onChange={(e) => setFormData({ ...formData, department: e.target.value })}
                      placeholder="e.g. Computer Science and Engineering"
                      className="w-full bg-[#101019] text-[#F5F5F7] text-sm rounded-xl px-4 py-2.5 border border-[#232336] focus:border-[#8B5CF6] focus:outline-none"
                    />
                  </div>

                  <div className="space-y-1.5">
                    <label className="text-xs font-medium text-[#9292A3]">Current Academic Year</label>
                    <select
                      value={formData.current_year}
                      onChange={(e) => setFormData({ ...formData, current_year: e.target.value })}
                      className="w-full bg-[#101019] text-[#F5F5F7] text-sm rounded-xl px-4 py-2.5 border border-[#232336] focus:border-[#8B5CF6] focus:outline-none"
                    >
                      {ACADEMIC_YEARS.map((y) => (
                        <option key={y} value={y} className="bg-[#101019] text-[#F5F5F7]">
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
                  <span>Learning Preferences</span>
                </div>
                <h2 className="text-xl sm:text-2xl font-bold font-mono text-[#F5F5F7]">
                  Your coding weapon & career goals
                </h2>
                <p className="text-xs text-[#9292A3]">
                  Select your primary language and technical benchmark.
                </p>
              </div>

              {/* Language selection */}
              <div className="space-y-2">
                <label className="text-xs font-medium text-[#9292A3]">
                  Preferred Programming Language (Allowed: C, C++, Java, Python, Go)
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
                        className={`p-3 rounded-xl border text-center transition-all ${
                          isSelected
                            ? 'bg-[#8B5CF6]/20 border-[#8B5CF6] text-white shadow-md shadow-[#8B5CF6]/20'
                            : 'bg-[#101019] border-[#232336] text-[#9292A3] hover:border-[#373752]'
                        }`}
                      >
                        <div className="font-mono font-bold text-sm">{lang.name}</div>
                        <div className="text-[10px] text-[#65657B] mt-0.5">{lang.ext}</div>
                      </button>
                    );
                  })}
                </div>
              </div>

              {/* Learning goal */}
              <div className="space-y-2">
                <label className="text-xs font-medium text-[#9292A3]">Learning Goal</label>
                <div className="grid grid-cols-1 sm:grid-cols-3 gap-2.5">
                  {LEARNING_GOALS.map((goal) => {
                    const isSelected = formData.learning_goal === goal.id;
                    return (
                      <button
                        type="button"
                        key={goal.id}
                        onClick={() => setFormData({ ...formData, learning_goal: goal.id })}
                        className={`p-3 rounded-xl border text-left transition-all ${
                          isSelected
                            ? 'bg-[#8B5CF6]/20 border-[#8B5CF6] text-white shadow-sm'
                            : 'bg-[#101019] border-[#232336] text-[#9292A3] hover:border-[#373752]'
                        }`}
                      >
                        <div className="text-xs font-bold text-[#F5F5F7]">{goal.title}</div>
                        <div className="text-[11px] text-[#65657B] mt-1">{goal.desc}</div>
                      </button>
                    );
                  })}
                </div>
              </div>

              {/* DSA Level */}
              <div className="space-y-2">
                <label className="text-xs font-medium text-[#9292A3]">Current DSA Level (Allowed: Beginner, Intermediate, Advanced)</label>
                <div className="grid grid-cols-1 sm:grid-cols-3 gap-2.5">
                  {DSA_LEVELS.map((lvl) => {
                    const isSelected = formData.dsa_level === lvl.id;
                    return (
                      <button
                        type="button"
                        key={lvl.id}
                        onClick={() => setFormData({ ...formData, dsa_level: lvl.id as DSALevel })}
                        className={`p-3 rounded-xl border text-left transition-all ${
                          isSelected
                            ? 'bg-[#8B5CF6]/20 border-[#8B5CF6] text-white'
                            : 'bg-[#101019] border-[#232336] text-[#9292A3] hover:border-[#373752]'
                        }`}
                      >
                        <div className="text-xs font-bold text-[#F5F5F7]">{lvl.title}</div>
                        <div className="text-[10px] text-[#65657B] mt-1">{lvl.desc}</div>
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
                <h2 className="text-xl sm:text-2xl font-bold font-mono text-[#F5F5F7]">
                  Verify Your Curriculum Configuration
                </h2>
                <p className="text-xs text-[#9292A3]">
                  Your data will be securely synced with the LearnX backend and Supabase PostgreSQL.
                </p>
              </div>

              <div className="space-y-4">
                {/* Academic summary */}
                <div className="p-4 bg-[#101019] rounded-xl border border-[#232336] space-y-2">
                  <div className="flex items-center justify-between pb-2 border-b border-[#232336]">
                    <span className="text-xs font-bold text-[#F5F5F7] flex items-center gap-1.5">
                      <GraduationCap className="w-3.5 h-3.5 text-[#8B5CF6]" /> Academic Background
                    </span>
                    <button
                      onClick={() => setStep(1)}
                      className="text-xs text-[#8B5CF6] hover:text-[#A78BFA] flex items-center gap-1"
                    >
                      <Edit2 className="w-3 h-3" /> Edit
                    </button>
                  </div>
                  <div className="grid grid-cols-2 gap-2 text-xs">
                    <div>
                      <span className="text-[#65657B] block text-[11px]">College</span>
                      <span className="font-medium text-[#F5F5F7]">{formData.college}</span>
                    </div>
                    <div>
                      <span className="text-[#65657B] block text-[11px]">Degree</span>
                      <span className="font-medium text-[#F5F5F7]">{formData.course}</span>
                    </div>
                    <div>
                      <span className="text-[#65657B] block text-[11px]">Department</span>
                      <span className="font-medium text-[#F5F5F7]">{formData.department}</span>
                    </div>
                    <div>
                      <span className="text-[#65657B] block text-[11px]">Current Year</span>
                      <span className="font-medium text-[#F5F5F7]">{formData.current_year}</span>
                    </div>
                  </div>
                </div>

                {/* Preferences summary */}
                <div className="p-4 bg-[#101019] rounded-xl border border-[#232336] space-y-2">
                  <div className="flex items-center justify-between pb-2 border-b border-[#232336]">
                    <span className="text-xs font-bold text-[#F5F5F7] flex items-center gap-1.5">
                      <Target className="w-3.5 h-3.5 text-[#8B5CF6]" /> Learning Calibration
                    </span>
                    <button
                      onClick={() => setStep(2)}
                      className="text-xs text-[#8B5CF6] hover:text-[#A78BFA] flex items-center gap-1"
                    >
                      <Edit2 className="w-3 h-3" /> Edit
                    </button>
                  </div>
                  <div className="grid grid-cols-3 gap-2 text-xs">
                    <div>
                      <span className="text-[#65657B] block text-[11px]">Preferred Language</span>
                      <Badge variant="purple" size="sm" className="mt-1">
                        {formData.preferred_language}
                      </Badge>
                    </div>
                    <div>
                      <span className="text-[#65657B] block text-[11px]">Goal</span>
                      <span className="font-medium text-[#F5F5F7] block mt-1">
                        {formData.learning_goal}
                      </span>
                    </div>
                    <div>
                      <span className="text-[#65657B] block text-[11px]">DSA Level</span>
                      <span className="font-medium text-[#F5F5F7] block mt-1">
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
        <div className="mt-8 pt-5 border-t border-[#232336] flex items-center justify-between">
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
