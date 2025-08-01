package main

import (
	"Tournament/internal/config"
	"context"
	"embed"
	"log"

	"github.com/uptrace/bun/migrate"

	_ "github.com/lib/pq"
)

/*
   Embed every migration file (*.sql) that lives in the migrations folder.
   Adjust the pattern if your folder name is different.
*/
//go:embed migrations/*.sql
var migrationFiles embed.FS

func main() {
	dbConfig := config.NewDatabaseConfig()
	db, err := dbConfig.NewBunDB()

	if err != nil {
		log.Fatalf("Failed to initialize Bun DB: %v", err)
	}
	defer db.Close()

	// 2. Create and populate *migrate.Migrations.
	migrations := migrate.NewMigrations()
	if err := migrations.Discover(migrationFiles); err != nil {
		log.Fatalf("discover migrations: %v", err)
	}

	// 3. Wire everything into a Migrator.
	migrator := migrate.NewMigrator(db, migrations)

	ctx := context.Background()

	// 4. Ensure the internal schema_migrations table exists.
	if err := migrator.Init(ctx); err != nil {
		log.Fatalf("init migrator: %v", err)
	}

	// 5. Apply all pending migration groups.
	for {
		group, err := migrator.Migrate(ctx)
		if err != nil {
			log.Fatalf("migrate: %v", err)
		}
		if group.IsZero() { // nothing left to run
			log.Println("database is up-to-date")
			break
		}
		log.Printf("applied %s\n", group)
	}
}
