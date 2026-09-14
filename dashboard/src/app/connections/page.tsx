'use client';

import { useState, useEffect } from 'react';
import { Users, Search, Activity, PowerOff } from 'lucide-react';

interface Connection {
  id: string;
  user_id: string;
  ip_address: string;
  connected: string;
  metadata: Record<string, string>;
}

export default function ConnectionsPage() {
  const [connections, setConnections] = useState<Connection[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');

  const fetchConnections = async () => {
    try {
      const res = await fetch('http://localhost:8080/api/admin/connections');
      const data = await res.json();
      setConnections(data.connections || []);
    } catch (err) {
      console.error('Failed to fetch connections:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchConnections();
    const interval = setInterval(fetchConnections, 2000);
    return () => clearInterval(interval);
  }, []);

  const filtered = connections.filter(c => 
    c.id.includes(searchTerm) || 
    (c.user_id && c.user_id.includes(searchTerm)) || 
    (c.ip_address && c.ip_address.includes(searchTerm))
  );

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
          Active Connections
        </h1>
        <div className="text-sm text-gray-500 dark:text-gray-400">
          Total: {connections.length}
        </div>
      </div>

      <div className="flex items-center px-4 py-3 bg-white dark:bg-gray-900 rounded-lg shadow-sm ring-1 ring-black/5 dark:ring-white/10">
        <Search className="h-5 w-5 text-gray-400 mr-3" />
        <input
          type="text"
          placeholder="Search by ID, User ID, or IP..."
          className="flex-1 bg-transparent border-0 focus:ring-0 text-sm text-gray-900 dark:text-white placeholder-gray-500"
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
        />
      </div>

      <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm ring-1 ring-black/5 dark:ring-white/10 overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-800">
          <thead className="bg-gray-50 dark:bg-gray-800/50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ID</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">User ID / IP</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Uptime</th>
              <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody className="bg-white dark:bg-gray-900 divide-y divide-gray-200 dark:divide-gray-800">
            {loading && connections.length === 0 ? (
              <tr>
                <td colSpan={4} className="px-6 py-8 text-center text-sm text-gray-500">Loading connections...</td>
              </tr>
            ) : filtered.length === 0 ? (
              <tr>
                <td colSpan={4} className="px-6 py-8 text-center text-sm text-gray-500">No active connections found.</td>
              </tr>
            ) : (
              filtered.map((conn) => (
                <tr key={conn.id} className="hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors">
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center">
                      <div className="h-2 w-2 rounded-full bg-green-500 mr-3"></div>
                      <span className="text-sm font-medium text-gray-900 dark:text-white font-mono">{conn.id}</span>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm text-gray-900 dark:text-white">{conn.user_id || 'Anonymous'}</div>
                    <div className="text-xs text-gray-500 dark:text-gray-400">{conn.ip_address || 'Unknown IP'}</div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                    {Math.floor((new Date().getTime() - new Date(conn.connected).getTime()) / 1000)}s
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                    <button className="text-red-600 hover:text-red-900 dark:text-red-400 dark:hover:text-red-300 transition-colors p-2 rounded-md hover:bg-red-50 dark:hover:bg-red-900/20" title="Disconnect">
                      <PowerOff className="h-4 w-4" />
                    </button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
