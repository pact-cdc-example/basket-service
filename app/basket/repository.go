package basket

import "context"

type Repository interface {
	CreateBasket(ctx context.Context, basket *Basket) (*Basket, error)
	GetBasketByID(ctx context.Context, basketID string) (*Basket, error)
	AddProductToBasket(ctx context.Context, product *Product) (*Basket, error)
}
