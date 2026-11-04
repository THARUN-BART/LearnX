import React from 'react';
import { NavLink, useNavigate } from 'react-router-dom';
import {
  LayoutDashboard,
  Map,
  Code2,
  Bot,
  AlertTriangle,
  BarChart3,
  User,
  Settings,
  LogOut,
  Flame,
  Sparkles,
  X,
  Brain,
  Compass,
} from 'lucide-react';
import { useAuth } from '../../hooks/useAuth';
import { useProfile } from '../../hooks/useProfile';
import { clsx } from 'clsx';

interface SidebarProps {
  mobileOpen?: boolean;
  onCloseMobile?: () => void;
}

export const Sidebar: React.FC<SidebarProps> = ({ mobileOpen, onCloseMobile }) => {
  const { user, logout } = useAuth();
  const { profile } = useProfile();
  const navigate = useNavigate();

  const navItems = [
    { label: 'Overview',        to: '/dashboard', icon: LayoutDashboard, badge: 'Active' },
    { label: 'Smart Picks',     to: '/recommendations', icon: Sparkles, badge: 'AI' },
    { label: 'Adaptive Path',   to: '/dashboard/learning-path', icon: Compass, badge: 'AI' },
    { label: 'My Roadmap',      to: '/roadmap',   icon: Map,             badge: '19' },
    { label: 'Practice Arena',  to: '/practice',  icon: Code2,           badge: 'IDE' },
    { label: 'AI Assistant',    to: '/assistant',  icon: Bot,             badge: 'PRO' },
    { label: 'Learning Memory', to: '/dashboard/learning-memory', icon: Brain, badge: 'NEW' },
    { label: 'Mistake Analysis',to: '/analysis',   icon: AlertTriangle },
    { label: 'Performance',     to: '/analytics',  icon: BarChart3 },
  ];

  const bottomNavItems = [
    { label: 'My Profile', to: '/profile', icon: User },
    { label: 'Settings', to: '/settings', icon: Settings },
  ];

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const content = (
    <div className="flex flex-col h-full bg-white border-r border-[#E2E8F0] w-64 select-none">
      {/* Brand Header */}
      <div className="h-16 flex items-center justify-between px-5 border-b border-[#E2E8F0]">
        <NavLink to="/dashboard" className="flex items-center gap-2.5">
          <div className="w-8 h-8 rounded-xl bg-gradient-to-tr from-[#EF4444] to-[#B91C1C] flex items-center justify-center shadow-md shadow-red-500/30">
            <Sparkles className="w-4 h-4 text-white" />
          </div>
          <div className="flex flex-col">
            <span className="font-bold text-base text-[#0F172A] tracking-tight font-mono">
              Learn<span className="text-[#EF4444]">X</span>
            </span>
            <span className="text-[9px] font-mono tracking-widest text-[#64748B] uppercase">
              Adaptive Engine
            </span>
          </div>
        </NavLink>
        {onCloseMobile && (
          <button
            onClick={onCloseMobile}
            className="md:hidden text-[#475569] hover:text-[#0F172A] p-1.5"
          >
            <X className="w-5 h-5" />
          </button>
        )}
      </div>

      {/* Streak & Consistency Bar */}
      <div className="p-3 mx-3 my-3 bg-[#F8FAFC] border border-[#E2E8F0] rounded-xl space-y-2">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="w-6 h-6 rounded-lg bg-amber-500/10 border border-amber-500/20 flex items-center justify-center text-amber-500">
              <Flame className="w-3.5 h-3.5 animate-pulse" />
            </div>
            <span className="text-xs font-mono font-semibold text-[#0F172A]">
              {profile?.streakDays ?? 0} Day{(profile?.streakDays ?? 0) !== 1 ? 's' : ''} Active
            </span>
          </div>
          <span className="text-[10px] font-mono text-emerald-600 font-medium">
            {(profile?.streakDays ?? 0) >= 7 ? 'Top 5%' : 'Building'}
          </span>
        </div>
        <div className="w-full bg-[#E2E8F0] h-1 rounded-full overflow-hidden">
          <div
            className="bg-gradient-to-r from-amber-500 to-[#EF4444] h-full transition-all duration-500"
            style={{ width: `${Math.min(100, Math.round(((profile?.streakDays ?? 0) / 30) * 100))}%` }}
          />
        </div>
      </div>

      {/* Navigation Links */}
      <div className="flex-1 px-3 py-1 space-y-1 overflow-y-auto">
        <div className="text-[10px] font-semibold uppercase tracking-widest text-[#64748B] px-3 py-1.5 font-mono">
          Cockpit
        </div>
        {navItems.map((item) => {
          const Icon = item.icon;
          return (
            <NavLink
              key={item.label}
              to={item.to}
              onClick={onCloseMobile}
              className={({ isActive }) =>
                clsx(
                  'flex items-center justify-between px-3 py-2 rounded-xl text-xs font-mono font-medium transition-all duration-150 group relative',
                  isActive
                    ? 'bg-[#FEF2F2] text-[#DC2626] border border-[#FCA5A5]/60 shadow-sm'
                    : 'text-[#475569] hover:text-[#0F172A] hover:bg-[#F1F5F9] border border-transparent'
                )
              }
            >
              <div className="flex items-center gap-2.5">
                <Icon className="w-4 h-4 transition-transform group-hover:scale-110" />
                <span>{item.label}</span>
              </div>
              {item.badge && (
                <span className="text-[9px] font-mono px-1.5 py-0.5 rounded bg-slate-100 text-[#475569] border border-[#E2E8F0]">
                  {item.badge}
                </span>
              )}
            </NavLink>
          );
        })}

        <div className="text-[10px] font-semibold uppercase tracking-widest text-[#64748B] px-3 pt-5 pb-1.5 font-mono">
          Preferences
        </div>
        {bottomNavItems.map((item) => {
          const Icon = item.icon;
          return (
            <NavLink
              key={item.label}
              to={item.to}
              onClick={onCloseMobile}
              className={({ isActive }) =>
                clsx(
                  'flex items-center gap-2.5 px-3 py-2 rounded-xl text-xs font-mono font-medium transition-all duration-150 group',
                  isActive
                    ? 'bg-[#FEF2F2] text-[#DC2626] border border-[#FCA5A5]/60'
                    : 'text-[#475569] hover:text-[#0F172A] hover:bg-[#F1F5F9] border border-transparent'
                )
              }
            >
              <Icon className="w-4 h-4 transition-transform group-hover:scale-110" />
              <span>{item.label}</span>
            </NavLink>
          );
        })}
      </div>

      {/* User profile footer & Logout */}
      <div className="p-3 border-t border-[#E2E8F0] bg-[#F8FAFC]">
        <div className="flex items-center justify-between p-2 rounded-xl hover:bg-[#F1F5F9] transition-colors">
          <div className="flex items-center gap-2.5 min-w-0">
            <div className="w-7 h-7 rounded-lg bg-gradient-to-tr from-[#EF4444] to-[#B91C1C] flex items-center justify-center text-white font-bold text-xs shrink-0 font-mono shadow-sm">
              {(user?.fullName || user?.full_name || 'U').charAt(0)}
            </div>
            <div className="min-w-0">
              <div className="text-xs font-semibold text-[#0F172A] truncate font-mono">
                {user?.fullName || user?.full_name || 'Student'}
              </div>
              <div className="text-[10px] text-[#64748B] truncate font-mono">
                {user?.email || 'student@gmail.com'}
              </div>
            </div>
          </div>
          <button
            onClick={handleLogout}
            title="Log Out"
            className="p-1.5 text-[#475569] hover:text-red-500 hover:bg-red-50 rounded-lg transition-colors cursor-pointer"
          >
            <LogOut className="w-3.5 h-3.5" />
          </button>
        </div>
      </div>
    </div>
  );

  return (
    <>
      {/* Desktop Persistent Sidebar */}
      <aside className="hidden md:block h-screen sticky top-0 shrink-0 z-30">
        {content}
      </aside>

      {/* Mobile Drawer */}
      {mobileOpen && (
        <div className="fixed inset-0 z-50 md:hidden flex">
          <div
            className="fixed inset-0 bg-black/70 backdrop-blur-sm"
            onClick={onCloseMobile}
          />
          <div className="relative z-10 animate-in slide-in-from-left duration-200">
            {content}
          </div>
        </div>
      )}
    </>
  );
};

export default Sidebar;
