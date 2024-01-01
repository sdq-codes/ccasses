package psql

import (
	"context"
	_ "github.com/golang-migrate/migrate/v4/source/file" // import file driver for migrate
	"github.com/sdq-codes/ccasses/gateway/core/errors"
	"github.com/sdq-codes/ccasses/gateway/core/logging"
)

const (
	// ErrDriverInit is returned when we cannot initialize the driver.
	ErrDriverInit = errors.Error("failed to initialize postgres driver")
	// ErrMigrateInit is returned when we cannot initialize the migrate driver.
	ErrMigrateInit = errors.Error("failed to initialize migration driver")
	// ErrMigration is returned when we cannot run a migration.
	ErrMigration = errors.Error("failed to migrate database")
)

// MigratePostgres migrates the database to the latest version.
func (d *Driver) MigratePostgres(ctx context.Context, migrationsPath string) error {
	err := d.GetDB().AutoMigrate()
	if err != nil {
		logging.From(ctx).Info("Migration failed")
		return err
	}

	logging.From(ctx).Info("migrations successfully run")

	return nil
}

// RevertMigrations reverts the database to the previous version.
func (d *Driver) RevertMigrations(ctx context.Context, migrationsPath string) error {
	panic("implement me")
}
