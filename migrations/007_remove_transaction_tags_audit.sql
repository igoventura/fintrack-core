-- Remove audit columns from transactions_tags
-- Rationale: Junction table logic simplifies to replace-all strategy, audit overhead not needed.

ALTER TABLE transactions_tags
DROP COLUMN created_at,
DROP COLUMN updated_at,
DROP COLUMN deactivated_at;

---- create above / drop below ----

ALTER TABLE transactions_tags
ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN deactivated_at TIMESTAMPTZ;
