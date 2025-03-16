-- name: GetHostServices :many
SELECT hs.*, h.host_name, s.service_name
FROM host_services as hs
JOIN hosts as h ON hs.host_id = h.id
JOIN services as s ON hs.service_id = s.id
WHERE hs.host_id = ?
AND hs.active = 1;

-- name: GetHostServicesByStatus :many
SELECT hs.*, h.host_name, s.service_name
FROM host_services as hs
JOIN hosts as h ON hs.host_id = h.id
JOIN services as s ON hs.service_id = s.id
WHERE hs.status = ?
AND hs.active = 1;

-- name: GetHostServicesStatuses :one
select 
	(select count(id) from host_services where active = 1 and status = 'pending') as pending,
	(select count(id) from host_services where active = 1 and status = 'healthy') as healthy,
	(select count(id) from host_services where active = 1 and status = 'warning') as warning,
	(select count(id) from host_services where active = 1 and status = 'problem') as problem;

-- name: CreateHostService :one
INSERT INTO host_services(host_id, service_id)
VALUES(?, ?)
RETURNING *;

-- name: GetHostService :one
SELECT hs.*, h.host_name, h.url, s.service_name
FROM host_services as hs
JOIN hosts as h ON hs.host_id = h.id
JOIN services as s ON hs.service_id = s.id
WHERE hs.id = ?
AND hs.active = 1;

-- name: DeleteHostService :one
UPDATE host_services SET 
    active = 0
WHERE id = ?
RETURNING *;

-- name: ReactivateHostService :one
UPDATE host_services SET 
    active = 1
WHERE id = ?
RETURNING *;

-- name: UpdteHostServiceSchedule :one
UPDATE host_services SET 
    schedule_number = ?,
    schedule_unit = ?
WHERE id = ?
RETURNING *;

-- name: GetInnactiveHostServiceByHostAndService :one
SELECT hs.*, h.host_name, s.service_name
FROM host_services as hs
JOIN hosts as h ON hs.host_id = h.id
JOIN services as s ON hs.service_id = s.id
WHERE hs.host_id = ?
AND hs.service_id = ?
AND hs.active = 0;

-- name: GetHostsServicesToMonitor :many
SELECT hs.*, h.host_name, s.service_name
FROM host_services as hs
JOIN hosts as h ON hs.host_id = h.id
JOIN services as s ON hs.service_id = s.id
Where hs.active = 1;
