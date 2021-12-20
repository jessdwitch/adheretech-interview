CREATE TABLE IF NOT EXISTS "secret_tokens" (
  "id" SERIAL PRIMARY KEY,
  "data" text NOT NULL CHECK (data NOT LIKE '%-%')
);