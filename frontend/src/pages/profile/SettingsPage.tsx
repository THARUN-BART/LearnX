import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Settings,
  Bell,
  Moon,
  Shield,
  Key,
  LogOut,
  Server,
  CheckCircle2,
  AlertCircle
} from 'lucide-react';
import { useAuth } from '../../hooks/useAuth';
import { API_BASE_URL } from '../../utils/constants';
import Card from '../../components/ui/Card';
import Button from '../../components/ui/Button';
import Badge from '../../components/ui/Badge';
import { toast } from 'sonner';

export const SettingsPage: React.FC = () => {
  const { logout, user } = useAuth();
  const navigate = useNavigate();

  const [emailAlerts, setEmailAlerts] = useState(true);
  const [weeklyDigest, setWeeklyDigest] = useState(true);
  const [mistakeNotifications, setMistakeNotifications] = useState(true);

  const handleLogout = () => {
    logout();
    toast.success('Signed out successfully');
    navigate('/login');
  };

  return (
    <div className="space-y-8 max-w-4xl mx-auto">
      <div>
        <h1 className="text-2xl sm:text-3xl font-bold tracking-tight text-[#F5F5F7] font-mono flex items-center gap-2">
          <Settings className="w-6 h-6 text-[#8B5CF6]" />
          <span>Platform Settings & Preferences</span>
        </h1>
        <p className="text-xs text-[#9292A3] mt-1">
          Configure notification dispatch rules, backend API endpoints, and authentication tokens.
        </p>
      </div>


      {/* Notification Preferences */}
      <Card className="p-6 space-y-4">
        <div className="flex items-center gap-2 pb-2 border-b border-[#232336] text-xs font-bold text-[#F5F5F7] font-mono">
          <Bell className="w-4 h-4 text-[#8B5CF6]" />
          <span>Notification & Telemetry Triggers</span>
        </div>

        <div className="space-y-3 text-xs">
          <label className="flex items-center justify-between p-3 rounded-xl hover:bg-[#101019] cursor-pointer transition-colors border border-transparent hover:border-[#232336]">
            <div>
              <span className="text-sm font-semibold text-[#F5F5F7] block">Daily Streak Reminders</span>
              <span className="text-[11px] text-[#65657B]">Receive an alert if your daily problem is not submitted by 8 PM</span>
            </div>
            <input
              type="checkbox"
              checked={emailAlerts}
              onChange={(e) => setEmailAlerts(e.target.checked)}
              className="w-4 h-4 rounded bg-[#101019] border-[#232336] accent-[#8B5CF6]"
            />
          </label>

          <label className="flex items-center justify-between p-3 rounded-xl hover:bg-[#101019] cursor-pointer transition-colors border border-transparent hover:border-[#232336]">
            <div>
              <span className="text-sm font-semibold text-[#F5F5F7] block">Mistake Analysis Forensic Digests</span>
              <span className="text-[11px] text-[#65657B]">Get weekly breakdown of recurring algorithmic anti-patterns</span>
            </div>
            <input
              type="checkbox"
              checked={mistakeNotifications}
              onChange={(e) => setMistakeNotifications(e.target.checked)}
              className="w-4 h-4 rounded bg-[#101019] border-[#232336] accent-[#8B5CF6]"
            />
          </label>

          <label className="flex items-center justify-between p-3 rounded-xl hover:bg-[#101019] cursor-pointer transition-colors border border-transparent hover:border-[#232336]">
            <div>
              <span className="text-sm font-semibold text-[#F5F5F7] block">Placement & Contest Announcements</span>
              <span className="text-[11px] text-[#65657B]">Curated company coding challenges and new practice sets</span>
            </div>
            <input
              type="checkbox"
              checked={weeklyDigest}
              onChange={(e) => setWeeklyDigest(e.target.checked)}
              className="w-4 h-4 rounded bg-[#101019] border-[#232336] accent-[#8B5CF6]"
            />
          </label>
        </div>
      </Card>

      {/* Account Security & Sign Out */}
      <Card className="p-6 space-y-4 border-red-500/20">
        <div className="flex items-center gap-2 pb-2 border-b border-[#232336] text-xs font-bold text-red-400 font-mono">
          <Shield className="w-4 h-4" />
          <span>Security & Session Termination</span>
        </div>

        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 text-xs">
          <div>
            <span className="font-semibold text-[#F5F5F7] block">Sign Out of LearnX</span>
            <span className="text-[#65657B] text-[11px]">
              Active account: {user?.email || 'student@learnx.dev'}
            </span>
          </div>

          <Button
            variant="danger"
            size="sm"
            onClick={handleLogout}
            leftIcon={<LogOut className="w-4 h-4" />}
          >
            Sign Out
          </Button>
        </div>
      </Card>
    </div>
  );
};

export default SettingsPage;

