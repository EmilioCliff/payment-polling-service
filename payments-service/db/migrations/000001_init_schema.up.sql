CREATE TABLE "transactions" (
    "transaction_id" uuid PRIMARY KEY,
    "user_id" bigint NOT NULL,
    "action" varchar NOT NULL,
    "amount" integer NOT NULL,
    "phone_number" varchar NOT NULL,
    "network_node" varchar NOT NULL,
    "narration" text NOT NULL,
    "status" boolean NOT NULL DEFAULT false,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT transaction_types CHECK (action IN ('payment', 'withdrawal'))
);