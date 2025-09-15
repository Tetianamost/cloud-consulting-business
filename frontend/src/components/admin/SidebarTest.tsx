import React from 'react';

/**
 * Simple test component to verify Tailwind responsive classes are working
 */
const SidebarTest: React.FC = () => {
  return (
    <div className="p-4">
      <h2 className="text-xl font-bold mb-4">Sidebar Responsive Test</h2>
      
      {/* Test the same responsive classes used in V0Sidebar */}
      <div className="hidden lg:flex lg:w-64 bg-blue-100 border border-blue-300 p-4 mb-4">
        <p>This should be hidden on mobile and visible on desktop (lg and up)</p>
      </div>
      
      {/* Test basic responsive visibility */}
      <div className="block lg:hidden bg-red-100 border border-red-300 p-4 mb-4">
        <p>This should be visible on mobile and hidden on desktop</p>
      </div>
      
      {/* Test flex responsive classes */}
      <div className="flex flex-col lg:flex-row gap-4">
        <div className="bg-green-100 border border-green-300 p-4 flex-1">
          <p>Column on mobile, row on desktop</p>
        </div>
        <div className="bg-yellow-100 border border-yellow-300 p-4 flex-1">
          <p>Column on mobile, row on desktop</p>
        </div>
      </div>
    </div>
  );
};

export default SidebarTest;