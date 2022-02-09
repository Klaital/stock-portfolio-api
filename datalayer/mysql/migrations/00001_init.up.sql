
CREATE TABLE users (
    user_id INTEGER PRIMARY KEY AUTO_INCREMENT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    email VARCHAR(256) UNIQUE NOT NULL UNIQUE,
    password_digest VARCHAR(256) NOT NULL
);


CREATE INDEX users_email ON users(email);

CREATE TABLE positions (
                           position_id INTEGER PRIMARY KEY AUTO_INCREMENT,
                           created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                           updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
                           user_id INTEGER NOT NULL,
                           symbol VARCHAR(8) NOT NULL,
                           bought_at DATETIME,
                           basis INTEGER NOT NULL DEFAULT 0,

                           FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE INDEX user_positions ON positions(user_id);
CREATE INDEX user_position_sym ON positions(user_id, symbol);

