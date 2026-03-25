package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
)

func RunMigrations(db *Postgres, dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var migrationFiles []string
	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".sql" {
			migrationFiles = append(migrationFiles, filepath.Join(dir, f.Name()))
		}
	}

	sort.Strings(migrationFiles)

	for _, file := range migrationFiles {
		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read error %s: %w", file, err)
		}

		if _, err := db.DB.Exec(string(sqlBytes)); err != nil {
			return fmt.Errorf("migration failed %s: %w", file, err)
		}

		log.Println("✅ Applied migration:", file)
	}

	return nil
}