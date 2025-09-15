import React, { useState, useEffect } from 'react';
import { Menu, X, Smartphone, Tablet, Monitor } from 'lucide-react';

/**
 * ResponsiveTest - Component to test responsive behavior of admin components
 * This component helps verify that all breakpoints work correctly
 */
const ResponsiveTest: React.FC = () => {
  const [screenSize, setScreenSize] = useState<'mobile' | 'tablet' | 'desktop'>('desktop');
  const [windowWidth, setWindowWidth] = useState(0);

  useEffect(() => {
    const handleResize = () => {
      const width = window.innerWidth;
      setWindowWidth(width);
      
      if (width < 768) {
        setScreenSize('mobile');
      } else if (width < 1024) {
        setScreenSize('tablet');
      } else {
        setScreenSize('desktop');
      }
    };

    handleResize();
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  const getScreenIcon = () => {
    switch (screenSize) {
      case 'mobile':
        return <Smartphone className="h-4 w-4" />;
      case 'tablet':
        return <Tablet className="h-4 w-4" />;
      default:
        return <Monitor className="h-4 w-4" />;
    }
  };

  const getBreakpointInfo = () => {
    return {
      mobile: windowWidth < 640 ? '✓' : '✗',
      sm: windowWidth >= 640 && windowWidth < 768 ? '✓' : '✗',
      md: windowWidth >= 768 && windowWidth < 1024 ? '✓' : '✗',
      lg: windowWidth >= 1024 && windowWidth < 1280 ? '✓' : '✗',
      xl: windowWidth >= 1280 ? '✓' : '✗',
    };
  };

  const breakpoints = getBreakpointInfo();

  return (
    <div className="fixed bottom-4 right-4 z-50 bg-white border border-gray-200 rounded-lg shadow-lg p-4 max-w-xs">
      <div className="flex items-center space-x-2 mb-3">
        {getScreenIcon()}
        <span className="font-medium text-sm">
          {screenSize.charAt(0).toUpperCase() + screenSize.slice(1)} ({windowWidth}px)
        </span>
      </div>
      
      <div className="space-y-1 text-xs">
        <div className="font-medium text-gray-700 mb-2">Tailwind Breakpoints:</div>
        <div className="grid grid-cols-2 gap-1">
          <div className={`flex justify-between ${breakpoints.mobile === '✓' ? 'text-green-600' : 'text-gray-400'}`}>
            <span>&lt;640px:</span>
            <span>{breakpoints.mobile}</span>
          </div>
          <div className={`flex justify-between ${breakpoints.sm === '✓' ? 'text-green-600' : 'text-gray-400'}`}>
            <span>sm:</span>
            <span>{breakpoints.sm}</span>
          </div>
          <div className={`flex justify-between ${breakpoints.md === '✓' ? 'text-green-600' : 'text-gray-400'}`}>
            <span>md:</span>
            <span>{breakpoints.md}</span>
          </div>
          <div className={`flex justify-between ${breakpoints.lg === '✓' ? 'text-green-600' : 'text-gray-400'}`}>
            <span>lg:</span>
            <span>{breakpoints.lg}</span>
          </div>
          <div className={`flex justify-between ${breakpoints.xl === '✓' ? 'text-green-600' : 'text-gray-400'}`}>
            <span>xl:</span>
            <span>{breakpoints.xl}</span>
          </div>
        </div>
      </div>

      <div className="mt-3 pt-3 border-t border-gray-100">
        <div className="text-xs text-gray-600">
          <div className="font-medium mb-1">Expected Behavior:</div>
          <ul className="space-y-1">
            <li>• Sidebar: Hidden &lt;lg, visible ≥lg</li>
            <li>• Metrics: 1 col &lt;sm, 2 cols sm-lg, 4 cols ≥lg</li>
            <li>• Tables: Horizontal scroll on mobile</li>
            <li>• Touch: 44px+ tap targets</li>
          </ul>
        </div>
      </div>
    </div>
  );
};

export default ResponsiveTest;