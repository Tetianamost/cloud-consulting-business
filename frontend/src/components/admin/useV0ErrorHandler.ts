import { useState, useCallback, useEffect } from 'react';

export interface V0ErrorState {
  error: Error | null;
  isError: boolean;
  errorType: 'api' | 'network' | 'component' | 'tailwind' | 'unknown';
  retryCount: number;
}

export interface V0ErrorHandlerOptions {
  maxRetries?: number;
  onError?: (error: Error, errorInfo: any) => void;
  resetOnPropsChange?: boolean;
}

/**
 * Custom hook for handling errors in V0 components
 * Provides error state management and retry functionality
 */
export function useV0ErrorHandler(options: V0ErrorHandlerOptions = {}) {
  const { maxRetries = 3, onError, resetOnPropsChange = true } = options;

  const [errorState, setErrorState] = useState<V0ErrorState>({
    error: null,
    isError: false,
    errorType: 'unknown',
    retryCount: 0,
  });

  // Determine error type based on error message/properties
  const getErrorType = useCallback((error: Error): V0ErrorState['errorType'] => {
    const message = error.message.toLowerCase();
    
    if (message.includes('network') || message.includes('fetch')) {
      return 'network';
    }
    if (message.includes('api') || message.includes('server') || message.includes('400') || message.includes('500')) {
      return 'api';
    }
    if (message.includes('tailwind') || message.includes('css') || message.includes('style')) {
      return 'tailwind';
    }
    if (error.name === 'ChunkLoadError' || message.includes('loading chunk')) {
      return 'component';
    }
    
    return 'unknown';
  }, []);

  // Set error state
  const setError = useCallback((error: Error | null, additionalInfo?: any) => {
    if (error) {
      const errorType = getErrorType(error);
      setErrorState(prev => ({
        error,
        isError: true,
        errorType,
        retryCount: prev.retryCount,
      }));

      // Call custom error handler if provided
      if (onError) {
        onError(error, { ...additionalInfo, errorType });
      }

      // Log error for debugging
      console.error(`V0 ${errorType} error:`, error, additionalInfo);
    } else {
      setErrorState({
        error: null,
        isError: false,
        errorType: 'unknown',
        retryCount: 0,
      });
    }
  }, [getErrorType, onError]);

  // Clear error state
  const clearError = useCallback(() => {
    setError(null);
  }, [setError]);

  // Retry with exponential backoff
  const retry = useCallback(async (retryFn?: () => Promise<void> | void) => {
    if (errorState.retryCount >= maxRetries) {
      console.warn('Max retries reached for V0 error handler');
      return false;
    }

    try {
      setErrorState(prev => ({
        ...prev,
        retryCount: prev.retryCount + 1,
      }));

      // Wait with exponential backoff
      const delay = Math.min(1000 * Math.pow(2, errorState.retryCount), 10000);
      await new Promise(resolve => setTimeout(resolve, delay));

      if (retryFn) {
        await retryFn();
      }

      // Clear error if retry was successful
      clearError();
      return true;
    } catch (error) {
      setError(error instanceof Error ? error : new Error('Retry failed'));
      return false;
    }
  }, [errorState.retryCount, maxRetries, clearError, setError]);

  // Async error handler wrapper
  const handleAsync = useCallback(async <T>(
    asyncFn: () => Promise<T>,
    errorMessage?: string
  ): Promise<T | null> => {
    try {
      clearError();
      return await asyncFn();
    } catch (error) {
      const errorToSet = error instanceof Error 
        ? error 
        : new Error(errorMessage || 'Async operation failed');
      setError(errorToSet);
      return null;
    }
  }, [clearError, setError]);

  // Sync error handler wrapper
  const handleSync = useCallback(<T>(
    syncFn: () => T,
    errorMessage?: string
  ): T | null => {
    try {
      clearError();
      return syncFn();
    } catch (error) {
      const errorToSet = error instanceof Error 
        ? error 
        : new Error(errorMessage || 'Sync operation failed');
      setError(errorToSet);
      return null;
    }
  }, [clearError, setError]);

  // Reset error state when component unmounts or props change
  useEffect(() => {
    // Removed cleanup to prevent infinite update loop on unmount
  }, [resetOnPropsChange]);

  return {
    ...errorState,
    setError,
    clearError,
    retry,
    handleAsync,
    handleSync,
    canRetry: errorState.retryCount < maxRetries,
  };
}

/**
 * Hook specifically for API error handling
 */
export function useV0ApiErrorHandler(options: V0ErrorHandlerOptions = {}) {
  const errorHandler = useV0ErrorHandler(options);

  const handleApiCall = useCallback(async <T>(
    apiCall: () => Promise<{ success: boolean; data?: T; error?: string }>,
    errorMessage?: string
  ): Promise<T | null> => {
    return errorHandler.handleAsync(async () => {
      const response = await apiCall();
      
      if (!response.success) {
        throw new Error(response.error || errorMessage || 'API call failed');
      }
      
      return response.data || null;
    }, errorMessage);
  }, [errorHandler]);

  return {
    ...errorHandler,
    handleApiCall,
  };
}

/**
 * Hook for handling Tailwind CSS loading errors
 */
export function useV0TailwindErrorHandler() {
  const [tailwindLoaded, setTailwindLoaded] = useState(true);
  const errorHandler = useV0ErrorHandler({
    maxRetries: 1,
    onError: (error) => {
      if (error.message.includes('tailwind') || error.message.includes('css')) {
        setTailwindLoaded(false);
      }
    }
  });

  // Check if Tailwind is loaded by testing a known class
  useEffect(() => {
    let didSet = false;
    const checkTailwind = () => {
      if (didSet) return;
      didSet = true;
      const testElement = document.createElement('div');
      testElement.className = 'bg-blue-500';
      testElement.style.display = 'none';
      document.body.appendChild(testElement);
      
      const computedStyle = window.getComputedStyle(testElement);
      const isLoaded = computedStyle.backgroundColor === 'rgb(59, 130, 246)'; // bg-blue-500
      
      document.body.removeChild(testElement);
      
      if (!isLoaded) {
        errorHandler.setError(new Error('Tailwind CSS failed to load'));
        setTailwindLoaded(false);
      } else {
        setTailwindLoaded(true);
        errorHandler.clearError();
      }
    };

    // Check only once on mount
    checkTailwind();

    return () => {};
  }, []);

  const retryTailwind = useCallback(() => {
    // Force reload the page to retry loading Tailwind
    window.location.reload();
  }, []);

  return {
    ...errorHandler,
    tailwindLoaded,
    retryTailwind,
  };
}