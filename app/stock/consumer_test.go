//go:build consumer

package stock_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gofiber/fiber/v2"
	"github.com/pact-cdc-example/basket-service/pkg/cerr"
	"github.com/pact-cdc-example/basket-service/pkg/httpclient"

	"github.com/pact-cdc-example/basket-service/app/stock"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/suite"
)

const (
	consumerBasketService = "BasketService"
	providerStockService  = "StockService"
)

const (
	isProductAvailableInStockPath = "/api/v1/stocks/availability"
)

type StockConsumerTestSuite struct {
	suite.Suite
	client        stock.Client
	pact          *dsl.Pact
	pactServerURL string
}

func TestStockConsumerTestSuite(t *testing.T) {
	suite.Run(t, new(StockConsumerTestSuite))
}

func (s *StockConsumerTestSuite) SetupSuite() {
	s.initPact()

	s.client = stock.NewClient(&stock.NewClientOpts{
		HTTPClient: httpclient.New(),
		BaseURL:    s.pactServerURL,
	})
}

func (s *StockConsumerTestSuite) TearDownSuite() {
	defer s.pact.Teardown()
}

func (s *StockConsumerTestSuite) TestGivenStockInquiryForProductReqThenItShouldReturnProductIDMustBeGivenErrWhenProductIDIsNotGiven() {
	quantity := int(gofakeit.Uint8())

	s.pact.
		AddInteraction().
		Given("i get product id must be given error if product id is not given").
		UponReceiving("A request for inquiry stock information about a product").
		WithRequest(dsl.Request{
			Method: http.MethodPost,
			Path:   dsl.String(fmt.Sprintf(isProductAvailableInStockPath)),
			Headers: map[string]dsl.Matcher{
				fiber.HeaderContentType: dsl.String(fiber.MIMEApplicationJSON),
				fiber.HeaderAccept:      dsl.String(fiber.MIMEApplicationJSON),
			},
			Body: dsl.StructMatcher{
				"quantity": dsl.Like(quantity),
			},
		}).
		WillRespondWith(dsl.Response{
			Status: http.StatusBadRequest,
			Headers: map[string]dsl.Matcher{
				fiber.HeaderContentType: dsl.String(fiber.MIMEApplicationJSON),
			},
			Body: dsl.StructMatcher{
				"code":    30000,
				"message": "Product id must be given to stock inquiry.",
			},
		})

	var test = func() error {
		_, err := s.client.IsProductAvailableInStock(context.Background(), stock.IsProductAvailableInStockRequest{
			ProductID: nil,
			Quantity:  &quantity,
		})

		return err
	}

	err := s.pact.Verify(test)

	s.Equal(err, cerr.Bag{Code: 30000, Message: "Product id must be given to stock inquiry."})
}

func (s *StockConsumerTestSuite) TestGivenStockInquiryForProductReqThenItShouldReturnQuantityMustBeGivenErrWhenQuantityIsNotGiven() {
	productID := gofakeit.UUID()

	s.pact.
		AddInteraction().
		Given("i get quantity must be given error if quantity is not given").
		UponReceiving("A request for inquiry stock information about a product").
		WithRequest(dsl.Request{
			Method: http.MethodPost,
			Path:   dsl.String(fmt.Sprintf(isProductAvailableInStockPath)),
			Headers: map[string]dsl.Matcher{
				fiber.HeaderContentType: dsl.String(fiber.MIMEApplicationJSON),
				fiber.HeaderAccept:      dsl.String(fiber.MIMEApplicationJSON),
			},
			Body: dsl.StructMatcher{
				"product_id": dsl.Like(productID),
			},
		}).
		WillRespondWith(dsl.Response{
			Status: http.StatusBadRequest,
			Headers: map[string]dsl.Matcher{
				fiber.HeaderContentType: dsl.String(fiber.MIMEApplicationJSON),
			},
			Body: dsl.StructMatcher{
				"code":    30002,
				"message": "Quantity must be given to stock inquiry.",
			},
		})

	var test = func() error {
		_, err := s.client.IsProductAvailableInStock(context.Background(), stock.IsProductAvailableInStockRequest{
			ProductID: &productID,
			Quantity:  nil,
		})

		return err
	}

	err := s.pact.Verify(test)

	s.Equal(err, cerr.Bag{Code: 30002, Message: "Quantity must be given to stock inquiry."})
}

func (s *StockConsumerTestSuite) TestGivenStockInquiryForProductReqThenItShouldReturnNoStockInfoFoundErrorWhenGivenProductIDNotHasStockInfo() {
	givenProductID := gofakeit.UUID()
	quantity := int(gofakeit.Uint8())

	s.pact.
		AddInteraction().
		Given("i get no stock information found error if no stock information found for given product id").
		UponReceiving("A request for inquiry stock information about a product").
		WithRequest(dsl.Request{
			Method: http.MethodPost,
			Path:   dsl.String(fmt.Sprintf(isProductAvailableInStockPath)),
			Headers: map[string]dsl.Matcher{
				fiber.HeaderContentType: dsl.String(fiber.MIMEApplicationJSON),
				fiber.HeaderAccept:      dsl.String(fiber.MIMEApplicationJSON),
			},
			Body: dsl.StructMatcher{
				"product_id": dsl.Like(givenProductID),
				"quantity":   dsl.Like(quantity),
			},
		}).
		WillRespondWith(dsl.Response{
			Status: http.StatusBadRequest,
			Headers: map[string]dsl.Matcher{
				fiber.HeaderContentType: dsl.String(fiber.MIMEApplicationJSON),
			},
			Body: dsl.StructMatcher{
				"code":    30001,
				"message": "No stock information found for given product id.",
			},
		})

	var test = func() error {
		_, err := s.client.IsProductAvailableInStock(context.Background(), stock.IsProductAvailableInStockRequest{
			ProductID: &givenProductID,
			Quantity:  &quantity,
		})

		return err
	}

	err := s.pact.Verify(test)

	s.Equal(err, cerr.Bag{Code: 30001, Message: "No stock information found for given product id."})
}

func (s *StockConsumerTestSuite) TestGivenStockInquiryForProductReqThenItShouldReturnFalseIfGivenProductIDNotInStockInGivenQuantity() {
	givenProductID := gofakeit.UUID()
	quantity := int(gofakeit.Uint8())

	s.pact.
		AddInteraction().
		Given("i get false").
		UponReceiving("A request for inquiry stock information about a product").
		WithRequest(dsl.Request{
			Method: http.MethodPost,
			Path:   dsl.String(fmt.Sprintf(isProductAvailableInStockPath)),
			Headers: map[string]dsl.Matcher{
				fiber.HeaderContentType: dsl.String(fiber.MIMEApplicationJSON),
				fiber.HeaderAccept:      dsl.String(fiber.MIMEApplicationJSON),
			},
			Body: dsl.StructMatcher{
				"product_id": dsl.Like(givenProductID),
				"quantity":   dsl.Like(quantity),
			},
		}).
		WillRespondWith(dsl.Response{
			Status: http.StatusOK,
			Headers: map[string]dsl.Matcher{
				fiber.HeaderContentType: dsl.String(fiber.MIMEApplicationJSON),
			},
			Body: dsl.StructMatcher{
				"is_available": false,
			},
		})

	var test = func() error {
		_, err := s.client.IsProductAvailableInStock(context.Background(), stock.IsProductAvailableInStockRequest{
			ProductID: &givenProductID,
			Quantity:  &quantity,
		})
		return err
	}

	err := s.pact.Verify(test)

	s.Nil(err)
}

func (s *StockConsumerTestSuite) TestGivenStockInquiryForProductReqThenItShouldReturnTrueIfGivenProductIDAvailableInStockInGivenQuantity() {
	givenProductID := gofakeit.UUID()
	quantity := int(gofakeit.Uint8())

	s.pact.
		AddInteraction().
		Given("i get true").
		UponReceiving("A request for inquiry stock information about a product").
		WithRequest(dsl.Request{
			Method: http.MethodPost,
			Path:   dsl.String(fmt.Sprintf(isProductAvailableInStockPath)),
			Headers: map[string]dsl.Matcher{
				fiber.HeaderContentType: dsl.String(fiber.MIMEApplicationJSON),
				fiber.HeaderAccept:      dsl.String(fiber.MIMEApplicationJSON),
			},
			Body: dsl.StructMatcher{
				"product_id": dsl.Like(givenProductID),
				"quantity":   dsl.Like(quantity),
			},
		}).
		WillRespondWith(dsl.Response{
			Status: http.StatusOK,
			Headers: map[string]dsl.Matcher{
				fiber.HeaderContentType: dsl.String(fiber.MIMEApplicationJSON),
			},
			Body: dsl.StructMatcher{
				"is_available": true,
			},
		})

	var test = func() error {
		_, err := s.client.IsProductAvailableInStock(context.Background(), stock.IsProductAvailableInStockRequest{
			ProductID: &givenProductID,
			Quantity:  &quantity,
		})
		return err
	}

	err := s.pact.Verify(test)

	s.Nil(err)
}

func (s *StockConsumerTestSuite) initPact() {
	s.pact = &dsl.Pact{
		Host:                     "127.0.0.1",
		Consumer:                 consumerBasketService,
		Provider:                 providerStockService,
		DisableToolValidityCheck: true,
		PactFileWriteMode:        "overwrite",
		LogDir:                   "./pacts/logs",
	}
	//it must be used otherwise it could not create pact file
	s.pact.Setup(true)

	s.pactServerURL = fmt.Sprintf("http://localhost:%d", s.pact.Server.Port)
}
