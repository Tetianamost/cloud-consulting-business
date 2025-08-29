import React, { useEffect } from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { useSelector, useDispatch } from 'react-redux';
import { AdminSidebar } from './sidebar';
import { InquiryList } from './inquiry-list';
import { MetricsDashboard } from './metrics-dashboard';
import { EmailMonitor } from './email-monitor';
import { InquiryAnalysisDashboard } from './inquiry-analysis-dashboard';
import ChatToggle from './ChatToggle';
import SimpleChat from './SimpleChat';
import ChatPage from './ChatPage';
import ConnectionStatus from './ConnectionStatus';
import SimpleWebSocketTest from './SimpleWebSocketTest';
import ChatModeToggle from './ChatModeToggle';
import { SimpleWorkingChat } from './SimpleWorkingChat';
import { RootState } from '../../store';
import { chatModeManager } from '../../services/chatModeManager';

interface IntegratedAdminDashboardProps {
  children?: React.ReactNode;
}

export const IntegratedAdminDashboard: React.FC<IntegratedAdminDashboardProps> = ({ children }) => {
  const dispatch = useDispatch();
  const { currentSession, messages } = useSelector((state: RootState) => state.chat);
  const { status: connectionStatus } = useSelector((state: RootState) => state.connection);

  // Initialize chat service when dashboard loads (single connection point)
  useEffect(() => {
    let isComponentMounted = true;

    const initializeConnection = async () => {
      if (!isComponentMounted) {
        return;
      }
      
      // Only initialize if we're not already connected or connecting
      if (connectionStatus === 'disconnected') {
        try {
          console.log('[Dashboard] Initializing chat service via mode manager');
          await chatModeManager.initializeChatService();
        } catch (error) {
          console.error('Failed to initialize chat service:', error);
        }
      }
    };

    // Initialize connection immediately
    initializeConnection();

    // Cleanup on unmount - but be careful with React StrictMode
    return () => {
      console.log('[Dashboard] Component cleanup called - STACK TRACE:', {
        stack: new Error().stack?.split('\n').slice(1, 6).join('\n')
      });
      isComponentMounted = false;
      
      // In development mode, DO NOT cleanup chat service on component unmount
      // This prevents React StrictMode from breaking the connection
      if (process.env.NODE_ENV === 'production') {
        console.log('[Dashboard] Production mode - cleaning up chat service on unmount');
        chatModeManager.cleanup();
      } else {
        console.log('[Dashboard] Development mode - keeping chat service alive to prevent StrictMode issues');
      }
    };
  }, []); // Only run once on mount

  const handleReconnect = () => {
    chatModeManager.forceReconnect();
  };

  return (
    <div className="flex min-h-screen bg-gray-100">
      <AdminSidebar />
      <main className="flex-1 flex flex-col overflow-hidden">
        {/* Dashboard Header with Chat Status */}
        <div className="bg-white border-b border-gray-200 px-6 py-3 flex items-center justify-between flex-shrink-0">
          <div className="flex items-center space-x-4">
            <h1 className="text-lg font-semibold text-gray-900">Admin Dashboard</h1>
            {currentSession && (
              <div className="flex items-center space-x-2 text-sm text-gray-600">
                <span>•</span>
                <span>Active Chat Session</span>
                <span className="px-2 py-1 bg-blue-100 text-blue-800 rounded-full text-xs font-medium">
                  {currentSession.client_name || `Session ${currentSession.id.slice(-8)}`}
                </span>
              </div>
            )}
          </div>
          
          <div className="flex items-center space-x-4">
            <ConnectionStatus 
              showDetails={false}
              onReconnect={handleReconnect}
              className="hidden sm:flex"
            />
            
            {/* Message count indicator */}
            {messages.length > 0 && (
              <div className="flex items-center space-x-2 text-sm text-gray-600">
                <span>{messages.length} messages</span>
              </div>
            )}
          </div>
        </div>

        {/* Main Content */}
        <div className="flex-1 p-6 overflow-y-auto">
          {children || (
            <Routes>
              <Route index element={<Navigate to="dashboard" replace />} />
              <Route path="dashboard" element={<InquiryAnalysisDashboard />} />
              <Route path="inquiries" element={<InquiryList />} />
              <Route path="chat" element={<ChatPage />} />
              <Route path="chat-mode" element={<ChatModeToggle />} />
              <Route path="metrics" element={<MetricsDashboard />} />
              <Route path="email-status" element={<EmailMonitor />} />
              <Route path="websocket-test" element={<SimpleWebSocketTest />} />
              <Route path="simple-chat" element={<SimpleWorkingChat />} />
            </Routes>
          )}
        </div>
      </main>
      
      {/* Consultant Chat Toggle - Always visible */}
      <ChatToggle />

      {/* Simple Working Chat Demo */}
      <div className="mt-8 bg-white rounded-lg shadow-lg">
        <div className="p-6">
          <h2 className="text-xl font-semibold mb-4 text-green-600 flex items-center">
            ✅ Working Chat Demo
            <span className="ml-2 text-sm text-gray-500 font-normal">(No complex polling - just works!)</span>
          </h2>
          <div className="h-96">
            <SimpleChat />
          </div>
        </div>
      </div>
    </div>
  );
};

export default IntegratedAdminDashboard;