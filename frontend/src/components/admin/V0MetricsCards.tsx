import React from 'react';
import { TrendingUp, TrendingDown, Minus, FileText, Target, Clock, AlertTriangle } from 'lucide-react';

// Types for metric card data
export interface MetricCardData {
  title: string;
  value: string | number;
  change: string;
  trend: 'up' | 'down' | 'neutral';
  icon: React.ComponentType<{ className?: string }>;
}

interface V0MetricsCardsProps {
  metrics: MetricCardData[];
  loading?: boolean;
}

/**
 * V0MetricsCards - Metric cards component matching v0.dev design
 * Features proper shadows, spacing, typography, and trend indicators
 */
const V0MetricsCards: React.FC<V0MetricsCardsProps> = ({ metrics, loading = false }) => {
  if (loading) {
    return <V0MetricsCardsSkeleton />;
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      {metrics.map((metric, index) => (
        <V0MetricCard key={index} {...metric} />
      ))}
    </div>
  );
};

/**
 * Individual metric card component
 */
interface V0MetricCardProps extends MetricCardData {}

const V0MetricCard: React.FC<V0MetricCardProps> = ({ 
  title, 
  value, 
  change, 
  trend, 
  icon: Icon 
}) => {
  const getTrendIcon = () => {
    switch (trend) {
      case 'up':
        return <TrendingUp className="h-4 w-4 text-green-500" />;
      case 'down':
        return <TrendingDown className="h-4 w-4 text-red-500" />;
      default:
        return <Minus className="h-4 w-4 text-gray-400" />;
    }
  };

  const getTrendColor = () => {
    switch (trend) {
      case 'up':
        return 'text-green-600';
      case 'down':
        return 'text-red-600';
      default:
        return 'text-gray-600';
    }
  };

  return (
    <div className="bg-white rounded-lg border border-gray-200 p-6 shadow-sm hover:shadow-md transition-shadow duration-200">
      {/* Header with title and icon */}
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-sm font-medium text-gray-600 truncate">{title}</h3>
        <div className="flex-shrink-0">
          <Icon className="h-5 w-5 text-gray-400" />
        </div>
      </div>

      {/* Main value */}
      <div className="mb-2">
        <div className="text-2xl font-bold text-gray-900">{value}</div>
      </div>

      {/* Trend indicator */}
      <div className={`flex items-center text-sm ${getTrendColor()}`}>
        {getTrendIcon()}
        <span className="ml-1">{change}</span>
      </div>
    </div>
  );
};

/**
 * Loading skeleton for metrics cards
 */
const V0MetricsCardsSkeleton: React.FC = () => {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      {[...Array(4)].map((_, index) => (
        <div key={index} className="bg-white rounded-lg border border-gray-200 p-6 shadow-sm animate-pulse">
          {/* Header skeleton */}
          <div className="flex items-center justify-between mb-4">
            <div className="h-4 bg-gray-200 rounded w-24"></div>
            <div className="h-5 w-5 bg-gray-200 rounded"></div>
          </div>

          {/* Value skeleton */}
          <div className="mb-2">
            <div className="h-8 bg-gray-200 rounded w-16"></div>
          </div>

          {/* Trend skeleton */}
          <div className="flex items-center">
            <div className="h-4 w-4 bg-gray-200 rounded mr-1"></div>
            <div className="h-4 bg-gray-200 rounded w-20"></div>
          </div>
        </div>
      ))}
    </div>
  );
};

/**
 * Default metric card configurations for common metrics
 */
export const defaultMetricCards = {
  aiReports: {
    title: "AI Reports Generated",
    icon: FileText,
  },
  confidence: {
    title: "Avg Confidence Score", 
    icon: Target,
  },
  processingTime: {
    title: "Avg Processing Time",
    icon: Clock,
  },
  opportunities: {
    title: "High-Value Opportunities",
    icon: AlertTriangle,
  },
};

export default V0MetricsCards;