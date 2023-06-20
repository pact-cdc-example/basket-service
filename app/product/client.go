package product

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pact-cdc-example/basket-service/pkg/httpclient"
)

type Client interface {
	GetProductByID(ctx context.Context, id string) (*Product, error)
	GetProductsByIDs(ctx context.Context, req GetProductByIDsRequest) ([]Product, error)
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

const (
	getProductByIDPath   = "%s/api/v1/products/%s"
	getProductsByIDsPath = "%s/api/v1/products/bulk"
)

func (c *client) GetProductByID(ctx context.Context, id string) (*Product, error) {
	url := fmt.Sprintf(getProductByIDPath, c.baseURL, id)

	resBytes, err := c.httpClient.Get(ctx, url, c.headers)
	if err != nil {
		return nil, err
	}

	var resp GetProductResponse
	if err := json.Unmarshal(resBytes, &resp); err != nil {
		return nil, err
	}

	return &resp.Product, nil
}

func (c *client) GetProductsByIDs(
	ctx context.Context, req GetProductByIDsRequest) ([]Product, error) {
	url := fmt.Sprintf(getProductsByIDsPath, c.baseURL)

	resBytes, err := c.httpClient.Post(ctx, url, c.headers, req)
	if err != nil {
		return nil, err
	}

	var resp GetProductsResponse
	if err := json.Unmarshal(resBytes, &resp); err != nil {
		return nil, err
	}

	return resp.Products, nil
}
