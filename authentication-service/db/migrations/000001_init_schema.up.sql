CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY,
    "full_name" varchar NOT NULL,
    "email" varchar UNIQUE NOT NULL,
    "password" varchar NOT NULL,
    "access_token" varchar,
    "created_at" timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX ON "users" ("email");