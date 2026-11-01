import React, { forwardRef, useState } from 'react';
import { Eye, EyeOff } from 'lucide-react';
import { clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

export interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
  hint?: string;
  leftIcon?: React.ReactNode;
  isPassword?: boolean;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, hint, leftIcon, isPassword, className, type = 'text', id, ...props }, ref) => {
    const [showPassword, setShowPassword] = useState(false);
    const inputId = id || (label ? label.toLowerCase().replace(/\s+/g, '-') : undefined);
    const resolvedType = isPassword ? (showPassword ? 'text' : 'password') : type;

    return (
      <div className="w-full flex flex-col gap-1.5">
        {label && (
          <label htmlFor={inputId} className="text-xs font-medium text-[#9292A3]">
            {label}
          </label>
        )}
        <div className="relative flex items-center">
          {leftIcon && (
            <div className="absolute left-3.5 text-[#9292A3] pointer-events-none flex items-center">
              {leftIcon}
            </div>
          )}
          <input
            id={inputId}
            ref={ref}
            type={resolvedType}
            className={twMerge(
              clsx(
                'w-full bg-[#101019] text-[#F5F5F7] placeholder-[#65657B] text-sm rounded-xl px-4 py-2.5 border border-[#232336] transition-all duration-200 outline-none',
                'focus:border-[#8B5CF6] focus:ring-1 focus:ring-[#8B5CF6]/50',
                leftIcon && 'pl-10',
                isPassword && 'pr-11',
                error && 'border-red-500/80 focus:border-red-500 focus:ring-red-500/30',
                className
              )
            )}
            {...props}
          />
          {isPassword && (
            <button
              type="button"
              tabIndex={-1}
              onClick={() => setShowPassword(!showPassword)}
              className="absolute right-3.5 text-[#9292A3] hover:text-[#F5F5F7] p-1 transition-colors"
            >
              {showPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
            </button>
          )}
        </div>
        {error && <span className="text-xs text-red-400 font-medium">{error}</span>}
        {hint && !error && <span className="text-xs text-[#65657B]">{hint}</span>}
      </div>
    );
  }
);

Input.displayName = 'Input';
export default Input;

