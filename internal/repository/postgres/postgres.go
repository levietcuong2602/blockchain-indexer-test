package postgres

import (
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/unanoc/blockchain-indexer/internal/repository/models"
)

type Database struct {
	Gorm *gorm.DB
}

func New(url string, log bool) (*Database, error) {
	var logMode logger.LogLevel
	if log {
		logMode = logger.Info
	}

	cfg := &gorm.Config{
		Logger:                 logger.Default.LogMode(logMode),
		SkipDefaultTransaction: true,
	}

	db, err := gorm.Open(postgres.Open(url), cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres conn: %w", err)
	}

	return &Database{Gorm: db}, nil
}

func Setup(db *Database) error {
	conn := db.Gorm.Begin()

	if err := conn.AutoMigrate(
		&models.BlockTracker{},
		&models.Transaction{},
		&models.Node{},
		&models.Block{},
		&models.Collection{},
		&models.Event{},
	); err != nil {
		return fmt.Errorf("failed to make auto migrations: %w", err)
	}

	if err := conn.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit a transaction: %w", err)
	}

	return nil
}

func IsErrNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
