package psql

import (
	"context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq" // imports the postgres driver
	"github.com/sdq-codes/claimclam/internal/core/errors"
)

const (
	// ErrConnect is returned when we cannot connect to the database.
	ErrConnect = errors.Error("failed to connect to postgres db")
	// ErrClose is returned when we cannot close the database.
	ErrClose = errors.Error("failed to close postgres db connection")
)

// Config represents the configuration for our postgres database.
type Config struct {
	DSN string `env:"POSTGRES_DSN" validate:"required"`
}

// Driver provides an implementation for connecting to a postgres database.
type Driver struct {
	cfg Config
	db  *gorm.DB
}

// New instantiates a instance of the Driver.
func New(cfg Config) *Driver {
	return &Driver{
		cfg: cfg,
	}
}

// Connect connects to the database.
func (d *Driver) Connect(ctx context.Context) error {
	dsn := "host=localhost user=guest password=guest dbname=luxridr port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return ErrConnect.Wrap(err)
	}
	d.db = db
	if err != nil {
		return ErrConnect.Wrap(err)
	}
	return nil
}

// Close closes the database connection.
func (d *Driver) Close(ctx context.Context) error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return ErrConnect.Wrap(err)
	}
	if err := sqlDB.Close(); err != nil {
		return ErrClose.Wrap(err)
	}
	return nil
}

// GetDB returns the underlying database connection.
func (d *Driver) GetDB() *gorm.DB {
	return d.db
}
