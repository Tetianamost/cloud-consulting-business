import React from 'react';
import { render, screen, fireEvent, waitFor, act } from '@testing-library/react';
import { Provider } from 'react-redux';
import { configureStore } from '@reduxjs/toolkit';
import '@testing-library/jest-dom';

import { ConsultantChat } from './ConsultantChat';
import chatReducer from '../../store/slices/chatSlice';
import connectionReducer from '../../store/slices/connectionSlice';

// Mock the enhanced AI service
jest.mock('../../services/simpleAIService', () => ({
  __esModule: true,
  default: {
    checkConnection: jest.fn().mockResolvedValue(true),
    isHealthy: jest.fn().mockReturnValue(true),
    sendMessage: jest.fn().mockResolvedValue({
      content: 'Mock AI response',
      timestamp: new Date().toISOString(),
    }),
    getMessages: jest.fn().mockResolvedValue([]),
    getSessionId: jest.fn().mockReturnValue('test-session-id'),
    resetSession: jest.fn(),
    forceReconnect: jest.fn().mockResolvedValue(true),
  },
}));

// Mock localStorage
const mockLocalStorage = {
  getItem: jest.fn(() => 'mock-admin-token'),
  setItem: jest.fn(),
  removeItem: jest.fn(),
  clear: jest.fn(),
};

Object.defineProperty(window, 'localStorage', {
  value: mockLocalStorage,
});

// Mock fetch for API calls
global.fetch = jest.fn().mockResolvedValue({
  ok: true,
  json: jest.fn().mockResolvedValue({
    success: true,
    messages: [],
  }),
});

// Create test store
const createTestStore = () => {
  return configureStore({
    reducer: {
      chat: chatReducer,
      connection: connectionReducer,
    },
  });
};

// Test wrapper component
const TestWrapper: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const store = createTestStore();
  return <Provider store={store}>{children}</Provider>;
};

describe('ConsultantChat', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockLocalStorage.getItem.mockReturnValue('mock-admin-token');
    
    // Reset AI service mocks to default values
    const mockAIService = require('../../services/simpleAIService').default;
    mockAIService.checkConnection.mockResolvedValue(true);
    mockAIService.isHealthy.mockReturnValue(true);
    mockAIService.sendMessage.mockResolvedValue({
      content: 'Mock AI response',
      timestamp: new Date().toISOString(),
    });
    mockAIService.getSessionId.mockReturnValue('test-session-id');
  });

  afterEach(() => {
    jest.clearAllTimers();
  });

  it('renders chat interface correctly', () => {
    render(
      <TestWrapper>
        <ConsultantChat />
      </TestWrapper>
    );

    expect(screen.getByText('AI Assistant')).toBeInTheDocument();
    expect(screen.getByPlaceholderText(/Ask about AWS services/)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /send/i })).toBeInTheDocument();
  });

  it('renders minimized state correctly', () => {
    const mockToggle = jest.fn();
    
    render(
      <TestWrapper>
        <ConsultantChat isMinimized={true} onToggleMinimize={mockToggle} />
      </TestWrapper>
    );

    const minimizedButton = screen.getByRole('button');
    expect(minimizedButton).toBeInTheDocument();
    
    fireEvent.click(minimizedButton);
    expect(mockToggle).toHaveBeenCalled();
  });

  it('displays quick action buttons', () => {
    render(
      <TestWrapper>
        <ConsultantChat />
      </TestWrapper>
    );

    expect(screen.getByText('Cost Estimate')).toBeInTheDocument();
    expect(screen.getByText('Security Review')).toBeInTheDocument();
    expect(screen.getByText('Best Practices')).toBeInTheDocument();
    expect(screen.getByText('Alternatives')).toBeInTheDocument();
    expect(screen.getByText('Next Steps')).toBeInTheDocument();
  });

  it('displays quick action buttons', () => {
    render(
      <TestWrapper>
        <ConsultantChat />
      </TestWrapper>
    );

    expect(screen.getByText('Cost Estimate')).toBeInTheDocument();
    expect(screen.getByText('Security Review')).toBeInTheDocument();
    expect(screen.getByText('Best Practices')).toBeInTheDocument();
    expect(screen.getByText('Alternatives')).toBeInTheDocument();
    expect(screen.getByText('Next Steps')).toBeInTheDocument();
  });

  it('opens full AI assistant when maximize button is clicked', () => {
    // Mock window.open
    const mockOpen = jest.fn();
    Object.defineProperty(window, 'open', {
      value: mockOpen,
      writable: true,
    });

    render(
      <TestWrapper>
        <ConsultantChat />
      </TestWrapper>
    );

    const maximizeButton = screen.getByTitle('Open Full AI Assistant');
    fireEvent.click(maximizeButton);

    expect(mockOpen).toHaveBeenCalledWith('/admin/ai-consultant', '_blank');
  });

  it('sends message when form is submitted', async () => {
    const mockSendMessage = require('../../services/simpleAIService').default.sendMessage;
    
    render(
      <TestWrapper>
        <ConsultantChat />
      </TestWrapper>
    );

    // Wait for service to initialize
    await waitFor(() => {
      expect(screen.getByPlaceholderText(/Ask about AWS services/)).not.toBeDisabled();
    });

    const input = screen.getByPlaceholderText(/Ask about AWS services/);
    const sendButton = screen.getByRole('button', { name: /send/i });

    // Type message
    fireEvent.change(input, { target: { value: 'What is EC2?' } });
    expect(input).toHaveValue('What is EC2?');

    // Submit form
    fireEvent.click(sendButton);

    // Check that input is cleared
    expect(input).toHaveValue('');

    // Wait for user message to appear
    await waitFor(() => {
      expect(screen.getByText('What is EC2?')).toBeInTheDocument();
    });

    // Verify the AI service was called
    expect(mockSendMessage).toHaveBeenCalledWith({
      message: 'What is EC2?',
      context: {
        clientName: undefined,
        meetingType: undefined,
      },
    });

    // Wait for AI response
    await waitFor(() => {
      expect(screen.getByText('Mock AI response')).toBeInTheDocument();
    });
  });

  it('sends message when Enter key is pressed', async () => {
    render(
      <TestWrapper>
        <ConsultantChat />
      </TestWrapper>
    );

    // Wait for service to initialize
    await waitFor(() => {
      expect(screen.getByPlaceholderText(/Ask about AWS services/)).not.toBeDisabled();
    });

    const input = screen.getByPlaceholderText(/Ask about AWS services/);

    // Type message
    fireEvent.change(input, { target: { value: 'Test message' } });

    // Press Enter (this will trigger form submission)
    fireEvent.keyDown(input, { key: 'Enter', code: 'Enter' });

    // Wait for message to appear
    await waitFor(() => {
      expect(screen.getByText('Test message')).toBeInTheDocument();
    });
  });

  it('handles quick action clicks', async () => {
    const mockSendMessage = require('../../services/simpleAIService').default.sendMessage;
    
    render(
      <TestWrapper>
        <ConsultantChat />
      </TestWrapper>
    );

    // Wait for service to initialize
    await waitFor(() => {
      expect(screen.getByPlaceholderText(/Ask about AWS services/)).not.toBeDisabled();
    });

    const costEstimateButton = screen.getByText('Cost Estimate');
    fireEvent.click(costEstimateButton);

    // Wait for quick action message to appear
    await waitFor(() => {
      expect(screen.getByText('Provide a cost estimate for this solution')).toBeInTheDocument();
    });

    // Verify the AI service was called with the quick action prompt
    expect(mockSendMessage).toHaveBeenCalledWith({
      message: 'Provide a cost estimate for this solution',
      context: {
        clientName: undefined,
        meetingType: undefined,
      },
    });

    // Wait for AI response
    await waitFor(() => {
      expect(screen.getByText('Mock AI response')).toBeInTheDocument();
    });
  });

  it('shows loading state while waiting for response', async () => {
    render(
      <TestWrapper>
        <ConsultantChat />
      </TestWrapper>
    );

    // Wait for service to initialize
    await waitFor(() => {
      expect(screen.getByPlaceholderText(/Ask about AWS services/)).not.toBeDisabled();
    });

    const input = screen.getByPlaceholderText(/Ask about AWS services/);
    const sendButton = screen.getByRole('button', { name: /send/i });

    // Send message
    fireEvent.change(input, { target: { value: 'Test message' } });
    fireEvent.click(sendButton);

    // Check for loading indicator (animated dots)
    await waitFor(() => {
      const loadingDots = screen.getAllByText('', { selector: '.animate-bounce' });
      expect(loadingDots.length).toBeGreaterThan(0);
    });
  });

  it('displays connection status indicator', async () => {
    render(
      <TestWrapper>
        <ConsultantChat />
      </TestWrapper>
    );

    // Should show connected status
    await waitFor(() => {
      expect(screen.getByText('Connected')).toBeInTheDocument();
    });

    // Status indicator should be green
    const statusIndicator = document.querySelector('.w-2.h-2.rounded-full');
    expect(statusIndicator).toBeInTheDocument();
    expect(statusIndicator).toHaveClass('bg-green-400');
  });

  it('shows welcome message when no messages exist', () => {
    render(
      <TestWrapper>
        <ConsultantChat />
      </TestWrapper>
    );

    expect(screen.getByText(/Start a conversation or click to open full AI Assistant/)).toBeInTheDocument();
    expect(screen.getByText('Open Full AI Assistant')).toBeInTheDocument();
  });

  it('formats timestamps correctly', async () => {
    render(
      <TestWrapper>
        <ConsultantChat />
      </TestWrapper>
    );

    // Wait for service to initialize
    await waitFor(() => {
      expect(screen.getByPlaceholderText(/Ask about AWS services/)).not.toBeDisabled();
    });

    const input = screen.getByPlaceholderText(/Ask about AWS services/);
    fireEvent.change(input, { target: { value: 'Test message' } });
    fireEvent.submit(input.closest('form')!);

    // Wait for message with timestamp
    await waitFor(() => {
      const timeElements = screen.getAllByText(/\d{1,2}:\d{2}/);
      expect(timeElements.length).toBeGreaterThan(0);
    });
  });

  it('disables input and buttons when service is unhealthy', async () => {
    // Mock the AI service to be unhealthy
    const mockAIService = require('../../services/simpleAIService').default;
    mockAIService.checkConnection.mockResolvedValue(false);
    mockAIService.isHealthy.mockReturnValue(false);

    render(
      <TestWrapper>
        <ConsultantChat />
      </TestWrapper>
    );

    // Wait for service check to complete
    await waitFor(() => {
      expect(screen.getByText('Offline')).toBeInTheDocument();
    });

    // Should show connection lost warning
    expect(screen.getByText(/Connection lost. Trying to reconnect/)).toBeInTheDocument();
    
    // Status indicator should be red
    const statusIndicator = document.querySelector('.w-2.h-2.rounded-full');
    expect(statusIndicator).toHaveClass('bg-red-400');
  });

  it('handles AI service connection errors gracefully', async () => {
    const consoleSpy = jest.spyOn(console, 'error').mockImplementation();
    
    // Mock the AI service to throw an error
    const mockAIService = require('../../services/simpleAIService').default;
    mockAIService.checkConnection.mockRejectedValue(new Error('Connection failed'));

    render(
      <TestWrapper>
        <ConsultantChat />
      </TestWrapper>
    );

    // Wait for error to be handled
    await waitFor(() => {
      expect(consoleSpy).toHaveBeenCalledWith('[ConsultantChat] Failed to initialize AI service:', expect.any(Error));
    });

    consoleSpy.mockRestore();
  });

  it('attempts to reconnect when service becomes unhealthy', async () => {
    jest.useFakeTimers();
    
    const mockAIService = require('../../services/simpleAIService').default;
    // Initially healthy, then becomes unhealthy
    mockAIService.isHealthy.mockReturnValueOnce(true).mockReturnValue(false);
    mockAIService.forceReconnect.mockResolvedValue(true);

    render(
      <TestWrapper>
        <ConsultantChat />
      </TestWrapper>
    );

    // Fast-forward time to trigger reconnection check
    act(() => {
      jest.advanceTimersByTime(30000);
    });

    // Wait for reconnection attempt
    await waitFor(() => {
      expect(mockAIService.forceReconnect).toHaveBeenCalled();
    });

    jest.useRealTimers();
  });

  it('handles AI service errors gracefully', async () => {
    const consoleSpy = jest.spyOn(console, 'error').mockImplementation();
    
    // Mock the AI service to throw an error when sending message
    const mockAIService = require('../../services/simpleAIService').default;
    mockAIService.sendMessage.mockRejectedValue(new Error('Service error'));

    render(
      <TestWrapper>
        <ConsultantChat />
      </TestWrapper>
    );

    // Wait for service to initialize
    await waitFor(() => {
      expect(screen.getByPlaceholderText(/Ask about AWS services/)).not.toBeDisabled();
    });

    // Send message to trigger error
    const input = screen.getByPlaceholderText(/Ask about AWS services/);
    fireEvent.change(input, { target: { value: 'Test' } });
    fireEvent.submit(input.closest('form')!);

    // Wait for error to be logged
    await waitFor(() => {
      expect(consoleSpy).toHaveBeenCalledWith('Failed to send message:', expect.any(Error));
    });

    consoleSpy.mockRestore();
  });

  it('handles close button click', () => {
    const mockClose = jest.fn();

    render(
      <TestWrapper>
        <ConsultantChat onClose={mockClose} />
      </TestWrapper>
    );

    const closeButton = screen.getByRole('button', { name: /x/i });
    fireEvent.click(closeButton);

    expect(mockClose).toHaveBeenCalled();
  });

  it('prevents sending empty messages', async () => {
    const mockSendMessage = require('../../services/simpleAIService').default.sendMessage;
    
    render(
      <TestWrapper>
        <ConsultantChat />
      </TestWrapper>
    );

    // Wait for service to initialize
    await waitFor(() => {
      expect(screen.getByPlaceholderText(/Ask about AWS services/)).not.toBeDisabled();
    });

    const sendButton = screen.getByRole('button', { name: /send/i });

    // Try to send empty message
    fireEvent.click(sendButton);

    // Verify no message was sent to the service
    expect(mockSendMessage).not.toHaveBeenCalled();
  });

  it('handles retry button click when connection is lost', async () => {
    // Mock the AI service to be initially unhealthy
    const mockAIService = require('../../services/simpleAIService').default;
    mockAIService.checkConnection.mockResolvedValueOnce(false);
    mockAIService.isHealthy.mockReturnValueOnce(false);
    mockAIService.forceReconnect.mockResolvedValue(true);

    render(
      <TestWrapper>
        <ConsultantChat />
      </TestWrapper>
    );

    // Wait for connection lost message
    await waitFor(() => {
      expect(screen.getByText(/Connection lost. Trying to reconnect/)).toBeInTheDocument();
    });

    // Click retry button
    const retryButton = screen.getByText('Retry');
    fireEvent.click(retryButton);

    // Verify forceReconnect was called
    expect(mockAIService.forceReconnect).toHaveBeenCalled();
  });

  it('displays user and assistant message types correctly', async () => {
    render(
      <TestWrapper>
        <ConsultantChat />
      </TestWrapper>
    );

    // Wait for service to initialize
    await waitFor(() => {
      expect(screen.getByPlaceholderText(/Ask about AWS services/)).not.toBeDisabled();
    });

    const input = screen.getByPlaceholderText(/Ask about AWS services/);
    fireEvent.change(input, { target: { value: 'Test question' } });
    fireEvent.submit(input.closest('form')!);

    // Wait for both user and assistant messages
    await waitFor(() => {
      expect(screen.getByText('Test question')).toBeInTheDocument();
    });

    await waitFor(() => {
      expect(screen.getByText('Mock AI response')).toBeInTheDocument();
    });

    // Check that messages have correct styling
    const userMessage = screen.getByText('Test question').closest('div');
    const assistantMessage = screen.getByText('Mock AI response').closest('div');

    expect(userMessage).toHaveClass('justify-end');
    expect(assistantMessage).toHaveClass('justify-start');
  });
});

// Integration tests with Redux store
describe('ConsultantChat Redux Integration', () => {
  it('integrates with Redux store correctly', () => {
    const store = createTestStore();

    render(
      <Provider store={store}>
        <ConsultantChat />
      </Provider>
    );

    // Verify component renders with store
    expect(screen.getByText('Consultant Assistant')).toBeInTheDocument();

    // Check initial store state
    const state = store.getState();
    expect(state.chat.messages).toEqual([]);
    expect(state.chat.isLoading).toBe(false);
  });
});

// Accessibility tests
describe('ConsultantChat Accessibility', () => {
  it('has proper ARIA labels and roles', () => {
    render(
      <TestWrapper>
        <ConsultantChat />
      </TestWrapper>
    );

    // Check for proper form elements
    const input = screen.getByPlaceholderText(/Ask about AWS services/);
    expect(input).toHaveAttribute('type', 'text');

    const sendButton = screen.getByRole('button', { name: /send/i });
    expect(sendButton).toHaveAttribute('type', 'submit');
  });

  it('supports keyboard navigation', () => {
    render(
      <TestWrapper>
        <ConsultantChat />
      </TestWrapper>
    );

    const input = screen.getByPlaceholderText(/Ask about AWS services/);
    
    // Focus should work
    input.focus();
    expect(document.activeElement).toBe(input);

    // Tab navigation should work
    fireEvent.keyDown(input, { key: 'Tab' });
    expect(document.activeElement).not.toBe(input);
  });
});