-- 1. ตาราง Accounts
CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL CHECK ("balance" >= 0), -- บังคับยอดเงินห้ามติดลบเด็ดขาด (No Overdraft)
  "currency" varchar(3) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- 2. ตาราง Entries (Ledger เงินเข้า-ออก)
CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "amount" bigint NOT NULL, -- บวกคือเงินเข้า, ลบคือเงินออก
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- 3. ตาราง Transfers (ประวัติการโอนเงิน)
CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "from_account_id" bigint NOT NULL,
  "to_account_id" bigint NOT NULL,
  "amount" bigint NOT NULL, -- ต้องมากกว่า 0 เสมอ
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- กำหนด Foreign Key ผูกความสัมพันธ์
ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE;
ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE;
ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE;

-- สร้าง Index เพื่อความเร็วสูงสุดในการค้นหา (Query Performance)
CREATE INDEX ON "accounts" ("owner");
CREATE INDEX ON "entries" ("account_id");
CREATE INDEX ON "transfers" ("from_account_id");
CREATE INDEX ON "transfers" ("to_account_id");
CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");