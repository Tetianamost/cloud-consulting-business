import React, { useEffect, useState } from "react";
import { fetchCostAnalysis, CostAnalysis } from "../../services/costAnalysisService";
import { Card, CardContent, CardHeader, CardTitle } from "../ui/card";
import { Button } from "../ui/button";
import { AlertCircle, RefreshCw } from "lucide-react";

const CostAnalysisPage: React.FC = () => {
  const [analyses, setAnalyses] = useState<CostAnalysis[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadAnalyses = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await fetchCostAnalysis();
      setAnalyses(data);
    } catch (err: any) {
      setError(err.message || "Failed to load cost analysis data");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadAnalyses();
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="text-center">
          <RefreshCw className="w-8 h-8 animate-spin mx-auto mb-4 text-blue-600" />
          <p className="text-gray-600">Loading cost analysis...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="text-center">
          <AlertCircle className="w-8 h-8 mx-auto mb-4 text-red-600" />
          <p className="text-red-600 mb-4">{error}</p>
          <Button onClick={loadAnalyses} variant="outline">
            <RefreshCw className="w-4 h-4 mr-2" />
            Retry
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Cost Analysis</h1>
          <p className="text-gray-600">View and analyze cloud cost breakdowns and optimization insights</p>
        </div>
        <Button onClick={loadAnalyses} variant="outline">
          <RefreshCw className="w-4 h-4 mr-2" />
          Refresh
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Cost Analysis Results</CardTitle>
        </CardHeader>
        <CardContent>
          {analyses.length === 0 ? (
            <p className="text-gray-600">No cost analysis data found.</p>
          ) : (
            <ul className="divide-y">
              {analyses.map((analysis) => (
                <li key={analysis.id} className="py-3">
                  <div className="font-semibold">
                    {new Date(analysis.analysisDate).toLocaleString()} &mdash; {analysis.currency} {analysis.totalMonthlyCost.toLocaleString()} /mo
                  </div>
                  <div className="text-sm text-gray-500">Annual: {analysis.currency} {analysis.totalAnnualCost.toLocaleString()}</div>
                  <div className="text-sm mt-1">Methodology: {analysis.methodology}</div>
                  <div className="text-sm mt-1">Top Cost Drivers: {analysis.costDrivers?.slice(0, 3).join(", ")}</div>
                  <div className="text-sm mt-1">Optimization Potential: {analysis.optimizationPotential}</div>
                </li>
              ))}
            </ul>
          )}
        </CardContent>
      </Card>
    </div>
  );
};

export default CostAnalysisPage;