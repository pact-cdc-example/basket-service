//go:build consumer

package product_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gofiber/fiber/v2"
	"github.com/pact-cdc-example/basket-service/pkg/httpclient"

	"github.com/pact-cdc-example/basket-service/app/product"
	"github.com/pact-cdc-example/basket-service/pkg/cerr"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/suite"
)

type ProductConsumerTestSuite struct {
	suite.Suite
	pact          *dsl.Pact
	pactServerURL string
	client        product.Client
}

const (
	providerProductService = "ProductService"
	consumerBasketService  = "BasketService"
)

const (
	getProductByIDPath   = "/api/v1/products/%s"
	getProductsByIDsPath = "/api/v1/products/bulk"
)

func (s *ProductConsumerTestSuite) SetupSuite() {
	s.initPact()

	s.client = product.NewClient(&product.NewClientOpts{
		HTTPClient: httpclient.New(),
		BaseURL:    s.pactServerURL,
	})
}

func (s *ProductConsumerTestSuite) TearDownSuite() {
	defer s.pact.Teardown()
}

func TestProductConsumerTestSuite(t *testing.T) {
	suite.Run(t, new(ProductConsumerTestSuite))
}

func (s *ProductConsumerTestSuite) TestGivenGetProductByIDRequestThenItShouldReturnProductNotFoundErrorWhenGivenProductIDNotExists() {
	givenProductID := gofakeit.UUID()

	s.pact.AddInteraction().
		Given("i get product not found error when the product with given id does not exists").
		UponReceiving("A request for product with a non exist product id").
		WithRequest(dsl.Request{
			Method:  http.MethodGet,
			Path:    dsl.String(fmt.Sprintf(getProductByIDPath, givenProductID)),
			Headers: map[string]dsl.Matcher{},
		}).
		WillRespondWith(dsl.Response{
			Status: http.StatusBadRequest,
			Headers: map[string]dsl.Matcher{
				fiber.HeaderContentType: dsl.String(fiber.MIMEApplicationJSON),
			},
			Body: dsl.StructMatcher{
				"code":    20001,
				"message": "Product not found.",
			},
		})

	var test = func() error {
		_, err := s.client.GetProductByID(context.Background(), givenProductID)
		return err
	}

	err := s.pact.Verify(test)

	s.Equal(err, cerr.Bag{Code: 20001, Message: "Product not found."})
}

func (s *ProductConsumerTestSuite) TestGivenGetProductByIDRequestThenItShouldReturnProductWhenGivenProductIDExists() {
	givenProductID := gofakeit.UUID()

	givenProduct := product.Product{
		ID:        givenProductID,
		Name:      gofakeit.Name(),
		Code:      gofakeit.Word(),
		Color:     gofakeit.Color(),
		CreatedAt: gofakeit.Date(),
		UpdatedAt: gofakeit.Date(),
		Price:     gofakeit.Price(10, 100),
		ImageURL:  gofakeit.ImageURL(200, 100),
		Type:      gofakeit.Word(),
	}

	s.pact.AddInteraction().
		Given("i get product with given id").
		UponReceiving("A request for product with a exist product id").
		WithRequest(dsl.Request{
			Method: http.MethodGet,
			Path:   dsl.String(fmt.Sprintf(getProductByIDPath, givenProductID)),
		}).
		WillRespondWith(dsl.Response{
			Status: http.StatusOK,
			Headers: map[string]dsl.Matcher{
				fiber.HeaderContentType: dsl.String(fiber.MIMEApplicationJSON),
			},
			Body: dsl.StructMatcher{
				"id":         dsl.Like(givenProduct.ID),
				"name":       dsl.Like(givenProduct.Name),
				"code":       dsl.Like(givenProduct.Code),
				"color":      dsl.Like(givenProduct.Color),
				"created_at": dsl.Like(givenProduct.CreatedAt),
				"updated_at": dsl.Like(givenProduct.UpdatedAt),
				"price":      dsl.Like(givenProduct.Price),
				"image_url":  dsl.Like(givenProduct.ImageURL),
				"type":       dsl.Like(givenProduct.Type),
			},
		})

	var test = func() error {
		_, err := s.client.GetProductByID(context.Background(), givenProductID)
		return err
	}

	err := s.pact.Verify(test)

	s.Nil(err)
}

func (s *ProductConsumerTestSuite) TestGivenGetProductsByIDsReqThenItShouldReturnBodyParserErrorWhenNoProductIDIsSent() {
	s.pact.AddInteraction().
		Given("i get body parser error when no product id is given").
		UponReceiving("A request for get products").
		WithRequest(dsl.Request{
			Method: http.MethodPost,
			Path:   dsl.String(getProductsByIDsPath),
			Headers: map[string]dsl.Matcher{
				fiber.HeaderContentType: dsl.String(fiber.MIMEApplicationJSON),
				fiber.HeaderAccept:      dsl.String(fiber.MIMEApplicationJSON),
			},
		}).
		WillRespondWith(dsl.Response{
			Status: http.StatusBadRequest,
			Headers: map[string]dsl.Matcher{
				fiber.HeaderContentType: dsl.String(fiber.MIMEApplicationJSON),
			},
			Body: dsl.StructMatcher{
				"code":    10001,
				"message": "could not parse request body.",
			},
		})

	var test = func() error {
		_, err := s.client.GetProductsByIDs(context.Background(), product.GetProductByIDsRequest{})
		return err
	}

	err := s.pact.Verify(test)

	s.Equal(err, cerr.Bag{Code: 10001, Message: "could not parse request body."})
}

func (s *ProductConsumerTestSuite) TestGivenGetProductsByIDsReqThenItShouldReturnProductNotFoundErrorWhenOneOrMoreGivenProductIDNotExists() {
	givenProductIDs := []string{gofakeit.UUID(), gofakeit.UUID(), gofakeit.UUID()}

	givenReq := product.GetProductByIDsRequest{
		IDs: givenProductIDs,
	}

	s.pact.AddInteraction().
		Given("i get product not found error when the one of product with given id does not exists").
		UponReceiving("A request for get products contains at least one not exist product id").
		WithRequest(dsl.Request{
			Method: http.MethodPost,
			Path:   dsl.String(getProductsByIDsPath),
			Headers: map[string]dsl.Matcher{
				fiber.HeaderContentType: dsl.String(fiber.MIMEApplicationJSON),
				fiber.HeaderAccept:      dsl.String(fiber.MIMEApplicationJSON),
			},
			Body: map[string]interface{}{
				"ids": dsl.EachLike(givenProductIDs[0], 1),
			},
		}).
		WillRespondWith(dsl.Response{
			Status: http.StatusBadRequest,
			Headers: map[string]dsl.Matcher{
				fiber.HeaderContentType: dsl.String(fiber.MIMEApplicationJSON),
			},
			Body: dsl.StructMatcher{
				"code":    20003,
				"message": "At least one of given product ids does not exist.",
			},
		})

	var test = func() error {
		_, err := s.client.GetProductsByIDs(context.Background(), givenReq)
		return err
	}

	err := s.pact.Verify(test)
	s.Equal(err, cerr.Bag{Code: 20003, Message: "At least one of given product ids does not exist."})
}

func (s *ProductConsumerTestSuite) NoTestGivenGetProductsByIDsReqThenItShouldReturnProductsWhenAllGivenProductIDsExists() {
	givenProductIDs := []string{gofakeit.UUID(), gofakeit.UUID(), gofakeit.UUID()}

	givenReq := product.GetProductByIDsRequest{
		IDs: givenProductIDs,
	}

	givenProducts := []product.Product{
		{
			ID:        givenProductIDs[0],
			Name:      gofakeit.Name(),
			Code:      gofakeit.Word(),
			Color:     gofakeit.Color(),
			CreatedAt: gofakeit.Date(),
			UpdatedAt: gofakeit.Date(),
			Price:     gofakeit.Price(10, 100),
			ImageURL:  gofakeit.ImageURL(200, 100),
			Type:      gofakeit.Word(),
		},
		{
			ID:        givenProductIDs[1],
			Name:      gofakeit.Name(),
			Code:      gofakeit.Word(),
			Color:     gofakeit.Color(),
			CreatedAt: gofakeit.Date(),
			UpdatedAt: gofakeit.Date(),
			Price:     gofakeit.Price(10, 100),
			ImageURL:  gofakeit.ImageURL(200, 100),
			Type:      gofakeit.Word(),
		},
		{
			ID:        givenProductIDs[2],
			Name:      gofakeit.Name(),
			Code:      gofakeit.Word(),
			Color:     gofakeit.Color(),
			CreatedAt: gofakeit.Date(),
			UpdatedAt: gofakeit.Date(),
			Price:     gofakeit.Price(10, 100),
			ImageURL:  gofakeit.ImageURL(200, 100),
			Type:      gofakeit.Word(),
		},
	}

	s.pact.AddInteraction().
		Given("i get products with given ids").
		UponReceiving("A request for get products with given ids").
		WithRequest(dsl.Request{
			Method: http.MethodGet,
			Path:   dsl.String(getProductsByIDsPath),
			Body: dsl.StructMatcher{
				"ids": dsl.Like(givenProductIDs),
			},
		}).
		WillRespondWith(dsl.Response{
			Status: http.StatusOK,
			Headers: map[string]dsl.Matcher{
				fiber.HeaderContentType: dsl.String(fiber.MIMEApplicationJSON),
			},
			Body: dsl.StructMatcher{
				"products": dsl.EachLike(dsl.StructMatcher{
					"id":        dsl.Like(givenProducts[0].ID),
					"name":      dsl.Like(givenProducts[0].Name),
					"code":      dsl.Like(givenProducts[0].Code),
					"color":     dsl.Like(givenProducts[0].Color),
					"createdAt": dsl.Like(givenProducts[0].CreatedAt),
					"updatedAt": dsl.Like(givenProducts[0].UpdatedAt),
					"price":     dsl.Like(givenProducts[0].Price),
					"imageURL":  dsl.Like(givenProducts[0].ImageURL),
					"type":      dsl.Like(givenProducts[0].Type),
				}, len(givenProducts)),
			},
		})

	var test = func() error {
		_, err := s.client.GetProductsByIDs(context.Background(), givenReq)
		return err
	}

	err := s.pact.Verify(test)

	s.Nil(err)
}

func (s *ProductConsumerTestSuite) initPact() {
	s.pact = &dsl.Pact{
		Host:                     "127.0.0.1",
		Consumer:                 consumerBasketService,
		Provider:                 providerProductService,
		DisableToolValidityCheck: true,
		PactFileWriteMode:        "overwrite",
		LogDir:                   "./pacts/logs",
		PactDir:                  ".././pacts",
	}
	//it must be used otherwise it could not create pact file
	s.pact.Setup(true)

	s.pactServerURL = fmt.Sprintf("http://localhost:%d", s.pact.Server.Port)
}
