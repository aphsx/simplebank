CREATE TABLE "users" (
	"username" varchar PRIMARY KEY,
	"hashed_password" varchar NOT NULL,
	"full_name" varchar NOT NULL,
	"email" varchar UNIQUE NOT NULL,
	"password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
	"created_at" timestamptz NOT NULL DEFAULT (now())
);

-- ผูกเจ้าของบัญชีเข้ากับ username ในตาราง users
ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

-- กฎ: ผู้ใช้ 1 คน ห้ามเปิดบัญชีสกุลเงินเดียวกันซ้ำซ้อน (เช่น ​aphsx มีบัญชี USD ได้แค่เล่มเดียว)
ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency");