package basket

import (
	"github.com/pact-cdc-example/basket-service/pkg/cerr"
)

// basket specific errorr

const (
	BasketNotFoundErrCode           cerr.Code = 10100
	ProductNotHasEnoughStockErrCode cerr.Code = 10101
)
