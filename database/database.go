// Package database: conexion a la base de datos sqlite
package database

import (
	"jos3lo89/go-image-service/config"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error

	db, err := gorm.Open(sqlite.Open(config.AppConfig.DBPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	log.Println("Conexión a la base de datos exitosa.")
	db.AutoMigrate()
	log.Println("Migración de la base de datos completada.")
}
