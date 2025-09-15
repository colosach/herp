-- name: CreateBusiness :one
INSERT INTO business (
    name, motto, email, website, tax_id, tax_rate, logo_url, rounding, currency, timezone, language,
    low_stock_threshold, allow_overselling, payment_type, font, primary_color, country
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
) RETURNING *;

-- name: GetBusiness :one
SELECT * FROM business WHERE id = $1;

-- name: ListBusinesses :many
SELECT * FROM business ORDER BY created_at DESC;

-- name: UpdateBusiness :one
UPDATE business SET
    name = $2,
    motto = $3,
    email = $4,
    website = $5,
    tax_id = $6,
    tax_rate = $7,
    logo_url = $8,
    rounding = $9,
    currency = $10,
    timezone = $11,
    language = $12,
    low_stock_threshold = $13,
    allow_overselling = $14,
    payment_type = $15,
    font = $16,
    primary_color = $17,
    country = $18,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteBusiness :one
DELETE FROM business WHERE id = $1
RETURNING *;

-- name: CreateBranch :one
INSERT INTO branch (
    business_id, name, address_one, addres_two, country, phone, email, website, city, state, zip_code
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetBranch :one
SELECT * FROM branch WHERE id = $1;

-- name: ListBranches :many
SELECT * FROM branch ORDER BY created_at DESC;

-- name: DeleteBranch :one
DELETE FROM branch WHERE id = $1
RETURNING *;

-- name: UpdateBranch :one
UPDATE branch SET
    name = $2,
    address_one = $3,
    addres_two = $4,
    country = $5,
    phone = $6,
    email = $7,
    website = $8,
    city = $9,
    state = $10,
    zip_code = $11,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;
