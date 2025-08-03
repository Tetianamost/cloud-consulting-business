import React, { useState } from 'react';
import { Wifi, WifiOff, RotateCcw, AlertTriangle, CheckCircle, Activity } from 'lucide-react';
import { useAppSelector } from '../../store/hooks';
import { ConnectionStatus as ConnectionStatusType } from '../../store/slices/connectionSlice';
import { connectionDiagnostics, DiagnosticReport } from '../../services/connectionDiagnostics';

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

  const getStatusIcon = (status: ConnectionStatusType, isHealthy: boolean) => {
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

  const getStatusColor = (status: ConnectionStatusType, isHealthy: boolean) => {
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

  if (!showDetails) {
    // Simple status indicator
    return (
      <div className={`flex items-center space-x-2 ${className}`}>
        {getStatusIcon(status, isHealthy)}
        <span className={`text-sm ${getStatusColor(status, isHealthy)}`}>
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
        {(status === 'disconnected' || status === 'failed') && onReconnect && (
          <button
            onClick={onReconnect}
            className="text-sm text-blue-600 hover:text-blue-800 transition-colors"
          >
            Reconnect
          </button>
        )}
      </div>

      <div className="space-y-3">
        {/* Status */}
        <div className="flex items-center justify-between">
          <span className="text-sm text-gray-600">Status:</span>
          <div className="flex items-center space-x-2">
            {getStatusIcon(status, isHealthy)}
            <span className={`text-sm font-medium ${getStatusColor(status, isHealthy)}`}>
              {getStatusText(status, reconnectAttempts)}
            </span>
          </div>
        </div>

        {/* Latency */}
        {status === 'connected' && (
          <div className="flex items-center justify-between">
            <span className="text-sm text-gray-600">Latency:</span>
            <span className={`text-sm ${
              latency && latency > 1000 ? 'text-yellow-600' : 'text-green-600'
            }`}>
              {formatLatency(latency)}
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
        {status === 'connected' && !isHealthy && (
          <div className="mt-3 p-2 bg-yellow-50 border border-yellow-200 rounded">
            <div className="flex items-start space-x-2">
              <AlertTriangle className="h-4 w-4 text-yellow-500 mt-0.5 flex-shrink-0" />
              <span className="text-sm text-yellow-700">
                Connection is unstable. Messages may be delayed.
              </span>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default ConnectionStatus;