package repository

const (
	queryInsertNewCreditLimit = `
		INSERT INTO credit_limitS (customer_id, tenor_month, limit_amount) VALUES (?, ?, ?)
	`

	queryFindCreditLimitByCustomerID = `
		SELECT
			tenor_month,
			limit_amount
		FROM credit_limits
		WHERE customer_id = ?
	`
)
