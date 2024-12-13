package dto

type CreateTransactionRequest struct {
	CustomerID        int    `json:"customer_id"`
	OnTheRoadPrice    int    `json:"on_the_road_price" validate:"required,numeric,amount_number"`
	InstallmentAmount int    `json:"installment_amount" validate:"required,numeric,amount_number"`
	InterestAmount    int    `json:"interest_amount" validate:"required,numeric,amount_number"`
	AssetName         string `json:"asset_name" validate:"required,valid_text,max=100"`
	TenorMonth        int    `json:"tenor_month" validate:"required,numeric,amount_number"`
}
