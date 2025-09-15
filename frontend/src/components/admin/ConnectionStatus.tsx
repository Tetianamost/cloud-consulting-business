import React, { useState, useEffect } from 'react';
import { Wifi, WifiOff, RotateCcw, AlertTriangle, CheckCircle, Clock } from 'lucide-react';
import { useAppSelector } from '../../store/hooks';
import { ConnectionStatus as ConnectionStatusType } from '../../store/slices/connectionSlice';
import DiagnosticButton from './DiagnosticButton';
import { pollingChatService } from '../../services/pollingChatService';

interface ConnectionStatusProps {
  showDetails?: boolean;
  className?: string;
  onReconnect?: () => void;
}

const ConnectionStatus: React.FC<ConnectionStatusProps> = ({
  showDetails = false,
  className = '',
  onReconnect,
}) => {
  const connectionState = useAppSelector(state => state.connection);
  const { status, isHealthy, latency, reconnectAttempts, maxReconnectAttempts, error } = connectionState;
  
  // Enhanced status from polling service
  const [pollingStatus, setPollingStatus] = useState<any>(null);
  const [statusMessage, setStatusMessage] = useState<string>('');

  useEffect(() => {
    // Get initial status
    const initialStatus = pollingChatService.getConnectionStatusInfo();
    const initialMessage = pollingChatService.getStatusMessage();
    setPollingStatus(initialStatus);
    setStatusMessage(initialMessage);

    // Subscribe to status changes
    const unsubscribe = pollingChatService.onStatusChange((status) => {
      setPollingStatus(status);
      setStatusMessage(pollingChatService.getStatusMessage());
    });

    return unsubscribe;
  }, []);

  const getStatusIcon = (status: ConnectionStatusType, isHealthy: boolean, pollingState?: string) => {
    // Use polling state if available for more accurate status
    if (pollingState) {
      switch (pollingState) {
        case 'connected':
          return <CheckCircle className="h-4 w-4 text-green-500" />;
        case 'polling':
          return <Clock className="h-4 w-4 text-blue-500 animate-pulse" />;
        case 'error':
          return <AlertTriangle className="h-4 w-4 text-yellow-500" />;
        case 'offline':
          return <WifiOff className="h-4 w-4 text-red-500" />;
      }
    }

    // Fallback to Redux status
    switch (status) {
      case 'connected':
        return isHealthy ? (
          <CheckCircle className="h-4 w-4 text-green-500" />
        ) : (
          <AlertTriangle className="h-4 w-4 text-yellow-500" />
        );
      case 'connecting':
      case 'reconnecting':
        return <RotateCcw className="h-4 w-4 text-blue-500 animate-spin" />;
      case 'disconnected':
      case 'failed':
        return <WifiOff className="h-4 w-4 text-red-500" />;
      default:
        return <Wifi className="h-4 w-4 text-gray-400" />;
    }
  };

  const getStatusText = (status: ConnectionStatusType, reconnectAttempts: number) => {
    // Use enhanced status message if available (from polling service)
    if (statusMessage) {
      return statusMessage;
    }

    // Use polling-specific status if available
    if (pollingStatus?.state) {
      switch (pollingStatus.state) {
        case 'connected':
          return 'Polling Active';
        case 'polling':
          return 'Polling for Messages...';
        case 'error':
          return 'Polling Error (Retrying)';
        case 'offline':
          return 'Offline (Messages Queued)';
        default:
          return 'Polling Status Unknown';
      }
    }

    // Fallback to Redux status
    switch (status) {
      case 'connected':
        return isHealthy ? 'Connected' : 'Connected (Unstable)';
      case 'connecting':
        return 'Connecting...';
      case 'reconnecting':
        return `Reconnecting... (${reconnectAttempts}/${maxReconnectAttempts})`;
      case 'disconnected':
        return 'Disconnected';
      case 'failed':
        return 'Connection Failed';
      default:
        return 'Unknown';
    }
  };

  const getStatusColor = (status: ConnectionStatusType, isHealthy: boolean, pollingState?: string) => {
    // Use polling state if available for more accurate coloring
    if (pollingState) {
      switch (pollingState) {
        case 'connected':
          return 'text-green-600';
        case 'polling':
          return 'text-blue-600';
        case 'error':
          return 'text-yellow-600';
        case 'offline':
          return 'text-red-600';
      }
    }

    // Fallback to Redux status
    switch (status) {
      case 'connected':
        return isHealthy ? 'text-green-600' : 'text-yellow-600';
      case 'connecting':
      case 'reconnecting':
        return 'text-blue-600';
      case 'disconnected':
      case 'failed':
        return 'text-red-600';
      default:
        return 'text-gray-600';
    }
  };

  const formatLatency = (latency: number | null): string => {
    if (latency === null) return 'N/A';
    if (latency < 1000) return `${latency}ms`;
    return `${(latency / 1000).toFixed(1)}s`;
  };

  const formatUptime = (uptime: number): string => {
    const seconds = Math.floor(uptime / 1000);
    const minutes = Math.floor(seconds / 60);
    const hours = Math.floor(minutes / 60);
    
    if (hours > 0) {
      return `${hours}h ${minutes % 60}m`;
    } else if (minutes > 0) {
      return `${minutes}m ${seconds % 60}s`;
    } else {
      return `${seconds}s`;
    }
  };

  const handleReconnect = () => {
    if (onReconnect) {
      onReconnect();
    } else {
      // Use polling service's force reconnect
      pollingChatService.forceReconnect();
    }
  };

  if (!showDetails) {
    // Simple status indicator
    return (
      <div className={`flex items-center space-x-2 ${className}`}>
        {getStatusIcon(status, isHealthy, pollingStatus?.state)}
        <span className={`text-sm ${getStatusColor(status, isHealthy, pollingStatus?.state)}`}>
          {getStatusText(status, reconnectAttempts)}
        </span>
      </div>
    );
  }

  // Detailed status display
  return (
    <div className={`bg-white rounded-lg border border-gray-200 p-4 ${className}`}>
      <div className="flex items-center justify-between mb-3">
        <h3 className="text-sm font-medium text-gray-900">Connection Status</h3>
        <div className="flex items-center space-x-2">
          {(status === 'disconnected' || status === 'failed' || pollingStatus?.state === 'offline' || pollingStatus?.state === 'error') && (
            <DiagnosticButton className="text-xs px-2 py-1" />
          )}
          {(status === 'disconnected' || status === 'failed' || pollingStatus?.state === 'offline' || pollingStatus?.state === 'error') && (
            <button
              onClick={handleReconnect}
              className="text-sm text-blue-600 hover:text-blue-800 transition-colors"
            >
              Reconnect
            </button>
          )}
        </div>
      </div>

      <div className="space-y-3">
        {/* Status */}
        <div className="flex items-center justify-between">
          <span className="text-sm text-gray-600">Status:</span>
          <div className="flex items-center space-x-2">
            {getStatusIcon(status, isHealthy, pollingStatus?.state)}
            <span className={`text-sm font-medium ${getStatusColor(status, isHealthy, pollingStatus?.state)}`}>
              {getStatusText(status, reconnectAttempts)}
            </span>
          </div>
        </div>

        {/* Latency */}
        {(status === 'connected' || pollingStatus?.state === 'connected') && (
          <div className="flex items-center justify-between">
            <span className="text-sm text-gray-600">Latency:</span>
            <span className={`text-sm ${
              latency && latency > 1000 ? 'text-yellow-600' : 'text-green-600'
            }`}>
              {formatLatency(latency)}
            </span>
          </div>
        )}

        {/* Error count for polling */}
        {pollingStatus && pollingStatus.errorCount > 0 && (
          <div className="flex items-center justify-between">
            <span className="text-sm text-gray-600">Errors:</span>
            <span className={`text-sm ${
              pollingStatus.errorCount >= 3 ? 'text-red-600' : 'text-yellow-600'
            }`}>
              {pollingStatus.errorCount}
            </span>
          </div>
        )}

        {/* Uptime */}
        {pollingStatus && pollingStatus.uptime && (
          <div className="flex items-center justify-between">
            <span className="text-sm text-gray-600">Uptime:</span>
            <span className="text-sm text-gray-900">
              {formatUptime(pollingStatus.uptime)}
            </span>
          </div>
        )}

        {/* Last successful connection */}
        {pollingStatus && pollingStatus.lastSuccessfulConnection && (
          <div className="flex items-center justify-between">
            <span className="text-sm text-gray-600">Last Success:</span>
            <span className="text-sm text-gray-900">
              {new Date(pollingStatus.lastSuccessfulConnection).toLocaleTimeString()}
            </span>
          </div>
        )}

        {/* Reconnection attempts */}
        {status === 'reconnecting' && (
          <div className="flex items-center justify-between">
            <span className="text-sm text-gray-600">Attempts:</span>
            <span className="text-sm text-gray-900">
              {reconnectAttempts} / {maxReconnectAttempts}
            </span>
          </div>
        )}

        {/* Error message */}
        {error && (
          <div className="mt-3 p-2 bg-red-50 border border-red-200 rounded">
            <div className="flex items-start space-x-2">
              <AlertTriangle className="h-4 w-4 text-red-500 mt-0.5 flex-shrink-0" />
              <span className="text-sm text-red-700">{error}</span>
            </div>
          </div>
        )}

        {/* Health warning */}
        {((status === 'connected' && !isHealthy) || (pollingStatus && !pollingStatus.isHealthy)) && (
          <div className="mt-3 p-2 bg-yellow-50 border border-yellow-200 rounded">
            <div className="flex items-start space-x-2">
              <AlertTriangle className="h-4 w-4 text-yellow-500 mt-0.5 flex-shrink-0" />
              <span className="text-sm text-yellow-700">
                Connection is unstable. Messages may be delayed.
              </span>
            </div>
          </div>
        )}

        {/* Offline warning */}
        {pollingStatus?.state === 'offline' && (
          <div className="mt-3 p-2 bg-red-50 border border-red-200 rounded">
            <div className="flex items-start space-x-2">
              <WifiOff className="h-4 w-4 text-red-500 mt-0.5 flex-shrink-0" />
              <span className="text-sm text-red-700">
                Network is offline. Messages will be queued and sent when connection is restored.
              </span>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default ConnectionStatus;