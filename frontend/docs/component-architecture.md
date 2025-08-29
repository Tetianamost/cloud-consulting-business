# Frontend Component Architecture

## Overview

The Cloud Consulting Platform frontend is built with React 19+ and TypeScript, following modern component architecture patterns with Redux Toolkit for state management. This document outlines the component structure, design patterns, and architectural decisions.

## Architecture Principles

### 1. Component Hierarchy
- **Atomic Design**: Components organized by complexity (atoms, molecules, organisms)
- **Feature-Based Structure**: Components grouped by functionality
- **Reusability**: Shared components in common directories
- **Separation of Concerns**: Clear separation between UI, logic, and data

### 2. State Management
- **Redux Toolkit**: Global state management
- **Local State**: Component-specific state with useState
- **Context API**: Authentication and theme management
- **Custom Hooks**: Reusable stateful logic

### 3. Type Safety
- **TypeScript**: Full type coverage
- **Interface Definitions**: Comprehensive type definitions
- **Prop Validation**: Runtime and compile-time validation
- **API Types**: Shared types between frontend and backend

## Directory Structure

```
frontend/src/
├── components/
│   ├── admin/              # Admin-specific components
│   │   ├── Login.tsx       # Authentication component
│   │   ├── IntegratedAdminDashboard.tsx
│   │   ├── AIConsultantPage.tsx  # Advanced AI consultant interface
│   │   ├── ChatToggle.tsx  # Chat functionality toggle
│   │   ├── SimpleWorkingChat.tsx
│   │   ├── ConnectionStatus.tsx
│   │   └── ...
│   ├── layout/             # Layout components
│   ├── sections/           # Page sections
│   └── ui/                 # Base UI components
├── contexts/               # React contexts
│   ├── AuthContext.tsx     # Authentication context
│   └── ThemeContext.tsx    # Theme management
├── hooks/                  # Custom hooks
│   ├── useAuth.ts          # Authentication hook
│   ├── useChat.ts          # Chat functionality
│   └── usePaginatedMessages.ts
├── services/               # API and external services
│   ├── api.ts              # Main API service
│   ├── chatModeManager.ts  # Chat mode management
│   ├── pollingChatService.ts
│   └── websocketService.ts
├── store/                  # Redux store
│   ├── index.ts            # Store configuration
│   └── slices/             # Redux slices
│       ├── authSlice.ts
│       ├── chatSlice.ts
│       └── connectionSlice.ts
├── types/                  # TypeScript definitions
├── utils/                  # Utility functions
└── styles/                 # Global styles and themes
```

## Core Components

### 1. Authentication System

#### Login Component
```typescript
// frontend/src/components/admin/Login.tsx
interface LoginProps {
  onSuccess?: () => void;
  redirectTo?: string;
}

const Login: React.FC<LoginProps> = ({ onSuccess, redirectTo }) => {
  // Component implementation with styled-components
  // Form validation and error handling
  // Integration with AuthContext
};
```

**Features**:
- Styled-components for consistent theming
- Form validation with error states
- Loading states during authentication
- Responsive design with mobile support
- Demo credentials display

#### AuthContext
```typescript
// frontend/src/contexts/AuthContext.tsx
interface AuthContextType {
  isAuthenticated: boolean;
  login: (username: string, password: string) => Promise<boolean>;
  logout: () => void;
  loading: boolean;
}

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  // JWT token management
  // Clean error handling without debug logging
  // Secure logout functionality
};
```

### 2. Dashboard Architecture

#### IntegratedAdminDashboard
```typescript
// frontend/src/components/admin/IntegratedAdminDashboard.tsx
interface IntegratedAdminDashboardProps {
  children?: React.ReactNode;
}

export const IntegratedAdminDashboard: React.FC<IntegratedAdminDashboardProps> = ({ children }) => {
  // Main dashboard container
  // Route management
  // Chat service initialization
  // Connection status monitoring
};
```

**Key Features**:
- Sidebar navigation with responsive design
- Header with chat status and controls
- Route-based content rendering
- Real-time connection monitoring
- Chat service integration

#### Component Layout
```
IntegratedAdminDashboard
├── AdminSidebar
│   ├── Navigation items
│   └── User profile section
├── Header
│   ├── Dashboard title
│   ├── Chat status indicator
│   ├── ConnectionStatus component
│   └── User controls
├── Main Content Area
│   ├── Routes (React Router)
│   │   ├── Dashboard (metrics)
│   │   ├── Inquiries management
│   │   ├── Reports management
│   │   ├── Chat interface
│   │   └── Settings
│   └── Loading states
└── Footer
    ├── System status
    └── Connection information
```

### 3. Chat System Components

#### AIConsultantPage (Primary Chat Interface)
```typescript
// frontend/src/components/admin/AIConsultantPage.tsx
interface AIConsultantPageProps {
  // No props - fully self-contained component
}

export const AIConsultantPage: React.FC = () => {
  // Advanced AI consultant interface with Redux integration
  // Context-aware conversations with client and meeting context
  // 8 pre-defined quick action buttons for common scenarios
  // Fullscreen mode support with toggle functionality
  // Connection health monitoring and status display
  // Debounced input for performance optimization
  // Session persistence across page reloads
};
```

**Key Features**:
- **Quick Actions System**: 8 pre-defined prompts for common consulting scenarios:
  - Cost Estimate, Security Review, Best Practices, Alternatives
  - Next Steps, Migration Plan, Compliance, Performance
- **Context Management**: Client name and meeting context for personalized AI responses
- **Connection Health**: Real-time monitoring with visual status indicators
- **Fullscreen Mode**: Distraction-free chat experience with maximize/minimize toggle
- **Session Persistence**: Maintains conversation history across page reloads using Redux
- **Debounced Input**: 300ms delay to prevent excessive API calls during typing
- **Settings Panel**: Configurable client context and connection testing
- **Professional UI**: Clean, consultant-focused interface with proper message formatting

#### SimpleWorkingChat
```typescript
// frontend/src/components/admin/SimpleWorkingChat.tsx
interface SimpleWorkingChatProps {
  sessionId?: string;
  onSessionChange?: (session: ChatSession) => void;
}

export const SimpleWorkingChat: React.FC<SimpleWorkingChatProps> = ({ sessionId, onSessionChange }) => {
  // Full-featured chat interface
  // Message history with virtual scrolling
  // Real-time message updates
  // WebSocket with polling fallback
};
```

**Features**:
- Real-time messaging with AI assistant
- Message history with pagination
- Typing indicators and presence
- File upload support (planned)
- Message status tracking

#### ChatToggle
```typescript
// frontend/src/components/admin/ChatToggle.tsx
interface ChatToggleProps {
  variant?: 'primary' | 'secondary';
  size?: 'small' | 'medium' | 'large';
  showStatus?: boolean;
  onToggle?: (enabled: boolean) => void;
  className?: string;
}

export const ChatToggle: React.FC<ChatToggleProps> = ({ variant, size, showStatus, onToggle, className }) => {
  // Chat functionality toggle
  // Connection status display
  // Visual feedback for state changes
};
```

#### ConnectionStatus
```typescript
// frontend/src/components/admin/ConnectionStatus.tsx
interface ConnectionStatusProps {
  showDetails?: boolean;
  onReconnect?: () => void;
  className?: string;
}

export const ConnectionStatus: React.FC<ConnectionStatusProps> = ({ showDetails, onReconnect, className }) => {
  // Real-time connection monitoring
  // Visual status indicators
  // Reconnection controls
  // Error state handling
};
```

## State Management

### Redux Store Structure

```typescript
// frontend/src/store/index.ts
interface RootState {
  auth: AuthState;
  chat: ChatState;
  connection: ConnectionState;
  inquiries: InquiryState;
  reports: ReportState;
  metrics: MetricsState;
}
```

### Chat Slice
```typescript
// frontend/src/store/slices/chatSlice.ts
interface ChatState {
  currentSession: ChatSession | null;
  messages: ChatMessage[];
  isLoading: boolean;
  error: string | null;
  typingUsers: string[];
  messageHistory: Record<string, ChatMessage[]>;
}

const chatSlice = createSlice({
  name: 'chat',
  initialState,
  reducers: {
    messageReceived: (state, action) => {
      state.messages.push(action.payload);
    },
    messageSent: (state, action) => {
      // Optimistic update
      state.messages.push({ ...action.payload, status: 'sending' });
    },
    sessionCreated: (state, action) => {
      state.currentSession = action.payload;
    },
    // ... other reducers
  },
});
```

### Connection Slice
```typescript
// frontend/src/store/slices/connectionSlice.ts
interface ConnectionState {
  status: 'disconnected' | 'connecting' | 'connected' | 'error';
  mode: 'websocket' | 'polling';
  lastConnected: Date | null;
  error: string | null;
  reconnectAttempts: number;
}
```

## Service Layer

### API Service
```typescript
// frontend/src/services/api.ts
class ApiService {
  private baseUrl: string;
  private authToken: string | null = null;

  // Authentication methods
  async login(username: string, password: string): Promise<LoginResponse>;
  setAuthToken(token: string | null): void;
  logout(): void;

  // Admin endpoints
  async listInquiries(filters?: InquiryFilters): Promise<AdminInquiriesResponse>;
  async getSystemMetrics(): Promise<SystemMetrics>;
  async downloadReport(inquiryId: string, format: 'pdf' | 'html'): Promise<Blob>;

  // Chat endpoints
  async createChatSession(): Promise<ChatSession>;
  async getChatHistory(sessionId: string): Promise<ChatMessage[]>;

  // Private methods
  private async request<T>(endpoint: string, options?: RequestInit): Promise<T>;
}
```

### Chat Mode Manager
```typescript
// frontend/src/services/chatModeManager.ts
class ChatModeManager {
  private currentMode: 'websocket' | 'polling' = 'websocket';
  private websocketService: WebSocketService;
  private pollingService: PollingChatService;

  async initializeChatService(): Promise<void>;
  async switchMode(mode: 'websocket' | 'polling'): Promise<void>;
  forceReconnect(): void;
  cleanup(): void;

  // Event handling
  private handleConnectionFailure(): void;
  private handleModeSwitch(): void;
}
```

## AI Consultant Page Features

### Quick Actions System
The AI Consultant Page includes pre-defined quick actions for common consulting scenarios:

```typescript
const QUICK_ACTIONS = [
  { id: 'cost-estimate', label: 'Cost Estimate', prompt: 'Provide a cost estimate for this solution' },
  { id: 'security-review', label: 'Security Review', prompt: 'What are the security considerations for this approach?' },
  { id: 'best-practices', label: 'Best Practices', prompt: 'What are the AWS best practices for this scenario?' },
  { id: 'alternatives', label: 'Alternatives', prompt: 'What are alternative approaches to consider?' },
  { id: 'next-steps', label: 'Next Steps', prompt: 'What are the recommended next steps?' },
  { id: 'migration-plan', label: 'Migration Plan', prompt: 'What is the recommended migration approach?' },
  { id: 'compliance', label: 'Compliance', prompt: 'What compliance considerations should we address?' },
  { id: 'performance', label: 'Performance', prompt: 'How can we optimize performance for this solution?' },
];
```

### Context Management
The component supports contextual conversations through:

- **Client Name**: Personalizes responses for specific clients
- **Meeting Context**: Provides situational awareness (e.g., "Migration planning", "Cost optimization")
- **Session Context**: Maintains conversation continuity across interactions

### Connection Management
Advanced connection handling with multiple modes:

```typescript
// Connection modes
type ConnectionMode = 'websocket' | 'polling' | 'auto';

// Automatic fallback logic
const chatModeManager = {
  getCurrentMode(): ConnectionMode,
  switchMode(mode: ConnectionMode): Promise<void>,
  isHealthy(): boolean,
  getStatusMessage(): string,
};
```

### User Interface Features

#### Fullscreen Mode
- Toggle between embedded and fullscreen views
- Optimized for focused conversations
- Maintains all functionality in both modes

#### Settings Panel
- Client name configuration
- Meeting context setup
- Connection mode selection
- Real-time status monitoring

#### Message Display
- User and AI message differentiation
- Timestamp formatting
- Loading states with animated indicators
- Error handling with retry options

## Custom Hooks

### useDebouncedInput Hook
```typescript
// frontend/src/hooks/useDebouncedInput.ts
interface UseDebouncedInputReturn {
  value: string;
  isDebouncing: boolean;
  setValue: (value: string) => void;
  clearValue: () => void;
}

export const useDebouncedInput = (options: {
  delay: number;
  minLength: number;
}): [UseDebouncedInputReturn, UseDebouncedInputActions] => {
  // Debounced input handling for chat messages
  // Prevents excessive API calls during typing
  // Provides visual feedback during debounce period
};
```

### useAuth Hook
```typescript
// frontend/src/hooks/useAuth.ts
export const useAuth = () => {
  const context = useContext(AuthContext);
  
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  
  return context;
};
```

### useChat Hook
```typescript
// frontend/src/hooks/useChat.ts
interface UseChatReturn {
  messages: ChatMessage[];
  sendMessage: (content: string) => Promise<void>;
  connectionStatus: ConnectionStatus;
  currentSession: ChatSession | null;
  isLoading: boolean;
  error: string | null;
}

export const useChat = (sessionId?: string): UseChatReturn => {
  // Chat functionality hook
  // Message management
  // Connection handling
  // Error management
};
```

### usePaginatedMessages Hook
```typescript
// frontend/src/hooks/usePaginatedMessages.ts
interface UsePaginatedMessagesReturn {
  messages: ChatMessage[];
  loadMore: () => void;
  hasMore: boolean;
  isLoading: boolean;
}

export const usePaginatedMessages = (sessionId: string): UsePaginatedMessagesReturn => {
  // Paginated message loading
  // Virtual scrolling support
  // Performance optimization
};
```

## Design Patterns

### 1. Component Composition
```typescript
// Composable components with render props
interface ChatContainerProps {
  children: (props: ChatRenderProps) => React.ReactNode;
}

const ChatContainer: React.FC<ChatContainerProps> = ({ children }) => {
  const chatProps = useChat();
  return <div className="chat-container">{children(chatProps)}</div>;
};

// Usage
<ChatContainer>
  {({ messages, sendMessage, connectionStatus }) => (
    <div>
      <MessageList messages={messages} />
      <MessageInput onSend={sendMessage} disabled={connectionStatus !== 'connected'} />
    </div>
  )}
</ChatContainer>
```

### 2. Higher-Order Components
```typescript
// Authentication HOC
function withAuth<P extends object>(Component: React.ComponentType<P>) {
  return function AuthenticatedComponent(props: P) {
    const { isAuthenticated, loading } = useAuth();
    
    if (loading) return <LoadingSpinner />;
    if (!isAuthenticated) return <Navigate to="/admin/login" />;
    
    return <Component {...props} />;
  };
}

// Usage
export default withAuth(AdminDashboard);
```

### 3. Custom Hook Patterns
```typescript
// Compound hook pattern
export const useChatWithConnection = (sessionId?: string) => {
  const chat = useChat(sessionId);
  const connection = useConnection();
  
  return {
    ...chat,
    connectionStatus: connection.status,
    reconnect: connection.reconnect,
  };
};
```

## Performance Optimization

### 1. React.memo Usage
```typescript
// Memoized components for performance
const MessageItem = React.memo<MessageItemProps>(({ message, isOwn }) => {
  return (
    <div className={`message ${isOwn ? 'own' : 'other'}`}>
      <div className="content">{message.content}</div>
      <div className="timestamp">{formatTime(message.createdAt)}</div>
    </div>
  );
});
```

### 2. Virtual Scrolling
```typescript
// Virtual scrolling for large message lists
const VirtualizedMessageList: React.FC<VirtualizedMessageListProps> = ({ messages }) => {
  const { virtualItems, totalSize, scrollElementRef } = useVirtualizer({
    count: messages.length,
    getScrollElement: () => scrollElementRef.current,
    estimateSize: () => 60,
  });

  return (
    <div ref={scrollElementRef} className="message-list">
      <div style={{ height: totalSize }}>
        {virtualItems.map((virtualItem) => (
          <MessageItem
            key={virtualItem.index}
            message={messages[virtualItem.index]}
            style={{
              position: 'absolute',
              top: 0,
              left: 0,
              width: '100%',
              transform: `translateY(${virtualItem.start}px)`,
            }}
          />
        ))}
      </div>
    </div>
  );
};
```

### 3. Code Splitting
```typescript
// Lazy loading for route components
const AdminDashboard = React.lazy(() => import('./components/admin/IntegratedAdminDashboard'));
const ChatPage = React.lazy(() => import('./components/admin/ChatPage'));

// Usage with Suspense
<Suspense fallback={<LoadingSpinner />}>
  <Routes>
    <Route path="/admin/dashboard" element={<AdminDashboard />} />
    <Route path="/admin/chat" element={<ChatPage />} />
  </Routes>
</Suspense>
```

## Testing Patterns

### 1. Component Testing
```typescript
// Component testing with React Testing Library
describe('Login Component', () => {
  it('should handle successful login', async () => {
    const mockLogin = jest.fn().mockResolvedValue(true);
    render(
      <AuthProvider value={{ login: mockLogin, isAuthenticated: false, loading: false, logout: jest.fn() }}>
        <Login />
      </AuthProvider>
    );

    const usernameInput = screen.getByLabelText(/username/i);
    const passwordInput = screen.getByLabelText(/password/i);
    const submitButton = screen.getByRole('button', { name: /sign in/i });

    fireEvent.change(usernameInput, { target: { value: 'admin' } });
    fireEvent.change(passwordInput, { target: { value: 'cloudadmin' } });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(mockLogin).toHaveBeenCalledWith('admin', 'cloudadmin');
    });
  });
});
```

### 2. Hook Testing
```typescript
// Custom hook testing
describe('useChat Hook', () => {
  it('should send message and update state', async () => {
    const { result } = renderHook(() => useChat('session-123'));

    act(() => {
      result.current.sendMessage('Hello');
    });

    expect(result.current.messages).toHaveLength(1);
    expect(result.current.messages[0].content).toBe('Hello');
  });
});
```

## Error Boundaries

### Global Error Boundary
```typescript
// Global error boundary for the application
class ErrorBoundary extends React.Component<ErrorBoundaryProps, ErrorBoundaryState> {
  constructor(props: ErrorBoundaryProps) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error('Error caught by boundary:', error, errorInfo);
    // Send to error reporting service
  }

  render() {
    if (this.state.hasError) {
      return <ErrorFallback error={this.state.error} />;
    }

    return this.props.children;
  }
}
```

## Accessibility

### ARIA Support
```typescript
// Accessible components with proper ARIA attributes
const ChatMessage: React.FC<ChatMessageProps> = ({ message, isOwn }) => {
  return (
    <div
      role="article"
      aria-label={`Message from ${isOwn ? 'you' : 'assistant'}`}
      className={`message ${isOwn ? 'own' : 'other'}`}
    >
      <div className="content" aria-describedby={`timestamp-${message.id}`}>
        {message.content}
      </div>
      <div id={`timestamp-${message.id}`} className="timestamp">
        {formatTime(message.createdAt)}
      </div>
    </div>
  );
};
```

### Keyboard Navigation
```typescript
// Keyboard navigation support
const MessageInput: React.FC<MessageInputProps> = ({ onSend, disabled }) => {
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  return (
    <textarea
      onKeyDown={handleKeyDown}
      disabled={disabled}
      aria-label="Type your message"
      placeholder="Type a message..."
    />
  );
};
```

## Code Quality Standards

### Logging Best Practices
- **Production Ready**: Clean console output without debug logging
- **Error Handling**: Structured error logging with context
- **Security**: No sensitive information in logs
- **Performance**: Minimal logging overhead in production

See [Logging Best Practices](./logging-best-practices.md) for detailed guidelines.

## Future Enhancements

### Planned Improvements
1. **Server-Side Rendering (SSR)**: Next.js migration for better SEO
2. **Progressive Web App (PWA)**: Offline support and mobile app features
3. **Micro-frontends**: Modular architecture for scalability
4. **Advanced State Management**: Zustand or Valtio for simpler state management
5. **Component Library**: Standalone component library with Storybook

### Performance Optimizations
1. **Bundle Splitting**: More granular code splitting
2. **Tree Shaking**: Better dead code elimination
3. **Image Optimization**: WebP support and lazy loading
4. **Caching Strategies**: Service worker implementation
5. **Memory Management**: Better cleanup and garbage collection

This component architecture provides a solid foundation for the Cloud Consulting Platform frontend, with clear separation of concerns, type safety, performance optimization, and production-ready code quality standards built in from the ground up.