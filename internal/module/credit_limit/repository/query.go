package repository

const (
	queryInsertNewCreditLimit = `
		INSERT INTO credit_limitS (customer_id, tenor_month, limit_amount) VALUES (?, ?, ?)
	`
)
