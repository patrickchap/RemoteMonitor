-- name: GetUserForAuth :one
SELECT email, password, id 
FROM users 
WHERE email = ?;
