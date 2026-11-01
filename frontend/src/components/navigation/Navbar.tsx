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
    <header className="h-16 border-b border-[#1C1C28] bg-[#07070B]/90 backdrop-blur-xl px-4 sm:px-6 flex items-center justify-between sticky top-0 z-20">
      <div className="flex items-center gap-3">
        <button
          onClick={onOpenMobile}
          className="md:hidden p-2 rounded-xl text-[#8E8EA0] hover:text-[#F4F4F6] hover:bg-[#14141E] transition-colors"
        >
          <Menu className="w-5 h-5" />
        </button>

        <div>
          {title && (
            <h1 className="text-base sm:text-lg font-bold text-[#F4F4F6] tracking-tight flex items-center gap-2 font-mono">
              <span className="text-[#8B5CF6]">/</span> {title}
            </h1>
          )}
          {subtitle && (
            <p className="text-[11px] font-mono text-[#5A5A6E] hidden sm:block">{subtitle}</p>
          )}
        </div>
      </div>

      <div className="flex items-center gap-3">
        {/* Real-time Kernel Latency indicator */}
        <div className="hidden sm:flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-[#12121B] border border-[#1C1C28] text-[11px] font-mono text-[#8E8EA0]">
          <span className="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse" />
          <span>Daemon: 18ms</span>
        </div>

        {/* Command Palette Trigger */}
        <button
          onClick={onOpenCommandPalette}
          className="flex items-center gap-2.5 px-3 py-1.5 rounded-xl bg-[#12121B] hover:bg-[#181824] border border-[#1C1C28] hover:border-[#2D2D3E] text-xs text-[#8E8EA0] transition-all cursor-pointer group"
        >
          <Search className="w-3.5 h-3.5 text-[#8B5CF6] group-hover:scale-110 transition-transform" />
          <span className="hidden md:inline font-mono">Quick search...</span>
          <kbd className="keycap hidden md:inline">⌘K</kbd>
        </button>

        <div className="h-4 w-px bg-[#1C1C28] hidden sm:block" />

        <button className="p-2 rounded-xl text-[#8E8EA0] hover:text-[#F4F4F6] hover:bg-[#14141E] relative transition-colors">
          <Bell className="w-4 h-4" />
          <span className="absolute top-2 right-2 w-2 h-2 rounded-full bg-[#8B5CF6] ring-2 ring-[#07070B]" />
        </button>
      </div>
    </header>
  );
};

export default Navbar;
