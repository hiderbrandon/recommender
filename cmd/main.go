package main

import (
	"log"
	"os"

	"recommender/config"
	"recommender/internal/adapters/clients"
	"recommender/internal/adapters/handlers"

	repository "recommender/internal/adapters/repositories"
	"recommender/internal/core/services"
	"recommender/routes"

	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno desde .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠ No se pudo cargar el archivo .env, usando variables del sistema")
	}

	db := config.InitDB()

	// Crear instancia del adaptador para la API externa
	apiClient := clients.NewExternalStockAPI()

	// Inyección de dependencias
	stockRepo := repository.NewCockroachStockRepository(db)
	stockService := services.NewStockService(stockRepo, apiClient)
	stockHandler := handlers.NewStockHandler(stockService)

	// Ejecutar la importación de datos solo una vez al inicio
	err := stockService.FetchAndStoreStocks()
	if err != nil {
		log.Println("Error importing stocks:", err)
	}

	r := routes.SetupRouter(stockHandler)

	// Obtener el puerto desde las variables de entorno
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8081" // Puerto por defecto si no se encuentra en .env
	}

	log.Println("🚀 Servidor corriendo en el puerto", port)
	r.Run(":" + port)
}
