package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"recommender/internal/core/ports"
	"time"
)

type ExternalStockAPI struct {
	client  *http.Client
	baseURL string
	apiKey  string
}

func NewExternalStockAPI() ports.StockAPIClient {

	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		panic("API_URL environment variable not set")
	}
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		panic("API_KEY environment variable not set")
	}
	return &ExternalStockAPI{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: apiURL,
		apiKey:  apiKey,
	}
}

func (a *ExternalStockAPI) FetchStocks(nextPage string) (*ports.APIResponse, error) {
	url := a.baseURL
	if nextPage != "" {
		url = fmt.Sprintf("%s?next_page=%s", a.baseURL, nextPage)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+a.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var apiResponse ports.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, err
	}

	return &apiResponse, nil
}
