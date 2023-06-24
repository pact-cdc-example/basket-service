package persistence

import (
	"context"
	"database/sql"

	"github.com/pact-cdc-example/basket-service/app/basket"

	"github.com/sirupsen/logrus"
)

type PostgresRepository interface {
	CreateBasket(
		ctx context.Context, bask *basket.Basket) (*basket.Basket, error)
	AddProductToBasket(
		ctx context.Context, product *basket.Product) (*basket.Basket, error)
	GetBasketByID(
		ctx context.Context, basketID string) (*basket.Basket, error)
}

type postgresRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

type NewPostgresRepositoryOpts struct {
	DB *sql.DB
	L  *logrus.Logger
}

func NewPostgresRepository(opts *NewPostgresRepositoryOpts) PostgresRepository {
	return &postgresRepository{
		db:     opts.DB,
		logger: opts.L,
	}
}

func (pr *postgresRepository) CreateBasket(
	ctx context.Context, bask *basket.Basket) (*basket.Basket, error) {
	err := pr.db.QueryRowContext(ctx,
		`INSERT INTO baskets (id, user_id)
		 VALUES ($1, $2) RETURNING created_at`,
		bask.ID, bask.UserID,
	).Err()

	if err != nil {
		pr.logger.Errorf("could not create basket :%v", err)
		return nil, err
	}

	return pr.getBasketByID(ctx, bask.ID)
}

func (pr *postgresRepository) getBasketByID(
	ctx context.Context, basketID string) (*basket.Basket, error) {

	row := pr.db.QueryRowContext(ctx,
		`SELECT id, user_id, created_at, updated_at 
		FROM baskets WHERE ID = $1`, basketID,
	)

	var bask basket.Basket
	if err := row.Scan(
		&bask.ID,
		&bask.UserID,
		&bask.CreatedAt,
		&bask.UpdatedAt,
	); err != nil {
		pr.logger.Errorf("could not get basket by id: %v", err)
		return nil, err
	}

	rows, err := pr.db.QueryContext(ctx,
		`SELECT product_id, quantity FROM basket_products
		WHERE basket_id = $1`, basketID)

	if err != nil {
		pr.logger.Errorf("could not get basket products: %v", err)
		return nil, err
	}

	var products []basket.Product
	for rows.Next() {
		var product basket.Product
		if err := rows.Scan(
			&product.ID,
			&product.Quantity,
		); err != nil {
			pr.logger.Errorf("could not scan basket product: %v", err)
			return nil, err
		}
		products = append(products, product)
	}

	bask.Products = products

	return &bask, nil
}

func (pr *postgresRepository) AddProductToBasket(
	ctx context.Context, product *basket.Product) (*basket.Basket, error) {
	err := pr.db.QueryRowContext(ctx,
		`INSERT INTO basket_products (product_id, quantity, basket_id)
		 VALUES ($1, $2, $3)`,
		product.ID, product.Quantity, product.BasketID,
	).Err()

	if err != nil {
		pr.logger.Errorf("could not add product to basket: %v", err)
		return nil, err
	}

	return pr.getBasketByID(ctx, product.BasketID)
}

func (pr *postgresRepository) GetBasketByID(
	ctx context.Context, basketID string) (*basket.Basket, error) {
	return pr.getBasketByID(ctx, basketID)
}
