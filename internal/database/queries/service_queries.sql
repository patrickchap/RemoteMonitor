-- name: GetServices :many
SELECT * FROM services
WHERE active = 1;
