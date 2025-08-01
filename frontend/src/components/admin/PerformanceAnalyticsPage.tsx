// Performance Analytics Page: Fetches and displays consultant/team analytics using the backend API

import React, { useEffect, useState } from "react";
import {
  fetchEngagementMetrics,
  fetchTeamAnalytics,
  fetchBenchmarkMetrics,
  EngagementMetrics,
  TeamAnalytics,
  BenchmarkMetrics,
  TimeRange,
} from "../../services/performanceAnalyticsService";

const defaultTimeRange: TimeRange = {
  start: new Date(new Date().setMonth(new Date().getMonth() - 1)).toISOString(),
  end: new Date().toISOString(),
};

const PerformanceAnalyticsPage: React.FC = () => {
  const [engagementMetrics, setEngagementMetrics] = useState<EngagementMetrics | null>(null);
  const [teamAnalytics, setTeamAnalytics] = useState<TeamAnalytics | null>(null);
  const [benchmarks, setBenchmarks] = useState<BenchmarkMetrics | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setLoading(true);
    setError(null);
    Promise.all([
      fetchEngagementMetrics({ time_range: defaultTimeRange }),
      fetchTeamAnalytics(defaultTimeRange),
      fetchBenchmarkMetrics(),
    ])
      .then(([engagement, team, benchmark]) => {
        setEngagementMetrics(engagement);
        setTeamAnalytics(team);
        setBenchmarks(benchmark);
      })
      .catch((err) => {
        setError(err.message || "Failed to load analytics data");
      })
      .finally(() => setLoading(false));
  }, []);

  // Only backend integration/data fetching logic is implemented here.
  // UI rendering is minimal and for demonstration/testing only.
  return (
    <div>
      <h1>Performance Analytics</h1>
      {loading && <div>Loading analytics...</div>}
      {error && <div style={{ color: "red" }}>Error: {error}</div>}
      {!loading && !error && (
        <pre style={{ background: "#f5f5f5", padding: 16, borderRadius: 8, fontSize: 13 }}>
          {JSON.stringify(
            {
              engagementMetrics,
              teamAnalytics,
              benchmarks,
            },
            null,
            2
          )}
        </pre>
      )}
    </div>
  );
};

export default PerformanceAnalyticsPage;