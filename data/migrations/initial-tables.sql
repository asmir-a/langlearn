CREATE TABLE IF NOT EXISTS users (
    username TEXT,
    password_hash TEXT,
    password_salt TEXT,
    PRIMARY KEY (username)
);

CREATE TABLE IF NOT EXISTS sessions (
    session_key TEXT,
    username TEXT,
    login_time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT Now(),
    last_seen_time TIMESTAMP WITHOUT TIME ZONE,
    PRIMARY KEY (session_key),
    CONSTRAINT fk_username
        FOREIGN KEY (username)
        REFERENCES users(username)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS korean_words (
    index INTEGER GENERATED ALWAYS AS IDENTITY,
    word TEXT,
    part_of_speech TEXT,
    defs TEXT[],
    freq_rank NUMERIC,
    PRIMARY KEY (index)
);

CREATE TABLE IF NOT EXISTS knows (
    username TEXT,
    word TEXT,
    current_count NUMERIC,
    CONSTRAINT fk_username
        FOREIGN KEY (username)
        REFERENCES users(username)
        ON DELETE CASCADE
);