-- 1. Make the column mandatory
ALTER TABLE transactions ALTER COLUMN category_id SET NOT NULL;

-- 2. Ensure the category actually exists in the categories table
ALTER TABLE transactions 
ADD CONSTRAINT fk_transaction_category 
FOREIGN KEY (category_id) 
REFERENCES categories (id);