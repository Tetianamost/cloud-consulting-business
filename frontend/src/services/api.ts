// API service for backend communication
const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8061';

export interface CreateInquiryRequest {
  name: string;
  email: string;
  company?: string;
  phone?: string;
  services: string[];
  message: string;
  source?: string;
}

export interface InquiryResponse {
  success: boolean;
  data?: {
    id: string;
    name: string;
    email: string;
    company: string;
    phone: string;
    services: string[];
    message: string;
    status: string;
    priority: string;
    source: string;
    created_at: string;
    updated_at: string;
  };
  message?: string;
  error?: string;
}

export interface ServiceConfig {
  success: boolean;
  data: {
    services: Array<{
      id: string;
      name: string;
      description: string;
    }>;
  };
}

class ApiService {
  private baseUrl: string;

  constructor(baseUrl: string = API_BASE_URL) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    
    const config: RequestInit = {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    };

    try {
      const response = await fetch(url, config);
      
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
      }

      return await response.json();
    } catch (error) {
      console.error('API request failed:', error);
      throw error;
    }
  }

  // Create a new inquiry
  async createInquiry(data: CreateInquiryRequest): Promise<InquiryResponse> {
    return this.request<InquiryResponse>('/api/v1/inquiries', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  // Get service configuration
  async getServiceConfig(): Promise<ServiceConfig> {
    return this.request<ServiceConfig>('/api/v1/config/services');
  }

  // Health check
  async healthCheck(): Promise<{ status: string; service: string; version: string; time: string }> {
    return this.request('/health');
  }
}

export const apiService = new ApiService();
export default apiService;