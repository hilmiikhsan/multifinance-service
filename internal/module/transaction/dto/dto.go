package dto

import "github.com/hilmiikhsan/multifinance-service/pkg/types"

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

type HistoryListTransactionItem struct {
	ID                int     `json:"id" db:"id"`
	CustomerID        int     `json:"customer_Id" db:"customer_id"`
	ContractNumber    string  `json:"contract_number" db:"contract_number"`
	OnTheRoadPrice    float64 `json:"on_the_road_price" db:"on_the_road_price"`
	AdminFee          float64 `json:"admin_fee" db:"admin_fee"`
	InstallmentAmount float64 `json:"installment_amount" db:"installment_amount"`
	InterestAmount    float64 `json:"interest_amount" db:"interest_amount"`
	AssetName         string  `json:"asset_name" db:"asset_name"`
	CreatedAt         string  `json:"created_at" db:"created_at"`
}

type GetHistoryListTransactionRequest struct {
	Page     int `query:"page" validate:"required,min=1"`
	Paginate int `query:"paginate" validate:"required,min=1,max=100"`
}

type GetHistoryListTransactionResponse struct {
	Items []HistoryListTransactionItem `json:"items"`
	Meta  types.Meta                   `json:"meta"`
}

func (r *GetHistoryListTransactionRequest) SetDefault() {
	if r.Page < 1 {
		r.Page = 1
	}

	if r.Paginate < 1 {
		r.Paginate = 10
	}
}
