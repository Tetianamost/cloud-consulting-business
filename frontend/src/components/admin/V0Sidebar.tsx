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
 * V0Sidebar - Navigation sidebar component matching v0.dev design
 */
const V0Sidebar: React.FC<V0SidebarProps> = ({ currentPath, onNavigate }) => {
  const location = useLocation();
  const activePath = currentPath || location.pathname;

  return (
    <div 
      className="hidden lg:flex lg:w-64 bg-white border-r border-gray-200 flex-col h-full"
      data-testid="v0-sidebar"
      style={{ 
        /* Debugging: Force visibility for testing */
        // display: 'flex' 
      }}
    >
      {/* Logo/Header Section */}
      <div className="flex items-center px-6 py-4 border-b border-gray-200">
        <div className="flex items-center space-x-2">
          <div className="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center">
            <Bot className="w-5 h-5 text-white" />
          </div>
          <div>
            <h1 className="text-lg font-semibold text-gray-900">AI Admin Portal</h1>
          </div>
        </div>
      </div>

      {/* Navigation Items */}
      <nav className="flex-1 px-4 py-6 space-y-2">
        {navItems.map((item) => {
          const isActive = activePath === item.href;
          const Icon = item.icon;
          
          return (
            <Link
              key={item.href}
              to={item.href}
              onClick={() => onNavigate?.(item.href)}
              className={`
                flex items-center px-3 py-2 rounded-lg text-sm font-medium transition-colors
                ${isActive 
                  ? 'bg-blue-50 text-blue-700 border-r-2 border-blue-700' 
                  : 'text-gray-700 hover:bg-gray-100 hover:text-gray-900'
                }
              `}
            >
              <Icon className={`w-5 h-5 mr-3 ${isActive ? 'text-blue-700' : 'text-gray-400'}`} />
              <div className="flex-1">
                <div className="font-medium">{item.title}</div>
                {item.description && (
                  <div className={`text-xs ${isActive ? 'text-blue-600' : 'text-gray-500'}`}>
                    {item.description}
                  </div>
                )}
              </div>
            </Link>
          );
        })}
      </nav>

      {/* User Profile Section at Bottom */}
      <div className="border-t border-gray-200 p-4">
        <div className="flex items-center space-x-3">
          <div className="w-10 h-10 bg-gray-300 rounded-full flex items-center justify-center">
            <span className="text-gray-600 text-sm font-medium">AU</span>
          </div>
          <div className="flex-1 min-w-0">
            <p className="text-sm font-medium text-gray-900 truncate">
              Admin User
            </p>
            <p className="text-xs text-gray-500 truncate">
              admin@example.com
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default V0Sidebar;