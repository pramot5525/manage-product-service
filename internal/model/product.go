package model

type Product struct {
	ID          int64    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	SalePrice   *float64 `json:"sale_price"`
	Price       float64  `json:"price"`
}

type ProductRequest struct {
	Name        string   `json:"name" validate:"required,omitempty"`
	Description *string  `json:"description" validate:"omitempty"`
	SalePrice   *float64 `json:"sale_price" validate:"omitempty"`
	Price       float64  `json:"price" validate:"required,omitempty"`
}

type APIResponse struct {
	Successful bool     `json:"successful"`
	ErrorCode  string   `json:"error_code"`
	Data       *Product `json:"data,omitempty"`
}
