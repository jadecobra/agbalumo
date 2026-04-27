-- delivery platforms migration (008)
ALTER TABLE listings ADD COLUMN delivery_platforms TEXT DEFAULT '';
