-- name: InsertUser :one
INSERT INTO users (name, email, password)
VALUES ($1, $2, $3)
RETURNING id;


-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;


-- name: GetUserWithPermissionsById :one
SELECT 
    u.id,
    u.name,
    u.email,
    r.name as role_name
FROM users u
LEFT JOIN roles r ON u.role_id = r.id
WHERE u.id = $1;


-- name: ListUser :many
SELECT 
    u.name,
    u.email,
    r.name as role_name
FROM users u
LEFT JOIN roles r on u.role_id = r.id;
