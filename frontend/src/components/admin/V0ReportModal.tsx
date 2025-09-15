import React from 'react';
import { 
  X, 
  Download, 
  FileText, 
  TrendingUp, 
  Target, 
  AlertTriangle,
  Clock,
  DollarSign,
  User,
  Building,
  Calendar,
  CheckCircle,
  XCircle
} from 'lucide-react';
import { AnalysisReport, RecommendedAction } from './V0InquiryAnalysisSection';

interface V0ReportModalProps {
  report: AnalysisReport | null;
  isOpen: boolean;
  onClose: () => void;
  onDownload?: (reportId: string, format: 'pdf' | 'html') => void;
}

const V0ReportModal: React.FC<V0ReportModalProps> = ({
  report,
  isOpen,
  onClose,
  onDownload
}) => {
  if (!isOpen || !report) return null;

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

  const handleDownload = (format: 'pdf' | 'html') => {
    onDownload?.(report.id, format);
  };

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto">
      {/* Backdrop */}
      <div 
        className="fixed inset-0 bg-black bg-opacity-50 transition-opacity"
        onClick={onClose}
      ></div>
      
      {/* Modal */}
      <div className="flex min-h-full items-center justify-center p-4">
        <div className="relative bg-white rounded-lg shadow-xl max-w-4xl w-full max-h-[90vh] overflow-hidden">
          {/* Header */}
          <div className="flex items-center justify-between p-6 border-b border-gray-200">
            <div className="flex items-center space-x-3">
              <div className="p-2 bg-blue-100 rounded-lg">
                <FileText className="w-5 h-5 text-blue-600" />
              </div>
              <div>
                <h2 className="text-xl font-semibold text-gray-900">{report.title}</h2>
                <p className="text-sm text-gray-600">
                  Generated on {new Date(report.generatedAt).toLocaleDateString()}
                </p>
              </div>
            </div>
            <div className="flex items-center space-x-2">
              <button
                onClick={() => handleDownload('pdf')}
                className="inline-flex items-center px-3 py-2 bg-green-600 text-white text-sm font-medium rounded-lg hover:bg-green-700 transition-colors"
              >
                <Download className="w-4 h-4 mr-2" />
                PDF
              </button>
              <button
                onClick={() => handleDownload('html')}
                className="inline-flex items-center px-3 py-2 bg-blue-600 text-white text-sm font-medium rounded-lg hover:bg-blue-700 transition-colors"
              >
                <Download className="w-4 h-4 mr-2" />
                HTML
              </button>
              <button
                onClick={onClose}
                className="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors"
              >
                <X className="w-5 h-5" />
              </button>
            </div>
          </div>

          {/* Content */}
          <div className="overflow-y-auto max-h-[calc(90vh-80px)]">
            <div className="p-6 space-y-6">
              {/* Overview Cards */}
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                <div className="bg-gray-50 rounded-lg p-4">
                  <div className="flex items-center space-x-2 mb-2">
                    <User className="w-4 h-4 text-gray-600" />
                    <span className="text-sm font-medium text-gray-600">Customer</span>
                  </div>
                  <p className="text-lg font-semibold text-gray-900">{report.customer}</p>
                </div>
                
                <div className="bg-gray-50 rounded-lg p-4">
                  <div className="flex items-center space-x-2 mb-2">
                    <Building className="w-4 h-4 text-gray-600" />
                    <span className="text-sm font-medium text-gray-600">Service</span>
                  </div>
                  <p className="text-lg font-semibold text-gray-900">{report.service}</p>
                </div>
                
                <div className="bg-gray-50 rounded-lg p-4">
                  <div className="flex items-center space-x-2 mb-2">
                    <DollarSign className="w-4 h-4 text-green-600" />
                    <span className="text-sm font-medium text-gray-600">Est. Value</span>
                  </div>
                  <p className="text-lg font-semibold text-green-600">{report.value}</p>
                </div>
                
                <div className="bg-gray-50 rounded-lg p-4">
                  <div className="flex items-center space-x-2 mb-2">
                    <Calendar className="w-4 h-4 text-gray-600" />
                    <span className="text-sm font-medium text-gray-600">Timeline</span>
                  </div>
                  <p className="text-lg font-semibold text-gray-900">{report.timeline}</p>
                </div>
              </div>

              {/* Risk and Confidence */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {/* Risk Assessment */}
                <div className="bg-white border border-gray-200 rounded-lg p-4">
                  <h3 className="text-lg font-semibold text-gray-900 mb-3 flex items-center">
                    <AlertTriangle className="w-5 h-5 mr-2 text-orange-600" />
                    Risk Assessment
                  </h3>
                  <div className="flex items-center justify-center">
                    <span className={`px-4 py-2 text-lg font-medium rounded-full border ${getRiskBadgeColor(report.risk)}`}>
                      {report.risk} Risk
                    </span>
                  </div>
                </div>

                {/* Confidence Score */}
                <div className="bg-white border border-gray-200 rounded-lg p-4">
                  <h3 className="text-lg font-semibold text-gray-900 mb-3 flex items-center">
                    <Target className="w-5 h-5 mr-2 text-blue-600" />
                    Confidence Score
                  </h3>
                  <div className="space-y-2">
                    <div className="flex items-center justify-between">
                      <span className="text-2xl font-bold text-gray-900">{report.confidence}%</span>
                      <span className="text-sm text-gray-600">
                        {report.confidence >= 80 ? 'High Confidence' : 
                         report.confidence >= 60 ? 'Medium Confidence' : 'Low Confidence'}
                      </span>
                    </div>
                    <div className="w-full bg-gray-200 rounded-full h-3">
                      <div
                        className={`h-3 rounded-full transition-all duration-300 ${getConfidenceColor(report.confidence)}`}
                        style={{ width: `${report.confidence}%` }}
                      ></div>
                    </div>
                  </div>
                </div>
              </div>

              {/* Key Insights */}
              <div className="bg-white border border-gray-200 rounded-lg p-6">
                <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                  <TrendingUp className="w-5 h-5 mr-2 text-blue-600" />
                  Key Insights
                </h3>
                <div className="space-y-3">
                  {report.insights.map((insight, index) => (
                    <div key={index} className="flex items-start space-x-3 p-3 bg-blue-50 rounded-lg">
                      <div className="w-2 h-2 bg-blue-600 rounded-full mt-2 flex-shrink-0"></div>
                      <p className="text-gray-700">{insight}</p>
                    </div>
                  ))}
                </div>
              </div>

              {/* Recommended Actions */}
              <div className="bg-white border border-gray-200 rounded-lg p-6">
                <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                  <Target className="w-5 h-5 mr-2 text-green-600" />
                  Recommended Actions
                </h3>
                <div className="space-y-4">
                  {report.actions.map((action) => (
                    <div key={action.id} className="border border-gray-200 rounded-lg p-4">
                      <div className="flex items-start justify-between mb-3">
                        <div className="flex items-start space-x-3">
                          <div className="flex-shrink-0 mt-0.5">
                            {getPriorityIcon(action.priority)}
                          </div>
                          <div>
                            <h4 className="text-lg font-medium text-gray-900">{action.title}</h4>
                            <span className={`inline-block px-2 py-1 text-xs font-medium rounded-full mt-1 ${
                              action.priority === 'High' ? 'bg-red-100 text-red-800' :
                              action.priority === 'Medium' ? 'bg-yellow-100 text-yellow-800' :
                              'bg-green-100 text-green-800'
                            }`}>
                              {action.priority} Priority
                            </span>
                          </div>
                        </div>
                        <div className="flex items-center space-x-3">
                          <span className="text-sm font-medium text-blue-600">{action.estimatedImpact}</span>
                          {action.completed ? (
                            <CheckCircle className="w-5 h-5 text-green-500" />
                          ) : (
                            <XCircle className="w-5 h-5 text-gray-400" />
                          )}
                        </div>
                      </div>
                      <p className="text-gray-600 ml-7">{action.description}</p>
                    </div>
                  ))}
                </div>
              </div>

              {/* Report Metadata */}
              <div className="bg-gray-50 rounded-lg p-4">
                <h3 className="text-sm font-semibold text-gray-900 mb-2">Report Information</h3>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
                  <div>
                    <span className="text-gray-600">Report ID:</span>
                    <p className="font-mono text-gray-900">{report.id}</p>
                  </div>
                  <div>
                    <span className="text-gray-600">Inquiry ID:</span>
                    <p className="font-mono text-gray-900">{report.inquiryId}</p>
                  </div>
                  <div>
                    <span className="text-gray-600">Generated:</span>
                    <p className="text-gray-900">
                      {new Date(report.generatedAt).toLocaleString()}
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default V0ReportModal;