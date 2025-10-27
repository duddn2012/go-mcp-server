package database

import (
	"fmt"
	"go_mcp_server/config"
	"log"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	instance *gorm.DB
	once     sync.Once
)

func GetDB(cfg *config.Config) *gorm.DB {
	once.Do(func() {
		var err error
		instance, err = connect(cfg)
		if err != nil {
			log.Fatalf("[Database] Failed to connect: %v", err)
		}
		log.Printf("[Database] Connected to %s:%s/%s", cfg.DBHost, cfg.DBPort, cfg.DBName)
	})
	return instance
}

func connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return db, nil
}
