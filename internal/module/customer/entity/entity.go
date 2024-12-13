package entity

import (
	"database/sql"
	"time"

	"github.com/hilmiikhsan/multifinance-service/internal/module/credit_limit/entity"
)

type Customer struct {
	ID              int             `db:"id"`
	Nik             string          `db:"nik"`
	Email           string          `db:"email"`
	Password        string          `db:"password"`
	FullName        string          `db:"full_name"`
	LegalName       string          `db:"legal_name"`
	BirthPlace      string          `db:"birth_place"`
	BirthDate       time.Time       `db:"birth_date"`
	Salary          float64         `db:"salary"`
	KtpPhotoPath    string          `db:"ktp_photo_path"`
	SelfiePhotoPath string          `db:"selfie_photo_path"`
	TenorMonth      int             `db:"tenor_month"`
	LimitAmount     float64         `db:"limit_amount"`
	Limits          []entity.Limits `db:"-"`
	CreatedAt       time.Time       `db:"created_at"`
	UpdatedAt       time.Time       `db:"updated_at"`
}

type CustomerWithLimits struct {
	CustomerID      int             `db:"id"`
	Nik             string          `db:"nik"`
	Email           string          `db:"email"`
	FullName        string          `db:"full_name"`
	LegalName       string          `db:"legal_name"`
	BirthPlace      string          `db:"birth_place"`
	BirthDate       time.Time       `db:"birth_date"`
	Salary          float64         `db:"salary"`
	KtpPhotoPath    string          `db:"ktp_photo_path"`
	SelfiePhotoPath string          `db:"selfie_photo_path"`
	CreatedAt       time.Time       `db:"created_at"`
	UpdatedAt       time.Time       `db:"updated_at"`
	TenorMonth      sql.NullInt64   `db:"tenor_month"`
	LimitAmount     sql.NullFloat64 `db:"limit_amount"`
}
