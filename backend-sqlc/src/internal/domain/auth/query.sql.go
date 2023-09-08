// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: query.sql

package auth

import (
	"context"
	"database/sql"
)

const createRole = `-- name: CreateRole :execresult
INSERT INTO roles (name, status, last_updated_by) VALUES(?, ?, ?)
`

type CreateRoleParams struct {
	Name          string
	Status        RolesStatus
	LastUpdatedBy sql.NullString
}

func (q *Queries) CreateRole(ctx context.Context, arg CreateRoleParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createRole, arg.Name, arg.Status, arg.LastUpdatedBy)
}

const deleteRole = `-- name: DeleteRole :exec
DELETE FROM roles WHERE id = ?
`

func (q *Queries) DeleteRole(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteRole, id)
	return err
}

const getRole = `-- name: GetRole :one
SELECT id, name, status, created_at, updated_at, last_updated_by FROM roles WHERE id = ? LIMIT 1
`

func (q *Queries) GetRole(ctx context.Context, id int64) (Role, error) {
	row := q.db.QueryRowContext(ctx, getRole, id)
	var i Role
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastUpdatedBy,
	)
	return i, err
}

const listRoles = `-- name: ListRoles :many
SELECT id, name, status, created_at, updated_at, last_updated_by FROM roles ORDER BY name
`

func (q *Queries) ListRoles(ctx context.Context) ([]Role, error) {
	rows, err := q.db.QueryContext(ctx, listRoles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Role
	for rows.Next() {
		var i Role
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.LastUpdatedBy,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}