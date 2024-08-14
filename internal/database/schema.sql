CREATE TABLE hosts (
    id INTEGER PRIMARY KEY,
    host_name VARCHAR(255) NOT NULL,
    canonical_name VARCHAR(255),
    url VARCHAR(255),
    ip VARCHAR(255),
    ipv6 VARCHAR(255),
    location VARCHAR(255),
    os VARCHAR(255),
    active INTEGER DEFAULT 1,
    last_updated DATETIME
);

CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    user_active INTEGER DEFAULT 0,
    access_level INTEGER DEFAULT 3,
    email VARCHAR(255),
    password VARCHAR(60),
    deleted_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
