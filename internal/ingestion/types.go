package ingestion

type TMDBMovie struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	Overview      string  `json:"overview"`
	OriginalLang  string  `json:"original_language"`
	ReleaseDate   string  `json:"release_date"`
	VoteAverage   float64 `json:"vote_average"`
	PosterPath    string  `json:"poster_path"`
	GenreIDs      []int   `json:"genre_ids"`
}

type DiscoverResponse struct {
	Page         int         `json:"page"`
	Results      []TMDBMovie `json:"results"`
	TotalPages   int         `json:"total_pages"`
	TotalResults int         `json:"total_results"`
}