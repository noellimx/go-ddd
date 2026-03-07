# PRD Title
Refactor for removing `echo` dependency and use `net/http` etc  in controller layer instead.

## 1. Purpose
- Eliminate the `github.com/labstack/echo/v4` dependency from the controller layer.
- Use only Go standard library (`net/http`, `encoding/json`) for HTTP handling.
- Reduce external dependencies and simplify the interface layer.

## 2. Scope
### In scope
- `internal/interface/api/rest/product_controller.go`
- `internal/interface/api/rest/seller_controller.go`
- `cmd/marketplace/main.go` — route registration currently uses `*echo.Echo`

### Out of scope
- Application/domain/infrastructure layers (no changes to services, repos, DTOs, mappers)
- Adding new endpoints or changing business logic

## 5. Technical Requirements

### 5.1 Code Changes

| File | Echo API to replace | `net/http` equivalent |
|------|--------------------|-----------------------|
| Both controllers | `c.Bind(&req)` | `json.NewDecoder(r.Body).Decode(&req)` |
| Both controllers | `c.JSON(status, body)` | `w.Header().Set("Content-Type","application/json")`, `w.WriteHeader(status)`, `json.NewEncoder(w).Encode(body)` |
| Both controllers | `c.Param("id")` | `r.PathValue("id")` (Go 1.22+ ServeMux) |
| SellerController | `c.Request().URL.Path[len("..."):]` (hack) | replace with `r.PathValue("id")` — fixes existing bug |
| SellerController | `c.NoContent(http.StatusNoContent)` | `w.WriteHeader(http.StatusNoContent)` |
| Both constructors | `func New*Controller(e *echo.Echo, ...)` | `func New*Controller(mux *http.ServeMux, ...)` |
| ProductController | `middleware.Recover()` | thin `recoveryMiddleware(next http.Handler) http.Handler` wrapper |
| `cmd/marketplace/main.go` | `echo.New()`, `e.Start(port)` | `http.NewServeMux()`, `http.ListenAndServe(port, mux)` |

Handler signatures change from:
```go
func (pc *ProductController) CreateProductController(c echo.Context) error
```
to:
```go
func (pc *ProductController) CreateProductController(w http.ResponseWriter, r *http.Request)
```

Route patterns change from Echo-style `/api/v1/sellers/:id` to Go 1.22 ServeMux style `GET /api/v1/sellers/{id}`.

### 5.2 go.mod
- Remove `github.com/labstack/echo/v4` and its transitive dependencies after migration (`go mod tidy`).

### 5.3 Tests and Implementation Approach

**Iteration order (one commit per controller):**
1. `ProductController` (3 handlers: CreateProduct, GetAllProducts, GetProductById)
2. `SellerController` (5 handlers: CreateSeller, GetAllSellers, GetSellerById, PutSeller, DeleteSeller)
3. `cmd/marketplace/main.go` wiring + `go mod tidy`

**Per-handler test strategy using `net/http/httptest`:**
- `httptest.NewRecorder()` as `ResponseWriter`, `httptest.NewRequest()` as `*http.Request`
- Assert response status code and JSON body
- Table-driven tests per handler covering:
  - Happy path (200/201/204)
  - Malformed JSON body → 400
  - Invalid UUID path param → 400
  - Service returns error → 500
  - Not-found result → 404 (GetById handlers)

**Coverage target:** ≥80% on both controller files.

## 6. Risks & Mitigations

| Risk | Mitigation |
|------|-----------|
| Breaking existing route behaviour | Write `httptest`-based tests before migrating each handler (TDD red→green) |
| `c.Param("id")` hack in SellerController (lines 76, 134) | Replace with `r.PathValue("id")` — Go 1.22+ ServeMux supports named path segments |
| Recover middleware loss | Implement a thin `recoveryMiddleware(next http.Handler) http.Handler` wrapper |
| Missed echo transitive dep | Run `go mod tidy` and verify with `go mod graph` after removal |

## 7. Success Criteria
- Zero imports of `github.com/labstack/echo` in `internal/interface/api/rest/`
- `go mod tidy` removes echo and its transitive deps from `go.sum`
- All existing routes respond identically (same status codes, same JSON shapes)
- ≥80% test coverage on both controller files
- `GetSellerById` and `DeleteSeller` no longer use the URL-parsing hack

## Example of current code and proposed change

**Before (echo):**
```go
func NewProductController(e *echo.Echo, service interfaces.ProductService) *ProductController {
    e.POST("/api/v1/products", controller.CreateProductController)
    e.GET("/api/v1/products/:id", controller.GetProductByIdController)
    e.Use(middleware.Recover())
    ...
}

func (pc *ProductController) GetProductByIdController(c echo.Context) error {
    id, err := uuid.Parse(c.Param("id"))
    ...
    return c.JSON(http.StatusOK, response)
}
```

**After (net/http):**
```go
func NewProductController(mux *http.ServeMux, service interfaces.ProductService) *ProductController {
    mux.HandleFunc("POST /api/v1/products", controller.CreateProductController)
    mux.HandleFunc("GET /api/v1/products/{id}", controller.GetProductByIdController)
    ...
}

func (pc *ProductController) GetProductByIdController(w http.ResponseWriter, r *http.Request) {
    id, err := uuid.Parse(r.PathValue("id"))
    ...
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}
```
