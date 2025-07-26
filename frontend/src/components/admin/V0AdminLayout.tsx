import React, { useState } from 'react';
import { Menu, X } from 'lucide-react';
import V0Sidebar from './V0Sidebar';

interface V0AdminLayoutProps {
  children: React.ReactNode;
  currentPath: string;
}

/**
 * V0AdminLayout - Main layout component for admin dashboard
 * Matches the v0.dev design with proper Tailwind styling
 */
const V0AdminLayout: React.FC<V0AdminLayoutProps> = ({ children, currentPath }) => {
  const [sidebarOpen, setSidebarOpen] = useState(false);

  return (
    <div className="admin-layout min-h-screen bg-gray-50">
      {/* Mobile sidebar overlay */}
      {sidebarOpen && (
        <div className="fixed inset-0 z-40 lg:hidden">
          <div className="fixed inset-0 bg-gray-600 bg-opacity-75" onClick={() => setSidebarOpen(false)} />
          <div className="relative flex-1 flex flex-col max-w-xs w-full bg-white">
            <div className="absolute top-0 right-0 -mr-12 pt-2">
              <button
                className="ml-1 flex items-center justify-center h-10 w-10 rounded-full focus:outline-none focus:ring-2 focus:ring-inset focus:ring-white"
                onClick={() => setSidebarOpen(false)}
              >
                <X className="h-6 w-6 text-white" />
              </button>
            </div>
            <V0Sidebar currentPath={currentPath} />
          </div>
        </div>
      )}

      {/* Main layout container with sidebar and content */}
      <div className="flex h-screen overflow-hidden">
        {/* Desktop Sidebar */}
        <V0Sidebar currentPath={currentPath} />
        {/* Main content area */}
        <div className="flex-1 flex flex-col overflow-hidden">
          {/* Header area - can be extended later */}
          <header className="bg-white border-b border-gray-200 px-6 py-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center">
                <button
                  className="lg:hidden mr-4 text-gray-500 hover:text-gray-700"
                  onClick={() => setSidebarOpen(true)}
                >
                  <Menu className="h-6 w-6" />
                </button>
                <h1 className="text-2xl font-semibold text-gray-900">
                  {getPageTitle(currentPath)}
                </h1>
              </div>
              <div className="flex items-center space-x-4">
                {/* User profile section - can be extended later */}
                <div className="flex items-center space-x-2">
                  <div className="w-8 h-8 bg-blue-500 rounded-full flex items-center justify-center">
                    <span className="text-white text-sm font-medium">A</span>
                  </div>
                  <span className="text-sm text-gray-700">Admin User</span>
                </div>
              </div>
            </div>
          </header>

          {/* Main content */}
          <main className="flex-1 overflow-y-auto bg-gray-50 p-6">
            <div className="max-w-7xl mx-auto">
              {children}
            </div>
          </main>
        </div>
      </div>
    </div>
  );
};

/**
 * Get page title based on current path
 */
function getPageTitle(path: string): string {
  switch (path) {
    case '/admin/dashboard':
      return 'AI Inquiry Analysis Dashboard';
    case '/admin/inquiries':
      return 'Inquiries';
    case '/admin/metrics':
      return 'Metrics';
    case '/admin/email-status':
      return 'Email Delivery';
    default:
      return 'Admin Portal';
  }
}

export default V0AdminLayout;