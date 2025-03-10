package routes

import (
    "github.com/gin-gonic/gin"
    "recommender/internal/adapters/handlers"
)

func SetupRouter(stockHandler *handlers.StockHandler) *gin.Engine {
    r := gin.Default()

    r.GET("/stocks", stockHandler.GetStocks)
    r.POST("/stocks", stockHandler.PostStock)
    r.GET("/stocks/recommendations", stockHandler.GetRecommendations) 
    r.GET("/stocks/:ticker", stockHandler.GetStockByTicker)



    return r
}
