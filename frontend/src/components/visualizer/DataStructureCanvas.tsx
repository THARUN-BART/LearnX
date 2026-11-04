import React from 'react';
import { ExecutionStep } from '../../types/visualizer.types';
import { ArrowRight, Sparkles } from 'lucide-react';

interface DataStructureCanvasProps {
  currentStepData: ExecutionStep | null;
  category: string;
}

export const DataStructureCanvas: React.FC<DataStructureCanvasProps> = ({
  currentStepData,
  category,
}) => {
  if (!currentStepData) {
    return (
      <div className="h-80 flex flex-col items-center justify-center text-xs font-mono text-[#94A3B8] gap-2">
        <Sparkles className="w-5 h-5 text-[#EF4444]/50 animate-pulse" />
        <span>Select an algorithm and click Run to start visualization</span>
      </div>
    );
  }

  const { state, operation, description, is_complete } = currentStepData;
  const arr = state.array || [];
  const matrix = state.matrix;
  const pointers = state.pointers || {};
  const comparing = new Set(state.comparing || []);
  const swapping = new Set(state.swapping || []);
  const highlighted = new Set(state.highlighted || []);

  // Compute pointer labels per index
  const indexPointers: Record<number, string[]> = {};
  Object.entries(pointers).forEach(([name, idx]) => {
    if (!indexPointers[idx]) indexPointers[idx] = [];
    indexPointers[idx].push(name);
  });

  const maxVal = arr.length > 0 ? Math.max(...arr.map((v) => Math.abs(v)), 1) : 1;

  return (
    <div className="flex flex-col h-full min-h-[340px] bg-[#0A0A10] border border-[#E2E8F0] rounded-2xl p-5 select-none font-mono relative overflow-hidden">
      {/* Canvas Top Bar */}
      <div className="flex items-center justify-between gap-3 border-b border-[#E2E8F0] pb-3 mb-4">
        <div className="flex items-center gap-2">
          <span
            className={`text-[10px] font-bold uppercase tracking-wider px-2.5 py-1 rounded-md border ${
              is_complete
                ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/30'
                : operation === 'swap'
                ? 'bg-rose-500/10 text-rose-400 border-rose-500/30'
                : operation === 'compare'
                ? 'bg-amber-500/10 text-amber-400 border-amber-500/30'
                : 'bg-[#EF4444]/10 text-[#DC2626] border-[#EF4444]/30'
            }`}
          >
            {operation.replace(/_/g, ' ')}
          </span>
          <span className="text-xs text-[#E2E8F0] font-sans font-medium line-clamp-1">
            {description}
          </span>
        </div>
        {state.status_text && (
          <span className="text-[11px] text-[#DC2626] bg-[#F1F5F9] px-2.5 py-1 rounded-lg border border-[#E2E8F0] shrink-0">
            {state.status_text}
          </span>
        )}
      </div>

      {/* Main Visualizer Area */}
      <div className="flex-1 flex items-center justify-center overflow-x-auto p-4">
        {/* Render 2D Matrix if present (DP Table / Grid) */}
        {matrix && matrix.length > 0 ? (
          <div className="flex flex-col items-center gap-1.5 overflow-x-auto max-w-full">
            <span className="text-[10px] uppercase tracking-wider text-[#64748B] mb-2 font-bold">
              Dynamic Programming Matrix State
            </span>
            <div className="grid gap-1.5 p-3 bg-white rounded-xl border border-[#E2E8F0]">
              {matrix.map((row, rIdx) => (
                <div key={rIdx} className="flex gap-1.5 items-center">
                  <span className="w-5 text-[10px] text-[#94A3B8] text-right mr-1 font-bold">
                    {rIdx}
                  </span>
                  {row.map((cell, cIdx) => {
                    const isActive =
                      state.active_cell &&
                      state.active_cell[0] === rIdx &&
                      state.active_cell[1] === cIdx;
                    return (
                      <div
                        key={cIdx}
                        className={`w-11 h-10 flex items-center justify-center rounded-lg text-xs font-bold border transition-all ${
                          isActive
                            ? 'bg-[#EF4444] text-white border-white scale-110 shadow-lg shadow-[#EF4444]/50 z-10'
                            : cell !== 0
                            ? 'bg-[#E2E8F0] text-[#E2E8F0] border-[#EF4444]/40'
                            : 'bg-[#F1F5F9] text-[#94A3B8] border-[#E2E8F0]'
                        }`}
                      >
                        {cell}
                      </div>
                    );
                  })}
                </div>
              ))}
            </div>
          </div>
        ) : category === 'Linked Lists' ? (
          /* Render Linked List Nodes & Arrows */
          <div className="flex items-center gap-2 overflow-x-auto py-8 px-4">
            {arr.map((val, idx) => {
              const isHighlight = highlighted.has(idx);
              const ptrs = indexPointers[idx] || [];

              return (
                <React.Fragment key={idx}>
                  <div className="flex flex-col items-center gap-2 relative">
                    {/* Pointer Badges on Top */}
                    <div className="h-6 flex items-center gap-1">
                      {ptrs.map((p) => (
                        <span
                          key={p}
                          className="text-[9px] uppercase font-bold bg-[#EF4444] text-white px-1.5 py-0.5 rounded shadow-sm"
                        >
                          {p}
                        </span>
                      ))}
                    </div>

                    {/* Linked List Node Box */}
                    <div
                      className={`w-16 h-16 rounded-xl flex flex-col items-center justify-center border-2 shadow-md transition-all ${
                        isHighlight
                          ? 'bg-[#EF4444]/20 border-[#EF4444] text-white scale-105'
                          : 'bg-[#F1F5F9] border-[#E2E8F0] text-[#E2E8F0]'
                      }`}
                    >
                      <span className="text-base font-bold">{val}</span>
                      <span className="text-[9px] text-[#94A3B8]">idx {idx}</span>
                    </div>
                  </div>

                  {/* Forward Arrow */}
                  {idx < arr.length - 1 ? (
                    <div className="flex items-center justify-center text-[#EF4444] px-1">
                      <ArrowRight className="w-5 h-5 animate-pulse" />
                    </div>
                  ) : (
                    <div className="flex items-center gap-1 text-[#94A3B8] text-xs font-bold pl-2">
                      <ArrowRight className="w-4 h-4" />
                      <span>NULL</span>
                    </div>
                  )}
                </React.Fragment>
              );
            })}
          </div>
        ) : (
          /* Render Standard Array / Bars with Pointers */
          <div className="flex items-end justify-center gap-2.5 sm:gap-3.5 h-56 w-full max-w-4xl px-2">
            {arr.map((val, idx) => {
              const isComparing = comparing.has(idx);
              const isSwapping = swapping.has(idx);
              const isHighlight = highlighted.has(idx);
              const ptrs = indexPointers[idx] || [];

              // Calculate bar height percentage (min 20%, max 100%)
              const heightPct = Math.max(25, Math.round((Math.abs(val) / maxVal) * 90));

              // Determine bar color
              let barBg = 'bg-[#E2E8F0] border-[#E2E8F0] text-[#E2E8F0]';
              if (is_complete) {
                barBg = 'bg-emerald-500/20 border-emerald-500 text-emerald-300';
              } else if (isSwapping) {
                barBg = 'bg-rose-500/30 border-rose-500 text-rose-300 scale-105 shadow-lg shadow-rose-500/20';
              } else if (isComparing) {
                barBg = 'bg-amber-500/30 border-amber-500 text-amber-300 scale-105 shadow-lg shadow-amber-500/20';
              } else if (isHighlight) {
                barBg = 'bg-[#EF4444]/30 border-[#EF4444] text-white shadow-lg shadow-[#EF4444]/30';
              }

              return (
                <div key={idx} className="flex-1 flex flex-col items-center justify-end h-full max-w-[54px] min-w-[32px]">
                  {/* Floating Pointer Badges */}
                  <div className="flex flex-col items-center gap-1 mb-2 h-14 justify-end">
                    {ptrs.map((p) => (
                      <span
                        key={p}
                        className="text-[9px] uppercase font-bold bg-[#EF4444] text-white px-1.5 py-0.5 rounded shadow"
                      >
                        {p}
                      </span>
                    ))}
                  </div>

                  {/* Animated Bar */}
                  <div
                    style={{ height: `${heightPct}%` }}
                    className={`w-full rounded-t-xl border flex flex-col items-center justify-between py-2 transition-all duration-300 ${barBg}`}
                  >
                    <span className="text-xs font-bold">{val}</span>
                  </div>

                  {/* Index Label */}
                  <div className="h-6 flex items-center justify-center text-[10px] text-[#94A3B8] font-bold">
                    {idx}
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>

      {/* Legend / Key Guide */}
      <div className="flex flex-wrap items-center justify-center gap-4 pt-3 border-t border-[#E2E8F0] text-[10px] text-[#64748B]">
        <div className="flex items-center gap-1.5">
          <span className="w-2.5 h-2.5 rounded-full bg-amber-400" />
          <span>Comparing</span>
        </div>
        <div className="flex items-center gap-1.5">
          <span className="w-2.5 h-2.5 rounded-full bg-rose-400" />
          <span>Swapping / Moving</span>
        </div>
        <div className="flex items-center gap-1.5">
          <span className="w-2.5 h-2.5 rounded-full bg-[#EF4444]" />
          <span>Active Pointer / Cell</span>
        </div>
        <div className="flex items-center gap-1.5">
          <span className="w-2.5 h-2.5 rounded-full bg-emerald-400" />
          <span>Locked / Sorted</span>
        </div>
      </div>
    </div>
  );
};
export default DataStructureCanvas;

