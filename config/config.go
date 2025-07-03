package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func DBInit() *gorm.DB {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, pass, name, port,
	)
	log.Println("DSN :", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Connect to database is Failed:", err)
	}

	var dbName string
	db.Raw("SELECT current_database()").Scan(&dbName)
	log.Println("PostgreSQL connect to:", dbName)

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Error get database:", err)
	}
	err = sqlDB.Ping()
	if err != nil {
		log.Fatal("Cannot access to database:", err)
	}

	log.Println("Success connect to database")
	return db
}
