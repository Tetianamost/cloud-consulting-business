import React, { useMemo, useCallback, useRef, useEffect } from 'react';
import { FixedSizeList as List } from 'react-window';
import { Bot, User } from 'lucide-react';
import { ChatMessage } from '../../store/slices/chatSlice';

interface VirtualizedMessageListProps {
  messages: ChatMessage[];
  height: number;
  isLoading?: boolean;
  onLoadMore?: () => void;
  hasMore?: boolean;
}

interface MessageItemProps {
  index: number;
  style: React.CSSProperties;
  data: {
    messages: ChatMessage[];
    formatTimestamp: (timestamp: string) => string;
  };
}

const MessageItem: React.FC<MessageItemProps> = ({ index, style, data }) => {
  const { messages, formatTimestamp } = data;
  const message = messages[index];

  if (!message) return null;

  return (
    <div style={style} className="px-4 py-2">
      <div
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
    </div>
  );
};

export const VirtualizedMessageList: React.FC<VirtualizedMessageListProps> = ({
  messages,
  height,
  isLoading = false,
  onLoadMore,
  hasMore = false,
}) => {
  const listRef = useRef<List>(null);
  const previousMessageCount = useRef(messages.length);

  const formatTimestamp = useCallback((timestamp: string) => {
    return new Date(timestamp).toLocaleTimeString([], { 
      hour: '2-digit', 
      minute: '2-digit' 
    });
  }, []);

  const itemData = useMemo(() => ({
    messages,
    formatTimestamp,
  }), [messages, formatTimestamp]);

  // Auto-scroll to bottom when new messages arrive
  useEffect(() => {
    if (messages.length > previousMessageCount.current && listRef.current) {
      listRef.current.scrollToItem(messages.length - 1, 'end');
    }
    previousMessageCount.current = messages.length;
  }, [messages.length]);

  // Handle scroll to load more messages
  const handleScroll = useCallback(({ scrollOffset }: { scrollOffset: number }) => {
    if (scrollOffset === 0 && hasMore && onLoadMore && !isLoading) {
      onLoadMore();
    }
  }, [hasMore, onLoadMore, isLoading]);

  if (messages.length === 0) {
    return (
      <div className="flex items-center justify-center h-full text-center text-gray-500 text-sm">
        <div>
          <Bot className="h-8 w-8 mx-auto mb-2 text-gray-400" />
          <p>Start a conversation to get real-time AWS consulting assistance during your client meeting.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="relative">
      {isLoading && (
        <div className="absolute top-0 left-0 right-0 z-10 bg-blue-50 border-b border-blue-200 p-2 text-center text-sm text-blue-600">
          Loading more messages...
        </div>
      )}
      <List
        ref={listRef}
        height={height}
        width="100%"
        itemCount={messages.length}
        itemSize={120} // Estimated height per message
        itemData={itemData}
        onScroll={handleScroll}
        overscanCount={5} // Render 5 extra items for smooth scrolling
      >
        {MessageItem}
      </List>
    </div>
  );
};

export default VirtualizedMessageList;