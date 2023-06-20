package stock

type IsProductAvailableInStockRequest struct {
	ProductID *string `json:"product_id,omitempty"`
	Quantity  *int    `json:"quantity,omitempty"`
}

type ReserveStockRequest struct {
	ProductID string `json:"product_id,omitempty"`
	Quantity  int    `json:"quantity,omitempty"`
}
