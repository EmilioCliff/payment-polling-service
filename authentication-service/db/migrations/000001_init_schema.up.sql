CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY,
    "username" varchar NOT NULL DEFAULT 'test',
    "full_name" varchar NOT NULL,
    "email" varchar UNIQUE NOT NULL,
    "password" varchar NOT NULL,
    "payd_username_key" varchar NOT NULL DEFAULT 'test',
    "payd_password_key" varchar NOT NULL DEFAULT 'test',
    "access_token" varchar,
    "created_at" timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX ON "users" ("email");