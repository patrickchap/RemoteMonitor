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

-- name: UpdateHost :one
UPDATE hosts SET 
    host_name = ?,
    canonical_name = ?,
    url = ?,
    ip = ?,
    ipv6 = ?,
    last_updated = ?
WHERE id = ?
RETURNING *;


-- name: GetHostsWithServices :many 
SELECT
  h.id,
  h.host_name,
  hs.status,
  s.service_name,
  hs.active,
  h.active AS host_active
FROM
  hosts h
  LEFT JOIN host_services AS hs ON h.id = hs.host_id
  LEFT JOIN services AS s ON hs.service_id = s.id
WHERE host_active = 1
AND hs.active = 1
AND h.active = 1
Limit ?
Offset ?;
