package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"

	"streamflix-backend/internal/ingestion"
	"streamflix-backend/internal/config"
	"streamflix-backend/internal/database"
	"streamflix-backend/internal/movie"
)

func main() {

	// Connect to PostgreSQL
	cfg := config.Load()
	database.Connect(cfg)

	repo := movie.NewRepository(database.DB)

	file, err := os.Open("data/movies_metadata.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true

	reader.FieldsPerRecord = -1

	// Skip header
	_, err = reader.Read()
	if err != nil {
		log.Fatal(err)
	}

	var batch []movie.Movie
	count := 0

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println(err)
			continue
		}

		m, err := ingestion.ParseMovie(record)
		if err != nil {
			log.Println(err)
			continue
		}

		// Skip invalid movies
		if m.Title == "" || m.TMDBID == 0 {
			continue
		}

		batch = append(batch, *m)

		if len(batch) == 1000 {
			if err := repo.CreateBatch(batch); err != nil {
				log.Fatal(err)
			}

			count += len(batch)
			log.Printf("Imported %d movies...\n", count)

			// Reuse the underlying array
			batch = batch[:0]
		}
	}

	// Insert remaining movies
	if len(batch) > 0 {
		if err := repo.CreateBatch(batch); err != nil {
			log.Fatal(err)
		}

		count += len(batch)
	}

	log.Printf("✅ Finished importing %d movies\n", count)

	log.Printf("Parsed %d movies\n", count)
}