import React from 'react';
import { Outlet, Link } from 'react-router-dom';
import { Sparkles, ArrowLeft, Terminal, CheckCircle2, Cpu } from 'lucide-react';

export const AuthLayout: React.FC = () => {
  return (
    <div className="min-h-screen bg-[#08080D] flex flex-col justify-between text-[#F5F5F7] relative overflow-hidden">
      {/* Background ambient radial glow */}
      <div className="absolute top-[-20%] left-[-10%] w-[500px] h-[500px] rounded-full bg-[#8B5CF6]/10 blur-[130px] pointer-events-none" />
      <div className="absolute bottom-[-15%] right-[-10%] w-[600px] h-[600px] rounded-full bg-[#6366F1]/10 blur-[150px] pointer-events-none" />

      {/* Top Header */}
      <header className="px-6 py-5 flex items-center justify-between z-10 max-w-7xl mx-auto w-full">
        <Link to="/" className="flex items-center gap-2.5 group">
          <div className="w-8 h-8 rounded-xl bg-gradient-to-tr from-[#8B5CF6] to-[#6366F1] flex items-center justify-center shadow-md shadow-[#8B5CF6]/30">
            <Sparkles className="w-4 h-4 text-white" />
          </div>
          <span className="font-bold text-xl text-[#F5F5F7] tracking-tight font-mono">
            Learn<span className="text-[#8B5CF6]">X</span>
          </span>
        </Link>
        <Link
          to="/"
          className="inline-flex items-center gap-1.5 text-xs text-[#9292A3] hover:text-[#F5F5F7] transition-colors"
        >
          <ArrowLeft className="w-3.5 h-3.5" />
          <span>Back to Home</span>
        </Link>
      </header>

      {/* Main Content Split or Centered */}
      <main className="flex-1 flex items-center justify-center p-4 sm:p-6 z-10">
        <div className="w-full max-w-5xl grid grid-cols-1 lg:grid-cols-12 gap-8 items-center">
          {/* Editorial / Technical Hero Info (Desktop Left) */}
          <div className="hidden lg:flex lg:col-span-6 flex-col justify-center space-y-6 pr-6">
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-[#171722] border border-[#232336] text-xs font-mono text-[#A78BFA] w-fit">
              <Cpu className="w-3.5 h-3.5" />
              <span>Real-Time Code Intelligence</span>
            </div>

            <h2 className="text-3xl xl:text-4xl font-bold tracking-tight text-[#F5F5F7] leading-tight font-mono">
              The fastest way to achieve <span className="text-transparent bg-clip-text bg-gradient-to-r from-[#A78BFA] to-[#818CF8]">algorithmic mastery</span>.
            </h2>

            <p className="text-sm text-[#9292A3] leading-relaxed">
              LearnX diagnoses algorithmic gaps in real time, generates tailored hints without spoiling solutions, and builds your technical interview confidence.
            </p>

            <div className="space-y-3 pt-2">
              <div className="flex items-center gap-3 text-xs text-[#9292A3]">
                <CheckCircle2 className="w-4 h-4 text-emerald-400 shrink-0" />
                <span>Automated mistake classification & anti-pattern analysis</span>
              </div>
              <div className="flex items-center gap-3 text-xs text-[#9292A3]">
                <CheckCircle2 className="w-4 h-4 text-emerald-400 shrink-0" />
                <span>Multi-language code execution: C++, Java, Python, Go, C</span>
              </div>
              <div className="flex items-center gap-3 text-xs text-[#9292A3]">
                <CheckCircle2 className="w-4 h-4 text-emerald-400 shrink-0" />
                <span>Dynamic roadmap calibrated to college placement & FAANG goals</span>
              </div>
            </div>

            {/* Terminal snippet preview */}
            <div className="p-4 rounded-xl bg-[#101019] border border-[#232336] font-mono text-xs text-[#9292A3] space-y-1 shadow-xl">
              <div className="flex items-center gap-1.5 pb-2 border-b border-[#232336] text-[11px] text-[#65657B]">
                <Terminal className="w-3.5 h-3.5 text-[#8B5CF6]" />
                <span>learnx-daemon ~ runtime-eval</span>
              </div>
              <div className="text-emerald-400 pt-1">✓ Optimized Time Complexity: O(N)</div>
              <div className="text-[#A78BFA]">→ Memory Footprint: 14.2 MB (Top 96.4%)</div>
            </div>
          </div>

          {/* Form Container (Right) */}
          <div className="w-full lg:col-span-6 flex justify-center">
            <div className="w-full max-w-md bg-[#171722]/90 backdrop-blur-xl border border-[#232336] rounded-2xl p-6 sm:p-8 shadow-2xl shadow-black/60 relative">
              <Outlet />
            </div>
          </div>
        </div>
      </main>

      {/* Footer */}
      <footer className="py-4 text-center text-xs text-[#65657B] z-10">
        © {new Date().getFullYear()} LearnX Inc. Built for ambitious software engineers.
      </footer>
    </div>
  );
};

export default AuthLayout;

