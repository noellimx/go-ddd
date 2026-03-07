package query

import (
	"github.com/google/uuid"
	"github.com/noellimx/go-ddd/internal/application/common"
)

type GetProductByIdQuery struct {
	Id uuid.UUID
}

type GetProductByIdQueryResult struct {
	Result *common.ProductResult
}
