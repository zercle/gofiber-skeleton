-- Create posts table
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes optimized for UUIDv7 (time-ordered UUIDs)
CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_status ON posts(status);
CREATE INDEX idx_posts_created_at_v7 ON posts(created_at DESC);
CREATE INDEX idx_posts_id_created_at ON posts(id, created_at DESC);
CREATE INDEX idx_posts_user_status_created_at ON posts(user_id, status, created_at DESC);
CREATE INDEX idx_posts_status_created_at_v7 ON posts(status, created_at DESC);

-- Partial indexes for common query patterns
CREATE INDEX idx_posts_published ON posts(created_at DESC) WHERE status = 'published';
CREATE INDEX idx_posts_user_active ON posts(user_id, created_at DESC)
WHERE status IN ('published', 'draft');

-- Full-text search indexes
CREATE INDEX idx_posts_title_gin ON posts USING gin(title gin_trgm_ops);
CREATE INDEX idx_posts_content_gin ON posts USING gin(content gin_trgm_ops);
CREATE INDEX idx_posts_search_gin ON posts USING gin(to_tsvector('english', title || ' ' || content));

-- Create trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_posts_updated_at
    BEFORE UPDATE ON posts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add table comment
COMMENT ON TABLE posts IS 'Posts created by users with draft/published status';