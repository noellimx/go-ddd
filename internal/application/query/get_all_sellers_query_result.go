package query

import "github.com/noellimx/go-ddd/internal/application/common"

type GetAllSellersQueryResult struct {
	Result []*common.SellerResult
}
