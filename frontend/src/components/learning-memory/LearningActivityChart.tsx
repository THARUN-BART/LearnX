import React from 'react';
import {
  ResponsiveContainer,
  AreaChart,
  Area,
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

export const LearningActivityChart: React.FC<Props> = ({ data }) => {
  if (!data || data.length === 0) {
    return (
      <div className="h-64 flex items-center justify-center text-xs text-[#8E8EA0] italic bg-[#0D0D14] rounded-2xl border border-[#1C1C28]">
        No activity history recorded for this date range.
      </div>
    );
  }

  return (
    <div className="p-5 rounded-2xl bg-[#0D0D14] border border-[#1C1C28]">
      <div className="flex items-center justify-between mb-4">
        <div>
          <h3 className="text-sm font-bold text-[#F4F4F6]">Learning Trajectory & Solves</h3>
          <p className="text-xs text-[#8E8EA0]">Independent vs assisted solves over recent sessions</p>
        </div>
      </div>

      <div className="h-64 w-full">
        <ResponsiveContainer width="100%" height="100%">
          <AreaChart data={data} margin={{ top: 10, right: 10, left: -20, bottom: 0 }}>
            <defs>
              <linearGradient id="indGrad" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="#10B981" stopOpacity={0.4} />
                <stop offset="95%" stopColor="#10B981" stopOpacity={0} />
              </linearGradient>
              <linearGradient id="assGrad" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="#3B82F6" stopOpacity={0.4} />
                <stop offset="95%" stopColor="#3B82F6" stopOpacity={0} />
              </linearGradient>
            </defs>
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
              formatter={(val) => <span className="text-[#8E8EA0] capitalize">{val}</span>}
            />
            <Area
              type="monotone"
              dataKey="independent_solves"
              name="Independent Solves"
              stroke="#10B981"
              strokeWidth={2}
              fillOpacity={1}
              fill="url(#indGrad)"
            />
            <Area
              type="monotone"
              dataKey="assisted_solves"
              name="Assisted Solves"
              stroke="#3B82F6"
              strokeWidth={2}
              fillOpacity={1}
              fill="url(#assGrad)"
            />
          </AreaChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
};
