import React, { useState, useEffect } from 'react';
import { 
  FileText, 
  TrendingUp, 
  AlertTriangle, 
  Clock, 
  DollarSign,
  Eye,
  Download,
  Plus,
  ChevronDown,
  ChevronUp,
  Target,
  CheckCircle,
  XCircle
} from 'lucide-react';
import { Inquiry } from '../../services/api';
import V0ReportModal from './V0ReportModal';
import { V0DataAdapter } from './V0DataAdapter';

// Types for analysis reports
export interface RecommendedAction {
  id: string;
  title: string;
  priority: 'High' | 'Medium' | 'Low';
  description: string;
  estimatedImpact: string;
  completed: boolean;
}

export interface AnalysisReport {
  id: string;
  title: string;
  customer: string;
  service: string;
  value: string;
  timeline: string;
  confidence: number;
  risk: 'High' | 'Medium' | 'Low';
  insights: string[];
  actions: RecommendedAction[];
  generatedAt: string;
  inquiryId: string;
}

interface V0InquiryAnalysisSectionProps {
  inquiries: Inquiry[];
  onGenerateReport?: (inquiryId: string) => void;
  onViewReport?: (reportId: string) => void;
  onDownloadReport?: (reportId: string) => void;
}

const V0InquiryAnalysisSection: React.FC<V0InquiryAnalysisSectionProps> = ({
  inquiries,
  onGenerateReport,
  onViewReport,
  onDownloadReport
}) => {
  const [reports, setReports] = useState<AnalysisReport[]>([]);
  const [expandedReport, setExpandedReport] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [selectedReport, setSelectedReport] = useState<AnalysisReport | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);

  // Generate analysis reports from inquiries using V0DataAdapter
  useEffect(() => {
    if (inquiries.length > 0) {
      const adaptedReports = inquiries.slice(0, 3).map((inquiry, index) => 
        V0DataAdapter.safeAdaptInquiryToAnalysisReport(inquiry, index)
      ).filter((report): report is AnalysisReport => report !== null);
      setReports(adaptedReports);
    }
  }, [inquiries]);

  const handleGenerateReport = () => {
    setLoading(true);
    // Simulate report generation
    setTimeout(() => {
      if (inquiries.length > 0) {
        const newReport = V0DataAdapter.safeAdaptInquiryToAnalysisReport(inquiries[0], reports.length);
        if (newReport) {
          setReports(prev => [newReport, ...prev]);
        }
      }
      setLoading(false);
      onGenerateReport?.(inquiries[0]?.id);
    }, 2000);
  };

  const handleViewReport = (report: AnalysisReport) => {
    setSelectedReport(report);
    setIsModalOpen(true);
    onViewReport?.(report.id);
  };

  const handleDownloadReport = (reportId: string, format: 'pdf' | 'html' = 'pdf') => {
    // Simulate download
    const report = reports.find(r => r.id === reportId);
    if (report) {
      const blob = new Blob([`Analysis Report: ${report.title}`], { type: 'text/plain' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `${report.title.replace(/[^a-zA-Z0-9]/g, '_')}.${format}`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
    }
    onDownloadReport?.(reportId);
  };

  const closeModal = () => {
    setIsModalOpen(false);
    setSelectedReport(null);
  };

  const toggleExpanded = (reportId: string) => {
    setExpandedReport(expandedReport === reportId ? null : reportId);
  };

  const getRiskBadgeColor = (risk: string) => {
    switch (risk) {
      case 'High': return 'bg-red-100 text-red-800 border-red-200';
      case 'Medium': return 'bg-yellow-100 text-yellow-800 border-yellow-200';
      case 'Low': return 'bg-green-100 text-green-800 border-green-200';
      default: return 'bg-gray-100 text-gray-800 border-gray-200';
    }
  };

  const getConfidenceColor = (confidence: number) => {
    if (confidence >= 80) return 'bg-green-500';
    if (confidence >= 60) return 'bg-yellow-500';
    return 'bg-red-500';
  };

  const getPriorityIcon = (priority: string) => {
    switch (priority) {
      case 'High': return <AlertTriangle className="w-4 h-4 text-red-500" />;
      case 'Medium': return <Clock className="w-4 h-4 text-yellow-500" />;
      case 'Low': return <Target className="w-4 h-4 text-green-500" />;
      default: return <Target className="w-4 h-4 text-gray-500" />;
    }
  };

  return (
    <div className="bg-white rounded-lg border border-gray-200 shadow-sm">
      {/* Header */}
      <div className="px-6 py-4 border-b border-gray-200">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <div className="p-2 bg-blue-100 rounded-lg">
              <FileText className="w-5 h-5 text-blue-600" />
            </div>
            <div>
              <h3 className="text-lg font-semibold text-gray-900">
                AI-Generated Inquiry Analysis Reports
              </h3>
              <p className="text-sm text-gray-600">
                Automated insights and recommendations for high-value opportunities
              </p>
            </div>
          </div>
          <button
            onClick={handleGenerateReport}
            disabled={loading}
            className="inline-flex items-center px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            {loading ? (
              <>
                <div className="animate-spin rounded-full h-4 w-4 border-2 border-white border-t-transparent mr-2"></div>
                Generating...
              </>
            ) : (
              <>
                <Plus className="w-4 h-4 mr-2" />
                Generate New Report
              </>
            )}
          </button>
        </div>
      </div>

      {/* Reports List */}
      <div className="p-6">
        {reports.length === 0 ? (
          <div className="text-center py-8">
            <FileText className="w-12 h-12 text-gray-400 mx-auto mb-4" />
            <h4 className="text-lg font-medium text-gray-900 mb-2">No Analysis Reports Yet</h4>
            <p className="text-gray-600 mb-4">
              Generate your first AI analysis report to get insights on high-value opportunities.
            </p>
            <button
              onClick={handleGenerateReport}
              disabled={loading || inquiries.length === 0}
              className="inline-flex items-center px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              <Plus className="w-4 h-4 mr-2" />
              Generate First Report
            </button>
          </div>
        ) : (
          <div className="space-y-4">
            {reports.map((report) => (
              <div key={report.id} className="border border-gray-200 rounded-lg overflow-hidden">
                {/* Report Card Header */}
                <div className="p-4 bg-gray-50 border-b border-gray-200">
                  <div className="flex items-center justify-between">
                    <div className="flex-1">
                      <div className="flex items-center space-x-3 mb-2">
                        <h4 className="text-lg font-semibold text-gray-900">{report.title}</h4>
                        <span className={`px-2 py-1 text-xs font-medium rounded-full border ${getRiskBadgeColor(report.risk)}`}>
                          {report.risk} Risk
                        </span>
                      </div>
                      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                        <div>
                          <span className="text-gray-600">Customer:</span>
                          <p className="font-medium text-gray-900">{report.customer}</p>
                        </div>
                        <div>
                          <span className="text-gray-600">Service:</span>
                          <p className="font-medium text-gray-900">{report.service}</p>
                        </div>
                        <div>
                          <span className="text-gray-600">Est. Value:</span>
                          <p className="font-medium text-green-600">{report.value}</p>
                        </div>
                        <div>
                          <span className="text-gray-600">Timeline:</span>
                          <p className="font-medium text-gray-900">{report.timeline}</p>
                        </div>
                      </div>
                    </div>
                    <div className="flex items-center space-x-3 ml-4">
                      <button
                        onClick={() => handleViewReport(report)}
                        className="p-2 text-gray-600 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
                        title="View Report"
                      >
                        <Eye className="w-4 h-4" />
                      </button>
                      <button
                        onClick={() => handleDownloadReport(report.id)}
                        className="p-2 text-gray-600 hover:text-green-600 hover:bg-green-50 rounded-lg transition-colors"
                        title="Download Report"
                      >
                        <Download className="w-4 h-4" />
                      </button>
                      <button
                        onClick={() => toggleExpanded(report.id)}
                        className="p-2 text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-colors"
                      >
                        {expandedReport === report.id ? (
                          <ChevronUp className="w-4 h-4" />
                        ) : (
                          <ChevronDown className="w-4 h-4" />
                        )}
                      </button>
                    </div>
                  </div>

                  {/* Confidence Bar */}
                  <div className="mt-4">
                    <div className="flex items-center justify-between text-sm mb-1">
                      <span className="text-gray-600">Confidence Score</span>
                      <span className="font-medium text-gray-900">{report.confidence}%</span>
                    </div>
                    <div className="w-full bg-gray-200 rounded-full h-2">
                      <div
                        className={`h-2 rounded-full transition-all duration-300 ${getConfidenceColor(report.confidence)}`}
                        style={{ width: `${report.confidence}%` }}
                      ></div>
                    </div>
                  </div>
                </div>

                {/* Expanded Content */}
                {expandedReport === report.id && (
                  <div className="p-4 space-y-4">
                    {/* Key Insights */}
                    <div>
                      <h5 className="text-sm font-semibold text-gray-900 mb-2 flex items-center">
                        <TrendingUp className="w-4 h-4 mr-2 text-blue-600" />
                        Key Insights
                      </h5>
                      <ul className="space-y-1">
                        {report.insights.map((insight, index) => (
                          <li key={index} className="text-sm text-gray-700 flex items-start">
                            <span className="w-1.5 h-1.5 bg-blue-600 rounded-full mt-2 mr-2 flex-shrink-0"></span>
                            {insight}
                          </li>
                        ))}
                      </ul>
                    </div>

                    {/* Recommended Actions */}
                    <div>
                      <h5 className="text-sm font-semibold text-gray-900 mb-2 flex items-center">
                        <Target className="w-4 h-4 mr-2 text-green-600" />
                        Recommended Actions
                      </h5>
                      <div className="space-y-2">
                        {report.actions.map((action) => (
                          <div key={action.id} className="flex items-start space-x-3 p-3 bg-gray-50 rounded-lg">
                            <div className="flex-shrink-0 mt-0.5">
                              {getPriorityIcon(action.priority)}
                            </div>
                            <div className="flex-1 min-w-0">
                              <div className="flex items-center justify-between mb-1">
                                <h6 className="text-sm font-medium text-gray-900">{action.title}</h6>
                                <div className="flex items-center space-x-2">
                                  <span className="text-xs text-gray-600">{action.estimatedImpact}</span>
                                  {action.completed ? (
                                    <CheckCircle className="w-4 h-4 text-green-500" />
                                  ) : (
                                    <XCircle className="w-4 h-4 text-gray-400" />
                                  )}
                                </div>
                              </div>
                              <p className="text-sm text-gray-600">{action.description}</p>
                            </div>
                          </div>
                        ))}
                      </div>
                    </div>

                    {/* Report Metadata */}
                    <div className="pt-3 border-t border-gray-200">
                      <p className="text-xs text-gray-500">
                        Generated on {new Date(report.generatedAt).toLocaleDateString()} at{' '}
                        {new Date(report.generatedAt).toLocaleTimeString()}
                      </p>
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Report Modal */}
      <V0ReportModal
        report={selectedReport}
        isOpen={isModalOpen}
        onClose={closeModal}
        onDownload={handleDownloadReport}
      />
    </div>
  );
};



export default V0InquiryAnalysisSection;