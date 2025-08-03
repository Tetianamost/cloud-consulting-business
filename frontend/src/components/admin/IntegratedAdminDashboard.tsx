import React, { useEffect } from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { useSelector, useDispatch } from 'react-redux';
import { AdminSidebar } from './sidebar';
import { InquiryList } from './inquiry-list';
import { MetricsDashboard } from './metrics-dashboard';
import { EmailMonitor } from './email-monitor';
import { InquiryAnalysisDashboard } from './inquiry-analysis-dashboard';
import ChatToggle from './ChatToggle';
import ChatPage from './ChatPage';
import ConnectionStatus from './ConnectionStatus';
import { RootState } from '../../store';
import websocketService from '../../services/websocketService';

interface IntegratedAdminDashboardProps {
  children?: React.ReactNode;
}

export const IntegratedAdminDashboard: React.FC<IntegratedAdminDashboardProps> = ({ children }) => {
  const dispatch = useDispatch();
  const { currentSession, messages } = useSelector((state: RootState) => state.chat);
  const { status: connectionStatus } = useSelector((state: RootState) => state.connection);

  // Initialize WebSocket connection when dashboard loads
  useEffect(() => {
    if (connectionStatus === 'disconnected') {
      const connectPromise = websocketService.connect();
      if (connectPromise && typeof connectPromise.catch === 'function') {
        connectPromise.catch(console.error);
      }
    }

    // Cleanup on unmount
    return () => {
      // Don't disconnect here as other components might be using the connection
    };
  }, [connectionStatus]);

  const handleReconnect = () => {
    websocketService.forceReconnect();
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
              <Route path="metrics" element={<MetricsDashboard />} />
              <Route path="email-status" element={<EmailMonitor />} />
            </Routes>
          )}
        </div>
      </main>
      
      {/* Consultant Chat Toggle - Always visible */}
      <ChatToggle />
    </div>
  );
};

export default IntegratedAdminDashboard;