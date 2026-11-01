export function formatDuration(minutes: number): string {
  const hours = Math.floor(minutes / 60);
  const mins = minutes % 60;
  if (hours === 0) return `${mins}m`;
  if (mins === 0) return `${hours}h`;
  return `${hours}h ${mins}m`;
}

export function formatPercentage(val: number, decimals: number = 0): string {
  return `${val.toFixed(decimals)}%`;
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

