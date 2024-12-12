package seeds

import (
	"fmt"

	"github.com/hilmiikhsan/multifinance-service/pkg/utils"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type Seed struct {
	db *sqlx.DB
}

// NewSeed initializes a new Seed instance with a database connection.
func newSeed(db *sqlx.DB) Seed {
	return Seed{db: db}
}

// Execute runs the seeder for the specified table with the given number of entries.
func Execute(db *sqlx.DB, table string, total int) {
	seed := newSeed(db)
	seed.run(table, total)
}

// run handles seeding based on the table name.
func (s *Seed) run(table string, total int) {
	switch table {
	case "customers":
		s.customersSeed(total)
	case "credit_limits":
		s.creditLimitsSeed()
	case "transactions":
		s.transactionsSeed()
	case "all":
		s.customersSeed(total)
		s.creditLimitsSeed()
		s.transactionsSeed()
	case "delete-all":
		s.deleteAll()
	default:
		log.Warn().Msg("No seed to run")
	}
}

// customersSeed seeds the customers table with the specified number of entries.
func (s *Seed) customersSeed(total int) {
	log.Info().Msg("Seeding customers table...")

	if s.db == nil {
		log.Error().Msg("Database connection is nil")
		return
	}

	for i := 0; i < total; i++ {
		query := `
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
			)`

		nik := fmt.Sprintf("%016d", i+1) // Generate unique NIK with 16 digits
		email := fmt.Sprintf("user%d@gmail.com", i+1)
		password, _ := utils.HashPassword("password") // Default password
		fullName := fmt.Sprintf("User %d", i+1)
		legalName := fmt.Sprintf("User Legal %d", i+1)
		birthPlace := "City"
		birthDate := "1990-01-01"
		salary := 5000000.00 // Default salary
		ktpPhotoPath := fmt.Sprintf("/path/to/ktp/user%d.jpg", i+1)
		selfiePhotoPath := fmt.Sprintf("/path/to/selfie/user%d.jpg", i+1)

		_, err := s.db.Exec(query, nik, email, password, fullName, legalName, birthPlace, birthDate, salary, ktpPhotoPath, selfiePhotoPath)
		if err != nil {
			log.Error().Err(err).Msg("Error seeding customers table")
			return
		}
	}

	log.Info().Msg("Customers table seeded successfully")
}

// creditLimitsSeed seeds the credit_limits table with the specified number of entries.
func (s *Seed) creditLimitsSeed() {
	log.Info().Msg("Seeding credit_limits table...")

	// Query to get all customer_ids from the customers table
	rows, err := s.db.Query("SELECT id FROM customers")
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving customer_ids from customers table")
		return
	}
	defer rows.Close()

	// Prepare the tenor_month and limit_amount arrays
	tenorMonths := []int{1, 2, 3, 4}
	limitAmounts := []float64{100000, 200000, 500000, 700000}

	// Iterate over the customers and insert data into credit_limits
	for rows.Next() {
		var customerId int64
		if err := rows.Scan(&customerId); err != nil {
			log.Error().Err(err).Msg("Error scanning customer_id")
			return
		}

		// Insert data for each tenor_month and limit_amount
		for i := 0; i < len(tenorMonths); i++ {
			query := `INSERT INTO credit_limits (customer_id, tenor_month, limit_amount) VALUES (?, ?, ?)`
			_, err := s.db.Exec(query, customerId, tenorMonths[i], limitAmounts[i])
			if err != nil {
				log.Error().Err(err).Msg("Error seeding credit_limits table")
				return
			}
		}
	}

	log.Info().Msg("Credit_limits table seeded successfully")
}

// transactionsSeed seeds the transactions table with the specified number of entries.
func (s *Seed) transactionsSeed() {
	log.Info().Msg("Seeding transactions table...")

	// Query to get all customer_ids from the customers table
	rows, err := s.db.Query("SELECT id FROM customers")
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving customer_ids from customers table")
		return
	}
	defer rows.Close()

	// Iterate over customers and insert data into transactions
	for rows.Next() {
		var customerId int64
		if err := rows.Scan(&customerId); err != nil {
			log.Error().Err(err).Msg("Error scanning customer_id")
			return
		}

		// Generate dummy data for transaction fields
		contractNumber := fmt.Sprintf("CN%010d", customerId)
		onTheRoadPrice := 5000000.0   // Example price
		adminFee := 50000.0           // Example admin fee
		installmentAmount := 200000.0 // Example installment amount
		interestAmount := 15000.0     // Example interest amount
		assetName := fmt.Sprintf("Asset_%d", customerId)

		query := `INSERT INTO transactions (customer_id, contract_number, on_the_road_price, admin_fee, installment_amount, interest_amount, asset_name) 
                  VALUES (?, ?, ?, ?, ?, ?, ?)`
		_, err := s.db.Exec(query, customerId, contractNumber, onTheRoadPrice, adminFee, installmentAmount, interestAmount, assetName)
		if err != nil {
			log.Error().Err(err).Msg("Error seeding transactions table")
			return
		}
	}

	log.Info().Msg("Transactions table seeded successfully")
}

// deleteAll deletes all data from the seeded tables.
func (s *Seed) deleteAll() {
	log.Info().Msg("Deleting all data from seeded tables...")

	tables := []string{"transactions", "credit_limits", "users"}
	for _, table := range tables {
		query := "DELETE FROM " + table
		_, err := s.db.Exec(query)
		if err != nil {
			log.Error().Err(err).Str("table", table).Msg("Error deleting data")
			return
		}
	}

	log.Info().Msg("All data deleted successfully")
}
