-- schema created at https://dbdiagram.io/d/fintrack-6980c9a9bd82f5fce2609709
CREATE TYPE IF NOT EXISTS "account_type" AS ENUM (
  'bank',
  'cash',
  'credit_card',
  'investment',
  'other'
);

CREATE TYPE IF NOT EXISTS "credit_card_brand" AS ENUM (
  'visa',
  'mastercard',
  'amex',
  'discover',
  'jcb',
  'unionpay',
  'diners_club',
  'maestro',
  'unknown'
);

CREATE TYPE IF NOT EXISTS "transaction_type" AS ENUM (
  'credit',
  'debit',
  'transfer',
  'payment'
);

CREATE TABLE IF NOT EXISTS "tenants" (
  "id" UUID PRIMARY KEY DEFAULT (gen_random_uuid()),
  "name" TEXT NOT NULL,
  "created_at" "TIMESTAMPTZ" NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" "TIMESTAMPTZ" NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "deactivated_at" "TIMESTAMPTZ"
);

CREATE TABLE IF NOT EXISTS "users" (
  "id" UUID PRIMARY KEY DEFAULT (gen_random_uuid()),
  "name" TEXT NOT NULL,
  "email" "VARCHAR(254)" NOT NULL,
  "created_at" "TIMESTAMPTZ" NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" "TIMESTAMPTZ" NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "deactivated_at" "TIMESTAMPTZ"
);

CREATE TABLE IF NOT EXISTS "users_tenants" (
  "user_id" UUID NOT NULL,
  "tenant_id" UUID NOT NULL,
  "created_at" "TIMESTAMPTZ" NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" "TIMESTAMPTZ" NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "deactivated_at" "TIMESTAMPTZ",
  PRIMARY KEY ("user_id", "tenant_id")
);

CREATE TABLE IF NOT EXISTS "accounts" (
  "id" UUID PRIMARY KEY DEFAULT (gen_random_uuid()),
  "tenant_id" UUID NOT NULL,
  "name" TEXT NOT NULL,
  "initial_balance" "DECIMAL(15,2)" NOT NULL DEFAULT 0,
  "color" "VARCHAR(9)" NOT NULL,
  "currency" "VARCHAR(3)" NOT NULL,
  "icon" "VARCHAR(256)" NOT NULL,
  "type" account_type NOT NULL,
  "created_at" "TIMESTAMPTZ" NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "created_by" UUID NOT NULL,
  "updated_at" "TIMESTAMPTZ" NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "updated_by" UUID NOT NULL,
  "deactivated_at" "TIMESTAMPTZ",
  "deactivated_by" UUID
);

CREATE TABLE IF NOT EXISTS "credit_card_info" (
  "id" UUID PRIMARY KEY DEFAULT (gen_random_uuid()),
  "account_id" UUID NOT NULL,
  "last_four" "VARCHAR(4)" NOT NULL,
  "name" TEXT NOT NULL,
  "brand" credit_card_brand NOT NULL,
  "closing_date" DATE NOT NULL,
  "due_date" DATE NOT NULL,
  "created_at" "TIMESTAMPTZ" NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "created_by" UUID NOT NULL,
  "updated_at" "TIMESTAMPTZ" NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "updated_by" UUID NOT NULL,
  "deactivated_at" "TIMESTAMPTZ",
  "deactivated_by" UUID
);

CREATE TABLE IF NOT EXISTS "tags" (
  "id" UUID PRIMARY KEY DEFAULT (gen_random_uuid()),
  "tenant_id" UUID NOT NULL,
  "name" TEXT NOT NULL,
  "deactivated_at" "TIMESTAMPTZ"
);

CREATE TABLE IF NOT EXISTS "categories" (
  "id" UUID PRIMARY KEY DEFAULT (gen_random_uuid()),
  "parent_category" UUID,
  "tenant_id" UUID NOT NULL,
  "name" TEXT NOT NULL,
  "deactivated_at" "TIMESTAMPTZ",
  "color" "VARCHAR(9)" NOT NULL,
  "icon" "VARCHAR(256)" NOT NULL
);

CREATE TABLE IF NOT EXISTS "transactions" (
  "id" UUID PRIMARY KEY DEFAULT (gen_random_uuid()),
  "previous_sibling_transaction_id" UUID,
  "next_sibling_transaction_id" UUID,
  "tenant_id" UUID NOT NULL,
  "from_account_id" UUID NOT NULL,
  "to_account_id" UUID,
  "amount" "NUMERIC(10,2)" NOT NULL,
  "accrual_month" "VARCHAR(6)" NOT NULL,
  "transaction_type" transaction_type NOT NULL,
  "category_id" UUID NOT NULL,
  "due_date" DATE NOT NULL,
  "payment_date" DATE,
  "created_at" "TIMESTAMPTZ" NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "created_by" UUID NOT NULL,
  "updated_at" "TIMESTAMPTZ" NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  "updated_by" UUID NOT NULL,
  "deactivated_at" "TIMESTAMPTZ",
  "deactivated_by" UUID
);

CREATE TABLE IF NOT EXISTS "transactions_tags" (
  "transaction_id" UUID NOT NULL,
  "tag_id" UUID NOT NULL,
  PRIMARY KEY ("transaction_id", "tag_id")
);

CREATE INDEX IF NOT EXISTS ON "tenants" USING HASH ("deactivated_at");

CREATE INDEX IF NOT EXISTS ON "users" USING HASH ("deactivated_at");

CREATE INDEX IF NOT EXISTS ON "users_tenants" USING HASH ("tenant_id");

CREATE INDEX IF NOT EXISTS ON "users_tenants" USING HASH ("user_id");

CREATE INDEX IF NOT EXISTS ON "users_tenants" USING HASH ("deactivated_at");

CREATE INDEX IF NOT EXISTS ON "accounts" USING HASH ("tenant_id");

CREATE INDEX IF NOT EXISTS ON "accounts" USING HASH ("deactivated_at");

CREATE UNIQUE INDEX IF NOT EXISTS ON "credit_card_info" ("account_id", "deactivated_at");

CREATE INDEX IF NOT EXISTS ON "credit_card_info" USING HASH ("deactivated_at");

CREATE INDEX IF NOT EXISTS ON "tags" USING HASH ("tenant_id");

CREATE INDEX IF NOT EXISTS ON "tags" USING HASH ("deactivated_at");

CREATE INDEX IF NOT EXISTS ON "categories" USING HASH ("tenant_id");

CREATE INDEX IF NOT EXISTS ON "categories" USING HASH ("deactivated_at");

CREATE INDEX IF NOT EXISTS ON "transactions" USING HASH ("tenant_id");

CREATE INDEX IF NOT EXISTS ON "transactions" USING HASH ("accrual_month");

CREATE INDEX IF NOT EXISTS ON "transactions" ("transaction_type");

CREATE INDEX IF NOT EXISTS ON "transactions" USING HASH ("deactivated_at");

COMMENT ON COLUMN "accounts"."color" IS 'RGBA color (eg.: #ffAABB11';

COMMENT ON COLUMN "accounts"."currency" IS 'ISO currency code (eg.: BRL, USD)';

COMMENT ON COLUMN "categories"."color" IS 'RGBA color (eg.: #ffAABB11';

COMMENT ON COLUMN "transactions"."accrual_month" IS 'Year and month (eg.: YYYYMM)';

ALTER TABLE "users_tenants" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "users_tenants" ADD FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id");

ALTER TABLE "accounts" ADD FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id");

ALTER TABLE "accounts" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");

ALTER TABLE "accounts" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "accounts" ADD FOREIGN KEY ("deactivated_by") REFERENCES "users" ("id");

ALTER TABLE "credit_card_info" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "credit_card_info" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");

ALTER TABLE "credit_card_info" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "credit_card_info" ADD FOREIGN KEY ("deactivated_by") REFERENCES "users" ("id");

ALTER TABLE "tags" ADD FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id");

ALTER TABLE "categories" ADD FOREIGN KEY ("parent_category") REFERENCES "categories" ("id");

ALTER TABLE "categories" ADD FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("previous_sibling_transaction_id") REFERENCES "transactions" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("next_sibling_transaction_id") REFERENCES "transactions" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("updated_by") REFERENCES "users" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("deactivated_by") REFERENCES "users" ("id");

ALTER TABLE "transactions_tags" ADD FOREIGN KEY ("transaction_id") REFERENCES "transactions" ("id");

ALTER TABLE "transactions_tags" ADD FOREIGN KEY ("tag_id") REFERENCES "tags" ("id");


---- create above / drop below ----

DROP TABLE IF EXISTS "transactions_tags";
DROP TABLE IF EXISTS "transactions";
DROP TABLE IF EXISTS "categories";
DROP TABLE IF EXISTS "tags";
DROP TABLE IF EXISTS "credit_card_info";
DROP TABLE IF EXISTS "accounts";
DROP TABLE IF EXISTS "users_tenants";
DROP TABLE IF EXISTS "users";
DROP TABLE IF EXISTS "tenants";

DROP TYPE IF EXISTS "transaction_type";
DROP TYPE IF EXISTS "credit_card_brand";
DROP TYPE IF EXISTS "account_type";
