package basket

import (
	"github.com/pact-cdc-example/basket-service/app/product"
)

type GetBasketResponse struct {
	ID        string                `json:"id"`
	UserID    string                `json:"user_id"`
	Products  []ProductQuantityPair `json:"products,omitempty"`
	CreatedAt string                `json:"created_at"`
	UpdatedAt string                `json:"updated_at"`
}

func NewBasketResponse(basket *Basket, products []product.Product) *GetBasketResponse {
	if basket == nil {
		return nil
	}

	productQuantityPairs := make([]ProductQuantityPair, len(basket.Products))
	for i := range basket.Products {
		for j := range products {
			if basket.Products[i].ID == products[j].ID {
				productQuantityPairs[i] = ProductQuantityPair{
					Product:  &products[j],
					Quantity: basket.Products[i].Quantity,
				}
			}
		}
	}

	return &GetBasketResponse{
		ID:        basket.ID,
		UserID:    basket.UserID,
		CreatedAt: basket.CreatedAt.Format(layoutISO),
		UpdatedAt: basket.UpdatedAt.Format(layoutISO),
		Products:  productQuantityPairs,
	}
}

type ProductQuantityPair struct {
	Product  *product.Product `json:"product,omitempty"`
	Quantity int              `json:"quantity"`
}
