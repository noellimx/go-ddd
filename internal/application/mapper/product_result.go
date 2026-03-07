package mapper

import (
	"github.com/noellimx/go-ddd/internal/application/common"
	"github.com/noellimx/go-ddd/internal/domain/entities"
)

func NewProductResultFromValidatedEntity(product *entities.ValidatedProduct) *common.ProductResult {
	return NewProductResultFromEntity(&product.Product)
}

func NewProductResultFromEntity(product *entities.Product) *common.ProductResult {
	if product == nil {
		return nil
	}

	return &common.ProductResult{
		Id:        product.Id,
		Name:      product.Name,
		Price:     product.Price,
		Seller:    NewSellerResultFromEntity(&product.Seller),
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
}
