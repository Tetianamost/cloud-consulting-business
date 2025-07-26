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

export interface Inquiry {
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
  reports?: Report[];
}

export interface Report {
  id: string;
  inquiry_id: string;
  type: string;
  title: string;
  content: string;
  status: string;
  generated_by: string;
  reviewed_by?: string;
  s3_key: string;
  created_at: string;
  updated_at: string;
}

export interface AdminInquiriesResponse {
  success: boolean;
  data: Inquiry[];
  count: number;
  total: number;
  page: number;
  pages: number;
}

export interface SystemMetrics {
  total_inquiries: number;
  reports_generated: number;
  emails_sent: number;
  email_delivery_rate: number;
  avg_report_gen_time_ms: number;
  system_uptime: string;
  last_processed_at?: string;
}

export interface EmailStatus {
  inquiry_id: string;
  customer_email: string;
  consultant_email: string;
  status: string;
  sent_at?: string;
  delivered_at?: string;
  error_message?: string;
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

  // Admin endpoints
  
  // List inquiries with filtering and pagination
  async listInquiries(filters?: {
    status?: string;
    priority?: string;
    service?: string;
    date_from?: string;
    date_to?: string;
    limit?: number;
    offset?: number;
  }): Promise<AdminInquiriesResponse> {
    const queryParams = new URLSearchParams();
    
    if (filters) {
      Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== '') {
          queryParams.append(key, String(value));
        }
      });
    }
    
    const queryString = queryParams.toString() ? `?${queryParams.toString()}` : '';
    return this.request<AdminInquiriesResponse>(`/api/v1/admin/inquiries${queryString}`);
  }
  
  // Get system metrics
  async getSystemMetrics(): Promise<{ success: boolean; data: SystemMetrics }> {
    return this.request<{ success: boolean; data: SystemMetrics }>('/api/v1/admin/metrics');
  }
  
  // Get email delivery status for an inquiry
  async getEmailStatus(inquiryId: string): Promise<{ success: boolean; data: EmailStatus }> {
    return this.request<{ success: boolean; data: EmailStatus }>(`/api/v1/admin/email-status/${inquiryId}`);
  }
  
  // Download report in specified format
  async downloadReport(inquiryId: string, format: 'pdf' | 'html'): Promise<Blob> {
    const url = `${this.baseUrl}/api/v1/admin/reports/${inquiryId}/download/${format}`;
    
    const response = await fetch(url);
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
    }
    
    return await response.blob();
  }
}

export const apiService = new ApiService();
export default apiService;