import { useState, useCallback } from 'react';

export function useApi<T, P extends any[]>(apiFunc: (...args: P) => Promise<T>) {
  const [data, setData] = useState<T | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const execute = useCallback(
    async (...args: P) => {
      setIsLoading(true);
      setError(null);
      try {
        const result = await apiFunc(...args);
        setData(result);
        return result;
      } catch (err: any) {
        const msg = err?.message || 'An unexpected error occurred';
        setError(msg);
        throw err;
      } finally {
        setIsLoading(false);
      }
    },
    [apiFunc]
  );

  return { data, isLoading, error, execute, setData };
}

export default useApi;

