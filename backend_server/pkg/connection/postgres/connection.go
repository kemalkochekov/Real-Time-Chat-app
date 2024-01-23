package postgres

import (
	"backend_server/configs"
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type DBops interface {
	// Database quires
	GetPool() *sqlx.DB
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	Close() error
}
type Database struct {
	db *sqlx.DB
}

func (s *Database) GetPool() *sqlx.DB {
	return s.db
}
func (s *Database) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return s.db.GetContext(ctx, dest, query, args...)
}
func (s *Database) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return s.db.QueryRowContext(ctx, query, args...)
}
func (s *Database) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return s.db.QueryContext(ctx, query, args...)
}
func (s *Database) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}
func (s *Database) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return s.db.SelectContext(ctx, dest, query, args...)
}
func (s *Database) Close() error {
	if err := goose.Down(s.db.DB, "./internal/repository/migrations"); err != nil {
		fmt.Printf("goose migration down failed: %v", err)
	}
	return s.db.Close()
}
func GenerateDsn(cfgs *configs.Config) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfgs.Postgres.Host, cfgs.Postgres.Port, cfgs.Postgres.User, cfgs.Postgres.Password, cfgs.Postgres.DBName)
}

// cfgs database configuration from env file
func NewDatabase(ctx context.Context, cfgs *configs.Config) (*Database, error) {
	db, err := sqlx.Connect("postgres", GenerateDsn(cfgs))
	if err != nil {
		return nil, fmt.Errorf("could not create connection pool: %v", err)
	}
	if err := goose.Up(db.DB, "./internal/migrations"); err != nil {
		return nil, fmt.Errorf("goose migration up failed: %v", err)
	}
	return &Database{db: db}, nil
}
