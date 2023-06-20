package basket

import "errors"

type CreateBasketRequest struct {
	UserID string `json:"user_id"`
}

func (cbr CreateBasketRequest) Validate() error {
	if cbr.UserID == "" {
		return errors.New("user id can not be empty")
	}

	return nil
}

type AddProductToBasketRequest struct {
	BasketID  string `json:"basket_id"`
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type AddBulkProductToBasketRequest struct {
	UserID   string `json:"user_id"`
	BasketID string `json:"basket_id"`
	Products []struct {
		ID       string `json:"id"`
		Quantity int    `json:"quantity"`
	} `json:"products"`
}
