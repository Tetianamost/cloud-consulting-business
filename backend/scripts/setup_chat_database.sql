-- Standalone chat database setup script
-- This script can be run independently to set up chat tables

-- Connect to the consulting database
\c consulting;

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Drop existing types if they exist (for clean reinstall)
DROP TYPE IF EXISTS session_status CASCADE;
DROP TYPE IF EXISTS message_type CASCADE;
DROP TYPE IF EXISTS message_status CASCADE;

-- Create enum types for better type safety
CREATE TYPE session_status AS ENUM ('active', 'inactive', 'expired', 'terminated');
CREATE TYPE message_type AS ENUM ('user', 'assistant', 'system', 'error');
CREATE TYPE message_status AS ENUM ('sent', 'delivered', 'read', 'failed');

-- Drop existing tables if they exist (for clean reinstall)
DROP TABLE IF EXISTS chat_messages CASCADE;
DROP TABLE IF EXISTS chat_sessions CASCADE;

-- Create chat_sessions table
CREATE TABLE chat_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR(100) NOT NULL,
    client_name VARCHAR(255),
    context TEXT,
    status session_status DEFAULT 'active',
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_activity TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE,
    
    -- Constraints
    CONSTRAINT chk_user_id_not_empty CHECK (user_id != ''),
    CONSTRAINT chk_client_name_length CHECK (LENGTH(client_name) <= 255),
    CONSTRAINT chk_user_id_length CHECK (LENGTH(user_id) <= 100)
);

-- Create chat_messages table
CREATE TABLE chat_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id UUID NOT NULL,
    type message_type NOT NULL DEFAULT 'user',
    content TEXT NOT NULL,
    metadata JSONB DEFAULT '{}',
    status message_status DEFAULT 'sent',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT chk_content_not_empty CHECK (content != ''),
    CONSTRAINT chk_content_length CHECK (LENGTH(content) <= 10000),
    
    -- Foreign key constraint
    CONSTRAINT fk_chat_messages_session_id 
        FOREIGN KEY (session_id) REFERENCES chat_sessions(id) ON DELETE CASCADE
);

-- Create indexes for optimal query performance
CREATE INDEX idx_chat_sessions_user_id ON chat_sessions(user_id);
CREATE INDEX idx_chat_sessions_status ON chat_sessions(status);
CREATE INDEX idx_chat_sessions_created_at ON chat_sessions(created_at);
CREATE INDEX idx_chat_sessions_last_activity ON chat_sessions(last_activity);
CREATE INDEX idx_chat_sessions_expires_at ON chat_sessions(expires_at);

CREATE INDEX idx_chat_messages_session_id ON chat_messages(session_id);
CREATE INDEX idx_chat_messages_type ON chat_messages(type);
CREATE INDEX idx_chat_messages_status ON chat_messages(status);
CREATE INDEX idx_chat_messages_created_at ON chat_messages(created_at);

-- Composite indexes for common query patterns
CREATE INDEX idx_chat_sessions_user_status ON chat_sessions(user_id, status);
CREATE INDEX idx_chat_messages_session_created ON chat_messages(session_id, created_at);

-- Create trigger functions
CREATE OR REPLACE FUNCTION update_chat_session_timestamps()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    NEW.last_activity = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE OR REPLACE FUNCTION update_session_activity_on_message()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE chat_sessions 
    SET last_activity = NOW(), updated_at = NOW()
    WHERE id = NEW.session_id;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers
CREATE TRIGGER update_chat_sessions_timestamps 
    BEFORE UPDATE ON chat_sessions 
    FOR EACH ROW EXECUTE FUNCTION update_chat_session_timestamps();

CREATE TRIGGER update_session_on_message 
    AFTER INSERT ON chat_messages 
    FOR EACH ROW EXECUTE FUNCTION update_session_activity_on_message();

-- Create utility function to cleanup expired sessions
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

-- Create views for common queries
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

CREATE OR REPLACE VIEW chat_session_stats AS
SELECT 
    status,
    COUNT(*) as session_count,
    AVG(EXTRACT(EPOCH FROM (NOW() - created_at))/3600) as avg_age_hours,
    AVG(EXTRACT(EPOCH FROM (last_activity - created_at))/60) as avg_duration_minutes
FROM chat_sessions 
GROUP BY status;

-- Add table comments for documentation
COMMENT ON TABLE chat_sessions IS 'Stores chat session information for AI consultant live chat';
COMMENT ON TABLE chat_messages IS 'Stores individual chat messages within sessions';

COMMENT ON COLUMN chat_sessions.user_id IS 'Identifier for the user (admin) who owns the session';
COMMENT ON COLUMN chat_sessions.client_name IS 'Optional client name for context';
COMMENT ON COLUMN chat_sessions.context IS 'Additional context information for the chat session';
COMMENT ON COLUMN chat_sessions.metadata IS 'JSON metadata for extensibility';
COMMENT ON COLUMN chat_sessions.expires_at IS 'When the session expires (NULL for no expiration)';

COMMENT ON COLUMN chat_messages.type IS 'Type of message: user, assistant, system, or error';
COMMENT ON COLUMN chat_messages.content IS 'The actual message content';
COMMENT ON COLUMN chat_messages.metadata IS 'JSON metadata for message-specific information';

-- Insert sample data for testing (optional)
INSERT INTO chat_sessions (user_id, client_name, context, status) VALUES
    ('admin1', 'Acme Corp', 'Initial consultation about AWS migration', 'active'),
    ('admin2', 'Tech Startup', 'Cost optimization discussion', 'active'),
    ('admin1', 'Enterprise Ltd', 'Architecture review session', 'inactive')
ON CONFLICT DO NOTHING;

-- Get the session IDs for sample messages
DO $$
DECLARE
    session1_id UUID;
    session2_id UUID;
BEGIN
    SELECT id INTO session1_id FROM chat_sessions WHERE client_name = 'Acme Corp' LIMIT 1;
    SELECT id INTO session2_id FROM chat_sessions WHERE client_name = 'Tech Startup' LIMIT 1;
    
    IF session1_id IS NOT NULL THEN
        INSERT INTO chat_messages (session_id, type, content) VALUES
            (session1_id, 'user', 'Hello, I need help with AWS migration planning'),
            (session1_id, 'assistant', 'I''d be happy to help you with AWS migration planning. Can you tell me about your current infrastructure?'),
            (session1_id, 'user', 'We have about 50 VMs running various applications');
    END IF;
    
    IF session2_id IS NOT NULL THEN
        INSERT INTO chat_messages (session_id, type, content) VALUES
            (session2_id, 'user', 'Our AWS costs are getting too high'),
            (session2_id, 'assistant', 'Let''s analyze your cost structure. What services are you currently using?');
    END IF;
END $$;

-- Success message
SELECT 'Chat database setup completed successfully!' AS setup_result,
       (SELECT COUNT(*) FROM chat_sessions) AS sessions_created,
       (SELECT COUNT(*) FROM chat_messages) AS messages_created;