package ingestion

import (
	"strconv"
	"time"

	"streamflix-backend/internal/movie"
)

func ParseMovie(record []string) (*movie.Movie, error) {
	m := &movie.Movie{}

	// TMDB ID
	if id, err := strconv.Atoi(record[5]); err == nil {
		m.TMDBID = id
	}

	m.Language = record[7]
	m.Overview = record[9]

	// Popularity
	if p, err := strconv.ParseFloat(record[10], 64); err == nil {
		m.Popularity = p
	}

	m.PosterPath = record[11]

	// Release Date
	if record[14] != "" {
		if t, err := time.Parse("2006-01-02", record[14]); err == nil {
			m.ReleaseDate = &t
		}
	}

	// Runtime
	if runtime, err := strconv.ParseFloat(record[16], 64); err == nil {
		m.Runtime = int(runtime)
	}

	m.Title = record[20]

	// Vote Average
	if rating, err := strconv.ParseFloat(record[22], 64); err == nil {
		m.VoteAverage = rating
	}

	// Vote Count
	if votes, err := strconv.Atoi(record[23]); err == nil {
		m.VoteCount = votes
	}

	// Raw genres for now
	m.Genres = ParseGenres(record[3])

	return m, nil
}