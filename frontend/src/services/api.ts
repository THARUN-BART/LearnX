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
      // Clear token if expired or unauthorized
      const currentPath = window.location.pathname;
      if (!currentPath.includes('/login') && !currentPath.includes('/signup')) {
        localStorage.removeItem(TOKEN_KEY);
        // Do not force reload if already on public pages
      }
    }
    return Promise.reject(error);
  }
);

export default apiClient;

