package query

import "github.com/noellimx/go-ddd/internal/application/common"

type GetAllProductsQueryResult struct {
	Result []*common.ProductResult
}
