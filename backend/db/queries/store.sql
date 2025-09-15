-- name: CreateStore :one
INSERT INTO store (
    name, description, branch_id, address, phone, email, is_active, store_code
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetStoreByID :one
SELECT * FROM store
WHERE id = $1;

-- name: GetStoreByCode :one
SELECT * FROM store
WHERE store_code = $1;

-- name: ListStores :many
SELECT * FROM store
WHERE is_active = COALESCE($1, is_active)
ORDER BY name
LIMIT $2 OFFSET $3;

-- name: UpdateStore :one
UPDATE store
SET 
    name = $1,
    description = $2,
    branch_id = $3,
    address = $4,
    phone = $5,
    email = $6,
    is_active = $7,
    store_code = $8,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $9
RETURNING *;

-- name: DeleteStore :exec
UPDATE store
SET is_active = FALSE,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: CreateStoreManager :one
INSERT INTO store_manager (
    store_id, user_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetStoreManagerByID :one
SELECT * FROM store_manager
WHERE id = $1;

-- name: ListManagersByStore :many
SELECT sm.*, u.username
FROM store_manager sm
JOIN users u ON sm.user_id = u.id
WHERE sm.store_id = $1
ORDER BY sm.assigned_at
LIMIT $2 OFFSET $3;

-- name: DeleteStoreManager :exec
DELETE FROM store_manager
WHERE id = $1;