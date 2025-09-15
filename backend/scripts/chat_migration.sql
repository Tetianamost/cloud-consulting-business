-- Chat system database migration
-- This migration adds chat_sessions and chat_messages tables for the AI consultant live chat feature

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create enum types for better type safety
CREATE TYPE session_status AS ENUM ('active', 'inactive', 'expired', 'terminated');
CREATE TYPE message_type AS ENUM ('user', 'assistant', 'system', 'error');
CREATE TYPE message_status AS ENUM ('sent', 'delivered', 'read', 'failed');

-- Create chat_sessions table
CREATE TABLE IF NOT EXISTS chat_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR(100) NOT NULL,
    client_name VARCHAR(255),
    context TEXT,
    status session_status DEFAULT 'active',
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_activity TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE
);

-- Create chat_messages table
CREATE TABLE IF NOT EXISTS chat_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id UUID NOT NULL REFERENCES chat_sessions(id) ON DELETE CASCADE,
    type message_type NOT NULL DEFAULT 'user',
    content TEXT NOT NULL,
    metadata JSONB DEFAULT '{}',
    status message_status DEFAULT 'sent',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for optimal query performance
CREATE INDEX IF NOT EXISTS idx_chat_sessions_user_id ON chat_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_chat_sessions_status ON chat_sessions(status);
CREATE INDEX IF NOT EXISTS idx_chat_sessions_created_at ON chat_sessions(created_at);
CREATE INDEX IF NOT EXISTS idx_chat_sessions_last_activity ON chat_sessions(last_activity);
CREATE INDEX IF NOT EXISTS idx_chat_sessions_expires_at ON chat_sessions(expires_at);

CREATE INDEX IF NOT EXISTS idx_chat_messages_session_id ON chat_messages(session_id);
CREATE INDEX IF NOT EXISTS idx_chat_messages_type ON chat_messages(type);
CREATE INDEX IF NOT EXISTS idx_chat_messages_status ON chat_messages(status);
CREATE INDEX IF NOT EXISTS idx_chat_messages_created_at ON chat_messages(created_at);

-- Composite indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_chat_sessions_user_status ON chat_sessions(user_id, status);
CREATE INDEX IF NOT EXISTS idx_chat_messages_session_created ON chat_messages(session_id, created_at);

-- Create trigger to automatically update updated_at and last_activity for sessions
CREATE OR REPLACE FUNCTION update_chat_session_timestamps()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    NEW.last_activity = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_chat_sessions_timestamps 
    BEFORE UPDATE ON chat_sessions 
    FOR EACH ROW EXECUTE FUNCTION update_chat_session_timestamps();

-- Create trigger to update session last_activity when messages are added
CREATE OR REPLACE FUNCTION update_session_activity_on_message()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE chat_sessions 
    SET last_activity = NOW(), updated_at = NOW()
    WHERE id = NEW.session_id;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_session_on_message 
    AFTER INSERT ON chat_messages 
    FOR EACH ROW EXECUTE FUNCTION update_session_activity_on_message();

-- Add foreign key constraints with proper referential integrity
ALTER TABLE chat_messages 
ADD CONSTRAINT fk_chat_messages_session_id 
FOREIGN KEY (session_id) REFERENCES chat_sessions(id) ON DELETE CASCADE;

-- Create function to cleanup expired sessions
CREATE OR REPLACE FUNCTION cleanup_expired_chat_sessions()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM chat_sessions 
    WHERE expires_at IS NOT NULL AND expires_at < NOW();
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ language 'plpgsql';

-- Create view for active chat sessions with message counts
CREATE OR REPLACE VIEW active_chat_sessions AS
SELECT 
    cs.id,
    cs.user_id,
    cs.client_name,
    cs.context,
    cs.status,
    cs.created_at,
    cs.last_activity,
    cs.expires_at,
    COUNT(cm.id) as message_count,
    MAX(cm.created_at) as last_message_at
FROM chat_sessions cs
LEFT JOIN chat_messages cm ON cs.id = cm.session_id
WHERE cs.status = 'active'
GROUP BY cs.id, cs.user_id, cs.client_name, cs.context, cs.status, cs.created_at, cs.last_activity, cs.expires_at;

-- Create view for chat session statistics
CREATE OR REPLACE VIEW chat_session_stats AS
SELECT 
    status,
    COUNT(*) as session_count,
    AVG(EXTRACT(EPOCH FROM (NOW() - created_at))/3600) as avg_age_hours,
    AVG(EXTRACT(EPOCH FROM (last_activity - created_at))/60) as avg_duration_minutes
FROM chat_sessions 
GROUP BY status;

-- Add comments for documentation
COMMENT ON TABLE chat_sessions IS 'Stores chat session information for AI consultant live chat';
COMMENT ON TABLE chat_messages IS 'Stores individual chat messages within sessions';
COMMENT ON COLUMN chat_sessions.user_id IS 'Identifier for the user (admin) who owns the session';
COMMENT ON COLUMN chat_sessions.client_name IS 'Optional client name for context';
COMMENT ON COLUMN chat_sessions.context IS 'Additional context information for the chat session';
COMMENT ON COLUMN chat_sessions.metadata IS 'JSON metadata for extensibility';
COMMENT ON COLUMN chat_messages.type IS 'Type of message: user, assistant, system, or error';
COMMENT ON COLUMN chat_messages.content IS 'The actual message content';
COMMENT ON COLUMN chat_messages.metadata IS 'JSON metadata for message-specific information';