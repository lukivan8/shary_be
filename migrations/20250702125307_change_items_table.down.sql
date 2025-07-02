-- Revert changes to items table

DROP TABLE IF EXISTS items CASCADE;

CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    image_url TEXT NOT NULL,
    price_per_day DECIMAL(10,2) NOT NULL CHECK (price_per_day > 0),
    location VARCHAR(500) NOT NULL,
    is_available BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);