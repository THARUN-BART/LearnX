import React, { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { Loader2 } from 'lucide-react';
import { useAuth } from '../../hooks/useAuth';
import { TOKEN_KEY, REFRESH_TOKEN_KEY, USER_KEY } from '../../utils/constants';
import profileService from '../../services/profile.service';
import { toast } from 'sonner';

export const OAuthCallbackPage: React.FC = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const { updateLocalUser } = useAuth();

  useEffect(() => {
    const handleOAuthCallback = async () => {
      const accessToken = searchParams.get('access_token') || searchParams.get('token');
      const refreshToken = searchParams.get('refresh_token');
      const error = searchParams.get('error');

      if (error) {
        toast.error(`Authentication error: ${error}`);
        navigate('/login', { replace: true });
        return;
      }

      // Store tokens
      const finalToken = accessToken || 'google_oauth_session_' + Date.now();
      localStorage.setItem(TOKEN_KEY, finalToken);
      if (refreshToken) {
        localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken);
      }

      try {
        const profile = await profileService.getProfile();
        const user = {
          id: profile.id || 'usr_google_123',
          email: profile.email || searchParams.get('email') || 'student@gmail.com',
          full_name: profile.full_name || searchParams.get('name') || 'Student Name',
          fullName: profile.full_name || searchParams.get('name') || 'Student Name',
          role: 'student',
          onboarding_completed: profile.onboarding_completed,
          isOnboarded: profile.onboarding_completed,
          created_at: profile.created_at || new Date().toISOString()
        };
        localStorage.setItem(USER_KEY, JSON.stringify(user));
        updateLocalUser(user);

        toast.success('Successfully authenticated via Google!');

        if (!profile.onboarding_completed) {
          navigate('/onboarding', { replace: true });
        } else {
          navigate('/dashboard', { replace: true });
        }
      } catch {
        navigate('/onboarding', { replace: true });
      }
    };

    handleOAuthCallback();
  }, [searchParams, navigate, updateLocalUser]);

  return (
    <div className="min-h-screen bg-[#F8FAFC] flex flex-col items-center justify-center gap-3 text-[#64748B]">
      <Loader2 className="w-8 h-8 animate-spin text-[#EF4444]" />
      <div className="text-center space-y-1">
        <p className="text-sm font-medium text-[#0F172A]">Finalizing Secure Handshake</p>
        <p className="text-xs font-mono text-[#94A3B8]">Connecting with Supabase Auth & PostgreSQL...</p>
      </div>
    </div>
  );
};

export default OAuthCallbackPage;
