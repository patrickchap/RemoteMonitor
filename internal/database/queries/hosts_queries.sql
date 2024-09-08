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

-- name: GetHostsWithServices :many 
SELECT h.id, h.host_name, hs.status, s.service_name
FROM hosts h LEFT JOIN host_services AS hs 
ON h.id = hs.host_id LEFT JOIN services AS s 
ON hs.service_id = s.id
Limit ?
Offset ?;
