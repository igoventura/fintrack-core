-- Write your migrate up statements here
ALTER TABLE "categories" ADD COLUMN "created_at" TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP);
ALTER TABLE "categories" ADD COLUMN "created_by" UUID NOT NULL;
ALTER TABLE "categories" ADD COLUMN "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP);
ALTER TABLE "categories" ADD COLUMN "updated_by" UUID NOT NULL;
ALTER TABLE "categories" ADD COLUMN "deactivated_by" UUID;

-- Add Foreign Keys
ALTER TABLE "categories" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");
ALTER TABLE "categories" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");
ALTER TABLE "categories" ADD FOREIGN KEY ("deactivated_by") REFERENCES "users" ("id");

---- create above / drop below ----

ALTER TABLE "categories" DROP COLUMN "deactivated_by";
ALTER TABLE "categories" DROP COLUMN "updated_by";
ALTER TABLE "categories" DROP COLUMN "updated_at";
ALTER TABLE "categories" DROP COLUMN "created_by";
ALTER TABLE "categories" DROP COLUMN "created_at";
