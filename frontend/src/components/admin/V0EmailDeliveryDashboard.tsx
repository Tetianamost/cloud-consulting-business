import React, { useState, useEffect } from 'react';
import { 
  Mail, 
  MailCheck, 
  MousePointer, 
  AlertTriangle,
  Clock,
  TrendingUp,
  TrendingDown,
  Minus
} from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../ui/card';
import { Progress } from '../ui/progress';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '../ui/select';
import { Badge } from '../ui/badge';
import { V0SkeletonCard, V0ProgressBar, V0LoadingSpinner } from './V0LoadingStates';
import { V0ApiErrorFallback } from './V0ErrorFallbacks';
import { useV0ApiErrorHandler } from './useV0ErrorHandler';
import apiService, { SystemMetrics, EmailStatus } from '../../services/api';

// Email metrics interface for v0 component format
export interface EmailMetrics {
  deliveryRate: number;
  openRate: number;
  clickRate: number;
  failedEmails: number;
  totalEmails: number;
  bounced: number;
  spam: number;
  delivered: number;
  opened: number;
  clicked: number;
}

// Email delivery card data interface
interface EmailDeliveryCardData {
  title: string;
  value: string | number;
  percentage: number;
  change: string;
  trend: 'up' | 'down' | 'neutral';
  icon: React.ComponentType<{ className?: string }>;
  color: string;
}

// Props interface for the dashboard component
interface V0EmailDeliveryDashboardProps {
  metrics?: EmailMetrics;
  timeRange?: string;
  onTimeRangeChange?: (range: string) => void;
  className?: string;
}

// Default email metrics for fallback
const defaultEmailMetrics: EmailMetrics = {
  deliveryRate: 0,
  openRate: 0,
  clickRate: 0,
  failedEmails: 0,
  totalEmails: 0,
  bounced: 0,
  spam: 0,
  delivered: 0,
  opened: 0,
  clicked: 0,
};

/**
 * V0EmailDeliveryDashboard - Email delivery monitoring dashboard with v0.dev styling
 * Recreates the email delivery metrics section from v0.dev with proper Tailwind classes
 */
export const V0EmailDeliveryDashboard: React.FC<V0EmailDeliveryDashboardProps> = ({
  metrics = defaultEmailMetrics,
  timeRange = '24h',
  onTimeRangeChange,
  className = '',
}) => {
  const [loading, setLoading] = useState(false);
  const { isError, error, clearError } = useV0ApiErrorHandler();

  // Transform metrics to card data format
  const getEmailDeliveryCards = (emailMetrics: EmailMetrics): EmailDeliveryCardData[] => {
    return [
      {
        title: 'Delivery Rate',
        value: `${emailMetrics.deliveryRate.toFixed(1)}%`,
        percentage: emailMetrics.deliveryRate,
        change: emailMetrics.deliveryRate >= 95 ? '+2.1% from last week' : 'Needs improvement',
        trend: emailMetrics.deliveryRate >= 95 ? 'up' : emailMetrics.deliveryRate >= 85 ? 'neutral' : 'down',
        icon: MailCheck,
        color: 'text-green-600',
      },
      {
        title: 'Open Rate',
        value: `${emailMetrics.openRate.toFixed(1)}%`,
        percentage: emailMetrics.openRate,
        change: emailMetrics.openRate >= 25 ? '+1.8% from last week' : 'Below average',
        trend: emailMetrics.openRate >= 25 ? 'up' : emailMetrics.openRate >= 15 ? 'neutral' : 'down',
        icon: Mail,
        color: 'text-blue-600',
      },
      {
        title: 'Click Rate',
        value: `${emailMetrics.clickRate.toFixed(1)}%`,
        percentage: emailMetrics.clickRate,
        change: emailMetrics.clickRate >= 5 ? '+0.9% from last week' : 'Room for improvement',
        trend: emailMetrics.clickRate >= 5 ? 'up' : emailMetrics.clickRate >= 2 ? 'neutral' : 'down',
        icon: MousePointer,
        color: 'text-purple-600',
      },
      {
        title: 'Failed Emails',
        value: emailMetrics.failedEmails,
        percentage: emailMetrics.totalEmails > 0 ? (emailMetrics.failedEmails / emailMetrics.totalEmails) * 100 : 0,
        change: emailMetrics.failedEmails <= 5 ? 'Within normal range' : 'Requires attention',
        trend: emailMetrics.failedEmails <= 5 ? 'up' : emailMetrics.failedEmails <= 15 ? 'neutral' : 'down',
        icon: AlertTriangle,
        color: 'text-red-600',
      },
    ];
  };

  // Get trend icon based on trend direction
  const getTrendIcon = (trend: 'up' | 'down' | 'neutral') => {
    switch (trend) {
      case 'up':
        return <TrendingUp className="h-4 w-4 text-green-500" />;
      case 'down':
        return <TrendingDown className="h-4 w-4 text-red-500" />;
      default:
        return <Minus className="h-4 w-4 text-gray-500" />;
    }
  };

  // Get trend color class
  const getTrendColor = (trend: 'up' | 'down' | 'neutral') => {
    switch (trend) {
      case 'up':
        return 'text-green-600';
      case 'down':
        return 'text-red-600';
      default:
        return 'text-gray-600';
    }
  };

  // Handle time range change
  const handleTimeRangeChange = (newRange: string) => {
    if (onTimeRangeChange) {
      onTimeRangeChange(newRange);
    }
  };

  const emailCards = getEmailDeliveryCards(metrics);

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Header with time range selector */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="min-w-0">
          <h2 className="text-xl sm:text-2xl font-bold tracking-tight">Email Delivery Monitoring</h2>
          <p className="text-sm sm:text-base text-muted-foreground">
            Track email delivery performance and engagement metrics
          </p>
        </div>
        <Select value={timeRange} onValueChange={handleTimeRangeChange}>
          <SelectTrigger className="w-full sm:w-[180px]">
            <Clock className="mr-2 h-4 w-4" />
            <SelectValue placeholder="Select time range" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="1h">Last Hour</SelectItem>
            <SelectItem value="24h">Last 24 Hours</SelectItem>
            <SelectItem value="7d">Last 7 Days</SelectItem>
            <SelectItem value="30d">Last 30 Days</SelectItem>
            <SelectItem value="90d">Last 90 Days</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {/* Error state */}
      {isError && (
        <V0ApiErrorFallback 
          message={error?.message || 'Failed to load email metrics'}
          onRetry={clearError}
        />
      )}

      {/* Loading state */}
      {loading && (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          {Array.from({ length: 4 }).map((_, index) => (
            <V0SkeletonCard key={index} />
          ))}
        </div>
      )}

      {/* Email delivery metrics cards */}
      {!loading && !isError && (
        <>
          {/* Show no data message if all metrics are zero */}
          {metrics.totalEmails === 0 && (
            <Card className="border-gray-200 bg-gray-50 mb-6">
              <CardContent className="pt-6">
                <div className="flex items-center space-x-2 text-gray-600">
                  <Mail className="h-5 w-5" />
                  <span className="font-medium">No Email Activity</span>
                </div>
                <p className="mt-2 text-sm text-gray-600">
                  No emails have been processed yet. Email metrics will appear once the system starts sending emails.
                </p>
              </CardContent>
            </Card>
          )}
          
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
            {emailCards.map((card, index) => {
            const IconComponent = card.icon;
            return (
              <Card key={index} className="relative overflow-hidden">
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium text-muted-foreground">
                    {card.title}
                  </CardTitle>
                  <IconComponent className={`h-4 w-4 ${card.color}`} />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">{card.value}</div>
                  <div className="mt-2 flex items-center space-x-1 text-xs">
                    {getTrendIcon(card.trend)}
                    <span className={getTrendColor(card.trend)}>
                      {card.change}
                    </span>
                  </div>
                  {/* Progress bar for percentage metrics */}
                  {typeof card.percentage === 'number' && card.title !== 'Failed Emails' && (
                    <div className="mt-3">
                      <V0ProgressBar 
                        value={card.percentage} 
                        size="sm"
                        color={card.trend === 'up' ? 'green' : card.trend === 'down' ? 'red' : 'blue'}
                      />
                    </div>
                  )}
                  {/* Special handling for failed emails - show as error progress */}
                  {card.title === 'Failed Emails' && (
                    <Progress 
                      value={card.percentage} 
                      className="mt-3 h-1 bg-red-100 [&>div]:bg-red-500" 
                    />
                  )}
                </CardContent>
              </Card>
            );
            })}
          </div>
        </>
      )}

      {/* Horizontal delivery status overview */}
      <Card>
        <CardHeader>
          <CardTitle>Email Delivery Status Overview</CardTitle>
          <CardDescription>
            Detailed breakdown of email delivery performance for the selected time period
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-6">
            {/* Delivered emails */}
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-2">
                  <MailCheck className="h-4 w-4 text-green-600" />
                  <span className="text-sm font-medium">Delivered</span>
                </div>
                <div className="flex items-center space-x-2">
                  <span className="text-sm text-muted-foreground">
                    {metrics.delivered} of {metrics.totalEmails}
                  </span>
                  <Badge variant="outline" className="bg-green-50 text-green-700 border-green-200">
                    {metrics.deliveryRate.toFixed(1)}%
                  </Badge>
                </div>
              </div>
              <Progress value={metrics.deliveryRate} className="h-2" />
            </div>

            {/* Opened emails */}
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-2">
                  <Mail className="h-4 w-4 text-blue-600" />
                  <span className="text-sm font-medium">Opened</span>
                </div>
                <div className="flex items-center space-x-2">
                  <span className="text-sm text-muted-foreground">
                    {metrics.opened} of {metrics.delivered}
                  </span>
                  <Badge variant="outline" className="bg-blue-50 text-blue-700 border-blue-200">
                    {metrics.openRate.toFixed(1)}%
                  </Badge>
                </div>
              </div>
              <Progress value={metrics.openRate} className="h-2" />
            </div>

            {/* Clicked emails */}
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-2">
                  <MousePointer className="h-4 w-4 text-purple-600" />
                  <span className="text-sm font-medium">Clicked</span>
                </div>
                <div className="flex items-center space-x-2">
                  <span className="text-sm text-muted-foreground">
                    {metrics.clicked} of {metrics.opened}
                  </span>
                  <Badge variant="outline" className="bg-purple-50 text-purple-700 border-purple-200">
                    {metrics.clickRate.toFixed(1)}%
                  </Badge>
                </div>
              </div>
              <Progress value={metrics.clickRate} className="h-2" />
            </div>

            {/* Failed emails breakdown */}
            {metrics.failedEmails > 0 && (
              <div className="space-y-4 pt-4 border-t">
                <h4 className="text-sm font-medium text-muted-foreground">Delivery Issues</h4>
                
                {/* Bounced emails */}
                <div className="space-y-2">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-2">
                      <AlertTriangle className="h-4 w-4 text-amber-600" />
                      <span className="text-sm font-medium">Bounced</span>
                    </div>
                    <div className="flex items-center space-x-2">
                      <span className="text-sm text-muted-foreground">{metrics.bounced}</span>
                      <Badge variant="outline" className="bg-amber-50 text-amber-700 border-amber-200">
                        {metrics.totalEmails > 0 ? ((metrics.bounced / metrics.totalEmails) * 100).toFixed(1) : 0}%
                      </Badge>
                    </div>
                  </div>
                  <Progress 
                    value={metrics.totalEmails > 0 ? (metrics.bounced / metrics.totalEmails) * 100 : 0} 
                    className="h-2 bg-amber-100 [&>div]:bg-amber-500" 
                  />
                </div>

                {/* Spam emails */}
                <div className="space-y-2">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-2">
                      <AlertTriangle className="h-4 w-4 text-red-600" />
                      <span className="text-sm font-medium">Marked as Spam</span>
                    </div>
                    <div className="flex items-center space-x-2">
                      <span className="text-sm text-muted-foreground">{metrics.spam}</span>
                      <Badge variant="outline" className="bg-red-50 text-red-700 border-red-200">
                        {metrics.totalEmails > 0 ? ((metrics.spam / metrics.totalEmails) * 100).toFixed(1) : 0}%
                      </Badge>
                    </div>
                  </div>
                  <Progress 
                    value={metrics.totalEmails > 0 ? (metrics.spam / metrics.totalEmails) * 100 : 0} 
                    className="h-2 bg-red-100 [&>div]:bg-red-500" 
                  />
                </div>
              </div>
            )}
          </div>
        </CardContent>
      </Card>

      {/* Loading state overlay */}
      {loading && (
        <div className="absolute inset-0 bg-white/50 flex items-center justify-center">
          <div className="flex items-center space-x-2">
            <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600"></div>
            <span className="text-sm text-muted-foreground">Loading email metrics...</span>
          </div>
        </div>
      )}
    </div>
  );
};

export default V0EmailDeliveryDashboard;