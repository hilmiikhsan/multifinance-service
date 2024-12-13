package repository

const (
	queryInsertNewUser = `
		INSERT INTO customers
		(
			nik,
			email,
			password,
			full_name,
			legal_name,
			birth_place,
			birth_date,
			salary,
			ktp_photo_path,
			selfie_photo_path
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ? 
		)
	`

	queryFindCustomerByEmail = `
		SELECT
			id,
			nik,
			email,
			password
		FROM customers
		WHERE email = ?
	`

	queryFindCustomer = `
		SELECT id, email FROM customers WHERE id = ?
	`

	queryFindCustomerByID = `
		SELECT
			c.id,
			c.nik,
			c.email,
			c.full_name,
			c.legal_name,
			c.birth_place,
			c.birth_date,
			c.salary,
			c.ktp_photo_path,
			c.selfie_photo_path,
			c.created_at,
			c.updated_at,
			cl.tenor_month,
			cl.limit_amount
		FROM customers c
		LEFT JOIN credit_limits cl ON c.id = cl.customer_id
		WHERE c.id = ?
	`
)
