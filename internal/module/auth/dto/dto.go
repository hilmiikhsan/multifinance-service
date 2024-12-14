package dto

type RegisterRequest struct {
	Nik             string `json:"nik" validate:"required,max=16,nik"`
	Email           string `json:"email" validate:"required,email,email_blacklist"`
	Password        string `json:"password" validate:"required,strong_password"`
	FullName        string `json:"full_name" validate:"required,max=100,valid_text"`
	LegalName       string `json:"legal_name" validate:"required,max=100,valid_text"`
	BirthPlace      string `json:"birth_place" validate:"required,max=100,valid_text"`
	BirthDate       string `json:"birth_date" validate:"required,birth_date"`
	Salary          int    `json:"salary" validate:"required,numeric,amount_number"`
	KtpPhotoPath    string `json:"ktp_photo_path" validate:"required,file_path"`
	SelfiePhotoPath string `json:"selfie_photo_path" validate:"required,file_path"`
}

type RegisterResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,email_blacklist"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	ID           int64  `json:"id"`
	Email        string `json:"email"`
	FullName     string `json:"full_name"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	Token string `json:"token"`
}
