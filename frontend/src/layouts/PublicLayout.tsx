import React from 'react';
import { Outlet } from 'react-router-dom';
import PublicNavbar from '../components/navigation/PublicNavbar';

export const PublicLayout: React.FC = () => {
  return (
    <div className="min-h-screen bg-[#08080D] flex flex-col text-[#F5F5F7]">
      <PublicNavbar />
      <main className="flex-1 pt-18">
        <Outlet />
      </main>
    </div>
  );
};

export default PublicLayout;

