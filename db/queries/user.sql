-- name: CreateUser :one
INSERT INTO users (first_name, last_name, email, password_hash, role_id, is_active, nin, gender, date_of_birth)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
    first_name = COALESCE($2, first_name),
    last_name = COALESCE($3, last_name),
    email = COALESCE($4, email),
    role_id = COALESCE($5, role_id),
    is_active = COALESCE($6, is_active),
    nin = COALESCE($7, nin),
    gender = COALESCE($8, gender),
    date_of_birth = COALESCE($9, date_of_birth),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: GetUserByID :one
SELECT u.*, r.name as role_name FROM users u
JOIN roles r ON u.role_id = r.id
WHERE u.id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT u.*, r.name as role_name FROM users u
JOIN roles r ON u.role_id = r.id
ORDER BY u.created_at DESC;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: CreateRole :one
INSERT INTO roles (name, description)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateRole :one
UPDATE roles
SET
    name = COALESCE($2, name),
    description = COALESCE($3, description)
WHERE id = $1
RETURNING *;

-- name: DeleteRole :exec
DELETE FROM roles WHERE id = $1;

-- name: GetRoleByID :one
SELECT * FROM roles WHERE id = $1 LIMIT 1;

-- name: ListRoles :many
SELECT * FROM roles ORDER BY name;

-- name: AddPermissionToRole :exec
INSERT INTO role_permissions (role_id, permission_id)
VALUES ($1, $2);

-- name: RemovePermissionFromRole :exec
DELETE FROM role_permissions
WHERE role_id = $1 AND permission_id = $2;

-- name: GetRolePermissions :many
SELECT p.* FROM permissions p
JOIN role_permissions rp ON p.id = rp.permission_id
WHERE rp.role_id = $1;

-- name: LogUserActivity :one
INSERT INTO user_activity_logs (user_id, action, description, ip_address, user_agent)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: LogLoginAttempt :one
INSERT INTO login_history (user_id, ip_address, user_agent, success)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserLoginHistory :many
SELECT * FROM login_history
WHERE user_id = $1
ORDER BY login_time DESC
LIMIT $2;

-- name: CreatePasswordResetToken :one
INSERT INTO password_reset_tokens (user_id, token, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetPasswordResetToken :one
SELECT * FROM password_reset_tokens
WHERE token = $1 AND expires_at > NOW() AND used = FALSE
LIMIT 1;

-- name: MarkTokenAsUsed :exec
UPDATE password_reset_tokens
SET used = TRUE
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT
    u.id,
    u.first_name,
    u.last_name,
    u.email,
    u.password_hash,
    u.is_active,
    r.name as role_name
FROM users u
JOIN roles r ON u.role_id = r.id
WHERE u.email = $1 LIMIT 1;

-- name: GetUserPermissions :many
SELECT p.code
FROM permissions p
JOIN role_permissions rp ON p.id = rp.permission_id
JOIN roles r ON rp.role_id = r.id
JOIN users u ON u.role_id = r.id
WHERE u.id = $1;
