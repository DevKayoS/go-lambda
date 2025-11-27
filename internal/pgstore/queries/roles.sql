-- name: GetUserWithPermissions :one
SELECT 
    u.id,
    u.name,
    u.email,
    u.password,
    u.role_id,
    r.name as role_name
FROM users u
LEFT JOIN roles r ON u.role_id = r.id
WHERE u.email = $1;

-- name: GetUserPermissions :many
SELECT p.name
FROM users u
JOIN roles r ON u.role_id = r.id
JOIN role_permissions rp ON r.id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
WHERE u.email = $1;

-- name: GetRoleByID :one
SELECT * FROM roles WHERE id = $1;

-- name: GetRoleByName :one
SELECT * FROM roles WHERE name = $1;

-- name: GetPermissionsByRoleID :many
SELECT p.* FROM permissions p
JOIN role_permissions rp ON p.id = rp.permission_id
WHERE rp.role_id = $1;
