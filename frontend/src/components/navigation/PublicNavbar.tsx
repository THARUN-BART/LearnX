import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Sparkles, ArrowRight, Menu, X } from 'lucide-react';
import Button from '../ui/Button';

export const PublicNavbar: React.FC = () => {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const navigate = useNavigate();

  return (
    <nav className="fixed top-0 left-0 right-0 z-50 bg-[#08080D]/80 backdrop-blur-xl border-b border-[#232336]/60">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-18 flex items-center justify-between">
        {/* Brand */}
        <Link to="/" className="flex items-center gap-2.5 group">
          <div className="w-8 h-8 rounded-xl bg-gradient-to-tr from-[#8B5CF6] to-[#6366F1] flex items-center justify-center shadow-md shadow-[#8B5CF6]/30 group-hover:scale-105 transition-transform">
            <Sparkles className="w-4 h-4 text-white" />
          </div>
          <span className="font-bold text-xl text-[#F5F5F7] tracking-tight font-mono">
            Learn<span className="text-[#8B5CF6]">X</span>
          </span>
        </Link>

        {/* Desktop Links */}
        <div className="hidden md:flex items-center gap-8 text-sm font-medium text-[#9292A3]">
          <a href="#problem" className="hover:text-[#F5F5F7] transition-colors">The Problem</a>
          <a href="#approach" className="hover:text-[#F5F5F7] transition-colors">Our Approach</a>
          <a href="#how-it-works" className="hover:text-[#F5F5F7] transition-colors">How It Works</a>
          <a href="#showcase" className="hover:text-[#F5F5F7] transition-colors">Product Showcase</a>
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
          className="md:hidden p-2 text-[#9292A3] hover:text-[#F5F5F7]"
        >
          {mobileMenuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
        </button>
      </div>

      {/* Mobile Drawer */}
      {mobileMenuOpen && (
        <div className="md:hidden bg-[#101019] border-b border-[#232336] px-6 py-5 space-y-4">
          <div className="flex flex-col space-y-3 text-sm font-medium text-[#9292A3]">
            <a
              href="#problem"
              onClick={() => setMobileMenuOpen(false)}
              className="hover:text-[#F5F5F7] py-1"
            >
              The Problem
            </a>
            <a
              href="#approach"
              onClick={() => setMobileMenuOpen(false)}
              className="hover:text-[#F5F5F7] py-1"
            >
              Our Approach
            </a>
            <a
              href="#how-it-works"
              onClick={() => setMobileMenuOpen(false)}
              className="hover:text-[#F5F5F7] py-1"
            >
              How It Works
            </a>
            <a
              href="#showcase"
              onClick={() => setMobileMenuOpen(false)}
              className="hover:text-[#F5F5F7] py-1"
            >
              Product Showcase
            </a>
          </div>
          <div className="pt-4 border-t border-[#232336] flex flex-col gap-2.5">
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

