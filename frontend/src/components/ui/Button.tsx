import React from 'react';
import { Loader2 } from 'lucide-react';
import { clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'outline' | 'ghost' | 'danger';
  size?: 'sm' | 'md' | 'lg';
  isLoading?: boolean;
  leftIcon?: React.ReactNode;
  rightIcon?: React.ReactNode;
}

export const Button: React.FC<ButtonProps> = ({
  children,
  variant = 'primary',
  size = 'md',
  isLoading = false,
  leftIcon,
  rightIcon,
  className,
  disabled,
  ...props
}) => {
  const baseStyles = 'inline-flex items-center justify-center font-medium transition-all duration-200 select-none cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-[#08080D] active:scale-[0.98]';

  const variants = {
    primary: 'bg-[#8B5CF6] hover:bg-[#7C3AED] text-white shadow-lg shadow-[#8B5CF6]/20 focus:ring-[#8B5CF6]',
    secondary: 'bg-[#171722] hover:bg-[#1E1E2D] text-[#F5F5F7] border border-[#232336] hover:border-[#373752] focus:ring-[#6366F1]',
    outline: 'bg-transparent border border-[#373752] text-[#F5F5F7] hover:bg-[#171722] hover:border-[#8B5CF6]/60 focus:ring-[#8B5CF6]',
    ghost: 'bg-transparent hover:bg-[#171722] text-[#9292A3] hover:text-[#F5F5F7] focus:ring-zinc-700',
    danger: 'bg-red-500/10 text-red-400 hover:bg-red-500/20 border border-red-500/20 focus:ring-red-500'
  };

  const sizes = {
    sm: 'text-xs px-3 py-1.5 rounded-lg gap-1.5',
    md: 'text-sm px-4 py-2.5 rounded-xl gap-2',
    lg: 'text-base px-6 py-3 rounded-xl gap-2.5 font-semibold'
  };

  return (
    <button
      className={twMerge(clsx(baseStyles, variants[variant], sizes[size], className))}
      disabled={disabled || isLoading}
      {...props}
    >
      {isLoading ? <Loader2 className="w-4 h-4 animate-spin text-current" /> : leftIcon}
      {children}
      {!isLoading && rightIcon}
    </button>
  );
};

export default Button;

