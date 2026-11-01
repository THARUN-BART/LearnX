/**
 * useProfile hook (re-export)
 *
 * Thin wrapper around ProfileContext. All existing call-sites keep working
 * unchanged, but now share a single cached profile instance via context.
 *
 * The hook also exposes `updateProfile` (alias for `patchProfile`) to
 * keep the legacy API surface intact.
 */
import { useProfile as useProfileContext } from '../contexts/ProfileContext';
import { StudentProfile } from '../types/profile.types';

export function useProfile() {
  const ctx = useProfileContext();

  /** Backwards-compat alias for patchProfile */
  const updateProfile = async (updates: Partial<StudentProfile>) => {
    await ctx.patchProfile(updates);
    return ctx.profile ? { ...ctx.profile, ...updates } : null;
  };

  return {
    profile: ctx.profile,
    isLoading: ctx.isLoading,
    error: ctx.error,
    refreshProfile: ctx.refreshProfile,
    updateProfile,
    patchProfile: ctx.patchProfile,
  };
}

export default useProfile;
