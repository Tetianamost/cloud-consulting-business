// Quality Assurance Service for recommendation accuracy, peer review, client outcomes, and continuous improvement

import apiService from "./api";

// --- Types ---
export interface RecommendationAccuracy {
  id: string;
  consultant_id: string;
  recommendation_type: string;
  accuracy_score: number;
  evaluated_at: string;
}

export interface PeerReview {
  id: string;
  reviewer_id: string;
  reviewee_id: string;
  review_notes: string;
  score: number;
  created_at: string;
}

export interface ClientOutcome {
  id: string;
  client_id: string;
  project_id: string;
  outcome_type: string;
  outcome_value: string;
  measured_at: string;
}

export interface ContinuousImprovementItem {
  id: string;
  description: string;
  status: string;
  owner: string;
  created_at: string;
  updated_at: string;
}

// --- API Functions ---

// Fetch recommendation accuracy tracking data
const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8061';

// Use apiService for admin QA endpoints
// Use apiService for admin QA endpoints
export async function fetchRecommendationAccuracy(params?: { consultant_id?: string; date_from?: string; date_to?: string }) : Promise<RecommendationAccuracy[]> {
  const query = params
    ? "?" + new URLSearchParams(Object.entries(params).filter(([_, v]) => v !== undefined && v !== null && v !== "") as [string, string][])
    : "";
  const res = await apiService.getRecommendationAccuracy(query);
  return res.data as RecommendationAccuracy[];
}

// Fetch peer review data
// Use apiService for admin QA endpoints
// Use apiService for admin QA endpoints
export async function fetchPeerReviews(params?: { reviewee_id?: string; reviewer_id?: string; date_from?: string; date_to?: string }) : Promise<PeerReview[]> {
  const query = params
    ? "?" + new URLSearchParams(Object.entries(params).filter(([_, v]) => v !== undefined && v !== null && v !== "") as [string, string][])
    : "";
  const res = await apiService.getPeerReviews(query);
  return res.data as PeerReview[];
}

// Fetch client outcome tracking data
// Use apiService for admin QA endpoints
// Use apiService for admin QA endpoints
export async function fetchClientOutcomes(params?: { client_id?: string; project_id?: string; date_from?: string; date_to?: string }) : Promise<ClientOutcome[]> {
  const query = params
    ? "?" + new URLSearchParams(Object.entries(params).filter(([_, v]) => v !== undefined && v !== null && v !== "") as [string, string][])
    : "";
  const res = await apiService.getClientOutcomes(query);
  return res.data as ClientOutcome[];
}

// Fetch continuous improvement items
// Use apiService for admin QA endpoints
// Use apiService for admin QA endpoints
export async function fetchContinuousImprovementItems(params?: { status?: string; owner?: string }) : Promise<ContinuousImprovementItem[]> {
  const query = params
    ? "?" + new URLSearchParams(Object.entries(params).filter(([_, v]) => v !== undefined && v !== null && v !== "") as [string, string][])
    : "";
  const res = await apiService.getContinuousImprovementItems(query);
  return res.data as ContinuousImprovementItem[];
}