import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Search,
  Map,
  Code2,
  Bot,
  AlertTriangle,
  BarChart3,
  User,
  Settings,
  ArrowRight,
  Terminal,
  Sparkles,
  Command,
  X
} from 'lucide-react';
import roadmapService from '../../services/roadmap.service';
import learningService from '../../services/learning.service';
import { Topic, Problem } from '../../types/learning.types';

interface CommandPaletteProps {
  isOpen: boolean;
  onClose: () => void;
}

export const CommandPalette: React.FC<CommandPaletteProps> = ({ isOpen, onClose }) => {
  const navigate = useNavigate();
  const [query, setQuery] = useState('');
  const [selectedIndex, setSelectedIndex] = useState(0);
  const [topics, setTopics] = useState<Topic[]>([]);
  const [problems, setProblems] = useState<Problem[]>([]);

  useEffect(() => {
    if (!isOpen) return;
    roadmapService.getRoadmap().then((tList) => {
      if (Array.isArray(tList)) setTopics(tList);
    });
  }, [isOpen]);

  useEffect(() => {
    if (!isOpen) return;
    const q = query.trim();
    if (q.length < 2) {
      setProblems([]);
      return;
    }
    const timer = setTimeout(async () => {
      const pList = await learningService.searchProblems(q, 12);
      if (Array.isArray(pList)) setProblems(pList);
    }, 250);
    return () => clearTimeout(timer);
  }, [isOpen, query]);

  // Quick navigation destinations
  const navigationItems = [
    { label: 'Overview Dashboard', category: 'Navigation', icon: BarChart3, path: '/dashboard' },
    { label: 'DSA Roadmap', category: 'Navigation', icon: Map, path: '/roadmap' },
    { label: 'Practice Arena & Compiler', category: 'Navigation', icon: Code2, path: '/practice' },
    { label: 'AI Coding Mentor', category: 'Navigation', icon: Bot, path: '/assistant' },
    { label: 'Mistake Diagnostics', category: 'Navigation', icon: AlertTriangle, path: '/analysis' },
    { label: 'Performance Analytics', category: 'Navigation', icon: BarChart3, path: '/analytics' },
    { label: 'Student Profile', category: 'Navigation', icon: User, path: '/profile' },
    { label: 'System Settings', category: 'Navigation', icon: Settings, path: '/settings' },
  ];

  // Map real topics and problems into searchable items
  const topicItems = topics.map((t) => ({
    label: `Roadmap: ${t.title}`,
    category: 'Curriculum Pillar',
    icon: Map,
    path: `/roadmap/${t.slug || t.id}`
  }));

  const problemItems = problems.map((p) => ({
    label: `Solve: ${p.title} (${p.difficulty})`,
    category: 'Coding Challenge',
    icon: Code2,
    path: `/practice/${p.id}`
  }));

  const allItems = [...navigationItems, ...topicItems, ...problemItems];

  const filtered = allItems.filter((item) =>
    item.label.toLowerCase().includes(query.toLowerCase()) ||
    item.category.toLowerCase().includes(query.toLowerCase())
  );

  useEffect(() => {
    setSelectedIndex(0);
  }, [query]);

  // Keyboard navigation
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (!isOpen) {
        if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
          e.preventDefault();
          onClose(); // toggle
        }
        return;
      }

      if (e.key === 'Escape') {
        e.preventDefault();
        onClose();
      } else if (e.key === 'ArrowDown') {
        e.preventDefault();
        setSelectedIndex((prev) => (prev + 1) % Math.max(1, filtered.length));
      } else if (e.key === 'ArrowUp') {
        e.preventDefault();
        setSelectedIndex((prev) => (prev - 1 + filtered.length) % Math.max(1, filtered.length));
      } else if (e.key === 'Enter') {
        e.preventDefault();
        if (filtered[selectedIndex]) {
          navigate(filtered[selectedIndex].path);
          onClose();
        }
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [isOpen, filtered, selectedIndex, navigate, onClose]);

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-start justify-center pt-20 px-4">
      {/* Backdrop */}
      <div
        className="fixed inset-0 bg-black/70 backdrop-blur-md animate-in fade-in duration-150"
        onClick={onClose}
      />

      {/* Palette Container */}
      <div className="relative w-full max-w-xl bg-white border border-[#2B2B3C] rounded-2xl shadow-2xl shadow-black overflow-hidden z-10 animate-in zoom-in-95 duration-150">
        {/* Search Input Bar */}
        <div className="flex items-center px-4 py-3.5 border-b border-[#E2E8F0] gap-3">
          <Search className="w-4 h-4 text-[#EF4444] shrink-0" />
          <input
            autoFocus
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Type a command, problem title, or DSA topic..."
            className="flex-1 bg-transparent text-sm text-[#0F172A] placeholder-[#5A5A6E] outline-none font-mono"
          />
          <button
            onClick={onClose}
            className="p-1 rounded-md text-[#94A3B8] hover:text-[#0F172A] transition-colors"
          >
            <X className="w-4 h-4" />
          </button>
        </div>

        {/* Results List */}
        <div className="max-h-80 overflow-y-auto p-2 space-y-1">
          {filtered.length === 0 ? (
            <div className="py-8 text-center text-xs font-mono text-[#94A3B8]">
              No matching commands or modules found.
            </div>
          ) : (
            filtered.slice(0, 10).map((item, index) => {
              const Icon = item.icon;
              const isSelected = selectedIndex === index;
              return (
                <div
                  key={item.label}
                  onClick={() => {
                    navigate(item.path);
                    onClose();
                  }}
                  className={`flex items-center justify-between px-3.5 py-2.5 rounded-xl cursor-pointer transition-all ${
                    isSelected
                      ? 'bg-[#1C1C2B] text-white border border-[#EF4444]/40'
                      : 'text-[#64748B] hover:bg-[#14141E] border border-transparent'
                  }`}
                >
                  <div className="flex items-center gap-3">
                    <div
                      className={`w-6 h-6 rounded-lg flex items-center justify-center ${
                        isSelected
                          ? 'bg-[#EF4444] text-white'
                          : 'bg-[#181824] text-[#64748B]'
                      }`}
                    >
                      <Icon className="w-3.5 h-3.5" />
                    </div>
                    <span className="text-xs font-medium">{item.label}</span>
                  </div>

                  <div className="flex items-center gap-2">
                    <span className="text-[10px] font-mono text-[#94A3B8] uppercase tracking-wider">
                      {item.category}
                    </span>
                    <ArrowRight className="w-3 h-3 text-[#94A3B8]" />
                  </div>
                </div>
              );
            })
          )}
        </div>

        {/* Footer shortcuts hint */}
        <div className="px-4 py-2.5 bg-[#F8FAFC] border-t border-[#E2E8F0] flex items-center justify-between text-[11px] font-mono text-[#94A3B8]">
          <div className="flex items-center gap-3">
            <span>
              <kbd className="keycap mr-1">↑</kbd>
              <kbd className="keycap mr-1">↓</kbd> Navigate
            </span>
            <span>
              <kbd className="keycap mr-1">↵</kbd> Select
            </span>
            <span>
              <kbd className="keycap mr-1">ESC</kbd> Close
            </span>
          </div>
          <div className="flex items-center gap-1.5 text-[#DC2626]">
            <Sparkles className="w-3 h-3" />
            <span>LearnX Global Action Bus</span>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CommandPalette;

