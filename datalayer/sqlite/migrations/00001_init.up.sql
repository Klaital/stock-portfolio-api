
CREATE TABLE users (
    user_id INTEGER PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    email VARCHAR(256) UNIQUE NOT NULL,
    password_digest VARCHAR(256) NOT NULL
);

CREATE UNIQUE INDEX users_email ON users(email);

CREATE TABLE positions (
    position_id INTEGER PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id INTEGER NOT NULL,
    symbol VARCHAR(8) NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE INDEX user_positions ON positions(user_id);


