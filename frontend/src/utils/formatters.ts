export function formatDuration(minutes: number): string {
  const hours = Math.floor(minutes / 60);
  const mins = minutes % 60;
  if (hours === 0) return `${mins}m`;
  if (mins === 0) return `${hours}h`;
  return `${hours}h ${mins}m`;
}

export function normalizePercentage(val: number | undefined | null): number {
  if (val === undefined || val === null || isNaN(val)) return 0;
  // If it's a decimal ratio (0 < val <= 1), convert to percentage (e.g. 0.6 -> 60)
  // If it's already on a 0-100 scale (e.g. 60 or 92), keep it.
  const scaled = val > 0 && val <= 1 ? val * 100 : val;
  // Guarantee the percentage is strictly clamped between 0 and 100
  return Math.min(100, Math.max(0, Math.round(scaled)));
}

export function formatPercentage(val: number | undefined | null, decimals: number = 0): string {
  const normalized = normalizePercentage(val);
  return `${decimals > 0 ? normalized.toFixed(decimals) : normalized}%`;
}

export function formatDate(dateString: string): string {
  const date = new Date(dateString);
  return new Intl.DateTimeFormat('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric'
  }).format(date);
}

export function truncateText(text: string, maxLen: number = 100): string {
  if (text.length <= maxLen) return text;
  return text.substring(0, maxLen).trim() + '...';
}
