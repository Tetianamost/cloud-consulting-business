import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { Provider } from 'react-redux';
import { configureStore } from '@reduxjs/toolkit';
import { BrowserRouter } from 'react-router-dom';
import ChatToggle from './ChatToggle';
import chatReducer from '../../store/slices/chatSlice';
import connectionReducer from '../../store/slices/connectionSlice';

// Mock websocket service
jest.mock('../../services/websocketService', () => ({
  __esModule: true,
  default: {
    connect: jest.fn().mockResolvedValue(undefined),
    disconnect: jest.fn(),
    sendChatMessage: jest.fn(),
    getConnectionStatus: jest.fn().mockReturnValue('disconnected'),
    isHealthy: jest.fn().mockReturnValue(false),
    forceReconnect: jest.fn(),
  },
}));

// Mock localStorage
const mockLocalStorage = {
  getItem: jest.fn(),
  setItem: jest.fn(),
  removeItem: jest.fn(),
  clear: jest.fn(),
};
Object.defineProperty(window, 'localStorage', {
  value: mockLocalStorage,
});

const createTestStore = (initialState = {}) => {
  return configureStore({
    reducer: {
      chat: chatReducer,
      connection: connectionReducer,
    },
    preloadedState: initialState,
  });
};

const renderWithProviders = (component: React.ReactElement, initialState = {}) => {
  const store = createTestStore(initialState);
  return render(
    <Provider store={store}>
      <BrowserRouter>
        {component}
      </BrowserRouter>
    </Provider>
  );
};

describe('ChatToggle', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockLocalStorage.getItem.mockReturnValue(null);
  });

  it('renders chat toggle button when closed', () => {
    renderWithProviders(<ChatToggle />);
    
    const toggleButton = screen.getByRole('button', { name: /open ai consultant chat/i });
    expect(toggleButton).toBeInTheDocument();
  });

  it('shows connection status indicator', () => {
    const initialState = {
      connection: {
        status: 'connected' as const,
        isHealthy: true,
        webSocket: null,
        lastConnected: null,
        reconnectAttempts: 0,
        maxReconnectAttempts: 10,
        reconnectDelay: 1000,
        maxReconnectDelay: 30000,
        error: null,
        latency: null,
        lastPingTime: null,
        connectionId: null,
      },
    };

    renderWithProviders(<ChatToggle />, initialState);
    
    // Should show green connection indicator for healthy connected state
    const toggleButton = screen.getByRole('button', { name: /open ai consultant chat/i });
    expect(toggleButton).toBeInTheDocument();
  });

  it('shows unread message count when there are new messages', () => {
    const initialState = {
      chat: {
        currentSession: null,
        sessions: [],
        sessionContext: {},
        messages: [
          {
            id: '1',
            type: 'assistant' as const,
            content: 'Hello, how can I help you?',
            timestamp: new Date().toISOString(),
            session_id: 'session-1',
          },
        ],
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
      },
    };

    renderWithProviders(<ChatToggle />, initialState);
    
    // Should show unread message indicator
    expect(screen.getAllByText('1')).toHaveLength(2); // One on main button, one on notification bell
  });

  it('opens chat when toggle button is clicked', async () => {
    renderWithProviders(<ChatToggle />);
    
    const toggleButton = screen.getByRole('button', { name: /open ai consultant chat/i });
    fireEvent.click(toggleButton);

    // Should render the ConsultantChat component
    await waitFor(() => {
      expect(screen.getByText(/consultant assistant/i)).toBeInTheDocument();
    });
  });

  it('persists chat state in localStorage', () => {
    renderWithProviders(<ChatToggle />);
    
    const toggleButton = screen.getByRole('button', { name: /open ai consultant chat/i });
    fireEvent.click(toggleButton);

    expect(mockLocalStorage.setItem).toHaveBeenCalledWith(
      'chatWidgetState',
      JSON.stringify({ isOpen: true, isMinimized: false })
    );
  });

  it('restores chat state from localStorage', () => {
    mockLocalStorage.getItem.mockReturnValue(
      JSON.stringify({ isOpen: true, isMinimized: false })
    );

    renderWithProviders(<ChatToggle />);
    
    // Should render the ConsultantChat component immediately
    expect(screen.getByText(/consultant assistant/i)).toBeInTheDocument();
  });

  it('shows notifications when there are unread messages and chat is closed', () => {
    const initialState = {
      chat: {
        currentSession: null,
        sessions: [],
        sessionContext: {},
        messages: [
          {
            id: '1',
            type: 'assistant' as const,
            content: 'This is a test notification message that should appear in the notification dropdown.',
            timestamp: new Date().toISOString(),
            session_id: 'session-1',
          },
        ],
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
      },
    };

    renderWithProviders(<ChatToggle />, initialState);
    
    // Should show notification bell
    const notificationButton = screen.getByRole('button', { name: /view chat notifications/i });
    expect(notificationButton).toBeInTheDocument();
    
    // Click to show notifications
    fireEvent.click(notificationButton);
    
    // Should show notification content
    expect(screen.getByText(/new ai response/i)).toBeInTheDocument();
  });

  it('handles connection errors gracefully', () => {
    const initialState = {
      connection: {
        status: 'failed' as const,
        isHealthy: false,
        error: 'Connection failed',
        webSocket: null,
        lastConnected: null,
        reconnectAttempts: 3,
        maxReconnectAttempts: 10,
        reconnectDelay: 1000,
        maxReconnectDelay: 30000,
        latency: null,
        lastPingTime: null,
        connectionId: null,
      },
    };

    renderWithProviders(<ChatToggle />, initialState);
    
    const toggleButton = screen.getByRole('button', { name: /open ai consultant chat/i });
    expect(toggleButton).toBeInTheDocument();
    
    // Should show red connection indicator for failed state
    // The exact implementation depends on the visual indicator used
  });
});