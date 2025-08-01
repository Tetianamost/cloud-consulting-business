import React from 'react';
import { AlertTriangle, RefreshCw, Wifi, WifiOff, Database, Palette } from 'lucide-react';

interface ErrorFallbackProps {
  onRetry?: () => void;
  message?: string;
}

/**
 * Generic error fallback with v0 styling
 */
export const V0GenericErrorFallback: React.FC<ErrorFallbackProps> = ({ 
  onRetry, 
  message = "Something went wrong" 
}) => (
  <div className="flex items-center justify-center p-8">
    <div className="text-center max-w-sm">
      <div className="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-red-100 mb-4">
        <AlertTriangle className="h-6 w-6 text-red-600" />
      </div>
      <h3 className="text-lg font-medium text-gray-900 mb-2">{message}</h3>
      <p className="text-sm text-gray-600 mb-4">
        Please try again or contact support if the problem persists.
      </p>
      {onRetry && (
        <button
          onClick={onRetry}
          className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors"
        >
          <RefreshCw className="h-4 w-4 mr-2" />
          Try Again
        </button>
      )}
    </div>
  </div>
);

/**
 * API error fallback for data loading failures
 */
export const V0ApiErrorFallback: React.FC<ErrorFallbackProps> = ({ 
  onRetry, 
  message = "Failed to load data" 
}) => (
  <div className="bg-white rounded-lg border border-red-200 p-6">
    <div className="flex items-center space-x-3">
      <div className="flex-shrink-0">
        <Database className="h-5 w-5 text-red-500" />
      </div>
      <div className="flex-1">
        <h3 className="text-sm font-medium text-red-800">{message}</h3>
        <p className="text-sm text-red-600 mt-1">
          Unable to connect to the server. Please check your connection and try again.
        </p>
      </div>
      {onRetry && (
        <div className="flex-shrink-0">
          <button
            onClick={onRetry}
            className="inline-flex items-center px-3 py-1.5 border border-red-300 text-xs font-medium rounded text-red-700 bg-red-50 hover:bg-red-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 transition-colors"
          >
            <RefreshCw className="h-3 w-3 mr-1" />
            Retry
          </button>
        </div>
      )}
    </div>
  </div>
);

/**
 * Network error fallback for connectivity issues
 */
export const V0NetworkErrorFallback: React.FC<ErrorFallbackProps> = ({ 
  onRetry, 
  message = "Connection lost" 
}) => (
  <div className="bg-white rounded-lg border border-yellow-200 p-6">
    <div className="flex items-center space-x-3">
      <div className="flex-shrink-0">
        <WifiOff className="h-5 w-5 text-yellow-500" />
      </div>
      <div className="flex-1">
        <h3 className="text-sm font-medium text-yellow-800">{message}</h3>
        <p className="text-sm text-yellow-600 mt-1">
          Please check your internet connection and try again.
        </p>
      </div>
      {onRetry && (
        <div className="flex-shrink-0">
          <button
            onClick={onRetry}
            className="inline-flex items-center px-3 py-1.5 border border-yellow-300 text-xs font-medium rounded text-yellow-700 bg-yellow-50 hover:bg-yellow-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-yellow-500 transition-colors"
          >
            <Wifi className="h-3 w-3 mr-1" />
            Reconnect
          </button>
        </div>
      )}
    </div>
  </div>
);

/**
 * Tailwind CSS loading failure fallback
 */
export const V0TailwindErrorFallback: React.FC<ErrorFallbackProps> = ({ 
  onRetry, 
  message = "Styling failed to load" 
}) => (
  <div style={{
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    padding: '2rem',
    backgroundColor: '#f9fafb',
    border: '1px solid #e5e7eb',
    borderRadius: '0.5rem',
    margin: '1rem'
  }}>
    <div style={{ textAlign: 'center', maxWidth: '24rem' }}>
      <div style={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        width: '3rem',
        height: '3rem',
        backgroundColor: '#fef3c7',
        borderRadius: '50%',
        margin: '0 auto 1rem'
      }}>
        <Palette style={{ width: '1.5rem', height: '1.5rem', color: '#d97706' }} />
      </div>
      <h3 style={{ 
        fontSize: '1.125rem', 
        fontWeight: '500', 
        color: '#111827', 
        marginBottom: '0.5rem' 
      }}>
        {message}
      </h3>
      <p style={{ 
        fontSize: '0.875rem', 
        color: '#6b7280', 
        marginBottom: '1rem' 
      }}>
        The page styling couldn't load properly. The content is still functional.
      </p>
      {onRetry && (
        <button
          onClick={onRetry}
          style={{
            display: 'inline-flex',
            alignItems: 'center',
            padding: '0.5rem 1rem',
            backgroundColor: '#3b82f6',
            color: 'white',
            border: 'none',
            borderRadius: '0.375rem',
            fontSize: '0.875rem',
            fontWeight: '500',
            cursor: 'pointer',
            transition: 'background-color 0.2s'
          }}
          onMouseOver={(e) => e.currentTarget.style.backgroundColor = '#2563eb'}
          onMouseOut={(e) => e.currentTarget.style.backgroundColor = '#3b82f6'}
        >
          <RefreshCw style={{ width: '1rem', height: '1rem', marginRight: '0.5rem' }} />
          Reload Page
        </button>
      )}
    </div>
  </div>
);

/**
 * Compact error fallback for smaller components
 */
export const V0CompactErrorFallback: React.FC<ErrorFallbackProps> = ({ 
  onRetry, 
  message = "Error" 
}) => (
  <div className="flex items-center justify-center p-4 bg-red-50 border border-red-200 rounded-md">
    <div className="flex items-center space-x-2 text-red-600">
      <AlertTriangle className="h-4 w-4" />
      <span className="text-sm font-medium">{message}</span>
      {onRetry && (
        <button
          onClick={onRetry}
          className="text-xs underline hover:no-underline ml-2"
        >
          Retry
        </button>
      )}
    </div>
  </div>
);

/**
 * Card-specific error fallback that maintains card layout
 */
export const V0CardErrorFallback: React.FC<ErrorFallbackProps> = ({ 
  onRetry, 
  message = "Failed to load" 
}) => (
  <div className="bg-white rounded-lg border border-red-200 shadow-sm p-6">
    <div className="text-center">
      <div className="mx-auto flex items-center justify-center h-10 w-10 rounded-full bg-red-100 mb-3">
        <AlertTriangle className="h-5 w-5 text-red-600" />
      </div>
      <h3 className="text-sm font-medium text-gray-900 mb-1">{message}</h3>
      <p className="text-xs text-gray-600 mb-3">
        This component couldn't load properly.
      </p>
      {onRetry && (
        <button
          onClick={onRetry}
          className="inline-flex items-center px-3 py-1.5 border border-transparent text-xs font-medium rounded text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 transition-colors"
        >
          <RefreshCw className="h-3 w-3 mr-1" />
          Retry
        </button>
      )}
    </div>
  </div>
);

/**
 * Higher-order component to wrap components with error boundaries
 */
export function withV0ErrorBoundary<P extends object>(
  Component: React.ComponentType<P>,
  fallback?: React.ComponentType<ErrorFallbackProps>
) {
  const WrappedComponent: React.FC<P> = (props) => {
    const [error, setError] = React.useState<Error | null>(null);

    const handleRetry = React.useCallback(() => {
      setError(null);
    }, []);

    if (error) {
      const FallbackComponent = fallback || V0GenericErrorFallback;
      return <FallbackComponent onRetry={handleRetry} message={error.message} />;
    }

    try {
      return <Component {...props} />;
    } catch (err) {
      setError(err instanceof Error ? err : new Error('Unknown error'));
      return null;
    }
  };

  WrappedComponent.displayName = `withV0ErrorBoundary(${Component.displayName || Component.name})`;
  return WrappedComponent;
}