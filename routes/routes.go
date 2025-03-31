package routes

import (
	"log"
	"recommender/internal/adapters/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(stockHandler *handlers.StockHandler) *gin.Engine {
	r := gin.Default()

	// Configurar CORS para aceptar cualquier origen
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 🔥 Permitir cualquier origen
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false, // No permitir credenciales por seguridad
	}))

	log.Println("✅ CORS configurado para permitir cualquier origen.")

	// Definir rutas
	r.GET("/stocks", stockHandler.GetStocks)
	r.POST("/stocks", stockHandler.PostStock)
	r.GET("/stocks/recommendations", stockHandler.GetRecommendations)
	r.GET("/stocks/:ticker", stockHandler.GetStockByTicker)

	return r
}
