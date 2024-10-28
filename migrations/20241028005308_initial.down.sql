-- Drop trigger and function
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column;

-- Drop indexes
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_status;
DROP INDEX IF EXISTS idx_users_created_at;

-- Drop users table
DROP TABLE IF EXISTS users;

-- Drop custom types
DROP TYPE IF EXISTS user_status;

-- Drop UUID extension
DROP EXTENSION IF EXISTS "uuid-ossp";
