-- name: GetUserByID :one
SELECT sqlc.embed(users) FROM users
WHERE users.id = $1;

-- name: GetUserByProviderID :one
SELECT sqlc.embed(users), sqlc.embed(user_auth) FROM users 
JOIN user_auth ON users.id = user_auth.user_id
WHERE user_auth.provider_id = $1
AND user_auth.provider = $2;

-- name: GetUserBySessionToken :one
SELECT sqlc.embed(users), sqlc.embed(user_session) FROM users
JOIN user_session ON users.id = user_session.user_id
WHERE user_session.token = $1;

-- name: GetUserAuthByUserID :many
SELECT sqlc.embed(user_auth) FROM user_auth
WHERE user_auth.user_id = $1;

-- name: GetUserByEmail :one
SELECT sqlc.embed(users) FROM users
JOIN user_auth ON users.id = user_auth.user_id
WHERE user_auth.email = sqlc.arg(email);

-- name: GetUserAuthByEmail :many
SELECT sqlc.embed(user_auth) FROM user_auth
WHERE user_auth.email = $1;

-- name: CreateUser :one
INSERT INTO users (id, given_name, family_name, primary_email) 
VALUES (
    sqlc.arg(id),
    sqlc.narg(given_name), 
    sqlc.narg(family_name),
    sqlc.narg(primary_email)
) RETURNING *;
 
-- name: UpdateUser :one
UPDATE users SET 
    id = sqlc.arg(id), 
    given_name = sqlc.narg(given_name), 
    family_name = sqlc.narg(family_name),
    primary_email = sqlc.narg(primary_email)
WHERE id = sqlc.arg(id) RETURNING *;

-- name: CreateUserAuth :one
INSERT INTO user_auth (user_id, value, provider, provider_id, provider_data, email, email_verified) 
VALUES (sqlc.arg(user_id), sqlc.arg(value), sqlc.arg(provider), sqlc.arg(provider_id), sqlc.arg(provider_data), sqlc.narg(email), sqlc.arg(email_verified)) RETURNING *;

-- name: CreateUserSession :one
INSERT INTO user_session (user_id, token, expires_at) VALUES (sqlc.arg(user_id), sqlc.arg(token), sqlc.arg(expires_at)) RETURNING *;

-- name: ExpireUserSession :exec
UPDATE user_session SET user_expired = TRUE WHERE user_session.id = $1;
