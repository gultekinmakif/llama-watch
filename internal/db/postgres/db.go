// Postgres handle and AutoMigrate. Opens the connection and runs Migrate at boot.
package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/gultekinmakif/llama-watch/internal/models"
	gormpg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Package handle. Set by New, read by Get.
var db *gorm.DB

const notInitMsg = "postgres: not initialized"

// New opens the connection and stores the package handle. Call once at startup.
func New(dsn string) error {
	gpg := gormpg.Open(dsn)
	d, err := gorm.Open(gpg, &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := d.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return err
	}

	db = d
	return nil
}

// Get returns the package handle. Panics if New has not run because that is a programming error.
func Get() *gorm.DB {
	if db == nil {
		panic(notInitMsg)
	}
	return db
}

func Close() error {
	if db == nil {
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Ping verifies the connection is alive.
func Ping(ctx context.Context) error {
	if db == nil {
		return errors.New(notInitMsg)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// Migrate drops the retired legacy tables, AutoMigrates the live models, then
// rebuilds dim_file_coverage as a view over matrix. All steps are idempotent.
func Migrate() error {
	if err := db.Exec("DROP TABLE IF EXISTS protocols, adapter_files, commit_refs, refresh_runs CASCADE").Error; err != nil {
		return err
	}
	return db.AutoMigrate(
		&models.Matrix{},
		&models.ProtocolIdentity{}
	)
}
