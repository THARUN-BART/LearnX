import React from 'react';
import { ShieldCheck, AlertTriangle, ArrowUpRight } from 'lucide-react';
import { Link } from 'react-router-dom';
import { ConceptStrength, ConceptWeakness } from '../../types/learning-memory.types';

interface Props {
  strengths: ConceptStrength[];
  weaknesses: ConceptWeakness[];
}

export const StrengthsWeaknessesSection: React.FC<Props> = ({ strengths, weaknesses }) => {
  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
      {/* Demonstrated Strengths */}
      <div className="p-5 rounded-2xl bg-[#0D0D14] border border-emerald-500/20 bg-gradient-to-b from-emerald-500/5 to-transparent flex flex-col justify-between">
        <div>
          <div className="flex items-center gap-2 mb-3">
            <div className="p-1.5 rounded-lg bg-emerald-500/10 text-emerald-400">
              <ShieldCheck className="w-5 h-5" />
            </div>
            <div>
              <h3 className="text-sm font-bold text-[#F4F4F6]">Demonstrated Strengths</h3>
              <p className="text-xs text-[#8E8EA0]">Concepts backed by consistent unassisted solves</p>
            </div>
          </div>

          <div className="space-y-3 mt-4">
            {strengths.length === 0 ? (
              <p className="text-xs text-[#8E8EA0] italic py-4">
                Solve problems independently without hints to establish verified strengths.
              </p>
            ) : (
              strengths.map((s) => (
                <div key={s.concept_key} className="p-3 rounded-xl bg-[#12121B] border border-[#1C1C28] flex items-center justify-between">
                  <div>
                    <h4 className="text-xs font-semibold text-[#F4F4F6]">{s.concept_name}</h4>
                    <p className="text-[11px] text-[#8E8EA0] mt-0.5">{s.reason}</p>
                  </div>
                  <Link
                    to={`/dashboard/learning-memory/concepts/${encodeURIComponent(s.concept_key)}`}
                    className="text-[#8E8EA0] hover:text-emerald-400 p-1"
                  >
                    <ArrowUpRight className="w-4 h-4" />
                  </Link>
                </div>
              ))
            )}
          </div>
        </div>
      </div>

      {/* Focus Areas (Weaknesses) */}
      <div className="p-5 rounded-2xl bg-[#0D0D14] border border-rose-500/20 bg-gradient-to-b from-rose-500/5 to-transparent flex flex-col justify-between">
        <div>
          <div className="flex items-center gap-2 mb-3">
            <div className="p-1.5 rounded-lg bg-rose-500/10 text-rose-400">
              <AlertTriangle className="w-5 h-5" />
            </div>
            <div>
              <h3 className="text-sm font-bold text-[#F4F4F6]">Areas for Targeted Focus</h3>
              <p className="text-xs text-[#8E8EA0]">Mistake patterns detected across multiple submissions</p>
            </div>
          </div>

          <div className="space-y-3 mt-4">
            {weaknesses.length === 0 ? (
              <p className="text-xs text-[#8E8EA0] italic py-4">
                No active recurring weaknesses! Your error patterns are currently well-balanced.
              </p>
            ) : (
              weaknesses.map((w) => (
                <div key={w.concept_key + w.mistake_subtype} className="p-3 rounded-xl bg-[#12121B] border border-[#1C1C28] flex items-center justify-between">
                  <div>
                    <div className="flex items-center gap-2">
                      <h4 className="text-xs font-semibold text-[#F4F4F6]">{w.concept_name}</h4>
                      <span className="px-2 py-0.2 rounded-full text-[9px] font-mono uppercase bg-rose-500/10 text-rose-400 border border-rose-500/20">
                        {w.mistake_subtype}
                      </span>
                    </div>
                    <p className="text-[11px] text-[#8E8EA0] mt-0.5">{w.reason}</p>
                  </div>
                  <Link
                    to={`/dashboard/learning-memory/concepts/${encodeURIComponent(w.concept_key)}`}
                    className="text-[#8E8EA0] hover:text-rose-400 p-1"
                  >
                    <ArrowUpRight className="w-4 h-4" />
                  </Link>
                </div>
              ))
            )}
          </div>
        </div>
      </div>
    </div>
  );
};
