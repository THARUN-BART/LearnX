import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';

// Layouts
import PublicLayout from '../layouts/PublicLayout';
import AuthLayout from '../layouts/AuthLayout';
import DashboardLayout from '../layouts/DashboardLayout';

// Route Guards
import ProtectedRoute from './ProtectedRoute';
import PublicRoute from './PublicRoute';

// Pages - Landing & Auth
import LandingPage from '../pages/landing/LandingPage';
import SignupPage from '../pages/auth/SignupPage';
import LoginPage from '../pages/auth/LoginPage';
import OTPVerificationPage from '../pages/auth/OTPVerificationPage';
import ForgotPasswordPage from '../pages/auth/ForgotPasswordPage';
import ResetPasswordPage from '../pages/auth/ResetPasswordPage';
import OAuthCallbackPage from '../pages/auth/OAuthCallbackPage';

// Pages - Onboarding
import OnboardingPage from '../pages/onboarding/OnboardingPage';

// Pages - Dashboard & Learning
import DashboardPage from '../pages/dashboard/DashboardPage';
import LearningPathPage from '../pages/learning-path/LearningPathPage';
import LearningMemoryPage from '../pages/learning-memory/LearningMemoryPage';
import ConceptDetailPage from '../pages/learning-memory/ConceptDetailPage';
import RoadmapPage from '../pages/roadmap/RoadmapPage';
import TopicDetailPage from '../pages/roadmap/TopicDetailPage';
import PracticePage from '../pages/practice/PracticePage';
import AssistantPage from '../pages/assistant/AssistantPage';
import MistakeAnalysisPage from '../pages/analysis/MistakeAnalysisPage';
import AnalyticsPage from '../pages/analytics/AnalyticsPage';
import ProfilePage from '../pages/profile/ProfilePage';
import SettingsPage from '../pages/profile/SettingsPage';
import NotFoundPage from '../pages/NotFoundPage';

export const AppRoutes: React.FC = () => {
  return (
    <Routes>
      {/* Landing Page (Public Layout) */}
      <Route element={<PublicLayout />}>
        <Route path="/" element={<LandingPage />} />
      </Route>

      {/* Auth Pages (Auth Layout with PublicRoute Guard) */}
      <Route
        element={
          <PublicRoute>
            <AuthLayout />
          </PublicRoute>
        }
      >
        <Route path="/signup" element={<SignupPage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/verify-signup" element={<OTPVerificationPage />} />
        <Route path="/forgot-password" element={<ForgotPasswordPage />} />
        <Route path="/verify-reset" element={<OTPVerificationPage />} />
        <Route path="/reset-password" element={<ResetPasswordPage />} />
        <Route path="/auth/callback" element={<OAuthCallbackPage />} />
      </Route>

      {/* Onboarding (Protected, requires login, but not yet onboarded) */}
      <Route
        path="/onboarding"
        element={
          <ProtectedRoute requiresOnboarded={false}>
            <OnboardingPage />
          </ProtectedRoute>
        }
      />

      {/* Protected Dashboard & Learning Modules (Dashboard Layout) */}
      <Route
        element={
          <ProtectedRoute requiresOnboarded={true}>
            <DashboardLayout />
          </ProtectedRoute>
        }
      >
        <Route path="/dashboard" element={<DashboardPage />} />
        <Route path="/dashboard/learning-path" element={<LearningPathPage />} />
        <Route path="/dashboard/learning-memory" element={<LearningMemoryPage />} />
        <Route path="/dashboard/learning-memory/concepts/:conceptKey" element={<ConceptDetailPage />} />
        <Route path="/roadmap" element={<RoadmapPage />} />
        <Route path="/roadmap/:topicId" element={<TopicDetailPage />} />
        <Route path="/practice" element={<PracticePage />} />
        <Route path="/practice/:problemId" element={<PracticePage />} />
        <Route path="/assistant" element={<AssistantPage />} />
        <Route path="/analysis" element={<MistakeAnalysisPage />} />
        <Route path="/analytics" element={<AnalyticsPage />} />
        <Route path="/profile" element={<ProfilePage />} />
        <Route path="/settings" element={<SettingsPage />} />
      </Route>

      {/* 404 Not Found Catch-All */}
      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  );
};

export default AppRoutes;

