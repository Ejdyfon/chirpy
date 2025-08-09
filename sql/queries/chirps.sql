-- name: CreatChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteChirps :exec
delete from chirps;

-- name: GetAllChirps :many
select * from chirps;

-- name: GetAllChirpsByAuthor :many
select * from chirps where user_id = $1;

-- name: GetChirpById :one
select * from chirps where id = $1;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1;