import React, { useState, useEffect, useMemo } from 'react';
import { 
  Plus, 
  History, 
  Search, 
  Filter, 
  Download, 
  Share2, 
  Trash2, 
  Clock, 
  User, 
  MessageSquare, 
  Calendar,
  BarChart3,
  FileText,
  X,
  ChevronDown,
  ChevronRight,
  Archive,
  Star,
  Tag
} from 'lucide-react';
import { useDispatch, useSelector } from 'react-redux';
import type { AppDispatch } from '../../store';
import { RootState } from '../../store';
import {
  setCurrentSession,
  createSession,
  loadSessionHistory,
  ChatSession,
  SessionContext,
} from '../../store/slices/chatSlice';

interface ChatSessionManagerProps {
  isOpen: boolean;
  onClose: () => void;
  onSessionSelect: (session: ChatSession) => void;
}

interface SessionFilters {
  search: string;
  status: 'all' | 'active' | 'inactive' | 'expired';
  dateRange: 'all' | 'today' | 'week' | 'month';
  clientName: string;
}

interface SessionStats {
  totalSessions: number;
  activeSessions: number;
  totalMessages: number;
  avgSessionDuration: number;
  topClients: Array<{ name: string; count: number }>;
}

export const ChatSessionManager: React.FC<ChatSessionManagerProps> = ({
  isOpen,
  onClose,
  onSessionSelect,
}) => {
  const dispatch = useDispatch<AppDispatch>();
  const { sessions, currentSession, messageHistory } = useSelector((state: RootState) => state.chat);
  
  // Local state
  const [filters, setFilters] = useState<SessionFilters>({
    search: '',
    status: 'all',
    dateRange: 'all',
    clientName: '',
  });
  const [selectedSessions, setSelectedSessions] = useState<Set<string>>(new Set());
  const [showCreateSession, setShowCreateSession] = useState(false);
  const [expandedSessions, setExpandedSessions] = useState<Set<string>>(new Set());
  const [sortBy, setSortBy] = useState<'created_at' | 'last_activity' | 'client_name'>('last_activity');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');

  // New session form state
  const [newSessionData, setNewSessionData] = useState<Partial<SessionContext>>({
    client_name: '',
    meeting_type: '',
    project_context: '',
  });

  // Load sessions on mount
  useEffect(() => {
    // In a real implementation, you would fetch sessions from the API
    // For now, we'll work with the sessions in the Redux store
  }, []);

  // Calculate session statistics
  const sessionStats = useMemo((): SessionStats => {
    const totalSessions = sessions.length;
    const activeSessions = sessions.filter(s => s.status === 'active').length;
    const totalMessages = Object.values(messageHistory).reduce((sum, messages) => sum + messages.length, 0);
    
    // Calculate average session duration (mock calculation)
    const avgSessionDuration = sessions.reduce((sum, session) => {
      const created = new Date(session.created_at).getTime();
      const lastActivity = new Date(session.last_activity).getTime();
      return sum + (lastActivity - created);
    }, 0) / (sessions.length || 1) / (1000 * 60); // Convert to minutes

    // Get top clients
    const clientCounts = sessions.reduce((acc, session) => {
      const client = session.client_name || 'Unknown';
      acc[client] = (acc[client] || 0) + 1;
      return acc;
    }, {} as Record<string, number>);

    const topClients = Object.entries(clientCounts)
      .map(([name, count]) => ({ name, count }))
      .sort((a, b) => b.count - a.count)
      .slice(0, 5);

    return {
      totalSessions,
      activeSessions,
      totalMessages,
      avgSessionDuration,
      topClients,
    };
  }, [sessions, messageHistory]);

  // Filter and sort sessions
  const filteredSessions = useMemo(() => {
    let filtered = sessions.filter(session => {
      // Search filter
      if (filters.search) {
        const searchLower = filters.search.toLowerCase();
        const matchesSearch = 
          session.client_name?.toLowerCase().includes(searchLower) ||
          session.context?.toLowerCase().includes(searchLower) ||
          session.id.toLowerCase().includes(searchLower);
        if (!matchesSearch) return false;
      }

      // Status filter
      if (filters.status !== 'all' && session.status !== filters.status) {
        return false;
      }

      // Client name filter
      if (filters.clientName && session.client_name !== filters.clientName) {
        return false;
      }

      // Date range filter
      if (filters.dateRange !== 'all') {
        const sessionDate = new Date(session.created_at);
        const now = new Date();
        const diffDays = Math.floor((now.getTime() - sessionDate.getTime()) / (1000 * 60 * 60 * 24));

        switch (filters.dateRange) {
          case 'today':
            if (diffDays > 0) return false;
            break;
          case 'week':
            if (diffDays > 7) return false;
            break;
          case 'month':
            if (diffDays > 30) return false;
            break;
        }
      }

      return true;
    });

    // Sort sessions
    filtered.sort((a, b) => {
      let aValue: string | number;
      let bValue: string | number;

      switch (sortBy) {
        case 'created_at':
          aValue = new Date(a.created_at).getTime();
          bValue = new Date(b.created_at).getTime();
          break;
        case 'last_activity':
          aValue = new Date(a.last_activity).getTime();
          bValue = new Date(b.last_activity).getTime();
          break;
        case 'client_name':
          aValue = a.client_name || '';
          bValue = b.client_name || '';
          break;
        default:
          aValue = new Date(a.last_activity).getTime();
          bValue = new Date(b.last_activity).getTime();
      }

      if (sortOrder === 'asc') {
        return aValue < bValue ? -1 : aValue > bValue ? 1 : 0;
      } else {
        return aValue > bValue ? -1 : aValue < bValue ? 1 : 0;
      }
    });

    return filtered;
  }, [sessions, filters, sortBy, sortOrder]);

  // Get unique client names for filter dropdown
  const uniqueClients = useMemo(() => {
    const clients = sessions
      .map(s => s.client_name)
      .filter((name): name is string => Boolean(name))
      .filter((name, index, arr) => arr.indexOf(name) === index)
      .sort();
    return clients;
  }, [sessions]);

  const handleCreateSession = async () => {
    try {
      const sessionData: Partial<ChatSession> = {
        consultant_id: 'current-user', // This would come from auth context
        client_name: newSessionData.client_name,
        context: newSessionData.project_context,
        status: 'active',
      };

      await dispatch(createSession(sessionData)).unwrap();
      setShowCreateSession(false);
      setNewSessionData({
        client_name: '',
        meeting_type: '',
        project_context: '',
      });
    } catch (error) {
      console.error('Failed to create session:', error);
    }
  };

  const handleSessionSelect = (session: ChatSession) => {
    dispatch(setCurrentSession(session));
    onSessionSelect(session);
    onClose();
  };

  const handleSessionToggle = (sessionId: string) => {
    const newSelected = new Set(selectedSessions);
    if (newSelected.has(sessionId)) {
      newSelected.delete(sessionId);
    } else {
      newSelected.add(sessionId);
    }
    setSelectedSessions(newSelected);
  };

  const handleExpandToggle = (sessionId: string) => {
    const newExpanded = new Set(expandedSessions);
    if (newExpanded.has(sessionId)) {
      newExpanded.delete(sessionId);
    } else {
      newExpanded.add(sessionId);
    }
    setExpandedSessions(newExpanded);
  };

  const handleBulkExport = () => {
    const sessionsToExport = sessions.filter(s => selectedSessions.has(s.id));
    const exportData = {
      sessions: sessionsToExport,
      messageHistory: Object.fromEntries(
        Object.entries(messageHistory).filter(([sessionId]) => selectedSessions.has(sessionId))
      ),
      exportedAt: new Date().toISOString(),
      stats: sessionStats,
    };

    const blob = new Blob([JSON.stringify(exportData, null, 2)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `chat-sessions-export-${new Date().toISOString().split('T')[0]}.json`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  const formatDuration = (minutes: number) => {
    if (minutes < 60) return `${Math.round(minutes)}m`;
    const hours = Math.floor(minutes / 60);
    const mins = Math.round(minutes % 60);
    return `${hours}h ${mins}m`;
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

    if (diffDays === 0) return 'Today';
    if (diffDays === 1) return 'Yesterday';
    if (diffDays < 7) return `${diffDays} days ago`;
    
    return date.toLocaleDateString([], { 
      month: 'short', 
      day: 'numeric',
      ...(date.getFullYear() !== now.getFullYear() && { year: 'numeric' })
    });
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 overflow-hidden" role="dialog" aria-modal="true" aria-labelledby="session-manager-title">
      <div 
        className="absolute inset-0 bg-gray-600 bg-opacity-75 transition-opacity" 
        onClick={onClose}
        aria-hidden="true"
      />
      
      <div className="absolute right-0 top-0 h-full w-full max-w-xs sm:max-w-2xl lg:max-w-4xl bg-white shadow-xl transform transition-transform"
           style={{ transform: 'translateX(0)' }}>
        <div className="flex h-full flex-col">
          {/* Header */}
          <header className="flex items-center justify-between p-4 sm:p-6 border-b border-gray-200 flex-shrink-0">
            <div className="flex items-center space-x-2 sm:space-x-3 min-w-0 flex-1">
              <History className="h-5 w-5 sm:h-6 sm:w-6 text-blue-600 flex-shrink-0" aria-hidden="true" />
              <h2 id="session-manager-title" className="text-lg sm:text-xl font-semibold text-gray-900 truncate">
                Session Manager
              </h2>
            </div>
            <button
              onClick={onClose}
              className="text-gray-400 hover:text-gray-600 transition-colors p-1 rounded focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
              aria-label="Close session manager"
            >
              <X className="h-5 w-5 sm:h-6 sm:w-6" />
            </button>
          </header>

          {/* Stats Overview */}
          <div className="p-6 border-b border-gray-200 bg-gray-50">
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              <div className="bg-white p-4 rounded-lg border">
                <div className="flex items-center space-x-2">
                  <MessageSquare className="h-5 w-5 text-blue-600" />
                  <span className="text-sm font-medium text-gray-600">Total Sessions</span>
                </div>
                <p className="text-2xl font-bold text-gray-900 mt-1">{sessionStats.totalSessions}</p>
              </div>
              <div className="bg-white p-4 rounded-lg border">
                <div className="flex items-center space-x-2">
                  <Clock className="h-5 w-5 text-green-600" />
                  <span className="text-sm font-medium text-gray-600">Active</span>
                </div>
                <p className="text-2xl font-bold text-gray-900 mt-1">{sessionStats.activeSessions}</p>
              </div>
              <div className="bg-white p-4 rounded-lg border">
                <div className="flex items-center space-x-2">
                  <BarChart3 className="h-5 w-5 text-purple-600" />
                  <span className="text-sm font-medium text-gray-600">Messages</span>
                </div>
                <p className="text-2xl font-bold text-gray-900 mt-1">{sessionStats.totalMessages}</p>
              </div>
              <div className="bg-white p-4 rounded-lg border">
                <div className="flex items-center space-x-2">
                  <Calendar className="h-5 w-5 text-orange-600" />
                  <span className="text-sm font-medium text-gray-600">Avg Duration</span>
                </div>
                <p className="text-2xl font-bold text-gray-900 mt-1">
                  {formatDuration(sessionStats.avgSessionDuration)}
                </p>
              </div>
            </div>
          </div>

          {/* Controls */}
          <div className="p-6 border-b border-gray-200 space-y-4">
            {/* Search and Filters */}
            <div className="flex flex-wrap items-center gap-4">
              <div className="flex-1 min-w-64">
                <div className="relative">
                  <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                  <input
                    type="text"
                    value={filters.search}
                    onChange={(e) => setFilters(prev => ({ ...prev, search: e.target.value }))}
                    placeholder="Search sessions..."
                    className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
              </div>
              
              <select
                value={filters.status}
                onChange={(e) => setFilters(prev => ({ ...prev, status: e.target.value as any }))}
                className="border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="all">All Status</option>
                <option value="active">Active</option>
                <option value="inactive">Inactive</option>
                <option value="expired">Expired</option>
              </select>

              <select
                value={filters.dateRange}
                onChange={(e) => setFilters(prev => ({ ...prev, dateRange: e.target.value as any }))}
                className="border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="all">All Time</option>
                <option value="today">Today</option>
                <option value="week">This Week</option>
                <option value="month">This Month</option>
              </select>

              <select
                value={filters.clientName}
                onChange={(e) => setFilters(prev => ({ ...prev, clientName: e.target.value }))}
                className="border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="">All Clients</option>
                {uniqueClients.map(client => (
                  <option key={client} value={client}>{client}</option>
                ))}
              </select>
            </div>

            {/* Action Buttons */}
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-2">
                <button
                  onClick={() => setShowCreateSession(true)}
                  className="flex items-center space-x-2 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
                >
                  <Plus className="h-4 w-4" />
                  <span>New Session</span>
                </button>
                
                {selectedSessions.size > 0 && (
                  <>
                    <button
                      onClick={handleBulkExport}
                      className="flex items-center space-x-2 px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition-colors"
                    >
                      <Download className="h-4 w-4" />
                      <span>Export ({selectedSessions.size})</span>
                    </button>
                    <button
                      onClick={() => setSelectedSessions(new Set())}
                      className="px-4 py-2 text-gray-600 hover:text-gray-800 transition-colors"
                    >
                      Clear Selection
                    </button>
                  </>
                )}
              </div>

              <div className="flex items-center space-x-2">
                <span className="text-sm text-gray-600">Sort by:</span>
                <select
                  value={sortBy}
                  onChange={(e) => setSortBy(e.target.value as any)}
                  className="border border-gray-300 rounded px-2 py-1 text-sm focus:outline-none focus:ring-1 focus:ring-blue-500"
                >
                  <option value="last_activity">Last Activity</option>
                  <option value="created_at">Created</option>
                  <option value="client_name">Client Name</option>
                </select>
                <button
                  onClick={() => setSortOrder(prev => prev === 'asc' ? 'desc' : 'asc')}
                  className="p-1 text-gray-500 hover:text-gray-700 transition-colors"
                >
                  {sortOrder === 'asc' ? '↑' : '↓'}
                </button>
              </div>
            </div>
          </div>

          {/* Sessions List */}
          <div className="flex-1 overflow-y-auto">
            {filteredSessions.length === 0 ? (
              <div className="text-center py-12">
                <History className="h-12 w-12 mx-auto text-gray-400 mb-4" />
                <h3 className="text-lg font-medium text-gray-900 mb-2">No sessions found</h3>
                <p className="text-gray-600">
                  {filters.search || filters.status !== 'all' || filters.dateRange !== 'all' || filters.clientName
                    ? 'Try adjusting your filters'
                    : 'Create your first session to get started'
                  }
                </p>
              </div>
            ) : (
              <div className="divide-y divide-gray-200">
                {filteredSessions.map((session) => {
                  const isExpanded = expandedSessions.has(session.id);
                  const isSelected = selectedSessions.has(session.id);
                  const isCurrent = currentSession?.id === session.id;
                  const messageCount = messageHistory[session.id]?.length || 0;

                  return (
                    <div
                      key={session.id}
                      className={`p-4 hover:bg-gray-50 transition-colors ${
                        isCurrent ? 'bg-blue-50 border-l-4 border-blue-500' : ''
                      }`}
                    >
                      <div className="flex items-center space-x-3">
                        <input
                          type="checkbox"
                          checked={isSelected}
                          onChange={() => handleSessionToggle(session.id)}
                          className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                        />
                        
                        <button
                          onClick={() => handleExpandToggle(session.id)}
                          className="text-gray-400 hover:text-gray-600 transition-colors"
                        >
                          {isExpanded ? (
                            <ChevronDown className="h-4 w-4" />
                          ) : (
                            <ChevronRight className="h-4 w-4" />
                          )}
                        </button>

                        <div className="flex-1 min-w-0">
                          <div className="flex items-center justify-between">
                            <div className="flex items-center space-x-3">
                              <h3 className="text-sm font-medium text-gray-900 truncate">
                                {session.client_name || `Session ${session.id.slice(-8)}`}
                              </h3>
                              <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                                session.status === 'active' 
                                  ? 'bg-green-100 text-green-800'
                                  : session.status === 'inactive'
                                  ? 'bg-gray-100 text-gray-800'
                                  : 'bg-red-100 text-red-800'
                              }`}>
                                {session.status}
                              </span>
                              {isCurrent && (
                                <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                                  Current
                                </span>
                              )}
                            </div>
                            
                            <div className="flex items-center space-x-2">
                              <span className="text-xs text-gray-500">
                                {messageCount} messages
                              </span>
                              <span className="text-xs text-gray-500">
                                {formatDate(session.last_activity)}
                              </span>
                              <button
                                onClick={() => handleSessionSelect(session)}
                                className="text-blue-600 hover:text-blue-800 text-sm font-medium transition-colors"
                              >
                                Open
                              </button>
                            </div>
                          </div>
                          
                          {session.context && (
                            <p className="text-sm text-gray-600 mt-1 truncate">
                              {session.context}
                            </p>
                          )}
                        </div>
                      </div>

                      {/* Expanded Details */}
                      {isExpanded && (
                        <div className="mt-4 ml-10 p-4 bg-gray-50 rounded-lg">
                          <div className="grid grid-cols-2 gap-4 text-sm">
                            <div>
                              <span className="font-medium text-gray-700">Session ID:</span>
                              <p className="text-gray-600 font-mono">{session.id}</p>
                            </div>
                            <div>
                              <span className="font-medium text-gray-700">Created:</span>
                              <p className="text-gray-600">{new Date(session.created_at).toLocaleString()}</p>
                            </div>
                            <div>
                              <span className="font-medium text-gray-700">Last Activity:</span>
                              <p className="text-gray-600">{new Date(session.last_activity).toLocaleString()}</p>
                            </div>
                            <div>
                              <span className="font-medium text-gray-700">Messages:</span>
                              <p className="text-gray-600">{messageCount}</p>
                            </div>
                          </div>
                          
                          <div className="flex items-center space-x-2 mt-4">
                            <button
                              onClick={() => {
                                // Export single session
                                const exportData = {
                                  session,
                                  messages: messageHistory[session.id] || [],
                                  exportedAt: new Date().toISOString(),
                                };
                                const blob = new Blob([JSON.stringify(exportData, null, 2)], { type: 'application/json' });
                                const url = URL.createObjectURL(blob);
                                const a = document.createElement('a');
                                a.href = url;
                                a.download = `session-${session.id}-export.json`;
                                document.body.appendChild(a);
                                a.click();
                                document.body.removeChild(a);
                                URL.revokeObjectURL(url);
                              }}
                              className="flex items-center space-x-1 px-3 py-1 text-sm text-gray-600 hover:text-gray-800 hover:bg-gray-200 rounded transition-colors"
                            >
                              <Download className="h-3 w-3" />
                              <span>Export</span>
                            </button>
                            <button className="flex items-center space-x-1 px-3 py-1 text-sm text-gray-600 hover:text-gray-800 hover:bg-gray-200 rounded transition-colors">
                              <Share2 className="h-3 w-3" />
                              <span>Share</span>
                            </button>
                            <button className="flex items-center space-x-1 px-3 py-1 text-sm text-red-600 hover:text-red-800 hover:bg-red-100 rounded transition-colors">
                              <Archive className="h-3 w-3" />
                              <span>Archive</span>
                            </button>
                          </div>
                        </div>
                      )}
                    </div>
                  );
                })}
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Create Session Modal */}
      {showCreateSession && (
        <div className="fixed inset-0 z-60 overflow-y-auto">
          <div className="flex items-center justify-center min-h-screen px-4">
            <div className="fixed inset-0 bg-gray-600 bg-opacity-75" onClick={() => setShowCreateSession(false)} />
            
            <div className="relative bg-white rounded-lg shadow-xl max-w-md w-full p-6">
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-medium text-gray-900">Create New Session</h3>
                <button
                  onClick={() => setShowCreateSession(false)}
                  className="text-gray-400 hover:text-gray-600 transition-colors"
                >
                  <X className="h-5 w-5" />
                </button>
              </div>

              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Client Name
                  </label>
                  <input
                    type="text"
                    value={newSessionData.client_name || ''}
                    onChange={(e) => setNewSessionData(prev => ({ ...prev, client_name: e.target.value }))}
                    placeholder="Enter client name..."
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Meeting Type
                  </label>
                  <input
                    type="text"
                    value={newSessionData.meeting_type || ''}
                    onChange={(e) => setNewSessionData(prev => ({ ...prev, meeting_type: e.target.value }))}
                    placeholder="e.g., Discovery, Architecture Review..."
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Project Context
                  </label>
                  <textarea
                    value={newSessionData.project_context || ''}
                    onChange={(e) => setNewSessionData(prev => ({ ...prev, project_context: e.target.value }))}
                    placeholder="Brief description of the project or meeting context..."
                    rows={3}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
              </div>

              <div className="flex items-center justify-end space-x-3 mt-6">
                <button
                  onClick={() => setShowCreateSession(false)}
                  className="px-4 py-2 text-gray-700 hover:text-gray-900 transition-colors"
                >
                  Cancel
                </button>
                <button
                  onClick={handleCreateSession}
                  className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
                >
                  Create Session
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default ChatSessionManager;