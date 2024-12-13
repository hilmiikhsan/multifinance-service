package dto

type CreateTransactionRequest struct {
	CustomerID        int    `json:"customer_id"`
	OnTheRoadPrice    int    `json:"on_the_road_price" validate:"required,numeric,amount_number"`
	InstallmentAmount int    `json:"installment_amount" validate:"required,numeric,amount_number"`
	InterestAmount    int    `json:"interest_amount" validate:"required,numeric,amount_number"`
	AssetName         string `json:"asset_name" validate:"required,valid_text,max=100"`
	TenorMonth        int    `json:"tenor_month" validate:"required,numeric,amount_number"`
}

type GetDetailTransactionResponse struct {
	ID                int     `json:"id"`
	CustomerID        int     `json:"customer_id"`
	ContractNumber    string  `json:"contract_number"`
	OnTheRoadPrice    float64 `json:"on_the_road_price"`
	AdminFee          float64 `json:"admin_fee"`
	InstallmentAmount float64 `json:"installment_amount"`
	InterestAmount    float64 `json:"interest_amount"`
	AssetName         string  `json:"asset_name"`
	CreatedAt         string  `json:"created_at"`
}
