-- name: CreateBusiness :one
INSERT INTO business (
    name, address_one, addres_two, country, phone, email, website, city, state, 
    zip_code, tax_id, tax_rate, logo_url, rounding, currency, timezone, language,
    low_stock_threshold, allow_overselling, payment_type, font, primary_color, motto
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9,
    $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23
) RETURNING *;

-- name: GetBusiness :one
SELECT * FROM business WHERE id = $1;

-- name: ListBusinesses :many
SELECT * FROM business ORDER BY created_at DESC;

-- name: UpdateBusiness :one
UPDATE business SET
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
    tax_id = $12,
    tax_rate = $13,
    logo_url = $14,
    rounding = $15,
    currency = $16,
    timezone = $17,
    language = $18,
    low_stock_threshold = $19,
    allow_overselling = $20,
    payment_type = $21,
    font = $22,
    primary_color = $23,
    motto = $24,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteBusiness :exec
DELETE FROM business WHERE id = $1;
