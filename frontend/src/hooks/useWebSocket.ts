import { useEffect, useCallback, useRef } from 'react';
import { useAppDispatch, useAppSelector } from '../store/hooks';
import { pollingChatService } from '../services/pollingChatService';
import { setError } from '../store/slices/chatSlice';

export interface ChatRequest {
  message: string;
  session_id?: string;
  client_name?: string;
  context?: string;
  quick_action?: string;
}

export interface UsePollingChatReturn {
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
  sendMessage: (request: ChatRequest) => Promise<string>;
  forceReconnect: () => void;
}

// Polling-based chat hook for reliable communication
export const usePollingChat = (autoConnect = true): UsePollingChatReturn => {
  const dispatch = useAppDispatch();
  const connectionState = useAppSelector(state => state.connection);
  const isInitializedRef = useRef(false);

  // Connection state derived from Redux store
  const isConnected = connectionState.status === 'connected' || connectionState.status === 'polling';
  const isConnecting = connectionState.status === 'connecting';
  const isReconnecting = connectionState.status === 'reconnecting';
  const connectionError = connectionState.error;
  const isHealthy = pollingChatService.isHealthy();
  const latency = connectionState.latency;
  const reconnectAttempts = connectionState.reconnectAttempts;

  // Memoized actions
  const connect = useCallback(async (): Promise<void> => {
    try {
      pollingChatService.startPolling();
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to connect';
      dispatch(setError(errorMessage));
      throw error;
    }
  }, [dispatch]);

  const disconnect = useCallback((): void => {
    pollingChatService.stopPolling();
  }, []);

  const sendMessage = useCallback(async (request: ChatRequest): Promise<string> => {
    if (!isConnected && !isReconnecting) {
      dispatch(setError('Not connected to chat service'));
      return '';
    }
    
    try {
      return await pollingChatService.sendMessage(request);
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to send message';
      dispatch(setError(errorMessage));
      throw error;
    }
  }, [isConnected, isReconnecting, dispatch]);

  const forceReconnect = useCallback((): void => {
    pollingChatService.forceReconnect();
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
      // Only disconnect if this is the last component using polling
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

// Keep the old export name for backward compatibility, but use polling
export const useWebSocket = usePollingChat;

export default usePollingChat;