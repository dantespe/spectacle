package testing

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/dantespe/spectacle/db"
	_ "github.com/mattn/go-sqlite3"
)

const baseSchema = "db/schema.sql"

// CreateTempSQLiteDB will attempt to create a temp file (database),
// load the Spectable schema into the database, and return a *os.File.
// It is the responsibility of the caller to remove the file by calling
// os.Remove(f.Name()) on the *os.File object.
func CreateTempSQLiteDB() (*os.File, error) {
	f, err := os.CreateTemp("", "spectacle_test_db")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp with err: %v", err)
	}

	db, err := sql.Open("sqlite3", f.Name())
	if err != nil {
		os.Remove(f.Name())
		return nil, fmt.Errorf("failed to create db handle with err: %v", err)
	}
	defer db.Close()

	schemaFile := os.Getenv("SPECTACLE_DIR") + "/" + baseSchema
	bytes, err := os.ReadFile(schemaFile)
	if err != nil {
		os.Remove(f.Name())
		return nil, fmt.Errorf("failed to read schema file: %s", schemaFile)
	}
	if _, err := db.Exec(string(bytes)); err != nil {
		os.Remove(f.Name())
		return nil, fmt.Errorf("failed to Exec(schema) with error: %v", err)
	}

	return f, nil
}

// CreateTempSQLiteEngine returns a new db.Engine from a temp file.
// If it succeeds, it's the responsibility of the caller to delete the temp
// file by calling os.Remove(fileName).
func CreateTempSQLiteEngine() (*db.Engine, string, error) {
	dbFile, err := CreateTempSQLiteDB()
	if err != nil {
		return nil, "", fmt.Errorf("failed to create temp SQLite DB: %v", err)
	}

	eng, err := db.New(
		db.WithSQLiteDatabaseFile(dbFile.Name()),
	)
	if err != nil {
		os.Remove(dbFile.Name())
		return nil, "", fmt.Errorf("failed to create DB engine: %v", err)
	}

	return eng, dbFile.Name(), nil
}
