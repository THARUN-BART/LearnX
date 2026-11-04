import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Terminal, ArrowLeft } from 'lucide-react';
import Button from '../components/ui/Button';

export const NotFoundPage: React.FC = () => {
  const navigate = useNavigate();

  return (
    <div className="min-h-screen bg-[#F8FAFC] flex flex-col items-center justify-center p-6 text-center text-[#0F172A]">
      <div className="w-16 h-16 rounded-2xl bg-[#EF4444]/15 border border-[#EF4444]/30 text-[#DC2626] flex items-center justify-center mb-6">
        <Terminal className="w-8 h-8" />
      </div>

      <span className="text-xs font-mono text-[#EF4444] uppercase tracking-widest block mb-2">
        Error 404 • Path Undefined
      </span>

      <h1 className="text-3xl sm:text-5xl font-bold font-mono tracking-tight mb-3">
        Page Not Found
      </h1>

      <p className="text-sm text-[#64748B] max-w-md mb-8">
        The route you are navigating to does not exist or has been relocated to another curriculum milestone.
      </p>

      <Button
        variant="primary"
        size="md"
        leftIcon={<ArrowLeft className="w-4 h-4" />}
        onClick={() => navigate('/dashboard')}
      >
        Return to Dashboard
      </Button>
    </div>
  );
};

export default NotFoundPage;

