-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteUsers :exec
delete from users;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUser :one
update users set email = $1, hashed_password = $2, updated_at=NOW() where id = $3
returning *;

-- name: UpgradeToChirpyRed :one
UPDATE users SET is_chirpy_red = true, updated_at = NOW()
WHERE id = $1
RETURNING *;