-- name: GetHosts :one
SELECT * FROM hosts 
Limit ?
Offset ?;
