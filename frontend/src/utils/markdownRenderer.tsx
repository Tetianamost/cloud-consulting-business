import React from 'react';

interface MarkdownRendererProps {
  content: string;
  className?: string;
}

export const MarkdownRenderer: React.FC<MarkdownRendererProps> = ({ content, className = '' }) => {
  // Simple markdown renderer for basic formatting
  const renderMarkdown = (text: string): React.ReactElement => {
    // Split by lines to handle different markdown elements
    const lines = text.split('\n');
    const elements: React.ReactElement[] = [];
    let currentList: string[] = [];
    let inCodeBlock = false;
    let codeBlockContent: string[] = [];

    const flushList = () => {
      if (currentList.length > 0) {
        elements.push(
          <ul key={`list-${elements.length}`} className="list-disc list-inside mb-4 space-y-1">
            {currentList.map((item, idx) => (
              <li key={idx} className="text-sm leading-relaxed">
                {renderInlineMarkdown(item)}
              </li>
            ))}
          </ul>
        );
        currentList = [];
      }
    };

    const flushCodeBlock = () => {
      if (codeBlockContent.length > 0) {
        elements.push(
          <pre key={`code-${elements.length}`} className="bg-gray-100 p-3 rounded text-sm overflow-x-auto mb-4">
            <code>{codeBlockContent.join('\n')}</code>
          </pre>
        );
        codeBlockContent = [];
      }
    };

    lines.forEach((line, index) => {
      // Handle code blocks
      if (line.startsWith('```')) {
        if (inCodeBlock) {
          flushCodeBlock();
          inCodeBlock = false;
        } else {
          flushList();
          inCodeBlock = true;
        }
        return;
      }

      if (inCodeBlock) {
        codeBlockContent.push(line);
        return;
      }

      // Handle headers
      if (line.startsWith('### ')) {
        flushList();
        elements.push(
          <h3 key={`h3-${index}`} className="text-lg font-semibold mt-6 mb-3 text-gray-800">
            {line.substring(4)}
          </h3>
        );
      } else if (line.startsWith('## ')) {
        flushList();
        elements.push(
          <h2 key={`h2-${index}`} className="text-xl font-bold mt-6 mb-4 text-gray-800">
            {line.substring(3)}
          </h2>
        );
      } else if (line.startsWith('# ')) {
        flushList();
        elements.push(
          <h1 key={`h1-${index}`} className="text-2xl font-bold mt-6 mb-4 text-gray-800">
            {line.substring(2)}
          </h1>
        );
      }
      // Handle list items
      else if (line.startsWith('- ') || line.startsWith('* ')) {
        currentList.push(line.substring(2));
      }
      // Handle numbered lists
      else if (/^\d+\.\s/.test(line)) {
        currentList.push(line.replace(/^\d+\.\s/, ''));
      }
      // Handle empty lines
      else if (line.trim() === '') {
        flushList();
        if (elements.length > 0 && elements[elements.length - 1].type !== 'br') {
          elements.push(<br key={`br-${index}`} />);
        }
      }
      // Handle regular paragraphs
      else {
        flushList();
        if (line.trim()) {
          elements.push(
            <p key={`p-${index}`} className="mb-3 text-sm leading-relaxed">
              {renderInlineMarkdown(line)}
            </p>
          );
        }
      }
    });

    // Flush any remaining content
    flushList();
    flushCodeBlock();

    return <div className={className}>{elements}</div>;
  };

  const renderInlineMarkdown = (text: string): React.ReactElement => {
    // Handle bold text
    let result = text.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>');
    
    // Handle italic text
    result = result.replace(/\*(.*?)\*/g, '<em>$1</em>');
    
    // Handle inline code
    result = result.replace(/`(.*?)`/g, '<code class="bg-gray-100 px-1 py-0.5 rounded text-xs">$1</code>');
    
    return <span dangerouslySetInnerHTML={{ __html: result }} />;
  };

  return renderMarkdown(content);
};

export default MarkdownRenderer;