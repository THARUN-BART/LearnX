import React, { useState } from 'react';
import { CodeIssue, IssueSeverity } from '../../types/codereview.types';
import { AlertCircle, AlertTriangle, Info, CheckCircle } from 'lucide-react';

interface IssueListProps {
  issues: CodeIssue[];
}

export const IssueList: React.FC<IssueListProps> = ({ issues }) => {
  const [filter, setFilter] = useState<string>('all');

  const filteredIssues = issues.filter((iss) => {
    if (filter === 'all') return true;
    return iss.severity === filter;
  });

  const getSeverityIcon = (sev: IssueSeverity) => {
    if (sev === 'critical') return <AlertCircle className="w-4 h-4 text-rose-400 shrink-0 mt-0.5" />;
    if (sev === 'warning') return <AlertTriangle className="w-4 h-4 text-amber-400 shrink-0 mt-0.5" />;
    return <Info className="w-4 h-4 text-blue-400 shrink-0 mt-0.5" />;
  };

  const getSeverityBadge = (sev: IssueSeverity) => {
    if (sev === 'critical') return 'bg-rose-500/10 text-rose-400 border-rose-500/20';
    if (sev === 'warning') return 'bg-amber-500/10 text-amber-400 border-amber-500/20';
    return 'bg-blue-500/10 text-blue-400 border-blue-500/20';
  };

  return (
    <div className="bg-[#0D0D14] border border-[#1C1C28] rounded-2xl p-5 flex flex-col gap-4 font-mono">
      {/* Header and Filter Tabs */}
      <div className="flex flex-wrap items-center justify-between gap-3 border-b border-[#1C1C28] pb-3">
        <div className="flex items-center gap-2">
          <span className="text-sm font-bold text-white font-sans">
            Detected Code Issues & Warnings
          </span>
          <span className="text-[10px] text-[#A78BFA] bg-[#8B5CF6]/10 px-2 py-0.5 rounded border border-[#8B5CF6]/20">
            {issues.length} Total
          </span>
        </div>

        <div className="flex items-center gap-1.5 text-xs font-sans">
          {['all', 'critical', 'warning', 'info'].map((f) => (
            <button
              key={f}
              onClick={() => setFilter(f)}
              className={`px-2.5 py-1 rounded-lg capitalize transition-colors ${
                filter === f
                  ? 'bg-[#8B5CF6] text-white font-bold'
                  : 'bg-[#12121B] text-[#8E8EA0] hover:text-[#F4F4F6] border border-[#1C1C28]'
              }`}
            >
              {f}
            </button>
          ))}
        </div>
      </div>

      {/* Issues List */}
      <div className="space-y-3">
        {filteredIssues.length === 0 ? (
          <div className="py-8 text-center text-xs text-emerald-400 font-sans flex flex-col items-center gap-2">
            <CheckCircle className="w-5 h-5 text-emerald-400" />
            <span>No issues found in this severity category. Great work!</span>
          </div>
        ) : (
          filteredIssues.map((iss) => (
            <div
              key={iss.id}
              className="p-4 bg-[#12121B] border border-[#1C1C28] rounded-xl flex flex-col gap-2.5"
            >
              <div className="flex items-start justify-between gap-3">
                <div className="flex items-start gap-2.5">
                  {getSeverityIcon(iss.severity)}
                  <div>
                    <div className="flex items-center gap-2">
                      <span className="text-xs font-bold text-[#F4F4F6] font-sans">
                        {iss.title}
                      </span>
                      {iss.line_number > 0 && (
                        <span className="text-[10px] text-[#8E8EA0] bg-[#181824] px-1.5 py-0.2 rounded border border-[#1C1C28]">
                          Line {iss.line_number}
                        </span>
                      )}
                    </div>
                    <p className="text-xs text-[#CBD5E1] font-sans mt-1 leading-relaxed">
                      {iss.description}
                    </p>
                  </div>
                </div>

                <div className="flex flex-col items-end gap-1 shrink-0">
                  <span
                    className={`text-[9px] uppercase font-bold px-2 py-0.5 rounded border ${getSeverityBadge(
                      iss.severity
                    )}`}
                  >
                    {iss.severity}
                  </span>
                  <span className="text-[9px] text-[#5A5A6E] capitalize font-mono">
                    {iss.source_type.replace(/_/g, ' ')}
                  </span>
                </div>
              </div>

              {/* Code Snippet */}
              {iss.code_snippet && (
                <pre className="p-2.5 bg-[#07070B] border border-[#1C1C28] rounded-lg text-[11px] text-[#E2E8F0] overflow-x-auto whitespace-pre-wrap">
                  {iss.code_snippet}
                </pre>
              )}

              {/* Recommendation */}
              <div className="p-2.5 bg-[#8B5CF6]/5 border border-[#8B5CF6]/20 rounded-lg text-xs font-sans text-[#D8B4FE]">
                <span className="font-bold text-[#A78BFA]">Recommendation: </span>
                {iss.recommendation}
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
};
export default IssueList;

