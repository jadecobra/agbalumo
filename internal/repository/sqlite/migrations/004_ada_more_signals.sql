-- Add more Ada-centric signals for the scraper
ALTER TABLE listings ADD COLUMN payment_methods TEXT;
ALTER TABLE listings ADD COLUMN menu_url TEXT;
