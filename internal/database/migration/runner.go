package migration

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {
	err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id SERIAL PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`).Error

	if err != nil {
		return err
	}

	files, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		return err
	}

	sort.Strings(files)

	for _, file := range files {
		name := filepath.Base(file)

		var count int64

		err := db.Table("schema_migrations").
			Where("name = ?", name).
			Count(&count).Error

		if err != nil {
			return err
		}

		if count > 0 {
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		log.Printf("Running migration: %s", name)

		if err := db.Exec(string(content)).Error; err != nil {
			return fmt.Errorf(
				"migration %s failed: %w",
				name,
				err,
			)
		}

		if err := db.Table("schema_migrations").
			Create(map[string]interface{}{
				"name": name,
			}).Error; err != nil {
			return err
		}
	}

	return nil
}