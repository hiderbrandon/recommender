package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"recommender/internal/core/domain"
	"recommender/internal/core/ports"
	"strings"
	"time"
)

type ExternalStockAPI struct {
	client  *http.Client
	baseURL string
	apiKey  string
}

func NewExternalStockAPI() ports.StockAPIClient {
	apiURL := strings.TrimSpace(os.Getenv("API_URL"))
	if apiURL == "" {
		panic("API_URL environment variable not set")
	}
	apiKey := strings.TrimSpace(os.Getenv("API_KEY"))
	if apiKey == "" {
		panic("API_KEY environment variable not set")
	}
	return &ExternalStockAPI{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: apiURL,
		apiKey:  apiKey,
	}
}

func (a *ExternalStockAPI) FetchStocks(nextPage string) (*domain.APIResponse, error) {
	url := a.baseURL
	if nextPage != "" {
		url = fmt.Sprintf("%s?next_page=%s", a.baseURL, nextPage)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Agregar encabezados, incluyendo User-Agent similar al de curl
	req.Header.Set("Authorization", "Bearer "+a.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "curl/7.68.0") // Opcional, para imitar la petición de curl

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var apiResponse domain.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, err
	}

	return &apiResponse, nil
}
