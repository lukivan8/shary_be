-- Create items table
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

-- Create index on title for search functionality
CREATE INDEX IF NOT EXISTS idx_items_title ON items USING gin(to_tsvector('english', title));

-- Create index on description for search functionality
CREATE INDEX IF NOT EXISTS idx_items_description ON items USING gin(to_tsvector('english', description));

-- Create index on location for location-based queries
CREATE INDEX IF NOT EXISTS idx_items_location ON items(location);

-- Create index on price_per_day for price filtering
CREATE INDEX IF NOT EXISTS idx_items_price_per_day ON items(price_per_day);

-- Create index on is_available for availability filtering
CREATE INDEX IF NOT EXISTS idx_items_is_available ON items(is_available);

-- Create index on created_at for sorting
CREATE INDEX IF NOT EXISTS idx_items_created_at ON items(created_at DESC);

-- Create a combined search index
CREATE INDEX IF NOT EXISTS idx_items_search ON items USING gin(to_tsvector('english', title || ' ' || description)); 