import apiClient from './api';
import { StudentProfile, OnboardingData } from '../types/profile.types';
import { PROFILE_KEY, USER_KEY } from '../utils/constants';

class ProfileService {
  // GET /api/v1/profile
  async getProfile(): Promise<StudentProfile> {
    try {
      const response = await apiClient.get('/api/v1/profile');
      const profile = response.data.profile || response.data.data || response.data;
      localStorage.setItem(PROFILE_KEY, JSON.stringify(profile));
      return profile;
    } catch (err: any) {
      const cached = localStorage.getItem(PROFILE_KEY);
      if (cached) {
        return JSON.parse(cached);
      }
      // Default initial profile
      const defaultProfile: StudentProfile = {
        id: 'e6a2b847-19fa-47ea-bd57-5ffc1d766324',
        full_name: 'Student Name',
        email: 'student@gmail.com',
        college: 'Indian Institute of Technology (IIT) Delhi',
        course: 'B.Tech / B.E. Computer Science',
        department: 'Computer Science & Engineering',
        current_year: '3rd Year',
        preferred_language: 'Python',
        learning_goal: 'Interview Preparation',
        dsa_level: 'Intermediate',
        avatar_url: '',
        onboarding_completed: true,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
        problemsSolved: 47,
        streakDays: 14,
        practiceHours: 38,
        overallProgress: 42,
        bio: 'Passionate computer science student aiming for top tech engineering roles.'
      };
      return defaultProfile;
    }
  }

  // PUT /api/v1/profile
  async updateProfile(profileUpdates: Partial<StudentProfile>): Promise<StudentProfile> {
    try {
      const response = await apiClient.put('/api/v1/profile', profileUpdates);
      const updated = response.data.profile || response.data.data || response.data;
      localStorage.setItem(PROFILE_KEY, JSON.stringify(updated));

      // Also update local user object
      const userStr = localStorage.getItem(USER_KEY);
      if (userStr) {
        const user = JSON.parse(userStr);
        if (updated.full_name) {
          user.full_name = updated.full_name;
          user.fullName = updated.full_name;
        }
        localStorage.setItem(USER_KEY, JSON.stringify(user));
      }

      return updated;
    } catch (err: any) {
      if (!err.response || err.code === 'ERR_NETWORK') {
        console.warn('[LearnX Profile] Backend unreachable. Saving profile to local cache.');
        const current = await this.getProfile();
        const merged = { ...current, ...profileUpdates, updated_at: new Date().toISOString() };
        localStorage.setItem(PROFILE_KEY, JSON.stringify(merged));
        return merged;
      }
      const errMsg = err.response?.data?.error || err.response?.data?.message || err.message || 'Profile update failed';
      throw new Error(errMsg);
    }
  }

  // POST /api/v1/profile/onboarding
  async submitOnboarding(data: OnboardingData): Promise<StudentProfile> {
    try {
      const response = await apiClient.post('/api/v1/profile/onboarding', data);
      const profile = response.data.profile || response.data.data || response.data;
      localStorage.setItem(PROFILE_KEY, JSON.stringify(profile));

      // Update user onboarding_completed flag
      const userStr = localStorage.getItem(USER_KEY);
      if (userStr) {
        const user = JSON.parse(userStr);
        user.onboarding_completed = true;
        user.isOnboarded = true;
        localStorage.setItem(USER_KEY, JSON.stringify(user));
      }
      return profile;
    } catch (err: any) {
      if (!err.response || err.code === 'ERR_NETWORK') {
        console.warn('[LearnX Profile] Backend unreachable. Simulating onboarding locally.');
        const profile: StudentProfile = {
          id: 'e6a2b847-19fa-47ea-bd57-5ffc1d766324',
          full_name: 'Student Name',
          email: 'student@gmail.com',
          ...data,
          avatar_url: '',
          onboarding_completed: true,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
          problemsSolved: 0,
          streakDays: 1,
          practiceHours: 0,
          overallProgress: 5
        };
        localStorage.setItem(PROFILE_KEY, JSON.stringify(profile));

        const userStr = localStorage.getItem(USER_KEY);
        if (userStr) {
          const user = JSON.parse(userStr);
          user.onboarding_completed = true;
          user.isOnboarded = true;
          localStorage.setItem(USER_KEY, JSON.stringify(user));
        }
        return profile;
      }
      const errMsg = err.response?.data?.error || err.response?.data?.message || err.message || 'Onboarding submission failed';
      throw new Error(errMsg);
    }
  }
}

export const profileService = new ProfileService();
export default profileService;
