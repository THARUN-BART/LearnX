import React from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import { Loader2 } from 'lucide-react';

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
      <div className="min-h-screen bg-[#08080D] flex flex-col items-center justify-center gap-3 text-[#9292A3]">
        <Loader2 className="w-8 h-8 animate-spin text-[#8B5CF6]" />
        <span className="text-sm font-mono tracking-wide">Restoring secure session...</span>
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  const isOnboarded = Boolean(user?.onboarding_completed ?? user?.isOnboarded);

  // If user hasn't onboarded yet and is attempting to access a route that requires onboarding
  if (requiresOnboarded && user && !isOnboarded && location.pathname !== '/onboarding') {
    return <Navigate to="/onboarding" replace />;
  }

  return <>{children}</>;
};

export default ProtectedRoute;

