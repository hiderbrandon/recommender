package services

import (
    "recommender/models"
    "recommender/internal/repositories"
)

func FetchStocks() ([]models.Stock, error) {
    return repositories.GetAllStocks()
}

func AddStock(stock *models.Stock) error {
    return repositories.CreateStock(stock)
}
