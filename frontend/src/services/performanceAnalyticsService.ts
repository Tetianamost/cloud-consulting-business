// Performance Analytics Service for consultant performance tracking, QA/validation, and analytics

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8061';

import apiService from "./api";

export interface EngagementMetrics {
  total_engagements: number;
  successful_engagements: number;
  success_rate: number;
  average_client_satisfaction: number;
  average_project_duration: number;
  average_time_to_value: number;
  total_cost_savings: number;
  total_revenue_impact: number;
  top_recommendation_types: RecommendationMetric[];
  industry_breakdown: Record<string, number>;
  trend_data: MetricTrend[];
}

export interface RecommendationMetric {
  type: string;
  count: number;
  success_rate: number;
  average_rating: number;
  average_impact: number;
  implementation_rate: number;
}

export interface MetricTrend {
  period: string;
  value: number;
  change: number;
}

export interface TeamAnalytics {
  team_size: number;
  total_engagements: number;
  average_client_satisfaction: number;
  top_performers: ConsultantRanking[];
  skill_distribution: Record<string, SkillDistribution>;
  knowledge_sharing: KnowledgeSharingMetrics;
  team_trends: MetricTrend[];
  benchmark_metrics: BenchmarkMetrics;
  report_period: TimeRange;
  generated_at: string;
}

export interface ConsultantRanking {
  consultant_id: string;
  name: string;
  performance_score: number;
  client_satisfaction: number;
  project_success_rate: number;
  specializations: string[];
}

export interface SkillDistribution {
  skill_area: string;
  average_level: number;
  high_performers: number;
  needs_development: number;
  critical_gaps: number;
}

export interface KnowledgeSharingMetrics {
  total_knowledge_items: number;
  active_contributors: number;
  average_rating: number;
  most_viewed_items: string[];
  recent_contributions: number;
  knowledge_utilization: number;
}

export interface BenchmarkMetrics {
  industry_averages: Record<string, number>;
  company_targets: Record<string, number>;
  best_in_class_metrics: Record<string, number>;
  performance_thresholds: Record<string, number>;
  last_updated: string;
}

export interface TimeRange {
  start: string;
  end: string;
}

// Fetch engagement success metrics (optionally filtered)
export async function fetchEngagementMetrics(
  filters?: { consultant_id?: string; industry?: string; project_type?: string; status?: string; time_range?: TimeRange }
): Promise<EngagementMetrics> {
  const params = new URLSearchParams();
  if (filters) {
    Object.entries(filters).forEach(([key, value]) => {
      if (value) {
        if (typeof value === "object" && value !== null && "start" in value && "end" in value) {
          params.append("start", value.start);
          params.append("end", value.end);
        } else {
          params.append(key, String(value));
        }
      }
    });
  }
  const res = await fetch(`${API_BASE_URL}/performance/engagement-metrics?${params.toString()}`, {
    credentials: "include",
  });
  if (!res.ok) throw new Error("Failed to fetch engagement metrics");
  return res.json();
}

// Fetch team-wide analytics
export async function fetchTeamAnalytics(timeRange?: TimeRange): Promise<TeamAnalytics> {
  const params = timeRange ? `?start=${timeRange.start}&end=${timeRange.end}` : "";
  const res = await fetch(`${API_BASE_URL}/performance/team-analytics${params}`, {
    credentials: "include",
  });
  if (!res.ok) throw new Error("Failed to fetch team analytics");
  return res.json();
}

// Fetch performance benchmarks
export async function fetchBenchmarkMetrics(): Promise<BenchmarkMetrics> {
  const res = await fetch(`${API_BASE_URL}/performance/benchmarks`, {
    credentials: "include",
  });
  if (!res.ok) throw new Error("Failed to fetch benchmark metrics");
  return res.json();
}