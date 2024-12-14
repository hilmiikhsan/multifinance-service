package repository

const (
	queryInsertNewCreditLimit = `
		INSERT INTO credit_limits (customer_id, tenor_month, limit_amount) VALUES (?, ?, ?)
	`

	queryFindCreditLimitByCustomerID = `
		SELECT
			tenor_month,
			limit_amount
		FROM credit_limits
		WHERE customer_id = ?
	`

	queryLockCreditLimitByCustomerAndTenor = `
		SELECT
			tenor_month,
			limit_amount
		FROM credit_limits
		WHERE customer_id = ? AND tenor_month = ?
		FOR UPDATE
	`
)
