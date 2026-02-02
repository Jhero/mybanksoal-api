ALTER TABLE users ADD COLUMN api_key VARCHAR(255) UNIQUE DEFAULT NULL;

-- Seed mobile user
INSERT INTO users (username, password, role, api_key, created_at, updated_at)
VALUES (
    'mobile_client', 
    '$2a$10$NotUsedButRequiredForConstraint................', 
    'mobile_reader', 
    'mobile-app-secure-api-key-2026', 
    NOW(), 
    NOW()
);
