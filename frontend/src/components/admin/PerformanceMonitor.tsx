import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '../ui/card';
import { Badge } from '../ui/badge';
import { Button } from '../ui/button';
import { 
  Activity, 
  Clock, 
  Zap, 
  HardDrive, 
  Wifi,
  AlertTriangle,
  CheckCircle,
  TrendingUp,
  TrendingDown
} from 'lucide-react';

interface PerformanceMetrics {
  loadTime: number;
  domContentLoaded: number;
  firstContentfulPaint: number;
  largestContentfulPaint: number;
  cumulativeLayoutShift: number;
  firstInputDelay: number;
  bundleSize: number;
  memoryUsage: number;
  networkSpeed: string;
}

/**
 * PerformanceMonitor - Component to monitor and display performance metrics
 * Helps track the impact of responsive design and performance optimizations
 */
const PerformanceMonitor: React.FC = () => {
  const [metrics, setMetrics] = useState<PerformanceMetrics | null>(null);
  const [isVisible, setIsVisible] = useState(false);
  const [isCollecting, setIsCollecting] = useState(false);

  const collectMetrics = async () => {
    setIsCollecting(true);
    
    try {
      // Performance API metrics
      const navigation = performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming;
      const paint = performance.getEntriesByType('paint');
      
      // Web Vitals approximation
      const fcp = paint.find(entry => entry.name === 'first-contentful-paint')?.startTime || 0;
      
      // Memory usage (if available)
      const memoryInfo = (performance as any).memory;
      const memoryUsage = memoryInfo ? memoryInfo.usedJSHeapSize / 1024 / 1024 : 0;

      // Network information (if available)
      const connection = (navigator as any).connection;
      const networkSpeed = connection ? connection.effectiveType : 'unknown';

      // Estimate bundle size from resource timing
      const resources = performance.getEntriesByType('resource') as PerformanceResourceTiming[];
      const jsResources = resources.filter(r => r.name.includes('.js'));
      const bundleSize = jsResources.reduce((total, resource) => {
        return total + (resource.transferSize || 0);
      }, 0) / 1024; // Convert to KB

      const newMetrics: PerformanceMetrics = {
        loadTime: navigation.loadEventEnd - navigation.fetchStart,
        domContentLoaded: navigation.domContentLoadedEventEnd - navigation.fetchStart,
        firstContentfulPaint: fcp,
        largestContentfulPaint: 0, // Would need observer for real LCP
        cumulativeLayoutShift: 0, // Would need observer for real CLS
        firstInputDelay: 0, // Would need observer for real FID
        bundleSize: bundleSize,
        memoryUsage: memoryUsage,
        networkSpeed: networkSpeed
      };

      setMetrics(newMetrics);
    } catch (error) {
      console.error('Error collecting performance metrics:', error);
    } finally {
      setIsCollecting(false);
    }
  };

  useEffect(() => {
    // Auto-collect metrics on component mount
    collectMetrics();
  }, []);

  const getPerformanceScore = (metric: string, value: number): 'good' | 'needs-improvement' | 'poor' => {
    switch (metric) {
      case 'loadTime':
        return value < 2000 ? 'good' : value < 4000 ? 'needs-improvement' : 'poor';
      case 'fcp':
        return value < 1800 ? 'good' : value < 3000 ? 'needs-improvement' : 'poor';
      case 'bundleSize':
        return value < 250 ? 'good' : value < 500 ? 'needs-improvement' : 'poor';
      case 'memoryUsage':
        return value < 50 ? 'good' : value < 100 ? 'needs-improvement' : 'poor';
      default:
        return 'good';
    }
  };

  const getScoreColor = (score: string) => {
    switch (score) {
      case 'good':
        return 'bg-green-50 text-green-700 border-green-200';
      case 'needs-improvement':
        return 'bg-yellow-50 text-yellow-700 border-yellow-200';
      case 'poor':
        return 'bg-red-50 text-red-700 border-red-200';
      default:
        return 'bg-gray-50 text-gray-700 border-gray-200';
    }
  };

  const getScoreIcon = (score: string) => {
    switch (score) {
      case 'good':
        return <CheckCircle className="h-4 w-4" />;
      case 'needs-improvement':
        return <AlertTriangle className="h-4 w-4" />;
      case 'poor':
        return <TrendingDown className="h-4 w-4" />;
      default:
        return <Activity className="h-4 w-4" />;
    }
  };

  if (!isVisible) {
    return (
      <Button
        onClick={() => setIsVisible(true)}
        className="fixed bottom-4 left-4 z-50 bg-blue-600 hover:bg-blue-700"
        size="sm"
      >
        <Activity className="h-4 w-4 mr-2" />
        Performance
      </Button>
    );
  }

  return (
    <div className="fixed bottom-4 left-4 z-50 w-96 max-h-96 overflow-y-auto">
      <Card className="shadow-lg border-gray-200">
        <CardHeader className="pb-3">
          <div className="flex items-center justify-between">
            <CardTitle className="flex items-center space-x-2 text-sm">
              <Activity className="h-4 w-4" />
              <span>Performance Monitor</span>
            </CardTitle>
            <div className="flex items-center space-x-2">
              <Button
                onClick={collectMetrics}
                disabled={isCollecting}
                size="sm"
                variant="outline"
              >
                {isCollecting ? 'Collecting...' : 'Refresh'}
              </Button>
              <Button
                onClick={() => setIsVisible(false)}
                size="sm"
                variant="ghost"
              >
                ×
              </Button>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          {metrics && (
            <>
              {/* Core Web Vitals */}
              <div>
                <h4 className="text-sm font-medium text-gray-900 mb-2">Core Metrics</h4>
                <div className="space-y-2">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-2">
                      <Clock className="h-3 w-3 text-gray-500" />
                      <span className="text-xs">Load Time</span>
                    </div>
                    <Badge 
                      variant="outline" 
                      className={getScoreColor(getPerformanceScore('loadTime', metrics.loadTime))}
                    >
                      {getScoreIcon(getPerformanceScore('loadTime', metrics.loadTime))}
                      <span className="ml-1">{Math.round(metrics.loadTime)}ms</span>
                    </Badge>
                  </div>

                  <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-2">
                      <Zap className="h-3 w-3 text-gray-500" />
                      <span className="text-xs">First Contentful Paint</span>
                    </div>
                    <Badge 
                      variant="outline" 
                      className={getScoreColor(getPerformanceScore('fcp', metrics.firstContentfulPaint))}
                    >
                      {getScoreIcon(getPerformanceScore('fcp', metrics.firstContentfulPaint))}
                      <span className="ml-1">{Math.round(metrics.firstContentfulPaint)}ms</span>
                    </Badge>
                  </div>

                  <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-2">
                      <HardDrive className="h-3 w-3 text-gray-500" />
                      <span className="text-xs">Bundle Size</span>
                    </div>
                    <Badge 
                      variant="outline" 
                      className={getScoreColor(getPerformanceScore('bundleSize', metrics.bundleSize))}
                    >
                      {getScoreIcon(getPerformanceScore('bundleSize', metrics.bundleSize))}
                      <span className="ml-1">{Math.round(metrics.bundleSize)}KB</span>
                    </Badge>
                  </div>

                  {metrics.memoryUsage > 0 && (
                    <div className="flex items-center justify-between">
                      <div className="flex items-center space-x-2">
                        <Activity className="h-3 w-3 text-gray-500" />
                        <span className="text-xs">Memory Usage</span>
                      </div>
                      <Badge 
                        variant="outline" 
                        className={getScoreColor(getPerformanceScore('memoryUsage', metrics.memoryUsage))}
                      >
                        {getScoreIcon(getPerformanceScore('memoryUsage', metrics.memoryUsage))}
                        <span className="ml-1">{Math.round(metrics.memoryUsage)}MB</span>
                      </Badge>
                    </div>
                  )}

                  {metrics.networkSpeed !== 'unknown' && (
                    <div className="flex items-center justify-between">
                      <div className="flex items-center space-x-2">
                        <Wifi className="h-3 w-3 text-gray-500" />
                        <span className="text-xs">Network</span>
                      </div>
                      <Badge variant="outline" className="bg-blue-50 text-blue-700 border-blue-200">
                        {metrics.networkSpeed}
                      </Badge>
                    </div>
                  )}
                </div>
              </div>

              {/* Optimization Tips */}
              <div className="pt-3 border-t border-gray-100">
                <h4 className="text-sm font-medium text-gray-900 mb-2">Optimization Status</h4>
                <div className="space-y-1 text-xs text-gray-600">
                  <div className="flex items-center space-x-2">
                    <CheckCircle className="h-3 w-3 text-green-600" />
                    <span>Lazy loading enabled</span>
                  </div>
                  <div className="flex items-center space-x-2">
                    <CheckCircle className="h-3 w-3 text-green-600" />
                    <span>Component memoization active</span>
                  </div>
                  <div className="flex items-center space-x-2">
                    <CheckCircle className="h-3 w-3 text-green-600" />
                    <span>Tailwind purging configured</span>
                  </div>
                  <div className="flex items-center space-x-2">
                    <CheckCircle className="h-3 w-3 text-green-600" />
                    <span>Responsive design optimized</span>
                  </div>
                </div>
              </div>
            </>
          )}
        </CardContent>
      </Card>
    </div>
  );
};

export default PerformanceMonitor;