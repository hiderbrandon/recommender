package services

import (
	"errors"
	"testing"
	"time"

	"gorm.io/gorm"
	"recommender/internal/core/domain"
)

// --- Mocks ---

type mockStockRepository struct {
	stocks []domain.Stock
}

func (m *mockStockRepository) GetTopStocksByTarget(limit int) ([]domain.Stock, error) {
	// Solo un stub para cumplir con la interfaz
	return []domain.Stock{}, nil
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
	return nil, gorm.ErrRecordNotFound
}
func (m *mockStockRepository) GetRecentStocks(limit int) ([]domain.Stock, error) {
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

type mockStockAPIClient struct {
	responses []*domain.APIResponse
	calls     int
}

func (m *mockStockAPIClient) FetchStocks(nextPage string) (*domain.APIResponse, error) {
	if m.calls < len(m.responses) {
		resp := m.responses[m.calls]
		m.calls++
		return resp, nil
	}
	return nil, errors.New("no more pages")
}

// --- Tests ---

func TestFetchAndStoreStocks(t *testing.T) {
	fixedTime := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

	mockAPI := &mockStockAPIClient{
		responses: []*domain.APIResponse{
			{
				Items: []domain.Stock{
					{
						Ticker:     "TSLA",
						TargetFrom: 200,
						TargetTo:   240,
						RatingFrom: "Neutral",
						RatingTo:   "Buy",
						Brokerage:  "The Goldman Sachs Group",
						Time:       fixedTime,
					},
				},
				NextPage: "",
			},
		},
	}

	repo := &mockStockRepository{}
	service := NewStockService(repo, mockAPI)

	err := service.FetchAndStoreStocks()
	if err != nil {
		t.Fatalf("FetchAndStoreStocks failed: %v", err)
	}

	if len(repo.stocks) != 1 {
		t.Errorf("expected 1 stock stored, got %d", len(repo.stocks))
	}
	if repo.stocks[0].Ticker != "TSLA" {
		t.Errorf("expected ticker TSLA, got %s", repo.stocks[0].Ticker)
	}
}

func TestAddStockAndFetchStocks(t *testing.T) {
	repo := &mockStockRepository{}
	service := NewStockService(repo, nil)
	fixedTime := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)


	stock := &domain.Stock{
		Ticker:     "AAPL",
		TargetFrom: 100,
		TargetTo:   120,
		RatingFrom: "Neutral",
		RatingTo:   "Buy",
		Brokerage:  "JP Morgan",
		Time:       fixedTime,
	}

	err := service.AddStock(stock)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stocks, err := service.FetchStocks(10, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stocks) != 1 || stocks[0].Ticker != "AAPL" {
		t.Errorf("expected to fetch the added stock")
	}
}

func TestGetTopRecommendedStocks(t *testing.T) {
	repo := &mockStockRepository{
		stocks: []domain.Stock{
			{Ticker: "A", TargetFrom: 100, TargetTo: 120, RatingFrom: "Sell", RatingTo: "Buy", Brokerage: "JP Morgan"},
			{Ticker: "B", TargetFrom: 100, TargetTo: 110, RatingFrom: "Neutral", RatingTo: "Buy", Brokerage: "Others"},
			{Ticker: "C", TargetFrom: 100, TargetTo: 90, RatingFrom: "Buy", RatingTo: "Sell", Brokerage: "Morgan Stanley"},
		},
	}
	service := NewStockService(repo, nil)

	top, err := service.GetTopRecommendedStocks(2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(top) != 2 {
		t.Errorf("expected 2 top stocks, got %d", len(top))
	}
	if top[0].Ticker != "A" {
		t.Errorf("expected stock 'A' to be top, got %s", top[0].Ticker)
	}
}

func TestGetStockByTicker(t *testing.T) {
	repo := &mockStockRepository{
		stocks: []domain.Stock{
			{Ticker: "MSFT"},
		},
	}
	service := NewStockService(repo, nil)

	stock, err := service.GetStockByTicker("MSFT")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stock.Ticker != "MSFT" {
		t.Errorf("expected ticker MSFT, got %s", stock.Ticker)
	}
}