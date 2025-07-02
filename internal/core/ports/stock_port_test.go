package ports

import (
	"errors"
	"recommender/internal/core/domain"
	"testing"
	"time"
)

// Mock implementation of StockRepository for testing
type mockStockRepository struct {
	stocks []domain.Stock
}

func (m *mockStockRepository) GetAll(limit, offset int) ([]domain.Stock, error) {
	if offset > len(m.stocks) {
		return []domain.Stock{}, nil
	}
	end := offset + limit
	if end > len(m.stocks) {
		end = len(m.stocks)
	}
	return m.stocks[offset:end], nil
}

func (m *mockStockRepository) Create(stock *domain.Stock) error {
	m.stocks = append(m.stocks, *stock)
	return nil
}

func (m *mockStockRepository) GetStockByTickerAndTime(ticker string, t time.Time) (*domain.Stock, error) {
	for _, s := range m.stocks {
		if s.Ticker == ticker && s.Time.Equal(t) {
			return &s, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockStockRepository) GetTopStocksByTarget(limit int) ([]domain.Stock, error) {
	if limit > len(m.stocks) {
		limit = len(m.stocks)
	}
	return m.stocks[:limit], nil
}

func (m *mockStockRepository) GetStockByTicker(ticker string) (*domain.Stock, error) {
	for _, s := range m.stocks {
		if s.Ticker == ticker {
			return &s, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockStockRepository) GetRecentStocks(limit int) ([]domain.Stock, error) {
	if limit > len(m.stocks) {
		limit = len(m.stocks)
	}
	return m.stocks[:limit], nil
}

func TestStockRepository_GetAll(t *testing.T) {
	repo := &mockStockRepository{
		stocks: []domain.Stock{
			{Ticker: "AAPL"},
			{Ticker: "MSFT"},
			{Ticker: "TSLA"},
		},
	}
	stocks, err := repo.GetAll(2, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stocks) != 2 || stocks[0].Ticker != "MSFT" {
		t.Errorf("unexpected stocks: %+v", stocks)
	}
}

func TestStockRepository_CreateAndGetStockByTicker(t *testing.T) {
	repo := &mockStockRepository{}
	stock := &domain.Stock{Ticker: "GOOG"}
	err := repo.Create(stock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := repo.GetStockByTicker("GOOG")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Ticker != "GOOG" {
		t.Errorf("expected ticker GOOG, got %s", got.Ticker)
	}
}

func TestStockRepository_GetStockByTickerAndTime(t *testing.T) {
	now := time.Now()
	repo := &mockStockRepository{
		stocks: []domain.Stock{
			{Ticker: "NFLX", Time: now},
		},
	}
	stock, err := repo.GetStockByTickerAndTime("NFLX", now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stock.Ticker != "NFLX" {
		t.Errorf("expected ticker NFLX, got %s", stock.Ticker)
	}
}

func TestStockRepository_GetTopStocksByTarget(t *testing.T) {
	repo := &mockStockRepository{
		stocks: []domain.Stock{
			{Ticker: "A"},
			{Ticker: "B"},
			{Ticker: "C"},
		},
	}
	top, err := repo.GetTopStocksByTarget(2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(top) != 2 {
		t.Errorf("expected 2 stocks, got %d", len(top))
	}
}

func TestStockRepository_GetRecentStocks(t *testing.T) {
	repo := &mockStockRepository{
		stocks: []domain.Stock{
			{Ticker: "X"},
			{Ticker: "Y"},
			{Ticker: "Z"},
		},
	}
	recent, err := repo.GetRecentStocks(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recent) != 1 || recent[0].Ticker != "X" {
		t.Errorf("unexpected recent stocks: %+v", recent)
	}
}