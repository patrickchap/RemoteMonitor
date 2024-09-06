// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: hosts_queries.sql

package database

import (
	"context"
	"database/sql"
)

const createHost = `-- name: CreateHost :one
INSERT INTO hosts (
    host_name,
    canonical_name,
    url,
    ip,
    ipv6,
    location,
    os,
    active,
    last_updated
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) 
RETURNING id, host_name, canonical_name, url, ip, ipv6, location, os, active, last_updated
`

type CreateHostParams struct {
	HostName      string         `json:"host_name"`
	CanonicalName sql.NullString `json:"canonical_name"`
	Url           sql.NullString `json:"url"`
	Ip            sql.NullString `json:"ip"`
	Ipv6          sql.NullString `json:"ipv6"`
	Location      sql.NullString `json:"location"`
	Os            sql.NullString `json:"os"`
	Active        sql.NullInt64  `json:"active"`
	LastUpdated   sql.NullTime   `json:"last_updated"`
}

func (q *Queries) CreateHost(ctx context.Context, arg CreateHostParams) (Host, error) {
	row := q.db.QueryRowContext(ctx, createHost,
		arg.HostName,
		arg.CanonicalName,
		arg.Url,
		arg.Ip,
		arg.Ipv6,
		arg.Location,
		arg.Os,
		arg.Active,
		arg.LastUpdated,
	)
	var i Host
	err := row.Scan(
		&i.ID,
		&i.HostName,
		&i.CanonicalName,
		&i.Url,
		&i.Ip,
		&i.Ipv6,
		&i.Location,
		&i.Os,
		&i.Active,
		&i.LastUpdated,
	)
	return i, err
}

const getHost = `-- name: GetHost :one
SELECT id, host_name, canonical_name, url, ip, ipv6, location, os, active, last_updated FROM hosts WHERE id = ?
`

func (q *Queries) GetHost(ctx context.Context, id int64) (Host, error) {
	row := q.db.QueryRowContext(ctx, getHost, id)
	var i Host
	err := row.Scan(
		&i.ID,
		&i.HostName,
		&i.CanonicalName,
		&i.Url,
		&i.Ip,
		&i.Ipv6,
		&i.Location,
		&i.Os,
		&i.Active,
		&i.LastUpdated,
	)
	return i, err
}

const getHostByHostname = `-- name: GetHostByHostname :one
SELECT id, host_name, canonical_name, url, ip, ipv6, location, os, active, last_updated FROM hosts WHERE host_name = ?
`

func (q *Queries) GetHostByHostname(ctx context.Context, hostName string) (Host, error) {
	row := q.db.QueryRowContext(ctx, getHostByHostname, hostName)
	var i Host
	err := row.Scan(
		&i.ID,
		&i.HostName,
		&i.CanonicalName,
		&i.Url,
		&i.Ip,
		&i.Ipv6,
		&i.Location,
		&i.Os,
		&i.Active,
		&i.LastUpdated,
	)
	return i, err
}

const getHosts = `-- name: GetHosts :many
SELECT id, host_name, canonical_name, url, ip, ipv6, location, os, active, last_updated FROM hosts 
Limit ?
Offset ?
`

type GetHostsParams struct {
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

func (q *Queries) GetHosts(ctx context.Context, arg GetHostsParams) ([]Host, error) {
	rows, err := q.db.QueryContext(ctx, getHosts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Host{}
	for rows.Next() {
		var i Host
		if err := rows.Scan(
			&i.ID,
			&i.HostName,
			&i.CanonicalName,
			&i.Url,
			&i.Ip,
			&i.Ipv6,
			&i.Location,
			&i.Os,
			&i.Active,
			&i.LastUpdated,
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