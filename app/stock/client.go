package stock

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pact-cdc-example/basket-service/pkg/httpclient"
)

const (
	isProductAvailableInStockPath = "%s/api/v1/stocks/availability"
	reserveStockPath              = "%s/api/v1/stocks/reserve"
)

type Client interface {
	IsProductAvailableInStock(
		ctx context.Context, req IsProductAvailableInStockRequest) (bool, error)
	ReserveStock(
		ctx context.Context, req ReserveStockRequest) (*Stock, error)
}

type client struct {
	httpClient httpclient.Client
	headers    map[string]string
	baseURL    string
}

type NewClientOpts struct {
	HTTPClient httpclient.Client
	BaseURL    string
}

func NewClient(opts *NewClientOpts) Client {
	return &client{
		httpClient: opts.HTTPClient,
		headers:    httpclient.DefaultHeaders,
		baseURL:    opts.BaseURL,
	}
}

func (c *client) IsProductAvailableInStock(ctx context.Context, req IsProductAvailableInStockRequest) (bool, error) {
	url := fmt.Sprintf(isProductAvailableInStockPath, c.baseURL)

	body, err := c.httpClient.GetWithBody(ctx, url, c.headers, req)
	if err != nil {
		return false, err
	}

	var resp IsProductAvailableInStockResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return false, err
	}

	return resp.IsAvailable, nil
}

func (c *client) ReserveStock(
	ctx context.Context, req ReserveStockRequest) (*Stock, error) {
	url := fmt.Sprintf(reserveStockPath, c.baseURL)

	body, err := c.httpClient.Put(ctx, url, c.headers, req)
	if err != nil {
		return nil, err
	}

	var resp Stock
	if err = json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
