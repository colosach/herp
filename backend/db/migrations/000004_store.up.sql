CREATE TABLE store (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(255),
    branch_id INTEGER NOT NULL,
    address VARCHAR(50),
    phone VARCHAR(15),
    email VARCHAR(50),
    is_active BOOLEAN DEFAULT TRUE,
    store_code VARCHAR(10) NOT NUll,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (branch_id) REFERENCES branch(id) ON DELETE CASCADE
);

CREATE TABLE store_manager (
    id SERIAL PRIMARY KEY,
    store_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (store_id) REFERENCES store(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);



CREATE INDEX idx_store_branch_id ON store (branch_id);
CREATE INDEX idx_store_name ON store (name);
CREATE INDEX idx_store_manager_store_id ON store_manager (store_id);
CREATE INDEX idx_store_manager_user_id ON store_manager (user_id);