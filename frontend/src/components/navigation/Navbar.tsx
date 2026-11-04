import React from 'react';
import { Menu, Bell, Sparkles, Terminal, Search, Activity } from 'lucide-react';
import { Badge } from '../ui/Badge';

interface NavbarProps {
  onOpenMobile?: () => void;
  onOpenCommandPalette?: () => void;
  title?: string;
  subtitle?: string;
}

export const Navbar: React.FC<NavbarProps> = ({
  onOpenMobile,
  onOpenCommandPalette,
  title,
  subtitle
}) => {
  return (
    <header className="h-16 border-b border-[#E2E8F0] bg-white/90 backdrop-blur-xl px-4 sm:px-6 flex items-center justify-between sticky top-0 z-20">
      <div className="flex items-center gap-3">
        <button
          onClick={onOpenMobile}
          className="md:hidden p-2 rounded-xl text-[#475569] hover:text-[#0F172A] hover:bg-[#F1F5F9] transition-colors"
        >
          <Menu className="w-5 h-5" />
        </button>

        <div>
          {title && (
            <h1 className="text-base sm:text-lg font-bold text-[#0F172A] tracking-tight flex items-center gap-2 font-mono">
              <span className="text-[#EF4444]">/</span> {title}
            </h1>
          )}
          {subtitle && (
            <p className="text-[11px] font-mono text-[#64748B] hidden sm:block">{subtitle}</p>
          )}
        </div>
      </div>

      <div className="flex items-center gap-3">
        {/* Real-time Kernel Latency indicator */}
        <div className="hidden sm:flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-[#F8FAFC] border border-[#E2E8F0] text-[11px] font-mono text-[#475569]">
          <span className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse" />
          <span>Daemon: 18ms</span>
        </div>

        {/* Command Palette Trigger */}
        <button
          onClick={onOpenCommandPalette}
          className="flex items-center gap-2.5 px-3 py-1.5 rounded-xl bg-white hover:bg-[#F8FAFC] border border-[#E2E8F0] hover:border-[#CBD5E1] text-xs text-[#475569] transition-all cursor-pointer group shadow-sm"
        >
          <Search className="w-3.5 h-3.5 text-[#EF4444] group-hover:scale-110 transition-transform" />
          <span className="hidden md:inline font-mono">Quick search...</span>
          <kbd className="keycap hidden md:inline">⌘K</kbd>
        </button>

        <div className="h-4 w-px bg-[#E2E8F0] hidden sm:block" />

        <button className="p-2 rounded-xl text-[#475569] hover:text-[#0F172A] hover:bg-[#F1F5F9] relative transition-colors">
          <Bell className="w-4 h-4" />
          <span className="absolute top-2 right-2 w-2 h-2 rounded-full bg-[#EF4444] ring-2 ring-white" />
        </button>
      </div>
    </header>
  );
};

export default Navbar;
