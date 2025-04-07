import React, { useState, useEffect } from 'react';
import { Link as ScrollLink, Events, scroller } from 'react-scroll';

/**
 * Custom hook to detect user scrolling actions
 * Used to temporarily disable react-scroll's spy feature during user scrolling
 * to prevent the page from getting stuck at elements with active links
 */
export const useScrollDetection = () => {
  const [isUserScrolling, setIsUserScrolling] = useState(false);
  const [scrollTimeout, setScrollTimeout] = useState<NodeJS.Timeout | null>(null);
  
  useEffect(() => {
    // Function to handle scrolling events
    const handleScrollStart = () => {
      setIsUserScrolling(true);
      
      // Clear any existing timeout
      if (scrollTimeout) {
        clearTimeout(scrollTimeout);
      }
      
      // Set a new timeout to re-enable spy after scrolling stops
      const timeout = setTimeout(() => {
        setIsUserScrolling(false);
      }, 500); // Wait 500ms after scrolling stops to re-enable spy
      
      setScrollTimeout(timeout);
    };
    
    // Add event listeners for scroll detection
    window.addEventListener('wheel', handleScrollStart, { passive: true });
    window.addEventListener('touchmove', handleScrollStart, { passive: true });
    
    // Handle keyboard scrolling
    window.addEventListener('keydown', (e) => {
      // Detect arrow keys, Page Up/Down, Space, Home, End
      if ([32, 33, 34, 35, 36, 37, 38, 39, 40].includes(e.keyCode)) {
        handleScrollStart();
      }
    });
    
    // Cleanup event listeners
    return () => {
      window.removeEventListener('wheel', handleScrollStart);
      window.removeEventListener('touchmove', handleScrollStart);
      
      if (scrollTimeout) clearTimeout(scrollTimeout);
    };
  }, [scrollTimeout]);
  
  return isUserScrolling;
};

/**
 * Improved ScrollLink component that temporarily disables the spy 
 * functionality during normal user scrolling and specifically addresses
 * mobile touch scrolling issues
 * 
 * This prevents the page from getting stuck when scrolling quickly and
 * fixes the "first scroll gets stuck" issue on mobile devices
 */
interface LinkProps {
  to: string;
  spy?: boolean;
  smooth?: boolean;
  offset?: number;
  duration?: number;
  children: React.ReactNode;
  className?: string;
  activeClass?: string;
  onClick?: () => void;
  [key: string]: any;
}

export const SmartScrollLink: React.FC<LinkProps> = ({ 
  children, 
  spy = true,
  onClick,
  ...props 
}) => {
  const isUserScrolling = useScrollDetection();
  const [isMobile, setIsMobile] = useState(false);
  
  // Detect mobile devices
  useEffect(() => {
    const checkMobile = () => {
      const mobile = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent) || 
                    (window.innerWidth <= 768);
      setIsMobile(mobile);
    };
    
    checkMobile();
    window.addEventListener('resize', checkMobile);
    
    return () => {
      window.removeEventListener('resize', checkMobile);
    };
  }, []);
  
  // Only enable spy when user isn't actively scrolling
  const spyEnabled = !isUserScrolling && spy;
  
  // Custom click handler to fix mobile scrolling issues
  const handleClick = (e: React.MouseEvent) => {
    // Execute original onClick if provided
    if (onClick) {
      onClick();
    }
    
    // Mobile-specific fixes
    if (isMobile) {
      // Add a tiny delay before allowing scrolling again
      // This prevents the "first scroll gets stuck" issue on mobile
      setTimeout(() => {
        // Force enable scrolling after navigation completes
        document.body.style.overflow = 'auto';
        document.documentElement.style.overflow = 'auto';
        
        // Remove any stuck touch events
        window.scrollBy(0, 1);
        window.scrollBy(0, -1);
      }, props.duration || 500);
    }
  };
  
  return (
    <ScrollLink
      {...props}
      spy={spyEnabled}
      onClick={handleClick}
      // Reduce duration slightly on mobile for better performance
      duration={isMobile ? (props.duration ? props.duration * 0.8 : 400) : props.duration}
    >
      {children}
    </ScrollLink>
  );
};

/**
 * Global scroll management - used to register global scroll listeners
 * Custom React Hook for App component with special mobile fixes
 */
export const useScrollManager = () => {
  // Track if device is mobile
  const [isMobile, setIsMobile] = useState(false);
  
  useEffect(() => {
    // Mobile detection function
    const checkMobile = () => {
      const mobile = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent) || 
                    (window.innerWidth <= 768);
      setIsMobile(mobile);
    };
    
    checkMobile();
    window.addEventListener('resize', checkMobile);
    
    // Set up event handlers for react-scroll issues
    const handleScrollEvent = () => {
      // This helps prevent scroll jank and "sticky scrolling"
      document.body.style.overflow = 'auto';
      document.documentElement.style.overflow = 'auto';
      
      // Mobile-specific fixes
      if (isMobile) {
        // Force a tiny scroll to reset mobile touch events
        // This is critical to fix the first-scroll-stuck behavior
        setTimeout(() => {
          window.scrollBy(0, 1);
          window.scrollBy(0, -1);
        }, 50);
      }
    };
    
    // Special handler to fix iOS momentum scrolling issues
    const fixMobileScrolling = () => {
      if (isMobile) {
        // Prevent body from becoming non-scrollable
        document.body.style.overflow = 'auto';
        document.documentElement.style.overflow = 'auto';
        
        // Enable -webkit-overflow-scrolling on iOS
        // Use type assertion to make TypeScript allow this vendor prefixed property
        (document.body.style as any)['-webkit-overflow-scrolling'] = 'touch';
      }
    };
    
    // Run the fix immediately
    fixMobileScrolling();
    
    // Register event listeners for React Scroll
    Events.scrollEvent.register('begin', fixMobileScrolling);
    Events.scrollEvent.register('end', handleScrollEvent);
    
    // Mobile-specific event listeners
    if (isMobile) {
      // Listen for the end of touch events
      document.addEventListener('touchend', () => {
        setTimeout(fixMobileScrolling, 100);
      });
    }
    
    return () => {
      // Clean up all registered events when unmounting
      Events.scrollEvent.remove('begin');
      Events.scrollEvent.remove('end');
      window.removeEventListener('resize', checkMobile);
      document.removeEventListener('touchend', fixMobileScrolling);
    };
  }, [isMobile]);
  
  // Helper methods for smooth scrolling
  const scrollToTop = () => {
    scroller.scrollTo('top', {
      duration: isMobile ? 400 : 500, // Shorter duration on mobile
      smooth: true
    });
  };
  
  return { scrollToTop };
};

export default SmartScrollLink;