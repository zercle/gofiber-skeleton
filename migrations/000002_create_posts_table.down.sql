DROP TRIGGER IF EXISTS update_posts_updated_at ON posts;
DROP INDEX IF EXISTS idx_posts_deleted_at;
DROP INDEX IF EXISTS idx_posts_created_at;
DROP INDEX IF EXISTS idx_posts_published_at;
DROP INDEX IF EXISTS idx_posts_is_published;
DROP INDEX IF EXISTS idx_posts_author_id;
DROP TABLE IF EXISTS posts;