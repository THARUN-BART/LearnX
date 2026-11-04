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

  return (
    <div className="w-full space-y-6">
      <div className="space-y-1.5">
        <h1 className="text-2xl font-bold tracking-tight text-[#0F172A] font-mono">
          Welcome Back
        </h1>
        <p className="text-xs text-[#64748B]">
          Enter your credentials to access your adaptive learning cockpit.
        </p>
      </div>

      {serverError && (
        <div className="p-3 text-xs bg-red-500/10 border border-red-500/20 text-red-400 rounded-xl">
          {serverError}
        </div>
      )}

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
            <label className="flex items-center gap-2 cursor-pointer select-none text-xs text-[#64748B]">
              <input
                type="checkbox"
                className="w-3.5 h-3.5 rounded bg-white border-[#E2E8F0] text-[#EF4444] focus:ring-0 focus:ring-offset-0 accent-[#EF4444]"
                {...register('rememberMe')}
              />
              <span>Remember me</span>
            </label>
            <Link
              to="/forgot-password"
              className="text-xs text-[#EF4444] hover:text-[#DC2626] transition-colors"
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

      <p className="text-center text-xs text-[#64748B]">
        New to LearnX?{' '}
        <Link to="/signup" className="text-[#EF4444] hover:text-[#DC2626] font-medium transition-colors">
          Create an account
        </Link>
      </p>
    </div>
  );
};

export default LoginPage;

