CREATE TYPE "category_type" AS ENUM (
  'income',
  'expense',
  'transfer'
);

ALTER TABLE "categories" ADD COLUMN "type" category_type NOT NULL DEFAULT 'expense';

ALTER TABLE "categories" ALTER COLUMN "type" DROP DEFAULT;

---- create above / drop below ----

ALTER TABLE "categories" DROP COLUMN "type";
DROP TYPE "category_type";
