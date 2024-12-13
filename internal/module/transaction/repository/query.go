package repository

const (
	queryInsertNewTransaction = `
		INSERT INTO transactions
		(
			customer_id,
			contract_number,
			on_the_road_price,
			admin_fee,
			installment_amount,
			interest_amount,
			asset_name
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`
)
