// Package db_test contains tests for db.
package db_test

import (
	"os"
	"testing"

	"github.com/dantespe/spectacle/db"
	spectesting "github.com/dantespe/spectacle/testing"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		desc                string
		opts                []db.Option
		expectedProvider    db.DatabaseProvider
		expectedEnvironment db.Environment
	}{
		{
			desc:                "default_no_opts",
			opts:                []db.Option{},
			expectedProvider:    db.DatabaseProvider_SQLITE,
			expectedEnvironment: db.Environment_DEVELOPMENT,
		},
		{
			desc: "sqlite_test",
			opts: []db.Option{
				db.WithDatabaseProvider(db.DatabaseProvider_SQLITE),
				db.WithEnvironment(db.Environment_TEST),
			},
			expectedProvider:    db.DatabaseProvider_SQLITE,
			expectedEnvironment: db.Environment_TEST,
		},
		{
			desc: "sqlite_test",
			opts: []db.Option{
				db.WithDatabaseProvider(db.DatabaseProvider_SQLITE),
				db.WithEnvironment(db.Environment_TEST),
			},
			expectedProvider:    db.DatabaseProvider_SQLITE,
			expectedEnvironment: db.Environment_TEST,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			eng, err := db.New(tc.opts...)
			if err != nil {
				t.Fatalf("got unexpected error for New(tc.opts...): %v", err)
			}
			if eng.DatabaseProvider != tc.expectedProvider {
				t.Errorf("got DatabaseProvider %d, want %d", eng.DatabaseProvider, tc.expectedProvider)
			}
			if eng.Environment != tc.expectedEnvironment {
				t.Errorf("got Environment %d, want %d", eng.Environment, tc.expectedEnvironment)
			}
			if eng.DatabaseHandle == nil {
				t.Errorf("got nil DatabaseHandle, want non-nil")
			}
		})
	}
}

func TestCustomSQLite(t *testing.T) {
	dbFile, err := spectesting.CreateTempSQLiteDB()
	if err != nil {
		t.Fatalf("failed to create temp SQLite DB: %v", err)
	}
	defer os.Remove(dbFile.Name())

	eng, err := db.New(
		db.WithSQLiteDatabaseFile(dbFile.Name()),
	)
	if err != nil {
		t.Fatalf("failed to create new DB Engine: %v", err)
	}

	if eng.DatabaseProvider != db.DatabaseProvider_SQLITE {
		t.Errorf("Got DatabaseProvider %d, want %d", eng.DatabaseProvider, db.DatabaseProvider_SQLITE)
	}
	if eng.DatabaseHandle == nil {
		t.Errorf("got nil DatabaseHandle, want non-nil")
	}
}
