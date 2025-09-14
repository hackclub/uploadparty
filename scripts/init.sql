-- PostgreSQL initialization script for Upload Party

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- Create database user with limited privileges (if not exists)
DO $$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_catalog.pg_user WHERE usename = 'uploadparty_app') THEN
      CREATE USER uploadparty_app WITH PASSWORD 'secure_app_password';
   END IF;
END
$$;

-- Grant necessary privileges
GRANT CONNECT ON DATABASE uploadparty_db TO uploadparty_app;
GRANT USAGE ON SCHEMA public TO uploadparty_app;
GRANT CREATE ON SCHEMA public TO uploadparty_app;

-- Indexes for performance
-- These will be created after GORM migrations, but we can prepare them

-- Create a function to add indexes after tables are created
CREATE OR REPLACE FUNCTION create_performance_indexes()
RETURNS void AS $$
BEGIN
    -- Only create indexes if tables exist
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users') THEN
        -- User indexes
        CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_username ON users(username);
        CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_email ON users(email);
        CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_total_points ON users(total_points DESC);
        CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_created_at ON users(created_at);
        CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_active ON users(is_active) WHERE is_active = true;
    END IF;

    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'beats') THEN
        -- Beat indexes
        CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_beats_user_id ON beats(user_id);
        CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_beats_genre ON beats(genre);
        CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_beats_created_at ON beats(created_at DESC);
        CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_beats_play_count ON beats(play_count DESC);
        CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_beats_like_count ON beats(like_count DESC);
        CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_beats_public_approved ON beats(is_public, is_approved) WHERE is_public = true AND is_approved = true;
        
        -- Full-text search index for beat titles and descriptions
        CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_beats_search ON beats USING gin(to_tsvector('english', title || ' ' || description));
    END IF;

    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'challenges') THEN
        -- Challenge indexes
        CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_challenges_active ON challenges(is_active) WHERE is_active = true;
        CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_challenges_dates ON challenges(start_date, end_date);
        CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_challenges_status ON challenges(status);
    END IF;

    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'challenge_submissions') THEN
        -- Challenge submission indexes
        CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_submissions_challenge_id ON challenge_submissions(challenge_id);
        CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_submissions_user_id ON challenge_submissions(user_id);
        CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_submissions_total_score ON challenge_submissions(total_score DESC);
    END IF;

    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'likes') THEN
        -- Like indexes
        CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_likes_user_beat ON likes(user_id, beat_id);
        CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_likes_beat_id ON likes(beat_id);
    END IF;

    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'comments') THEN
        -- Comment indexes
        CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_comments_beat_id ON comments(beat_id);
        CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_comments_user_id ON comments(user_id);
        CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_comments_created_at ON comments(created_at DESC);
    END IF;

    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'scores') THEN
        -- Score indexes
        CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_scores_user_id ON scores(user_id);
        CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_scores_source ON scores(source);
        CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_scores_created_at ON scores(created_at DESC);
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Security: Row Level Security policies (can be added after tables are created)
-- This is a placeholder for future RLS implementation

-- Create a function to enable RLS (to be called after migrations)
CREATE OR REPLACE FUNCTION enable_row_level_security()
RETURNS void AS $$
BEGIN
    -- Enable RLS on sensitive tables
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users') THEN
        ALTER TABLE users ENABLE ROW LEVEL SECURITY;
        
        -- Policy: Users can only see their own data or public profiles
        DROP POLICY IF EXISTS users_select_policy ON users;
        CREATE POLICY users_select_policy ON users
            FOR SELECT
            USING (id = current_setting('app.current_user_id')::uuid OR is_active = true);
            
        -- Policy: Users can only update their own data
        DROP POLICY IF EXISTS users_update_policy ON users;
        CREATE POLICY users_update_policy ON users
            FOR UPDATE
            USING (id = current_setting('app.current_user_id')::uuid);
    END IF;

    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'beats') THEN
        ALTER TABLE beats ENABLE ROW LEVEL SECURITY;
        
        -- Policy: Anyone can see public and approved beats, owners can see their own
        DROP POLICY IF EXISTS beats_select_policy ON beats;
        CREATE POLICY beats_select_policy ON beats
            FOR SELECT
            USING (
                (is_public = true AND is_approved = true) 
                OR user_id = current_setting('app.current_user_id')::uuid
            );
            
        -- Policy: Only owners can update their beats
        DROP POLICY IF EXISTS beats_update_policy ON beats;
        CREATE POLICY beats_update_policy ON beats
            FOR UPDATE
            USING (user_id = current_setting('app.current_user_id')::uuid);
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Create a monitoring view for admin dashboard
CREATE OR REPLACE VIEW admin_stats AS
SELECT 
    'users' as metric,
    COUNT(*) as total,
    COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '24 hours') as last_24h,
    COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '7 days') as last_7d
FROM users WHERE deleted_at IS NULL
UNION ALL
SELECT 
    'beats' as metric,
    COUNT(*) as total,
    COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '24 hours') as last_24h,
    COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '7 days') as last_7d
FROM beats WHERE deleted_at IS NULL
UNION ALL
SELECT 
    'challenges' as metric,
    COUNT(*) as total,
    COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '24 hours') as last_24h,
    COUNT(*) FILTER (WHERE created_at >= NOW() - INTERVAL '7 days') as last_7d
FROM challenges WHERE deleted_at IS NULL;

-- Comment for future reference
COMMENT ON FUNCTION create_performance_indexes() IS 'Creates performance indexes after GORM migrations complete';
COMMENT ON FUNCTION enable_row_level_security() IS 'Enables row-level security policies for data protection';
COMMENT ON VIEW admin_stats IS 'Administrative statistics view for monitoring';
