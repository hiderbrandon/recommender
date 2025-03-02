package repository

import (
	"recommender/internal/core/domain"
	port "recommender/internal/core/ports"

	"gorm.io/gorm"
)

type CockroachStockRepository struct {
	db *gorm.DB
}

func NewCockroachStockRepository(db *gorm.DB) port.StockRepository {
	return &CockroachStockRepository{db: db}
}

func (r *CockroachStockRepository) GetAll() ([]domain.Stock, error) {
	var stocks []domain.Stock
	result := r.db.Find(&stocks)
	return stocks, result.Error
}

func (r *CockroachStockRepository) Create(stock *domain.Stock) error {
	return r.db.Create(stock).Error
}
