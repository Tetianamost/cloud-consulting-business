import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { 
  BarChart3, 
  MessageSquare, 
  FileText, 
  Mail, 
  Settings,
  Bot
} from 'lucide-react';

interface V0SidebarProps {
  currentPath: string;
  onNavigate?: (path: string) => void;
  isMobile?: boolean;
}

interface NavItem {
  title: string;
  href: string;
  icon: React.ComponentType<{ className?: string }>;
  description?: string;
}

const navItems: NavItem[] = [
  {
    title: "AI Dashboard",
    href: "/admin/dashboard",
    icon: BarChart3,
    description: "Overview and analytics"
  },
  {
    title: "Inquiries",
    href: "/admin/inquiries",
    icon: MessageSquare,
    description: "Customer inquiries"
  },
  {
    title: "AI Consultant",
    href: "/admin/ai-consultant",
    icon: Bot,
    description: "Full AI Assistant"
  },
  {
    title: "AI Reports",
    href: "/admin/reports",
    icon: FileText,
    description: "Generated reports"
  },
  {
    title: "Email Monitor",
    href: "/admin/email-status",
    icon: Mail,
    description: "Email delivery status"
  },
  {
    title: "Settings",
    href: "/admin/settings",
    icon: Settings,
    description: "System settings"
  }
];

/**
 * V0Sidebar - Navigation sidebar component matching v0.dev design with enhanced mobile support
 * Memoized for performance optimization
 */
const V0Sidebar: React.FC<V0SidebarProps> = React.memo(({ currentPath, onNavigate, isMobile = false }) => {
  const location = useLocation();
  const activePath = currentPath || location.pathname;

  // Handle navigation with optional callback for mobile
  const handleNavigation = (path: string) => {
    if (onNavigate) {
      onNavigate(path);
    }
  };

  return (
    <div 
      className={`
        ${isMobile 
          ? 'flex w-full' 
          : 'hidden lg:flex lg:w-64'
        } 
        bg-white border-r border-gray-200 flex-col h-full transition-all duration-200
      `}
      data-testid="v0-sidebar"
    >
      {/* Logo/Header Section with responsive padding */}
      <div className={`flex items-center border-b border-gray-200 ${isMobile ? 'px-4 py-3' : 'px-6 py-4'}`}>
        <div className="flex items-center space-x-2">
          <div className={`bg-blue-600 rounded-lg flex items-center justify-center ${isMobile ? 'w-7 h-7' : 'w-8 h-8'}`}>
            <Bot className={`text-white ${isMobile ? 'w-4 h-4' : 'w-5 h-5'}`} />
          </div>
          <div>
            <h1 className={`font-semibold text-gray-900 ${isMobile ? 'text-base' : 'text-lg'}`}>
              AI Admin Portal
            </h1>
          </div>
        </div>
      </div>

      {/* Navigation Items with responsive spacing */}
      <nav className={`flex-1 space-y-1 ${isMobile ? 'px-3 py-4' : 'px-4 py-6'}`}>
        {navItems.map((item) => {
          const isActive = activePath === item.href;
          const Icon = item.icon;
          
          return (
            <Link
              key={item.href}
              to={item.href}
              onClick={() => handleNavigation(item.href)}
              className={`
                flex items-center rounded-lg font-medium transition-all duration-200
                ${isMobile ? 'px-3 py-3 text-sm' : 'px-3 py-2 text-sm lg:text-base'}
                ${isActive 
                  ? 'bg-blue-50 text-blue-700 border-r-2 border-blue-700 shadow-sm' 
                  : 'text-gray-700 hover:bg-gray-100 hover:text-gray-900 hover:shadow-sm'
                }
              `}
            >
              <Icon className={`mr-3 transition-colors ${isMobile ? 'w-5 h-5' : 'w-5 h-5'} ${
                isActive ? 'text-blue-700' : 'text-gray-400'
              }`} />
              <div className="flex-1 min-w-0">
                <div className="font-medium truncate">{item.title}</div>
                {item.description && !isMobile && (
                  <div className={`text-xs truncate ${isActive ? 'text-blue-600' : 'text-gray-500'}`}>
                    {item.description}
                  </div>
                )}
              </div>
            </Link>
          );
        })}
      </nav>

      {/* User Profile Section at Bottom with responsive design */}
      <div className={`border-t border-gray-200 ${isMobile ? 'p-3' : 'p-4'}`}>
        <div className="flex items-center space-x-3">
          <div className={`bg-gray-300 rounded-full flex items-center justify-center ${isMobile ? 'w-9 h-9' : 'w-10 h-10'}`}>
            <span className={`text-gray-600 font-medium ${isMobile ? 'text-xs' : 'text-sm'}`}>AU</span>
          </div>
          <div className="flex-1 min-w-0">
            <p className={`font-medium text-gray-900 truncate ${isMobile ? 'text-sm' : 'text-sm'}`}>
              Admin User
            </p>
            {!isMobile && (
              <p className="text-xs text-gray-500 truncate">
                admin@example.com
              </p>
            )}
          </div>
        </div>
      </div>
    </div>
  );
});

export default V0Sidebar;