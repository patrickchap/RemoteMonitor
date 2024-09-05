-- name: GetUserForAuth :one
SELECT email, password, id 
FROM users 
WHERE email = ?;

-- name: GetUserByEmail :one
SELECT * 
FROM users
WHERE email = ?;
