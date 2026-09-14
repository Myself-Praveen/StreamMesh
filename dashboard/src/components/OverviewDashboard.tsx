'use client';

import { useMetricsStream } from '@/hooks/useMetricsStream';
import { MetricsCards } from './MetricsCards';
import { ThroughputChart } from './ThroughputChart';
import { LatencyChart } from './LatencyChart';
import { AlertCircle } from 'lucide-react';

export function OverviewDashboard() {
  // Assuming the StreamMesh backend is running on 8080 locally for demo
  const { metrics, history, error, isConnected } = useMetricsStream('http://localhost:8080/api/metrics/stream');

  return (
    <div className="space-y-6">
      {!isConnected && !error && (
        <div className="rounded-md bg-blue-50 p-4 dark:bg-blue-900/20">
          <div className="flex">
            <div className="ml-3">
              <h3 className="text-sm font-medium text-blue-800 dark:text-blue-300">Connecting to StreamMesh...</h3>
            </div>
          </div>
        </div>
      )}

      {error && (
        <div className="rounded-md bg-red-50 p-4 dark:bg-red-900/20">
          <div className="flex">
            <div className="flex-shrink-0">
              <AlertCircle className="h-5 w-5 text-red-400" aria-hidden="true" />
            </div>
            <div className="ml-3">
              <h3 className="text-sm font-medium text-red-800 dark:text-red-300">Connection Error</h3>
              <div className="mt-2 text-sm text-red-700 dark:text-red-400">
                <p>Failed to connect to the StreamMesh metrics stream. Is the backend running on port 8080?</p>
              </div>
            </div>
          </div>
        </div>
      )}

      <MetricsCards metrics={metrics} />

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        <ThroughputChart history={history} />
        <LatencyChart history={history} />
      </div>
    </div>
  );
}
