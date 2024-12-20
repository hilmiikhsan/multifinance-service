package repository

import (
	"context"
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/hilmiikhsan/multifinance-service/constants"
	creditLimitEntity "github.com/hilmiikhsan/multifinance-service/internal/module/credit_limit/entity"
	"github.com/hilmiikhsan/multifinance-service/internal/module/customer/entity"
	"github.com/hilmiikhsan/multifinance-service/internal/module/customer/ports"
	"github.com/hilmiikhsan/multifinance-service/pkg/err_msg"
	"github.com/hilmiikhsan/multifinance-service/pkg/utils"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.CustomerRepository = &customerRepository{}

type customerRepository struct {
	db *sqlx.DB
}

func NewCustomerRepository(db *sqlx.DB) *customerRepository {
	return &customerRepository{
		db: db,
	}
}

func (r *customerRepository) InsertNewUser(ctx context.Context, tx *sql.Tx, data *entity.Customer) (*entity.Customer, error) {
	var res = new(entity.Customer)

	result, err := tx.ExecContext(ctx, r.db.Rebind(queryInsertNewUser),
		data.Nik,
		data.Email,
		data.Password,
		data.FullName,
		data.LegalName,
		data.BirthPlace,
		data.BirthDate,
		data.Salary,
		data.KtpPhotoPath,
		data.SelfiePhotoPath,
	)
	if err != nil {
		uniqueConstraints := map[string]string{
			"nik":   constants.ErrNikAlreadyRegistered,
			"email": constants.ErrEmailAlreadyRegistered,
		}

		val, handleErr := utils.HandleInsertUniqueError(err, data, uniqueConstraints)
		if handleErr != nil {
			log.Error().Err(handleErr).Any("payload", data).Msg("repository::InsertNewUser - Failed to insert new user")
			return nil, handleErr
		}

		if customer, ok := val.(*entity.Customer); ok {
			log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewUser - Failed to insert new user")
			return customer, nil
		}

		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))

	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Error().Err(err).Msg("repository::InsertNewUser - Failed to retrieve last inserted ID")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	err = tx.QueryRowContext(ctx, queryFindCustomer, lastInsertID).Scan(&res.ID, &res.Email)
	if err != nil {
		log.Error().Err(err).Msg("repository::InsertNewUser - Failed to fetch inserted user details")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return res, nil
}

func (r *customerRepository) FindCustomerByEmail(ctx context.Context, email string) (*entity.Customer, error) {
	var res = new(entity.Customer)

	err := r.db.GetContext(ctx, res, r.db.Rebind(queryFindCustomerByEmail), email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Any("email", email).Msg("repository::FindCustomerByEmail - Email not found")
			return nil, nil
		}

		log.Error().Err(err).Any("email", email).Msg("repository::FindCustomerByEmail - Failed to find user by email")
		return nil, err
	}

	return res, nil
}

func (r *customerRepository) FindCustomerByID(ctx context.Context, id int) (*entity.Customer, error) {
	var rows []entity.CustomerWithLimits

	err := r.db.SelectContext(ctx, &rows, r.db.Rebind(queryFindCustomerByID), id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Int("id", id).Msg("repository::FindCustomerByID - ID not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrUserNotFound))
		}
		log.Error().Err(err).Int("id", id).Msg("repository::FindCustomerByID - Failed to find user by ID")
		return nil, err
	}

	if len(rows) == 0 {
		return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrUserNotFound))
	}

	customer := &entity.Customer{
		ID:              rows[0].CustomerID,
		Nik:             rows[0].Nik,
		Email:           rows[0].Email,
		FullName:        rows[0].FullName,
		LegalName:       rows[0].LegalName,
		BirthPlace:      rows[0].BirthPlace,
		BirthDate:       rows[0].BirthDate,
		Salary:          rows[0].Salary,
		KtpPhotoPath:    rows[0].KtpPhotoPath,
		SelfiePhotoPath: rows[0].SelfiePhotoPath,
		CreatedAt:       rows[0].CreatedAt,
		UpdatedAt:       rows[0].UpdatedAt,
	}

	for _, row := range rows {
		if row.TenorMonth.Valid && row.LimitAmount.Valid {
			customer.Limits = append(customer.Limits, creditLimitEntity.Limits{
				TenorMonth:  int(row.TenorMonth.Int64),
				LimitAmount: row.LimitAmount.Float64,
			})
		}
	}

	return customer, nil
}
