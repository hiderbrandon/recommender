package repository


import (
    "recommender/config"
    "recommender/internal/core/domain"
    "recommender/internal/core/ports"
)

type CockroachStockRepository struct{}

func NewCockroachStockRepository() port.StockRepository {
    return &CockroachStockRepository{}
}

func (r *CockroachStockRepository) GetAll() ([]domain.Stock, error) {
    var stocks []domain.Stock
    result := config.DB.Find(&stocks)
    return stocks, result.Error
}

func (r *CockroachStockRepository) Create(stock *domain.Stock) error {
    return config.DB.Create(stock).Error
}
