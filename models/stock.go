package models

import "time"

type Stock struct {
    ID         uint      `gorm:"primaryKey" json:"id"`
    Ticker     string    `json:"ticker"`
    Company    string    `json:"company"`
    Brokerage  string    `json:"brokerage"` // Nuevo campo
    Action     string    `json:"action"`
    RatingFrom string    `json:"rating_from"`
    RatingTo   string    `json:"rating_to"`
    TargetFrom float64   `json:"target_from"` // Cambiado de string a float64
    TargetTo   float64   `json:"target_to"`   // Cambiado de string a float64
    Time       time.Time `json:"time"`        // Agregado para reflejar la fecha de recomendación
}
