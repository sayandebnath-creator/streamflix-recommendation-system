CREATE INDEX IF NOT EXISTS movie_embedding_idx
ON movies
USING ivfflat (embedding vector_cosine_ops)
WITH (lists = 100);