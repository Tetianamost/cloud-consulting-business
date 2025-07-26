import React from 'react';

// Simple component to test if Tailwind CSS is working
const TailwindTest: React.FC = () => {
  return (
    <div className="admin-layout">
      <div className="p-4 bg-blue-500 text-white rounded-lg shadow-lg max-w-md mx-auto mt-8">
        <h2 className="text-xl font-bold mb-2">Tailwind CSS Test</h2>
        <p className="text-blue-100">
          If you can see this styled properly with blue background, white text, 
          padding, rounded corners, and shadow, then Tailwind CSS is working correctly!
        </p>
        <div className="mt-4 flex gap-2">
          <button className="admin-button-primary">Primary Button</button>
          <button className="admin-button-secondary">Secondary Button</button>
        </div>
        <div className="mt-4">
          <div className="admin-badge admin-badge-success">Success Badge</div>
          <div className="admin-badge admin-badge-warning ml-2">Warning Badge</div>
          <div className="admin-badge admin-badge-danger ml-2">Danger Badge</div>
        </div>
        <div className="mt-4">
          <div className="admin-progress-bar">
            <div className="admin-progress-fill success" style={{ width: '75%' }}></div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default TailwindTest;