import { useEffect, useCallback, useRef } from 'react';
import { useAppDispatch, useAppSelector } from '../store/hooks';
import websocketService, { ChatRequest } from '../services/websocketService';
import { setError } from '../store/slices/chatSlice';

export interface UseWebSocketReturn {
  // Connection state
  isConnected: boolean;
  isConnecting: boolean;
  isReconnecting: boolean;
  connectionError: string | null;
  isHealthy: boolean;
  latency: number | null;
  reconnectAttempts: number;
  
  // Actions
  connect: () => Promise<void>;
  disconnect: () => void;
  sendMessage: (request: ChatRequest) => string;
  forceReconnect: () => void;
}

export const useWebSocket = (autoConnect = true): UseWebSocketReturn => {
  const dispatch = useAppDispatch();
  const connectionState = useAppSelector(state => state.connection);
  const isInitializedRef = useRef(false);

  // Connection state derived from Redux store
  const isConnected = connectionState.status === 'connected';
  const isConnecting = connectionState.status === 'connecting';
  const isReconnecting = connectionState.status === 'reconnecting';
  const connectionError = connectionState.error;
  const isHealthy = connectionState.isHealthy;
  const latency = connectionState.latency;
  const reconnectAttempts = connectionState.reconnectAttempts;

  // Memoized actions
  const connect = useCallback(async (): Promise<void> => {
    try {
      await websocketService.connect();
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to connect';
      dispatch(setError(errorMessage));
      throw error;
    }
  }, [dispatch]);

  const disconnect = useCallback((): void => {
    websocketService.disconnect();
  }, []);

  const sendMessage = useCallback((request: ChatRequest): string => {
    if (!isConnected && !isReconnecting) {
      dispatch(setError('Not connected to chat service'));
      return '';
    }
    
    return websocketService.sendChatMessage(request);
  }, [isConnected, isReconnecting, dispatch]);

  const forceReconnect = useCallback((): void => {
    websocketService.forceReconnect();
  }, []);

  // Auto-connect on mount if enabled
  useEffect(() => {
    if (autoConnect && !isInitializedRef.current) {
      isInitializedRef.current = true;
      
      // Only connect if not already connected or connecting
      if (connectionState.status === 'disconnected') {
        connect().catch(error => {
          console.error('Auto-connect failed:', error);
        });
      }
    }
  }, [autoConnect, connect, connectionState.status]);

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      // Only disconnect if this is the last component using the WebSocket
      // In a real app, you might want to implement reference counting
      if (!autoConnect) {
        disconnect();
      }
    };
  }, [autoConnect, disconnect]);

  return {
    // Connection state
    isConnected,
    isConnecting,
    isReconnecting,
    connectionError,
    isHealthy,
    latency,
    reconnectAttempts,
    
    // Actions
    connect,
    disconnect,
    sendMessage,
    forceReconnect,
  };
};

export default useWebSocket;