import axios, { AxiosError, InternalAxiosRequestConfig } from 'axios';
import { API_BASE_URL, TOKEN_KEY } from '../utils/constants';

const AUTH_PATHS = [
  '/auth/login',
  '/auth/signup',
  '/auth/logout',
  '/auth/verify',
  '/profile',
];

function isAuthCritical(url = ''): boolean {
  return AUTH_PATHS.some((p) => url.includes(p));
}

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 20000,
});

export const aiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 60000,
});

function attachToken(config: InternalAxiosRequestConfig) {
  const token = localStorage.getItem(TOKEN_KEY);
  if (token && config.headers) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
}

apiClient.interceptors.request.use(attachToken, (error: AxiosError) => Promise.reject(error));
aiClient.interceptors.request.use(attachToken, (error: AxiosError) => Promise.reject(error));

function handleUnauthorized(error: AxiosError) {
  if (error.response?.status === 401) {
    const currentPath = window.location.pathname;
    const reqUrl = error.config?.url || '';
    const hadAuth = !!error.config?.headers?.Authorization;
    const onAuthPage =
      currentPath.includes('/login') ||
      currentPath.includes('/signup') ||
      currentPath.includes('/verify') ||
      currentPath.includes('/reset');

    if (hadAuth && !onAuthPage && isAuthCritical(reqUrl)) {
      localStorage.removeItem(TOKEN_KEY);
    }
  }
  return Promise.reject(error);
}

apiClient.interceptors.response.use((response) => response, handleUnauthorized);
aiClient.interceptors.response.use((response) => response, handleUnauthorized);

export function unwrap<T = any>(payload: any, ...keys: string[]): T {
  if (payload == null) return payload;
  const data = payload.data !== undefined && payload.success !== undefined ? payload : payload;
  for (const key of keys) {
    if (data?.[key] !== undefined) return data[key] as T;
  }
  if (data?.data !== undefined) return data.data as T;
  return data as T;
}

export default apiClient;
