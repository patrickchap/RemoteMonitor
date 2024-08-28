// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: hosts_queries.sql

package database

import (
	"context"
)

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
