import { OverviewDashboard } from '@/components/OverviewDashboard';

export default function Home() {
  return (
    <div className="flex flex-col space-y-8">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
          Overview
        </h1>
        <div className="flex items-center space-x-2">
          <span className="relative flex h-3 w-3">
            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
            <span className="relative inline-flex rounded-full h-3 w-3 bg-green-500"></span>
          </span>
          <span className="text-sm font-medium text-gray-600 dark:text-gray-400">Live</span>
        </div>
      </div>
      
      <OverviewDashboard />
    </div>
  );
}
