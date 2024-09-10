-- name: GetServices :one
SELECT * FROM services
WHERE active = 1;
