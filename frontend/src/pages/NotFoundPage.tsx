import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Terminal, ArrowLeft } from 'lucide-react';
import Button from '../components/ui/Button';

export const NotFoundPage: React.FC = () => {
  const navigate = useNavigate();

  return (
    <div className="min-h-screen bg-[#08080D] flex flex-col items-center justify-center p-6 text-center text-[#F5F5F7]">
      <div className="w-16 h-16 rounded-2xl bg-[#8B5CF6]/15 border border-[#8B5CF6]/30 text-[#A78BFA] flex items-center justify-center mb-6">
        <Terminal className="w-8 h-8" />
      </div>

      <span className="text-xs font-mono text-[#8B5CF6] uppercase tracking-widest block mb-2">
        Error 404 • Path Undefined
      </span>

      <h1 className="text-3xl sm:text-5xl font-bold font-mono tracking-tight mb-3">
        Page Not Found
      </h1>

      <p className="text-sm text-[#9292A3] max-w-md mb-8">
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

