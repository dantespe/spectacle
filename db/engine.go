// Package db contains database configuration.
package db

import (
	"database/sql"
	"fmt"
	"os"

	pq "github.com/lib/pq"
)

// Postgres Defaults
const (
	DefaultPostgresHost = "localhost"
	DefaultPostgresPort = 5432
	DefaultPostgresUser = "postgres"
	// This is not the password. This is the environment variable that stores the password.
	DefaultPostgresPassword     = "PGPASSWORD"
	DefaultPostgresDatabaseName = "test"
)

// DatabaseProvider
type DatabaseProvider int

const (
	DatabaseProvider_UNKNOWN DatabaseProvider = iota
	DatabaseProvider_POSTGRES
)

// Environment
type Environment int

const (
	Environment_UNKNOWN Environment = iota
	Environment_DEVELOPMENT
	Environment_TEST
	// Custom should remain last
	Environment_CUSTOM
)

type Engine struct {
	// Stores the type of Database.
	DatabaseProvider DatabaseProvider

	// Environment that the Engine belongs to.
	Environment Environment

	// DatabaseHandle
	DatabaseHandle *sql.DB

	// PostgresConfig
	pc *postgresConfig
}

type postgresConfig struct {
	Host                string
	Port                int
	User                string
	Password            string
	DatabaseName        string
	SSLMode             string
	WithoutDatabaseName bool
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

func WithHost(host string) Option {
	return func(e *Engine) error {
		e.pc.Host = host
		return nil
	}
}

func WithPort(port int) Option {
	return func(e *Engine) error {
		e.pc.Port = port
		return nil
	}
}

func WithUser(user string) Option {
	return func(e *Engine) error {
		e.pc.User = user
		return nil
	}
}

func WithDefaultPassword() Option {
	return WithPasswordFromEnvironment(DefaultPostgresPassword)
}

func WithPasswordFromEnvironment(env string) Option {
	return func(e *Engine) error {
		passwd := os.Getenv(env)
		if passwd == "" {
			return fmt.Errorf("got empty password from environment variable: %s, want: non-empty", env)
		}
		e.pc.Password = passwd
		return nil
	}
}

func WithoutDatabaseName() Option {
	return func(e *Engine) error {
		e.pc.DatabaseName = ""
		e.pc.WithoutDatabaseName = true
		return nil
	}
}

func WithDatabaseName(dbName string) Option {
	return func(e *Engine) error {
		e.pc.DatabaseName = dbName
		return nil
	}
}

func WithSSLMode(ssl string) Option {
	return func(e *Engine) error {
		e.pc.SSLMode = ssl
		return nil
	}
}

func Default() (*Engine, error) {
	opts := []Option{
		WithDatabaseProvider(DatabaseProvider_POSTGRES),
		WithEnvironment(Environment_TEST),
		WithHost(DefaultPostgresHost),
		WithPort(DefaultPostgresPort),
		WithDatabaseName(DefaultPostgresDatabaseName),
		WithUser(DefaultPostgresUser),
		WithDefaultPassword(),
		WithSSLMode("disable"),
	}

	eng := &Engine{
		pc: &postgresConfig{},
	}
	for _, o := range opts {
		o(eng)
	}

	return eng, nil
}

// Connection returns a string of the Postgres Info.
func (e *Engine) Connection() (string, error) {
	conn := fmt.Sprintf("user=%s sslmode=%s ", e.pc.User, e.pc.SSLMode)
	if e.pc.Host == "" {
		return "", fmt.Errorf("host cannot be empty")
	}
	conn += fmt.Sprintf("host=%s ", e.pc.Host)

	if e.pc.Password == "" {
		// Password is technically optional, but we will
		// enforce that we don't use a password-less
		return "", fmt.Errorf("postgres password cannot be empty")
	}
	conn += fmt.Sprintf("password=%s ", e.pc.Password)

	if e.pc.DatabaseName != "" {
		conn += fmt.Sprintf("dbname=%s ", e.pc.DatabaseName)
	}
	return conn, nil
}

func (e *Engine) createDatabaseHandler() error {
	if e.DatabaseHandle != nil {
		return nil
	}
	if e.DatabaseProvider == DatabaseProvider_POSTGRES {
		return e.createPostgresHandler()
	}
	return fmt.Errorf("Unsupported DatabaseProvider: %d", e.DatabaseProvider)
}

func (e *Engine) getPostgresDatabaseName() error {
	if e.pc.WithoutDatabaseName {
		e.pc.DatabaseName = ""
		return nil
	}
	if e.Environment == Environment_DEVELOPMENT {
		e.pc.DatabaseName = "dev"
		return nil
	}
	if e.Environment == Environment_TEST {
		e.pc.DatabaseName = "test"
		return nil
	}
	if e.Environment == Environment_CUSTOM {
		return nil
	}
	return fmt.Errorf("unsupported environment: %d", e.Environment)
}

func (e *Engine) createPostgresHandler() error {
	if e.Environment <= Environment_UNKNOWN || e.Environment > Environment_CUSTOM {
		return fmt.Errorf("Unknown Environment: %d", e.Environment)
	}
	if err := e.getPostgresDatabaseName(); err != nil {
		return fmt.Errorf("failed to get postgres database name with err: %v", err)
	}
	conn, err := e.Connection()
	if err != nil {
		return fmt.Errorf("failed to get Postgres connection with err: %v", err)
	}
	dh, err := sql.Open("postgres", conn)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres database: %v", err)
	}
	e.DatabaseHandle = dh
	return nil
}

type Tx struct {
	engine *Engine
	tx     *sql.Tx
	stmt   *sql.Stmt
	table  string
	args   []string
	buf    int
}

func NewTx(e *Engine, table string, args ...string) (*Tx, error) {
	if e == nil {
		return nil, fmt.Errorf("cannot create new transcation with nil engine")
	}

	tx, err := e.DatabaseHandle.Begin()
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare(pq.CopyIn(table, args...))
	if err != nil {
		return nil, err
	}

	return &Tx{
		engine: e,
		tx:     tx,
		stmt:   stmt,
		table:  table,
		args:   args,
	}, nil
}

func (t *Tx) Exec(args ...interface{}) error {
	if _, err := t.stmt.Exec(args...); err != nil {
		return err
	}
	t.buf++

	if t.isFull() {
		t.flush()
		return t.reset()
	}
	return nil
}

func (t *Tx) isFull() bool {
	return t.buf > 500
}

func (t *Tx) flush() error {
	if t.buf == 0 {
		return nil
	}
	if err := t.stmt.Close(); err != nil {
		return err
	}
	return t.tx.Commit()
}

func (t *Tx) reset() error {
	tx, err := t.engine.DatabaseHandle.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(pq.CopyIn(t.table, t.args...))
	if err != nil {
		return err
	}

	t.tx = tx
	t.stmt = stmt
	t.buf = 0
	return nil
}

func (t *Tx) Close() error {
	return t.flush()
}
