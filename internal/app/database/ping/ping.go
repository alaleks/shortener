package ping

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func Run(dsn string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("database connection error: %w", err)
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		return fmt.Errorf("ping error db: %w", err)
	}

	return nil
}
