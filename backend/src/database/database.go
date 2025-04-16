package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var DB *sql.DB

func Connect(dsn string) {
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Error pinging database:", err)
	}

	err = setUpDatabase()
	if err != nil {
		log.Fatal("Error setting up database:", err)
	}
}

func setUpDatabase() error {
	sqlPath := "src/database/migrations"
	sqlFiles, err := filepath.Glob(filepath.Join(sqlPath, "*.sql"))
	if err != nil {
		return fmt.Errorf("error finding migration files: %w", err)
	}

	for _, file := range sqlFiles {
		if err := executeSQLFile(file); err != nil {
			return err
		}
	}
	return nil
}

func executeSQLFile(filePath string) error {
	sqlContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading %s: %w", filePath, err)
	}
	if _, err := DB.Exec(string(sqlContent)); err != nil {
		return fmt.Errorf("error executing %s: %w", filePath, err)
	}
	return nil
}
