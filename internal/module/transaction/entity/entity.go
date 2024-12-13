package entity

type Transaction struct {
	ID                int     `db:"id"`
	CustomerID        int     `db:"customer_id"`
	ContractNumber    string  `db:"contract_number"`
	OnTheRoadPrice    float64 `db:"on_the_road_price"`
	AdminFee          float64 `db:"admin_fee"`
	InstallmentAmount float64 `db:"installment_amount"`
	InterestAmount    float64 `db:"interest_amount"`
	AssetName         string  `db:"asset_name"`
	CreatedAt         string  `db:"created_at"`
	UpdatedAt         string  `db:"updated_at"`
}
