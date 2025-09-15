import { configureStore } from '@reduxjs/toolkit';
import chatReducer, {
  setCurrentSession,
  addMessage,
  clearMessages,
  ChatMessage,
  ChatSession,
} from './chatSlice';

// Create a properly typed store for testing
const createTestStore = () => configureStore({
  reducer: {
    chat: chatReducer,
  },
});

type TestStore = ReturnType<typeof createTestStore>;

describe('chatSlice - Essential Tests', () => {
  let store: TestStore;

  const getChatState = () => store.getState().chat;

  beforeEach(() => {
    store = createTestStore();
  });

  describe('initial state', () => {
    it('should have correct initial state', () => {
      const state = getChatState();
      
      expect(state.currentSession).toBeNull();
      expect(state.sessions).toEqual([]);
      expect(state.messages).toEqual([]);
      expect(state.isLoading).toBe(false);
      expect(state.error).toBeNull();
    });
  });

  describe('session management', () => {
    it('should set current session', () => {
      const mockSession: ChatSession = {
        id: 'session-1',
        consultant_id: 'consultant-1',
        client_name: 'Test Client',
        status: 'active',
        created_at: '2023-01-01T00:00:00Z',
        updated_at: '2023-01-01T00:00:00Z',
        last_activity: '2023-01-01T00:00:00Z',
      };

      store.dispatch(setCurrentSession(mockSession));
      
      const state = getChatState();
      expect(state.currentSession).toEqual(mockSession);
    });
  });

  describe('message management', () => {
    it('should add message', () => {
      const mockMessage: ChatMessage = {
        id: 'msg-1',
        session_id: 'session-1',
        type: 'user',
        content: 'Hello',
        timestamp: '2023-01-01T00:00:00Z',
      };

      store.dispatch(addMessage(mockMessage));
      
      const state = getChatState();
      expect(state.messages).toContain(mockMessage);
    });

    it('should clear messages', () => {
      const mockMessage: ChatMessage = {
        id: 'msg-1',
        session_id: 'session-1',
        type: 'user',
        content: 'Hello',
        timestamp: '2023-01-01T00:00:00Z',
      };

      store.dispatch(addMessage(mockMessage));
      store.dispatch(clearMessages());
      
      const state = getChatState();
      expect(state.messages).toEqual([]);
    });
  });
});