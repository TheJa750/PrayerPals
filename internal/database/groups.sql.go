// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: groups.sql

package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const addUserToGroup = `-- name: AddUserToGroup :exec
INSERT INTO users_groups (user_id, group_id, role)
VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING
`

type AddUserToGroupParams struct {
	UserID  uuid.UUID
	GroupID uuid.UUID
	Role    string
}

func (q *Queries) AddUserToGroup(ctx context.Context, arg AddUserToGroupParams) error {
	_, err := q.db.ExecContext(ctx, addUserToGroup, arg.UserID, arg.GroupID, arg.Role)
	return err
}

const createGroup = `-- name: CreateGroup :one
INSERT INTO groups (name, description, owner_id)
VALUES ($1, $2, $3)
RETURNING id, name, description, created_at, updated_at, owner_id
`

type CreateGroupParams struct {
	Name        string
	Description sql.NullString
	OwnerID     uuid.NullUUID
}

func (q *Queries) CreateGroup(ctx context.Context, arg CreateGroupParams) (Group, error) {
	row := q.db.QueryRowContext(ctx, createGroup, arg.Name, arg.Description, arg.OwnerID)
	var i Group
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.OwnerID,
	)
	return i, err
}

const deleteGroup = `-- name: DeleteGroup :exec
DELETE FROM groups
WHERE id = $1
`

func (q *Queries) DeleteGroup(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteGroup, id)
	return err
}

const getGroupByID = `-- name: GetGroupByID :one
SELECT id, name, description, created_at, updated_at, owner_id
FROM groups
WHERE id = $1
`

func (q *Queries) GetGroupByID(ctx context.Context, id uuid.UUID) (Group, error) {
	row := q.db.QueryRowContext(ctx, getGroupByID, id)
	var i Group
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.OwnerID,
	)
	return i, err
}

const getGroupMembersIDs = `-- name: GetGroupMembersIDs :many
SELECT user_id
FROM users_groups
WHERE group_id = $1
`

func (q *Queries) GetGroupMembersIDs(ctx context.Context, groupID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := q.db.QueryContext(ctx, getGroupMembersIDs, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []uuid.UUID
	for rows.Next() {
		var user_id uuid.UUID
		if err := rows.Scan(&user_id); err != nil {
			return nil, err
		}
		items = append(items, user_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getGroupSpecialRoles = `-- name: GetGroupSpecialRoles :many
SELECT user_id, role
FROM users_groups
WHERE group_id = $1 AND role != 'member'
`

type GetGroupSpecialRolesRow struct {
	UserID uuid.UUID
	Role   string
}

func (q *Queries) GetGroupSpecialRoles(ctx context.Context, groupID uuid.UUID) ([]GetGroupSpecialRolesRow, error) {
	rows, err := q.db.QueryContext(ctx, getGroupSpecialRoles, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetGroupSpecialRolesRow
	for rows.Next() {
		var i GetGroupSpecialRolesRow
		if err := rows.Scan(&i.UserID, &i.Role); err != nil {
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

const getGroupsForUser = `-- name: GetGroupsForUser :many
SELECT group_id
FROM users_groups
WHERE user_id = $1
`

func (q *Queries) GetGroupsForUser(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := q.db.QueryContext(ctx, getGroupsForUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []uuid.UUID
	for rows.Next() {
		var group_id uuid.UUID
		if err := rows.Scan(&group_id); err != nil {
			return nil, err
		}
		items = append(items, group_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const removeUserFromGroup = `-- name: RemoveUserFromGroup :exec
DELETE FROM users_groups
WHERE user_id = $1 AND group_id = $2
`

type RemoveUserFromGroupParams struct {
	UserID  uuid.UUID
	GroupID uuid.UUID
}

func (q *Queries) RemoveUserFromGroup(ctx context.Context, arg RemoveUserFromGroupParams) error {
	_, err := q.db.ExecContext(ctx, removeUserFromGroup, arg.UserID, arg.GroupID)
	return err
}

const resetGroups = `-- name: ResetGroups :exec
TRUNCATE TABLE groups CASCADE
`

func (q *Queries) ResetGroups(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, resetGroups)
	return err
}
