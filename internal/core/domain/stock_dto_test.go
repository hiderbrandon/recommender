package domain_test

import (
	"encoding/json"
	"testing"

	"recommender/internal/core/domain"

	"github.com/stretchr/testify/assert"
)

func TestStockDTO_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"ticker": "AAPL",
		"target_from": "150.00",
		"target_to": "180.00",
		"company": "Apple Inc.",
		"brokerage": "Goldman Sachs",
		"action": "Buy",
		"rating_from": "Neutral",
		"rating_to": "Buy",
		"time": "2024-06-01T15:04:05Z"
	}`

	var dto domain.StockDTO
	err := json.Unmarshal([]byte(jsonData), &dto)

	assert.NoError(t, err)
	assert.Equal(t, "AAPL", dto.Ticker)
	assert.Equal(t, "150.00", dto.TargetFrom)
	assert.Equal(t, "180.00", dto.TargetTo)
	assert.Equal(t, "Apple Inc.", dto.Company)
	assert.Equal(t, "Goldman Sachs", dto.Brokerage)
	assert.Equal(t, "Buy", dto.Action)
	assert.Equal(t, "Neutral", dto.RatingFrom)
	assert.Equal(t, "Buy", dto.RatingTo)
	assert.Equal(t, "2024-06-01T15:04:05Z", dto.Time)
}

func TestAPIResponseDTO_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"items": [{
			"ticker": "GOOGL",
			"target_from": "2500",
			"target_to": "2800",
			"company": "Google LLC",
			"brokerage": "JP Morgan",
			"action": "Hold",
			"rating_from": "Sell",
			"rating_to": "Hold",
			"time": "2024-05-01T10:00:00Z"
		}],
		"next_page": "page_2"
	}`

	var response domain.APIResponseDTO
	err := json.Unmarshal([]byte(jsonData), &response)

	assert.NoError(t, err)
	assert.Len(t, response.Items, 1)
	assert.Equal(t, "GOOGL", response.Items[0].Ticker)
	assert.Equal(t, "page_2", response.NextPage)
}