package main

import (
	"log"

	"github.com/pact-cdc-example/basket-service/app/basket"
	"github.com/pact-cdc-example/basket-service/app/persistence"
	"github.com/pact-cdc-example/basket-service/app/product"
	"github.com/pact-cdc-example/basket-service/app/stock"
	"github.com/pact-cdc-example/basket-service/config"
	"github.com/pact-cdc-example/basket-service/pkg/httpclient"
	"github.com/pact-cdc-example/basket-service/pkg/postgres"
	"github.com/pact-cdc-example/basket-service/pkg/server"
	"github.com/sirupsen/logrus"
)

func main() {
	c := config.New()

	db := postgres.New(&postgres.NewPostgresOpts{
		Host:     c.Postgres().Host,
		Port:     c.Postgres().Port,
		DBName:   c.Postgres().DBName,
		Password: c.Postgres().Password,
		Username: c.Postgres().Username,
	})

	logger := logrus.New()

	repository := persistence.NewPostgresRepository(&persistence.NewPostgresRepositoryOpts{
		DB: db,
		L:  logger,
	})

	httpClient := httpclient.New()

	productClient := product.NewClient(&product.NewClientOpts{
		HTTPClient: httpClient,
		BaseURL:    c.ExternalURL().ProductAPI,
	})

	stockClient := stock.NewClient(&stock.NewClientOpts{
		HTTPClient: httpClient,
		BaseURL:    c.ExternalURL().StockAPI,
	})

	basketService := basket.NewService(&basket.NewServiceOpts{
		R: repository, L: logger, PC: productClient, SC: stockClient,
	})

	basketHandler := basket.NewHandler(&basket.NewHandlerOpts{
		S: basketService, L: logger,
	})

	app := server.New(&server.NewServerOpts{
		Port: c.Server().Port,
	}, []server.RouteHandler{
		basketHandler,
	})

	if err := app.Run(); err != nil {
		log.Fatalf("server is closed: %v", err)
	}
}
