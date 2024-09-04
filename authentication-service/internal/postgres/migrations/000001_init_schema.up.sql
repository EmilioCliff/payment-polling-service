CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY,
    "payd_username" varchar NOT NULL,
    "full_name" varchar NOT NULL,
    "email" varchar UNIQUE NOT NULL,
    "password" varchar NOT NULL,
    "payd_username_key" varchar NOT NULL,
    "payd_password_key" varchar NOT NULL,
    "access_token" varchar,
    "created_at" timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX ON "users" ("email");