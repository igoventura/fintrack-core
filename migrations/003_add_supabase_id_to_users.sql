-- Write your migrate up statements here
ALTER TABLE "users" ADD COLUMN "supabase_id" VARCHAR(128) UNIQUE;
CREATE INDEX ON "users" ("supabase_id");

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
DROP INDEX IF EXISTS "users_supabase_id_idx";
ALTER TABLE "users" DROP COLUMN IF EXISTS "supabase_id";
