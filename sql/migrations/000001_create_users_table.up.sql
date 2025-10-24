CREATE TABLE users (
  id VARCHAR(36) PRIMARY KEY COMMENT 'UUID v7',
  email VARCHAR(255) NOT NULL UNIQUE COMMENT 'User email address',
  password_hash VARCHAR(255) NOT NULL COMMENT 'Bcrypt hashed password',
  first_name VARCHAR(100) COMMENT 'User first name',
  last_name VARCHAR(100) COMMENT 'User last name',
  is_active BOOLEAN DEFAULT TRUE COMMENT 'Whether the user is active',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation timestamp',
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Last update timestamp',
  deleted_at TIMESTAMP NULL COMMENT 'Soft delete timestamp',

  INDEX idx_email (email),
  INDEX idx_created_at (created_at),
  INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Users table for authentication and basic user info';
