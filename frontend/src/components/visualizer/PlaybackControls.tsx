import React from 'react';
import {
  Play,
  Pause,
  SkipBack,
  SkipForward,
  RotateCcw,
  Gauge,
} from 'lucide-react';

interface PlaybackControlsProps {
  isPlaying: boolean;
  onPlayPause: () => void;
  onNext: () => void;
  onPrev: () => void;
  onRestart: () => void;
  currentStep: number;
  totalSteps: number;
  onSeek: (step: number) => void;
  speed: number;
  onSpeedChange: (speed: number) => void;
}

export const PlaybackControls: React.FC<PlaybackControlsProps> = ({
  isPlaying,
  onPlayPause,
  onNext,
  onPrev,
  onRestart,
  currentStep,
  totalSteps,
  onSeek,
  speed,
  onSpeedChange,
}) => {
  const speeds = [0.5, 1, 1.5, 2, 3];

  return (
    <div className="bg-white border border-[#E2E8F0] rounded-xl p-3.5 flex flex-col gap-3 font-mono">
      {/* Progress Bar / Scrub Slider */}
      <div className="flex items-center gap-3">
        <span className="text-[11px] text-[#64748B] w-12 text-right">
          {totalSteps > 0 ? currentStep + 1 : 0} / {totalSteps}
        </span>
        <input
          type="range"
          min={0}
          max={Math.max(0, totalSteps - 1)}
          value={currentStep}
          onChange={(e) => onSeek(Number(e.target.value))}
          className="flex-1 accent-[#8B5CF6] h-1.5 bg-[#F1F5F9] rounded-lg cursor-pointer"
        />
        <span className="text-[11px] text-[#DC2626] w-12 text-left">
          {totalSteps > 0 ? `${Math.round(((currentStep + 1) / totalSteps) * 100)}%` : '0%'}
        </span>
      </div>

      {/* Main Buttons */}
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div className="flex items-center gap-1.5">
          <button
            onClick={onRestart}
            title="Restart from beginning"
            className="p-2 text-[#64748B] hover:text-[#0F172A] bg-[#F1F5F9] hover:bg-[#181824] border border-[#E2E8F0] rounded-lg transition-colors"
          >
            <RotateCcw className="w-4 h-4" />
          </button>
          <button
            onClick={onPrev}
            disabled={currentStep <= 0}
            title="Previous Step"
            className="p-2 text-[#64748B] hover:text-[#0F172A] disabled:opacity-30 disabled:pointer-events-none bg-[#F1F5F9] hover:bg-[#181824] border border-[#E2E8F0] rounded-lg transition-colors"
          >
            <SkipBack className="w-4 h-4" />
          </button>
          <button
            onClick={onPlayPause}
            className={`px-4 py-2 rounded-lg flex items-center gap-2 text-xs font-bold transition-all shadow-md ${
              isPlaying
                ? 'bg-amber-500 text-black hover:bg-amber-400'
                : 'bg-[#EF4444] text-white hover:bg-[#7C3AED] shadow-[#8B5CF6]/30'
            }`}
          >
            {isPlaying ? (
              <>
                <Pause className="w-4 h-4 fill-current" /> Pause
              </>
            ) : (
              <>
                <Play className="w-4 h-4 fill-current" /> Play
              </>
            )}
          </button>
          <button
            onClick={onNext}
            disabled={currentStep >= totalSteps - 1}
            title="Next Step"
            className="p-2 text-[#64748B] hover:text-[#0F172A] disabled:opacity-30 disabled:pointer-events-none bg-[#F1F5F9] hover:bg-[#181824] border border-[#E2E8F0] rounded-lg transition-colors"
          >
            <SkipForward className="w-4 h-4" />
          </button>
        </div>

        {/* Playback Speed Selector */}
        <div className="flex items-center gap-2 bg-[#F1F5F9] border border-[#E2E8F0] px-2.5 py-1.5 rounded-lg">
          <Gauge className="w-3.5 h-3.5 text-[#64748B]" />
          <span className="text-[10px] text-[#94A3B8] uppercase font-bold">Speed:</span>
          <div className="flex items-center gap-1">
            {speeds.map((s) => (
              <button
                key={s}
                onClick={() => onSpeedChange(s)}
                className={`text-[10px] px-2 py-0.5 rounded transition-colors ${
                  speed === s
                    ? 'bg-[#EF4444] text-white font-bold'
                    : 'text-[#64748B] hover:text-[#0F172A]'
                }`}
              >
                {s}x
              </button>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};
export default PlaybackControls;

