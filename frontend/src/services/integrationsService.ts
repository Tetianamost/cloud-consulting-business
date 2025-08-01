 // Integrations service for backend API communication

import { Integration, IntegrationListResponse } from "../types/integrations";
export type { Integration, IntegrationListResponse };

// Use apiService for admin endpoints
import apiService from "./api";

export async function fetchIntegrations(): Promise<IntegrationListResponse> {
  return apiService.getIntegrations();
}

// Use apiService for admin endpoints
export async function fetchIntegrationData(
  integrationId: string,
  dataType?: string
): Promise<any> {
  return apiService.getIntegrationData(integrationId, dataType);
}

// Use apiService for admin endpoints
export async function testIntegration(
  integrationId: string
): Promise<{ success: boolean; result: any }> {
  return apiService.testIntegration(integrationId);
}