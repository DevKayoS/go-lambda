-- name: InsertUser :one
INSERT INTO users (name, email, password)
VALUES ($1, $2, $3)
RETURNING id;


-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;
