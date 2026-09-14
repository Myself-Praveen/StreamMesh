'use client';

import { MetricsSnapshot } from '@/hooks/useMetricsStream';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';

interface LatencyChartProps {
  history: MetricsSnapshot[];
}

export function LatencyChart({ history }: LatencyChartProps) {
  const data = history.map((snapshot) => {
    return {
      time: snapshot.timestamp ? new Date(snapshot.timestamp).toLocaleTimeString() : '',
      p50: snapshot.latencies.p50,
      p90: snapshot.latencies.p90,
      p99: snapshot.latencies.p99,
    };
  });

  return (
    <div className="rounded-xl bg-white p-6 shadow-sm ring-1 ring-black/5 dark:bg-gray-900 dark:ring-white/10 h-[400px]">
      <h3 className="text-base font-semibold leading-6 text-gray-900 dark:text-white mb-6">
        Latency Distribution (ms)
      </h3>
      <ResponsiveContainer width="100%" height="100%">
        <LineChart
          data={data}
          margin={{ top: 10, right: 30, left: 0, bottom: 0 }}
        >
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
          <Legend wrapperStyle={{ paddingTop: '20px' }} />
          <Line
            type="monotone"
            dataKey="p50"
            name="p50"
            stroke="#10b981"
            strokeWidth={2}
            dot={false}
            isAnimationActive={false}
          />
          <Line
            type="monotone"
            dataKey="p90"
            name="p90"
            stroke="#f59e0b"
            strokeWidth={2}
            dot={false}
            isAnimationActive={false}
          />
          <Line
            type="monotone"
            dataKey="p99"
            name="p99"
            stroke="#ef4444"
            strokeWidth={2}
            dot={false}
            isAnimationActive={false}
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}
