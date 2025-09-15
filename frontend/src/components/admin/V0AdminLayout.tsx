import React, { useState } from 'react';
import { Menu, X } from 'lucide-react';
import V0Sidebar from './V0Sidebar';
import V0ErrorBoundary from './V0ErrorBoundary';
import { V0TailwindErrorFallback } from './V0ErrorFallbacks';
import { useV0TailwindErrorHandler } from './useV0ErrorHandler';
import ChatToggle from './ChatToggle';

interface V0AdminLayoutProps {
  children: React.ReactNode;
  currentPath: string;
}

/**
 * V0AdminLayout - Main layout component for admin dashboard
 * Matches the v0.dev design with proper Tailwind styling and enhanced mobile responsiveness
 */
const V0AdminLayout: React.FC<V0AdminLayoutProps> = ({ children, currentPath }) => {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const { tailwindLoaded, retryTailwind } = useV0TailwindErrorHandler();

  // Close sidebar when clicking outside on mobile
  const closeSidebar = () => setSidebarOpen(false);

  // Handle escape key to close sidebar
  React.useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && sidebarOpen) {
        closeSidebar();
      }
    };

    document.addEventListener('keydown', handleEscape);
    return () => document.removeEventListener('keydown', handleEscape);
  }, [sidebarOpen]);

  // Show Tailwind error fallback if CSS failed to load
  if (!tailwindLoaded) {
    return <V0TailwindErrorFallback onRetry={retryTailwind} />;
  }

  return (
    <V0ErrorBoundary>
      <div className="admin-layout min-h-screen bg-gray-50">
      {/* Mobile sidebar overlay with improved animations and accessibility */}
      {sidebarOpen && (
        <div className="fixed inset-0 z-50 lg:hidden">
          {/* Backdrop with fade animation */}
          <div 
            className="fixed inset-0 bg-gray-600 bg-opacity-75 transition-opacity duration-300 ease-in-out"
            onClick={closeSidebar}
            aria-hidden="true"
          />
          
          {/* Mobile sidebar with slide animation */}
          <div className="relative flex-1 flex flex-col max-w-xs w-full bg-white shadow-xl transform transition-transform duration-300 ease-in-out">
            {/* Close button with better positioning */}
            <div className="absolute top-0 right-0 -mr-12 pt-2">
              <button
                className="ml-1 flex items-center justify-center h-10 w-10 rounded-full focus:outline-none focus:ring-2 focus:ring-inset focus:ring-white hover:bg-gray-600 hover:bg-opacity-20 transition-colors"
                onClick={closeSidebar}
                aria-label="Close sidebar"
              >
                <X className="h-6 w-6 text-white" />
              </button>
            </div>
            
            {/* Mobile sidebar content */}
            <V0ErrorBoundary>
              <V0Sidebar currentPath={currentPath} isMobile={true} onNavigate={closeSidebar} />
            </V0ErrorBoundary>
          </div>
        </div>
      )}

      {/* Main layout container with improved responsive behavior */}
      <div className="flex h-screen overflow-hidden">
        {/* Desktop Sidebar - hidden on mobile, visible on lg+ */}
        <V0ErrorBoundary>
          <V0Sidebar currentPath={currentPath} isMobile={false} />
        </V0ErrorBoundary>
        
        {/* Main content area with responsive padding */}
        <div className="flex-1 flex flex-col overflow-hidden min-w-0">
          {/* Header with responsive design */}
          <header className="bg-white border-b border-gray-200 px-4 sm:px-6 py-3 sm:py-4 flex-shrink-0">
            <div className="flex items-center justify-between">
              <div className="flex items-center min-w-0 flex-1">
                {/* Mobile menu button */}
                <button
                  className="lg:hidden mr-3 sm:mr-4 text-gray-500 hover:text-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 rounded-md p-1 transition-colors"
                  onClick={() => setSidebarOpen(true)}
                  aria-label="Open sidebar"
                >
                  <Menu className="h-5 w-5 sm:h-6 sm:w-6" />
                </button>
                
                {/* Page title with responsive text size */}
                <h1 className="text-lg sm:text-xl lg:text-2xl font-semibold text-gray-900 truncate">
                  {getPageTitle(currentPath)}
                </h1>
              </div>
              
              {/* User profile section with responsive design */}
              <div className="flex items-center space-x-2 sm:space-x-4 flex-shrink-0">
                <div className="hidden sm:flex items-center space-x-2">
                  <div className="w-7 h-7 sm:w-8 sm:h-8 bg-blue-500 rounded-full flex items-center justify-center">
                    <span className="text-white text-xs sm:text-sm font-medium">A</span>
                  </div>
                  <span className="text-sm text-gray-700 hidden md:inline">Admin User</span>
                </div>
                
                {/* Mobile-only user avatar */}
                <div className="sm:hidden w-7 h-7 bg-blue-500 rounded-full flex items-center justify-center">
                  <span className="text-white text-xs font-medium">A</span>
                </div>
              </div>
            </div>
          </header>

          {/* Main content with responsive padding and scrolling */}
          <main className={`flex-1 overflow-y-auto bg-gray-50 ${
            currentPath === '/admin/chat' 
              ? 'p-0' // No padding for chat page to allow full height
              : 'px-4 py-4 sm:px-6 sm:py-6'
          }`}>
            <div className={currentPath === '/admin/chat' ? 'h-full' : 'max-w-7xl mx-auto'}>
              <V0ErrorBoundary>
                {children}
              </V0ErrorBoundary>
            </div>
          </main>
        </div>
      </div>
      
      {/* Consultant Chat Toggle - only show if not on chat page */}
      {currentPath !== '/admin/chat' && <ChatToggle />}
    </div>
    </V0ErrorBoundary>
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
    case '/admin/chat':
      return 'AI Consultant Chat';
    case '/admin/metrics':
      return 'Metrics';
    case '/admin/email-status':
      return 'Email Delivery';
    default:
      return 'Admin Portal';
  }
}

export default V0AdminLayout;