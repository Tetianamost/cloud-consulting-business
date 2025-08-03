import React, { useState, useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { MessageSquare, Bell, AlertCircle } from 'lucide-react';
import { RootState } from '../../store';
import ConsultantChat from './ConsultantChat';
import websocketService from '../../services/websocketService';

interface ChatNotification {
  id: string;
  type: 'message' | 'error' | 'connection';
  title: string;
  message: string;
  timestamp: string;
  read: boolean;
}

export const ChatToggle: React.FC = () => {
  const dispatch = useDispatch();
  
  // Redux state
  const { currentSession, messages, isLoading, error } = useSelector((state: RootState) => state.chat);
  const { status: connectionStatus, isHealthy } = useSelector((state: RootState) => state.connection);
  
  // Local state for widget
  const [isChatOpen, setIsChatOpen] = useState(() => {
    // Persist chat state across page navigation
    const saved = localStorage.getItem('chatWidgetState');
    return saved ? JSON.parse(saved).isOpen : false;
  });
  
  const [isMinimized, setIsMinimized] = useState(() => {
    const saved = localStorage.getItem('chatWidgetState');
    return saved ? JSON.parse(saved).isMinimized : false;
  });
  
  const [notifications, setNotifications] = useState<ChatNotification[]>([]);
  const [showNotifications, setShowNotifications] = useState(false);

  // Persist widget state
  useEffect(() => {
    const state = { isOpen: isChatOpen, isMinimized };
    localStorage.setItem('chatWidgetState', JSON.stringify(state));
  }, [isChatOpen, isMinimized]);

  // Initialize WebSocket connection when component mounts
  useEffect(() => {
    if (connectionStatus === 'disconnected') {
      const connectPromise = websocketService.connect();
      if (connectPromise && typeof connectPromise.catch === 'function') {
        connectPromise.catch(console.error);
      }
    }
  }, [connectionStatus]);

  // Handle new messages for notifications
  useEffect(() => {
    if (messages.length > 0 && !isChatOpen) {
      const lastMessage = messages[messages.length - 1];
      if (lastMessage.type === 'assistant') {
        const notification: ChatNotification = {
          id: `msg-${Date.now()}`,
          type: 'message',
          title: 'New AI Response',
          message: lastMessage.content.slice(0, 100) + (lastMessage.content.length > 100 ? '...' : ''),
          timestamp: new Date().toISOString(),
          read: false,
        };
        
        setNotifications(prev => [notification, ...prev.slice(0, 4)]); // Keep only 5 notifications
      }
    }
  }, [messages, isChatOpen]);

  // Handle connection status notifications
  useEffect(() => {
    if (connectionStatus === 'failed') {
      const notification: ChatNotification = {
        id: `conn-${Date.now()}`,
        type: 'connection',
        title: 'Connection Failed',
        message: 'Unable to connect to chat service. Please try again.',
        timestamp: new Date().toISOString(),
        read: false,
      };
      
      setNotifications(prev => [notification, ...prev.slice(0, 4)]);
    }
  }, [connectionStatus]);

  // Handle errors
  useEffect(() => {
    if (error) {
      const notification: ChatNotification = {
        id: `error-${Date.now()}`,
        type: 'error',
        title: 'Chat Error',
        message: error,
        timestamp: new Date().toISOString(),
        read: false,
      };
      
      setNotifications(prev => [notification, ...prev.slice(0, 4)]);
    }
  }, [error]);

  const handleToggleChat = () => {
    if (isChatOpen && !isMinimized) {
      setIsMinimized(true);
    } else if (isChatOpen && isMinimized) {
      setIsMinimized(false);
      // Mark notifications as read when expanding
      setNotifications(prev => prev.map(n => ({ ...n, read: true })));
    } else {
      setIsChatOpen(true);
      setIsMinimized(false);
      // Mark notifications as read when opening
      setNotifications(prev => prev.map(n => ({ ...n, read: true })));
    }
  };

  const handleCloseChat = () => {
    setIsChatOpen(false);
    setIsMinimized(false);
  };

  const unreadCount = notifications.filter(n => !n.read).length;
  const hasUnreadMessages = unreadCount > 0;

  // Get connection status indicator
  const getConnectionIndicator = () => {
    switch (connectionStatus) {
      case 'connected':
        return isHealthy ? 'bg-green-500' : 'bg-yellow-500';
      case 'connecting':
      case 'reconnecting':
        return 'bg-yellow-500 animate-pulse';
      case 'failed':
      case 'disconnected':
        return 'bg-red-500';
      default:
        return 'bg-gray-500';
    }
  };

  if (!isChatOpen) {
    return (
      <div className="fixed bottom-4 right-4 z-50">
        {/* Notifications dropdown */}
        {showNotifications && notifications.length > 0 && (
          <div className="absolute bottom-16 right-0 w-80 bg-white rounded-lg shadow-xl border border-gray-200 mb-2">
            <div className="p-3 border-b border-gray-200">
              <div className="flex items-center justify-between">
                <h3 className="text-sm font-medium text-gray-900">Chat Notifications</h3>
                <button
                  onClick={() => setNotifications([])}
                  className="text-xs text-gray-500 hover:text-gray-700"
                >
                  Clear all
                </button>
              </div>
            </div>
            <div className="max-h-64 overflow-y-auto">
              {notifications.map((notification) => (
                <div
                  key={notification.id}
                  className={`p-3 border-b border-gray-100 hover:bg-gray-50 cursor-pointer ${
                    !notification.read ? 'bg-blue-50' : ''
                  }`}
                  onClick={() => {
                    setIsChatOpen(true);
                    setIsMinimized(false);
                    setShowNotifications(false);
                    setNotifications(prev => prev.map(n => 
                      n.id === notification.id ? { ...n, read: true } : n
                    ));
                  }}
                >
                  <div className="flex items-start space-x-2">
                    <div className="flex-shrink-0 mt-0.5">
                      {notification.type === 'message' && (
                        <MessageSquare className="h-4 w-4 text-blue-600" />
                      )}
                      {notification.type === 'error' && (
                        <AlertCircle className="h-4 w-4 text-red-600" />
                      )}
                      {notification.type === 'connection' && (
                        <AlertCircle className="h-4 w-4 text-yellow-600" />
                      )}
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className="text-sm font-medium text-gray-900 truncate">
                        {notification.title}
                      </p>
                      <p className="text-xs text-gray-600 mt-1 line-clamp-2">
                        {notification.message}
                      </p>
                      <p className="text-xs text-gray-500 mt-1">
                        {new Date(notification.timestamp).toLocaleTimeString([], {
                          hour: '2-digit',
                          minute: '2-digit'
                        })}
                      </p>
                    </div>
                    {!notification.read && (
                      <div className="w-2 h-2 bg-blue-600 rounded-full flex-shrink-0 mt-2" />
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Chat toggle button */}
        <div className="relative">
          <button
            onClick={handleToggleChat}
            className="bg-blue-600 hover:bg-blue-700 text-white p-3 rounded-full shadow-lg transition-colors group"
            title={`Open Consultant Assistant${currentSession ? ` (Session: ${currentSession.id.slice(-8)})` : ''}`}
            aria-label="Open AI consultant chat"
          >
            <MessageSquare className="h-6 w-6" />
            
            {/* Connection status indicator */}
            <div 
              className={`absolute -top-1 -right-1 w-3 h-3 rounded-full border-2 border-white ${getConnectionIndicator()}`}
              title={`Connection: ${connectionStatus}`}
            />
            
            {/* Unread message indicator */}
            {hasUnreadMessages && (
              <div className="absolute -top-2 -left-2 bg-red-500 text-white text-xs rounded-full h-5 w-5 flex items-center justify-center font-medium">
                {unreadCount > 9 ? '9+' : unreadCount}
              </div>
            )}
          </button>

          {/* Notifications bell */}
          {notifications.length > 0 && (
            <button
              onClick={(e) => {
                e.stopPropagation();
                setShowNotifications(!showNotifications);
              }}
              className="absolute -top-2 -left-12 bg-gray-600 hover:bg-gray-700 text-white p-2 rounded-full shadow-lg transition-colors"
              title="View notifications"
              aria-label="View chat notifications"
            >
              <Bell className="h-4 w-4" />
              {hasUnreadMessages && (
                <div className="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full h-4 w-4 flex items-center justify-center font-medium">
                  {unreadCount > 9 ? '9+' : unreadCount}
                </div>
              )}
            </button>
          )}
        </div>
      </div>
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