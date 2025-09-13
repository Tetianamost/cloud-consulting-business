import React, { useState } from 'react';
import { pollingChatService } from '../../services/pollingChatService';

interface DiagnosticButtonProps {
  className?: string;
}

const DiagnosticButton: React.FC<DiagnosticButtonProps> = ({ className = '' }) => {
  const [isRunning, setIsRunning] = useState(false);
  const [diagnosticResult, setDiagnosticResult] = useState<string | null>(null);

  const handleRunDiagnostics = async () => {
    setIsRunning(true);
    setDiagnosticResult(null);
    
    try {
      // Run polling chat diagnostics
      const connectionStatus = pollingChatService.getConnectionStatus();
      const statusInfo = pollingChatService.getConnectionStatusInfo();
      const isHealthy = pollingChatService.isHealthy();
      const statusMessage = pollingChatService.getStatusMessage();
      
      // Create diagnostic report
      const report = {
        timestamp: new Date().toISOString(),
        connectionStatus,
        isHealthy,
        statusMessage,
        statusInfo
      };
      
      console.log('🔍 Polling Chat Diagnostic Report:', report);
      
      // Set user-friendly result message
      if (isHealthy) {
        setDiagnosticResult('✅ Polling chat is working correctly');
      } else {
        setDiagnosticResult(`⚠️ ${statusMessage}`);
      }
      
      // Force a reconnection attempt if unhealthy
      if (!isHealthy) {
        console.log('🔄 Attempting to reconnect...');
        pollingChatService.forceReconnect();
      }
      
    } catch (error) {
      console.error('Failed to run diagnostics:', error);
      setDiagnosticResult('❌ Diagnostic check failed');
    } finally {
      setIsRunning(false);
    }
  };

  return (
    <div className="space-y-2">
      <button
        onClick={handleRunDiagnostics}
        disabled={isRunning}
        className={`
          px-4 py-2 text-sm font-medium rounded-md
          ${isRunning 
            ? 'bg-gray-300 text-gray-500 cursor-not-allowed' 
            : 'bg-blue-600 text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500'
          }
          ${className}
        `}
        title="Run polling chat connection diagnostics"
      >
        {isRunning ? (
          <>
            <svg className="animate-spin -ml-1 mr-2 h-4 w-4 text-gray-500 inline" fill="none" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Running...
          </>
        ) : (
          <>
            🔍 Run Diagnostics
          </>
        )}
      </button>
      
      {diagnosticResult && (
        <div className="text-sm p-2 rounded bg-gray-50 border">
          {diagnosticResult}
        </div>
      )}
    </div>
  );
};

export default DiagnosticButton;