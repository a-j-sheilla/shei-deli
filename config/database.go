// config/database.go
package config

import (
    "log"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "shei-deli/models"
)

var DB *gorm.DB

func InitDatabase() {
    var err error
    DB, err = gorm.Open(sqlite.Open("shei_deli.db"), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Auto-migrate the schema
    err = DB.AutoMigrate(&models.Recipe{}, &models.Feedback{}, &models.User{})
    if err != nil {
        log.Fatalf("Failed to migrate database: %v", err)
    }

    log.Println("Database connected and migrated!")
}
