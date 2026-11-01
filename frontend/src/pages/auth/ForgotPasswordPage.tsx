import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { Mail, ArrowRight, KeyRound } from 'lucide-react';
import { forgotPasswordSchema, ForgotPasswordFormData } from '../../utils/validation';
import authService from '../../services/auth.service';
import Input from '../../components/ui/Input';
import Button from '../../components/ui/Button';
import { toast } from 'sonner';

export const ForgotPasswordPage: React.FC = () => {
  const navigate = useNavigate();
  const [serverError, setServerError] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting }
  } = useForm<ForgotPasswordFormData>({
    resolver: zodResolver(forgotPasswordSchema),
    mode: 'onBlur',
    defaultValues: {
      email: ''
    }
  });

  const onSubmit = async (data: ForgotPasswordFormData) => {
    setServerError(null);
    try {
      await authService.forgotPassword(data);
      localStorage.setItem('learnx_reset_email', data.email);
      toast.success('Reset verification code sent to your email.');
      navigate('/verify-reset');
    } catch (err: any) {
      setServerError(err?.message || 'Failed to dispatch reset code. Please check your email.');
      toast.error('Could not send code');
    }
  };

  return (
    <div className="w-full space-y-6">
      <div className="space-y-1.5 text-center sm:text-left">
        <div className="w-10 h-10 rounded-xl bg-[#8B5CF6]/15 border border-[#8B5CF6]/30 text-[#A78BFA] flex items-center justify-center mb-3">
          <KeyRound className="w-5 h-5" />
        </div>
        <h1 className="text-2xl font-bold tracking-tight text-[#F5F5F7] font-mono">
          Reset Your Password
        </h1>
        <p className="text-xs text-[#9292A3]">
          Enter your registered email address and we'll dispatch a 6-digit verification code.
        </p>
      </div>

      {serverError && (
        <div className="p-3 text-xs bg-red-500/10 border border-red-500/20 text-red-400 rounded-xl">
          {serverError}
        </div>
      )}

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        <Input
          label="Registered Email"
          type="email"
          placeholder="student@university.edu"
          leftIcon={<Mail className="w-4 h-4" />}
          error={errors.email?.message}
          {...register('email')}
        />

        <Button
          type="submit"
          className="w-full mt-2"
          size="md"
          isLoading={isSubmitting}
          rightIcon={<ArrowRight className="w-4 h-4" />}
        >
          Send Verification Code
        </Button>
      </form>

      <p className="text-center text-xs text-[#9292A3]">
        Remembered your password?{' '}
        <Link to="/login" className="text-[#8B5CF6] hover:text-[#A78BFA] font-medium transition-colors">
          Return to login
        </Link>
      </p>
    </div>
  );
};

export default ForgotPasswordPage;

