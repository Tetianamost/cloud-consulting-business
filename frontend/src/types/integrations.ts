// Shared types for integrations

export interface Integration {
  id: string;
  name: string;
  type: string;
  status: string;
  description?: string;
  config?: Record<string, any>;
  lastSync?: string;
}

export interface IntegrationListResponse {
  success: boolean;
  data: Integration[];
}