import React from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import { Loader2 } from 'lucide-react';
import { USER_KEY, PROFILE_KEY } from '../utils/constants';

interface ProtectedRouteProps {
  children: React.ReactNode;
  requiresOnboarded?: boolean;
}

export const ProtectedRoute: React.FC<ProtectedRouteProps> = ({
  children,
  requiresOnboarded = true,
}) => {
  const { isAuthenticated, isLoading, user } = useAuth();
  const location = useLocation();

  if (isLoading) {
    return (
      <div className="min-h-screen bg-[#F8FAFC] flex flex-col items-center justify-center gap-3 text-[#64748B]">
        <Loader2 className="w-8 h-8 animate-spin text-[#EF4444]" />
        <span className="text-xs font-mono tracking-wide">Restoring secure session...</span>
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  // Resilient check across user state, localStorage user, and localStorage profile
  let isOnboarded = Boolean(user?.onboarding_completed ?? user?.isOnboarded);
  if (!isOnboarded) {
    try {
      const storedUser = localStorage.getItem(USER_KEY);
      if (storedUser) {
        const u = JSON.parse(storedUser);
        if (u.onboarding_completed || u.isOnboarded) isOnboarded = true;
      }
      const storedProfile = localStorage.getItem(PROFILE_KEY);
      if (storedProfile) {
        const p = JSON.parse(storedProfile);
        if (p.onboarding_completed || p.college) isOnboarded = true;
      }
    } catch {}
  }

  // If attempting to access an onboarding-dependent route and user is genuinely not onboarded
  if (requiresOnboarded && !isOnboarded && location.pathname !== '/onboarding') {
    return <Navigate to="/onboarding" replace />;
  }

  return <>{children}</>;
};

export default ProtectedRoute;
