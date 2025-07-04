-- Add back image_url column to items table
ALTER TABLE items ADD COLUMN image_url TEXT NOT NULL DEFAULT ''; 