// technicalToolsService.ts
/**
 * Uses fetch API instead of axios for compatibility with Node.js v18.
 */

export interface TechnicalAnalysisResult {
  id: string;
  type: 'analysis' | 'code_review' | 'security' | 'performance' | 'compliance';
  status: string;
  summary: string;
  details?: string;
  createdAt: string;
}

const API_BASE = '/api/technical-tools';

export async function fetchTechnicalAnalysis(): Promise<TechnicalAnalysisResult[]> {
  const res = await fetch(`${API_BASE}/analysis`);
  if (!res.ok) throw new Error('Failed to fetch technical analysis');
  return await res.json();
}

export async function fetchCodeReview(): Promise<TechnicalAnalysisResult[]> {
  const res = await fetch(`${API_BASE}/code-review`);
  if (!res.ok) throw new Error('Failed to fetch code review');
  return await res.json();
}

export async function fetchSecurityResults(): Promise<TechnicalAnalysisResult[]> {
  const res = await fetch(`${API_BASE}/security`);
  if (!res.ok) throw new Error('Failed to fetch security results');
  return await res.json();
}

export async function fetchPerformanceResults(): Promise<TechnicalAnalysisResult[]> {
  const res = await fetch(`${API_BASE}/performance`);
  if (!res.ok) throw new Error('Failed to fetch performance results');
  return await res.json();
}

export async function fetchComplianceResults(): Promise<TechnicalAnalysisResult[]> {
  const res = await fetch(`${API_BASE}/compliance`);
  if (!res.ok) throw new Error('Failed to fetch compliance results');
  return await res.json();
}