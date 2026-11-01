import { z } from 'zod';

export const signupSchema = z
  .object({
    fullName: z
      .string()
      .min(2, 'Full name must be at least 2 characters')
      .max(60, 'Full name must be less than 60 characters'),
    email: z
      .string()
      .min(1, 'Email is required')
      .email('Please enter a valid email address'),
    password: z
      .string()
      .min(8, 'Password must be at least 8 characters')
      .regex(/[A-Z]/, 'Must include at least one uppercase letter')
      .regex(/[a-z]/, 'Must include at least one lowercase letter')
      .regex(/[0-9]/, 'Must include at least one number'),
    confirmPassword: z.string().min(1, 'Please confirm your password')
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: 'Passwords do not match',
    path: ['confirmPassword']
  });

export type SignupFormData = z.infer<typeof signupSchema>;

export const loginSchema = z.object({
  email: z
    .string()
    .min(1, 'Email is required')
    .email('Please enter a valid email address'),
  password: z
    .string()
    .min(1, 'Password is required'),
  rememberMe: z.boolean().optional()
});

export type LoginFormData = z.infer<typeof loginSchema>;

export const otpSchema = z.object({
  otp: z
    .string()
    .length(6, 'Verification code must be exactly 6 digits')
    .regex(/^\d+$/, 'Code must contain numbers only')
});

export type OtpFormData = z.infer<typeof otpSchema>;

export const forgotPasswordSchema = z.object({
  email: z
    .string()
    .min(1, 'Email is required')
    .email('Please enter a valid email address')
});

export type ForgotPasswordFormData = z.infer<typeof forgotPasswordSchema>;

export const resetPasswordSchema = z
  .object({
    password: z
      .string()
      .min(8, 'Password must be at least 8 characters')
      .regex(/[A-Z]/, 'Must include at least one uppercase letter')
      .regex(/[a-z]/, 'Must include at least one lowercase letter')
      .regex(/[0-9]/, 'Must include at least one number'),
    confirmPassword: z.string().min(1, 'Please confirm your password')
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: 'Passwords do not match',
    path: ['confirmPassword']
  });

export type ResetPasswordFormData = z.infer<typeof resetPasswordSchema>;

export const onboardingAcademicSchema = z.object({
  college: z.string().min(2, 'College name is required'),
  course: z.string().min(2, 'Course is required'),
  department: z.string().min(2, 'Department is required'),
  academicYear: z.string().min(1, 'Academic year is required')
});

export const onboardingPreferencesSchema = z.object({
  preferredLanguage: z.enum(['python', 'cpp', 'java', 'c', 'go'] as const),
  learningGoal: z.enum(['faang_prep', 'college_exams', 'competitive_programming', 'career_switch', 'foundations'] as const),
  dsaLevel: z.enum(['beginner', 'intermediate', 'advanced'] as const)
});

export function calculatePasswordStrength(pass: string): {
  score: number;
  label: string;
  color: string;
} {
  let score = 0;
  if (!pass) return { score: 0, label: 'Empty', color: 'bg-zinc-700' };
  if (pass.length >= 8) score += 1;
  if (pass.length >= 12) score += 1;
  if (/[A-Z]/.test(pass)) score += 1;
  if (/[0-9]/.test(pass)) score += 1;
  if (/[^A-Za-z0-9]/.test(pass)) score += 1;

  if (score <= 2) return { score, label: 'Weak', color: 'bg-red-500' };
  if (score <= 3) return { score, label: 'Fair', color: 'bg-amber-500' };
  if (score <= 4) return { score, label: 'Good', color: 'bg-indigo-400' };
  return { score, label: 'Strong', color: 'bg-emerald-500' };
}

