import React, { useState, useEffect, useRef, useCallback, useMemo } from 'react';
import { Send, MessageSquare, User, Bot, Settings, RefreshCw, AlertCircle, Loader, Maximize2, Minimize2 } from 'lucide-react';
import { useDispatch, useSelector } from 'react-redux';
import { RootState } from '../../store';
import {
  setCurrentSession,
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
import MarkdownRenderer from '../../utils/markdownRenderer';

interface ChatMessage {
  id: string;
  type: 'user' | 'assistant' | 'system';
  content: string;
  timestamp: string;
  session_id: string;
}

interface ChatRequest {
  message: string;
  session_id?: string;
  client_name?: string;
  context?: string;
  quick_action?: string;
}

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

const TABS = [
  { id: 'chat', label: 'AI Chat', icon: MessageSquare },
  { id: 'analysis', label: 'Analysis', icon: Bot },
  { id: 'reports', label: 'Reports', icon: Settings },
  { id: 'templates', label: 'Templates', icon: RefreshCw },
];

export const AIConsultantPage: React.FC = () => {
  const dispatch = useDispatch();
  const connectionState = useSelector((state: RootState) => state.connection);
  const chatState = useSelector((state: RootState) => state.chat);
  
  const [clientName, setClientName] = useState('');
  const [meetingContext, setMeetingContext] = useState('');
  const [showSettings, setShowSettings] = useState(false);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [isServiceHealthy, setIsServiceHealthy] = useState(false);
  const [activeTab, setActiveTab] = useState('chat');
  
  const isConnected = isServiceHealthy;
  const isLoading = chatState.isLoading;
  const sessionId = chatState.currentSession?.id || null;
  
  const inputRef = useRef<HTMLInputElement>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  
  // Use debounced input
  const [inputState, inputActions] = useDebouncedInput({
    delay: 300,
    minLength: 1,
  });

  // Get messages from Redux store
  const messages = useMemo(() => chatState.messages, [chatState.messages]);

  // Auto-scroll to bottom when new messages arrive
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

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

  // Initialize simple AI service
  useEffect(() => {
    const initializeChat = async () => {
      try {
        console.log('[AIConsultantPage] Initializing AI service...');
        const connected = await enhancedAIService.checkConnection();
        setIsServiceHealthy(connected);
        if (connected) {
          console.log('[AIConsultantPage] AI service connected');
        } else {
          console.warn('[AIConsultantPage] AI service offline');
          dispatch(setError('AI service is offline'));
        }
      } catch (error) {
        console.error('[AIConsultantPage] Failed to initialize AI service:', error);
        dispatch(setError('Failed to initialize AI service'));
        setIsServiceHealthy(false);
      }
    };

    initializeChat();

    // Set up periodic connection check
    const connectionCheckInterval = setInterval(async () => {
      const currentHealth = enhancedAIService.isHealthy();
      if (!currentHealth) {
        console.log('[AIConsultantPage] Attempting to reconnect...');
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

  const sendMessage = useCallback(async (message: string, quickAction?: string) => {
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
      
      // Add error message to UI
      const errorMessage: StoreChatMessage = {
        id: `msg-${Date.now()}-error`,
        type: 'assistant',
        content: 'Sorry, I encountered an error processing your request. Please try again.',
        timestamp: new Date().toISOString(),
        session_id: enhancedAIService.getSessionId(),
        status: 'failed',
      };
      dispatch(addMessage(errorMessage));
    } finally {
      dispatch(setLoading(false));
    }
  }, [clientName, meetingContext, dispatch, inputActions]);

  const handleSubmit = useCallback((e: React.FormEvent) => {
    e.preventDefault();
    sendMessage(inputState.value);
  }, [sendMessage, inputState.value]);

  const handleQuickAction = (action: typeof QUICK_ACTIONS[0]) => {
    sendMessage(action.prompt, action.id);
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

  const toggleFullscreen = () => {
    setIsFullscreen(!isFullscreen);
  };

  return (
    <div className={`flex flex-col bg-white ${isFullscreen ? 'fixed inset-0 z-50' : 'h-full'}`}>
      {/* Header */}
      <div className="bg-blue-600 text-white flex-shrink-0">
        <div className="flex items-center justify-between p-4">
          <div className="flex items-center space-x-3">
            <MessageSquare className="h-6 w-6" />
            <div className="flex flex-col">
              <h1 className="text-xl font-semibold">AI Consultant Assistant</h1>
              <span className="text-sm opacity-75">
                {isServiceHealthy ? 'Connected' : 'Offline'}
              </span>
            </div>
            <div className={`w-3 h-3 rounded-full ${isServiceHealthy ? 'bg-green-400' : 'bg-red-400'}`} />
          </div>
        <div className="flex items-center space-x-2">
          <button
            onClick={() => setShowSettings(!showSettings)}
            className="text-white hover:text-gray-200 transition-colors p-2 rounded"
          >
            <Settings className="h-5 w-5" />
          </button>
          <button
            onClick={toggleFullscreen}
            className="text-white hover:text-gray-200 transition-colors p-2 rounded"
          >
            {isFullscreen ? <Minimize2 className="h-5 w-5" /> : <Maximize2 className="h-5 w-5" />}
          </button>
          <button
            onClick={clearChat}
            className="text-white hover:text-gray-200 transition-colors p-2 rounded"
          >
            <RefreshCw className="h-5 w-5" />
          </button>
        </div>
      </div>
      </div>

      {/* Tabs */}
      <div className="flex border-b border-gray-200 bg-white flex-shrink-0">
        {[
          { id: 'chat', label: 'Chat Assistant', icon: MessageSquare },
          { id: 'analysis', label: 'Architecture Analysis', icon: Settings },
          { id: 'reports', label: 'Report Generation', icon: RefreshCw },
          { id: 'settings', label: 'Settings', icon: Settings },
        ].map((tab) => {
          const Icon = tab.icon;
          return (
            <button
              key={tab.id}
              onClick={() => {
                setActiveTab(tab.id);
                if (tab.id === 'settings') {
                  setShowSettings(true);
                } else {
                  setShowSettings(false);
                }
              }}
              className={`flex items-center space-x-2 px-4 py-3 text-sm font-medium border-b-2 transition-colors ${
                activeTab === tab.id
                  ? 'border-blue-500 text-blue-600 bg-blue-50'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              }`}
            >
              <Icon className="h-4 w-4" />
              <span>{tab.label}</span>
            </button>
          );
        })}
      </div>

      {/* Settings Panel */}
      {showSettings && (
        <div className="p-4 border-b border-gray-200 bg-gray-50 flex-shrink-0">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Client Name
              </label>
              <input
                type="text"
                value={clientName}
                onChange={(e) => setClientName(e.target.value)}
                placeholder="Enter client name..."
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Meeting Context
              </label>
              <input
                type="text"
                value={meetingContext}
                onChange={(e) => setMeetingContext(e.target.value)}
                placeholder="e.g., Migration planning, Cost optimization..."
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
          </div>
          <div className="mt-4">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Connection Mode
            </label>
            <div className="flex space-x-2">
              <span className={`px-3 py-1 text-sm rounded ${
                isServiceHealthy 
                  ? 'bg-green-100 text-green-700' 
                  : 'bg-red-100 text-red-700'
              }`}>
                {isServiceHealthy ? 'Connected' : 'Offline'}
              </span>
              <button
                onClick={async () => {
                  const connected = await enhancedAIService.forceReconnect();
                  setIsServiceHealthy(connected);
                  if (!connected) {
                    dispatch(setError('Unable to connect to AI service'));
                  } else {
                    console.log('[AIConsultantPage] Connection test successful');
                  }
                }}
                className="px-3 py-1 text-sm rounded bg-blue-600 text-white hover:bg-blue-700"
              >
                Test Connection
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Tab Content */}
      {activeTab === 'chat' && (
        <>
          {/* Quick Actions */}
          <div className="p-3 border-b border-gray-200 bg-gray-50 flex-shrink-0">
            <div className="flex flex-wrap gap-2">
              {QUICK_ACTIONS.map((action) => (
                <button
                  key={action.id}
                  onClick={() => handleQuickAction(action)}
                  disabled={isLoading}
                  className="px-3 py-1 text-sm bg-blue-100 text-blue-700 rounded hover:bg-blue-200 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                  {action.label}
                </button>
              ))}
            </div>
          </div>

          {/* Messages */}
      <div className="flex-1 overflow-y-auto p-4 space-y-4 min-h-0">
        {!isServiceHealthy && (
          <div className="bg-red-50 border border-red-200 rounded-lg p-4 mb-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-2">
                <AlertCircle className="h-5 w-5 text-red-600" />
                <div>
                  <h3 className="text-sm font-medium text-red-800">Connection Issue</h3>
                  <p className="text-sm text-red-700">
                    Unable to connect to the chat service. Please check your connection and try again.
                  </p>
                </div>
              </div>
              <button
                onClick={async () => {
                  console.log('[AIConsultantPage] Manual reconnect attempt');
                  const connected = await enhancedAIService.forceReconnect();
                  setIsServiceHealthy(connected);
                  if (connected) {
                    console.log('[AIConsultantPage] Reconnected successfully');
                  } else {
                    console.log('[AIConsultantPage] Reconnection failed');
                  }
                }}
                className="px-3 py-2 bg-red-600 text-white rounded hover:bg-red-700 text-sm"
              >
                Retry Connection
              </button>
            </div>
          </div>
        )}

        {messages.length === 0 && !isLoading && (
          <div className="text-center text-gray-500 py-12">
            <MessageSquare className="h-12 w-12 mx-auto mb-4 text-gray-400" />
            <h3 className="text-lg font-medium mb-2">Welcome to AI Consultant Assistant</h3>
            <p className="text-sm max-w-md mx-auto">
              Get instant, expert-level AWS consulting assistance. Ask questions about architecture, costs, security, best practices, and more.
            </p>
            <div className="mt-6">
              <p className="text-xs text-gray-400 mb-2">Try asking:</p>
              <div className="flex flex-wrap justify-center gap-2">
                <span className="px-2 py-1 bg-gray-100 text-gray-600 rounded text-xs">
                  "How much would it cost to migrate to AWS?"
                </span>
                <span className="px-2 py-1 bg-gray-100 text-gray-600 rounded text-xs">
                  "What's the best architecture for high availability?"
                </span>
                <span className="px-2 py-1 bg-gray-100 text-gray-600 rounded text-xs">
                  "How do I secure my AWS environment?"
                </span>
              </div>
            </div>
          </div>
        )}
        
        {messages.map((message) => (
          <div
            key={message.id}
            className={`flex ${message.type === 'user' ? 'justify-end' : 'justify-start'}`}
          >
            <div
              className={`max-w-4xl rounded-lg p-4 shadow-sm ${
                message.type === 'user'
                  ? 'bg-blue-600 text-white'
                  : 'bg-gray-100 text-gray-900 border border-gray-200'
              }`}
            >
              <div className="flex items-start space-x-3">
                {message.type === 'assistant' && (
                  <Bot className="h-5 w-5 mt-0.5 text-blue-600 flex-shrink-0" />
                )}
                {message.type === 'user' && (
                  <User className="h-5 w-5 mt-0.5 text-white flex-shrink-0" />
                )}
                <div className="flex-1 min-w-0">
                  <div className="prose prose-sm max-w-none">
                    {message.type === 'assistant' ? (
                      <MarkdownRenderer content={message.content} />
                    ) : (
                      <p className="whitespace-pre-wrap break-words leading-relaxed">
                        {message.content}
                      </p>
                    )}
                  </div>
                  <p className={`text-xs mt-2 ${
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
            <div className="bg-gray-100 border border-gray-200 rounded-lg p-4 max-w-4xl shadow-sm">
              <div className="flex items-center space-x-3">
                <Bot className="h-5 w-5 text-blue-600 flex-shrink-0" />
                <div className="flex space-x-1">
                  <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" />
                  <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '0.1s' }} />
                  <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '0.2s' }} />
                </div>
                <span className="text-sm text-gray-600">AI is thinking...</span>
              </div>
            </div>
          </div>
        )}
        
        <div ref={messagesEndRef} />
      </div>

          {/* Input */}
          <form onSubmit={handleSubmit} className="p-4 border-t border-gray-200 flex-shrink-0">
            <div className="flex space-x-3">
              <input
                ref={inputRef}
                type="text"
                value={inputState.value}
                onChange={(e) => inputActions.setValue(e.target.value)}
                placeholder="Ask about AWS services, costs, best practices, architecture..."
                disabled={isLoading}
                className="flex-1 px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:bg-gray-50"
              />
              {inputState.isDebouncing && (
                <div className="flex items-center px-3 text-gray-400">
                  <Loader className="h-5 w-5 animate-spin" />
                </div>
              )}
              <button
                type="submit"
                disabled={isLoading || !inputState.value.trim()}
                className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors flex-shrink-0"
              >
                <Send className="h-5 w-5" />
              </button>
            </div>
          </form>
        </>
      )}

      {/* Architecture Analysis Tab */}
      {activeTab === 'analysis' && (
        <div className="flex-1 p-6 overflow-y-auto">
          <div className="max-w-4xl mx-auto">
            <h2 className="text-2xl font-bold text-gray-900 mb-6">Architecture Analysis</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="bg-white border border-gray-200 rounded-lg p-6">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">Current Architecture Review</h3>
                <p className="text-gray-600 mb-4">Upload your current architecture diagrams or describe your infrastructure for AI-powered analysis.</p>
                <button className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700">
                  Start Analysis
                </button>
              </div>
              <div className="bg-white border border-gray-200 rounded-lg p-6">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">Well-Architected Review</h3>
                <p className="text-gray-600 mb-4">Get recommendations based on AWS Well-Architected Framework principles.</p>
                <button className="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700">
                  Begin Review
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Report Generation Tab */}
      {activeTab === 'reports' && (
        <div className="flex-1 p-6 overflow-y-auto">
          <div className="max-w-4xl mx-auto">
            <h2 className="text-2xl font-bold text-gray-900 mb-6">Report Generation</h2>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div className="bg-white border border-gray-200 rounded-lg p-6">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">Migration Assessment</h3>
                <p className="text-gray-600 mb-4">Generate comprehensive migration reports with cost estimates and timelines.</p>
                <button className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700">
                  Generate Report
                </button>
              </div>
              <div className="bg-white border border-gray-200 rounded-lg p-6">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">Cost Optimization</h3>
                <p className="text-gray-600 mb-4">Analyze current spending and get optimization recommendations.</p>
                <button className="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700">
                  Analyze Costs
                </button>
              </div>
              <div className="bg-white border border-gray-200 rounded-lg p-6">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">Security Assessment</h3>
                <p className="text-gray-600 mb-4">Review security posture and compliance requirements.</p>
                <button className="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700">
                  Security Review
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Settings content is already handled by the showSettings state */}
    </div>
  );
};

export default AIConsultantPage;