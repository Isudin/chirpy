-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES ($1, NOW(), NOW(), $2, $3)
RETURNING *;

-- name: RefreshToken :exec
UPDATE refresh_tokens
SET expires_at = $1
WHERE token = $2;

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1;

-- name: GetTokenData :one
SELECT * FROM refresh_tokens
WHERE token = $1 AND revoked_at IS NULL;