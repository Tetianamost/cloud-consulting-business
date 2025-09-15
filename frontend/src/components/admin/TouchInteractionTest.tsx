import React, { useState, useRef, useEffect } from 'react';
import { Button } from '../ui/Button';
import { Card, CardContent, CardHeader, CardTitle } from '../ui/Card';
import { Badge } from '../ui/Badge';
import { 
  Hand, 
  MousePointer, 
  Smartphone, 
  CheckCircle, 
  XCircle,
  AlertTriangle 
} from 'lucide-react';

interface TouchTestResult {
  element: string;
  size: { width: number; height: number };
  passed: boolean;
  recommendation?: string;
}

/**
 * TouchInteractionTest - Component to test touch interaction compliance
 * Verifies that interactive elements meet mobile usability standards
 */
const TouchInteractionTest: React.FC = () => {
  const [testResults, setTestResults] = useState<TouchTestResult[]>([]);
  const [isRunning, setIsRunning] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);

  // Minimum touch target size (44px x 44px is Apple's recommendation)
  const MIN_TOUCH_SIZE = 44;

  const runTouchTests = () => {
    setIsRunning(true);
    
    setTimeout(() => {
      const results: TouchTestResult[] = [];
      
      if (containerRef.current) {
        // Test buttons
        const buttons = containerRef.current.querySelectorAll('button');
        buttons.forEach((button, index) => {
          const rect = button.getBoundingClientRect();
          const passed = rect.width >= MIN_TOUCH_SIZE && rect.height >= MIN_TOUCH_SIZE;
          
          results.push({
            element: `Button ${index + 1}`,
            size: { width: Math.round(rect.width), height: Math.round(rect.height) },
            passed,
            recommendation: !passed ? `Increase size to ${MIN_TOUCH_SIZE}px minimum` : undefined
          });
        });

        // Test links
        const links = containerRef.current.querySelectorAll('a');
        links.forEach((link, index) => {
          const rect = link.getBoundingClientRect();
          const passed = rect.width >= MIN_TOUCH_SIZE && rect.height >= MIN_TOUCH_SIZE;
          
          results.push({
            element: `Link ${index + 1}`,
            size: { width: Math.round(rect.width), height: Math.round(rect.height) },
            passed,
            recommendation: !passed ? `Add padding to reach ${MIN_TOUCH_SIZE}px` : undefined
          });
        });

        // Test select triggers
        const selects = containerRef.current.querySelectorAll('[role="combobox"]');
        selects.forEach((select, index) => {
          const rect = select.getBoundingClientRect();
          const passed = rect.height >= MIN_TOUCH_SIZE;
          
          results.push({
            element: `Select ${index + 1}`,
            size: { width: Math.round(rect.width), height: Math.round(rect.height) },
            passed,
            recommendation: !passed ? `Increase height to ${MIN_TOUCH_SIZE}px` : undefined
          });
        });
      }
      
      setTestResults(results);
      setIsRunning(false);
    }, 1000);
  };

  const getStatusIcon = (passed: boolean) => {
    return passed ? (
      <CheckCircle className="h-4 w-4 text-green-600" />
    ) : (
      <XCircle className="h-4 w-4 text-red-600" />
    );
  };

  const passedCount = testResults.filter(r => r.passed).length;
  const totalCount = testResults.length;

  return (
    <div ref={containerRef} className="p-6 space-y-6">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center space-x-2">
            <Hand className="h-5 w-5" />
            <span>Touch Interaction Test</span>
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-2">
                <Smartphone className="h-4 w-4 text-blue-600" />
                <span className="text-sm">Testing for mobile usability compliance</span>
              </div>
              <Button 
                onClick={runTouchTests} 
                disabled={isRunning}
                size="sm"
              >
                {isRunning ? 'Testing...' : 'Run Tests'}
              </Button>
            </div>

            {testResults.length > 0 && (
              <div className="space-y-4">
                <div className="flex items-center space-x-4">
                  <Badge 
                    variant={passedCount === totalCount ? "default" : "destructive"}
                    className="flex items-center space-x-1"
                  >
                    {passedCount === totalCount ? (
                      <CheckCircle className="h-3 w-3" />
                    ) : (
                      <AlertTriangle className="h-3 w-3" />
                    )}
                    <span>{passedCount}/{totalCount} Passed</span>
                  </Badge>
                  <span className="text-sm text-gray-600">
                    Minimum touch target: {MIN_TOUCH_SIZE}px × {MIN_TOUCH_SIZE}px
                  </span>
                </div>

                <div className="space-y-2">
                  {testResults.map((result, index) => (
                    <div 
                      key={index}
                      className={`flex items-center justify-between p-3 rounded-lg border ${
                        result.passed ? 'bg-green-50 border-green-200' : 'bg-red-50 border-red-200'
                      }`}
                    >
                      <div className="flex items-center space-x-3">
                        {getStatusIcon(result.passed)}
                        <div>
                          <div className="font-medium text-sm">{result.element}</div>
                          <div className="text-xs text-gray-600">
                            {result.size.width}px × {result.size.height}px
                          </div>
                        </div>
                      </div>
                      {result.recommendation && (
                        <div className="text-xs text-red-600 max-w-xs text-right">
                          {result.recommendation}
                        </div>
                      )}
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        </CardContent>
      </Card>

      {/* Sample interactive elements for testing */}
      <Card>
        <CardHeader>
          <CardTitle>Sample Interactive Elements</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex flex-wrap gap-2">
              <Button size="sm">Small Button</Button>
              <Button>Default Button</Button>
              <Button size="lg">Large Button</Button>
            </div>
            
            <div className="flex flex-wrap gap-2">
              <a href="#" className="text-blue-600 hover:underline text-sm">Small Link</a>
              <a href="#" className="text-blue-600 hover:underline">Default Link</a>
              <a href="#" className="text-blue-600 hover:underline text-lg">Large Link</a>
            </div>

            <div className="flex flex-wrap gap-2">
              <select className="border border-gray-300 rounded px-2 py-1 text-sm">
                <option>Small Select</option>
              </select>
              <select className="border border-gray-300 rounded px-3 py-2">
                <option>Default Select</option>
              </select>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

export default TouchInteractionTest;