# Apply all migrations

go run cmd/migrate/main.go -up

# Apply specific number of migrations

go run cmd/migrate/main.go -up -step 2

# Rollback all migrations

go run cmd/migrate/main.go -down

# Rollback specific number of migrations

go run cmd/migrate/main.go -down -step 1

# Migrate to specific version

go run cmd/migrate/main.go -version 3

# Check current migration version

go run cmd/migrate/main.go

# Seed all data

go run cmd/seed/main.go -all

# Seed specific data types

go run cmd/seed/main.go -users
go run cmd/seed/main.go -categories
go run cmd/seed/main.go -products

# Clear all data before seeding

go run cmd/seed/main.go -clear -all

# Clear and seed specific data

go run cmd/seed/main.go -clear -users -categories

### Running Tests

To run the entire test suite:

```shellscript
go test ./...
```

To run just unit tests (skipping integration tests):

```shellscript
go test -short ./...
```

To run tests for a specific package:

```shellscript
go test github.com/zenfulcode/commercify/internal/domain/entity
```
