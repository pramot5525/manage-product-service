package database

import (
	"fmt"

	"manage-product-service/internal/config"
	"manage-product-service/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cfg config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&model.Product{}); err != nil {
		return nil, err
	}

	return db, nil
}
