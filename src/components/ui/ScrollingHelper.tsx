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
 * functionality during normal user scrolling
 * 
 * This prevents the page from getting stuck when scrolling quickly
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
  ...props 
}) => {
  const isUserScrolling = useScrollDetection();
  
  // Only enable spy when user isn't actively scrolling
  const spyEnabled = !isUserScrolling && spy;
  
  return (
    <ScrollLink {...props} spy={spyEnabled}>
      {children}
    </ScrollLink>
  );
};

/**
 * Global scroll management - used to register global scroll listeners
 * Custom React Hook for App component
 */
export const useScrollManager = () => {
  useEffect(() => {
    // Set up event handlers for react-scroll issues
    const handleScrollEvent = () => {
      // This helps prevent scroll jank and "sticky scrolling"
      document.body.style.overflow = 'auto';
    };
    
    // Register event listeners
    Events.scrollEvent.register('begin', () => {
      // Handle scroll start
    });
    
    Events.scrollEvent.register('end', handleScrollEvent);
    
    return () => {
      // Clean up all registered events when unmounting
      Events.scrollEvent.remove('begin');
      Events.scrollEvent.remove('end');
    };
  }, []);
  
  // Helper methods for smooth scrolling
  const scrollToTop = () => {
    scroller.scrollTo('top', {
      duration: 500,
      smooth: true
    });
  };
  
  return { scrollToTop };
};

export default SmartScrollLink;