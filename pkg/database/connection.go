package database

import (
	"fmt"
	"log"

	"github.com/glebarez/sqlite"
	"github.com/jovan/mybanksoal-api/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDatabase(cfg *config.Config) *gorm.DB {
	var db *gorm.DB
	var err error

	switch cfg.Database.Driver {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Name,
		)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
	default:
		log.Fatalf("Invalid database driver: %s", cfg.Database.Driver)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return db
}
