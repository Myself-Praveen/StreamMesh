'use client';

import { useState, useEffect } from 'react';
import { Server, Activity, Clock } from 'lucide-react';

interface Node {
  id: string;
  hostname: string;
  start_time: string;
}

export default function ClusterPage() {
  const [nodes, setNodes] = useState<Node[]>([]);
  const [selfNode, setSelfNode] = useState<string>('');
  const [loading, setLoading] = useState(true);

  const fetchNodes = async () => {
    try {
      const res = await fetch('http://localhost:8080/api/admin/cluster');
      const data = await res.json();
      setNodes(data.nodes || []);
      setSelfNode(data.self || '');
    } catch (err) {
      console.error('Failed to fetch cluster nodes:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchNodes();
    const interval = setInterval(fetchNodes, 5000); // refresh every 5s
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
          Cluster Topology
        </h1>
        <div className="flex items-center space-x-2">
          <span className="relative flex h-3 w-3">
            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
            <span className="relative inline-flex rounded-full h-3 w-3 bg-green-500"></span>
          </span>
          <span className="text-sm font-medium text-gray-600 dark:text-gray-400">{nodes.length} Nodes Online</span>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {loading && nodes.length === 0 ? (
          <div className="col-span-full py-12 text-center text-sm text-gray-500">Discovering cluster nodes...</div>
        ) : (
          nodes.map((node) => (
            <div key={node.id} className={`relative overflow-hidden rounded-xl bg-white shadow-sm ring-1 dark:bg-gray-900 p-6 ${node.id === selfNode ? 'ring-blue-500 dark:ring-blue-500 shadow-md' : 'ring-black/5 dark:ring-white/10'}`}>
              {node.id === selfNode && (
                <div className="absolute top-0 right-0 -mt-1 -mr-1">
                  <span className="inline-flex items-center px-2.5 py-0.5 rounded-bl-lg text-xs font-medium bg-blue-100 text-blue-800 dark:bg-blue-900/40 dark:text-blue-300">
                    Current Node
                  </span>
                </div>
              )}
              
              <div className="flex items-center space-x-4">
                <div className="flex-shrink-0 bg-blue-50 dark:bg-blue-900/20 p-3 rounded-lg">
                  <Server className="h-6 w-6 text-blue-600 dark:text-blue-400" />
                </div>
                <div>
                  <h3 className="text-sm font-semibold text-gray-900 dark:text-white">{node.hostname}</h3>
                  <p className="text-xs text-gray-500 dark:text-gray-400 font-mono mt-1 break-all">{node.id}</p>
                </div>
              </div>

              <div className="mt-6 border-t border-gray-100 dark:border-gray-800 pt-4">
                <dl className="grid grid-cols-2 gap-4">
                  <div>
                    <dt className="flex items-center text-xs font-medium text-gray-500 dark:text-gray-400">
                      <Activity className="h-3.5 w-3.5 mr-1" />
                      Status
                    </dt>
                    <dd className="mt-1 flex items-center text-sm font-medium text-green-600 dark:text-green-400">
                      Healthy
                    </dd>
                  </div>
                  <div>
                    <dt className="flex items-center text-xs font-medium text-gray-500 dark:text-gray-400">
                      <Clock className="h-3.5 w-3.5 mr-1" />
                      Uptime
                    </dt>
                    <dd className="mt-1 text-sm font-medium text-gray-900 dark:text-white">
                      {Math.floor((new Date().getTime() - new Date(node.start_time).getTime()) / 3600000)}h {Math.floor(((new Date().getTime() - new Date(node.start_time).getTime()) % 3600000) / 60000)}m
                    </dd>
                  </div>
                </dl>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
