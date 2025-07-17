-- name: CreateUser :one
INSERT INTO users (
  name,
  email,
  password,
  token
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: CreateTransaction :one
INSERT INTO transactions (
  user_id,
  amount
) VALUES (
  $1, $2
)
RETURNING *;

-- name: ListTransactionsByUserID :many
SELECT * FROM transactions
WHERE user_id = $1
ORDER BY created_at DESC;
