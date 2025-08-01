import React from 'react';
import { Loader2, RefreshCw } from 'lucide-react';

interface LoadingProps {
  size?: 'sm' | 'md' | 'lg';
  message?: string;
  className?: string;
}

/**
 * V0LoadingSpinner - Consistent loading spinner with v0 styling
 */
export const V0LoadingSpinner: React.FC<LoadingProps> = ({ 
  size = 'md', 
  message,
  className = '' 
}) => {
  const sizeClasses = {
    sm: 'h-4 w-4',
    md: 'h-6 w-6',
    lg: 'h-8 w-8'
  };

  return (
    <div className={`flex items-center justify-center ${className}`}>
      <div className="flex flex-col items-center space-y-2">
        <Loader2 className={`${sizeClasses[size]} animate-spin text-blue-600`} />
        {message && (
          <p className="text-sm text-gray-600 animate-pulse">{message}</p>
        )}
      </div>
    </div>
  );
};

/**
 * V0InlineLoader - Small inline loading indicator
 */
export const V0InlineLoader: React.FC<{ message?: string }> = ({ message = 'Loading...' }) => (
  <div className="flex items-center space-x-2 text-sm text-gray-600">
    <Loader2 className="h-4 w-4 animate-spin" />
    <span>{message}</span>
  </div>
);

/**
 * V0ButtonLoader - Loading state for buttons
 */
export const V0ButtonLoader: React.FC<{ message?: string }> = ({ message = 'Loading...' }) => (
  <div className="flex items-center space-x-2">
    <Loader2 className="h-4 w-4 animate-spin" />
    <span>{message}</span>
  </div>
);

/**
 * V0SkeletonCard - Skeleton loader for metric cards
 */
export const V0SkeletonCard: React.FC = () => (
  <div className="bg-white rounded-lg border border-gray-200 p-4 sm:p-6 shadow-sm animate-pulse">
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
);

/**
 * V0SkeletonTable - Skeleton loader for tables
 */
export const V0SkeletonTable: React.FC<{ rows?: number; columns?: number }> = ({ 
  rows = 5, 
  columns = 4 
}) => (
  <div className="bg-white rounded-lg border border-gray-200 overflow-hidden">
    {/* Header skeleton */}
    <div className="border-b border-gray-200 p-4">
      <div className="grid gap-4" style={{ gridTemplateColumns: `repeat(${columns}, 1fr)` }}>
        {Array.from({ length: columns }).map((_, index) => (
          <div key={index} className="h-4 bg-gray-200 rounded animate-pulse"></div>
        ))}
      </div>
    </div>

    {/* Rows skeleton */}
    {Array.from({ length: rows }).map((_, rowIndex) => (
      <div key={rowIndex} className="border-b border-gray-100 p-4 last:border-b-0">
        <div className="grid gap-4" style={{ gridTemplateColumns: `repeat(${columns}, 1fr)` }}>
          {Array.from({ length: columns }).map((_, colIndex) => (
            <div 
              key={colIndex} 
              className="h-4 bg-gray-100 rounded animate-pulse"
              style={{ animationDelay: `${(rowIndex * columns + colIndex) * 0.1}s` }}
            ></div>
          ))}
        </div>
      </div>
    ))}
  </div>
);

/**
 * V0SkeletonList - Skeleton loader for lists
 */
export const V0SkeletonList: React.FC<{ items?: number }> = ({ items = 3 }) => (
  <div className="space-y-4">
    {Array.from({ length: items }).map((_, index) => (
      <div key={index} className="bg-white rounded-lg border border-gray-200 p-4 animate-pulse">
        <div className="flex items-center space-x-4">
          <div className="h-10 w-10 bg-gray-200 rounded-full"></div>
          <div className="flex-1 space-y-2">
            <div className="h-4 bg-gray-200 rounded w-3/4"></div>
            <div className="h-3 bg-gray-100 rounded w-1/2"></div>
          </div>
          <div className="h-8 w-16 bg-gray-200 rounded"></div>
        </div>
      </div>
    ))}
  </div>
);

/**
 * V0SkeletonChart - Skeleton loader for charts and graphs
 */
export const V0SkeletonChart: React.FC<{ height?: string }> = ({ height = 'h-64' }) => (
  <div className={`bg-white rounded-lg border border-gray-200 p-6 ${height} animate-pulse`}>
    <div className="h-4 bg-gray-200 rounded w-1/4 mb-4"></div>
    <div className="flex items-end space-x-2 h-full">
      {Array.from({ length: 7 }).map((_, index) => (
        <div 
          key={index}
          className="bg-gray-200 rounded-t flex-1"
          style={{ 
            height: `${Math.random() * 60 + 20}%`,
            animationDelay: `${index * 0.1}s`
          }}
        ></div>
      ))}
    </div>
  </div>
);

/**
 * V0SkeletonDashboard - Complete dashboard skeleton
 */
export const V0SkeletonDashboard: React.FC = () => (
  <div className="space-y-6">
    {/* Header skeleton */}
    <div className="animate-pulse">
      <div className="h-8 bg-gray-200 rounded w-1/3 mb-2"></div>
      <div className="h-4 bg-gray-100 rounded w-1/2"></div>
    </div>

    {/* Metrics cards skeleton */}
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 sm:gap-6">
      {Array.from({ length: 4 }).map((_, index) => (
        <V0SkeletonCard key={index} />
      ))}
    </div>

    {/* Content sections skeleton */}
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <V0SkeletonChart />
      <V0SkeletonList items={4} />
    </div>
  </div>
);

/**
 * V0ProgressBar - Progress indicator with v0 styling
 */
interface ProgressProps {
  value: number;
  max?: number;
  size?: 'sm' | 'md' | 'lg';
  color?: 'blue' | 'green' | 'yellow' | 'red';
  showLabel?: boolean;
  label?: string;
}

export const V0ProgressBar: React.FC<ProgressProps> = ({
  value,
  max = 100,
  size = 'md',
  color = 'blue',
  showLabel = false,
  label
}) => {
  const percentage = Math.min((value / max) * 100, 100);
  
  const sizeClasses = {
    sm: 'h-1',
    md: 'h-2',
    lg: 'h-3'
  };

  const colorClasses = {
    blue: 'bg-blue-600',
    green: 'bg-green-600',
    yellow: 'bg-yellow-600',
    red: 'bg-red-600'
  };

  return (
    <div className="w-full">
      {(showLabel || label) && (
        <div className="flex justify-between items-center mb-1">
          <span className="text-sm font-medium text-gray-700">
            {label || 'Progress'}
          </span>
          {showLabel && (
            <span className="text-sm text-gray-600">
              {Math.round(percentage)}%
            </span>
          )}
        </div>
      )}
      <div className={`w-full bg-gray-200 rounded-full ${sizeClasses[size]}`}>
        <div
          className={`${sizeClasses[size]} ${colorClasses[color]} rounded-full transition-all duration-300 ease-in-out`}
          style={{ width: `${percentage}%` }}
        ></div>
      </div>
    </div>
  );
};

/**
 * V0LoadingOverlay - Full-screen loading overlay
 */
export const V0LoadingOverlay: React.FC<LoadingProps> = ({ 
  message = 'Loading...', 
  size = 'lg' 
}) => (
  <div className="fixed inset-0 bg-white bg-opacity-90 flex items-center justify-center z-50">
    <div className="text-center">
      <V0LoadingSpinner size={size} message={message} />
    </div>
  </div>
);

/**
 * V0RefreshButton - Button with loading state for refresh actions
 */
interface RefreshButtonProps {
  onClick: () => void;
  loading?: boolean;
  disabled?: boolean;
  children?: React.ReactNode;
  size?: 'sm' | 'md' | 'lg';
}

export const V0RefreshButton: React.FC<RefreshButtonProps> = ({
  onClick,
  loading = false,
  disabled = false,
  children = 'Refresh',
  size = 'md'
}) => {
  const sizeClasses = {
    sm: 'px-2 py-1 text-xs',
    md: 'px-3 py-1.5 text-sm',
    lg: 'px-4 py-2 text-base'
  };

  return (
    <button
      onClick={onClick}
      disabled={disabled || loading}
      className={`
        inline-flex items-center space-x-2 border border-gray-300 rounded-md
        bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500
        disabled:opacity-50 disabled:cursor-not-allowed transition-colors
        ${sizeClasses[size]}
      `}
    >
      <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
      <span>{children}</span>
    </button>
  );
};

/**
 * Higher-order component to add loading states to any component
 */
export function withV0Loading<P extends object>(
  Component: React.ComponentType<P>,
  LoadingSkeleton: React.ComponentType = V0LoadingSpinner
) {
  const WrappedComponent: React.FC<P & { loading?: boolean }> = ({ loading, ...props }) => {
    if (loading) {
      return <LoadingSkeleton />;
    }
    
    return <Component {...(props as P)} />;
  };

  WrappedComponent.displayName = `withV0Loading(${Component.displayName || Component.name})`;
  return WrappedComponent;
}