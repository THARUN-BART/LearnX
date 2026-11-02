import axios, { AxiosError, InternalAxiosRequestConfig } from 'axios';
import { API_BASE_URL, TOKEN_KEY } from '../utils/constants';

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 10000,
});

// Attach JWT access token to outgoing requests
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem(TOKEN_KEY);
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error: AxiosError) => Promise.reject(error)
);

// Response interceptor for token refresh & unauthorized handling
apiClient.interceptors.response.use(
  (response) => response,
  (error: AxiosError) => {
    if (error.response?.status === 401) {
      const currentPath = window.location.pathname;
      const reqUrl = error.config?.url || '';
      const hadAuth = !!error.config?.headers?.Authorization;

      // Only purge token if request actually carried an Authorization header and was rejected
      // (avoid clearing when unauthenticated probing happens on catalog endpoints)
      if (hadAuth && !currentPath.includes('/login') && !currentPath.includes('/signup')) {
        if (!reqUrl.includes('/roadmaps')) {
          localStorage.removeItem(TOKEN_KEY);
        }
      }
    }
    return Promise.reject(error);
  }
);

export default apiClient;

