package dto

type CreditLimit struct {
	Tenor       int     `json:"tenor"`
	LimitAmount float64 `json:"limit_amount"`
}
