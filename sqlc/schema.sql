CREATE TABLE "user" (
  "id" uuid PRIMARY KEY,
  "name" text NOT NULL,
  "email" text UNIQUE NOT NULL,
  "password" text NOT NULL,
  "token" text NOT NULL
);

CREATE TABLE "transaction" (
  "id" uuid PRIMARY KEY,
  "user_id" uuid NOT NULL,
  "amount" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "transaction" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");
