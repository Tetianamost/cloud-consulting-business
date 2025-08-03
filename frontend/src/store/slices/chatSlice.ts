import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';

// Types
export interface ChatMessage {
  id: string;
  type: 'user' | 'assistant' | 'system';
  content: string;
  timestamp: string;
  session_id: string;
  status?: 'sending' | 'sent' | 'delivered' | 'failed';
  metadata?: Record<string, any>;
}

export interface ChatSession {
  id: string;
  consultant_id: string;
  client_name?: string;
  context?: string;
  status: 'active' | 'inactive' | 'expired';
  created_at: string;
  updated_at: string;
  last_activity: string;
  expires_at?: string;
  metadata?: Record<string, any>;
}

export interface SessionContext {
  client_name?: string;
  meeting_type?: string;
  project_context?: string;
  service_types?: string[];
  cloud_providers?: string[];
  custom_fields?: Record<string, string>;
}

interface ChatState {
  // Session management
  currentSession: ChatSession | null;
  sessions: ChatSession[];
  sessionContext: SessionContext;
  
  // Message management
  messages: ChatMessage[];
  messageHistory: Record<string, ChatMessage[]>; // sessionId -> messages
  
  // UI state
  isLoading: boolean;
  isTyping: boolean;
  error: string | null;
  
  // Optimistic updates
  pendingMessages: ChatMessage[];
  
  // Settings
  settings: {
    showTimestamps: boolean;
    autoScroll: boolean;
    soundEnabled: boolean;
  };
}

const initialState: ChatState = {
  currentSession: null,
  sessions: [],
  sessionContext: {},
  messages: [],
  messageHistory: {},
  isLoading: false,
  isTyping: false,
  error: null,
  pendingMessages: [],
  settings: {
    showTimestamps: true,
    autoScroll: true,
    soundEnabled: false,
  },
};

// Async thunks for API calls
export const createSession = createAsyncThunk<ChatSession, Partial<ChatSession>>(
  'chat/createSession',
  async (sessionData: Partial<ChatSession>) => {
    const token = localStorage.getItem('adminToken');
    const response = await fetch('/api/v1/admin/chat/sessions', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
      body: JSON.stringify(sessionData),
    });
    
    if (!response.ok) {
      throw new Error('Failed to create session');
    }
    
    return response.json() as Promise<ChatSession>;
  }
);

export const loadSessionHistory = createAsyncThunk<ChatMessage[], string>(
  'chat/loadSessionHistory',
  async (sessionId: string) => {
    const token = localStorage.getItem('adminToken');
    const response = await fetch(`/api/v1/admin/chat/sessions/${sessionId}/history`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });
    
    if (!response.ok) {
      throw new Error('Failed to load session history');
    }
    
    return response.json() as Promise<ChatMessage[]>;
  }
);

export const updateSessionContext = createAsyncThunk<
  SessionContext,
  { sessionId: string; context: SessionContext }
>(
  'chat/updateSessionContext',
  async ({ sessionId, context }: { sessionId: string; context: SessionContext }) => {
    const token = localStorage.getItem('adminToken');
    const response = await fetch(`/api/v1/admin/chat/sessions/${sessionId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
      body: JSON.stringify({ context }),
    });
    
    if (!response.ok) {
      throw new Error('Failed to update session context');
    }
    
    return response.json() as Promise<SessionContext>;
  }
);

const chatSlice = createSlice({
  name: 'chat',
  initialState,
  reducers: {
    // Session management
    setCurrentSession: (state, action: PayloadAction<ChatSession>) => {
      state.currentSession = action.payload;
      // Load messages for this session
      if (state.messageHistory[action.payload.id]) {
        state.messages = state.messageHistory[action.payload.id];
      } else {
        state.messages = [];
      }
    },
    
    clearCurrentSession: (state) => {
      state.currentSession = null;
      state.messages = [];
      state.sessionContext = {};
    },
    
    updateSessionContext: (state, action: PayloadAction<Partial<SessionContext>>) => {
      state.sessionContext = { ...state.sessionContext, ...action.payload };
    },
    
    addSession: (state, action: PayloadAction<ChatSession>) => {
      const existingIndex = state.sessions.findIndex(s => s.id === action.payload.id);
      if (existingIndex >= 0) {
        state.sessions[existingIndex] = action.payload;
      } else {
        state.sessions.push(action.payload);
      }
    },
    
    // Message management with optimistic updates
    addMessage: (state, action: PayloadAction<ChatMessage>) => {
      const message = action.payload;
      
      // Add to current messages
      state.messages.push(message);
      
      // Add to message history
      if (!state.messageHistory[message.session_id]) {
        state.messageHistory[message.session_id] = [];
      }
      state.messageHistory[message.session_id].push(message);
      
      // Remove from pending if it was there
      state.pendingMessages = state.pendingMessages.filter(m => m.id !== message.id);
    },
    
    addOptimisticMessage: (state, action: PayloadAction<ChatMessage>) => {
      const message = { ...action.payload, status: 'sending' as const };
      
      // Add to current messages immediately
      state.messages.push(message);
      
      // Add to pending messages for tracking
      state.pendingMessages.push(message);
    },
    
    updateMessageStatus: (state, action: PayloadAction<{ id: string; status: ChatMessage['status'] }>) => {
      const { id, status } = action.payload;
      
      // Update in current messages
      const messageIndex = state.messages.findIndex(m => m.id === id);
      if (messageIndex >= 0) {
        state.messages[messageIndex].status = status;
      }
      
      // Update in message history
      Object.keys(state.messageHistory).forEach(sessionId => {
        const historyIndex = state.messageHistory[sessionId].findIndex(m => m.id === id);
        if (historyIndex >= 0) {
          state.messageHistory[sessionId][historyIndex].status = status;
        }
      });
      
      // Update in pending messages
      const pendingIndex = state.pendingMessages.findIndex(m => m.id === id);
      if (pendingIndex >= 0) {
        state.pendingMessages[pendingIndex].status = status;
      }
    },
    
    removeFailedMessage: (state, action: PayloadAction<string>) => {
      const messageId = action.payload;
      
      // Remove from current messages
      state.messages = state.messages.filter(m => m.id !== messageId);
      
      // Remove from pending messages
      state.pendingMessages = state.pendingMessages.filter(m => m.id !== messageId);
      
      // Remove from message history
      Object.keys(state.messageHistory).forEach(sessionId => {
        state.messageHistory[sessionId] = state.messageHistory[sessionId].filter(m => m.id !== messageId);
      });
    },
    
    clearMessages: (state) => {
      state.messages = [];
      if (state.currentSession) {
        state.messageHistory[state.currentSession.id] = [];
      }
    },
    
    setMessages: (state, action: PayloadAction<ChatMessage[]>) => {
      state.messages = action.payload;
      if (state.currentSession) {
        state.messageHistory[state.currentSession.id] = action.payload;
      }
    },
    
    // UI state management
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.isLoading = action.payload;
    },
    
    setTyping: (state, action: PayloadAction<boolean>) => {
      state.isTyping = action.payload;
    },
    
    setError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    
    clearError: (state) => {
      state.error = null;
    },
    
    // Settings management
    updateSettings: (state, action: PayloadAction<Partial<ChatState['settings']>>) => {
      state.settings = { ...state.settings, ...action.payload };
    },
    
    // Recovery actions
    retryFailedMessages: (state) => {
      state.pendingMessages.forEach(message => {
        if (message.status === 'failed') {
          message.status = 'sending';
        }
      });
    },
  },
  
  extraReducers: (builder) => {
    // Create session
    builder
      .addCase(createSession.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(createSession.fulfilled, (state, action) => {
        state.isLoading = false;
        state.currentSession = action.payload;
        state.sessions.push(action.payload);
      })
      .addCase(createSession.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || 'Failed to create session';
      });
    
    // Load session history
    builder
      .addCase(loadSessionHistory.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(loadSessionHistory.fulfilled, (state, action) => {
        state.isLoading = false;
        const messages = action.payload || [];
        state.messages = messages;
        if (state.currentSession) {
          state.messageHistory[state.currentSession.id] = messages;
        }
      })
      .addCase(loadSessionHistory.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || 'Failed to load session history';
      });
    
    // Update session context
    builder
      .addCase(updateSessionContext.pending, (state) => {
        state.error = null;
      })
      .addCase(updateSessionContext.fulfilled, (state, action) => {
        if (state.currentSession) {
          state.currentSession = { ...state.currentSession, ...action.payload };
        }
      })
      .addCase(updateSessionContext.rejected, (state, action) => {
        state.error = action.error.message || 'Failed to update session context';
      });
  },
});

export const {
  setCurrentSession,
  clearCurrentSession,
  updateSessionContext: updateSessionContextAction,
  addSession,
  addMessage,
  addOptimisticMessage,
  updateMessageStatus,
  removeFailedMessage,
  clearMessages,
  setMessages,
  setLoading,
  setTyping,
  setError,
  clearError,
  updateSettings,
  retryFailedMessages,
} = chatSlice.actions;

export default chatSlice.reducer;