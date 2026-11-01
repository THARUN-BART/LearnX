import { useState, useEffect, useCallback } from 'react';
import { StudentProfile } from '../types/profile.types';
import profileService from '../services/profile.service';
import { useAuth } from './useAuth';

export const useProfile = () => {
  const { isAuthenticated, updateLocalUser } = useAuth();
  const [profile, setProfile] = useState<StudentProfile | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const fetchProfile = useCallback(async () => {
    if (!isAuthenticated) {
      setIsLoading(false);
      return;
    }
    setIsLoading(true);
    setError(null);
    try {
      const data = await profileService.getProfile();
      setProfile(data);
    } catch (err: any) {
      setError(err?.message || 'Failed to fetch student profile');
    } finally {
      setIsLoading(false);
    }
  }, [isAuthenticated]);

  useEffect(() => {
    fetchProfile();
  }, [fetchProfile]);

  const updateProfile = async (updates: Partial<StudentProfile>) => {
    try {
      const updated = await profileService.updateProfile(updates);
      setProfile(updated);
      updateLocalUser({ isOnboarded: true });
      return updated;
    } catch (err: any) {
      throw err;
    }
  };

  return { profile, isLoading, error, refreshProfile: fetchProfile, updateProfile };
};

export default useProfile;

