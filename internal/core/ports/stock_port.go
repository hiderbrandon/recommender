package port

import "recommender/internal/core/domain"

type StockRepository interface {
	GetAll() ([]domain.Stock, error)
	Create(stock *domain.Stock) error
}
