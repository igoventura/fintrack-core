-- Resize color columns to support up to 128 characters
ALTER TABLE accounts ALTER COLUMN color TYPE VARCHAR(128);
ALTER TABLE categories ALTER COLUMN color TYPE VARCHAR(128);

---- create above / drop below ----

-- Revert color columns to original size
ALTER TABLE categories ALTER COLUMN color TYPE VARCHAR(9);
ALTER TABLE accounts ALTER COLUMN color TYPE VARCHAR(9);

