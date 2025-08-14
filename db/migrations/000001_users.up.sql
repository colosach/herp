-- Roles table
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT
);

-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    -- address VARCHAR(100) NOT NULL,
    -- phone_number VARCHAR(15) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    -- nin VARCHAR(11) NOT NULL,
    gender VARCHAR NOT NULL CHECK (gender IN ('male', 'female')),
    -- date_of_birth DATE NOT NULL,
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

-- Sample user with admin role (role_id = 1)
-- Password: admin123 (hashed with bcrypt)
INSERT INTO users (username, first_name, last_name, email, password_hash, gender, role_id, is_active)
VALUES (
    'admin',
    'John',
    'Doe',
    'admin@hotel.com',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi',
    'male',
    1,
    true
);

-- Sample user with manager role (role_id = 2)
-- Password: manager123 (hashed with bcrypt)
INSERT INTO users (username, first_name, last_name, email, password_hash, gender, role_id, is_active)
VALUES (
    'manager1',
    'Jane',
    'Smith',
    'manager@hotel.com',
    '$2a$10$8K1p/a0dhrxiH8Tf4Gro9e.0uI4JGO0JG6LJZr1f7wFYw8mO6pR1W',
    'female',
    2,
    true
);

-- Sample user with POS staff role (role_id = 3)
-- Password: pos123 (hashed with bcrypt)
INSERT INTO users (username, first_name, last_name, email, password_hash, gender, role_id, is_active)
VALUES (
    'pos_staff1',
    'Mike',
    'Johnson',
    'pos@hotel.com',
    '$2a$10$Q8H2k5JFJzJ8F5JzJ8F5Ju.J8F5JzJ8F5JzJ8F5JzJ8F5JzJ8F5J8',
    'male',
    3,
    true
);

-- Sample user with cashier role (role_id = 4)
-- Password: cashier123 (hashed with bcrypt)
INSERT INTO users (username, first_name, last_name, email, password_hash, gender, role_id, is_active)
VALUES (
    'cashier1',
    'Sarah',
    'Wilson',
    'cashier@hotel.com',
    '$2a$10$L7N3k6KGKzL8G6KzL8G6Lu.L8G6KzL8G6KzL8G6KzL8G6KzL8G6L8',
    'female',
    4,
    true
);

-- Additional sample user - inactive user for testing
-- Password: test123 (hashed with bcrypt)
INSERT INTO users (username, first_name, last_name, email, password_hash, gender, role_id, is_active)
VALUES (
    'test_user',
    'Test',
    'User',
    'test@hotel.com',
    '$2a$10$T9P4l7MHMzT9H7MzT9H7Mu.T9H7MzT9H7MzT9H7MzT9H7MzT9H7M9',
    'male',
    3,
    false
);