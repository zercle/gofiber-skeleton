-- Create view for posts with author information (optimized for UUIDv7)
CREATE OR REPLACE VIEW posts_with_authors AS
SELECT
    p.id,
    p.title,
    p.content,
    p.status,
    p.user_id,
    u.full_name,
    u.email,
    p.created_at,
    p.updated_at
FROM posts p
JOIN users u ON p.user_id = u.id
WHERE u.is_active = true;

-- Add comment to the view
COMMENT ON VIEW posts_with_authors IS 'Optimized view for posts with active author information';

-- Create index to support the view (through underlying tables)
-- The indexes already exist on the underlying tables, but we can create a materialized view if needed
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_posts_with_authors_status_created
ON posts(status, created_at DESC)
WHERE user_id IN (SELECT id FROM users WHERE is_active = true);