CREATE TABLE companies (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE users (
    id   BIGSERIAL PRIMARY KEY,
    name text      NOT NULL,
    email text UNIQUE NOT NULL,
    bio  text
);