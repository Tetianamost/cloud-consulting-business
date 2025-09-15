import React, { useState, useRef, useEffect } from 'react';
import { Send, MessageSquare, User, Bot, Maximize2, Minimize2, X, AlertCircle } from 'lucide-react';
import simpleAIService from '../../services/simpleAIService';
import MarkdownRenderer from '../../utils/markdownRenderer';

interface Message {
  id: string;
  type: 'user' | 'assistant';
  content: string;
  timestamp: string;
}

interface SimpleAIWidgetProps {
  isMinimized?: boolean;
  onToggleMinimize?: () => void;
  onClose?: () => void;
}

export const SimpleAIWidget: React.FC<SimpleAIWidgetProps> = ({
  isMinimized = false,
  onToggleMinimize,
  onClose
}) => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [inputValue, setInputValue] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [isConnected, setIsConnected] = useState(true); // Start optimistic
  const [clientName, setClientName] = useState('');
  const [showSettings, setShowSettings] = useState(false);
  
  const inputRef = useRef<HTMLInputElement>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  // Check connection on mount
  useEffect(() => {
    const checkConnection = async () => {
      try {
        const connected = await simpleAIService.checkConnection();
        setIsConnected(connected);
      } catch (error) {
        setIsConnected(false);
      }
    };
    checkConnection();
  }, []);

  // Auto-scroll to bottom when new messages arrive
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const sendMessage = async (message: string) => {
    if (!message.trim() || isLoading) return;

    const userMessage: Message = {
      id: `user-${Date.now()}`,
      type: 'user',
      content: message.trim(),
      timestamp: new Date().toISOString()
    };

    setMessages(prev => [...prev, userMessage]);
    setInputValue('');
    setIsLoading(true);

    try {
      const response = await simpleAIService.sendMessage({
        message: message.trim(),
        context: {
          clientName: clientName || undefined
        }
      });

      const aiMessage: Message = {
        id: `ai-${Date.now()}`,
        type: 'assistant',
        content: response.content,
        timestamp: response.timestamp
      };

      setMessages(prev => [...prev, aiMessage]);
    } catch (error) {
      const errorMessage: Message = {
        id: `error-${Date.now()}`,
        type: 'assistant',
        content: 'Sorry, I encountered an error. Please try again.',
        timestamp: new Date().toISOString()
      };
      setMessages(prev => [...prev, errorMessage]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    sendMessage(inputValue);
  };

  const formatTimestamp = (timestamp: string) => {
    return new Date(timestamp).toLocaleTimeString([], { 
      hour: '2-digit', 
      minute: '2-digit' 
    });
  };

  const clearChat = () => {
    setMessages([]);
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
    <div className="fixed bottom-4 right-4 w-80 h-[500px] bg-white rounded-lg shadow-xl border border-gray-200 flex flex-col z-50">
      {/* Header */}
      <div className="flex items-center justify-between p-2 border-b border-gray-200 bg-blue-600 text-white rounded-t-lg flex-shrink-0">
        <div className="flex items-center space-x-2 min-w-0">
          <MessageSquare className="h-4 w-4 flex-shrink-0" />
          <div className="flex flex-col min-w-0">
            <h3 className="font-semibold text-xs truncate">AI Assistant</h3>
            <span className="text-xs opacity-75 truncate">
              {isConnected ? 'Ready' : 'Offline'}
            </span>
          </div>
          <div className={`w-2 h-2 rounded-full flex-shrink-0 ${isConnected ? 'bg-green-400' : 'bg-red-400'}`} />
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

      {/* Settings */}
      {showSettings && (
        <div className="p-2 border-b border-gray-200 bg-gray-50 flex-shrink-0">
          <input
            type="text"
            value={clientName}
            onChange={(e) => setClientName(e.target.value)}
            placeholder="Client name (optional)"
            className="w-full px-2 py-1 text-xs border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-blue-500"
          />
        </div>
      )}

      {/* Quick Actions */}
      <div className="p-2 border-b border-gray-200 bg-gray-50 flex-shrink-0">
        <div className="flex flex-wrap gap-1">
          <button
            onClick={() => sendMessage('Provide a cost estimate for this solution')}
            disabled={isLoading}
            className="px-2 py-1 text-xs bg-blue-100 text-blue-700 rounded hover:bg-blue-200 disabled:opacity-50"
          >
            Cost
          </button>
          <button
            onClick={() => sendMessage('What are the security considerations?')}
            disabled={isLoading}
            className="px-2 py-1 text-xs bg-blue-100 text-blue-700 rounded hover:bg-blue-200 disabled:opacity-50"
          >
            Security
          </button>
          <button
            onClick={() => sendMessage('What are the best practices?')}
            disabled={isLoading}
            className="px-2 py-1 text-xs bg-blue-100 text-blue-700 rounded hover:bg-blue-200 disabled:opacity-50"
          >
            Best Practices
          </button>
        </div>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-hidden relative">
        <div className="h-full overflow-y-auto p-2 space-y-2">
          {!isConnected && (
            <div className="bg-yellow-50 border border-yellow-200 rounded p-2 mb-2">
              <div className="flex items-center space-x-1">
                <AlertCircle className="h-3 w-3 text-yellow-600 flex-shrink-0" />
                <p className="text-xs text-yellow-700">Working offline with demo responses</p>
              </div>
            </div>
          )}

          {messages.length === 0 && (
            <div className="text-center text-gray-500 text-xs py-4">
              <MessageSquare className="h-6 w-6 mx-auto mb-2 text-gray-400" />
              <p className="px-2 mb-2">Ask me about AWS costs, security, architecture, or best practices.</p>
              <button
                onClick={() => window.open('/admin/ai-consultant', '_blank')}
                className="text-blue-600 hover:text-blue-800 underline text-xs"
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
                    {message.type === 'assistant' ? (
                      <MarkdownRenderer 
                        content={message.content.length > 500 ? `${message.content.substring(0, 500)}...` : message.content}
                        className="text-xs"
                      />
                    ) : (
                      <p className="whitespace-pre-wrap break-words leading-tight text-xs">
                        {message.content}
                      </p>
                    )}
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
          
          <div ref={messagesEndRef} />
        </div>
      </div>

      {/* Input */}
      <form onSubmit={handleSubmit} className="p-2 border-t border-gray-200 flex-shrink-0">
        <div className="flex space-x-2">
          <input
            ref={inputRef}
            type="text"
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            placeholder="Ask about AWS..."
            disabled={isLoading}
            className="flex-1 px-2 py-1 text-xs border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-blue-500 disabled:opacity-50"
          />
          <button
            type="submit"
            disabled={isLoading || !inputValue.trim()}
            className="px-2 py-1 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors flex-shrink-0"
          >
            <Send className="h-3 w-3" />
          </button>
        </div>
      </form>
    </div>
  );
};

export default SimpleAIWidget;