CREATE TABLE "transactions" (
    "transaction_id" uuid PRIMARY KEY,
    "payd_transaction_ref" varchar NOT NULL,
    "user_id" bigint NOT NULL,
    "action" varchar NOT NULL,
    "amount" integer NOT NULL,
    "phone_number" varchar NOT NULL,
    "network_node" varchar NOT NULL,
    "narration" text NOT NULL,
    "status" boolean NOT NULL DEFAULT false,
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "created_at" timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT transaction_types CHECK (action = 'payment' OR action = 'withdrawal' )
);