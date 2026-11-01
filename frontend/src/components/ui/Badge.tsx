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
    default: 'bg-[#232336] text-[#9292A3] border border-[#373752]/40',
    easy: 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20',
    medium: 'bg-amber-500/10 text-amber-400 border border-amber-500/20',
    hard: 'bg-red-500/10 text-red-400 border border-red-500/20',
    purple: 'bg-[#8B5CF6]/15 text-[#A78BFA] border border-[#8B5CF6]/30',
    indigo: 'bg-[#6366F1]/15 text-[#818CF8] border border-[#6366F1]/30',
    cyan: 'bg-cyan-500/15 text-cyan-400 border border-cyan-500/30',
    outline: 'bg-transparent text-[#9292A3] border border-[#232336]'
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

