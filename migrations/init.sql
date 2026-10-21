BEGIN;

CREATE TABLE quiestion(
    id INTEGER SERIAL PRIMARY,
    text TEXT NOT NULL,
    options TEXT[] NOT NULL,
    correct_option INTEGER NOT NULL,
    time_limit INTEGER NOT NULL
);

CREATE TABLE quiz(
    id INTEGER SERIAL PRIMARY,
    title TEXT NOT NULL,
    quiestions_id  INTEGER NOT NULL REFERENCES quiestion (id) ON DELETE CASCADE,
);

COMMIT;