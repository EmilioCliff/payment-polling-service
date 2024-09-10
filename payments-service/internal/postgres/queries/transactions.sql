-- name: CreateTransaction :one
INSERT INTO transactions (
    transaction_id, payd_transaction_ref,user_id, action, amount, phone_number, network_node, narration
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetTransaction :one
SELECT * FROM transactions
WHERE transaction_id = $1;

-- name: UpdateTransaction :one
UPDATE transactions
SET status = $2,
    updated_at = now()
WHERE transaction_id = $1
RETURNING *;
