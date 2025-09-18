-- name: CreateStore :one
INSERT INTO store (
    name, description, branch_id, address, phone, email, 
    is_active, store_code, store_type, assigned_user, manager_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING *;

-- name: GetStoreByID :one
SELECT * FROM store WHERE id = $1 LIMIT 1;

-- name: GetStoresByBranch :many
SELECT * FROM store WHERE branch_id = $1 ORDER BY name;

-- name: GetCentralStoreByBranch :one
SELECT * FROM store WHERE branch_id = $1 AND store_type = 'central' LIMIT 1;

-- name: ListStores :many
SELECT * FROM store ORDER BY created_at DESC;

-- name: UpdateStore :one
UPDATE store
SET name = $2,
    description = $3,
    address = $4,
    phone = $5,
    email = $6,
    is_active = $7,
    store_code = $8,
    store_type = $9,
    assigned_user = $10,
    manager_id = $11,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteStore :exec
DELETE FROM store WHERE id = $1;

-- name: DeactivateStore :one
UPDATE store
SET is_active = FALSE,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: SearchStoresByName :many
SELECT * FROM store
WHERE name ILIKE '%' || $1 || '%'
ORDER BY name;
