-- Brand
-- name: CreateBrand :one
INSERT INTO brand (name, description, logo)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetBrand :one
SELECT * FROM brand WHERE id = $1 LIMIT 1;

-- name: ListBrands :many
SELECT * FROM brand ORDER BY name;

-- name: UpdateBrand :one
UPDATE brand
SET name = $2,
    description = $3,
    logo = $4,
    is_active = $5,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteBrand :exec
DELETE FROM brand WHERE id = $1;


-- Category
-- name: CreateCategory :one
INSERT INTO category (name, parent_id, description)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetCategory :one
SELECT * FROM category WHERE id = $1 LIMIT 1;

-- name: ListCategories :many
SELECT * FROM category ORDER BY name;

-- name: UpdateCategory :one
UPDATE category
SET name = $2,
    parent_id = $3,
    description = $4,
    is_active = $5,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM category WHERE id = $1;


-- Item
-- name: CreateItem :one
INSERT INTO item (brand_id, category_id, name, description)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetItem :one
SELECT * FROM item WHERE id = $1 LIMIT 1;

-- name: ListItems :many
SELECT * FROM item ORDER BY name;

-- name: ListItemsByCategory :many
SELECT * FROM item WHERE category_id = $1 ORDER BY name;

-- name: UpdateItem :one
UPDATE item
SET brand_id = $2,
    category_id = $3,
    name = $4,
    description = $5,
    is_active = $6,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteItem :exec
DELETE FROM item WHERE id = $1;


-- Variation
-- name: CreateVariation :one
INSERT INTO variation (item_id, sku, name, unit, size, color, barcode, price)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetVariation :one
SELECT * FROM variation WHERE id = $1 LIMIT 1;

-- name: ListVariationsByItem :many
SELECT * FROM variation WHERE item_id = $1 ORDER BY name;

-- name: UpdateVariation :one
UPDATE variation
SET sku = $2,
    name = $3,
    unit = $4,
    size = $5,
    color = $6,
    barcode = $7,
    price = $8,
    is_active = $9,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteVariation :exec
DELETE FROM variation WHERE id = $1;


-- Image
-- name: CreateItemImage :one
INSERT INTO item_image (item_id, variation_id, url, is_primary)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetItemImagesByItem :many
SELECT * FROM item_image WHERE item_id = $1;

-- name: GetItemImagesByVariation :many
SELECT * FROM item_image WHERE variation_id = $1;

-- name: DeleteItemImage :exec
DELETE FROM item_image WHERE id = $1;


-- Inventory
-- name: UpsertInventory :one
INSERT INTO inventory (store_id, variation_id, quantity, reorder_level, max_level)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (store_id, variation_id)
DO UPDATE SET
    quantity = EXCLUDED.quantity,
    reorder_level = EXCLUDED.reorder_level,
    max_level = EXCLUDED.max_level,
    last_updated = NOW()
RETURNING *;

-- name: GetInventoryByStore :many
SELECT * FROM inventory WHERE store_id = $1;

-- name: GetInventoryItem :one
SELECT * FROM inventory
WHERE store_id = $1 AND variation_id = $2
LIMIT 1;

-- name: UpdateInventoryQuantity :one
UPDATE inventory
SET quantity = $3,
    last_updated = NOW()
WHERE store_id = $1 AND variation_id = $2
RETURNING *;

-- name: DeleteInventory :exec
DELETE FROM inventory WHERE id = $1;

-- Units
-- name: CreateUnit :one
INSERT INTO unit (name, short_code)
VALUES ($1, $2)
RETURNING *;

-- name: GetUnitByID :one
SELECT * FROM unit
WHERE id = $1;

-- name: ListUnits :many
SELECT * FROM unit
ORDER BY id;

-- name: UpdateUnit :one
UPDATE unit
SET name = $1,
    short_code = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $3
RETURNING *;

-- name: DeleteUnit :exec
DELETE FROM unit
WHERE id = $1
RETURNING *;


-- Color
-- name: CreateColor :one
INSERT INTO color (name)
VALUES ($1)
RETURNING *;

-- name: GetColorByID :one
SELECT * FROM color
WHERE id = $1;

-- name: GetColorByName :one
SELECT * FROM color
WHERE name = $1;

-- name: ListColors :many
SELECT * FROM color
ORDER BY id;

-- name: UpdateColor :one
UPDATE color
SET name = $1,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $2
RETURNING *;

-- name: DeleteColor :exec
DELETE FROM color
WHERE id = $1
RETURNING *;
