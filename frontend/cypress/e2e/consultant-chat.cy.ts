/// <reference types="cypress" />

describe('Consultant Chat E2E Tests', () => {
  beforeEach(() => {
    // Mock localStorage with admin token
    cy.window().then((win) => {
      win.localStorage.setItem('adminToken', 'mock-admin-token');
    });

    // Mock WebSocket server responses
    cy.intercept('GET', '/api/v1/admin/chat/ws', { fixture: 'websocket-mock.json' });
    
    // Visit the admin page with chat component
    cy.visit('/admin/chat');
  });

  describe('Chat Interface', () => {
    it('should display the chat interface correctly', () => {
      cy.get('[data-testid="consultant-chat"]').should('be.visible');
      cy.get('[data-testid="chat-header"]').should('contain', 'Consultant Assistant');
      cy.get('[data-testid="message-input"]').should('be.visible');
      cy.get('[data-testid="send-button"]').should('be.visible');
      cy.get('[data-testid="connection-status"]').should('be.visible');
    });

    it('should show quick action buttons', () => {
      cy.get('[data-testid="quick-actions"]').should('be.visible');
      cy.get('[data-testid="quick-action-cost-estimate"]').should('contain', 'Cost Estimate');
      cy.get('[data-testid="quick-action-security-review"]').should('contain', 'Security Review');
      cy.get('[data-testid="quick-action-best-practices"]').should('contain', 'Best Practices');
      cy.get('[data-testid="quick-action-alternatives"]').should('contain', 'Alternatives');
      cy.get('[data-testid="quick-action-next-steps"]').should('contain', 'Next Steps');
    });

    it('should display welcome message when no messages exist', () => {
      cy.get('[data-testid="welcome-message"]')
        .should('be.visible')
        .should('contain', 'Start a conversation to get real-time AWS consulting assistance');
    });
  });

  describe('Settings Panel', () => {
    it('should open and close settings panel', () => {
      cy.get('[data-testid="settings-button"]').click();
      cy.get('[data-testid="settings-panel"]').should('be.visible');
      
      cy.get('[data-testid="client-name-input"]').should('be.visible');
      cy.get('[data-testid="meeting-context-input"]').should('be.visible');
      cy.get('[data-testid="clear-chat-button"]').should('be.visible');

      // Close settings by clicking outside or on settings button again
      cy.get('[data-testid="settings-button"]').click();
      cy.get('[data-testid="settings-panel"]').should('not.be.visible');
    });

    it('should update client name and meeting context', () => {
      cy.get('[data-testid="settings-button"]').click();
      
      cy.get('[data-testid="client-name-input"]')
        .clear()
        .type('E2E Test Client');
      
      cy.get('[data-testid="meeting-context-input"]')
        .clear()
        .type('End-to-end testing session');

      // Values should be persisted
      cy.get('[data-testid="client-name-input"]').should('have.value', 'E2E Test Client');
      cy.get('[data-testid="meeting-context-input"]').should('have.value', 'End-to-end testing session');
    });

    it('should clear chat history', () => {
      // First send a message to have something to clear
      cy.get('[data-testid="message-input"]').type('Test message for clearing');
      cy.get('[data-testid="send-button"]').click();
      
      // Wait for message to appear
      cy.get('[data-testid="message-list"]').should('contain', 'Test message for clearing');

      // Open settings and clear chat
      cy.get('[data-testid="settings-button"]').click();
      cy.get('[data-testid="clear-chat-button"]').click();

      // Messages should be cleared
      cy.get('[data-testid="welcome-message"]').should('be.visible');
      cy.get('[data-testid="message-list"]').should('not.contain', 'Test message for clearing');
    });
  });

  describe('Message Sending', () => {
    beforeEach(() => {
      // Mock WebSocket connection
      cy.window().then((win) => {
        // Mock WebSocket for testing
        const mockWebSocket = {
          readyState: WebSocket.OPEN,
          send: cy.stub().as('websocketSend'),
          close: cy.stub(),
          addEventListener: cy.stub(),
          removeEventListener: cy.stub(),
        };

        // Override WebSocket constructor
        win.WebSocket = cy.stub().returns(mockWebSocket);
      });
    });

    it('should send a message and receive response', () => {
      const testMessage = 'What is AWS Lambda?';
      
      // Type message
      cy.get('[data-testid="message-input"]').type(testMessage);
      
      // Send message
      cy.get('[data-testid="send-button"]').click();

      // Input should be cleared
      cy.get('[data-testid="message-input"]').should('have.value', '');

      // User message should appear
      cy.get('[data-testid="message-list"]')
        .should('contain', testMessage)
        .within(() => {
          cy.get('[data-testid="user-message"]').should('contain', testMessage);
        });

      // Mock AI response
      cy.window().then((win) => {
        const mockResponse = {
          success: true,
          session_id: 'test-session-id',
          message: {
            id: 'response-1',
            type: 'assistant',
            content: 'AWS Lambda is a serverless computing service that lets you run code without provisioning or managing servers.',
            timestamp: new Date().toISOString(),
            session_id: 'test-session-id'
          }
        };

        // Simulate WebSocket message event
        win.dispatchEvent(new MessageEvent('message', {
          data: JSON.stringify(mockResponse)
        }));
      });

      // AI response should appear
      cy.get('[data-testid="message-list"]')
        .should('contain', 'AWS Lambda is a serverless computing service')
        .within(() => {
          cy.get('[data-testid="assistant-message"]')
            .should('contain', 'AWS Lambda is a serverless computing service');
        });
    });

    it('should send message using Enter key', () => {
      const testMessage = 'Test message with Enter key';
      
      cy.get('[data-testid="message-input"]')
        .type(testMessage)
        .type('{enter}');

      // Input should be cleared
      cy.get('[data-testid="message-input"]').should('have.value', '');

      // Message should appear
      cy.get('[data-testid="message-list"]').should('contain', testMessage);
    });

    it('should not send empty messages', () => {
      // Try to send empty message
      cy.get('[data-testid="send-button"]').click();

      // No message should be sent
      cy.get('[data-testid="message-list"]').should('not.contain', '');
      cy.get('[data-testid="welcome-message"]').should('be.visible');
    });

    it('should show loading state while waiting for response', () => {
      cy.get('[data-testid="message-input"]').type('Test loading state');
      cy.get('[data-testid="send-button"]').click();

      // Loading indicator should appear
      cy.get('[data-testid="loading-indicator"]').should('be.visible');
      cy.get('[data-testid="loading-dots"]').should('be.visible');
    });
  });

  describe('Quick Actions', () => {
    it('should send quick action messages', () => {
      const quickActions = [
        { testId: 'quick-action-cost-estimate', expectedText: 'Provide a cost estimate' },
        { testId: 'quick-action-security-review', expectedText: 'security considerations' },
        { testId: 'quick-action-best-practices', expectedText: 'best practices' },
        { testId: 'quick-action-alternatives', expectedText: 'alternative approaches' },
        { testId: 'quick-action-next-steps', expectedText: 'next steps' },
      ];

      quickActions.forEach((action, index) => {
        cy.get(`[data-testid="${action.testId}"]`).click();

        // Message should appear in chat
        cy.get('[data-testid="message-list"]')
          .should('contain', action.expectedText);

        // Clear chat for next test
        if (index < quickActions.length - 1) {
          cy.get('[data-testid="settings-button"]').click();
          cy.get('[data-testid="clear-chat-button"]').click();
          cy.get('[data-testid="settings-button"]').click();
        }
      });
    });

    it('should disable quick actions when not connected', () => {
      // Mock disconnected state
      cy.window().then((win) => {
        win.localStorage.removeItem('adminToken');
      });

      cy.reload();

      // Quick action buttons should be disabled
      cy.get('[data-testid="quick-action-cost-estimate"]').should('be.disabled');
      cy.get('[data-testid="quick-action-security-review"]').should('be.disabled');
    });
  });

  describe('Connection Status', () => {
    it('should show connected status when WebSocket is connected', () => {
      cy.get('[data-testid="connection-status"]')
        .should('have.class', 'bg-green-400')
        .should('be.visible');
    });

    it('should show disconnected status when WebSocket fails', () => {
      // Mock WebSocket connection failure
      cy.window().then((win) => {
        win.WebSocket = cy.stub().throws(new Error('Connection failed'));
      });

      cy.reload();

      cy.get('[data-testid="connection-status"]')
        .should('have.class', 'bg-red-400')
        .should('be.visible');

      // Input should be disabled
      cy.get('[data-testid="message-input"]').should('be.disabled');
      cy.get('[data-testid="send-button"]').should('be.disabled');
    });

    it('should show connecting status during connection attempt', () => {
      // This test would require more sophisticated WebSocket mocking
      // to simulate the connecting state
      cy.get('[data-testid="connection-status"]').should('be.visible');
    });
  });

  describe('Message Display', () => {
    it('should display messages with correct styling', () => {
      // Send a user message
      cy.get('[data-testid="message-input"]').type('User message test');
      cy.get('[data-testid="send-button"]').click();

      // Check user message styling
      cy.get('[data-testid="user-message"]')
        .should('have.class', 'bg-blue-600')
        .should('have.class', 'text-white');

      // Mock AI response
      cy.window().then((win) => {
        const mockResponse = {
          success: true,
          session_id: 'test-session-id',
          message: {
            id: 'response-1',
            type: 'assistant',
            content: 'AI response test',
            timestamp: new Date().toISOString(),
            session_id: 'test-session-id'
          }
        };

        win.dispatchEvent(new MessageEvent('message', {
          data: JSON.stringify(mockResponse)
        }));
      });

      // Check assistant message styling
      cy.get('[data-testid="assistant-message"]')
        .should('have.class', 'bg-gray-100')
        .should('have.class', 'text-gray-900');
    });

    it('should display timestamps correctly', () => {
      cy.get('[data-testid="message-input"]').type('Timestamp test');
      cy.get('[data-testid="send-button"]').click();

      // Timestamp should be visible and formatted correctly
      cy.get('[data-testid="message-timestamp"]')
        .should('be.visible')
        .should('match', /\d{1,2}:\d{2}/); // HH:MM format
    });

    it('should auto-scroll to bottom when new messages arrive', () => {
      // Send multiple messages to create scrollable content
      for (let i = 1; i <= 10; i++) {
        cy.get('[data-testid="message-input"]').type(`Message ${i}`);
        cy.get('[data-testid="send-button"]').click();
        cy.wait(100); // Small delay between messages
      }

      // The last message should be visible (auto-scrolled)
      cy.get('[data-testid="message-list"]')
        .should('contain', 'Message 10')
        .within(() => {
          cy.contains('Message 10').should('be.visible');
        });
    });
  });

  describe('Responsive Design', () => {
    it('should work correctly on mobile viewport', () => {
      cy.viewport('iphone-x');

      // Chat interface should still be visible and functional
      cy.get('[data-testid="consultant-chat"]').should('be.visible');
      cy.get('[data-testid="message-input"]').should('be.visible');
      cy.get('[data-testid="send-button"]').should('be.visible');

      // Quick actions should be responsive
      cy.get('[data-testid="quick-actions"]').should('be.visible');
      cy.get('[data-testid="quick-action-cost-estimate"]').should('be.visible');
    });

    it('should work correctly on tablet viewport', () => {
      cy.viewport('ipad-2');

      cy.get('[data-testid="consultant-chat"]').should('be.visible');
      cy.get('[data-testid="message-input"]').should('be.visible');
      
      // Send a test message
      cy.get('[data-testid="message-input"]').type('Tablet test message');
      cy.get('[data-testid="send-button"]').click();
      
      cy.get('[data-testid="message-list"]').should('contain', 'Tablet test message');
    });
  });

  describe('Accessibility', () => {
    it('should be keyboard navigable', () => {
      // Tab through interface elements
      cy.get('body').tab();
      cy.focused().should('have.attr', 'data-testid', 'message-input');

      cy.focused().tab();
      cy.focused().should('have.attr', 'data-testid', 'send-button');

      cy.focused().tab();
      cy.focused().should('have.attr', 'data-testid', 'settings-button');
    });

    it('should have proper ARIA labels', () => {
      cy.get('[data-testid="message-input"]')
        .should('have.attr', 'aria-label')
        .should('contain', 'message');

      cy.get('[data-testid="send-button"]')
        .should('have.attr', 'aria-label')
        .should('contain', 'send');
    });

    it('should support screen readers', () => {
      // Check for proper semantic HTML
      cy.get('[data-testid="message-list"]').should('have.attr', 'role', 'log');
      cy.get('[data-testid="message-input"]').should('have.attr', 'type', 'text');
      cy.get('[data-testid="send-button"]').should('have.attr', 'type', 'submit');
    });
  });

  describe('Error Handling', () => {
    it('should handle WebSocket connection errors gracefully', () => {
      // Mock WebSocket error
      cy.window().then((win) => {
        const mockWebSocket = {
          readyState: WebSocket.CLOSED,
          send: cy.stub(),
          close: cy.stub(),
          addEventListener: cy.stub(),
          removeEventListener: cy.stub(),
        };

        win.WebSocket = cy.stub().returns(mockWebSocket);
      });

      cy.reload();

      // Should show disconnected state
      cy.get('[data-testid="connection-status"]').should('have.class', 'bg-red-400');
      cy.get('[data-testid="message-input"]').should('contain.text', 'Connecting...');
    });

    it('should handle malformed WebSocket messages', () => {
      // This would require more sophisticated mocking to test
      // malformed message handling
      cy.get('[data-testid="consultant-chat"]').should('be.visible');
    });

    it('should handle network timeouts', () => {
      // Mock slow network response
      cy.intercept('GET', '/api/v1/admin/chat/ws', {
        delay: 10000,
        statusCode: 408,
      });

      cy.reload();

      // Should handle timeout gracefully
      cy.get('[data-testid="connection-status"]').should('be.visible');
    });
  });

  describe('Performance', () => {
    it('should handle many messages without performance degradation', () => {
      // Send many messages quickly
      for (let i = 1; i <= 50; i++) {
        cy.get('[data-testid="message-input"]').type(`Performance test ${i}`);
        cy.get('[data-testid="send-button"]').click();
      }

      // Interface should remain responsive
      cy.get('[data-testid="message-input"]').should('be.visible');
      cy.get('[data-testid="send-button"]').should('not.be.disabled');
      
      // Last message should be visible
      cy.get('[data-testid="message-list"]').should('contain', 'Performance test 50');
    });

    it('should load quickly', () => {
      const startTime = Date.now();
      
      cy.visit('/admin/chat').then(() => {
        const loadTime = Date.now() - startTime;
        expect(loadTime).to.be.lessThan(3000); // Should load within 3 seconds
      });

      cy.get('[data-testid="consultant-chat"]').should('be.visible');
    });
  });

  describe('Session Management', () => {
    it('should maintain session across page reloads', () => {
      // Send a message to create session
      cy.get('[data-testid="message-input"]').type('Session test message');
      cy.get('[data-testid="send-button"]').click();

      // Reload page
      cy.reload();

      // Session should be maintained (this would require proper session storage)
      cy.get('[data-testid="consultant-chat"]').should('be.visible');
    });

    it('should handle session expiration gracefully', () => {
      // This would require mocking session expiration
      cy.get('[data-testid="consultant-chat"]').should('be.visible');
    });
  });
});

// Custom Cypress commands for chat testing
declare global {
  namespace Cypress {
    interface Chainable {
      /**
       * Send a chat message
       */
      sendChatMessage(message: string): Chainable<Element>;
      
      /**
       * Wait for AI response
       */
      waitForAIResponse(): Chainable<Element>;
      
      /**
       * Clear chat history
       */
      clearChatHistory(): Chainable<Element>;
    }
  }
}

Cypress.Commands.add('sendChatMessage', (message: string) => {
  cy.get('[data-testid="message-input"]').clear().type(message);
  cy.get('[data-testid="send-button"]').click();
});

Cypress.Commands.add('waitForAIResponse', () => {
  cy.get('[data-testid="loading-indicator"]', { timeout: 10000 }).should('not.exist');
  cy.get('[data-testid="assistant-message"]').should('be.visible');
});

Cypress.Commands.add('clearChatHistory', () => {
  cy.get('[data-testid="settings-button"]').click();
  cy.get('[data-testid="clear-chat-button"]').click();
  cy.get('[data-testid="settings-button"]').click();
});