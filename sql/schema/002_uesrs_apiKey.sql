-- +goose Up
ALTER TABLE users 
ADD COLUMN api_key VARCHAR(64);
UPDATE users SET api_key = encode(sha256(random()::text::bytea), 'hex');
ALTER TABLE users ALTER COLUMN api_key SET NOT NULL;
ALTER TABLE users ADD CONSTRAINT unique_api_key UNIQUE(api_key);

-- +goose Down
ALTER TABLE users
DROP COLUMN api_key;