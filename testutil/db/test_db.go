package db

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// TestDB provides a PostgreSQL database for testing
type TestDB struct {
	DB *sql.DB
}

// NewTestDB creates a new test database
func NewTestDB() (*TestDB, error) {
	// Get connection details from environment variables or use defaults for testing
	host := getEnv("TEST_DB_HOST", "localhost")
	port := getEnv("TEST_DB_PORT", "5432")
	user := getEnv("TEST_DB_USER", "postgres")
	password := getEnv("TEST_DB_PASSWORD", "postgres")
	dbname := getEnv("TEST_DB_NAME", "commercify_test")

	// Create connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Ping database to verify connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	// Run migrations
	if err := runMigrations(db, dbname); err != nil {
		db.Close()
		return nil, err
	}

	return &TestDB{DB: db}, nil
}

// Close closes the database connection
func (tdb *TestDB) Close() {
	if tdb.DB != nil {
		tdb.DB.Close()
	}
}

// Clean cleans all data from the database
func (tdb *TestDB) Clean() error {
	// Tables to clean in reverse order of dependencies
	tables := []string{
		"order_items",
		"orders",
		"cart_items",
		"carts",
		"products",
		"categories",
		"users",
	}

	// Begin transaction
	tx, err := tdb.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Disable foreign key checks
	if _, err := tx.Exec("SET CONSTRAINTS ALL DEFERRED"); err != nil {
		return err
	}

	// Delete data from tables
	for _, table := range tables {
		if _, err := tx.Exec(fmt.Sprintf("DELETE FROM %s", table)); err != nil {
			return err
		}
		// Reset sequences
		if _, err := tx.Exec(fmt.Sprintf("ALTER SEQUENCE %s_id_seq RESTART WITH 1", table)); err != nil {
			return err
		}
	}

	// Re-enable foreign key checks
	if _, err := tx.Exec("SET CONSTRAINTS ALL IMMEDIATE"); err != nil {
		return err
	}

	return nil
}

// runMigrations runs database migrations
func runMigrations(db *sql.DB, dbName string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../../migrations",
		dbName,
		driver,
	)
	if err != nil {
		return err
	}

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

// SetupTestDB sets up a test database for use in tests
func SetupTestDB(t *testing.T) *TestDB {
	testDB, err := NewTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}

	// Clean the database
	if err := testDB.Clean(); err != nil {
		testDB.Close()
		t.Fatalf("Failed to clean test database: %v", err)
	}

	return testDB
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
