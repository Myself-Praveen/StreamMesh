import Link from 'next/link';
import { LayoutDashboard, Radio, Users, Activity, Settings } from 'lucide-react';

const navigation = [
  { name: 'Overview', href: '/', icon: LayoutDashboard },
  { name: 'Channels', href: '/channels', icon: Radio },
  { name: 'Connections', href: '/connections', icon: Users },
  { name: 'Cluster', href: '/cluster', icon: Activity },
  { name: 'Settings', href: '/settings', icon: Settings },
];

export function Sidebar() {
  return (
    <div className="flex h-full w-64 flex-col border-r border-gray-200 bg-white/50 backdrop-blur-xl dark:border-gray-800 dark:bg-gray-900/50">
      <div className="flex h-16 items-center gap-2 border-b border-gray-200 px-6 dark:border-gray-800">
        <div className="h-8 w-8 rounded-lg bg-blue-600 flex items-center justify-center text-white font-bold">
          S
        </div>
        <span className="text-lg font-bold tracking-tight">StreamMesh</span>
      </div>
      <nav className="flex-1 space-y-1 px-3 py-4">
        {navigation.map((item) => {
          const Icon = item.icon;
          return (
            <Link
              key={item.name}
              href={item.href}
              className="group flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-800 transition-colors"
            >
              <Icon className="h-5 w-5 flex-shrink-0 text-gray-400 group-hover:text-gray-500 dark:group-hover:text-gray-300" />
              {item.name}
            </Link>
          );
        })}
      </nav>
      <div className="border-t border-gray-200 p-4 dark:border-gray-800">
        <div className="flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium text-green-600 bg-green-50 dark:bg-green-900/20 dark:text-green-400">
          <div className="h-2 w-2 rounded-full bg-green-500 animate-pulse" />
          Cluster Online
        </div>
      </div>
    </div>
  );
}
