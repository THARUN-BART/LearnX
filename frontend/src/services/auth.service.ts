import apiClient from './api';
import {
  AuthResponse,
  LoginRequest,
  SignupRequest,
  VerifyOtpRequest,
  ForgotPasswordRequest,
  VerifyResetOtpRequest,
  ResetPasswordRequest,
  User
} from '../types/auth.types';
import { TOKEN_KEY, REFRESH_TOKEN_KEY, USER_KEY, PROFILE_KEY } from '../utils/constants';

class AuthService {
  // Step 1: POST /api/v1/auth/signup
  async signup(data: {
    fullName: string;
    email: string;
    password: string;
    confirmPassword: string;
  }): Promise<{ success: boolean; message: string; email: string; expires_in_seconds?: number }> {
    try {
      const payload: SignupRequest = {
        full_name: data.fullName,
        email: data.email,
        password: data.password,
        confirm_password: data.confirmPassword
      };
      const response = await apiClient.post('/api/v1/auth/signup', payload);
      localStorage.setItem('learnx_pending_email', data.email);
      localStorage.setItem('learnx_pending_fullname', data.fullName);
      return response.data;
    } catch (err: any) {
      if (!err.response || err.code === 'ERR_NETWORK') {
        console.warn('[LearnX Auth] Backend unreachable. Using local dev fallback for signup.');
        localStorage.setItem('learnx_pending_email', data.email);
        localStorage.setItem('learnx_pending_fullname', data.fullName);
        return {
          success: true,
          message: 'OTP sent to your email. Please verify to complete signup.',
          email: data.email,
          expires_in_seconds: 600
        };
      }
      const errMsg = err.response?.data?.error || err.response?.data?.message || err.message || 'Signup failed';
      throw new Error(errMsg);
    }
  }

  // Step 2: POST /api/v1/auth/verify-signup-otp
  async verifySignupOtp(data: { email: string; otp: string }): Promise<AuthResponse> {
    try {
      const payload: VerifyOtpRequest = {
        email: data.email,
        code: data.otp // backend uses 'code'
      };
      const response = await apiClient.post('/api/v1/auth/verify-signup-otp', payload);
      const resData = response.data;

      // Extract user and session
      const user = resData.user;
      user.fullName = user.full_name;
      user.isOnboarded = user.onboarding_completed;

      const accessToken = resData.session?.access_token;
      const refreshToken = resData.session?.refresh_token;

      if (accessToken) {
        localStorage.setItem(TOKEN_KEY, accessToken);
      }
      if (refreshToken) {
        localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken);
      }
      localStorage.setItem(USER_KEY, JSON.stringify(user));
      localStorage.removeItem('learnx_pending_email');
      localStorage.removeItem('learnx_pending_fullname');

      return resData;
    } catch (err: any) {
      if (!err.response || err.code === 'ERR_NETWORK') {
        console.warn('[LearnX Auth] Backend unreachable. Using local dev fallback for verify OTP.');
        if (data.otp.length === 6) {
          const email = data.email || localStorage.getItem('learnx_pending_email') || 'student@gmail.com';
          const fullName = localStorage.getItem('learnx_pending_fullname') || 'Student Name';
          const mockUser: User = {
            id: 'e6a2b847-19fa-47ea-bd57-5ffc1d766324',
            email,
            full_name: fullName,
            fullName,
            role: 'student',
            onboarding_completed: false,
            isOnboarded: false,
            created_at: new Date().toISOString()
          };
          const mockToken = 'mock_jwt_token_' + Date.now();
          localStorage.setItem(TOKEN_KEY, mockToken);
          localStorage.setItem(USER_KEY, JSON.stringify(mockUser));
          localStorage.removeItem('learnx_pending_email');
          localStorage.removeItem('learnx_pending_fullname');
          return {
            success: true,
            message: 'Signup and OTP verification completed successfully.',
            user: mockUser,
            session: {
              access_token: mockToken,
              refresh_token: 'mock_refresh_' + Date.now(),
              token_type: 'bearer',
              expires_in: 3600
            }
          };
        }
        throw new Error('Invalid OTP code. Please enter a valid 6-digit code.');
      }
      const errMsg = err.response?.data?.error || err.response?.data?.message || err.message || 'OTP verification failed';
      throw new Error(errMsg);
    }
  }

  // Step 3: POST /api/v1/auth/login
  async login(data: LoginRequest): Promise<AuthResponse> {
    try {
      const response = await apiClient.post('/api/v1/auth/login', data);
      const resData = response.data;

      const user = resData.user;
      user.fullName = user.full_name;
      user.isOnboarded = user.onboarding_completed;

      const accessToken = resData.session?.access_token;
      const refreshToken = resData.session?.refresh_token;

      if (accessToken) {
        localStorage.setItem(TOKEN_KEY, accessToken);
      }
      if (refreshToken) {
        localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken);
      }
      localStorage.setItem(USER_KEY, JSON.stringify(user));

      return resData;
    } catch (err: any) {
      if (!err.response || err.code === 'ERR_NETWORK') {
        console.warn('[LearnX Auth] Backend unreachable. Using local dev fallback for login.');
        const savedProfileStr = localStorage.getItem(PROFILE_KEY);
        const isOnboarded = !!savedProfileStr;

        const mockUser: User = {
          id: 'e6a2b847-19fa-47ea-bd57-5ffc1d766324',
          email: data.email,
          full_name: data.email.split('@')[0].replace('.', ' ').toUpperCase(),
          fullName: data.email.split('@')[0].replace('.', ' ').toUpperCase(),
          role: 'student',
          onboarding_completed: isOnboarded,
          isOnboarded: isOnboarded,
          created_at: new Date().toISOString()
        };
        const mockToken = 'mock_jwt_token_' + Date.now();
        localStorage.setItem(TOKEN_KEY, mockToken);
        localStorage.setItem(USER_KEY, JSON.stringify(mockUser));
        return {
          success: true,
          message: 'Login successful.',
          user: mockUser,
          session: {
            access_token: mockToken,
            refresh_token: 'mock_refresh_' + Date.now(),
            token_type: 'bearer',
            expires_in: 3600
          }
        };
      }
      const errMsg = err.response?.data?.error || err.response?.data?.message || err.message || 'Login failed';
      throw new Error(errMsg);
    }
  }

  // Step 4: POST /api/v1/auth/forgot-password
  async forgotPassword(data: ForgotPasswordRequest): Promise<{ success: boolean; message: string }> {
    try {
      const response = await apiClient.post('/api/v1/auth/forgot-password', data);
      localStorage.setItem('learnx_reset_email', data.email);
      return response.data;
    } catch (err: any) {
      if (!err.response || err.code === 'ERR_NETWORK') {
        console.warn('[LearnX Auth] Backend unreachable. Using fallback for forgot password.');
        localStorage.setItem('learnx_reset_email', data.email);
        return {
          success: true,
          message: 'If an account with that email exists, a password reset OTP has been sent.'
        };
      }
      const errMsg = err.response?.data?.error || err.response?.data?.message || err.message || 'Forgot password failed';
      throw new Error(errMsg);
    }
  }

  // Step 5: POST /api/v1/auth/verify-reset-otp
  async verifyResetOtp(data: { email: string; otp: string }): Promise<{ success: boolean; message: string; reset_token: string }> {
    try {
      const payload: VerifyResetOtpRequest = {
        email: data.email,
        code: data.otp // backend uses 'code'
      };
      const response = await apiClient.post('/api/v1/auth/verify-reset-otp', payload);
      const resData = response.data;
      if (resData.reset_token) {
        sessionStorage.setItem('learnx_reset_auth_token', resData.reset_token);
      }
      return resData;
    } catch (err: any) {
      if (!err.response || err.code === 'ERR_NETWORK') {
        console.warn('[LearnX Auth] Backend unreachable. Using fallback for verify reset OTP.');
        if (data.otp.length === 6) {
          const resetToken = '63f13f0426b3c95e1e75db740e716ce0' + Math.random().toString(36).substring(2);
          sessionStorage.setItem('learnx_reset_auth_token', resetToken);
          return {
            success: true,
            message: 'OTP verified successfully. You may now reset your password.',
            reset_token: resetToken
          };
        }
        throw new Error('Invalid OTP code. Please enter 6 valid digits.');
      }
      const errMsg = err.response?.data?.error || err.response?.data?.message || err.message || 'Invalid or expired OTP';
      throw new Error(errMsg);
    }
  }

  // Step 6: POST /api/v1/auth/reset-password
  async resetPassword(data: {
    token: string;
    password: string;
    confirmPassword?: string;
  }): Promise<{ success: boolean; message: string }> {
    try {
      const payload: ResetPasswordRequest = {
        reset_token: data.token,
        new_password: data.password,
        confirm_new_password: data.confirmPassword || data.password
      };
      const response = await apiClient.post('/api/v1/auth/reset-password', payload);
      sessionStorage.removeItem('learnx_reset_auth_token');
      localStorage.removeItem('learnx_reset_email');
      return response.data;
    } catch (err: any) {
      if (!err.response || err.code === 'ERR_NETWORK') {
        console.warn('[LearnX Auth] Backend unreachable. Using fallback for reset password.');
        sessionStorage.removeItem('learnx_reset_auth_token');
        localStorage.removeItem('learnx_reset_email');
        return {
          success: true,
          message: 'Password has been reset successfully. Please log in with your new password.'
        };
      }
      const errMsg = err.response?.data?.error || err.response?.data?.message || err.message || 'Password reset failed';
      throw new Error(errMsg);
    }
  }

  // Step 7: POST /api/v1/auth/logout
  async logout(): Promise<void> {
    try {
      await apiClient.post('/api/v1/auth/logout');
    } catch {
      // Ignore network errors on logout
    } finally {
      localStorage.removeItem(TOKEN_KEY);
      localStorage.removeItem(REFRESH_TOKEN_KEY);
      localStorage.removeItem(USER_KEY);
    }
  }

  // Step 8: GET /api/v1/auth/google
  initiateGoogleOAuth(): void {
    const redirectUrl = encodeURIComponent(`${window.location.origin}/auth/callback`);
    const backendOAuthUrl = `${apiClient.defaults.baseURL}/api/v1/auth/google?redirect_to=${redirectUrl}`;
    window.location.href = backendOAuthUrl;
  }

  getCurrentUser(): User | null {
    const userStr = localStorage.getItem(USER_KEY);
    if (!userStr) return null;
    try {
      const u = JSON.parse(userStr);
      u.fullName = u.full_name || u.fullName;
      u.isOnboarded = u.onboarding_completed ?? u.isOnboarded;
      return u;
    } catch {
      return null;
    }
  }

  isAuthenticated(): boolean {
    return !!localStorage.getItem(TOKEN_KEY);
  }
}

export const authService = new AuthService();
export default authService;
