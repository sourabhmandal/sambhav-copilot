CREATE TABLE translations (
    id SERIAL PRIMARY KEY,
    normalized_hash CHAR(64) NOT NULL UNIQUE,
    language_code VARCHAR(10) NOT NULL,
    original_text TEXT NOT NULL,
    translated_text TEXT NOT NULL,
    confidence_score NUMERIC(4,3),
    provider VARCHAR(50),
    created_at TIMESTAMP DEFAULT now(),
    UNIQUE(normalized_hash, language_code)
);

CREATE INDEX idx_translations_lookup 
ON translations(normalized_hash, language_code);

CREATE TABLE users (
    id   BIGSERIAL PRIMARY KEY,
    name text      NOT NULL,
    email text UNIQUE NOT NULL,
    bio  text
);