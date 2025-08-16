package database

import (
	"log"

	"github.com/Motionists/mt-order/server/internal/config"
	"github.com/Motionists/mt-order/server/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Init(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(mysql.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移数据库表
	err = db.AutoMigrate(
		&models.User{},
		&models.Merchant{},
		&models.Dish{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database connected and migrated successfully")
	return db
}
