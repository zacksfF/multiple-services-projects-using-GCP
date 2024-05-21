package migratedb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Darwin is a helper struct to access the Validate and migration functions.
type darwin struct {
	driver     Driver
	migrations []Migration
}

func Seed(ctx context.Context, db *sqlx.DB) (err error) {
	if err := sqldb.StatusCheck(ctx, db); err != nil {
		return fmt.Errorf("status check database: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if errTx := tx.Rollback(); errTx != nil {
			if errors.Is(errTx, sql.ErrTxDone) {
				return
			}

			err = fmt.Errorf("rollback: %w", errTx)
			return
		}
	}()

	if _, err := tx.Exec(seedDoc); err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

// Migrate attempts to bring the database up to date with the migrations
// defined in this package.
func Migrate(ctx context.Context, db *sqlx.DB) error {
	if err := sqldb.StatusCheck(ctx, db); err != nil {
		return fmt.Errorf("status check database: %w", err)
	}

	driver, err := generic.New(db.DB, postgres.Dialect{})
	if err != nil {
		return fmt.Errorf("construct darwin driver: %w", err)
	}

	d := darwin.New(driver, darwin.ParseMigrations(migrateDoc))
	return d.Migrate()
}
