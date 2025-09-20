
-- How this hierarchy works:

-- Brand: Coca-Cola, Johnnie Walker, Samsung.

-- Category: Drinks → Soft Drinks → Cola.

-- Item: Coca-Cola.

-- Variation: Coca-Cola 500ml Bottle, Coca-Cola 1L Bottle.

-- Inventory: “Branch 1 Central Store has 50 of Coca-Cola 500ml.”

CREATE TABLE brand (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    logo VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE category (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    parent_id INT REFERENCES category(id) ON DELETE SET NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (name, parent_id) -- avoids duplicate names under same parent
);


CREATE TABLE item (
    id SERIAL PRIMARY KEY,
    brand_id INT REFERENCES brand(id) ON DELETE SET NULL,
    category_id INT REFERENCES category(id) ON DELETE SET NULL,
    name VARCHAR(150) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Variations represent different versions of an item (e.g. size, color)
CREATE TABLE variation (
    id SERIAL PRIMARY KEY,
    item_id INT NOT NULL REFERENCES item(id) ON DELETE CASCADE,
    sku VARCHAR(50) NOT NULL UNIQUE, -- stock keeping unit
    name VARCHAR(100)NOT NULL, -- e.g. '500ml Bottle', '1kg Pack'
    unit VARCHAR(20) NOT NULL, -- e.g. 'ml', 'kg', 'pcs'
    size VARCHAR(50), -- optional, e.g. '500', 'Large'
    color VARCHAR(30), -- optional, e.g. 'Red', 'Blue'
    barcode VARCHAR(50) UNIQUE,
    price NUMERIC(12,2),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE item_image (
    id SERIAL PRIMARY KEY,
    item_id INT REFERENCES item(id) ON DELETE CASCADE,
    variation_id INT REFERENCES variation(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT item_or_variation CHECK (
        (item_id IS NOT NULL AND variation_id IS NULL)
        OR (item_id IS NULL AND variation_id IS NOT NULL)
    )
);


-- Inventory table to track stock levels in each store
-- Each record represents the stock of a specific variation in a specific store
CREATE TABLE inventory (
    id SERIAL PRIMARY KEY,
    store_id INT NOT NULL REFERENCES store(id) ON DELETE CASCADE,
    variation_id INT NOT NULL REFERENCES variation(id) ON DELETE CASCADE,
    quantity INT NOT NULL DEFAULT 0,
    reorder_level INT DEFAULT 0,
    max_level INT,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (store_id, variation_id)
);


