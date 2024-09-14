-- name: GetHostServices :many
SELECT hs.*, h.host_name, s.service_name
FROM host_services as hs
JOIN hosts as h ON hs.host_id = h.id
JOIN services as s ON hs.service_id = s.id
WHERE hs.host_id = ?
AND hs.active = 1;


-- name: CreateHostService :one
INSERT INTO host_services(host_id, service_id)
VALUES(?, ?)
RETURNING *;
