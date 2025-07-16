-- name: CreateUser :one
INSERT INTO "user" (
  id,
  name,
  email,
  password,
  token
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM "user"
WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM "user"
WHERE id = $1 LIMIT 1;

-- name: CreateTransaction :one
INSERT INTO "transaction" (
  id,
  user_id,
  amount
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: ListTransactionsByUserID :many
SELECT * FROM "transaction"
WHERE user_id = $1
ORDER BY created_at DESC;
