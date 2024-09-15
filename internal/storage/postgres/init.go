package psq

import (
	"context"
	"fmt"

	"gitlab.com/Satori27/avito/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Storage struct {
	db *gorm.DB
}

func New(cancel context.CancelFunc, s *Storage, cfg *config.Config) error {
	const op = "storage.postgres.New"

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s target_session_attrs=read-write", cfg.DB.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer cancel()
	db.Logger = logger.Default.LogMode(logger.Info)

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	Exec(db)

	// err = db.AutoMigrate(&models.Tender{}, &models.Bid{}, &models.TenderVesion{}, &models.BidVersion{}, &models.BidFeedback{})
	// if err != nil {
	// 	return fmt.Errorf("%s: %w", op, err)
	// }
	s.db = db
	return nil
}

func Exec(db *gorm.DB) {

	db.Exec(`CREATE TYPE organization_type AS ENUM (
		'IE',
		'LLC',
		'JSC'
	);
	`)

	db.Exec(
		`CREATE TYPE tender_service_type AS ENUM (
	'Construction',
	'Delivery',
	'Manufacture'
	);
	`)

	db.Exec(
		`CREATE TYPE tender_status AS ENUM (
		'Created',
		'Published',
		'Closed'
		);
	`)

	db.Exec(
		`CREATE TYPE bid_status AS ENUM (
		'Created',
		'Published',
		'Canceled',
		'Approved',
		'Rejected'
		);
	`)

	db.Exec(
		`CREATE TYPE bid_author_type AS ENUM (
		'Organization',
		'User'
		);
	`)

}
