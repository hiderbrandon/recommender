package main

import (
	"log"
	"os"

	"recommender/config"
	"recommender/routes"

	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno desde .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠ No se pudo cargar el archivo .env, usando variables del sistema")
	}

	config.InitDB()
	r := routes.SetupRouter()

	// Obtener el puerto desde las variables de entorno
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8081" // Puerto por defecto si no se encuentra en .env
	}

	log.Println("🚀 Servidor corriendo en el puerto", port)
	r.Run(":" + port)
}
