-- name: CreateBusiness :one
INSERT INTO business (
    owner_id, name, motto, email, website, tax_id, tax_rate,
    country, logo_url, rounding, currency, timezone, language,
    low_stock_threshold, allow_overselling, payment_type,
    font, primary_color
) VALUES (
    $1, $2, $3, $4, $5, $6, $7,
    $8, $9, $10, $11, $12, $13,
    $14, $15, $16, $17, $18
) RETURNING *;

-- name: GetBusiness :one
SELECT *
FROM business
WHERE id = $1 AND owner_id = $2;

-- name: ListBusinesses :many
SELECT *
FROM business
WHERE owner_id = $1
ORDER BY created_at;

-- name: UpdateBusiness :one
UPDATE business SET
    name = COALESCE(sqlc.narg(name), name),
    motto = COALESCE(sqlc.narg(motto), motto),
    email = COALESCE(sqlc.narg(email), email),
    website = COALESCE(sqlc.narg(website), website),
    tax_id = COALESCE(sqlc.narg(tax_id), tax_id),
    tax_rate = COALESCE(sqlc.narg(tax_rate), tax_rate),
    logo_url = COALESCE(sqlc.narg(logo_url), logo_url),
    rounding = COALESCE(sqlc.narg(rounding), rounding),
    currency = COALESCE(sqlc.narg(currency), currency),
    timezone = COALESCE(sqlc.narg(timezone), timezone),
    language = COALESCE(sqlc.narg(language), language),
    low_stock_threshold = COALESCE(sqlc.narg(low_stock_threshold), low_stock_threshold),
    allow_overselling = COALESCE(sqlc.narg(allow_overselling), allow_overselling),
    payment_type = COALESCE(sqlc.narg(payment_type), payment_type),
    font = COALESCE(sqlc.narg(font), font),
    primary_color = COALESCE(sqlc.narg(primary_color), primary_color),
    country = COALESCE(sqlc.narg(country), country),
    updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg(id) AND owner_id = sqlc.arg(owner_id)
RETURNING *;

-- name: DeleteBusiness :one
DELETE FROM business
WHERE id = $1 AND owner_id = $2
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
