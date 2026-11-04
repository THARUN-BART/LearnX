import React from 'react';
import { clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

interface BadgeProps {
  children: React.ReactNode;
  variant?: 'default' | 'easy' | 'medium' | 'hard' | 'purple' | 'indigo' | 'cyan' | 'outline';
  size?: 'sm' | 'md';
  className?: string;
  icon?: React.ReactNode;
}

export const Badge: React.FC<BadgeProps> = ({
  children,
  variant = 'default',
  size = 'md',
  className,
  icon
}) => {
  const variants = {
    default: 'bg-[#F1F5F9] text-[#475569] border border-[#E2E8F0]',
    easy: 'bg-emerald-50 text-emerald-600 border border-emerald-200',
    medium: 'bg-amber-50 text-amber-600 border border-amber-200',
    hard: 'bg-red-50 text-red-600 border border-red-200',
    purple: 'bg-red-50 text-red-600 border border-red-200',
    indigo: 'bg-rose-50 text-rose-600 border border-rose-200',
    cyan: 'bg-cyan-50 text-cyan-600 border border-cyan-200',
    outline: 'bg-transparent text-[#475569] border border-[#E2E8F0]'
  };

  const sizes = {
    sm: 'text-[11px] px-2 py-0.5 rounded-md font-medium tracking-wide',
    md: 'text-xs px-2.5 py-1 rounded-lg font-medium'
  };

  return (
    <span
      className={twMerge(
        clsx(
          'inline-flex items-center gap-1.5 font-mono select-none',
          variants[variant],
          sizes[size],
          className
        )
      )}
    >
      {icon}
      {children}
    </span>
  );
};

export default Badge;

