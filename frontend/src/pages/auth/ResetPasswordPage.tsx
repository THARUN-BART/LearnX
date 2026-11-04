import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { Lock, ArrowRight, CheckCircle2 } from 'lucide-react';
import { resetPasswordSchema, ResetPasswordFormData, calculatePasswordStrength } from '../../utils/validation';
import authService from '../../services/auth.service';
import Input from '../../components/ui/Input';
import Button from '../../components/ui/Button';
import { toast } from 'sonner';

export const ResetPasswordPage: React.FC = () => {
  const navigate = useNavigate();
  const [serverError, setServerError] = useState<string | null>(null);

  const resetToken =
    sessionStorage.getItem('learnx_reset_auth_token') || 'rst_valid_token_default';

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors, isSubmitting }
  } = useForm<ResetPasswordFormData>({
    resolver: zodResolver(resetPasswordSchema),
    mode: 'onBlur',
    defaultValues: {
      password: '',
      confirmPassword: ''
    }
  });

  const passwordValue = watch('password', '');
  const strength = calculatePasswordStrength(passwordValue);

  const onSubmit = async (data: ResetPasswordFormData) => {
    setServerError(null);
    try {
      await authService.resetPassword({
        token: resetToken,
        password: data.password,
        confirmPassword: data.confirmPassword
      });
      toast.success('Password updated successfully! Please log in.');
      navigate('/login');
    } catch (err: any) {
      setServerError(err?.message || 'Password reset failed. Please start the process again.');
      toast.error('Reset failed');
    }
  };

  return (
    <div className="w-full space-y-6">
      <div className="space-y-1.5">
        <h1 className="text-2xl font-bold tracking-tight text-[#0F172A] font-mono">
          Set New Password
        </h1>
        <p className="text-xs text-[#64748B]">
          Create a secure password with at least 8 characters, numbers and uppercase letters.
        </p>
      </div>

      {serverError && (
        <div className="p-3 text-xs bg-red-500/10 border border-red-500/20 text-red-400 rounded-xl">
          {serverError}
        </div>
      )}

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        <div className="space-y-1.5">
          <Input
            label="New Password"
            isPassword
            placeholder="Min. 8 characters"
            leftIcon={<Lock className="w-4 h-4" />}
            error={errors.password?.message}
            {...register('password')}
          />
          {passwordValue && (
            <div className="pt-1 space-y-1">
              <div className="flex items-center justify-between text-[11px] text-[#64748B]">
                <span>Strength:</span>
                <span className="font-medium text-[#0F172A]">{strength.label}</span>
              </div>
              <div className="grid grid-cols-5 gap-1 h-1 w-full bg-white rounded-full overflow-hidden">
                {[1, 2, 3, 4, 5].map((level) => (
                  <div
                    key={level}
                    className={`h-full transition-all duration-300 ${
                      level <= strength.score ? strength.color : 'bg-transparent'
                    }`}
                  />
                ))}
              </div>
            </div>
          )}
        </div>

        <Input
          label="Confirm New Password"
          isPassword
          placeholder="Re-enter password"
          leftIcon={<Lock className="w-4 h-4" />}
          error={errors.confirmPassword?.message}
          {...register('confirmPassword')}
        />

        <Button
          type="submit"
          className="w-full mt-2"
          size="md"
          isLoading={isSubmitting}
          rightIcon={<ArrowRight className="w-4 h-4" />}
        >
          Update Password & Log In
        </Button>
      </form>

      <p className="text-center text-xs text-[#64748B]">
        <Link to="/login" className="text-[#EF4444] hover:text-[#DC2626] transition-colors">
          Cancel and return to login
        </Link>
      </p>
    </div>
  );
};

export default ResetPasswordPage;

