CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'buyer', -- 'buyer', 'seller', 'admin'
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE oauth_clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id VARCHAR(255) UNIQUE NOT NULL,
    client_secret VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL, -- 'Mobile App', 'Web App', 'Admin Portal'
    redirect_uris TEXT[], -- Array of allowed redirect URIs
    grant_types TEXT[] DEFAULT ARRAY['authorization_code', 'refresh_token'],
    scope VARCHAR(255) DEFAULT 'read write',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE authorization_codes (
    code VARCHAR(255) PRIMARY KEY,
    client_id VARCHAR(255) REFERENCES oauth_clients(client_id),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    redirect_uri TEXT,
    scope VARCHAR(255),
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE access_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token VARCHAR(500) UNIQUE NOT NULL,
    client_id VARCHAR(255) REFERENCES oauth_clients(client_id),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    scope VARCHAR(255),
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token VARCHAR(500) UNIQUE NOT NULL,
    client_id VARCHAR(255) REFERENCES oauth_clients(client_id),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    access_token_id UUID REFERENCES access_tokens(id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_auth_codes_expires ON authorization_codes(expires_at);
CREATE INDEX idx_access_tokens_user ON access_tokens(user_id);
CREATE INDEX idx_access_tokens_expires ON access_tokens(expires_at);
CREATE INDEX idx_refresh_tokens_user ON refresh_tokens(user_id);

