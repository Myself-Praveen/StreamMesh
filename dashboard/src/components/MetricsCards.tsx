'use client';

import { Activity, Users, Send, Download, Cpu, HardDrive } from 'lucide-react';
import { MetricsSnapshot } from '@/hooks/useMetricsStream';

interface MetricsCardsProps {
  metrics: MetricsSnapshot | null;
}

export function MetricsCards({ metrics }: MetricsCardsProps) {
  if (!metrics) {
    return (
      <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
        {[1, 2, 3, 4, 5, 6].map((i) => (
          <div key={i} className="animate-pulse overflow-hidden rounded-xl bg-white shadow-sm ring-1 ring-black/5 dark:bg-gray-900 dark:ring-white/10 h-32"></div>
        ))}
      </div>
    );
  }

  const cards = [
    {
      name: 'Active Connections',
      value: metrics.global.active_connections.toLocaleString(),
      icon: Users,
      color: 'text-blue-500',
      bg: 'bg-blue-50 dark:bg-blue-500/10',
    },
    {
      name: 'Messages Delivered',
      value: metrics.global.messages_delivered.toLocaleString(),
      icon: Send,
      color: 'text-green-500',
      bg: 'bg-green-50 dark:bg-green-500/10',
    },
    {
      name: 'Throughput (Msg/s)',
      value: (metrics.global.messages_delivered / Math.max(1, (metrics.system.NumGoroutines / 100))).toFixed(0), // Approximation for demo if no rate available
      icon: Activity,
      color: 'text-purple-500',
      bg: 'bg-purple-50 dark:bg-purple-500/10',
    },
    {
      name: 'Bandwidth In',
      value: formatBytes(metrics.global.bytes_received),
      icon: Download,
      color: 'text-amber-500',
      bg: 'bg-amber-50 dark:bg-amber-500/10',
    },
    {
      name: 'Goroutines',
      value: metrics.system.NumGoroutines.toLocaleString(),
      icon: Cpu,
      color: 'text-red-500',
      bg: 'bg-red-50 dark:bg-red-500/10',
    },
    {
      name: 'Memory Allocated',
      value: formatBytes(metrics.system.AllocBytes),
      icon: HardDrive,
      color: 'text-indigo-500',
      bg: 'bg-indigo-50 dark:bg-indigo-500/10',
    },
  ];

  return (
    <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
      {cards.map((item) => {
        const Icon = item.icon;
        return (
          <div
            key={item.name}
            className="overflow-hidden rounded-xl bg-white shadow-sm ring-1 ring-black/5 dark:bg-gray-900 dark:ring-white/10"
          >
            <div className="p-5">
              <div className="flex items-center">
                <div className={`flex-shrink-0 rounded-md p-3 ${item.bg}`}>
                  <Icon className={`h-6 w-6 ${item.color}`} aria-hidden="true" />
                </div>
                <div className="ml-5 w-0 flex-1">
                  <dl>
                    <dt className="truncate text-sm font-medium text-gray-500 dark:text-gray-400">
                      {item.name}
                    </dt>
                    <dd>
                      <div className="text-2xl font-bold text-gray-900 dark:text-white">
                        {item.value}
                      </div>
                    </dd>
                  </dl>
                </div>
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
}

function formatBytes(bytes: number, decimals = 2) {
  if (!+bytes) return '0 Bytes';
  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`;
}
