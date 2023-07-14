// Package db_test contains tests for db.
package db_test

import (
	"strings"
	"testing"

	"github.com/dantespe/spectacle/db"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		desc                 string
		opts                 []db.Option
		expectedProvider     db.DatabaseProvider
		expectedEnvironment  db.Environment
		expectedDatabaseName string
	}{
		{
			desc:                "default_no_opts",
			opts:                []db.Option{},
			expectedProvider:    db.DatabaseProvider_POSTGRES,
			expectedEnvironment: db.Environment_TEST,
		},
		{
			desc: "postgres_dev",
			opts: []db.Option{
				db.WithEnvironment(db.Environment_DEVELOPMENT),
				db.WithoutDatabaseName(),
			},
			expectedProvider:    db.DatabaseProvider_POSTGRES,
			expectedEnvironment: db.Environment_DEVELOPMENT,
		},
		{
			desc: "postgres_custom",
			opts: []db.Option{
				db.WithEnvironment(db.Environment_CUSTOM),
				db.WithDatabaseName("tmp_0001"),
			},
			expectedProvider:    db.DatabaseProvider_POSTGRES,
			expectedEnvironment: db.Environment_CUSTOM,
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
			conn, err := eng.Connection()
			if err != nil {
				t.Fatalf("got unexpected error for Connection(): %v", err)
			}
			if !strings.Contains(conn, tc.expectedDatabaseName) {
				t.Errorf("failed to find database name in Connection: %s", tc.expectedDatabaseName)
			}
		})
	}
}
