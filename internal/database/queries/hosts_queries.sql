-- name: GetHosts :many
SELECT * FROM hosts 
Limit ?
Offset ?;
