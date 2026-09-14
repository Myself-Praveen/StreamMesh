'use client';

import { useState, useEffect } from 'react';
import { useMetricsStream } from '@/hooks/useMetricsStream';
import { Search, Radio, Users } from 'lucide-react';

export default function ChannelsPage() {
  const [searchTerm, setSearchTerm] = useState('');
  const { metrics } = useMetricsStream('http://localhost:8080/api/metrics/stream');

  const topics = metrics?.topics || {};
  const topicList = Object.entries(topics).map(([name, stats]) => ({
    name,
    ...stats,
  }));

  const filteredTopics = topicList.filter(t => t.name.toLowerCase().includes(searchTerm.toLowerCase()));

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
          Channels
        </h1>
      </div>

      <div className="flex items-center px-4 py-3 bg-white dark:bg-gray-900 rounded-lg shadow-sm ring-1 ring-black/5 dark:ring-white/10">
        <Search className="h-5 w-5 text-gray-400 mr-3" />
        <input
          type="text"
          placeholder="Search channels..."
          className="flex-1 bg-transparent border-0 focus:ring-0 text-sm text-gray-900 dark:text-white placeholder-gray-500"
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
        />
      </div>

      <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm ring-1 ring-black/5 dark:ring-white/10 overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-800">
          <thead className="bg-gray-50 dark:bg-gray-800/50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Channel Name</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Subscribers</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Messages Delivered</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
            </tr>
          </thead>
          <tbody className="bg-white dark:bg-gray-900 divide-y divide-gray-200 dark:divide-gray-800">
            {filteredTopics.map((topic) => (
              <tr key={topic.name} className="hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors">
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="flex items-center">
                    <Radio className="h-5 w-5 text-purple-500 mr-3" />
                    <span className="text-sm font-medium text-gray-900 dark:text-white">{topic.name}</span>
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="flex items-center text-sm text-gray-500 dark:text-gray-400">
                    <Users className="h-4 w-4 mr-2" />
                    {topic.subscribers}
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                  {topic.messages_delivered.toLocaleString()}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${topic.subscribers > 0 ? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400' : 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-400'}`}>
                    {topic.subscribers > 0 ? 'Active' : 'Idle'}
                  </span>
                </td>
              </tr>
            ))}
            {filteredTopics.length === 0 && (
              <tr>
                <td colSpan={4} className="px-6 py-8 text-center text-sm text-gray-500 dark:text-gray-400">
                  No channels found matching "{searchTerm}"
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
