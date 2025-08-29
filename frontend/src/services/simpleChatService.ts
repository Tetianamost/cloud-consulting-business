/**
 * Simple Chat Service - A working chat implementation that just works!
 * No complex polling, no WebSocket issues - just simple HTTP requests that work.
 */

interface SimpleChatMessage {
  id: string;
  content: string;
  role: 'user' | 'assistant';
  timestamp: string;
  session_id: string;
}

interface SendMessageRequest {
  content: string;
  session_id: string;
}

interface SendMessageResponse {
  success: boolean;
  message_id: string;
  error?: string;
}

interface GetMessagesResponse {
  success: boolean;
  messages: SimpleChatMessage[];
}

class SimpleChatService {
  private baseUrl = '/api/v1/admin/simple-chat';
  private sessionId: string;

  constructor() {
    // Generate a simple session ID
    this.sessionId = `session-${Date.now()}-${Math.random().toString(36).substring(2, 9)}`;
    console.log('[SimpleChat] Initialized with session:', this.sessionId);
  }

  /**
   * Send a message and get immediate response
   */
  async sendMessage(content: string): Promise<SimpleChatMessage[]> {
    try {
      console.log('[SimpleChat] Sending message:', content);

      const token = localStorage.getItem('token');
      if (!token) {
        throw new Error('No authentication token found');
      }

      // Send the message
      const response = await fetch(`${this.baseUrl}/messages`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({
          content,
          session_id: this.sessionId,
        } as SendMessageRequest),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const result: SendMessageResponse = await response.json();
      
      if (!result.success) {
        throw new Error(result.error || 'Failed to send message');
      }

      console.log('[SimpleChat] Message sent successfully:', result.message_id);

      // Immediately get all messages (including the AI response)
      return await this.getMessages();
    } catch (error) {
      console.error('[SimpleChat] Error sending message:', error);
      throw error;
    }
  }

  /**
   * Get all messages for the current session
   */
  async getMessages(): Promise<SimpleChatMessage[]> {
    try {
      const token = localStorage.getItem('token');
      if (!token) {
        throw new Error('No authentication token found');
      }

      const response = await fetch(`${this.baseUrl}/messages?session_id=${this.sessionId}`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const result: GetMessagesResponse = await response.json();
      
      if (!result.success) {
        throw new Error('Failed to get messages');
      }

      console.log('[SimpleChat] Retrieved messages:', result.messages.length);
      return result.messages;
    } catch (error) {
      console.error('[SimpleChat] Error getting messages:', error);
      throw error;
    }
  }

  /**
   * Get the current session ID
   */
  getSessionId(): string {
    return this.sessionId;
  }

  /**
   * Reset the session (start a new conversation)
   */
  resetSession(): void {
    this.sessionId = `session-${Date.now()}-${Math.random().toString(36).substring(2, 9)}`;
    console.log('[SimpleChat] Reset to new session:', this.sessionId);
  }
}

// Export a singleton instance
export const simpleChatService = new SimpleChatService();
export default SimpleChatService;