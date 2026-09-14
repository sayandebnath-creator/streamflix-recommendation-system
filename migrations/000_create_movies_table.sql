CREATE TABLE IF NOT EXISTS movies (
    id UUID PRIMARY KEY,
    tmdb_id INTEGER NOT NULL UNIQUE,
    title TEXT NOT NULL,
    overview TEXT,
    genres TEXT,
    language VARCHAR(10),
    release_date DATE,
    runtime INTEGER,
    popularity DOUBLE PRECISION,
    vote_average DOUBLE PRECISION,
    vote_count INTEGER,
    poster_path TEXT,
    embedding vector(384),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_movies_tmdb_id
ON movies (tmdb_id);
