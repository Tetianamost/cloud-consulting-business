import React, { useState } from 'react';
import ConsultantChat from './ConsultantChat';

export const ChatToggle: React.FC = () => {
  const [isChatOpen, setIsChatOpen] = useState(false);
  const [isMinimized, setIsMinimized] = useState(false);

  const handleToggleChat = () => {
    if (isChatOpen && !isMinimized) {
      setIsMinimized(true);
    } else if (isChatOpen && isMinimized) {
      setIsMinimized(false);
    } else {
      setIsChatOpen(true);
      setIsMinimized(false);
    }
  };

  const handleCloseChat = () => {
    setIsChatOpen(false);
    setIsMinimized(false);
  };

  if (!isChatOpen) {
    return (
      <button
        onClick={handleToggleChat}
        className="fixed bottom-4 right-4 bg-blue-600 hover:bg-blue-700 text-white p-3 rounded-full shadow-lg transition-colors z-50"
        title="Open Consultant Assistant"
      >
        <svg className="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
        </svg>
      </button>
    );
  }

  return (
    <ConsultantChat
      isMinimized={isMinimized}
      onToggleMinimize={handleToggleChat}
      onClose={handleCloseChat}
    />
  );
};

export default ChatToggle;