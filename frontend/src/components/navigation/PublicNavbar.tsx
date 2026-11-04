import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Sparkles, ArrowRight, Menu, X } from 'lucide-react';
import Button from '../ui/Button';

export const PublicNavbar: React.FC = () => {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const navigate = useNavigate();

  return (
    <nav className="fixed top-0 left-0 right-0 z-50 bg-white/90 backdrop-blur-xl border-b border-[#E2E8F0]">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-18 flex items-center justify-between">
        {/* Brand */}
        <Link to="/" className="flex items-center gap-2.5 group">
          <div className="w-8 h-8 rounded-xl bg-gradient-to-tr from-[#EF4444] to-[#B91C1C] flex items-center justify-center shadow-md shadow-[#EF4444]/30 group-hover:scale-105 transition-transform">
            <Sparkles className="w-4 h-4 text-white" />
          </div>
          <span className="font-bold text-xl text-[#0F172A] tracking-tight font-mono">
            Learn<span className="text-[#EF4444]">X</span>
          </span>
        </Link>

        {/* Desktop Links */}
        <div className="hidden md:flex items-center gap-8 text-sm font-medium text-[#64748B]">
          <a href="#problem" className="hover:text-[#0F172A] transition-colors">The Problem</a>
          <a href="#approach" className="hover:text-[#0F172A] transition-colors">Our Approach</a>
          <a href="#how-it-works" className="hover:text-[#0F172A] transition-colors">How It Works</a>
          <a href="#showcase" className="hover:text-[#0F172A] transition-colors">Product Showcase</a>
        </div>

        {/* Action Buttons */}
        <div className="hidden md:flex items-center gap-3">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => navigate('/login')}
          >
            Log In
          </Button>
          <Button
            variant="primary"
            size="sm"
            rightIcon={<ArrowRight className="w-4 h-4" />}
            onClick={() => navigate('/signup')}
          >
            Start Learning
          </Button>
        </div>

        {/* Mobile menu button */}
        <button
          onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
          className="md:hidden p-2 text-[#64748B] hover:text-[#0F172A]"
        >
          {mobileMenuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
        </button>
      </div>

      {/* Mobile Drawer */}
      {mobileMenuOpen && (
        <div className="md:hidden bg-white border-b border-[#E2E8F0] px-6 py-5 space-y-4">
          <div className="flex flex-col space-y-3 text-sm font-medium text-[#64748B]">
            <a
              href="#problem"
              onClick={() => setMobileMenuOpen(false)}
              className="hover:text-[#0F172A] py-1"
            >
              The Problem
            </a>
            <a
              href="#approach"
              onClick={() => setMobileMenuOpen(false)}
              className="hover:text-[#0F172A] py-1"
            >
              Our Approach
            </a>
            <a
              href="#how-it-works"
              onClick={() => setMobileMenuOpen(false)}
              className="hover:text-[#0F172A] py-1"
            >
              How It Works
            </a>
            <a
              href="#showcase"
              onClick={() => setMobileMenuOpen(false)}
              className="hover:text-[#0F172A] py-1"
            >
              Product Showcase
            </a>
          </div>
          <div className="pt-4 border-t border-[#E2E8F0] flex flex-col gap-2.5">
            <Button
              variant="secondary"
              className="w-full justify-center"
              onClick={() => {
                setMobileMenuOpen(false);
                navigate('/login');
              }}
            >
              Log In
            </Button>
            <Button
              variant="primary"
              className="w-full justify-center"
              onClick={() => {
                setMobileMenuOpen(false);
                navigate('/signup');
              }}
            >
              Start Learning
            </Button>
          </div>
        </div>
      )}
    </nav>
  );
};

export default PublicNavbar;

