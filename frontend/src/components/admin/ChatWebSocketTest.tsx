import React, { useEffect, useRef } from "react";

// Try to get JWT token from localStorage (adminToken)
function getToken() {
  return localStorage.getItem("adminToken") || "";
}

const WS_BASE = "ws://localhost:8061/api/v1/admin/chat/ws";

const TEST_SESSION_ID = "";
const TEST_METADATA = {
  client_name: "Test Client",
  context: "Testing chat integration from frontend",
};

export const ChatWebSocketTest: React.FC = () => {
  const ws = useRef<WebSocket | null>(null);

  useEffect(() => {
    const token = getToken();
    const wsUrl =
      token && token.length > 0
        ? `${WS_BASE}?token=${encodeURIComponent(token)}`
        : WS_BASE;

    ws.current = new window.WebSocket(wsUrl);

    ws.current.onopen = () => {
      console.log("[ChatWebSocketTest] WebSocket opened");

      ws.current?.send(
        JSON.stringify({
          type: "message",
          session_id: TEST_SESSION_ID,
          message_id: "",
          content: "Hello from frontend (full protocol)!",
          metadata: TEST_METADATA,
          timestamp: new Date().toISOString(),
        })
      );
    };

    ws.current.onmessage = (event) => {
      console.log("[ChatWebSocketTest] Received:", event.data);
    };

    ws.current.onerror = (event) => {
      console.error("[ChatWebSocketTest] WebSocket error:", event);
    };

    ws.current.onclose = (event) => {
      console.log("[ChatWebSocketTest] WebSocket closed:", event);
    };

    return () => {
      ws.current?.close();
    };
  }, []);

  return <div>Chat WebSocket Test Component (see console for logs)</div>;
};