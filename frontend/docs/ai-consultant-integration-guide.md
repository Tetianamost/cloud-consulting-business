# AI Consultant Integration Guide

## Overview

This guide provides comprehensive instructions for integrating with the AI Consultant system, including frontend components, backend services, and API endpoints.

## Frontend Integration

### Component Usage

#### Basic Integration

```tsx
import { AIConsultantPage } from '../components/admin/AIConsultantPage';

// Basic usage in admin dashboard
export const AdminDashboard: React.FC = () => {
  return (
    <div className="admin-layout">
      <AIConsultantPage />
    </div>
  );
};
```

#### Custom Context Integration

```tsx
// With custom client context
const ClientConsultationPage: React.FC<{ clientId: string }> = ({ clientId }) => {
  const [clientData, setClientData] = useState(null);
  
  useEffect(() => {
    // Load client data from your system
    loadClientData(clientId).then(setClientData);
  }, [clientId]);

  return (
    <AIConsultantPage 
      initialContext={{
        clientName: clientData?.name,
        meetingType: clientData?.consultationType
      }}
    />
  );
};
```

### Service Integration

#### Using SimpleAIService

```typescript
import simpleAIService from '../services/simpleAIService';

// Send a message with context
const sendConsultingMessage = async (message: string, context: any) => {
  try {
    const response = await simpleAIService.sendMessage({
      message,
      context: {
        clientName: context.clientName,
        meetingType: context.meetingType
      }
    });
    
    return response;
  } catch (error) {
    console.error('Failed to send message:', error);
    throw error;
  }
};

// Check service health
const checkAIServiceHealth = async () => {
  const isHealthy = await simpleAIService.checkConnection();
  return isHealthy;
};
```

## Backend Integration

### Handler Registration

```go
// Register AI consultant routes
func RegisterAIConsultantRoutes(router *gin.RouterGroup, handler *SimpleChatHandler) {
    aiRoutes := router.Group("/simple-chat")
    aiRoutes.Use(AuthMiddleware()) // Ensure authentication
    
    aiRoutes.POST("/messages", handler.SendMessage)
    aiRoutes.GET("/messages", handler.GetMessages)
}
```

### Custom Handler Implementation

```go
// Extend the simple chat handler
type EnhancedChatHandler struct {
    *SimpleChatHandler
    aiService AIService
    contextManager ContextManager
}

func (h *EnhancedChatHandler) SendMessageWithContext(c *gin.Context) {
    // Custom implementation with AI integration
    var req SimpleSendMessageRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }
    
    // Add context processing
    context := h.contextManager.GetContext(req.SessionID)
    
    // Generate AI response
    aiResponse, err := h.aiService.GenerateResponse(req.Content, context)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "AI service error"})
        return
    }
    
    // Store and return response
    // ... implementation details
}
```

## API Integration

### Authentication Setup

```typescript
// Configure API client with authentication
const apiClient = axios.create({
  baseURL: process.env.REACT_APP_API_URL,
  headers: {
    'Content-Type': 'application/json'
  }
});

// Add JWT token to requests
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('adminToken');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});
```

### Error Handling

```typescript
// Comprehensive error handling
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Handle authentication errors
      localStorage.removeItem('adminToken');
      window.location.href = '/admin/login';
    }
    
    if (error.response?.status >= 500) {
      // Handle server errors
      console.error('Server error:', error);
    }
    
    return Promise.reject(error);
  }
);
```

## Redux Integration

### State Management

```typescript
// Chat slice integration
import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface AIConsultantState {
  sessions: Record<string, ChatSession>;
  activeSessionId: string | null;
  isLoading: boolean;
  error: string | null;
}

const aiConsultantSlice = createSlice({
  name: 'aiConsultant',
  initialState,
  reducers: {
    setActiveSession: (state, action: PayloadAction<string>) => {
      state.activeSessionId = action.payload;
    },
    addMessage: (state, action: PayloadAction<ChatMessage>) => {
      const sessionId = action.payload.session_id;
      if (!state.sessions[sessionId]) {
        state.sessions[sessionId] = { id: sessionId, messages: [] };
      }
      state.sessions[sessionId].messages.push(action.payload);
    }
  }
});
```

### Async Actions

```typescript
// Thunk actions for async operations
export const sendAIMessage = createAsyncThunk(
  'aiConsultant/sendMessage',
  async ({ message, sessionId }: { message: string; sessionId: string }) => {
    const response = await simpleAIService.sendMessage({
      message,
      session_id: sessionId
    });
    
    // Get updated messages
    const messages = await simpleAIService.getMessages();
    return messages;
  }
);
```

This integration guide provides the foundation for implementing AI consultant functionality across different parts of your application.