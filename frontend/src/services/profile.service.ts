import apiClient from './api';
import { StudentProfile, OnboardingData } from '../types/profile.types';
import { USER_KEY } from '../utils/constants';
import { profileCache } from '../utils/profileCache';

// ---------------------------------------------------------------------------
// Default fallback profile (used when both backend and cache are unavailable)
// ---------------------------------------------------------------------------
const DEFAULT_PROFILE: StudentProfile = {
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
  bio: 'Passionate computer science student aiming for top tech engineering roles.',
};

class ProfileService {
  // -------------------------------------------------------------------------
  // GET /api/v1/profile
  //
  // Strategy:
  //  1. If the in-memory / localStorage cache is still fresh → return it
  //     immediately (zero network round-trip).
  //  2. Otherwise fetch from the backend, update the cache, and return.
  //  3. On network error fall back to whatever is in the cache (may be
  //     stale) or to the hard-coded default profile so the UI always
  //     has something to render.
  //
  //  Pass `{ forceRefresh: true }` to bypass the freshness check and always
  //  hit the backend (useful after a profile update to confirm server state).
  // -------------------------------------------------------------------------
  async getProfile(opts: { forceRefresh?: boolean } = {}): Promise<StudentProfile> {
    // Return fresh cache immediately
    if (!opts.forceRefresh && profileCache.isFresh()) {
      const cached = profileCache.get();
      if (cached) return cached;
    }

    try {
      const response = await apiClient.get('/api/v1/profile');
      const profile: StudentProfile =
        response.data.profile ?? response.data.data ?? response.data;

      profileCache.set(profile);
      this._syncUserStore(profile);
      return profile;
    } catch {
      // Return stale cache on network error so UI never breaks
      const stale = profileCache.get();
      if (stale) return stale;

      // Absolute last resort
      return DEFAULT_PROFILE;
    }
  }

  // -------------------------------------------------------------------------
  // PUT /api/v1/profile
  // -------------------------------------------------------------------------
  async updateProfile(updates: Partial<StudentProfile>): Promise<StudentProfile> {
    try {
      const response = await apiClient.put('/api/v1/profile', updates);
      const updated: StudentProfile =
        response.data.profile ?? response.data.data ?? response.data;

      profileCache.set(updated);
      this._syncUserStore(updated);
      return updated;
    } catch (err: any) {
      // Network-only error: apply the update locally and keep working offline
      if (!err.response || err.code === 'ERR_NETWORK') {
        console.warn('[LearnX Profile] Backend unreachable – applying patch locally.');
        const patched = profileCache.patch(updates);
        if (patched) {
          this._syncUserStore(patched);
          return patched;
        }
      }

      const msg =
        err.response?.data?.error ??
        err.response?.data?.message ??
        err.message ??
        'Profile update failed';
      throw new Error(msg);
    }
  }

  // -------------------------------------------------------------------------
  // POST /api/v1/profile/onboarding
  // -------------------------------------------------------------------------
  async submitOnboarding(data: OnboardingData): Promise<StudentProfile> {
    try {
      const response = await apiClient.post('/api/v1/profile/onboarding', data);
      const profile: StudentProfile =
        response.data.profile ?? response.data.data ?? response.data;

      profile.onboarding_completed = true;
      profileCache.set(profile);
      this._syncUserStore(profile);
      return profile;
    } catch (err: any) {
      if (!err.response || err.code === 'ERR_NETWORK') {
        console.warn('[LearnX Profile] Backend unreachable – simulating onboarding locally.');
        const profile: StudentProfile = {
          ...DEFAULT_PROFILE,
          ...data,
          onboarding_completed: true,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
          problemsSolved: 0,
          streakDays: 1,
          practiceHours: 0,
          overallProgress: 5,
        };
        profileCache.set(profile);
        this._syncUserStore(profile);
        return profile;
      }

      const msg =
        err.response?.data?.error ??
        err.response?.data?.message ??
        err.message ??
        'Onboarding submission failed';
      throw new Error(msg);
    }
  }

  // -------------------------------------------------------------------------
  // Invalidate the profile cache (e.g. after logout)
  // -------------------------------------------------------------------------
  clearCache(): void {
    profileCache.clear();
  }

  // -------------------------------------------------------------------------
  // Private: keep the lightweight auth user store in sync with the profile
  // so components that read USER_KEY get up-to-date full_name, etc.
  // -------------------------------------------------------------------------
  private _syncUserStore(profile: StudentProfile): void {
    try {
      const raw = localStorage.getItem(USER_KEY);
      if (!raw) return;
      const user = JSON.parse(raw);
      if (profile.full_name) {
        user.full_name = profile.full_name;
        user.fullName = profile.full_name;
      }
      if (profile.onboarding_completed !== undefined) {
        user.onboarding_completed = profile.onboarding_completed;
        user.isOnboarded = profile.onboarding_completed;
      }
      localStorage.setItem(USER_KEY, JSON.stringify(user));
    } catch {
      // Non-critical – swallow
    }
  }
}

export const profileService = new ProfileService();
export default profileService;
