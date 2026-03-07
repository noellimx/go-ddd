// Package rest contains HTTP handler tests for the product controller.
//
// These tests are written in the TDD RED phase against the future net/http
// handler signatures BEFORE the Echo → net/http migration is complete.
// They will NOT compile or pass until each handler is changed from:
//
//	func (pc *ProductController) XxxController(c echo.Context) error
//
// to:
//
//	func (pc *ProductController) XxxController(w http.ResponseWriter, r *http.Request)
//
// GitHub Copilot assisted in generating this file.
package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/noellimx/go-ddd/internal/application/command"
	"github.com/noellimx/go-ddd/internal/application/common"
	"github.com/noellimx/go-ddd/internal/application/query"
	"github.com/noellimx/go-ddd/internal/interface/api/rest/dto/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── Mock service ─────────────────────────────────────────────────────────────

// mockProductService is a controllable test double for interfaces.ProductService.
//
// Populate the Result/Err fields for each method before a test case to drive
// the desired behaviour. Fields left at their zero value are safe to call —
// the corresponding method will return (nil, nil).
type mockProductService struct {
	CreateProductResult *command.CreateProductCommandResult
	CreateProductErr    error

	UpdateProductResult *command.UpdateProductCommandResult
	UpdateProductErr    error

	DeleteProductResult *command.DeleteProductCommandResult
	DeleteProductErr    error

	FindAllProductsResult *query.GetAllProductsQueryResult
	FindAllProductsErr    error

	FindProductByIdResult *query.GetProductByIdQueryResult
	FindProductByIdErr    error
}

// CreateProduct implements interfaces.ProductService.
func (m *mockProductService) CreateProduct(cmd *command.CreateProductCommand) (*command.CreateProductCommandResult, error) {
	return m.CreateProductResult, m.CreateProductErr
}

// UpdateProduct implements interfaces.ProductService.
func (m *mockProductService) UpdateProduct(cmd *command.UpdateProductCommand) (*command.UpdateProductCommandResult, error) {
	return m.UpdateProductResult, m.UpdateProductErr
}

// DeleteProduct implements interfaces.ProductService.
func (m *mockProductService) DeleteProduct(cmd *command.DeleteProductCommand) (*command.DeleteProductCommandResult, error) {
	return m.DeleteProductResult, m.DeleteProductErr
}

// FindAllProducts implements interfaces.ProductService.
func (m *mockProductService) FindAllProducts() (*query.GetAllProductsQueryResult, error) {
	return m.FindAllProductsResult, m.FindAllProductsErr
}

// FindProductById implements interfaces.ProductService.
func (m *mockProductService) FindProductById(q *query.GetProductByIdQuery) (*query.GetProductByIdQueryResult, error) {
	return m.FindProductByIdResult, m.FindProductByIdErr
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

// newTestProductController creates a ProductController wired to the given mock
// service without registering any routes, keeping tests independent of any
// specific router implementation.
func newTestProductController(svc *mockProductService) *ProductController {
	return &ProductController{service: svc}
}

// newProductMux registers the three product handlers on a fresh ServeMux using
// the Go 1.22+ path-parameter syntax so that r.PathValue("id") resolves
// correctly inside GetProductByIdController.
func newProductMux(pc *ProductController) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/products", pc.CreateProductController)
	mux.HandleFunc("GET /api/v1/products", pc.GetAllProductsController)
	mux.HandleFunc("GET /api/v1/products/{id}", pc.GetProductByIdController)
	return mux
}

// decodeJSON is a test helper that deserialises the recorder body into dst and
// fails the test immediately if decoding returns an error.
func decodeJSON(t *testing.T, rec *httptest.ResponseRecorder, dst any) {
	t.Helper()
	require.NoError(t, json.NewDecoder(rec.Body).Decode(dst), "response body must be valid JSON")
}

// ─── Fixtures ─────────────────────────────────────────────────────────────────

var (
	fixedTime      = time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	fixedProductID = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	fixedSellerID  = uuid.MustParse("22222222-2222-2222-2222-222222222222")

	// sampleProduct is a fully-populated ProductResult used across happy-path cases.
	sampleProduct = &common.ProductResult{
		Id:        fixedProductID,
		Name:      "Test Product",
		Price:     9.99,
		CreatedAt: fixedTime,
		UpdatedAt: fixedTime,
	}
)

// ─── TestCreateProductController ──────────────────────────────────────────────

// TestCreateProductController verifies all three branches of
// CreateProductController using table-driven sub-tests.
//
// Error paths are listed first; the happy path is last.
func TestCreateProductController(t *testing.T) {
	type testCase struct {
		name           string
		body           string
		setupService   func(m *mockProductService)
		wantStatusCode int
		// wantBody receives the recorder after ServeHTTP so each case can make
		// targeted assertions about the response payload.
		wantBody func(t *testing.T, rec *httptest.ResponseRecorder)
	}

	tests := []testCase{
		{
			name: "malformed JSON body returns 400",
			body: `{invalid-json`,
			setupService: func(m *mockProductService) {
				// the service must never be reached for malformed input
			},
			wantStatusCode: http.StatusBadRequest,
			wantBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var got map[string]string
				decodeJSON(t, rec, &got)
				assert.Equal(t, "Failed to parse request body", got["error"])
			},
		},
		{
			name: "invalid SellerId UUID in body returns 400",
			// SellerId is a valid JSON value but not a UUID, so
			// CreateProductRequest.ToCreateProductCommand() should fail.
			body: `{"idempotency_key":"k1","Name":"Widget","Price":1.5,"SellerId":"not-a-uuid"}`,
			setupService: func(m *mockProductService) {
				// the service must never be reached when command construction fails
			},
			wantStatusCode: http.StatusBadRequest,
			wantBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var got map[string]string
				decodeJSON(t, rec, &got)
				assert.Equal(t, "Invalid product Id format", got["error"])
			},
		},
		{
			name: "service error returns 500",
			body: `{"idempotency_key":"k1","Name":"Widget","Price":1.5,"SellerId":"` + fixedSellerID.String() + `"}`,
			setupService: func(m *mockProductService) {
				m.CreateProductErr = errors.New("db connection lost")
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var got map[string]string
				decodeJSON(t, rec, &got)
				assert.Equal(t, "Failed to create product", got["error"])
			},
		},
		{
			name: "valid request returns 201 with created product",
			body: `{"idempotency_key":"k1","Name":"Widget","Price":1.5,"SellerId":"` + fixedSellerID.String() + `"}`,
			setupService: func(m *mockProductService) {
				m.CreateProductResult = &command.CreateProductCommandResult{
					Result: sampleProduct,
				}
			},
			wantStatusCode: http.StatusCreated,
			wantBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var got response.ProductResponse
				decodeJSON(t, rec, &got)
				assert.Equal(t, fixedProductID.String(), got.Id)
				assert.Equal(t, "Test Product", got.Name)
				assert.InDelta(t, 9.99, got.Price, 0.001)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := &mockProductService{}
			tc.setupService(svc)

			pc := newTestProductController(svc)
			mux := newProductMux(pc)

			req := httptest.NewRequest(
				http.MethodPost,
				"/api/v1/products",
				bytes.NewBufferString(tc.body),
			)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantStatusCode, rec.Code)
			tc.wantBody(t, rec)
		})
	}
}

// ─── TestGetAllProductsController ─────────────────────────────────────────────

// TestGetAllProductsController verifies all branches of
// GetAllProductsController using table-driven sub-tests.
//
// Error paths are listed first; the happy path is last.
func TestGetAllProductsController(t *testing.T) {
	type testCase struct {
		name           string
		setupService   func(m *mockProductService)
		wantStatusCode int
		wantBody       func(t *testing.T, rec *httptest.ResponseRecorder)
	}

	tests := []testCase{
		{
			name: "service error returns 500",
			setupService: func(m *mockProductService) {
				m.FindAllProductsErr = errors.New("db timeout")
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var got map[string]string
				decodeJSON(t, rec, &got)
				assert.Equal(t, "Failed to fetch products", got["error"])
			},
		},
		{
			name: "valid request returns 200 with product list",
			setupService: func(m *mockProductService) {
				m.FindAllProductsResult = &query.GetAllProductsQueryResult{
					Result: []*common.ProductResult{sampleProduct},
				}
			},
			wantStatusCode: http.StatusOK,
			wantBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var got response.ListProductsResponse
				decodeJSON(t, rec, &got)
				require.Len(t, got.Products, 1, "expected exactly one product in the list")
				assert.Equal(t, fixedProductID.String(), got.Products[0].Id)
				assert.Equal(t, "Test Product", got.Products[0].Name)
				assert.InDelta(t, 9.99, got.Products[0].Price, 0.001)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := &mockProductService{}
			tc.setupService(svc)

			pc := newTestProductController(svc)
			mux := newProductMux(pc)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantStatusCode, rec.Code)
			tc.wantBody(t, rec)
		})
	}
}

// ─── TestGetProductByIdController ─────────────────────────────────────────────

// TestGetProductByIdController verifies all branches of
// GetProductByIdController using table-driven sub-tests.
//
// Error paths are listed first; the happy path is last.
// Routing is performed through a ServeMux so that r.PathValue("id") is
// populated correctly inside the handler.
func TestGetProductByIdController(t *testing.T) {
	type testCase struct {
		name           string
		pathID         string // value placed in the {id} path segment
		setupService   func(m *mockProductService)
		wantStatusCode int
		wantBody       func(t *testing.T, rec *httptest.ResponseRecorder)
	}

	tests := []testCase{
		{
			name:   "invalid UUID path param returns 400",
			pathID: "not-a-uuid",
			setupService: func(m *mockProductService) {
				// the service must never be reached when UUID parsing fails
			},
			wantStatusCode: http.StatusBadRequest,
			wantBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var got map[string]string
				decodeJSON(t, rec, &got)
				assert.Equal(t, "Invalid product Id format", got["error"])
			},
		},
		{
			name:   "product not found (nil result, no error) returns 404",
			pathID: fixedProductID.String(),
			setupService: func(m *mockProductService) {
				// service returns (nil, nil) — query succeeded but no record exists
				m.FindProductByIdResult = nil
				m.FindProductByIdErr = nil
			},
			wantStatusCode: http.StatusNotFound,
			wantBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var got map[string]string
				decodeJSON(t, rec, &got)
				assert.Equal(t, "Product not found", got["error"])
			},
		},
		{
			name:   "service error returns 500",
			pathID: fixedProductID.String(),
			setupService: func(m *mockProductService) {
				m.FindProductByIdErr = errors.New("db unreachable")
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var got map[string]string
				decodeJSON(t, rec, &got)
				assert.Equal(t, "Failed to fetch product", got["error"])
			},
		},
		{
			name:   "valid UUID returns 200 with product",
			pathID: fixedProductID.String(),
			setupService: func(m *mockProductService) {
				m.FindProductByIdResult = &query.GetProductByIdQueryResult{
					Result: sampleProduct,
				}
			},
			wantStatusCode: http.StatusOK,
			wantBody: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var got response.ProductResponse
				decodeJSON(t, rec, &got)
				assert.Equal(t, fixedProductID.String(), got.Id)
				assert.Equal(t, "Test Product", got.Name)
				assert.InDelta(t, 9.99, got.Price, 0.001)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := &mockProductService{}
			tc.setupService(svc)

			pc := newTestProductController(svc)
			mux := newProductMux(pc)

			// Route through the mux so the Go 1.22 path-parameter matching
			// populates r.PathValue("id") before the handler is called.
			req := httptest.NewRequest(
				http.MethodGet,
				"/api/v1/products/"+tc.pathID,
				nil,
			)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantStatusCode, rec.Code)
			tc.wantBody(t, rec)
		})
	}
}
