-- Remove role column from users table
ALTER TABLE users DROP COLUMN IF EXISTS role;

-- Drop the index
DROP INDEX IF EXISTS idx_users_role;