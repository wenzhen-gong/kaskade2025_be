package db

import (
	"fmt"
	"kaskade_backend/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() *gorm.DB {

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	fmt.Print(dsn)
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatalf("❌ Failed to connect database: %v", err)
	}

	// 自动迁移数据库结构
	err = db.AutoMigrate(&models.User{}, &models.RequestHistory{}, &models.BenchmarkResult{})
	if err != nil {
		log.Fatalf("❌ Failed to migrate database: %v", err)
	}

	log.Println("✅ Database connected and migrated successfully")
	return db
}
