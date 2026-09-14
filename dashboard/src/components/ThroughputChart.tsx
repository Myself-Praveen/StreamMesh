'use client';

import { MetricsSnapshot } from '@/hooks/useMetricsStream';
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from 'recharts';

interface ThroughputChartProps {
  history: MetricsSnapshot[];
}

export function ThroughputChart({ history }: ThroughputChartProps) {
  // Transform history into data points
  // We need to calculate rate (messages per second) from cumulative counts
  const data = history.map((snapshot, index) => {
    let msgRate = 0;
    if (index > 0) {
      const prev = history[index - 1];
      msgRate = snapshot.global.messages_delivered - prev.global.messages_delivered;
      // Ensure no negative values if server restarted
      if (msgRate < 0) msgRate = 0;
    }

    return {
      time: snapshot.timestamp ? new Date(snapshot.timestamp).toLocaleTimeString() : '',
      rate: msgRate,
    };
  });

  return (
    <div className="rounded-xl bg-white p-6 shadow-sm ring-1 ring-black/5 dark:bg-gray-900 dark:ring-white/10 h-[400px]">
      <h3 className="text-base font-semibold leading-6 text-gray-900 dark:text-white mb-6">
        Message Throughput (msg/sec)
      </h3>
      <ResponsiveContainer width="100%" height="100%">
        <AreaChart
          data={data}
          margin={{ top: 10, right: 30, left: 0, bottom: 0 }}
        >
          <defs>
            <linearGradient id="colorRate" x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor="#8b5cf6" stopOpacity={0.3} />
              <stop offset="95%" stopColor="#8b5cf6" stopOpacity={0} />
            </linearGradient>
          </defs>
          <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#374151" opacity={0.2} />
          <XAxis 
            dataKey="time" 
            axisLine={false} 
            tickLine={false} 
            tick={{ fill: '#6b7280', fontSize: 12 }} 
            dy={10} 
            minTickGap={30}
          />
          <YAxis 
            axisLine={false} 
            tickLine={false} 
            tick={{ fill: '#6b7280', fontSize: 12 }} 
          />
          <Tooltip 
            contentStyle={{ borderRadius: '8px', border: 'none', boxShadow: '0 4px 6px -1px rgb(0 0 0 / 0.1)' }}
          />
          <Area
            type="monotone"
            dataKey="rate"
            stroke="#8b5cf6"
            strokeWidth={2}
            fillOpacity={1}
            fill="url(#colorRate)"
            isAnimationActive={false} // Disable animation for smoother live streaming
          />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
}
