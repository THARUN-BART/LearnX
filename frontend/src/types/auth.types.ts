export interface User {
  id: string;
  email: string;
  full_name: string;
  fullName?: string; // alias for UI convenience
  role?: string;
  onboarding_completed: boolean;
  isOnboarded?: boolean; // alias for UI convenience
  avatar_url?: string;
  avatarUrl?: string;
  created_at?: string;
}

export interface Session {
  access_token: string;
  refresh_token: string;
  token_type: string;
  expires_in: number;
}

export interface AuthResponse {
  success: boolean;
  message?: string;
  user: User;
  session: Session;
}

export interface SignupRequest {
  full_name: string;
  email: string;
  password: string;
  confirm_password: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface VerifyOtpRequest {
  email: string;
  code: string;
}

export interface ForgotPasswordRequest {
  email: string;
}

export interface VerifyResetOtpRequest {
  email: string;
  code: string;
}

export interface ResetPasswordRequest {
  reset_token: string;
  new_password: string;
  confirm_new_password: string;
}
