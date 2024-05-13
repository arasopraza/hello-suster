CREATE TABLE user_it (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    nip BIGINT NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)