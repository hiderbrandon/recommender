package config

import (
    "fmt"
    "log"
    "os"

    "github.com/joho/godotenv"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "recommender/models"
)

var DB *gorm.DB

func InitDB() {
    // Cargar variables de entorno
    if err := godotenv.Load(); err != nil {
        log.Println("⚠ No se pudo cargar el archivo .env, usando variables del sistema")
    }

    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("SSL_MODE"),
    )

    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("❌ Error conectando a la base de datos:", err)
    }

    log.Println("✅ Conectado a la base de datos")

    // Migraciones automáticas
    err = DB.AutoMigrate(&models.Stock{})
    if err != nil {
        log.Fatal("❌ Error al migrar la base de datos:", err)
    }
    log.Println("✅ Migraciones completadas")
}
