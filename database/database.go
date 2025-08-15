// // Package database: conexion a la base de datos sqlite
// package database
//
// import (
// 	"jos3lo89/go-image-service/models"
// 	"log"
// 	"os"
//
// 	"github.com/gofiber/fiber/v2/middleware/logger"
// 	"gorm.io/driver/sqlite"
// 	"gorm.io/gorm"
// )
//
// var DB *gorm.DB
//
// func ConnectDB(dbPath  string) {
// 	var err error
//
// 	// db, err := gorm.Open(sqlite.Open(config.AppConfig.DBPath), &gorm.Config{
// 	// 	Logger: logger.Default.LogMode(logger.Info),
// 	// })
//
//
// 	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
// 		Logger: logger.Default.LogMode(logger.Info),
// 	})
// 	if err != nil {
// 		log.Fatalf("Error al conectar con la base de datos: %v", err)
// 		os.Exit(1)
// 	}
//
// 	log.Println("Conexión a la base de datos exitosa.")
// 	DB.AutoMigrate(&models.User{}, models.Image{})
// 	log.Println("Migración de la base de datos completada.")
// }

// Package database gestiona la conexión a la base de datos y la migración de modelos.
package database

import (
	"jos3lo89/go-image-service/models"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDB inicializa la conexión con la base de datos SQLite y realiza la migración automática.
func ConnectDB(dbPath string) {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
		os.Exit(1)
	}

	log.Println("Conexión a la base de datos exitosa.")
	DB.AutoMigrate(&models.User{}, &models.Image{})
	log.Println("Migración de la base de datos completada.")
}
