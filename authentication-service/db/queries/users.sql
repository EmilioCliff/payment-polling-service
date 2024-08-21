-- name: RegisterUser :one
INSERT INTO users (
    full_name, payd_username, email, password, payd_username_key, payd_password_key
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1
LIMIT 1;