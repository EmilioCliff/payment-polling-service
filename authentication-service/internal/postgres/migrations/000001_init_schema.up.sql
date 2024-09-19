CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY,
    "full_name" varchar NOT NULL,
    "email" varchar UNIQUE NOT NULL,
    "password" varchar NOT NULL,
    "payd_username" varchar NOT NULL,
    "payd_account_id" varchar NOT NULL,
    "payd_username_key" varchar NOT NULL,
    "payd_password_key" varchar NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX ON "users" ("email");