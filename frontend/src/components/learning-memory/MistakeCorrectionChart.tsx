import React from 'react';
import {
  ResponsiveContainer,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  CartesianGrid,
  Legend,
} from 'recharts';
import { TimelinePoint } from '../../types/learning-memory.types';

interface Props {
  data: TimelinePoint[];
}

export const MistakeCorrectionChart: React.FC<Props> = ({ data }) => {
  if (!data || data.length === 0) {
    return (
      <div className="h-64 flex items-center justify-center text-xs text-[#64748B] italic bg-white rounded-2xl border border-[#E2E8F0]">
        No mistake telemetry available.
      </div>
    );
  }

  return (
    <div className="p-5 rounded-2xl bg-white border border-[#E2E8F0]">
      <div className="flex items-center justify-between mb-4">
        <div>
          <h3 className="text-sm font-bold text-[#0F172A]">Mistake Detection vs Resolution</h3>
          <p className="text-xs text-[#64748B]">Monitoring whether errors are successfully corrected</p>
        </div>
      </div>

      <div className="h-64 w-full">
        <ResponsiveContainer width="100%" height="100%">
          <BarChart data={data} margin={{ top: 10, right: 10, left: -20, bottom: 0 }}>
            <CartesianGrid strokeDasharray="3 3" stroke="#1C1C28" vertical={false} />
            <XAxis
              dataKey="date"
              stroke="#5A5A6E"
              tick={{ fontSize: 10, fill: '#8E8EA0' }}
              tickFormatter={(v) => v.slice(5)}
            />
            <YAxis stroke="#5A5A6E" tick={{ fontSize: 10, fill: '#8E8EA0' }} allowDecimals={false} />
            <Tooltip
              contentStyle={{
                backgroundColor: '#12121B',
                borderColor: '#1C1C28',
                borderRadius: '0.75rem',
                fontSize: '11px',
                color: '#F4F4F6',
              }}
            />
            <Legend
              wrapperStyle={{ fontSize: '11px', paddingTop: '10px' }}
              formatter={(val) => <span className="text-[#64748B] capitalize">{val}</span>}
            />
            <Bar
              dataKey="mistakes_detected"
              name="Mistakes Detected"
              fill="#F43F5E"
              radius={[4, 4, 0, 0]}
            />
            <Bar
              dataKey="mistakes_corrected"
              name="Mistakes Corrected"
              fill="#10B981"
              radius={[4, 4, 0, 0]}
            />
          </BarChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
};
