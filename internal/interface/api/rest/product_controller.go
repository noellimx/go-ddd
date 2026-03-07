package rest

import (
	"encoding/json"
	"net/http"

	"github.com/noellimx/go-ddd/internal/application/interfaces"
	"github.com/noellimx/go-ddd/internal/application/query"
	"github.com/noellimx/go-ddd/internal/interface/api/rest/dto/mapper"
	"github.com/noellimx/go-ddd/internal/interface/api/rest/dto/request"

	"github.com/google/uuid"
)

type ProductController struct {
	service interfaces.ProductService
}

func NewProductController(service interfaces.ProductService) *ProductController {
	controller := &ProductController{
		service: service,
	}

	return controller
}

func (pc *ProductController) CreateProductController(w http.ResponseWriter, r *http.Request) {
	var createProductRequest request.CreateProductRequest

	if err := json.NewDecoder(r.Body).Decode(&createProductRequest); err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Failed to parse request body",
		})
		return
	}

	productCommand, err := createProductRequest.ToCreateProductCommand()
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid product Id format",
		})
		return
	}

	result, err := pc.service.CreateProduct(productCommand)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to create product",
		})
		return
	}

	response := mapper.ToProductResponse(result.Result)

	WriteJSON(w, http.StatusCreated, response)
}

func (pc *ProductController) GetAllProductsController(w http.ResponseWriter, r *http.Request) {
	products, err := pc.service.FindAllProducts()
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch products",
		})
		return
	}

	response := mapper.ToProductListResponse(products.Result)

	WriteJSON(w, http.StatusOK, response)
}

func (pc *ProductController) GetProductByIdController(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid product Id format",
		})
		return
	}

	product, err := pc.service.FindProductById(&query.GetProductByIdQuery{Id: id})
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch product",
		})
		return
	}

	if product == nil {
		WriteJSON(w, http.StatusNotFound, map[string]string{
			"error": "Product not found",
		})
		return
	}

	response := mapper.ToProductResponse(product.Result)

	WriteJSON(w, http.StatusOK, response)
}
