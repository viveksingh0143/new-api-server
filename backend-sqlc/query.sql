-- name: GetRole :one
SELECT * FROM roles WHERE id = ? LIMIT 1;

-- name: ListRoles :many
SELECT * FROM roles ORDER BY name;

-- name: CreateRole :execresult
INSERT INTO roles (name, status, last_updated_by) VALUES(?, ?, ?);

-- name: DeleteRole :exec
DELETE FROM roles WHERE id = ?;