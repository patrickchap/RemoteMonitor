-- name: GetHosts :many
SELECT * FROM hosts 
Limit ?
Offset ?;

-- name: GetHost :one
SELECT * FROM hosts WHERE id = ?;

-- name: GetHostByHostname :one
SELECT * FROM hosts WHERE host_name = ?;

-- name: CreateHost :one
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
RETURNING *;

