CREATE TABLE "users" (
  "user_name" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00z',
  "created_at" timestamptz NOT NULL DEFAULT 'now()'
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("user_name");
CREATE UNIQUE INDEX ON "accounts" ("owner", "currency");
-- this query automatically creates a unique composite index under the hood so is basically the same
--ALTER TABLE "accounts" ADD CONSTRAINT ("owner_currency_key") UNIQUE ("owner", "currency");