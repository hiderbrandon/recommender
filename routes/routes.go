package routes

import (
	"recommender/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/stocks", handlers.GetStocks)
	r.POST("/stocks", handlers.PostStock)
	return r
}
