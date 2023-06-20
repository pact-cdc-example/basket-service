package basket

import "time"

type Basket struct {
	ID        string    `json:"-"`
	UserID    string    `json:"-"`
	Products  []Product `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type Product struct {
	ID        string    `json:"-"`
	Quantity  int       `json:"-"`
	BasketID  string    `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func getIDsOfProducts(products []Product) []string {
	ids := make([]string, len(products))
	for i, p := range products {
		ids[i] = p.ID
	}

	return ids
}
