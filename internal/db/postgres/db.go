package postgres

import (
	"context"
	"time"

	"github.com/gultekinmakif/go-http-server/internal/models"
	gormpg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Package handle. Set by New, read by Get.
var db *gorm.DB

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

// Get returns the package handle. Panics if New hasn't run — misuse is a programming error.
func Get() *gorm.DB {
	if db == nil {
		panic("postgres: not initialized")
	}
	return db
}

func Close() error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// Migrate applies the current schema via AutoMigrate.
func Migrate() error {
	return db.AutoMigrate(&models.Post{})
}
