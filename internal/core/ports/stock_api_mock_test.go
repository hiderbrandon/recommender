package ports_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"recommender/internal/core/domain"
)

type MockStockAPIClient struct {
	mock.Mock
}

func (m *MockStockAPIClient) FetchStocks(nextPage string) (*domain.APIResponse, error) {
	args := m.Called(nextPage)
	return args.Get(0).(*domain.APIResponse), args.Error(1)
}

func TestFetchStocksSuccess(t *testing.T) {
	mockClient := new(MockStockAPIClient)

	expected := &domain.APIResponse{
		Items: []domain.Stock{{
			ID:         1,
			Ticker:     "AAPL",
			Company:    "Apple Inc.",
			Brokerage:  "Goldman Sachs",
			Action:     "Buy",
			RatingFrom: "Neutral",
			RatingTo:   "Buy",
			TargetFrom: 140.0,
			TargetTo:   160.0,
			Time:       time.Now(),
		}},
		NextPage: "page2",
	}

	mockClient.On("FetchStocks", "").Return(expected, nil)

	resp, err := mockClient.FetchStocks("")

	assert.NoError(t, err)
	assert.Equal(t, "AAPL", resp.Items[0].Ticker)
	assert.Equal(t, "Apple Inc.", resp.Items[0].Company)
	assert.Equal(t, "page2", resp.NextPage)

	mockClient.AssertExpectations(t)
}

func TestFetchStocksError(t *testing.T) {
	mockClient := new(MockStockAPIClient)

	mockClient.On("FetchStocks", "").Return(&domain.APIResponse{}, errors.New("api error"))

	resp, err := mockClient.FetchStocks("")

	assert.Error(t, err)
	assert.Equal(t, "api error", err.Error())
	assert.Empty(t, resp.Items)

	mockClient.AssertExpectations(t)
}
