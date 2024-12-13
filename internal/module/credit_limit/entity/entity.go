package entity

type CreditLimit struct {
	CustomerID  int     `db:"customer_id"`
	TenorMonth  int     `db:"tenor_month"`
	LimitAmount float64 `db:"limit_amount"`
}

type Limits struct {
	TenorMonth  int     `db:"tenor_month"`
	LimitAmount float64 `db:"limit_amount"`
}
