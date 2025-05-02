package db

import (
	"fmt"
	"log"
	"wider-circle-backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitializeDatabase(driver, dsn string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	switch driver {
	case "sqlite", "sqlite3":
		log.Printf("Connecting to SQLite database: %s", dsn)
		db, err = gorm.Open(sqlite.Open(dsn), config)
	case "postgres", "postgresql":
		log.Printf("Connecting to PostgreSQL database: %s", dsn)
		db, err = gorm.Open(postgres.Open(dsn), config)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = db.AutoMigrate(&models.Employee{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database schema: %w", err)
	}

	return db, nil
}

func FetchEmployees(db *gorm.DB) ([]models.Employee, error) {
	var employees []models.Employee
	err := db.Find(&employees).Error
	if err != nil {
		return nil, err
	}
	return employees, nil
}

func StoreEmployees(db *gorm.DB, employees []models.Employee) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for _, emp := range employees {
			if err := tx.Create(&emp).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func EmployeeExists(db *gorm.DB) (bool, error) {
	var count int64
	err := db.Model(&models.Employee{}).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func SaveEmployees(db *gorm.DB, employees []models.Employee) error {
	return db.Create(&employees).Error
}

func GetEmployees(db *gorm.DB) ([]models.Employee, error) {
	var employees []models.Employee
	err := db.Find(&employees).Error
	return employees, err
}
