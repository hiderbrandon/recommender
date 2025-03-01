package domain

import "time"

type Stock struct {
    ID         uint      `json:"id"`
    Ticker     string    `json:"ticker"`
    Company    string    `json:"company"`
    Brokerage  string    `json:"brokerage"`
    Action     string    `json:"action"`
    RatingFrom string    `json:"rating_from"`
    RatingTo   string    `json:"rating_to"`
    TargetFrom float64   `json:"target_from"`
    TargetTo   float64   `json:"target_to"`
    Time       time.Time `json:"time"`
}
