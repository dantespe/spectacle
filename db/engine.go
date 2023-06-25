// Package db contains database configuration.
package db

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// DatabaseProvider
type DatabaseProvider int

const (
	DatabaseProvider_UNKNOWN DatabaseProvider = iota
	DatabaseProvider_SQLITE
)

// Environment
type Environment int

const (
	Environment_UNKNOWN Environment = iota
	Environment_CUSTOM
	Environment_DEVELOPMENT
	Environment_TEST
)

type Engine struct {
	// Stores the type of Database.
	DatabaseProvider DatabaseProvider

	// Environment that the Engine belongs to.
	Environment Environment

	// DatabaseHandle
	DatabaseHandle *sql.DB
}

// Option for creating DB Engine
type Option func(*Engine) error

// New creates a new Engine with the given options.
func New(opts ...Option) (*Engine, error) {
	eng, err := Default()
	if err != nil {
		return nil, err
	}
	for _, o := range opts {
		o(eng)
	}
	if err := eng.createDatabaseHandler(); err != nil {
		return nil, err
	}
	return eng, nil
}

func WithDatabaseProvider(dp DatabaseProvider) Option {
	return func(e *Engine) error {
		e.DatabaseProvider = dp
		return nil
	}
}

func WithEnvironment(env Environment) Option {
	return func(e *Engine) error {
		e.Environment = env
		return nil
	}
}

func WithSQLiteDatabaseFile(dbPath string) Option {
	return func(e *Engine) error {
		e.DatabaseProvider = DatabaseProvider_SQLITE
		return e.createSQLiteDatabaseHandler(dbPath)
	}
}

func Default() (*Engine, error) {
	return &Engine{
		DatabaseProvider: DatabaseProvider_SQLITE,
		Environment:      Environment_DEVELOPMENT,
	}, nil
}

func (e *Engine) createDatabaseHandler() error {
	if e.DatabaseHandle != nil {
		return nil
	}

	if e.DatabaseProvider == DatabaseProvider_SQLITE {
		return e.createSQLiteHandler()
	}
	return fmt.Errorf("Unknown Engine Provider: %d", e.DatabaseProvider)
}

func (e *Engine) createSQLiteHandler() error {
	if e.Environment != Environment_DEVELOPMENT && e.Environment != Environment_TEST {
		return fmt.Errorf("Unknown Environment: %d", e.Environment)
	}

	// When we test, our working directory is spectacle/db.
	// When we run in production, it is spectacle/.
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	if strings.HasSuffix(wd, "/spectacle") {
		wd = wd + "/db"
	}

	envToDatabase := map[Environment]string{
		Environment_DEVELOPMENT: "dev.db",
		Environment_TEST:        "test.db",
	}
	databasePath := fmt.Sprintf("%s/%s", wd, envToDatabase[e.Environment])
	return e.createSQLiteDatabaseHandler(databasePath)
}

func (e *Engine) createSQLiteDatabaseHandler(dbPath string) error {
	if e.DatabaseProvider != DatabaseProvider_SQLITE && e.DatabaseProvider != DatabaseProvider_UNKNOWN {
		return fmt.Errorf("running SQLite handler on non-SQLite engine")
	}
	if _, err := os.Stat(dbPath); err != nil {
		return fmt.Errorf("failed to open databasePath (%s) with err: %v", dbPath, err)
	}
	dh, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("faield to create database handle with err: %v", err)
	}
	e.DatabaseHandle = dh
	return nil
}
