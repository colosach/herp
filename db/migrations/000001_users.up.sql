-- Roles table
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT
);

-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    -- address VARCHAR(100) NOT NULL,
    -- phone_number VARCHAR(15) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    nin VARCHAR(11) NOT NULL,
    gender VARCHAR NOT NULL CHECK (gender IN ('male', 'female')),
    date_of_birth DATE NOT NULL,
    role_id INTEGER NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id)
);

-- Permissions table
CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    description TEXT
);

-- Role-Permissions junction table
CREATE TABLE role_permissions (
    role_id INTEGER NOT NULL,
    permission_id INTEGER NOT NULL,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(id),
    FOREIGN KEY (permission_id) REFERENCES permissions(id)
);

-- Sample data
INSERT INTO roles (name, description) VALUES
('admin', 'System administrator with full access'),
('manager', 'Hotel manager with broad access'),
('pos_staff', 'POS system user'),
('cashier', 'Cashier with limited POS access');

INSERT INTO permissions (code, description) VALUES
('pos:sell', 'Create new sales in POS'),
('pos:view', 'View sales history in POS'),
('pos:manage_items', 'Manage POS items'),
('booking:create', 'Create new bookings'),
('booking:manage', 'Manage all bookings');

-- Assign permissions to roles
INSERT INTO role_permissions (role_id, permission_id) VALUES
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5), -- admin has all permissions
(2, 1), (2, 2), (2, 3), (2, 4),         -- manager has all POS and booking permissions
(3, 1), (3, 2),                          -- pos_staff can sell and view
(4, 1);                                  -- cashier can only sell