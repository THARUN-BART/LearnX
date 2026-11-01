import React, { createContext, useContext, useState, useEffect } from 'react';
import { User, LoginRequest } from '../types/auth.types';
import authService from '../services/auth.service';
import { TOKEN_KEY, USER_KEY, PROFILE_KEY } from '../utils/constants';
import { profileCache } from '../utils/profileCache';

interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  pendingEmail: string | null;
  setPendingEmail: (email: string | null) => void;
  login: (data: LoginRequest) => Promise<void>;
  signup: (data: {
    fullName: string;
    email: string;
    password: string;
    confirmPassword: string;
  }) => Promise<{ success: boolean; message: string; email: string }>;
  verifySignupOtp: (data: { email: string; otp: string }) => Promise<void>;
  logout: () => void;
  updateLocalUser: (user: Partial<User>) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [pendingEmail, setPendingEmail] = useState<string | null>(() => {
    return localStorage.getItem('learnx_pending_email');
  });

  useEffect(() => {
    // Restore authentication state
    const token = localStorage.getItem(TOKEN_KEY);
    const storedUser = localStorage.getItem(USER_KEY);
    const storedProfile = localStorage.getItem(PROFILE_KEY);

    if (token && storedUser) {
      try {
        const parsed = JSON.parse(storedUser);
        let isOnboarded = Boolean(parsed.onboarding_completed ?? parsed.isOnboarded);
        
        // Also check stored profile
        if (!isOnboarded && storedProfile) {
          try {
            const profile = JSON.parse(storedProfile);
            if (profile.onboarding_completed || profile.college) {
              isOnboarded = true;
            }
          } catch {}
        }

        parsed.fullName = parsed.full_name || parsed.fullName;
        parsed.onboarding_completed = isOnboarded;
        parsed.isOnboarded = isOnboarded;
        setUser(parsed);
      } catch (e) {
        localStorage.removeItem(TOKEN_KEY);
        localStorage.removeItem(USER_KEY);
      }
    }
    setIsLoading(false);
  }, []);

  const login = async (data: LoginRequest) => {
    setIsLoading(true);
    try {
      const response = await authService.login(data);
      const u = response.user;
      const storedProfile = localStorage.getItem(PROFILE_KEY);
      let isOnboarded = Boolean(u.onboarding_completed ?? u.isOnboarded);
      if (!isOnboarded && storedProfile) {
        try {
          const profile = JSON.parse(storedProfile);
          if (profile.onboarding_completed || profile.college) {
            isOnboarded = true;
          }
        } catch {}
      }

      u.fullName = u.full_name || u.fullName;
      u.onboarding_completed = isOnboarded;
      u.isOnboarded = isOnboarded;
      setUser(u);
    } finally {
      setIsLoading(false);
    }
  };

  const signup = async (data: {
    fullName: string;
    email: string;
    password: string;
    confirmPassword: string;
  }) => {
    setIsLoading(true);
    try {
      const res = await authService.signup(data);
      setPendingEmail(data.email);
      localStorage.setItem('learnx_pending_email', data.email);
      return res;
    } finally {
      setIsLoading(false);
    }
  };

  const verifySignupOtp = async (data: { email: string; otp: string }) => {
    setIsLoading(true);
    try {
      const res = await authService.verifySignupOtp(data);
      const u = res.user;
      u.fullName = u.full_name || u.fullName;
      u.isOnboarded = Boolean(u.onboarding_completed ?? u.isOnboarded);
      setUser(u);
      localStorage.removeItem('learnx_pending_email');
      setPendingEmail(null);
    } finally {
      setIsLoading(false);
    }
  };

  const logout = () => {
    authService.logout();
    profileCache.clear();
    setUser(null);
  };

  const updateLocalUser = (updatedProps: Partial<User>) => {
    setUser((prevUser) => {
      const base = prevUser || (localStorage.getItem(USER_KEY) ? JSON.parse(localStorage.getItem(USER_KEY)!) : {
        id: 'usr_default',
        email: 'student@gmail.com',
        full_name: 'Student Name'
      });

      const updated = {
        ...base,
        ...updatedProps,
        full_name: updatedProps.full_name || updatedProps.fullName || base.full_name || base.fullName,
        fullName: updatedProps.full_name || updatedProps.fullName || base.full_name || base.fullName,
        onboarding_completed: updatedProps.onboarding_completed ?? updatedProps.isOnboarded ?? base.onboarding_completed ?? true,
        isOnboarded: updatedProps.onboarding_completed ?? updatedProps.isOnboarded ?? base.isOnboarded ?? true
      };

      localStorage.setItem(USER_KEY, JSON.stringify(updated));
      return updated;
    });
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        isAuthenticated: !!user && !!localStorage.getItem(TOKEN_KEY),
        isLoading,
        pendingEmail,
        setPendingEmail,
        login,
        signup,
        verifySignupOtp,
        logout,
        updateLocalUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

export default AuthContext;
