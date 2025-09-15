// Service for Cost Analysis API integration

export interface CostAnalysis {
  id: string;
  analysisDate: string;
  totalMonthlyCost: number;
  totalAnnualCost: number;
  currency: string;
  serviceBreakdown: Record<string, number>;
  categoryBreakdown: Record<string, number>;
  regionBreakdown: Record<string, number>;
  costDrivers: string[];
  costTrends: string[];
  benchmarkComparison: string;
  optimizationPotential: string;
  assumptions: string[];
  methodology: string;
}

// Use apiService for admin endpoints
import apiService from "./api";

export async function fetchCostAnalysis(): Promise<CostAnalysis[]> {
// Use apiService for admin endpoints
  return apiService.getCostAnalysis();
}