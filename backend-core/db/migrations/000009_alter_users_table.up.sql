ALTER TABLE IF EXISTS users

ADD COLUMN is_verified BOOLEAN DEFAULT false;
