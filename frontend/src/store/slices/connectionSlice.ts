import { createSlice, PayloadAction } from '@reduxjs/toolkit';

export type ConnectionStatus = 'disconnected' | 'connecting' | 'connected' | 'reconnecting' | 'failed';

interface ConnectionState {
  status: ConnectionStatus;
  webSocket: WebSocket | null;
  lastConnected: string | null;
  reconnectAttempts: number;
  maxReconnectAttempts: number;
  reconnectDelay: number;
  maxReconnectDelay: number;
  error: string | null;
  latency: number | null;
  lastPingTime: number | null;
  isHealthy: boolean;
  connectionId: string | null;
}

const initialState: ConnectionState = {
  status: 'disconnected',
  webSocket: null,
  lastConnected: null,
  reconnectAttempts: 0,
  maxReconnectAttempts: 10,
  reconnectDelay: 1000, // Start with 1 second
  maxReconnectDelay: 30000, // Max 30 seconds
  error: null,
  latency: null,
  lastPingTime: null,
  isHealthy: false,
  connectionId: null,
};

const connectionSlice = createSlice({
  name: 'connection',
  initialState,
  reducers: {
    // Connection status management
    setConnectionStatus: (state, action: PayloadAction<ConnectionStatus>) => {
      state.status = action.payload;
      
      if (action.payload === 'connected') {
        state.lastConnected = new Date().toISOString();
        state.reconnectAttempts = 0;
        state.reconnectDelay = 1000; // Reset delay
        state.error = null;
        state.isHealthy = true;
      } else if (action.payload === 'disconnected' || action.payload === 'failed') {
        state.isHealthy = false;
      }
    },
    
    setWebSocket: (state, action: PayloadAction<WebSocket | null>) => {
      state.webSocket = action.payload;
    },
    
    setConnectionId: (state, action: PayloadAction<string | null>) => {
      state.connectionId = action.payload;
    },
    
    // Reconnection management with exponential backoff
    incrementReconnectAttempts: (state) => {
      state.reconnectAttempts += 1;
      
      // Exponential backoff with jitter
      const baseDelay = Math.min(
        state.reconnectDelay * Math.pow(2, state.reconnectAttempts - 1),
        state.maxReconnectDelay
      );
      
      // Add jitter (±25%)
      const jitter = baseDelay * 0.25 * (Math.random() - 0.5);
      state.reconnectDelay = Math.max(1000, baseDelay + jitter);
    },
    
    resetReconnectAttempts: (state) => {
      state.reconnectAttempts = 0;
      state.reconnectDelay = 1000;
    },
    
    setMaxReconnectAttempts: (state, action: PayloadAction<number>) => {
      state.maxReconnectAttempts = action.payload;
    },
    
    // Error handling
    setConnectionError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
      if (action.payload) {
        state.isHealthy = false;
      }
    },
    
    clearConnectionError: (state) => {
      state.error = null;
    },
    
    // Health monitoring
    updateLatency: (state, action: PayloadAction<number>) => {
      state.latency = action.payload;
      state.isHealthy = action.payload < 5000; // Consider unhealthy if latency > 5s
    },
    
    setPingTime: (state, action: PayloadAction<number>) => {
      state.lastPingTime = action.payload;
    },
    
    setHealthStatus: (state, action: PayloadAction<boolean>) => {
      state.isHealthy = action.payload;
    },
    
    // Connection cleanup
    cleanup: (state) => {
      if (state.webSocket) {
        try {
          state.webSocket.close();
        } catch (error) {
          console.warn('Error closing WebSocket:', error);
        }
      }
      
      state.webSocket = null;
      state.status = 'disconnected';
      state.connectionId = null;
      state.latency = null;
      state.lastPingTime = null;
      state.isHealthy = false;
    },
    
    // Force reconnection
    forceReconnect: (state) => {
      state.status = 'reconnecting';
      state.reconnectAttempts = 0;
      state.reconnectDelay = 1000;
      state.error = null;
    },
  },
});

export const {
  setConnectionStatus,
  setWebSocket,
  setConnectionId,
  incrementReconnectAttempts,
  resetReconnectAttempts,
  setMaxReconnectAttempts,
  setConnectionError,
  clearConnectionError,
  updateLatency,
  setPingTime,
  setHealthStatus,
  cleanup,
  forceReconnect,
} = connectionSlice.actions;

export default connectionSlice.reducer;