import React, { useState, useEffect } from 'react';

interface Message {
  id: string;
  type: 'user' | 'assistant';
  content: string;
  timestamp: string;
}

export const SimpleWorkingChat: React.FC = () => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState('');
  const [loading, setLoading] = useState(false);
  const [sessionId] = useState(`session-${Date.now()}-${Math.random().toString(36).substring(2, 9)}`);

  const sendMessage = async () => {
    if (!input.trim() || loading) {
      console.log('Send blocked:', { inputEmpty: !input.trim(), loading });
      return;
    }

    console.log('Sending message:', input.trim());

    const userMessage: Message = {
      id: `user-${Date.now()}`,
      type: 'user',
      content: input.trim(),
      timestamp: new Date().toISOString()
    };

    setMessages(prev => [...prev, userMessage]);
    setInput('');
    setLoading(true);

    try {
      // Send message to backend
      const token = localStorage.getItem('token');
      console.log('Token:', token ? 'exists' : 'missing');
      
      const response = await fetch('/api/v1/admin/chat/messages', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
          content: userMessage.content,
          session_id: sessionId
        })
      });

      console.log('Response status:', response.status);

      if (!response.ok) {
        const errorText = await response.text();
        console.error('Response error:', errorText);
        throw new Error(`Failed to send message: ${response.status} ${errorText}`);
      }

      const result = await response.json();
      console.log('Response result:', result);
      
      if (result.success && result.content) {
        // Add AI response directly from the send response
        const aiMessage: Message = {
          id: result.message_id,
          type: 'assistant',
          content: result.content,
          timestamp: new Date().toISOString()
        };
        setMessages(prev => [...prev, aiMessage]);
        console.log('AI response added:', aiMessage.content);
      } else {
        console.warn('No AI response in result:', result);
      }

    } catch (error) {
      console.error('Failed to send message:', error);
      // Add error message to chat
      const errorMessage: Message = {
        id: `error-${Date.now()}`,
        type: 'assistant',
        content: `Error: ${error instanceof Error ? error.message : 'Unknown error'}`,
        timestamp: new Date().toISOString()
      };
      setMessages(prev => [...prev, errorMessage]);
    } finally {
      setLoading(false);
      console.log('Loading set to false');
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  };

  return (
    <div style={{ 
      display: 'flex', 
      flexDirection: 'column', 
      height: '500px', 
      border: '1px solid #ccc', 
      borderRadius: '8px',
      backgroundColor: 'white'
    }}>
      <div style={{ 
        padding: '16px', 
        borderBottom: '1px solid #eee',
        backgroundColor: '#f8f9fa',
        fontWeight: 'bold'
      }}>
        🤖 Simple Working Chat
      </div>
      
      <div style={{ 
        flex: 1, 
        padding: '16px', 
        overflowY: 'auto',
        display: 'flex',
        flexDirection: 'column',
        gap: '12px'
      }}>
        {messages.map((message) => (
          <div
            key={message.id}
            style={{
              padding: '12px',
              borderRadius: '8px',
              backgroundColor: message.type === 'user' ? '#e3f2fd' : '#f5f5f5',
              alignSelf: message.type === 'user' ? 'flex-end' : 'flex-start',
              maxWidth: '80%',
              border: `1px solid ${message.type === 'user' ? '#2196f3' : '#ddd'}`
            }}
          >
            <div style={{ fontWeight: 'bold', marginBottom: '4px', fontSize: '12px' }}>
              {message.type === 'user' ? '👤 You' : '🤖 AI Assistant'}
            </div>
            <div>{message.content}</div>
          </div>
        ))}
        
        {loading && (
          <div style={{
            padding: '12px',
            borderRadius: '8px',
            backgroundColor: '#f5f5f5',
            alignSelf: 'flex-start',
            border: '1px solid #ddd'
          }}>
            <div style={{ fontWeight: 'bold', marginBottom: '4px', fontSize: '12px' }}>
              🤖 AI Assistant
            </div>
            <div>Thinking...</div>
          </div>
        )}
      </div>
      
      <div style={{ 
        padding: '16px', 
        borderTop: '1px solid #eee',
        display: 'flex',
        gap: '8px'
      }}>
        <input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyPress={handleKeyPress}
          placeholder="Type your message..."
          disabled={loading}
          style={{
            flex: 1,
            padding: '8px 12px',
            border: '1px solid #ddd',
            borderRadius: '4px',
            fontSize: '14px'
          }}
        />
        <button
          onClick={sendMessage}
          disabled={loading || !input.trim()}
          style={{
            padding: '8px 16px',
            backgroundColor: '#2196f3',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            cursor: loading ? 'not-allowed' : 'pointer',
            opacity: loading || !input.trim() ? 0.6 : 1
          }}
        >
          {loading ? '...' : 'Send'}
        </button>
      </div>
    </div>
  );
};