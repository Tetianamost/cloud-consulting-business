import React from 'react';
import { TrendingUp, TrendingDown, Minus, FileText, Target, Clock, AlertTriangle } from 'lucide-react';
import V0ErrorBoundary from './V0ErrorBoundary';
import { V0CardErrorFallback } from './V0ErrorFallbacks';
import { useV0ErrorHandler } from './useV0ErrorHandler';

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
 * Memoized for performance optimization with error handling
 */
function V0MetricsCardsInner({ metrics, loading = false }: V0MetricsCardsProps) {
  const { handleSync, isError } = useV0ErrorHandler();

  if (loading) {
    return <V0MetricsCardsSkeleton />;
  }

  if (isError) {
    return <V0CardErrorFallback message="Failed to load metrics" />;
  }

  let content;
  if (!metrics || metrics.length === 0) {
    content = (
      <div className="col-span-full flex items-center justify-center p-8 text-gray-500">
        No metrics data available
      </div>
    );
  } else {
    content = metrics.map((metric, index) => (
      <V0ErrorBoundary
        key={index}
        fallback={<V0CardErrorFallback message="Metric card error" />}
      >
        <V0MetricCard {...metric} />
      </V0ErrorBoundary>
    ));
  }

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 sm:gap-6">
      {content}
    </div>
  );
}
const V0MetricsCards = React.memo(V0MetricsCardsInner);

/**
 * Individual metric card component
 */
interface V0MetricCardProps extends MetricCardData {}

function V0MetricCardInner({
  title,
  value,
  change,
  trend,
  icon: Icon
}: V0MetricCardProps) {
  const { handleSync } = useV0ErrorHandler();

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

  // Validate required props
  if (!title || value === undefined || value === null) {
    return <V0CardErrorFallback message="Missing required metric data" />;
  }

  return (
    <div className="bg-white rounded-lg border border-gray-200 p-4 sm:p-6 shadow-sm hover:shadow-md transition-shadow duration-200">
      {/* Header with title and icon */}
      <div className="flex items-center justify-between mb-3 sm:mb-4">
        <h3 className="text-sm font-medium text-gray-600 truncate pr-2">{title}</h3>
        <div className="flex-shrink-0">
          {Icon && <Icon className="h-4 w-4 sm:h-5 sm:w-5 text-gray-400" />}
        </div>
      </div>

      {/* Main value */}
      <div className="mb-2">
        <div className="text-xl sm:text-2xl font-bold text-gray-900">{value}</div>
      </div>

      {/* Trend indicator */}
      {change && (
        <div className={`flex items-center text-xs sm:text-sm ${getTrendColor()}`}>
          {getTrendIcon()}
          <span className="ml-1 truncate">{change}</span>
        </div>
      )}
    </div>
  );
}
const V0MetricCard = React.memo(V0MetricCardInner);

/**
 * Loading skeleton for metrics cards
 */
const V0MetricsCardsSkeleton: React.FC = () => {
  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 sm:gap-6">
      {[...Array(4)].map((_, index) => (
        <div key={index} className="bg-white rounded-lg border border-gray-200 p-4 sm:p-6 shadow-sm animate-pulse">
          {/* Header skeleton */}
          <div className="flex items-center justify-between mb-3 sm:mb-4">
            <div className="h-4 bg-gray-200 rounded w-20 sm:w-24"></div>
            <div className="h-4 w-4 sm:h-5 sm:w-5 bg-gray-200 rounded"></div>
          </div>

          {/* Value skeleton */}
          <div className="mb-2">
            <div className="h-6 sm:h-8 bg-gray-200 rounded w-12 sm:w-16"></div>
          </div>

          {/* Trend skeleton */}
          <div className="flex items-center">
            <div className="h-3 w-3 sm:h-4 sm:w-4 bg-gray-200 rounded mr-1"></div>
            <div className="h-3 sm:h-4 bg-gray-200 rounded w-16 sm:w-20"></div>
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