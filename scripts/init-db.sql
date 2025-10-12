-- PostgreSQL 18 has native UUIDv7 support, no extensions needed
-- This script is kept for backwards compatibility but UUIDv7 is used directly in migrations

-- Create custom function for UUID generation (PostgreSQL 18+ uses uuidv7())
CREATE OR REPLACE FUNCTION generate_uuid() RETURNS UUID AS $$
BEGIN
    RETURN uuidv7();
END;
$$ LANGUAGE plpgsql;