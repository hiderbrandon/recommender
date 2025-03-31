package ports

import "recommender/internal/core/domain"

type APIResponse struct {
	Items    []domain.Stock `json:"items"`
	NextPage string         `json:"next_page"`
}

type StockAPIClient interface {
	FetchStocks(nextPage string) (*domain.APIResponse, error)
}
