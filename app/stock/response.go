package stock

import "time"

type IsProductAvailableInStockResponse struct {
	IsAvailable bool `json:"is_available,omitempty"`
}

type Stock struct {
	ID               string    `json:"id"`
	ProductID        string    `json:"product_id"`
	Quantity         int       `json:"quantity"`
	ReservedQuantity int       `json:"reserved_quantity"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
