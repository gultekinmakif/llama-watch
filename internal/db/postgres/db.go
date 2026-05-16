package postgres

import (
	"context"
	"time"

	"github.com/gultekinmakif/llama-watch/internal/models"
	gormpg "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

// Get returns the package handle. Panics if New has not run because that is a programming error.
func Get() *gorm.DB {
	if db == nil {
		panic("postgres: not initialized")
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

// dimensionSeed is the static set of dimension kinds the schema expects to
// exist before any adapter_files rows are inserted. New kinds are added by
// hand. Keep this list in lockstep with adapter_files.dimension_kind values.
var dimensionSeed = []models.Dimension{
	{Kind: "tvl", DisplayName: "TVL"},
	{Kind: "dailyFees", DisplayName: "Daily Fees"},
	{Kind: "dailyRevenue", DisplayName: "Daily Revenue"},
	{Kind: "dailyVolume", DisplayName: "Daily Volume"},
	{Kind: "dailyNotionalVolume", DisplayName: "Daily Notional Volume"},
	{Kind: "dailyPremiumVolume", DisplayName: "Daily Premium Volume"},
	{Kind: "openInterestAtEnd", DisplayName: "Open Interest"},
	{Kind: "dailyBridgeVolume", DisplayName: "Daily Bridge Volume"},
	{Kind: "dailyActiveUsers", DisplayName: "Daily Active Users"},
}

// Migrate applies the current schema via AutoMigrate, then UPSERTs the static
// dimensions rows. Order matters: dimensions and protocols are parents for the
// FKs in adapter_files; commit_refs FK-references adapter_files.
func Migrate() error {
	if err := db.AutoMigrate(
		&models.Dimension{},
		&models.Protocol{},
		&models.AdapterFile{},
		&models.CommitRef{},
		&models.RefreshRun{},
	); err != nil {
		return err
	}

	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "kind"}},
		DoUpdates: clause.AssignmentColumns([]string{"display_name"}),
	}).Create(&dimensionSeed).Error
}
