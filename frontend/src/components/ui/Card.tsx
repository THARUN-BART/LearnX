import React from 'react';
import { clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

interface CardProps extends React.HTMLAttributes<HTMLDivElement> {
  hoverable?: boolean;
  glow?: boolean;
}

export const Card: React.FC<CardProps> = ({
  children,
  className,
  hoverable = false,
  glow = false,
  ...props
}) => {
  return (
    <div
      className={twMerge(
        clsx(
          'bg-[#171722] border border-[#232336] rounded-2xl p-5 transition-all duration-200',
          hoverable && 'hover:border-[#373752] hover:bg-[#1A1A28] cursor-pointer hover:shadow-xl hover:shadow-black/40',
          glow && 'glow-purple',
          className
        )
      )}
      {...props}
    >
      {children}
    </div>
  );
};

export default Card;

