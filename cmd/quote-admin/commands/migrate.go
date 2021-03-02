package commands

import (
	"fmt"

	"github.com/johanronkko/quote-service/internal/business/data/schema"
	"github.com/johanronkko/quote-service/internal/foundation/database"
)

// Migrate creates the schema in the database.
func Migrate(cfg database.Config) error {
	db, err := database.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	if err := schema.Migrate(db); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	fmt.Println("migrations complete")
	return nil
}
