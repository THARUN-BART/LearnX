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
          'bg-white border border-[#E2E8F0] rounded-2xl p-5 transition-all duration-200 shadow-sm',
          hoverable && 'hover:border-[#CBD5E1] hover:bg-[#F8FAFC] cursor-pointer hover:shadow-md',
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

