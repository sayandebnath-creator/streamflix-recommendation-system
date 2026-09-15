ALTER TABLE movies
ADD COLUMN poster_valid BOOLEAN;

CREATE INDEX idx_movies_poster_valid
ON movies (poster_valid);