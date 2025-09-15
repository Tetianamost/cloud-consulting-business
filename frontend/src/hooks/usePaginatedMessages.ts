import { useState, useCallback, useEffect, useRef } from 'react';
import { ChatMessage } from '../store/slices/chatSlice';

interface UsePaginatedMessagesOptions {
  sessionId: string;
  pageSize?: number;
  initialLoad?: boolean;
}

interface PaginatedMessagesState {
  messages: ChatMessage[];
  isLoading: boolean;
  hasMore: boolean;
  error: string | null;
  totalCount: number;
}

interface PaginatedMessagesActions {
  loadMore: () => Promise<void>;
  refresh: () => Promise<void>;
  reset: () => void;
  addMessage: (message: ChatMessage) => void;
}

export const usePaginatedMessages = ({
  sessionId,
  pageSize = 50,
  initialLoad = true,
}: UsePaginatedMessagesOptions): [PaginatedMessagesState, PaginatedMessagesActions] => {
  const [state, setState] = useState<PaginatedMessagesState>({
    messages: [],
    isLoading: false,
    hasMore: true,
    error: null,
    totalCount: 0,
  });

  const offsetRef = useRef(0);
  const loadingRef = useRef(false);

  const loadMessages = useCallback(async (offset: number = 0, reset: boolean = false) => {
    if (loadingRef.current || !sessionId) return;

    loadingRef.current = true;
    setState(prev => ({ ...prev, isLoading: true, error: null }));

    try {
      const token = localStorage.getItem('adminToken');
      if (!token) {
        throw new Error('No authentication token found');
      }

      const response = await fetch(
        `/api/v1/admin/chat/sessions/${sessionId}/history?limit=${pageSize}&offset=${offset}`,
        {
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        }
      );

      if (!response.ok) {
        throw new Error(`Failed to load messages: ${response.statusText}`);
      }

      const data = await response.json();
      const newMessages: ChatMessage[] = data.messages || [];
      const totalCount = data.total_count || 0;

      setState(prev => ({
        ...prev,
        messages: reset ? newMessages : [...newMessages, ...prev.messages],
        hasMore: offset + newMessages.length < totalCount,
        totalCount,
        isLoading: false,
      }));

      offsetRef.current = reset ? newMessages.length : offset + newMessages.length;
    } catch (error) {
      console.error('Failed to load messages:', error);
      setState(prev => ({
        ...prev,
        error: error instanceof Error ? error.message : 'Failed to load messages',
        isLoading: false,
      }));
    } finally {
      loadingRef.current = false;
    }
  }, [sessionId, pageSize]);

  const loadMore = useCallback(async () => {
    if (!state.hasMore || state.isLoading) return;
    await loadMessages(offsetRef.current);
  }, [loadMessages, state.hasMore, state.isLoading]);

  const refresh = useCallback(async () => {
    offsetRef.current = 0;
    await loadMessages(0, true);
  }, [loadMessages]);

  const reset = useCallback(() => {
    setState({
      messages: [],
      isLoading: false,
      hasMore: true,
      error: null,
      totalCount: 0,
    });
    offsetRef.current = 0;
  }, []);

  const addMessage = useCallback((message: ChatMessage) => {
    setState(prev => ({
      ...prev,
      messages: [...prev.messages, message],
      totalCount: prev.totalCount + 1,
    }));
  }, []);

  // Initial load
  useEffect(() => {
    if (initialLoad && sessionId) {
      loadMessages(0, true);
    }
  }, [sessionId, initialLoad, loadMessages]);

  return [
    state,
    {
      loadMore,
      refresh,
      reset,
      addMessage,
    },
  ];
};

export default usePaginatedMessages;