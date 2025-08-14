-- Add role column to users table
ALTER TABLE users ADD COLUMN role VARCHAR(50) NOT NULL DEFAULT 'customer' CHECK (role IN ('admin', 'customer'));

-- Create index for role-based queries
CREATE INDEX idx_users_role ON users(role);