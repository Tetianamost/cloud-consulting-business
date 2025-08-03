-- Chat data seeding script for testing and development
-- This script creates sample chat sessions and messages for testing the enhanced repositories

-- Clear existing test data (optional - uncomment if needed)
-- DELETE FROM chat_messages WHERE session_id IN (SELECT id FROM chat_sessions WHERE user_id LIKE 'test_%');
-- DELETE FROM chat_sessions WHERE user_id LIKE 'test_%';

-- Insert test chat sessions with various scenarios
INSERT INTO chat_sessions (id, user_id, client_name, context, status, metadata, created_at, updated_at, last_activity, expires_at) VALUES
    -- Active sessions
    (uuid_generate_v4(), 'test_admin1', 'Acme Corporation', 'AWS migration planning session', 'active', 
     '{"meeting_type": "consultation", "priority": "high", "services": ["EC2", "RDS", "S3"]}',
     NOW() - INTERVAL '2 hours', NOW() - INTERVAL '5 minutes', NOW() - INTERVAL '5 minutes', NOW() + INTERVAL '2 hours'),
    
    (uuid_generate_v4(), 'test_admin2', 'TechStart Inc', 'Cost optimization review', 'active',
     '{"meeting_type": "review", "priority": "medium", "services": ["Lambda", "DynamoDB", "CloudWatch"]}',
     NOW() - INTERVAL '1 hour', NOW() - INTERVAL '10 minutes', NOW() - INTERVAL '10 minutes', NOW() + INTERVAL '3 hours'),
    
    (uuid_generate_v4(), 'test_admin1', 'Global Enterprises', 'Security audit consultation', 'active',
     '{"meeting_type": "audit", "priority": "high", "services": ["IAM", "CloudTrail", "GuardDuty"]}',
     NOW() - INTERVAL '30 minutes', NOW() - INTERVAL '2 minutes', NOW() - INTERVAL '2 minutes', NOW() + INTERVAL '4 hours'),
    
    -- Inactive sessions
    (uuid_generate_v4(), 'test_admin3', 'Small Business Co', 'Initial AWS setup', 'inactive',
     '{"meeting_type": "setup", "priority": "low", "services": ["EC2", "S3"]}',
     NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day', NULL),
    
    (uuid_generate_v4(), 'test_admin2', 'Mid-size Corp', 'Architecture review completed', 'inactive',
     '{"meeting_type": "review", "priority": "medium", "services": ["ECS", "ALB", "RDS"]}',
     NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days', NULL),
    
    -- Expired session
    (uuid_generate_v4(), 'test_admin1', 'Old Client Ltd', 'Expired consultation', 'expired',
     '{"meeting_type": "consultation", "priority": "low"}',
     NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days', NOW() - INTERVAL '1 day'),
    
    -- Terminated session
    (uuid_generate_v4(), 'test_admin3', 'Cancelled Project', 'Terminated early', 'terminated',
     '{"meeting_type": "consultation", "priority": "medium", "termination_reason": "client_request"}',
     NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days', NULL)
ON CONFLICT DO NOTHING;

-- Get session IDs for message insertion
DO $
DECLARE
    session_acme UUID;
    session_techstart UUID;
    session_global UUID;
    session_small UUID;
    session_midsize UUID;
    session_old UUID;
    session_cancelled UUID;
BEGIN
    -- Get session IDs
    SELECT id INTO session_acme FROM chat_sessions WHERE client_name = 'Acme Corporation' AND user_id = 'test_admin1' LIMIT 1;
    SELECT id INTO session_techstart FROM chat_sessions WHERE client_name = 'TechStart Inc' AND user_id = 'test_admin2' LIMIT 1;
    SELECT id INTO session_global FROM chat_sessions WHERE client_name = 'Global Enterprises' AND user_id = 'test_admin1' LIMIT 1;
    SELECT id INTO session_small FROM chat_sessions WHERE client_name = 'Small Business Co' AND user_id = 'test_admin3' LIMIT 1;
    SELECT id INTO session_midsize FROM chat_sessions WHERE client_name = 'Mid-size Corp' AND user_id = 'test_admin2' LIMIT 1;
    SELECT id INTO session_old FROM chat_sessions WHERE client_name = 'Old Client Ltd' AND user_id = 'test_admin1' LIMIT 1;
    SELECT id INTO session_cancelled FROM chat_sessions WHERE client_name = 'Cancelled Project' AND user_id = 'test_admin3' LIMIT 1;
    
    -- Insert messages for Acme Corporation session (active, high-volume)
    IF session_acme IS NOT NULL THEN
        INSERT INTO chat_messages (id, session_id, type, content, metadata, status, created_at) VALUES
            (uuid_generate_v4(), session_acme, 'user', 'Hello, I need help planning our AWS migration. We have about 50 VMs currently running on-premises.', 
             '{"user_name": "John Smith", "company": "Acme Corporation"}', 'delivered', NOW() - INTERVAL '2 hours'),
            
            (uuid_generate_v4(), session_acme, 'assistant', 'Hello John! I''d be happy to help you plan your AWS migration. With 50 VMs, we''ll want to do a thorough assessment. Can you tell me what types of applications you''re running?',
             '{"response_time_ms": 1200, "confidence": 0.95}', 'delivered', NOW() - INTERVAL '2 hours' + INTERVAL '30 seconds'),
            
            (uuid_generate_v4(), session_acme, 'user', 'We have a mix of web applications, databases, and some legacy systems. Most are Windows-based with a few Linux servers.',
             '{"user_name": "John Smith"}', 'delivered', NOW() - INTERVAL '2 hours' + INTERVAL '2 minutes'),
            
            (uuid_generate_v4(), session_acme, 'assistant', 'Great! For a mixed Windows/Linux environment, I recommend starting with AWS Application Discovery Service to inventory your systems. We can then use AWS Migration Hub to track progress. Would you like me to walk you through the migration strategies?',
             '{"response_time_ms": 1500, "confidence": 0.92, "suggested_services": ["Application Discovery Service", "Migration Hub"]}', 'delivered', NOW() - INTERVAL '2 hours' + INTERVAL '3 minutes'),
            
            (uuid_generate_v4(), session_acme, 'user', 'Yes, please explain the different migration strategies available.',
             '{"user_name": "John Smith"}', 'delivered', NOW() - INTERVAL '2 hours' + INTERVAL '5 minutes'),
            
            (uuid_generate_v4(), session_acme, 'assistant', 'There are 6 main migration strategies (the 6 Rs): 1) Rehost (lift-and-shift), 2) Replatform (lift-tinker-shift), 3) Repurchase (move to SaaS), 4) Refactor/Re-architect, 5) Retire, 6) Retain. For your 50 VMs, I''d recommend starting with Rehost for quick wins, then Replatform for optimization.',
             '{"response_time_ms": 2100, "confidence": 0.98, "migration_strategies": ["rehost", "replatform"], "recommended_first": "rehost"}', 'delivered', NOW() - INTERVAL '2 hours' + INTERVAL '6 minutes'),
            
            (uuid_generate_v4(), session_acme, 'user', 'What about costs? We need to understand the financial impact.',
             '{"user_name": "John Smith"}', 'delivered', NOW() - INTERVAL '1 hour'),
            
            (uuid_generate_v4(), session_acme, 'assistant', 'Excellent question! I recommend using AWS Pricing Calculator for estimates. Generally, you can expect 20-30% cost savings with right-sizing. We should also consider Reserved Instances for predictable workloads. Would you like me to help you estimate costs for your specific workloads?',
             '{"response_time_ms": 1800, "confidence": 0.94, "cost_savings_estimate": "20-30%", "tools": ["AWS Pricing Calculator"]}', 'delivered', NOW() - INTERVAL '1 hour' + INTERVAL '45 seconds'),
            
            (uuid_generate_v4(), session_acme, 'user', 'That would be very helpful. Can you also recommend a timeline for the migration?',
             '{"user_name": "John Smith"}', 'delivered', NOW() - INTERVAL '30 minutes'),
            
            (uuid_generate_v4(), session_acme, 'assistant', 'For 50 VMs, I''d suggest a phased approach over 6-9 months: Phase 1 (Months 1-2): Assessment and planning, Phase 2 (Months 3-5): Migrate non-critical systems, Phase 3 (Months 6-8): Migrate critical systems, Phase 4 (Month 9): Optimization and cleanup. This allows for learning and reduces risk.',
             '{"response_time_ms": 2200, "confidence": 0.96, "timeline_months": "6-9", "phases": 4}', 'delivered', NOW() - INTERVAL '30 minutes' + INTERVAL '1 minute'),
            
            (uuid_generate_v4(), session_acme, 'user', 'Perfect! Can you send me a summary of our discussion?',
             '{"user_name": "John Smith"}', 'read', NOW() - INTERVAL '5 minutes'),
            
            (uuid_generate_v4(), session_acme, 'system', 'Generating migration summary report...',
             '{"action": "generate_report", "report_type": "migration_summary"}', 'sent', NOW() - INTERVAL '4 minutes');
    END IF;
    
    -- Insert messages for TechStart Inc session (cost optimization focus)
    IF session_techstart IS NOT NULL THEN
        INSERT INTO chat_messages (id, session_id, type, content, metadata, status, created_at) VALUES
            (uuid_generate_v4(), session_techstart, 'user', 'Our AWS costs have increased significantly. We need help optimizing our spending.',
             '{"user_name": "Sarah Johnson", "company": "TechStart Inc"}', 'delivered', NOW() - INTERVAL '1 hour'),
            
            (uuid_generate_v4(), session_techstart, 'assistant', 'I understand your concern about rising costs. Let''s start by analyzing your current usage. Are you using AWS Cost Explorer and have you set up any budgets or alerts?',
             '{"response_time_ms": 1100, "confidence": 0.93}', 'delivered', NOW() - INTERVAL '1 hour' + INTERVAL '20 seconds'),
            
            (uuid_generate_v4(), session_techstart, 'user', 'We have Cost Explorer but haven''t set up detailed budgets. Our main costs seem to be from EC2 and data transfer.',
             '{"user_name": "Sarah Johnson"}', 'delivered', NOW() - INTERVAL '50 minutes'),
            
            (uuid_generate_v4(), session_techstart, 'assistant', 'Great starting point! For EC2 optimization: 1) Right-size instances using AWS Compute Optimizer, 2) Consider Reserved Instances for steady workloads, 3) Use Spot Instances for fault-tolerant workloads. For data transfer: review your architecture for unnecessary cross-AZ transfers.',
             '{"response_time_ms": 1600, "confidence": 0.95, "optimization_areas": ["right-sizing", "reserved_instances", "spot_instances", "data_transfer"]}', 'delivered', NOW() - INTERVAL '48 minutes'),
            
            (uuid_generate_v4(), session_techstart, 'user', 'We''re running several Lambda functions too. Any optimization tips there?',
             '{"user_name": "Sarah Johnson"}', 'delivered', NOW() - INTERVAL '30 minutes'),
            
            (uuid_generate_v4(), session_techstart, 'assistant', 'Absolutely! Lambda optimization tips: 1) Right-size memory allocation (CPU scales with memory), 2) Minimize cold starts with provisioned concurrency if needed, 3) Use ARM-based Graviton2 processors for up to 34% better price performance, 4) Optimize your code for faster execution.',
             '{"response_time_ms": 1400, "confidence": 0.97, "lambda_optimizations": ["memory_sizing", "provisioned_concurrency", "graviton2", "code_optimization"]}', 'delivered', NOW() - INTERVAL '29 minutes'),
            
            (uuid_generate_v4(), session_techstart, 'user', 'This is very helpful. Can you help us set up a cost monitoring strategy?',
             '{"user_name": "Sarah Johnson"}', 'read', NOW() - INTERVAL '10 minutes'),
            
            (uuid_generate_v4(), session_techstart, 'assistant', 'Certainly! I recommend: 1) Set up AWS Budgets with alerts at 50%, 80%, and 100% of your target, 2) Use Cost Anomaly Detection for unusual spending patterns, 3) Implement cost allocation tags for better tracking, 4) Schedule monthly cost reviews. Would you like me to walk you through setting these up?',
             '{"response_time_ms": 1700, "confidence": 0.94, "monitoring_strategy": ["budgets", "anomaly_detection", "cost_tags", "monthly_reviews"]}', 'sent', NOW() - INTERVAL '9 minutes');
    END IF;
    
    -- Insert messages for Global Enterprises session (security focus)
    IF session_global IS NOT NULL THEN
        INSERT INTO chat_messages (id, session_id, type, content, metadata, status, created_at) VALUES
            (uuid_generate_v4(), session_global, 'user', 'We need to conduct a security audit of our AWS environment. Where should we start?',
             '{"user_name": "Michael Chen", "company": "Global Enterprises", "role": "CISO"}', 'delivered', NOW() - INTERVAL '30 minutes'),
            
            (uuid_generate_v4(), session_global, 'assistant', 'Excellent initiative! For a comprehensive security audit, I recommend starting with AWS Security Hub for centralized security findings, then AWS Config for compliance monitoring, and AWS CloudTrail for activity logging. Have you enabled these services?',
             '{"response_time_ms": 1300, "confidence": 0.96, "security_services": ["Security Hub", "Config", "CloudTrail"]}', 'delivered', NOW() - INTERVAL '29 minutes'),
            
            (uuid_generate_v4(), session_global, 'user', 'We have CloudTrail enabled but not the others. What about IAM best practices?',
             '{"user_name": "Michael Chen"}', 'delivered', NOW() - INTERVAL '25 minutes'),
            
            (uuid_generate_v4(), session_global, 'assistant', 'Great question! Key IAM best practices: 1) Enable MFA for all users, 2) Use IAM roles instead of long-term access keys, 3) Apply least privilege principle, 4) Regular access reviews, 5) Use AWS IAM Access Analyzer to identify unused permissions. I can help you audit your current IAM setup.',
             '{"response_time_ms": 1500, "confidence": 0.98, "iam_best_practices": ["mfa", "roles_over_keys", "least_privilege", "access_reviews", "access_analyzer"]}', 'delivered', NOW() - INTERVAL '24 minutes'),
            
            (uuid_generate_v4(), session_global, 'user', 'We also need to ensure our data is properly encrypted. Any recommendations?',
             '{"user_name": "Michael Chen"}', 'delivered', NOW() - INTERVAL '15 minutes'),
            
            (uuid_generate_v4(), session_global, 'assistant', 'Absolutely! Encryption recommendations: 1) Enable encryption at rest for all storage services (S3, EBS, RDS), 2) Use AWS KMS for key management, 3) Enable encryption in transit with TLS/SSL, 4) Consider AWS CloudHSM for high-security requirements. AWS Macie can help discover and classify sensitive data.',
             '{"response_time_ms": 1600, "confidence": 0.97, "encryption_services": ["KMS", "CloudHSM", "Macie"]}', 'delivered', NOW() - INTERVAL '14 minutes'),
            
            (uuid_generate_v4(), session_global, 'user', 'Perfect. Can you create a security audit checklist for us?',
             '{"user_name": "Michael Chen"}', 'read', NOW() - INTERVAL '2 minutes'),
            
            (uuid_generate_v4(), session_global, 'system', 'Generating comprehensive security audit checklist...',
             '{"action": "generate_checklist", "checklist_type": "security_audit"}', 'sent', NOW() - INTERVAL '1 minute');
    END IF;
    
    -- Insert messages for completed sessions
    IF session_small IS NOT NULL THEN
        INSERT INTO chat_messages (id, session_id, type, content, metadata, status, created_at) VALUES
            (uuid_generate_v4(), session_small, 'user', 'We''re a small business looking to move to AWS. What''s the best starting point?',
             '{"user_name": "Lisa Wong", "company": "Small Business Co"}', 'delivered', NOW() - INTERVAL '1 day'),
            
            (uuid_generate_v4(), session_small, 'assistant', 'Welcome to AWS! For small businesses, I recommend starting with AWS Lightsail for simple workloads, or EC2 with S3 for more flexibility. What type of applications are you looking to host?',
             '{"response_time_ms": 1000, "confidence": 0.90}', 'delivered', NOW() - INTERVAL '1 day' + INTERVAL '30 seconds'),
            
            (uuid_generate_v4(), session_small, 'user', 'Just a simple website and email hosting.',
             '{"user_name": "Lisa Wong"}', 'delivered', NOW() - INTERVAL '1 day' + INTERVAL '2 minutes'),
            
            (uuid_generate_v4(), session_small, 'assistant', 'Perfect! AWS Lightsail would be ideal for you. It includes everything you need: virtual server, SSD storage, data transfer, DNS management, and static IP. For email, consider Amazon WorkMail or integrate with a third-party service.',
             '{"response_time_ms": 1200, "confidence": 0.95, "recommended_services": ["Lightsail", "WorkMail"]}', 'delivered', NOW() - INTERVAL '1 day' + INTERVAL '3 minutes');
    END IF;
    
    IF session_midsize IS NOT NULL THEN
        INSERT INTO chat_messages (id, session_id, type, content, metadata, status, created_at) VALUES
            (uuid_generate_v4(), session_midsize, 'user', 'We completed our architecture review. Thank you for all the recommendations!',
             '{"user_name": "David Park", "company": "Mid-size Corp"}', 'delivered', NOW() - INTERVAL '2 days'),
            
            (uuid_generate_v4(), session_midsize, 'assistant', 'You''re very welcome! I''m glad we could help optimize your architecture. Don''t hesitate to reach out if you need assistance during implementation. Good luck with your deployment!',
             '{"response_time_ms": 800, "confidence": 0.99, "session_outcome": "successful"}', 'delivered', NOW() - INTERVAL '2 days' + INTERVAL '1 minute'),
            
            (uuid_generate_v4(), session_midsize, 'system', 'Session marked as completed successfully.',
             '{"action": "session_complete", "outcome": "successful", "satisfaction": "high"}', 'delivered', NOW() - INTERVAL '2 days' + INTERVAL '2 minutes');
    END IF;
    
    -- Insert messages for expired session
    IF session_old IS NOT NULL THEN
        INSERT INTO chat_messages (id, session_id, type, content, metadata, status, created_at) VALUES
            (uuid_generate_v4(), session_old, 'user', 'Hello, we need some basic AWS consultation.',
             '{"user_name": "Robert Taylor", "company": "Old Client Ltd"}', 'delivered', NOW() - INTERVAL '3 days'),
            
            (uuid_generate_v4(), session_old, 'assistant', 'Hello Robert! I''d be happy to help with your AWS consultation. What specific areas are you interested in?',
             '{"response_time_ms": 900, "confidence": 0.85}', 'delivered', NOW() - INTERVAL '3 days' + INTERVAL '20 seconds'),
            
            (uuid_generate_v4(), session_old, 'system', 'Session expired due to inactivity.',
             '{"action": "session_expired", "reason": "inactivity", "last_activity": "3 days ago"}', 'delivered', NOW() - INTERVAL '1 day');
    END IF;
    
    -- Insert messages for terminated session
    IF session_cancelled IS NOT NULL THEN
        INSERT INTO chat_messages (id, session_id, type, content, metadata, status, created_at) VALUES
            (uuid_generate_v4(), session_cancelled, 'user', 'We need to discuss our cloud strategy.',
             '{"user_name": "Jennifer Adams", "company": "Cancelled Project"}', 'delivered', NOW() - INTERVAL '5 days'),
            
            (uuid_generate_v4(), session_cancelled, 'assistant', 'I''d be happy to help with your cloud strategy. What are your main objectives?',
             '{"response_time_ms": 1000, "confidence": 0.88}', 'delivered', NOW() - INTERVAL '5 days' + INTERVAL '15 seconds'),
            
            (uuid_generate_v4(), session_cancelled, 'user', 'Actually, we need to postpone this. The project has been cancelled.',
             '{"user_name": "Jennifer Adams"}', 'delivered', NOW() - INTERVAL '5 days' + INTERVAL '5 minutes'),
            
            (uuid_generate_v4(), session_cancelled, 'system', 'Session terminated at client request.',
             '{"action": "session_terminated", "reason": "client_request", "termination_type": "project_cancelled"}', 'delivered', NOW() - INTERVAL '5 days' + INTERVAL '6 minutes');
    END IF;
    
END $;

-- Update session statistics after seeding
SELECT 'Chat data seeding completed successfully!' AS seeding_result,
       (SELECT COUNT(*) FROM chat_sessions WHERE user_id LIKE 'test_%') AS test_sessions_created,
       (SELECT COUNT(*) FROM chat_messages WHERE session_id IN (SELECT id FROM chat_sessions WHERE user_id LIKE 'test_%')) AS test_messages_created;

-- Display summary of seeded data
SELECT 
    'Data Summary' as category,
    status,
    COUNT(*) as session_count,
    (SELECT COUNT(*) FROM chat_messages WHERE session_id IN (SELECT id FROM chat_sessions WHERE status = cs.status AND user_id LIKE 'test_%')) as message_count
FROM chat_sessions cs 
WHERE user_id LIKE 'test_%'
GROUP BY status
ORDER BY 
    CASE status 
        WHEN 'active' THEN 1 
        WHEN 'inactive' THEN 2 
        WHEN 'expired' THEN 3 
        WHEN 'terminated' THEN 4 
    END;