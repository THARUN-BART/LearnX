import React, { useState } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { Mail, Lock, ArrowRight } from 'lucide-react';
import { loginSchema, LoginFormData } from '../../utils/validation';
import { useAuth } from '../../hooks/useAuth';
import authService from '../../services/auth.service';
import Input from '../../components/ui/Input';
import Button from '../../components/ui/Button';
import { toast } from 'sonner';

export const LoginPage: React.FC = () => {
  const { login } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const [serverError, setServerError] = useState<string | null>(null);

  const from = (location.state as any)?.from?.pathname || '/dashboard';

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting }
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
    mode: 'onBlur',
    defaultValues: {
      email: '',
      password: '',
      rememberMe: true
    }
  });

  const onSubmit = async (data: LoginFormData) => {
    setServerError(null);
    try {
      await login(data);
      toast.success('Welcome back to LearnX!');
      navigate(from, { replace: true });
    } catch (err: any) {
      setServerError(err?.message || 'Invalid email or password. Please try again.');
      toast.error('Authentication failed');
    }
  };

  const handleGoogleLogin = () => {
    authService.initiateGoogleOAuth();
  };

  return (
    <div className="w-full space-y-6">
      <div className="space-y-1.5">
        <h1 className="text-2xl font-bold tracking-tight text-[#F5F5F7] font-mono">
          Welcome Back
        </h1>
        <p className="text-xs text-[#9292A3]">
          Enter your credentials to access your adaptive learning cockpit.
        </p>
      </div>

      {serverError && (
        <div className="p-3 text-xs bg-red-500/10 border border-red-500/20 text-red-400 rounded-xl">
          {serverError}
        </div>
      )}

      {/* Google OAuth Button */}
      <button
        type="button"
        onClick={handleGoogleLogin}
        className="w-full flex items-center justify-center gap-3 px-4 py-2.5 rounded-xl bg-[#101019] hover:bg-[#1B1B26] border border-[#232336] text-xs font-medium text-[#F5F5F7] transition-all duration-200"
      >
        <svg className="w-4 h-4" viewBox="0 0 24 24">
          <path
            fill="#4285F4"
            d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
          />
          <path
            fill="#34A853"
            d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
          />
          <path
            fill="#FBBC05"
            d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.06H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.94l2.85-2.22.81-.63z"
          />
          <path
            fill="#EA4335"
            d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.06l3.66 2.84c.87-2.6 3.3-4.52 6.16-4.52z"
          />
        </svg>
        <span>Log in with Google</span>
      </button>

      <div className="relative flex items-center justify-center">
        <div className="border-t border-[#232336] w-full" />
        <span className="bg-[#171722] px-3 text-[11px] text-[#65657B] uppercase tracking-wider font-mono absolute">
          or with email
        </span>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        <Input
          label="Email Address"
          type="email"
          placeholder="student@university.edu"
          leftIcon={<Mail className="w-4 h-4" />}
          error={errors.email?.message}
          {...register('email')}
        />

        <div className="space-y-1">
          <Input
            label="Password"
            isPassword
            placeholder="••••••••"
            leftIcon={<Lock className="w-4 h-4" />}
            error={errors.password?.message}
            {...register('password')}
          />
          <div className="flex items-center justify-between pt-1">
            <label className="flex items-center gap-2 cursor-pointer select-none text-xs text-[#9292A3]">
              <input
                type="checkbox"
                className="w-3.5 h-3.5 rounded bg-[#101019] border-[#232336] text-[#8B5CF6] focus:ring-0 focus:ring-offset-0 accent-[#8B5CF6]"
                {...register('rememberMe')}
              />
              <span>Remember me</span>
            </label>
            <Link
              to="/forgot-password"
              className="text-xs text-[#8B5CF6] hover:text-[#A78BFA] transition-colors"
            >
              Forgot password?
            </Link>
          </div>
        </div>

        <Button
          type="submit"
          className="w-full mt-2"
          size="md"
          isLoading={isSubmitting}
          rightIcon={<ArrowRight className="w-4 h-4" />}
        >
          Sign In
        </Button>
      </form>

      <p className="text-center text-xs text-[#9292A3]">
        New to LearnX?{' '}
        <Link to="/signup" className="text-[#8B5CF6] hover:text-[#A78BFA] font-medium transition-colors">
          Create an account
        </Link>
      </p>
    </div>
  );
};

export default LoginPage;

