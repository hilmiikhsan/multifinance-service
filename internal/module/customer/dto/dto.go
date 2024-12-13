package dto

import "github.com/hilmiikhsan/multifinance-service/internal/module/credit_limit/dto"

type GetCustomerProfileResponse struct {
	ID              int               `json:"id"`
	Nik             string            `json:"nik"`
	FullName        string            `json:"full_name"`
	LegalName       string            `json:"legal_name"`
	BirthPlace      string            `json:"birth_place"`
	BirthDate       string            `json:"birth_date"`
	Salary          float64           `json:"salary"`
	KtpPhotoPath    string            `json:"ktp_photo_path"`
	SelfiePhotoPath string            `json:"selfie_photo_path"`
	Limits          []dto.CreditLimit `json:"limits"`
	CreatedAt       string            `json:"created_at"`
	UpdatedAt       string            `json:"updated_at"`
}
