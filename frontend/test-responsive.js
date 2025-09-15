/**
 * Simple script to test responsive design implementation
 * Run this in the browser console to verify breakpoints
 */

function testResponsiveBreakpoints() {
  const breakpoints = {
    mobile: 640,
    sm: 640,
    md: 768,
    lg: 1024,
    xl: 1280
  };

  console.log('=== Responsive Design Test ===');
  console.log(`Current window width: ${window.innerWidth}px`);
  
  // Test sidebar visibility
  const sidebar = document.querySelector('[data-testid="v0-sidebar"]');
  if (sidebar) {
    const sidebarStyles = window.getComputedStyle(sidebar);
    console.log(`Sidebar display: ${sidebarStyles.display}`);
    console.log(`Sidebar visibility expected: ${window.innerWidth >= breakpoints.lg ? 'visible' : 'hidden'}`);
  }

  // Test metrics grid
  const metricsGrid = document.querySelector('.grid.grid-cols-1.sm\\:grid-cols-2.lg\\:grid-cols-4');
  if (metricsGrid) {
    const gridStyles = window.getComputedStyle(metricsGrid);
    console.log(`Metrics grid columns: ${gridStyles.gridTemplateColumns}`);
  }

  // Test mobile menu button
  const mobileMenuButton = document.querySelector('button[aria-label="Open sidebar"]');
  if (mobileMenuButton) {
    const buttonStyles = window.getComputedStyle(mobileMenuButton);
    console.log(`Mobile menu button display: ${buttonStyles.display}`);
    console.log(`Mobile menu expected: ${window.innerWidth < breakpoints.lg ? 'visible' : 'hidden'}`);
  }

  // Test table responsiveness
  const tableWrapper = document.querySelector('.overflow-x-auto');
  if (tableWrapper) {
    console.log('Table wrapper found - horizontal scroll enabled for mobile');
  }

  console.log('=== Test Complete ===');
}

// Test at different breakpoints
function testAllBreakpoints() {
  const testSizes = [375, 640, 768, 1024, 1280];
  
  testSizes.forEach(width => {
    console.log(`\n--- Testing at ${width}px ---`);
    // Note: This would require actually resizing the window
    // In real testing, you'd use browser dev tools or automated testing
    console.log(`Expected behavior at ${width}px:`);
    
    if (width < 640) {
      console.log('- Sidebar: Hidden');
      console.log('- Metrics: 1 column');
      console.log('- Mobile menu: Visible');
    } else if (width < 768) {
      console.log('- Sidebar: Hidden');
      console.log('- Metrics: 2 columns');
      console.log('- Mobile menu: Visible');
    } else if (width < 1024) {
      console.log('- Sidebar: Hidden');
      console.log('- Metrics: 2 columns');
      console.log('- Mobile menu: Visible');
    } else {
      console.log('- Sidebar: Visible');
      console.log('- Metrics: 4 columns');
      console.log('- Mobile menu: Hidden');
    }
  });
}

// Export functions for browser console use
window.testResponsiveBreakpoints = testResponsiveBreakpoints;
window.testAllBreakpoints = testAllBreakpoints;

console.log('Responsive test functions loaded. Run testResponsiveBreakpoints() or testAllBreakpoints() in console.');