import React, { useState, useEffect, useRef, useCallback, useMemo } from 'react';
import { Send, MessageSquare, User, Bot, Search, Filter, Settings, History, Download, Copy, RefreshCw, AlertCircle, CheckCircle, Loader, X, Plus } from 'lucide-react';
import { useDispatch, useSelector } from 'react-redux';
import { RootState } from '../../store';
import {
  setCurrentSession,
  updateSessionContextAction,
  addOptimisticMessage,
  setLoading,
  clearMessages,
  updateSettings,
  setError,
  ChatMessage,
  SessionContext,
  ChatSession,
} from '../../store/slices/chatSlice';
import { pollingChatService } from '../../services/pollingChatService';
import ChatSessionManager from './ChatSessionManager';

interface ChatRequest {
  message: string;
  session_id?: string;
  client_name?: string;
  context?: string;
  quick_action?: string;
}

const QUICK_ACTIONS = [
  { id: 'cost-estimate', label: 'Cost Estimate', prompt: 'Provide a cost estimate for this solution', category: 'analysis' },
  { id: 'security-review', label: 'Security Review', prompt: 'What are the security considerations for this approach?', category: 'security' },
  { id: 'best-practices', label: 'Best Practices', prompt: 'What are the AWS best practices for this scenario?', category: 'guidance' },
  { id: 'alternatives', label: 'Alternatives', prompt: 'What are alternative approaches to consider?', category: 'analysis' },
  { id: 'next-steps', label: 'Next Steps', prompt: 'What are the recommended next steps?', category: 'guidance' },
  { id: 'architecture-review', label: 'Architecture Review', prompt: 'Review the architecture and suggest improvements', category: 'architecture' },
  { id: 'performance-optimization', label: 'Performance', prompt: 'How can we optimize performance for this solution?', category: 'optimization' },
  { id: 'disaster-recovery', label: 'Disaster Recovery', prompt: 'What disaster recovery options should we consider?', category: 'resilience' },
];

const QUICK_ACTION_CATEGORIES = [
  { id: 'all', label: 'All Actions' },
  { id: 'analysis', label: 'Analysis' },
  { id: 'security', label: 'Security' },
  { id: 'guidance', label: 'Guidance' },
  { id: 'architecture', label: 'Architecture' },
  { id: 'optimization', label: 'Optimization' },
  { id: 'resilience', label: 'Resilience' },
];

export const ChatPage: React.FC = () => {
  // Redux state
  const dispatch = useDispatch();
  const { 
    currentSession, 
    messages, 
    sessionContext, 
    isLoading, 
    isTyping, 
    error,
    settings 
  } = useSelector((state: RootState) => state.chat);
  const { status: connectionStatus, isHealthy } = useSelector((state: RootState) => state.connection);

  // Local state
  const [inputMessage, setInputMessage] = useState('');
  const [searchQuery, setSearchQuery] = useState('');
  const [showSettings, setShowSettings] = useState(false);
  const [showSessionContext, setShowSessionContext] = useState(false);
  const [showSessionManager, setShowSessionManager] = useState(false);
  const [selectedQuickActionCategory, setSelectedQuickActionCategory] = useState('all');
  const [messageFilter, setMessageFilter] = useState<'all' | 'user' | 'assistant'>('all');
  
  // Refs
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const messagesContainerRef = useRef<HTMLDivElement>(null);

  // Auto-scroll to bottom when new messages arrive
  const scrollToBottom = useCallback(() => {
    if (settings.autoScroll) {
      messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    }
  }, [settings.autoScroll]);

  useEffect(() => {
    scrollToBottom();
  }, [messages, scrollToBottom]);


  // Focus input when component mounts
  useEffect(() => {
    inputRef.current?.focus();
  }, []);

  // Start/stop polling based on session
  useEffect(() => {
    if (currentSession) {
      console.log('[ChatPage] Starting polling for session:', currentSession.id);
      pollingChatService.startPolling();
    } else {
      console.log('[ChatPage] No session, stopping polling');
      pollingChatService.stopPolling();
    }

    // Cleanup on unmount
    return () => {
      console.log('[ChatPage] Component unmounting, stopping polling');
      pollingChatService.stopPolling();
    };
  }, [currentSession]);

  // Keyboard navigation
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      // Escape to close modals
      if (e.key === 'Escape') {
        if (showSessionManager) {
          setShowSessionManager(false);
        } else if (showSessionContext) {
          setShowSessionContext(false);
        }
      }
      
      // Ctrl/Cmd + K to open session manager
      if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
        e.preventDefault();
        setShowSessionManager(true);
      }
      
      // Ctrl/Cmd + / to focus search
      if ((e.ctrlKey || e.metaKey) && e.key === '/') {
        e.preventDefault();
        const searchInput = document.querySelector('input[placeholder*="Search"]') as HTMLInputElement;
        searchInput?.focus();
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [showSessionManager, showSessionContext]);

  // Filtered messages based on search and filter
  const filteredMessages = useMemo(() => {
    let filtered = messages;

    // Apply message type filter
    if (messageFilter !== 'all') {
      filtered = filtered.filter(msg => msg.type === messageFilter);
    }

    // Apply search filter
    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter(msg => 
        msg.content.toLowerCase().includes(query)
      );
    }

    return filtered;
  }, [messages, messageFilter, searchQuery]);

  // Filtered quick actions based on category
  const filteredQuickActions = useMemo(() => {
    if (selectedQuickActionCategory === 'all') {
      return QUICK_ACTIONS;
    }
    return QUICK_ACTIONS.filter(action => action.category === selectedQuickActionCategory);
  }, [selectedQuickActionCategory]);

  const sendMessage = useCallback(async (message: string, quickAction?: string) => {
    if (!message.trim()) return;

    const request: ChatRequest = {
      message: message.trim(),
      session_id: currentSession?.id,
      client_name: sessionContext.client_name,
      context: sessionContext.project_context,
      quick_action: quickAction,
    };

    try {
      // Send via polling service
      await pollingChatService.sendMessage(request);
      setInputMessage('');
    } catch (error) {
      console.error('Failed to send message:', error);
      dispatch(setError('Failed to send message. Please try again.'));
    }
  }, [currentSession?.id, sessionContext, dispatch]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    sendMessage(inputMessage);
  };

  const handleQuickAction = (action: typeof QUICK_ACTIONS[0]) => {
    sendMessage(action.prompt, action.id);
  };

  const formatTimestamp = (timestamp: string) => {
    const date = new Date(timestamp);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / (1000 * 60));
    const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
    const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

    if (diffMins < 1) return 'Just now';
    if (diffMins < 60) return `${diffMins}m ago`;
    if (diffHours < 24) return `${diffHours}h ago`;
    if (diffDays < 7) return `${diffDays}d ago`;
    
    return date.toLocaleDateString([], { 
      month: 'short', 
      day: 'numeric',
      ...(date.getFullYear() !== now.getFullYear() && { year: 'numeric' })
    });
  };

  const clearChat = () => {
    dispatch(clearMessages());
  };

  const copyMessage = (content: string) => {
    navigator.clipboard.writeText(content);
  };

  const exportChat = () => {
    const chatData = {
      session: currentSession,
      context: sessionContext,
      messages: messages,
      exportedAt: new Date().toISOString(),
    };
    
    const blob = new Blob([JSON.stringify(chatData, null, 2)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `chat-export-${currentSession?.id || 'session'}-${new Date().toISOString().split('T')[0]}.json`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  const updateSessionContextField = (field: keyof SessionContext, value: string) => {
    dispatch(updateSessionContextAction({ [field]: value }));
  };

  const handleSessionSelect = (session: ChatSession) => {
    dispatch(setCurrentSession(session));
  };

  return (
    <div className="flex h-full bg-white" role="main" aria-label="AI Consultant Chat">
      {/* Main Chat Area */}
      <div className="flex-1 flex flex-col min-w-0">
        {/* Header */}
        <header className="flex items-center justify-between p-3 sm:p-4 border-b border-gray-200 bg-white flex-shrink-0">
          <div className="flex items-center space-x-2 sm:space-x-3 min-w-0 flex-1">
            <MessageSquare className="h-5 w-5 sm:h-6 sm:w-6 text-blue-600 flex-shrink-0" aria-hidden="true" />
            <div className="min-w-0 flex-1">
              <h1 className="text-lg sm:text-xl font-semibold text-gray-900 truncate">
                AI Consultant Chat
              </h1>
              <div className="flex items-center space-x-2 text-xs sm:text-sm text-gray-500">
                <div 
                  className={`w-2 h-2 rounded-full flex-shrink-0 ${
                    connectionStatus === 'connected' ? 'bg-green-500' : 
                    connectionStatus === 'connecting' ? 'bg-yellow-500' : 'bg-red-500'
                  }`}
                  aria-hidden="true"
                />
                <span className="capitalize" aria-live="polite">
                  {connectionStatus}
                </span>
                {currentSession && (
                  <>
                    <span aria-hidden="true">•</span>
                    <span className="hidden sm:inline">
                      Session: {currentSession.id.slice(-8)}
                    </span>
                  </>
                )}
              </div>
            </div>
          </div>
          
          <div className="flex items-center space-x-1 sm:space-x-2 flex-shrink-0">
            <button
              onClick={() => setShowSessionManager(true)}
              className="p-1.5 sm:p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-md transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
              title="Session Manager (Ctrl+K)"
              aria-label="Open session manager"
            >
              <History className="h-4 w-4 sm:h-5 sm:w-5" />
            </button>
            <button
              onClick={() => setShowSessionContext(!showSessionContext)}
              className="p-1.5 sm:p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-md transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
              title="Session Context"
              aria-label={`${showSessionContext ? 'Hide' : 'Show'} session context`}
              aria-expanded={showSessionContext}
            >
              <Settings className="h-4 w-4 sm:h-5 sm:w-5" />
            </button>
            <button
              onClick={exportChat}
              disabled={messages.length === 0}
              className="p-1.5 sm:p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-md transition-colors disabled:opacity-50 disabled:cursor-not-allowed focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
              title="Export Chat"
              aria-label="Export chat history"
            >
              <Download className="h-4 w-4 sm:h-5 sm:w-5" />
            </button>
            <button
              onClick={clearChat}
              className="p-1.5 sm:p-2 text-gray-500 hover:text-red-600 hover:bg-red-50 rounded-md transition-colors focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
              title="Clear Chat"
              aria-label="Clear all messages"
            >
              <RefreshCw className="h-4 w-4 sm:h-5 sm:w-5" />
            </button>
          </div>
        </header>

        {/* Session Context Panel */}
        {showSessionContext && (
          <section 
            className="p-3 sm:p-4 border-b border-gray-200 bg-gray-50"
            aria-label="Session context settings"
          >
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3 sm:gap-4">
              <div>
                <label 
                  htmlFor="client-name"
                  className="block text-sm font-medium text-gray-700 mb-1"
                >
                  Client Name
                </label>
                <input
                  id="client-name"
                  type="text"
                  value={sessionContext.client_name || ''}
                  onChange={(e) => updateSessionContextField('client_name', e.target.value)}
                  placeholder="Enter client name..."
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  aria-describedby="client-name-help"
                />
                <div id="client-name-help" className="sr-only">
                  Enter the name of the client for this session
                </div>
              </div>
              <div>
                <label 
                  htmlFor="meeting-type"
                  className="block text-sm font-medium text-gray-700 mb-1"
                >
                  Meeting Type
                </label>
                <input
                  id="meeting-type"
                  type="text"
                  value={sessionContext.meeting_type || ''}
                  onChange={(e) => updateSessionContextField('meeting_type', e.target.value)}
                  placeholder="e.g., Discovery, Architecture Review..."
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  aria-describedby="meeting-type-help"
                />
                <div id="meeting-type-help" className="sr-only">
                  Specify the type of meeting or consultation
                </div>
              </div>
              <div className="sm:col-span-2 lg:col-span-1">
                <label 
                  htmlFor="project-context"
                  className="block text-sm font-medium text-gray-700 mb-1"
                >
                  Project Context
                </label>
                <input
                  id="project-context"
                  type="text"
                  value={sessionContext.project_context || ''}
                  onChange={(e) => updateSessionContextField('project_context', e.target.value)}
                  placeholder="e.g., Migration, Cost optimization..."
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  aria-describedby="project-context-help"
                />
                <div id="project-context-help" className="sr-only">
                  Provide context about the project or consultation topic
                </div>
              </div>
            </div>
          </section>
        )}

        {/* Search and Filter Bar */}
        <section className="p-3 sm:p-4 border-b border-gray-200 bg-white" aria-label="Message search and filters">
          <div className="flex flex-col sm:flex-row items-stretch sm:items-center gap-3 sm:gap-4">
            <div className="flex-1 relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" aria-hidden="true" />
              <input
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="Search messages... (Ctrl+/)"
                className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                aria-label="Search messages"
              />
            </div>
            <div className="flex items-center space-x-2 flex-shrink-0">
              <Filter className="h-4 w-4 text-gray-500" aria-hidden="true" />
              <label htmlFor="message-filter" className="sr-only">Filter messages by type</label>
              <select
                id="message-filter"
                value={messageFilter}
                onChange={(e) => setMessageFilter(e.target.value as 'all' | 'user' | 'assistant')}
                className="border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-sm"
              >
                <option value="all">All Messages</option>
                <option value="user">My Messages</option>
                <option value="assistant">AI Responses</option>
              </select>
            </div>
          </div>
        </section>

        {/* Quick Actions */}
        <section className="p-3 sm:p-4 border-b border-gray-200 bg-gray-50" aria-label="Quick action buttons">
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2 sm:gap-3 mb-3">
            <h3 className="text-sm font-medium text-gray-700">Quick Actions</h3>
            <div className="flex items-center space-x-2">
              <label htmlFor="action-category" className="sr-only">Filter actions by category</label>
              <select
                id="action-category"
                value={selectedQuickActionCategory}
                onChange={(e) => setSelectedQuickActionCategory(e.target.value)}
                className="text-xs border border-gray-300 rounded px-2 py-1 focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
              >
                {QUICK_ACTION_CATEGORIES.map(category => (
                  <option key={category.id} value={category.id}>{category.label}</option>
                ))}
              </select>
            </div>
          </div>
          <div className="flex flex-wrap gap-2" role="group" aria-label="Quick action buttons">
            {filteredQuickActions.map((action) => (
              <button
                key={action.id}
                onClick={() => handleQuickAction(action)}
                disabled={isLoading || (connectionStatus !== 'connected' && connectionStatus !== 'polling')}
                className="px-3 py-1.5 text-sm bg-blue-100 text-blue-700 rounded-md hover:bg-blue-200 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-colors touch-manipulation"
                aria-label={`Quick action: ${action.label}`}
                title={action.prompt}
              >
                {action.label}
              </button>
            ))}
          </div>
        </section>

        {/* Messages */}
        <div 
          ref={messagesContainerRef}
          className="flex-1 overflow-y-auto p-3 sm:p-4 space-y-3 sm:space-y-4"
          role="log"
          aria-label="Chat messages"
          aria-live="polite"
          aria-atomic="false"
        >
          {filteredMessages.length === 0 && !searchQuery && (
            <div className="text-center text-gray-500 py-8 sm:py-12" role="status">
              <MessageSquare className="h-10 w-10 sm:h-12 sm:w-12 mx-auto mb-3 sm:mb-4 text-gray-400" aria-hidden="true" />
              <h3 className="text-base sm:text-lg font-medium text-gray-900 mb-2">Welcome to AI Consultant Chat</h3>
              <p className="text-sm sm:text-base text-gray-600 max-w-md mx-auto px-4">
                Start a conversation to get real-time AWS consulting assistance. Use quick actions above or type your question below.
              </p>
            </div>
          )}

          {filteredMessages.length === 0 && searchQuery && (
            <div className="text-center text-gray-500 py-8 sm:py-12" role="status">
              <Search className="h-10 w-10 sm:h-12 sm:w-12 mx-auto mb-3 sm:mb-4 text-gray-400" aria-hidden="true" />
              <h3 className="text-base sm:text-lg font-medium text-gray-900 mb-2">No messages found</h3>
              <p className="text-sm sm:text-base text-gray-600 px-4">
                No messages match your search for "{searchQuery}"
              </p>
            </div>
          )}
          
          {filteredMessages.map((message) => (
            <div
              key={message.id}
              className={`flex ${message.type === 'user' ? 'justify-end' : 'justify-start'}`}
              role="article"
              aria-label={`${message.type === 'user' ? 'Your message' : 'AI response'} at ${formatTimestamp(message.timestamp)}`}
            >
              <div
                className={`max-w-[85%] sm:max-w-[80%] rounded-lg p-3 sm:p-4 group relative ${
                  message.type === 'user'
                    ? 'bg-blue-600 text-white'
                    : 'bg-gray-100 text-gray-900 border border-gray-200'
                }`}
              >
                <div className="flex items-start space-x-2 sm:space-x-3">
                  {message.type === 'assistant' && (
                    <Bot 
                      className="h-4 w-4 sm:h-5 sm:w-5 mt-0.5 text-blue-600 flex-shrink-0" 
                      aria-hidden="true"
                    />
                  )}
                  {message.type === 'user' && (
                    <User 
                      className="h-4 w-4 sm:h-5 sm:w-5 mt-0.5 text-white flex-shrink-0" 
                      aria-hidden="true"
                    />
                  )}
                  <div className="flex-1 min-w-0">
                    <div className="whitespace-pre-wrap break-words text-sm sm:text-base">
                      {message.content}
                    </div>
                    <div className="flex items-center justify-between mt-2">
                      <div className={`text-xs ${
                        message.type === 'user' ? 'text-blue-200' : 'text-gray-500'
                      }`}>
                        {settings.showTimestamps && (
                          <time dateTime={message.timestamp}>
                            {formatTimestamp(message.timestamp)}
                          </time>
                        )}
                        {message.status && (
                          <span className="ml-2 inline-flex items-center" aria-label={`Message status: ${message.status}`}>
                            {message.status === 'sending' && (
                              <Loader className="h-3 w-3 animate-spin" aria-hidden="true" />
                            )}
                            {message.status === 'sent' && (
                              <CheckCircle className="h-3 w-3" aria-hidden="true" />
                            )}
                            {message.status === 'delivered' && (
                              <CheckCircle className="h-3 w-3" aria-hidden="true" />
                            )}
                            {message.status === 'failed' && (
                              <AlertCircle className="h-3 w-3 text-red-500" aria-hidden="true" />
                            )}
                          </span>
                        )}
                      </div>
                      <button
                        onClick={() => copyMessage(message.content)}
                        className={`opacity-0 group-hover:opacity-100 focus:opacity-100 transition-opacity p-1 rounded focus:outline-none focus:ring-2 focus:ring-offset-2 ${
                          message.type === 'user' 
                            ? 'hover:bg-blue-700 text-blue-200 focus:ring-blue-300' 
                            : 'hover:bg-gray-200 text-gray-500 focus:ring-gray-400'
                        }`}
                        title="Copy message"
                        aria-label="Copy message to clipboard"
                      >
                        <Copy className="h-3 w-3" aria-hidden="true" />
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          ))}
          
          {isTyping && (
            <div className="flex justify-start" role="status" aria-label="AI is typing">
              <div className="bg-gray-100 rounded-lg p-3 sm:p-4 max-w-[85%] sm:max-w-[80%] border border-gray-200">
                <div className="flex items-center space-x-2 sm:space-x-3">
                  <Bot className="h-4 w-4 sm:h-5 sm:w-5 text-blue-600" aria-hidden="true" />
                  <div className="flex space-x-1" aria-hidden="true">
                    <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" />
                    <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '0.1s' }} />
                    <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '0.2s' }} />
                  </div>
                  <span className="sr-only">AI is typing a response</span>
                </div>
              </div>
            </div>
          )}
          
          <div ref={messagesEndRef} />
        </div>

        {/* Error Display */}
        {error && (
          <div className="p-3 sm:p-4 bg-red-50 border-t border-red-200" role="alert" aria-live="assertive">
            <div className="flex items-center space-x-2 text-red-700">
              <AlertCircle className="h-4 w-4 flex-shrink-0" aria-hidden="true" />
              <span className="text-sm flex-1">{error}</span>
              <button
                onClick={() => dispatch({ type: 'chat/clearError' })}
                className="text-red-500 hover:text-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 rounded p-1"
                aria-label="Dismiss error"
              >
                <X className="h-4 w-4" />
              </button>
            </div>
          </div>
        )}

        {/* Input */}
        <form 
          onSubmit={handleSubmit} 
          className="p-3 sm:p-4 border-t border-gray-200 bg-white flex-shrink-0"
          aria-label="Send message form"
        >
          <div className="flex space-x-2 sm:space-x-3">
            <label htmlFor="message-input" className="sr-only">
              Type your message
            </label>
            <input
              id="message-input"
              ref={inputRef}
              type="text"
              value={inputMessage}
              onChange={(e) => setInputMessage(e.target.value)}
              placeholder={
                (connectionStatus === 'connected' || connectionStatus === 'polling')
                  ? "Ask about AWS services, costs, best practices..." 
                  : "Connecting..."
              }
              disabled={isLoading || (connectionStatus !== 'connected' && connectionStatus !== 'polling')}
              className="flex-1 px-3 sm:px-4 py-2 sm:py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 disabled:opacity-50 disabled:bg-gray-50 text-sm sm:text-base"
              aria-describedby="message-input-help"
            />
            <div id="message-input-help" className="sr-only">
              Type your question or message and press Enter or click Send to submit
            </div>
            <button
              type="submit"
              disabled={isLoading || (connectionStatus !== 'connected' && connectionStatus !== 'polling') || !inputMessage.trim()}
              className="px-4 sm:px-6 py-2 sm:py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-colors flex items-center space-x-2 touch-manipulation"
              aria-label={isLoading ? 'Sending message' : 'Send message'}
            >
              {isLoading ? (
                <Loader className="h-4 w-4 animate-spin" aria-hidden="true" />
              ) : (
                <Send className="h-4 w-4" aria-hidden="true" />
              )}
              <span className="hidden sm:inline">Send</span>
            </button>
          </div>
        </form>
      </div>

      {/* Session Manager */}
      <ChatSessionManager
        isOpen={showSessionManager}
        onClose={() => setShowSessionManager(false)}
        onSessionSelect={handleSessionSelect}
      />
    </div>
  );
};

export default ChatPage;