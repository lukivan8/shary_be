-- Drop favorite_items table
DROP TABLE IF EXISTS favorite_items CASCADE;

-- Drop rents table
DROP TABLE IF EXISTS rents CASCADE;

-- Drop item_photos table
DROP TABLE IF EXISTS item_photos CASCADE;

-- Drop statuses table
DROP TABLE IF EXISTS statuses CASCADE;

-- Drop categories table
DROP TABLE IF EXISTS categories CASCADE;

-- Drop users table
DROP TABLE IF EXISTS users CASCADE; 

-- Drop users table
DROP TABLE IF EXISTS users CASCADE;



-- SQL in this section is executed when the migration is rolled back.

-- Delete indexes before deleting tables
DROP INDEX IF EXISTS idx_bookings_item_id;
DROP INDEX IF EXISTS idx_blocked_periods_item_id;
DROP INDEX IF EXISTS idx_bookings_user_id;

-- Drop bookings table
DROP TABLE IF EXISTS blocked_periods CASCADE;
DROP TABLE IF EXISTS bookings CASCADE;
