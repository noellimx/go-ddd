package rest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/noellimx/go-ddd/internal/application/command"
	"github.com/noellimx/go-ddd/internal/application/common"
	"github.com/noellimx/go-ddd/internal/application/interfaces"
	"github.com/noellimx/go-ddd/internal/application/query"
	"github.com/noellimx/go-ddd/internal/domain/entities"
	"github.com/noellimx/go-ddd/internal/interface/api/rest/dto/response"
	"github.com/stretchr/testify/mock"

	"github.com/labstack/echo/v4"
	"github.com/noellimx/go-ddd/internal/interface/api/rest"
	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
	// Setup
	mockService := interfaces.NewMockProductService(t)
	reqBody := map[string]interface{}{"Name": "TestProduct", "Price": 9.99, "SellerId": "123e4567-e89b-12d3-a456-426614174000"}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewReader(reqBodyBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctrl := rest.NewProductController(mockService)

	createProductCommandResult := &command.CreateProductCommandResult{
		Result: &common.ProductResult{
			Id:    uuid.New(),
			Name:  "TestProduct",
			Price: 9.99,
		},
	}
	mockService.On("CreateProduct", mock.Anything).Return(createProductCommandResult, nil)

	// Execute
	ctrl.CreateProductController(rec, req)

	// Deserialize the response body
	var responseBody map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Remove fields from responseBody that are not present in reqBody
	// For example, remove Id and Seller fields
	delete(responseBody, "Id")
	delete(responseBody, "Seller")
	delete(reqBody, "SellerId")
	delete(responseBody, "CreatedAt")
	delete(responseBody, "UpdatedAt")

	// Assertions
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, reqBody, responseBody)
	mockService.AssertExpectations(t)
}

func TestGetAllProducts(t *testing.T) {
	// Setup
	mockService := interfaces.NewMockProductService(t) // Assuming you have a mock of ProductService

	expectedProducts := []*entities.Product{
		{
			Id:    uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"),
			Name:  "TestProduct1",
			Price: 9.99,
		}, {
			Id:    uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c9"),
			Name:  "TestProduct2",
			Price: 14.99,
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
	rec := httptest.NewRecorder()

	ctrl := rest.NewProductController(mockService)
	result := &query.GetAllProductsQueryResult{
		Result: []*common.ProductResult{
			{
				Id:    uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"),
				Name:  "TestProduct1",
				Price: 9.99,
			},
			{
				Id:    uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c9"),
				Name:  "TestProduct2",
				Price: 14.99,
			},
		},
	}

	mockService.On("FindAllProducts").Return(result, nil)

	var expectedListResponse response.ListProductsResponse
	for _, product := range expectedProducts {
		expectedListResponse.Products = append(expectedListResponse.Products,
			&response.ProductResponse{
				Id:    product.Id.String(),
				Name:  product.Name,
				Price: product.Price,
			})
	}

	// Assertions
	ctrl.GetAllProductsController(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var receivedListResponse response.ListProductsResponse
	err := json.Unmarshal(rec.Body.Bytes(), &receivedListResponse)
	if assert.NoError(t, err) {
		assert.ElementsMatch(t, expectedListResponse.Products, receivedListResponse.Products)
	}
}
