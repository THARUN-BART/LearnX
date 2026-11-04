import React from 'react';

export const Skeleton: React.FC<{ className?: string }> = ({ className = '' }) => (
  <div className={`animate-pulse rounded-xl bg-[#E2E8F0] border border-[#E2E8F0] ${className}`} />
);

export const PageSkeleton: React.FC<{ rows?: number }> = ({ rows = 4 }) => (
  <div className="space-y-4">
    <Skeleton className="h-24 w-full" />
    <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
      {Array.from({ length: 4 }).map((_, i) => (
        <Skeleton key={i} className="h-24" />
      ))}
    </div>
    {Array.from({ length: rows }).map((_, i) => (
      <Skeleton key={`r-${i}`} className="h-40 w-full" />
    ))}
  </div>
);

export default Skeleton;
