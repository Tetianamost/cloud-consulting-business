import React, { useState, useEffect, useRef, useCallback } from 'react';
import { Send, MessageSquare, User, Bot, Clock, Minimize2, Maximize2, X } from 'lucide-react';

interface ChatMessage {
  id: string;
  type: 'user' | 'assistant' | 'system';
  content: string;
  timestamp: string;
  session_id: string;
}

interface ChatSession {
  id: string;
  consultant_id: string;
  client_name?: string;
  context?: string;
  messages: ChatMessage[];
  created_at: string;
  last_activity: string;
}

interface ChatRequest {
  message: string;
  session_id?: string;
  client_name?: string;
  context?: string;
  quick_action?: string;
}

interface ChatResponse {
  message: ChatMessage;
  session_id: string;
  success: boolean;
  error?: string;
}

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
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [inputMessage, setInputMessage] = useState('');
  const [clientName, setClientName] = useState('');
  const [meetingContext, setMeetingContext] = useState('');
  const [sessionId, setSessionId] = useState<string | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [showSettings, setShowSettings] = useState(false);
  
  const wsRef = useRef<WebSocket | null>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  // Auto-scroll to bottom when new messages arrive
  const scrollToBottom = useCallback(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, []);

  useEffect(() => {
    scrollToBottom();
  }, [messages, scrollToBottom]);

  // Initialize WebSocket connection
  useEffect(() => {
    connectWebSocket();
    return () => {
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, []);

  const connectWebSocket = () => {
    const token = localStorage.getItem('adminToken');
    if (!token) {
      console.error('No admin token found');
      return;
    }

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/api/v1/admin/chat/ws`;
    
    try {
      wsRef.current = new WebSocket(wsUrl);
      
      wsRef.current.onopen = () => {
        console.log('WebSocket connected');
        setIsConnected(true);
      };
      
      wsRef.current.onmessage = (event) => {
        try {
          const response: ChatResponse = JSON.parse(event.data);
          if (response.success) {
            setMessages(prev => [...prev, response.message]);
            setSessionId(response.session_id);
          } else {
            console.error('Chat error:', response.error);
          }
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error);
        }
        setIsLoading(false);
      };
      
      wsRef.current.onclose = () => {
        console.log('WebSocket disconnected');
        setIsConnected(false);
        // Attempt to reconnect after 3 seconds
        setTimeout(connectWebSocket, 3000);
      };
      
      wsRef.current.onerror = (error) => {
        console.error('WebSocket error:', error);
        setIsConnected(false);
      };
    } catch (error) {
      console.error('Failed to create WebSocket connection:', error);
    }
  };

  const sendMessage = (message: string, quickAction?: string) => {
    if (!wsRef.current || wsRef.current.readyState !== WebSocket.OPEN) {
      console.error('WebSocket not connected');
      return;
    }

    if (!message.trim()) return;

    const request: ChatRequest = {
      message: message.trim(),
      session_id: sessionId || undefined,
      client_name: clientName || undefined,
      context: meetingContext || undefined,
      quick_action: quickAction || undefined,
    };

    // Add user message to UI immediately
    const userMessage: ChatMessage = {
      id: `temp-${Date.now()}`,
      type: 'user',
      content: message.trim(),
      timestamp: new Date().toISOString(),
      session_id: sessionId || '',
    };
    
    setMessages(prev => [...prev, userMessage]);
    setInputMessage('');
    setIsLoading(true);

    // Send to WebSocket
    wsRef.current.send(JSON.stringify(request));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    sendMessage(inputMessage);
  };

  const handleQuickAction = (action: typeof QUICK_ACTIONS[0]) => {
    sendMessage(action.prompt, action.id);
  };

  const formatTimestamp = (timestamp: string) => {
    return new Date(timestamp).toLocaleTimeString([], { 
      hour: '2-digit', 
      minute: '2-digit' 
    });
  };

  const clearChat = () => {
    setMessages([]);
    setSessionId(null);
  };

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
    <div className="fixed bottom-4 right-4 w-96 h-[600px] bg-white rounded-lg shadow-xl border border-gray-200 flex flex-col z-50">
      {/* Header */}
      <div className="flex items-center justify-between p-4 border-b border-gray-200 bg-blue-600 text-white rounded-t-lg">
        <div className="flex items-center space-x-2">
          <MessageSquare className="h-5 w-5" />
          <h3 className="font-semibold">Consultant Assistant</h3>
          <div className={`w-2 h-2 rounded-full ${isConnected ? 'bg-green-400' : 'bg-red-400'}`} />
        </div>
        <div className="flex items-center space-x-2">
          <button
            onClick={() => setShowSettings(!showSettings)}
            className="text-white hover:text-gray-200 transition-colors"
          >
            <Clock className="h-4 w-4" />
          </button>
          {onToggleMinimize && (
            <button
              onClick={onToggleMinimize}
              className="text-white hover:text-gray-200 transition-colors"
            >
              <Minimize2 className="h-4 w-4" />
            </button>
          )}
          {onClose && (
            <button
              onClick={onClose}
              className="text-white hover:text-gray-200 transition-colors"
            >
              <X className="h-4 w-4" />
            </button>
          )}
        </div>
      </div>

      {/* Settings Panel */}
      {showSettings && (
        <div className="p-3 border-b border-gray-200 bg-gray-50 space-y-2">
          <div>
            <label className="block text-xs font-medium text-gray-700 mb-1">
              Client Name
            </label>
            <input
              type="text"
              value={clientName}
              onChange={(e) => setClientName(e.target.value)}
              placeholder="Enter client name..."
              className="w-full px-2 py-1 text-sm border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-blue-500"
            />
          </div>
          <div>
            <label className="block text-xs font-medium text-gray-700 mb-1">
              Meeting Context
            </label>
            <input
              type="text"
              value={meetingContext}
              onChange={(e) => setMeetingContext(e.target.value)}
              placeholder="e.g., Migration planning, Cost optimization..."
              className="w-full px-2 py-1 text-sm border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-blue-500"
            />
          </div>
          <button
            onClick={clearChat}
            className="text-xs text-red-600 hover:text-red-800 transition-colors"
          >
            Clear Chat History
          </button>
        </div>
      )}

      {/* Quick Actions */}
      <div className="p-2 border-b border-gray-200 bg-gray-50">
        <div className="flex flex-wrap gap-1">
          {QUICK_ACTIONS.map((action) => (
            <button
              key={action.id}
              onClick={() => handleQuickAction(action)}
              disabled={isLoading || !isConnected}
              className="px-2 py-1 text-xs bg-blue-100 text-blue-700 rounded hover:bg-blue-200 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {action.label}
            </button>
          ))}
        </div>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {messages.length === 0 && (
          <div className="text-center text-gray-500 text-sm">
            <MessageSquare className="h-8 w-8 mx-auto mb-2 text-gray-400" />
            <p>Start a conversation to get real-time AWS consulting assistance during your client meeting.</p>
          </div>
        )}
        
        {messages.map((message) => (
          <div
            key={message.id}
            className={`flex ${message.type === 'user' ? 'justify-end' : 'justify-start'}`}
          >
            <div
              className={`max-w-[80%] rounded-lg p-3 ${
                message.type === 'user'
                  ? 'bg-blue-600 text-white'
                  : 'bg-gray-100 text-gray-900'
              }`}
            >
              <div className="flex items-start space-x-2">
                {message.type === 'assistant' && (
                  <Bot className="h-4 w-4 mt-0.5 text-blue-600" />
                )}
                {message.type === 'user' && (
                  <User className="h-4 w-4 mt-0.5 text-white" />
                )}
                <div className="flex-1">
                  <p className="text-sm whitespace-pre-wrap">{message.content}</p>
                  <p className={`text-xs mt-1 ${
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
            <div className="bg-gray-100 rounded-lg p-3 max-w-[80%]">
              <div className="flex items-center space-x-2">
                <Bot className="h-4 w-4 text-blue-600" />
                <div className="flex space-x-1">
                  <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" />
                  <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '0.1s' }} />
                  <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '0.2s' }} />
                </div>
              </div>
            </div>
          </div>
        )}
        
        <div ref={messagesEndRef} />
      </div>

      {/* Input */}
      <form onSubmit={handleSubmit} className="p-4 border-t border-gray-200">
        <div className="flex space-x-2">
          <input
            ref={inputRef}
            type="text"
            value={inputMessage}
            onChange={(e) => setInputMessage(e.target.value)}
            placeholder={isConnected ? "Ask about AWS services, costs, best practices..." : "Connecting..."}
            disabled={isLoading || !isConnected}
            className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
          />
          <button
            type="submit"
            disabled={isLoading || !isConnected || !inputMessage.trim()}
            className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            <Send className="h-4 w-4" />
          </button>
        </div>
      </form>
    </div>
  );
};

export default ConsultantChat;