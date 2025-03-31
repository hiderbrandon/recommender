package domain

type StockDTO struct {
	Ticker     string `json:"ticker"`
	TargetFrom string `json:"target_from"` // String como viene de la API
	TargetTo   string `json:"target_to"`   // String como viene de la API
	Company    string `json:"company"`
	Brokerage  string `json:"brokerage"`
	Action     string `json:"action"`
	RatingFrom string `json:"rating_from"`
	RatingTo   string `json:"rating_to"`
	Time       string `json:"time"` // String antes de parsear a `time.Time`
}

// APIResponseDTO representa la estructura de respuesta de la API externa
type APIResponseDTO struct {
	Items    []StockDTO `json:"items"`
	NextPage string     `json:"next_page"`
}
