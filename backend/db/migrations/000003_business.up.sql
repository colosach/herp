

CREATE TYPE payment_type AS ENUM ('cash', 'pos', 'room_charge', 'transfer');

-- Create the business table
CREATE TABLE business (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    motto VARCHAR(50),
    email VARCHAR(50) UNIQUE,
    website VARCHAR(20),
    tax_id VARCHAR(100),
    tax_rate DECIMAL(5,2) DEFAULT 0.00,
    country VARCHAR(100) NOT NULL,
    logo_url VARCHAR(255),
    rounding VARCHAR(50) DEFAULT 'nearest',
    currency VARCHAR(10) DEFAULT 'NGN',
    timezone VARCHAR(10) DEFAULT 'UTC',
    language VARCHAR(50) DEFAULT 'en',
    low_stock_threshold INT DEFAULT 5,
    allow_overselling BOOLEAN DEFAULT FALSE,
    payment_type payment_type[] DEFAULT ARRAY['cash']::payment_type[], 
    font VARCHAR(100) DEFAULT 'Arial',
    primary_color VARCHAR(7) DEFAULT '#000000',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_business_name ON business (name);
CREATE INDEX idx_business_email ON business (email);
CREATE INDEX idx_business_tax_id ON business (tax_id);


CREATE TABLE branch (
    id SERIAL PRIMARY KEY,
    business_id INTEGER NOT NULL,
    name VARCHAR(50) NOT NULL,
    address_one VARCHAR(20),
    addres_two VARCHAR(20),
    country VARCHAR(10),
    phone VARCHAR(15),
    email VARCHAR(20),
    website VARCHAR(20),
    city VARCHAR(50),
    state VARCHAR(10),
    zip_code VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (business_id) REFERENCES business(id) ON DELETE CASCADE
);

CREATE INDEX idx_branch_business_id ON branch (business_id);
CREATE INDEX idx_branch_name ON branch (name);
CREATE INDEX idx_branch_city_state ON branch (city, state);