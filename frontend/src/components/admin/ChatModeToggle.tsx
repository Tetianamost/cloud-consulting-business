import React, { useState, useEffect } from 'react';
import { useSelector } from 'react-redux';
import { RootState } from '../../store';
import { pollingChatService } from '../../services/pollingChatService';

interface ChatModeToggleProps {
  className?: string;
}

interface PollingConfig {
  polling_interval: number;
  max_reconnect_attempts: number;
}

interface PerformanceStats {
  polling: {
    errorCount: number;
    lastError: Date | null;
    successRate: number;
    averageResponseTime: number;
  };
}

const ChatModeToggle: React.FC<ChatModeToggleProps> = ({ className = '' }) => {
  const [config, setConfig] = useState<PollingConfig>({
    polling_interval: 3000,
    max_reconnect_attempts: 3
  });
  const [performanceStats, setPerformanceStats] = useState<PerformanceStats | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [showAdvanced, setShowAdvanced] = useState(false);

  const connectionStatus = useSelector((state: RootState) => state.connection.status);
  const connectionError = useSelector((state: RootState) => state.connection.error);

  // Update state from polling chat service
  useEffect(() => {
    const updateState = () => {
      // Get connection status info from polling service
      const statusInfo = pollingChatService.getConnectionStatusInfo();
      
      // Create performance stats from connection info
      setPerformanceStats({
        polling: {
          errorCount: statusInfo.errorCount || 0,
          lastError: statusInfo.lastError ? new Date(statusInfo.lastError) : null,
          successRate: statusInfo.isHealthy ? 100 : 0, // Simple health-based success rate
          averageResponseTime: 0 // Not available in current interface
        }
      });
    };

    // Initial update
    updateState();

    // Update every 2 seconds
    const interval = setInterval(updateState, 2000);

    return () => clearInterval(interval);
  }, []);

  const handleConfigUpdate = async (updates: Partial<PollingConfig>) => {
    setIsLoading(true);
    setError(null);

    try {
      // Update local config
      setConfig(prev => ({ ...prev, ...updates }));
      
      // Update polling service configuration
      if (updates.polling_interval) {
        pollingChatService.setPollingInterval(updates.polling_interval);
      }
      
      console.log('Polling configuration updated:', updates);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update configuration');
    } finally {
      setIsLoading(false);
    }
  };

  const handleForceReconnect = () => {
    pollingChatService.forceReconnect();
  };

  const getStatusColor = () => {
    switch (connectionStatus) {
      case 'connected':
      case 'polling':
        return 'text-green-600';
      case 'connecting':
      case 'reconnecting':
        return 'text-yellow-600';
      case 'disconnected':
      case 'failed':
        return 'text-red-600';
      default:
        return 'text-gray-600';
    }
  };

  const getStatusIcon = () => {
    switch (connectionStatus) {
      case 'connected':
      case 'polling':
        return '🟢';
      case 'connecting':
      case 'reconnecting':
        return '🟡';
      case 'disconnected':
      case 'failed':
        return '🔴';
      default:
        return '⚪';
    }
  };

  return (
    <div className={`bg-white rounded-lg shadow-sm border border-gray-200 p-4 ${className}`}>
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-lg font-semibold text-gray-900">Polling Chat Control</h3>
        <button
          onClick={() => setShowAdvanced(!showAdvanced)}
          className="text-sm text-blue-600 hover:text-blue-800"
        >
          {showAdvanced ? 'Hide Advanced' : 'Show Advanced'}
        </button>
      </div>

      {/* Current Status */}
      <div className="mb-4 p-3 bg-gray-50 rounded-lg">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <span className="text-lg">{getStatusIcon()}</span>
            <div>
              <div className={`font-medium ${getStatusColor()}`}>
                {pollingChatService.getStatusMessage()}
              </div>
              <div className="text-sm text-gray-600">
                Mode: Polling | Status: {connectionStatus}
              </div>
            </div>
          </div>
          <button
            onClick={handleForceReconnect}
            disabled={isLoading}
            className="px-3 py-1 text-sm bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
          >
            Reconnect
          </button>
        </div>

        {/* Connection error */}
        {connectionError && (
          <div className="mt-2 p-2 bg-red-50 border border-red-200 rounded text-sm text-red-700">
            {connectionError}
          </div>
        )}
      </div>

      {/* Chat Mode Info */}
      <div className="mb-4 p-3 bg-blue-50 border border-blue-200 rounded-lg">
        <div className="flex items-center space-x-2">
          <span className="text-blue-600">ℹ️</span>
          <div className="text-sm text-blue-800">
            <strong>Polling Mode Active:</strong> Chat uses HTTP polling for reliable communication.
            This provides better stability and reliability.
          </div>
        </div>
      </div>

      {/* Advanced Configuration */}
      {showAdvanced && (
        <div className="space-y-4 border-t pt-4">
          <h4 className="font-medium text-gray-900">Polling Configuration</h4>
          
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Polling Interval (ms)
              </label>
              <input
                type="number"
                min="1000"
                max="30000"
                step="1000"
                value={config.polling_interval}
                onChange={(e) => handleConfigUpdate({ polling_interval: parseInt(e.target.value) })}
                disabled={isLoading}
                className="w-full px-3 py-1 border border-gray-300 rounded text-sm"
              />
              <p className="text-xs text-gray-500 mt-1">
                How often to check for new messages (1-30 seconds)
              </p>
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Max Reconnect Attempts
              </label>
              <input
                type="number"
                min="1"
                max="10"
                value={config.max_reconnect_attempts}
                onChange={(e) => handleConfigUpdate({ max_reconnect_attempts: parseInt(e.target.value) })}
                disabled={isLoading}
                className="w-full px-3 py-1 border border-gray-300 rounded text-sm"
              />
              <p className="text-xs text-gray-500 mt-1">
                Maximum retry attempts on connection failure
              </p>
            </div>
          </div>
        </div>
      )}

      {/* Performance Metrics */}
      {showAdvanced && performanceStats && (
        <div className="border-t pt-4">
          <h4 className="font-medium text-gray-900 mb-2">Performance Metrics</h4>
          <div className="bg-gray-50 p-3 rounded">
            <h5 className="font-medium text-gray-700 mb-2">Polling Performance</h5>
            <div className="grid grid-cols-2 gap-4 text-sm">
              <div>
                <div>Success Rate: {performanceStats.polling.successRate.toFixed(1)}%</div>
                <div>Errors: {performanceStats.polling.errorCount}</div>
              </div>
              <div>
                <div>Avg Response: {performanceStats.polling.averageResponseTime}ms</div>
                {performanceStats.polling.lastError && (
                  <div className="text-gray-600">
                    Last Error: {performanceStats.polling.lastError.toLocaleTimeString()}
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Error Display */}
      {error && (
        <div className="mt-4 p-3 bg-red-50 border border-red-200 rounded text-sm text-red-700">
          {error}
        </div>
      )}

      {/* Loading Indicator */}
      {isLoading && (
        <div className="mt-4 flex items-center justify-center">
          <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
          <span className="ml-2 text-sm text-gray-600">Updating...</span>
        </div>
      )}
    </div>
  );
};

export default ChatModeToggle;