-- This document is licensed under the Creative Commons Attribution-ShareAlike 4.0 International License.
-- To view a copy of this license, visit http://creativecommons.org/licenses/by-sa/4.0/ or send a letter to Creative Commons, PO Box 1866, Mountain View, CA 94042, USA.
-- You are free to share and adapt this material for any purpose, even commercially, under the following terms:
-- Attribution: You must give appropriate credit, provide a link to the license, and indicate if changes were made. You may do so in any reasonable manner, but not in any way that suggests the licensor endorses you or your use.
-- ShareAlike: If you remix, transform, or build upon the material, you must distribute your contributions under the same license as the original.  


CREATE TYPE payment_type AS ENUM ('cash', 'pos', 'room_charge', 'transfer');

-- Create the business table
CREATE TABLE business (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    motto VARCHAR(255),
    address_one VARCHAR(255) NOT NULL,
    addres_two VARCHAR(255),
    country VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    email VARCHAR(255) UNIQUE,
    website VARCHAR(255),
    city VARCHAR(100),
    state VARCHAR(100),
    zip_code VARCHAR(20),   
    tax_id VARCHAR(100),
    tax_rate DECIMAL(5,2) DEFAULT 0.00,
    logo_url VARCHAR(255),
    rounding VARCHAR(50) DEFAULT 'nearest',
    currency VARCHAR(10) DEFAULT 'NGN',
    timezone VARCHAR(100) DEFAULT 'UTC',
    language VARCHAR(50) DEFAULT 'en',
    low_stock_threshold INT DEFAULT 5,
    allow_overselling BOOLEAN DEFAULT FALSE,
    payment_type payment_type[] DEFAULT ARRAY['cash']::payment_type[], 
    font VARCHAR(100) DEFAULT 'Arial',
    primary_color VARCHAR(7) DEFAULT '#000000',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
