DROP TRIGGER IF EXISTS update_posts_updated_at ON posts;
DROP INDEX IF EXISTS idx_posts_published_at;
DROP INDEX IF EXISTS idx_posts_is_published;
DROP INDEX IF EXISTS idx_posts_slug;
DROP INDEX IF EXISTS idx_posts_user_id;
DROP TABLE IF EXISTS posts;