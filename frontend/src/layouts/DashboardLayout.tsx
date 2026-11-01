import React, { useState, useEffect } from 'react';
import { Outlet, useLocation } from 'react-router-dom';
import Sidebar from '../components/navigation/Sidebar';
import Navbar from '../components/navigation/Navbar';
import CommandPalette from '../components/common/CommandPalette';

export const DashboardLayout: React.FC = () => {
  const [mobileSidebarOpen, setMobileSidebarOpen] = useState(false);
  const [commandPaletteOpen, setCommandPaletteOpen] = useState(false);
  const location = useLocation();

  // Listen for Cmd+K or Ctrl+K globally
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
        e.preventDefault();
        setCommandPaletteOpen((prev) => !prev);
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, []);

  // Map route to human title
  const getPageInfo = () => {
    const path = location.pathname;
    if (path === '/dashboard') return { title: 'Dashboard', subtitle: 'Real-time telemetry and active study modules' };
    if (path.startsWith('/roadmap')) return { title: 'DSA Roadmap', subtitle: 'Curriculum from fundamentals to advanced trees & DP' };
    if (path.startsWith('/practice')) return { title: 'Practice Arena', subtitle: 'Production Monaco editor with live compiler & test runner' };
    if (path.startsWith('/assistant')) return { title: 'AI Assistant', subtitle: 'Context-aware hints, algorithmic explainers & debugger' };
    if (path.startsWith('/analysis')) return { title: 'Mistake Analysis', subtitle: 'Automated breakdown of incorrect submissions & edge cases' };
    if (path.startsWith('/analytics')) return { title: 'Performance Analytics', subtitle: 'Detailed breakdown of time, accuracy, and topic strengths' };
    if (path.startsWith('/profile')) return { title: 'Student Profile', subtitle: 'Academic credentials, target roles, and language settings' };
    if (path.startsWith('/settings')) return { title: 'Platform Settings', subtitle: 'Theme, notification triggers, and API credentials' };
    return { title: 'LearnX Command Center', subtitle: 'Intelligent coding environment' };
  };

  const { title, subtitle } = getPageInfo();

  return (
    <div className="min-h-screen bg-[#07070B] flex text-[#F4F4F6] relative selection:bg-[#8B5CF6]/30">
      <Sidebar
        mobileOpen={mobileSidebarOpen}
        onCloseMobile={() => setMobileSidebarOpen(false)}
      />
      <div className="flex-1 flex flex-col min-w-0 overflow-x-hidden">
        <Navbar
          onOpenMobile={() => setMobileSidebarOpen(true)}
          onOpenCommandPalette={() => setCommandPaletteOpen(true)}
          title={title}
          subtitle={subtitle}
        />
        <main className="flex-1 p-4 sm:p-6 lg:p-8 max-w-7xl mx-auto w-full">
          <Outlet />
        </main>
      </div>

      {/* Global Command Palette Modal */}
      <CommandPalette
        isOpen={commandPaletteOpen}
        onClose={() => setCommandPaletteOpen(false)}
      />
    </div>
  );
};

export default DashboardLayout;
