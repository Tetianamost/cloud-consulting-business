import { createSlice, PayloadAction } from '@reduxjs/toolkit';

export type ConnectionStatus = 'disconnected' | 'connecting' | 'connected' | 'reconnecting' | 'failed' | 'polling';

interface ConnectionState {
  status: ConnectionStatus;
  lastConnected: string | null;
  reconnectAttempts: number;
  maxReconnectAttempts: number;
  reconnectDelay: number;
  maxReconnectDelay: number;
  error: string | null;
  latency: number | null;
  isHealthy: boolean;
  connectionId: string | null;
  // Polling-specific state
  isPolling: boolean;
  lastPollTime: number | null;
  pollInterval: number;
  errorCount: number;
  lastSuccessfulPoll: string | null;
}

const initialState: ConnectionState = {
  status: 'disconnected',
  lastConnected: null,
  reconnectAttempts: 0,
  maxReconnectAttempts: 10,
  reconnectDelay: 1000, // Start with 1 second
  maxReconnectDelay: 30000, // Max 30 seconds
  error: null,
  latency: null,
  isHealthy: false,
  connectionId: null,
  // Polling-specific state
  isPolling: false,
  lastPollTime: null,
  pollInterval: 3000, // 3 seconds default
  errorCount: 0,
  lastSuccessfulPoll: null,
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
    

    
    setHealthStatus: (state, action: PayloadAction<boolean>) => {
      state.isHealthy = action.payload;
    },
    
    // Connection cleanup
    cleanup: (state) => {
      state.status = 'disconnected';
      state.connectionId = null;
      state.latency = null;
      state.isHealthy = false;
      // Reset polling state
      state.isPolling = false;
      state.lastPollTime = null;
      state.errorCount = 0;
      state.lastSuccessfulPoll = null;
    },
    
    // Force reconnection
    forceReconnect: (state) => {
      state.status = 'reconnecting';
      state.reconnectAttempts = 0;
      state.reconnectDelay = 1000;
      state.error = null;
    },

    // Polling-specific actions
    setPollingStatus: (state, action: PayloadAction<boolean>) => {
      state.isPolling = action.payload;
      if (action.payload) {
        state.status = 'polling';
        state.lastPollTime = Date.now();
      } else {
        state.status = 'disconnected';
      }
    },

    updatePollTime: (state) => {
      state.lastPollTime = Date.now();
    },

    setPollInterval: (state, action: PayloadAction<number>) => {
      state.pollInterval = action.payload;
    },

    incrementErrorCount: (state) => {
      state.errorCount += 1;
    },

    resetErrorCount: (state) => {
      state.errorCount = 0;
    },

    setLastSuccessfulPoll: (state, action: PayloadAction<string>) => {
      state.lastSuccessfulPoll = action.payload;
      state.errorCount = 0; // Reset error count on successful poll
      state.isHealthy = true;
    },
  },
});

export const {
  setConnectionStatus,
  setConnectionId,
  incrementReconnectAttempts,
  resetReconnectAttempts,
  setMaxReconnectAttempts,
  setConnectionError,
  clearConnectionError,
  updateLatency,
  setHealthStatus,
  cleanup,
  forceReconnect,
  setPollingStatus,
  updatePollTime,
  setPollInterval,
  incrementErrorCount,
  resetErrorCount,
  setLastSuccessfulPoll,
} = connectionSlice.actions;

export default connectionSlice.reducer;