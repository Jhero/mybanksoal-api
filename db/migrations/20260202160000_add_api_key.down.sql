DELETE FROM users WHERE username = 'mobile_client';
ALTER TABLE users DROP COLUMN api_key;
