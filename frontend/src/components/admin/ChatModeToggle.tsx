import React, { useState, useEffect } from 'react';
import { useSelector } from 'react-redux';
import { RootState } from '../../store';
import { chatModeManager, ChatMode, ChatConfig } from '../../services/chatModeManager';

interface ChatModeToggleProps {
  className?: string;
}

interface PerformanceStats {
  websocket: {
    connectionTime: number;
    disconnectionCount: number;
    lastDisconnection: Date | null;
  };
  polling: {
    errorCount: number;
    lastError: Date | null;
  };
}

const ChatModeToggle: React.FC<ChatModeToggleProps> = ({ className = '' }) => {
  const [currentMode, setCurrentMode] = useState<ChatMode>('auto');
  const [activeService, setActiveService] = useState<'websocket' | 'polling' | null>(null);
  const [config, setConfig] = useState<ChatConfig | null>(null);
  const [fallbackState, setFallbackState] = useState<any>(null);
  const [performanceStats, setPerformanceStats] = useState<PerformanceStats | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [showAdvanced, setShowAdvanced] = useState(false);

  const connectionStatus = useSelector((state: RootState) => state.connection.status);
  const connectionError = useSelector((state: RootState) => state.connection.error);

  // Update state from chat mode manager
  useEffect(() => {
    const updateState = () => {
      setCurrentMode(chatModeManager.getCurrentMode());
      setActiveService(chatModeManager.getActiveService());
      setConfig(chatModeManager.getConfiguration());
      setFallbackState(chatModeManager.getFallbackState());
      setPerformanceStats(chatModeManager.getPerformanceMetrics());
    };

    // Initial update
    updateState();

    // Update every 2 seconds
    const interval = setInterval(updateState, 2000);

    return () => clearInterval(interval);
  }, []);

  const handleModeChange = async (newMode: ChatMode) => {
    setIsLoading(true);
    setError(null);

    try {
      await chatModeManager.switchMode(newMode);
      setCurrentMode(newMode);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to switch mode');
    } finally {
      setIsLoading(false);
    }
  };

  const handleConfigUpdate = async (updates: Partial<ChatConfig>) => {
    setIsLoading(true);
    setError(null);

    try {
      await chatModeManager.updateConfiguration(updates);
      setConfig(chatModeManager.getConfiguration());
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update configuration');
    } finally {
      setIsLoading(false);
    }
  };

  const handleForceReconnect = () => {
    chatModeManager.forceReconnect();
  };

  const getStatusColor = () => {
    if (fallbackState?.isInFallback) return 'text-yellow-600';
    
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
    if (fallbackState?.isInFallback) return '⚠️';
    
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
        <h3 className="text-lg font-semibold text-gray-900">Chat Mode Control</h3>
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
                {chatModeManager.getStatusMessage()}
              </div>
              <div className="text-sm text-gray-600">
                Mode: {currentMode} | Active: {activeService || 'none'}
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

        {/* Fallback notification */}
        {fallbackState?.isInFallback && (
          <div className="mt-2 p-2 bg-yellow-50 border border-yellow-200 rounded text-sm">
            <strong>Fallback Mode:</strong> {fallbackState.fallbackReason}
            <br />
            <span className="text-gray-600">
              Started: {new Date(fallbackState.fallbackStartTime).toLocaleTimeString()}
            </span>
          </div>
        )}

        {/* Connection error */}
        {connectionError && (
          <div className="mt-2 p-2 bg-red-50 border border-red-200 rounded text-sm text-red-700">
            {connectionError}
          </div>
        )}
      </div>

      {/* Mode Selection */}
      <div className="mb-4">
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Chat Mode
        </label>
        <div className="grid grid-cols-3 gap-2">
          {(['websocket', 'polling', 'auto'] as ChatMode[]).map((mode) => (
            <button
              key={mode}
              onClick={() => handleModeChange(mode)}
              disabled={isLoading || currentMode === mode}
              className={`px-3 py-2 text-sm rounded border ${
                currentMode === mode
                  ? 'bg-blue-600 text-white border-blue-600'
                  : 'bg-white text-gray-700 border-gray-300 hover:bg-gray-50'
              } disabled:opacity-50`}
            >
              {mode.charAt(0).toUpperCase() + mode.slice(1)}
            </button>
          ))}
        </div>
      </div>

      {/* Advanced Configuration */}
      {showAdvanced && config && (
        <div className="space-y-4 border-t pt-4">
          <h4 className="font-medium text-gray-900">Advanced Configuration</h4>
          
          {/* Fallback Settings */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Enable WebSocket Fallback
              </label>
              <input
                type="checkbox"
                checked={config.enable_websocket_fallback}
                onChange={(e) => handleConfigUpdate({ enable_websocket_fallback: e.target.checked })}
                disabled={isLoading}
                className="rounded border-gray-300"
              />
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
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                WebSocket Timeout (seconds)
              </label>
              <input
                type="number"
                min="1"
                max="60"
                value={config.websocket_timeout}
                onChange={(e) => handleConfigUpdate({ websocket_timeout: parseInt(e.target.value) })}
                disabled={isLoading}
                className="w-full px-3 py-1 border border-gray-300 rounded text-sm"
              />
            </div>
            
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
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Fallback Delay (ms)
            </label>
            <input
              type="number"
              min="1000"
              max="30000"
              step="1000"
              value={config.fallback_delay}
              onChange={(e) => handleConfigUpdate({ fallback_delay: parseInt(e.target.value) })}
              disabled={isLoading}
              className="w-full px-3 py-1 border border-gray-300 rounded text-sm"
            />
          </div>
        </div>
      )}

      {/* Performance Metrics */}
      {showAdvanced && performanceStats && (
        <div className="border-t pt-4">
          <h4 className="font-medium text-gray-900 mb-2">Performance Metrics</h4>
          <div className="grid grid-cols-2 gap-4 text-sm">
            <div className="bg-gray-50 p-3 rounded">
              <h5 className="font-medium text-gray-700 mb-1">WebSocket</h5>
              <div>Connection Time: {performanceStats.websocket.connectionTime}ms</div>
              <div>Disconnections: {performanceStats.websocket.disconnectionCount}</div>
              {performanceStats.websocket.lastDisconnection && (
                <div className="text-gray-600">
                  Last: {new Date(performanceStats.websocket.lastDisconnection).toLocaleTimeString()}
                </div>
              )}
            </div>
            
            <div className="bg-gray-50 p-3 rounded">
              <h5 className="font-medium text-gray-700 mb-1">Polling</h5>
              <div>Errors: {performanceStats.polling.errorCount}</div>
              {performanceStats.polling.lastError && (
                <div className="text-gray-600">
                  Last Error: {new Date(performanceStats.polling.lastError).toLocaleTimeString()}
                </div>
              )}
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