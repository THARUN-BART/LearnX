export type DSALevel = 'Beginner' | 'Intermediate' | 'Advanced';
export type ProgrammingLanguage = 'C' | 'C++' | 'Java' | 'Python' | 'Go';

export interface StudentProfile {
  id: string;
  full_name: string;
  email: string;
  college: string;
  course: string;
  department: string;
  current_year: string;
  preferred_language: ProgrammingLanguage;
  learning_goal: string;
  dsa_level: DSALevel;
  avatar_url: string;
  onboarding_completed: boolean;
  created_at: string;
  updated_at: string;
  // Aliases for UI convenience
  fullName?: string;
  currentYear?: string;
  preferredLanguage?: ProgrammingLanguage;
  learningGoal?: string;
  dsaLevel?: DSALevel;
  // UI metrics extensions
  problemsSolved?: number;
  streakDays?: number;
  practiceHours?: number;
  overallProgress?: number;
  bio?: string;
  githubUrl?: string;
  linkedinUrl?: string;
}

export interface OnboardingData {
  college: string;
  course: string;
  department: string;
  current_year: string;
  preferred_language: ProgrammingLanguage;
  learning_goal: string;
  dsa_level: DSALevel;
}
