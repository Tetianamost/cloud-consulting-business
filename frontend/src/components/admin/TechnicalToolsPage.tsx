// TechnicalToolsPage.tsx
import React, { useEffect, useState } from 'react';
import {
  fetchTechnicalAnalysis,
  fetchCodeReview,
  fetchSecurityResults,
  fetchPerformanceResults,
  fetchComplianceResults,
  TechnicalAnalysisResult,
} from '../../services/technicalToolsService';

const TechnicalToolsPage: React.FC = () => {
  const [analysis, setAnalysis] = useState<TechnicalAnalysisResult[]>([]);
  const [codeReview, setCodeReview] = useState<TechnicalAnalysisResult[]>([]);
  const [security, setSecurity] = useState<TechnicalAnalysisResult[]>([]);
  const [performance, setPerformance] = useState<TechnicalAnalysisResult[]>([]);
  const [compliance, setCompliance] = useState<TechnicalAnalysisResult[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setLoading(true);
    Promise.all([
      fetchTechnicalAnalysis(),
      fetchCodeReview(),
      fetchSecurityResults(),
      fetchPerformanceResults(),
      fetchComplianceResults(),
    ])
      .then(([analysis, codeReview, security, performance, compliance]) => {
        setAnalysis(analysis);
        setCodeReview(codeReview);
        setSecurity(security);
        setPerformance(performance);
        setCompliance(compliance);
        setLoading(false);
      })
      .catch((err) => {
        setError('Failed to fetch technical tools data.');
        setLoading(false);
      });
  }, []);

  if (loading) return <div>Loading technical tools data...</div>;
  if (error) return <div>{error}</div>;

  // Placeholder: Replace with actual UI rendering as needed
  return (
    <div>
      <h1>Technical Tools Results</h1>
      <section>
        <h2>Technical Analysis</h2>
        <pre>{JSON.stringify(analysis, null, 2)}</pre>
      </section>
      <section>
        <h2>Code Review</h2>
        <pre>{JSON.stringify(codeReview, null, 2)}</pre>
      </section>
      <section>
        <h2>Security</h2>
        <pre>{JSON.stringify(security, null, 2)}</pre>
      </section>
      <section>
        <h2>Performance</h2>
        <pre>{JSON.stringify(performance, null, 2)}</pre>
      </section>
      <section>
        <h2>Compliance</h2>
        <pre>{JSON.stringify(compliance, null, 2)}</pre>
      </section>
    </div>
  );
};

export default TechnicalToolsPage;