-- 1. Remove the old case-sensitive constraint
ALTER TABLE categories DROP CONSTRAINT IF EXISTS unique_user_category_name;

-- 2. Create the case-insensitive unique index
CREATE UNIQUE INDEX unique_user_category_name_idx
ON categories (user_id, LOWER(name));