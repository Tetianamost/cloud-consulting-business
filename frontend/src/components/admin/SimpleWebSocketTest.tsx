import React, { useState, useEffect } from 'react';
import { simpleWebSocketService } from '../../services/simpleWebSocketService';

interface SimpleWebSocketTestProps {
  className?: string;
}

export const SimpleWebSocketTest: React.FC<SimpleWebSocketTestProps> = ({ className }) => {
  const [isConnected, setIsConnected] = useState(false);
  const [connectionStatus, setConnectionStatus] = useState<string>('disconnected');
  const [testMessage, setTestMessage] = useState('');
  const [logs, setLogs] = useState<string[]>([]);

  const addLog = (message: string) => {
    const timestamp = new Date().toLocaleTimeString();
    setLogs(prev => [...prev.slice(-9), `[${timestamp}] ${message}`]);
  };

  useEffect(() => {
    // Check connection status periodically
    const statusInterval = setInterval(() => {
      const connected = simpleWebSocketService.isConnected();
      setIsConnected(connected);
      setConnectionStatus(connected ? 'connected' : 'disconnected');
    }, 1000);

    return () => {
      clearInterval(statusInterval);
    };
  }, []);

  const handleConnect = async () => {
    try {
      setConnectionStatus('connecting');
      addLog('Attempting to connect...');
      await simpleWebSocketService.connect();
      addLog('Connected successfully!');
      setConnectionStatus('connected');
      setIsConnected(true);
    } catch (error) {
      addLog(`Connection failed: ${error instanceof Error ? error.message : 'Unknown error'}`);
      setConnectionStatus('failed');
      setIsConnected(false);
    }
  };

  const handleDisconnect = () => {
    addLog('Disconnecting...');
    simpleWebSocketService.disconnect();
    setConnectionStatus('disconnected');
    setIsConnected(false);
    addLog('Disconnected');
  };

  const handleSendMessage = () => {
    if (!testMessage.trim()) {
      addLog('Please enter a message to send');
      return;
    }

    const success = simpleWebSocketService.send({
      type: 'message',
      content: testMessage,
      timestamp: new Date().toISOString()
    });

    if (success) {
      addLog(`Sent: ${testMessage}`);
      setTestMessage('');
    } else {
      addLog('Failed to send message - not connected');
    }
  };

  const handleSendHeartbeat = () => {
    const success = simpleWebSocketService.send({
      type: 'heartbeat',
      timestamp: new Date().toISOString()
    });

    if (success) {
      addLog('Sent heartbeat');
    } else {
      addLog('Failed to send heartbeat - not connected');
    }
  };

  const getStatusColor = () => {
    switch (connectionStatus) {
      case 'connected': return 'text-green-600';
      case 'connecting': return 'text-yellow-600';
      case 'failed': return 'text-red-600';
      default: return 'text-gray-600';
    }
  };

  return (
    <div className={`p-6 bg-white rounded-lg shadow-md ${className}`}>
      <h3 className="text-lg font-semibold mb-4">Simple WebSocket Test</h3>
      
      {/* Connection Status */}
      <div className="mb-4">
        <div className="flex items-center gap-2 mb-2">
          <span className="font-medium">Status:</span>
          <span className={`font-semibold ${getStatusColor()}`}>
            {connectionStatus}
          </span>
          <div className={`w-3 h-3 rounded-full ${isConnected ? 'bg-green-500' : 'bg-red-500'}`} />
        </div>
      </div>

      {/* Connection Controls */}
      <div className="flex gap-2 mb-4">
        <button
          onClick={handleConnect}
          disabled={isConnected || connectionStatus === 'connecting'}
          className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:bg-gray-400 disabled:cursor-not-allowed"
        >
          Connect
        </button>
        <button
          onClick={handleDisconnect}
          disabled={!isConnected}
          className="px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600 disabled:bg-gray-400 disabled:cursor-not-allowed"
        >
          Disconnect
        </button>
      </div>

      {/* Message Testing */}
      <div className="mb-4">
        <div className="flex gap-2 mb-2">
          <input
            type="text"
            value={testMessage}
            onChange={(e) => setTestMessage(e.target.value)}
            placeholder="Enter test message..."
            className="flex-1 px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
            onKeyPress={(e) => e.key === 'Enter' && handleSendMessage()}
          />
          <button
            onClick={handleSendMessage}
            disabled={!isConnected}
            className="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600 disabled:bg-gray-400 disabled:cursor-not-allowed"
          >
            Send
          </button>
        </div>
        <button
          onClick={handleSendHeartbeat}
          disabled={!isConnected}
          className="px-4 py-2 bg-purple-500 text-white rounded hover:bg-purple-600 disabled:bg-gray-400 disabled:cursor-not-allowed"
        >
          Send Heartbeat
        </button>
      </div>

      {/* Logs */}
      <div className="mb-4">
        <h4 className="font-medium mb-2">Connection Logs:</h4>
        <div className="bg-gray-100 p-3 rounded max-h-40 overflow-y-auto">
          {logs.length === 0 ? (
            <p className="text-gray-500 text-sm">No logs yet...</p>
          ) : (
            logs.map((log, index) => (
              <div key={index} className="text-sm font-mono text-gray-700 mb-1">
                {log}
              </div>
            ))
          )}
        </div>
      </div>

      {/* Clear Logs */}
      <button
        onClick={() => setLogs([])}
        className="px-3 py-1 text-sm bg-gray-500 text-white rounded hover:bg-gray-600"
      >
        Clear Logs
      </button>
    </div>
  );
};

export default SimpleWebSocketTest;