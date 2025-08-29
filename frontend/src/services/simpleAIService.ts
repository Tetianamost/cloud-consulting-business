import apiService from "./api";

interface SimpleAIRequest {
  message: string;
  context?: {
    clientName?: string;
    meetingType?: string;
  };
}

interface SimpleAIResponse {
  content: string;
  timestamp: string;
}

interface SimpleChatMessage {
  id: string;
  content: string;
  role: "user" | "assistant";
  timestamp: string;
  session_id: string;
}

class EnhancedAIService {
  private static instance: EnhancedAIService;
  private isConnected = false;
  private sessionId: string;
  private baseUrl: string;

  constructor() {
    this.sessionId = `session-${Date.now()}-${Math.random()
      .toString(36)
      .substring(2, 9)}`;
    this.baseUrl = process.env.REACT_APP_API_URL || "http://localhost:8061";
  }

  static getInstance(): EnhancedAIService {
    if (!EnhancedAIService.instance) {
      EnhancedAIService.instance = new EnhancedAIService();
    }
    return EnhancedAIService.instance;
  }

  async checkConnection(): Promise<boolean> {
    try {
      console.log('[SimpleAIService] Checking connection to:', this.baseUrl);
      
      // First check if we have an auth token
      const token = localStorage.getItem("adminToken");
      if (!token) {
        console.warn('[SimpleAIService] No auth token found');
        this.isConnected = false;
        return false;
      }

      // Test the actual chat API endpoint that requires authentication
      const response = await fetch(
        `${this.baseUrl}/api/v1/admin/simple-chat/messages?session_id=health-check`,
        {
          method: "GET",
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (response.ok) {
        this.isConnected = true;
        console.log('[SimpleAIService] Connection successful');
        return true;
      } else {
        console.error('[SimpleAIService] API responded with status:', response.status);
        this.isConnected = false;
        return false;
      }
    } catch (error) {
      console.error('[SimpleAIService] Connection failed:', error);
      this.isConnected = false;
      return false;
    }
  }

  isHealthy(): boolean {
    return this.isConnected;
  }

  async sendMessage(request: SimpleAIRequest): Promise<SimpleAIResponse> {
    try {
      const token = localStorage.getItem("adminToken");
      if (!token) {
        throw new Error("Authentication required");
      }

      // Send message to the simple chat endpoint
      const response = await fetch(
        `${this.baseUrl}/api/v1/admin/simple-chat/messages`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            content: request.message,
            session_id: this.sessionId,
          }),
        }
      );

      if (!response.ok) {
        if (response.status === 401) {
          throw new Error("Authentication failed");
        }
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const result = await response.json();

      if (!result.success) {
        throw new Error(result.error || "Failed to send message");
      }

      // Get the latest messages to find the AI response
      const messagesResponse = await this.getMessages();
      const aiMessage = messagesResponse.find(
        (msg) => msg.role === "assistant" && msg.session_id === this.sessionId
      );

      if (aiMessage) {
        this.isConnected = true;
        return {
          content: aiMessage.content,
          timestamp: aiMessage.timestamp,
        };
      } else {
        throw new Error("No AI response received");
      }
    } catch (error) {
      console.error("SimpleAIService: Error sending message:", error);
      this.isConnected = false;
      throw error;
    }
  }

  async getMessages(): Promise<SimpleChatMessage[]> {
    try {
      const token = localStorage.getItem("adminToken");
      if (!token) {
        throw new Error("Authentication required");
      }

      const response = await fetch(
        `${this.baseUrl}/api/v1/admin/simple-chat/messages?session_id=${this.sessionId}`,
        {
          method: "GET",
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const result = await response.json();

      if (!result.success) {
        throw new Error(result.error || "Failed to get messages");
      }

      return result.messages || [];
    } catch (error) {
      console.error("SimpleAIService: Error getting messages:", error);
      throw error;
    }
  }

  getSessionId(): string {
    return this.sessionId;
  }

  resetSession(): void {
    this.sessionId = `session-${Date.now()}-${Math.random()
      .toString(36)
      .substring(2, 9)}`;
  }

  async forceReconnect(): Promise<boolean> {
    this.isConnected = false;
    return await this.checkConnection();
  }
}

export const enhancedAIService = EnhancedAIService.getInstance();
export default enhancedAIService;
