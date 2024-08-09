-- name: CreateUser :one
INSERT INTO
    users (
        username,
        first_name,
        last_name,
        email,
        passwhash,
        is_active,
        created_at
    )
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING
    username,
    first_name,
    last_name,
    email;

-- name: GetUserCredentialByEmail :one
SELECT
    user_id,
    username,
    email,
    passwhash
FROM
    users
WHERE
    email = $1 AND deleted_at IS NULL;

-- name: DeactivateTokenSession :exec
UPDATE
    users_sessions
SET
    is_active = FALSE,
    used_at = $1
WHERE
    user_id = $2;

-- name: InsertTokenSession :one
INSERT INTO
    users_sessions (
        user_id,
        token,
        created_at,
        is_active
    )
VALUES ($1, $2, $3, $4)
RETURNING
    *;
