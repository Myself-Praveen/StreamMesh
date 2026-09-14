import { useEffect, useState } from 'react';

export interface GlobalMetrics {
  active_connections: number;
  messages_published: number;
  messages_delivered: number;
  bytes_sent: number;
  bytes_received: number;
}

export interface TopicMetrics {
  messages_published: number;
  messages_delivered: number;
  subscribers: number;
}

export interface LatencyMetrics {
  count: number;
  min: number;
  max: number;
  mean: number;
  p50: number;
  p90: number;
  p99: number;
}

export interface SystemMetrics {
  NumGoroutines: number;
  AllocBytes: number;
  SysBytes: number;
  NumGC: number;
}

export interface MetricsSnapshot {
  global: GlobalMetrics;
  topics: Record<string, TopicMetrics>;
  latencies: LatencyMetrics;
  system: SystemMetrics;
  timestamp?: number;
}

export function useMetricsStream(url: string = 'http://localhost:8080/api/metrics/stream') {
  const [metrics, setMetrics] = useState<MetricsSnapshot | null>(null);
  const [history, setHistory] = useState<MetricsSnapshot[]>([]);
  const [isConnected, setIsConnected] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    const eventSource = new EventSource(url);

    eventSource.onopen = () => {
      setIsConnected(true);
      setError(null);
    };

    eventSource.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data) as MetricsSnapshot;
        data.timestamp = Date.now();
        
        setMetrics(data);
        setHistory((prev) => {
          const newHistory = [...prev, data];
          // Keep last 60 data points (1 minute of history at 1 tick/sec)
          if (newHistory.length > 60) {
            newHistory.shift();
          }
          return newHistory;
        });
      } catch (err) {
        console.error('Failed to parse metrics stream:', err);
      }
    };

    eventSource.onerror = (err) => {
      setIsConnected(false);
      setError(new Error('EventSource failed'));
      console.error('Metrics stream error:', err);
      eventSource.close();
    };

    return () => {
      eventSource.close();
      setIsConnected(false);
    };
  }, [url]);

  return { metrics, history, isConnected, error };
}
