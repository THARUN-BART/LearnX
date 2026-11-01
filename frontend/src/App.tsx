import React from 'react';
import { BrowserRouter } from 'react-router-dom';
import { Toaster } from 'sonner';
import { AuthProvider } from './contexts/AuthContext';
import { ProfileProvider } from './contexts/ProfileContext';
import AppRoutes from './routes/AppRoutes';

export const App: React.FC = () => {
  return (
    <BrowserRouter>
      <AuthProvider>
        <ProfileProvider>
          <AppRoutes />
          <Toaster
            position="top-right"
            theme="dark"
            richColors
            toastOptions={{
              style: {
                background: '#171722',
                border: '1px solid #232336',
                color: '#F5F5F7',
                fontFamily: 'Inter, sans-serif'
              }
            }}
          />
        </ProfileProvider>
      </AuthProvider>
    </BrowserRouter>
  );
};

export default App;
