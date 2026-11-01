/**
 * ProfileContext
 *
 * Provides the student profile to the whole component tree through a React
 * context so pages never have to call profileService.getProfile() themselves
 * unless they need a forced refresh.
 *
 * Cache strategy (mirrors profileCache.ts):
 *  - On mount, serve the cached version immediately (zero-latency first
 *    render), then silently revalidate in the background.
 *  - After a successful backend fetch the context state is updated, which
 *    triggers a re-render only where the profile is actually consumed.
 *  - `refreshProfile(forceRefresh?)` can be called from any child to
 *    invalidate and re-fetch (e.g. after the user edits their profile).
 *  - `patchProfile(updates)` applies an optimistic local update and then
 *    persists it to the backend. The UI stays responsive even on slow
 *    connections.
 */

import React, {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useRef,
  useState,
} from 'react';
import { StudentProfile } from '../types/profile.types';
import profileService from '../services/profile.service';
import { profileCache } from '../utils/profileCache';
import { useAuth } from './AuthContext';

// ---------------------------------------------------------------------------
// Context shape
// ---------------------------------------------------------------------------
interface ProfileContextType {
  /** The student profile; null while the first fetch is in-flight. */
  profile: StudentProfile | null;
  /** True only during the very first fetch (no cached data available yet). */
  isLoading: boolean;
  /** Non-null when the last fetch ended with an error and no cache was found. */
  error: string | null;
  /**
   * Trigger a profile fetch.
   * @param forceRefresh - bypass the TTL check and always hit the backend.
   */
  refreshProfile: (forceRefresh?: boolean) => Promise<void>;
  /**
   * Apply an optimistic local patch and persist it to the backend.
   * On network error the patch is kept locally so the UI remains consistent.
   */
  patchProfile: (updates: Partial<StudentProfile>) => Promise<void>;
}

const ProfileContext = createContext<ProfileContextType | undefined>(undefined);

// ---------------------------------------------------------------------------
// Provider
// ---------------------------------------------------------------------------
export const ProfileProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const { isAuthenticated } = useAuth();

  // Seed state from cache immediately so the UI has data on the first render
  const [profile, setProfile] = useState<StudentProfile | null>(
    () => profileCache.get()
  );
  const [isLoading, setIsLoading] = useState<boolean>(!profileCache.isFresh());
  const [error, setError] = useState<string | null>(null);

  // Prevent concurrent fetches
  const fetchingRef = useRef(false);

  // -------------------------------------------------------------------------
  const refreshProfile = useCallback(
    async (forceRefresh = false) => {
      if (fetchingRef.current) return; // deduplicate concurrent calls
      fetchingRef.current = true;

      // Only show the loader if we have nothing to show yet
      if (!profile) setIsLoading(true);
      setError(null);

      try {
        const fresh = await profileService.getProfile({ forceRefresh });
        setProfile(fresh);
      } catch (err: any) {
        const msg = err?.message ?? 'Failed to load profile';
        setError(msg);
        // Keep whatever we have in state – don't blank the UI
      } finally {
        setIsLoading(false);
        fetchingRef.current = false;
      }
    },
    [profile]
  );

  // -------------------------------------------------------------------------
  const patchProfile = useCallback(
    async (updates: Partial<StudentProfile>) => {
      // Optimistic update: apply immediately so the UI feels instant
      setProfile((prev) => {
        if (!prev) return prev;
        return { ...prev, ...updates, updated_at: new Date().toISOString() };
      });

      try {
        const saved = await profileService.updateProfile(updates);
        setProfile(saved);
      } catch (err: any) {
        // The service already applied the patch locally via profileCache.patch,
        // so we can stay with the optimistic state. Surface the error.
        const msg = err?.message ?? 'Profile update failed';
        setError(msg);
      }
    },
    []
  );

  // -------------------------------------------------------------------------
  // Fetch on mount (and whenever authentication changes)
  // -------------------------------------------------------------------------
  useEffect(() => {
    if (!isAuthenticated) {
      // Clear profile state on logout
      setProfile(null);
      return;
    }

    // If cache is fresh skip the network call but still warm state
    if (profileCache.isFresh()) {
      const cached = profileCache.get();
      if (cached) {
        setProfile(cached);
        setIsLoading(false);
        return;
      }
    }

    refreshProfile(false);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isAuthenticated]);

  // -------------------------------------------------------------------------
  return (
    <ProfileContext.Provider
      value={{ profile, isLoading, error, refreshProfile, patchProfile }}
    >
      {children}
    </ProfileContext.Provider>
  );
};

// ---------------------------------------------------------------------------
// Hook
// ---------------------------------------------------------------------------
export function useProfile(): ProfileContextType {
  const ctx = useContext(ProfileContext);
  if (!ctx) {
    throw new Error('useProfile must be used inside <ProfileProvider>');
  }
  return ctx;
}

export default ProfileContext;

