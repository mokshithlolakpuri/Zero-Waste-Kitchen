-- Create users table first
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    fcm_token VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create receipts table second
CREATE TABLE receipts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    image_path VARCHAR(255) NOT NULL,
    total_amount DECIMAL(10, 2),
    store_name VARCHAR(255),
    purchase_date TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create grocery_items table last (with foreign key to receipts)
CREATE TABLE grocery_items (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receipt_id INTEGER REFERENCES receipts(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    quantity DECIMAL(10, 2) NOT NULL,
    unit VARCHAR(50),
    barcode VARCHAR(255),
    batch_number VARCHAR(255),
    manufacture_date TIMESTAMP,
    expiry_date TIMESTAMP,
    storage_location VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- Create recipes table
CREATE TABLE recipes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    ingredients TEXT NOT NULL, -- Each ingredient on a new line
    instructions TEXT NOT NULL, -- Step-by-step instructions
    prep_time INTEGER NOT NULL, -- In minutes
    cook_time INTEGER NOT NULL, -- In minutes
    servings INTEGER NOT NULL,
    difficulty VARCHAR(50) NOT NULL, -- Easy, Medium, or Hard
    cuisine VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for recipes table
CREATE INDEX idx_recipes_user_id ON recipes(user_id);

-- Create indexes
CREATE INDEX idx_grocery_items_user_expiry ON grocery_items(user_id, expiry_date);