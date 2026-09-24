-- Account table

CREATE TABLE "accounts" (
    "id" bigserial PRIMARY KEY,
    "owner" varchar NOT NULL,
    "balance" bigint NOT NULL,
    "currency" varchar(3) NOT NULL,
    "create_at" timestamptz NOT NULL DEFAULT (now())
);

-- Ledger / Entries table

CREATE TABLE "entries" (
    "id" bigserial PRIMARY KEY,
    "account_id" bigint NOT NULL,
    "amount" bigint NOT NULL,
    "create_at" timestamptz NOT NULL DEFAULT (now(now))
);

-- Tranfers table

CREATE TABLE "tranfers" (
    "id" bigserial PRIMARY KEY,
    "from_account_id" bigserial NOT NULL,
    "to_account_id" bigserial NOT NULL,
    "amout" bigint NOT NULL,
    "create_at" timestamptz NOT NULL DEFAULT (now())
);

-- Assign foreign key

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE;
ALTER TABLE "tranfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE;
ALTER TABLE "tranfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "account_id" ("id") ON DELETE CASCADE;

-- create index
CREATE INDEX ON "accounts" ("owner");
CREATE INDEX ON "entries" ("account_id")
CREATE INDEX ON "tranfers" ("from_account_id");
CREATE INDEX ON "tranfers" ("to_account_id");
CREATE INDEX ON "tranfers" ("from_account_id","to_account_id")