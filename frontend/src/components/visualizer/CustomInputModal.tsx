import React, { useState } from 'react';
import { X, Dices } from 'lucide-react';
import Button from '../ui/Button';

interface CustomInputModalProps {
  isOpen: boolean;
  onClose: () => void;
  inputType: string;
  defaultInput: Record<string, any>;
  onApply: (newInput: Record<string, any>) => void;
}

export const CustomInputModal: React.FC<CustomInputModalProps> = ({
  isOpen,
  onClose,
  inputType,
  defaultInput,
  onApply,
}) => {
  const [arrayStr, setArrayStr] = useState<string>(
    Array.isArray(defaultInput.array)
      ? defaultInput.array.join(', ')
      : Array.isArray(defaultInput.values)
      ? defaultInput.values.join(', ')
      : '29, 10, 14, 37, 13, 22'
  );
  const [target, setTarget] = useState<number>(defaultInput.target ?? 23);
  const [arraySize, setArraySize] = useState<number>(6);

  if (!isOpen) return null;

  const handleGenerateRandom = () => {
    const generated = Array.from({ length: arraySize }, () =>
      Math.floor(Math.random() * 90) + 5
    );
    if (inputType === 'array_target' || inputType === 'binary_search') {
      generated.sort((a, b) => a - b);
      setTarget(generated[Math.floor(generated.length / 2)]);
    }
    setArrayStr(generated.join(', '));
  };

  const handleSave = () => {
    const parsedNums = arrayStr
      .split(',')
      .map((s) => parseInt(s.trim(), 10))
      .filter((n) => !isNaN(n));

    // Limit to safe rendering bound
    const safeSlice = parsedNums.slice(0, 16);

    const config: Record<string, any> = {};
    if (inputType === 'linked_list') {
      config.values = safeSlice;
    } else {
      config.array = safeSlice;
    }

    if (inputType === 'array_target') {
      config.target = Number(target);
    }

    onApply(config);
    onClose();
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm font-mono">
      <div className="w-full max-w-md bg-[#0D0D14] border border-[#1C1C28] rounded-2xl p-5 flex flex-col gap-4">
        {/* Header */}
        <div className="flex items-center justify-between border-b border-[#1C1C28] pb-3">
          <h3 className="text-sm font-bold text-[#F4F4F6] font-sans">Configure Custom Input</h3>
          <button onClick={onClose} className="p-1.5 text-[#8E8EA0] hover:text-[#F4F4F6] rounded-lg">
            <X className="w-4 h-4" />
          </button>
        </div>

        {/* Array Input */}
        <div className="space-y-1.5">
          <label className="text-[11px] text-[#8E8EA0] font-sans">
            Comma-separated values (Max 16 elements):
          </label>
          <input
            value={arrayStr}
            onChange={(e) => setArrayStr(e.target.value)}
            placeholder="e.g. 15, 8, 42, 4, 16, 23"
            className="w-full bg-[#12121B] text-[#F4F4F6] text-xs px-3 py-2 rounded-lg border border-[#1C1C28] outline-none font-mono"
          />
        </div>

        {/* Optional Target Input */}
        {inputType === 'array_target' && (
          <div className="space-y-1.5">
            <label className="text-[11px] text-[#8E8EA0] font-sans">Search Target Value:</label>
            <input
              type="number"
              value={target}
              onChange={(e) => setTarget(Number(e.target.value))}
              className="w-full bg-[#12121B] text-[#F4F4F6] text-xs px-3 py-2 rounded-lg border border-[#1C1C28] outline-none font-mono"
            />
          </div>
        )}

        {/* Random Generator */}
        <div className="p-3 bg-[#12121B] rounded-xl border border-[#1C1C28] space-y-2">
          <div className="flex items-center justify-between text-[11px] text-[#8E8EA0] font-sans">
            <span>Random Array Size: {arraySize}</span>
            <button
              onClick={handleGenerateRandom}
              className="text-[#A78BFA] hover:text-white flex items-center gap-1 text-[11px]"
            >
              <Dices className="w-3.5 h-3.5" /> Generate
            </button>
          </div>
          <input
            type="range"
            min={4}
            max={12}
            value={arraySize}
            onChange={(e) => setArraySize(Number(e.target.value))}
            className="w-full accent-[#8B5CF6] h-1.5 bg-[#1C1C28] rounded cursor-pointer"
          />
        </div>

        {/* Footer */}
        <div className="flex items-center justify-end gap-2 pt-2 border-t border-[#1C1C28]">
          <Button variant="ghost" size="sm" onClick={onClose} className="text-xs">
            Cancel
          </Button>
          <Button variant="primary" size="sm" onClick={handleSave} className="text-xs">
            Apply Input
          </Button>
        </div>
      </div>
    </div>
  );
};
export default CustomInputModal;

