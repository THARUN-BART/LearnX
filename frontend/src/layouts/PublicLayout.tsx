import React from 'react';
import { Outlet } from 'react-router-dom';
import PublicNavbar from '../components/navigation/PublicNavbar';

export const PublicLayout: React.FC = () => {
  return (
    <div className="min-h-screen bg-[#F8FAFC] flex flex-col text-[#0F172A]">
      <PublicNavbar />
      <main className="flex-1 pt-18">
        <Outlet />
      </main>
    </div>
  );
};

export default PublicLayout;

