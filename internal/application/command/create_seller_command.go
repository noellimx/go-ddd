package command

import "github.com/noellimx/go-ddd/internal/application/common"

type CreateSellerCommand struct {
	IdempotencyKey string
	Name           string
}

type CreateSellerCommandResult struct {
	Result *common.SellerResult
}
