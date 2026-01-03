-- Ensure a user cannot have two categories with the same name
ALTER TABLE categories
ADD CONSTRAINT unique_user_category_name UNIQUE (user_id, name);