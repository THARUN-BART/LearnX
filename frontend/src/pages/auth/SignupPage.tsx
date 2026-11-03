import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { User, Mail, Lock, Sparkles, ArrowRight } from 'lucide-react';
import { signupSchema, SignupFormData, calculatePasswordStrength } from '../../utils/validation';
import { useAuth } from '../../hooks/useAuth';
import Input from '../../components/ui/Input';
import Button from '../../components/ui/Button';
import authService from '../../services/auth.service';
import { toast } from 'sonner';

export const SignupPage: React.FC = () => {
  const { signup } = useAuth();
  const navigate = useNavigate();
  const [serverError, setServerError] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors, isSubmitting }
  } = useForm<SignupFormData>({
    resolver: zodResolver(signupSchema),
    mode: 'onBlur',
    defaultValues: {
      fullName: '',
      email: '',
      password: '',
      confirmPassword: ''
    }
  });

  const passwordValue = watch('password', '');
  const strength = calculatePasswordStrength(passwordValue);

  const onSubmit = async (data: SignupFormData) => {
    setServerError(null);
    try {
      await signup(data);
      toast.success('Account created! Please verify the code sent to your email.');
      navigate('/verify-signup');
    } catch (err: any) {
      setServerError(err?.message || 'Failed to create account. Please try again.');
      toast.error('Signup failed');
    }
  };

  return (
    <div className="w-full space-y-6">
      <div className="space-y-1.5">
        <h1 className="text-2xl font-bold tracking-tight text-[#F5F5F7] font-mono">
          Create an Account
        </h1>
        <p className="text-xs text-[#9292A3]">
          Start your personalized algorithmic mastery journey today.
        </p>
      </div>

      {serverError && (
        <div className="p-3 text-xs bg-red-500/10 border border-red-500/20 text-red-400 rounded-xl">
          {serverError}
        </div>
      )}

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        <Input
          label="Full Name"
          placeholder="e.g. Alex Chen"
          leftIcon={<User className="w-4 h-4" />}
          error={errors.fullName?.message}
          {...register('fullName')}
        />

        <Input
          label="Email Address"
          type="email"
          placeholder="alex@university.edu"
          leftIcon={<Mail className="w-4 h-4" />}
          error={errors.email?.message}
          {...register('email')}
        />

        <div className="space-y-1.5">
          <Input
            label="Password"
            isPassword
            placeholder="Min. 8 characters"
            leftIcon={<Lock className="w-4 h-4" />}
            error={errors.password?.message}
            {...register('password')}
          />
          {passwordValue && (
            <div className="pt-1 space-y-1">
              <div className="flex items-center justify-between text-[11px] text-[#9292A3]">
                <span>Strength:</span>
                <span className="font-medium text-[#F5F5F7]">{strength.label}</span>
              </div>
              <div className="grid grid-cols-5 gap-1 h-1 w-full bg-[#101019] rounded-full overflow-hidden">
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
          label="Confirm Password"
          isPassword
          placeholder="Re-enter your password"
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
          Create Account
        </Button>
      </form>

      <p className="text-center text-xs text-[#9292A3]">
        Already have an account?{' '}
        <Link to="/login" className="text-[#8B5CF6] hover:text-[#A78BFA] font-medium transition-colors">
          Log In
        </Link>
      </p>
    </div>
  );
};

export default SignupPage;

