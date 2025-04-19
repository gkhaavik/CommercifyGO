package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/infrastructure/database"
)

func main() {
	// Define command line flags
	upFlag := flag.Bool("up", false, "Run migrations up")
	downFlag := flag.Bool("down", false, "Rollback migrations")
	versionFlag := flag.Int("version", -1, "Migrate to specific version")
	stepFlag := flag.Int("step", 0, "Number of migrations to apply (up) or rollback (down)")
	flag.Parse()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.NewPostgresConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create migration instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		cfg.Database.DBName,
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}

	// Execute migration command based on flags
	if *upFlag {
		if *stepFlag > 0 {
			if err := m.Steps(*stepFlag); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("Failed to apply %d migrations: %v", *stepFlag, err)
			}
			fmt.Printf("Applied %d migrations\n", *stepFlag)
		} else {
			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("Failed to apply migrations: %v", err)
			}
			fmt.Println("Applied all migrations")
		}
	} else if *downFlag {
		if *stepFlag > 0 {
			if err := m.Steps(-(*stepFlag)); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("Failed to rollback %d migrations: %v", *stepFlag, err)
			}
			fmt.Printf("Rolled back %d migrations\n", *stepFlag)
		} else {
			if err := m.Down(); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("Failed to rollback migrations: %v", err)
			}
			fmt.Println("Rolled back all migrations")
		}
	} else if *versionFlag >= 0 {
		if err := m.Migrate(uint(*versionFlag)); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to migrate to version %d: %v", *versionFlag, err)
		}
		fmt.Printf("Migrated to version %d\n", *versionFlag)
	} else {
		// If no flags provided, print current version
		version, dirty, err := m.Version()
		if err != nil && err != migrate.ErrNilVersion {
			log.Fatalf("Failed to get migration version: %v", err)
		}

		if err == migrate.ErrNilVersion {
			fmt.Println("No migrations applied yet")
		} else {
			fmt.Printf("Current migration version: %d (dirty: %t)\n", version, dirty)
		}

		// Print usage
		fmt.Println("\nUsage:")
		flag.PrintDefaults()
	}
}
