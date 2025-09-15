import React, { useState, useEffect, useRef, useCallback, useMemo } from 'react';
import { Send, MessageSquare, User, Bot, Minimize2, Maximize2, X, AlertCircle, Loader } from 'lucide-react';
import { useDispatch, useSelector } from 'react-redux';
import { RootState } from '../../store';
import {
  clearCurrentSession,
  updateSessionContextAction,
  addMessage,
  setLoading,
  clearMessages,
  setError,
  ChatMessage as StoreChatMessage,
  SessionContext,
} from '../../store/slices/chatSlice';
import enhancedAIService from '../../services/simpleAIService';
import useDebouncedInput from '../../hooks/useDebouncedInput';

interface ConsultantChatProps {
  isMinimized?: boolean;
  onToggleMinimize?: () => void;
  onClose?: () => void;
}

const QUICK_ACTIONS = [
  { id: 'cost-estimate', label: 'Cost Estimate', prompt: 'Provide a cost estimate for this solution' },
  { id: 'security-review', label: 'Security Review', prompt: 'What are the security considerations for this approach?' },
  { id: 'best-practices', label: 'Best Practices', prompt: 'What are the AWS best practices for this scenario?' },
  { id: 'alternatives', label: 'Alternatives', prompt: 'What are alternative approaches to consider?' },
  { id: 'next-steps', label: 'Next Steps', prompt: 'What are the recommended next steps?' },
];

export const ConsultantChat: React.FC<ConsultantChatProps> = ({
  isMinimized = false,
  onToggleMinimize,
  onClose
}) => {
  // Get state from Redux store first
  const dispatch = useDispatch();
  const connectionState = useSelector((state: RootState) => state.connection);
  const chatState = useSelector((state: RootState) => state.chat);
  
  const isConnected = connectionState.status === 'connected' || connectionState.status === 'polling';
  const isLoading = chatState.isLoading;
  const sessionId = chatState.currentSession?.id || null;

  const [clientName, setClientName] = useState('');
  const [meetingContext, setMeetingContext] = useState('');
  const [isServiceHealthy, setIsServiceHealthy] = useState(false);
  
  const inputRef = useRef<HTMLInputElement>(null);
  
  // Update session context when client name or meeting context changes
  useEffect(() => {
    if (sessionId && (clientName || meetingContext)) {
      const contextUpdate: SessionContext = {
        client_name: clientName,
        meeting_type: meetingContext,
      };
      dispatch(updateSessionContextAction(contextUpdate));
    }
  }, [clientName, meetingContext, sessionId, dispatch]);

  // Use debounced input for search and typing
  const [inputState, inputActions] = useDebouncedInput({
    delay: 300,
    minLength: 1,
  });

  // Get messages from Redux store
  const messages = useMemo(() => chatState.messages, [chatState.messages]);

  // Initialize simple AI service
  useEffect(() => {
    const initializeChat = async () => {
      try {
        console.log('[ConsultantChat] Initializing AI service...');
        const connected = await enhancedAIService.checkConnection();
        setIsServiceHealthy(connected);
        if (connected) {
          console.log('[ConsultantChat] AI service initialized successfully');
        } else {
          console.warn('[ConsultantChat] AI service connection failed');
          dispatch(setError('Failed to connect to AI service'));
        }
      } catch (error) {
        console.error('[ConsultantChat] Failed to initialize AI service:', error);
        dispatch(setError('Failed to initialize AI service'));
        setIsServiceHealthy(false);
      }
    };

    initializeChat();

    // Set up periodic connection check
    const connectionCheckInterval = setInterval(async () => {
      const currentHealth = enhancedAIService.isHealthy();
      if (!currentHealth) {
        console.log('[ConsultantChat] Attempting to reconnect...');
        const reconnected = await enhancedAIService.forceReconnect();
        setIsServiceHealthy(reconnected);
      } else {
        setIsServiceHealthy(true);
      }
    }, 30000); // Check every 30 seconds

    return () => {
      clearInterval(connectionCheckInterval);
    };
  }, [dispatch]);

  const sendMessage = useCallback(async (message: string) => {
    if (!message.trim()) return;

    try {
      dispatch(setLoading(true));
      inputActions.clearValue();

      // Add user message to UI immediately
      const userMessage: StoreChatMessage = {
        id: `msg-${Date.now()}-user`,
        type: 'user',
        content: message.trim(),
        timestamp: new Date().toISOString(),
        session_id: enhancedAIService.getSessionId(),
        status: 'sending',
      };
      dispatch(addMessage(userMessage));

      // Send using enhanced AI service
      const response = await enhancedAIService.sendMessage({
        message: message.trim(),
        context: {
          clientName: clientName || undefined,
          meetingType: meetingContext || undefined,
        }
      });

      // Add AI response to UI
      const aiMessage: StoreChatMessage = {
        id: `msg-${Date.now()}-ai`,
        type: 'assistant',
        content: response.content,
        timestamp: response.timestamp,
        session_id: enhancedAIService.getSessionId(),
        status: 'delivered',
      };
      dispatch(addMessage(aiMessage));
      
    } catch (error) {
      console.error('Failed to send message:', error);
      dispatch(setError('Failed to send message. Please try again.'));
    } finally {
      dispatch(setLoading(false));
    }
  }, [clientName, meetingContext, dispatch, inputActions]);

  const handleSubmit = useCallback((e: React.FormEvent) => {
    e.preventDefault();
    sendMessage(inputState.value);
  }, [sendMessage, inputState.value]);

  const handleQuickAction = (action: typeof QUICK_ACTIONS[0]) => {
    sendMessage(action.prompt);
  };

  const formatTimestamp = (timestamp: string) => {
    return new Date(timestamp).toLocaleTimeString([], { 
      hour: '2-digit', 
      minute: '2-digit' 
    });
  };

  const clearChat = useCallback(() => {
    dispatch(clearMessages());
    dispatch(clearCurrentSession());
    enhancedAIService.resetSession();
  }, [dispatch]);

  if (isMinimized) {
    return (
      <div className="fixed bottom-4 right-4 z-50">
        <button
          onClick={onToggleMinimize}
          className="bg-blue-600 hover:bg-blue-700 text-white p-3 rounded-full shadow-lg transition-colors"
        >
          <MessageSquare className="h-6 w-6" />
        </button>
      </div>
    );
  }

  return (
    <div className="fixed bottom-4 right-4 w-80 h-[500px] bg-white rounded-lg shadow-xl border border-gray-200 flex flex-col z-50">
      {/* Header */}
      <div className="flex items-center justify-between p-2 border-b border-gray-200 bg-blue-600 text-white rounded-t-lg flex-shrink-0">
        <div className="flex items-center space-x-2 min-w-0">
          <MessageSquare className="h-4 w-4 flex-shrink-0" />
          <div className="flex flex-col min-w-0">
            <h3 className="font-semibold text-xs truncate">AI Assistant</h3>
            <span className="text-xs opacity-75 truncate">
              {isServiceHealthy ? 'Connected' : 'Offline'}
            </span>
          </div>
          <div className={`w-2 h-2 rounded-full flex-shrink-0 ${isServiceHealthy ? 'bg-green-400' : 'bg-red-400'}`} />
        </div>
        <div className="flex items-center space-x-1 flex-shrink-0">
          <button
            onClick={() => window.open('/admin/ai-consultant', '_blank')}
            className="text-white hover:text-gray-200 transition-colors p-1 bg-blue-700 rounded text-xs"
            title="Open Full AI Assistant"
          >
            <Maximize2 className="h-3 w-3" />
          </button>
          {onToggleMinimize && (
            <button
              onClick={onToggleMinimize}
              className="text-white hover:text-gray-200 transition-colors p-1"
            >
              <Minimize2 className="h-3 w-3" />
            </button>
          )}
          {onClose && (
            <button
              onClick={onClose}
              className="text-white hover:text-gray-200 transition-colors p-1"
            >
              <X className="h-3 w-3" />
            </button>
          )}
        </div>
      </div>



      {/* Quick Actions */}
      <div className="p-2 border-b border-gray-200 bg-gray-50">
        <div className="flex flex-wrap gap-1">
          {QUICK_ACTIONS.map((action) => (
            <button
              key={action.id}
              onClick={() => handleQuickAction(action)}
              disabled={isLoading}
              className="px-2 py-1 text-xs bg-blue-100 text-blue-700 rounded hover:bg-blue-200 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {action.label}
            </button>
          ))}
        </div>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-hidden relative">
        <div className="h-full overflow-y-auto p-2 space-y-2">
          {!isServiceHealthy && (
            <div className="bg-yellow-50 border border-yellow-200 rounded p-2 mb-2">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-1">
                  <AlertCircle className="h-3 w-3 text-yellow-600 flex-shrink-0" />
                  <p className="text-xs text-yellow-700">Connection lost. Trying to reconnect...</p>
                </div>
                <button
                  onClick={async () => {
                    console.log('[ConsultantChat] Manual reconnect attempt');
                    const connected = await enhancedAIService.forceReconnect();
                    setIsServiceHealthy(connected);
                    if (connected) {
                      console.log('[ConsultantChat] Reconnected successfully');
                    } else {
                      console.log('[ConsultantChat] Reconnection failed');
                    }
                  }}
                  className="text-xs bg-yellow-600 text-white px-2 py-1 rounded hover:bg-yellow-700"
                >
                  Retry
                </button>
              </div>
            </div>
          )}

          {messages.length === 0 && !isLoading && (
            <div className="text-center text-gray-500 text-xs py-4">
              <MessageSquare className="h-6 w-6 mx-auto mb-2 text-gray-400" />
              <p className="px-2">Start a conversation or click to open full AI Assistant.</p>
              <button
                onClick={() => window.open('/admin/ai-consultant', '_blank')}
                className="mt-2 text-blue-600 hover:text-blue-800 underline text-xs"
              >
                Open Full AI Assistant
              </button>
            </div>
          )}
          
          {messages.map((message) => (
            <div
              key={message.id}
              className={`flex ${message.type === 'user' ? 'justify-end' : 'justify-start'}`}
            >
              <div
                className={`max-w-[90%] rounded p-2 text-xs ${
                  message.type === 'user'
                    ? 'bg-blue-600 text-white'
                    : 'bg-gray-100 text-gray-900'
                }`}
              >
                <div className="flex items-start space-x-1">
                  {message.type === 'assistant' && (
                    <Bot className="h-3 w-3 mt-0.5 text-blue-600 flex-shrink-0" />
                  )}
                  {message.type === 'user' && (
                    <User className="h-3 w-3 mt-0.5 text-white flex-shrink-0" />
                  )}
                  <div className="flex-1 min-w-0">
                    <p className="whitespace-pre-wrap break-words leading-tight">
                      {message.content.length > 200 ? `${message.content.substring(0, 200)}...` : message.content}
                    </p>
                    <p className={`text-xs mt-1 opacity-75 ${
                      message.type === 'user' ? 'text-blue-200' : 'text-gray-500'
                    }`}>
                      {formatTimestamp(message.timestamp)}
                    </p>
                  </div>
                </div>
              </div>
            </div>
          ))}
          
          {isLoading && (
            <div className="flex justify-start">
              <div className="bg-gray-100 rounded p-2 max-w-[90%]">
                <div className="flex items-center space-x-1">
                  <Bot className="h-3 w-3 text-blue-600 flex-shrink-0" />
                  <div className="flex space-x-1">
                    <div className="w-1 h-1 bg-gray-400 rounded-full animate-bounce" />
                    <div className="w-1 h-1 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '0.1s' }} />
                    <div className="w-1 h-1 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '0.2s' }} />
                  </div>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Input */}
      <form onSubmit={handleSubmit} className="p-3 border-t border-gray-200 flex-shrink-0">
        <div className="flex space-x-2">
          <input
            ref={inputRef}
            type="text"
            value={inputState.value}
            onChange={(e) => inputActions.setValue(e.target.value)}
            placeholder={isConnected ? "Ask about AWS services, costs, best practices..." : "Connecting..."}
            disabled={isLoading || !isConnected}
            className="flex-1 px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:bg-gray-50"
          />
          {inputState.isDebouncing && (
            <div className="flex items-center px-2 text-gray-400">
              <Loader className="h-4 w-4 animate-spin" />
            </div>
          )}
          <button
            type="submit"
            disabled={isLoading || !isConnected || !inputState.value.trim()}
            className="px-3 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors flex-shrink-0"
          >
            <Send className="h-4 w-4" />
          </button>
        </div>
      </form>
    </div>
  );
};

export default ConsultantChat;