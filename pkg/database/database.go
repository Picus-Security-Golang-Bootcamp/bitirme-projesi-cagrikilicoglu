package database

import (
	"time"

	"github.com/cagrikilicoglu/shopping-basket/pkg/config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect connects to the database with configuration input
func Connect(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DBConfig.DataSourceName), &gorm.Config{})
	if err != nil {
		zap.L().Fatal("Cannot connect to database", zap.Error(err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		zap.L().Fatal("Cannot get sqlDB from database", zap.Error(err))
	}
	err = sqlDB.Ping()
	if err != nil {
		zap.L().Fatal("Database connection is not alive", zap.Error(err))
	}

	sqlDB.SetMaxOpenConns(cfg.DBConfig.MaxOpen)
	sqlDB.SetMaxIdleConns(cfg.DBConfig.MaxIdle)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.DBConfig.MaxLifetime) * time.Second)

	return db
}
