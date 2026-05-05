package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB
var JWTSecret = []byte("company-management-secret-key")

type Config struct {
	Port        int
	DBType      string
	DBConn      string
	JWTSecret   []byte
}

func InitConfig() *Config {
	port, _ := strconv.Atoi(getEnv("PORT", "8080"))
	
	return &Config{
		Port:      port,
		DBType:    getEnv("DB_TYPE", "sqlite3"),
		DBConn:    getEnv("DB_CONN", "./company.db"),
		JWTSecret: []byte(getEnv("JWT_SECRET", "company-management-secret-key")),
	}
}

func InitDB(cfg *Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DBConn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %v", err)
	}
	
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	
	DB = db
	return db, nil
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
