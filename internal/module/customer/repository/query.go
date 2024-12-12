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
)
