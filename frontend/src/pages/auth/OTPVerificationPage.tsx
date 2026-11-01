import React, { useState, useRef, useEffect } from 'react';
import { useNavigate, useLocation, Link } from 'react-router-dom';
import { ShieldCheck, ArrowRight, RotateCcw, CheckCircle2 } from 'lucide-react';
import { useAuth } from '../../hooks/useAuth';
import authService from '../../services/auth.service';
import Button from '../../components/ui/Button';
import { toast } from 'sonner';

export const OTPVerificationPage: React.FC = () => {
  const { verifySignupOtp, pendingEmail } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  const isResetMode = location.pathname.includes('verify-reset');
  const email =
    pendingEmail ||
    localStorage.getItem(isResetMode ? 'learnx_reset_email' : 'learnx_pending_email') ||
    'student@learnx.dev';

  const [otp, setOtp] = useState<string[]>(['', '', '', '', '', '']);
  const [activeInput, setActiveInput] = useState<number>(0);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [isSuccess, setIsSuccess] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [cooldown, setCooldown] = useState<number>(60);

  const inputRefs = useRef<(HTMLInputElement | null)[]>([]);

  // Auto focus first input on mount
  useEffect(() => {
    inputRefs.current[0]?.focus();
  }, []);

  // Cooldown countdown timer
  useEffect(() => {
    if (cooldown <= 0) return;
    const timer = setInterval(() => {
      setCooldown((prev) => prev - 1);
    }, 1000);
    return () => clearInterval(timer);
  }, [cooldown]);

  const handleChange = (index: number, val: string) => {
    // Only accept numeric
    const cleanVal = val.replace(/\D/g, '');
    const newOtp = [...otp];

    if (cleanVal.length > 1) {
      // Multiple characters entered / typed
      const digits = cleanVal.slice(0, 6).split('');
      for (let i = 0; i < 6; i++) {
        newOtp[i] = digits[i] || '';
      }
      setOtp(newOtp);
      const nextIndex = Math.min(digits.length, 5);
      inputRefs.current[nextIndex]?.focus();
      return;
    }

    newOtp[index] = cleanVal;
    setOtp(newOtp);
    setError(null);

    // Auto advance
    if (cleanVal && index < 5) {
      inputRefs.current[index + 1]?.focus();
      setActiveInput(index + 1);
    }
  };

  const handleKeyDown = (index: number, e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Backspace') {
      if (!otp[index] && index > 0) {
        inputRefs.current[index - 1]?.focus();
        setActiveInput(index - 1);
      }
    } else if (e.key === 'ArrowLeft' && index > 0) {
      inputRefs.current[index - 1]?.focus();
      setActiveInput(index - 1);
    } else if (e.key === 'ArrowRight' && index < 5) {
      inputRefs.current[index + 1]?.focus();
      setActiveInput(index + 1);
    }
  };

  const handlePaste = (e: React.ClipboardEvent) => {
    e.preventDefault();
    const pastedData = e.clipboardData.getData('text').trim().replace(/\D/g, '');
    if (!pastedData) return;

    const digits = pastedData.slice(0, 6).split('');
    const newOtp = [...otp];
    digits.forEach((digit, idx) => {
      if (idx < 6) newOtp[idx] = digit;
    });
    setOtp(newOtp);
    const focusIdx = Math.min(digits.length, 5);
    inputRefs.current[focusIdx]?.focus();
    setActiveInput(focusIdx);
  };

  const handleResend = async () => {
    if (cooldown > 0) return;
    try {
      if (isResetMode) {
        await authService.forgotPassword({ email });
      } else {
        await authService.signup({
          email,
          fullName: 'Student Coder',
          password: 'Password123',
          confirmPassword: 'Password123'
        });
      }
      setCooldown(60);
      toast.success('New 6-digit code has been sent!');
    } catch {
      toast.error('Failed to resend code');
    }
  };

  const handleVerify = async () => {
    const fullOtp = otp.join('');
    if (fullOtp.length !== 6) {
      setError('Please enter all 6 digits of the code.');
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      if (isResetMode) {
        const res = await authService.verifyResetOtp({ email, otp: fullOtp });
        setIsSuccess(true);
        toast.success(res.message || 'Code verified successfully');
        setTimeout(() => {
          navigate('/reset-password');
        }, 800);
      } else {
        await verifySignupOtp({ email, otp: fullOtp });
        setIsSuccess(true);
        toast.success('Account successfully verified!');
        setTimeout(() => {
          navigate('/onboarding');
        }, 800);
      }
    } catch (err: any) {
      setError(err?.message || 'Invalid or expired code. Please verify and try again.');
      toast.error('Verification failed');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="w-full space-y-6">
      <div className="space-y-1.5 text-center">
        <div className="w-12 h-12 rounded-2xl bg-[#8B5CF6]/15 border border-[#8B5CF6]/30 text-[#A78BFA] flex items-center justify-center mx-auto mb-3 shadow-lg shadow-[#8B5CF6]/10">
          {isSuccess ? (
            <CheckCircle2 className="w-6 h-6 text-emerald-400" />
          ) : (
            <ShieldCheck className="w-6 h-6" />
          )}
        </div>
        <h1 className="text-2xl font-bold tracking-tight text-[#F5F5F7] font-mono">
          {isResetMode ? 'Verify Reset Code' : 'Verify Your Email'}
        </h1>
        <p className="text-xs text-[#9292A3]">
          We have dispatched a 6-digit verification code to
        </p>
        <p className="text-xs font-mono font-medium text-[#A78BFA]">{email}</p>
      </div>

      {error && (
        <div className="p-3 text-xs bg-red-500/10 border border-red-500/20 text-red-400 rounded-xl text-center">
          {error}
        </div>
      )}

      {/* 6-Digit OTP Inputs */}
      <div className="flex justify-center gap-2 sm:gap-3 py-2" onPaste={handlePaste}>
        {otp.map((digit, idx) => (
          <input
            key={idx}
            ref={(el) => {
              inputRefs.current[idx] = el;
            }}
            type="text"
            inputMode="numeric"
            maxLength={1}
            value={digit}
            onChange={(e) => handleChange(idx, e.target.value)}
            onKeyDown={(e) => handleKeyDown(idx, e)}
            onFocus={() => setActiveInput(idx)}
            className={`w-11 h-13 sm:w-12 sm:h-14 text-center text-xl font-mono font-bold rounded-xl bg-[#101019] border transition-all duration-150 outline-none select-none ${
              error
                ? 'border-red-500 text-red-400'
                : digit
                ? 'border-[#8B5CF6] text-white bg-[#171722] shadow-sm shadow-[#8B5CF6]/20'
                : activeInput === idx
                ? 'border-[#8B5CF6] ring-1 ring-[#8B5CF6]/50'
                : 'border-[#232336] text-[#F5F5F7]'
            }`}
          />
        ))}
      </div>

      <Button
        onClick={handleVerify}
        className="w-full"
        size="md"
        isLoading={isLoading}
        disabled={otp.join('').length !== 6 || isSuccess}
        rightIcon={<ArrowRight className="w-4 h-4" />}
      >
        {isSuccess ? 'Verified!' : 'Verify Code'}
      </Button>

      {/* Resend cooldown & back navigation */}
      <div className="flex flex-col items-center gap-3 pt-2 text-xs text-[#9292A3]">
        <div className="flex items-center gap-1.5">
          <span>Didn't receive the code?</span>
          {cooldown > 0 ? (
            <span className="font-mono text-[#65657B]">Resend in {cooldown}s</span>
          ) : (
            <button
              onClick={handleResend}
              className="text-[#8B5CF6] hover:text-[#A78BFA] font-medium flex items-center gap-1 transition-colors"
            >
              <RotateCcw className="w-3 h-3" />
              <span>Resend OTP</span>
            </button>
          )}
        </div>

        <Link
          to={isResetMode ? '/forgot-password' : '/signup'}
          className="text-[#65657B] hover:text-[#9292A3] transition-colors"
        >
          ← Change email address
        </Link>
      </div>
    </div>
  );
};

export default OTPVerificationPage;
