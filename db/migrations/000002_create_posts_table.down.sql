-- Drop posts table and related objects
DROP TRIGGER IF EXISTS update_posts_updated_at ON posts;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS posts;