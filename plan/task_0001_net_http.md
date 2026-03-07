# Task List: epic_0001_net_http — Remove Echo, Use net/http

## Status Legend
- [ ] pending
- [~] in progress
- [x] done

---

## Phase 1 — ProductController

### T1.1 Implement `writeJSON` helper
- [ ] T1.1.1 Create `internal/interface/api/rest/helpers.go`
- [ ] T1.1.2 Implement `writeJSON(w http.ResponseWriter, status int, body any)` — sets `Content-Type: application/json`, calls `w.WriteHeader(status)`, encodes body with `json.NewEncoder(w).Encode(body)`

### T1.2 Migrate `NewProductController` constructor
- [ ] T1.2.1 Change parameter: `e *echo.Echo` → `mux *http.ServeMux`
- [ ] T1.2.2 Replace `e.POST("/api/v1/products", ...)` → `mux.HandleFunc("POST /api/v1/products", ...)`
- [ ] T1.2.3 Replace `e.GET("/api/v1/products", ...)` → `mux.HandleFunc("GET /api/v1/products", ...)`
- [ ] T1.2.4 Replace `e.GET("/api/v1/products/:id", ...)` → `mux.HandleFunc("GET /api/v1/products/{id}", ...)`
- [ ] T1.2.5 Remove `e.Use(middleware.Recover())` (recovery applied at main.go level)
- [ ] T1.2.6 Remove imports: `"github.com/labstack/echo/v4"`, `"github.com/labstack/echo/v4/middleware"`

### T1.3 Migrate `CreateProductController`
- [ ] T1.3.1 Change signature: `(c echo.Context) error` → `(w http.ResponseWriter, r *http.Request)`
- [ ] T1.3.2 Replace `c.Bind(&req)` → `json.NewDecoder(r.Body).Decode(&req)`
- [ ] T1.3.3 Replace `return c.JSON(http.StatusBadRequest, ...)` (bind error) → `writeJSON(w, http.StatusBadRequest, ...)`
- [ ] T1.3.4 Replace `return c.JSON(http.StatusBadRequest, ...)` (UUID error) → `writeJSON(w, http.StatusBadRequest, ...)`
- [ ] T1.3.5 Replace `return c.JSON(http.StatusInternalServerError, ...)` → `writeJSON(w, http.StatusInternalServerError, ...)`
- [ ] T1.3.6 Replace `return c.JSON(http.StatusCreated, response)` → `writeJSON(w, http.StatusCreated, response)`

### T1.4 Migrate `GetAllProductsController`
- [ ] T1.4.1 Change signature: `(c echo.Context) error` → `(w http.ResponseWriter, r *http.Request)`
- [ ] T1.4.2 Replace `return c.JSON(http.StatusInternalServerError, ...)` → `writeJSON(w, http.StatusInternalServerError, ...)`
- [ ] T1.4.3 Replace `return c.JSON(http.StatusOK, response)` → `writeJSON(w, http.StatusOK, response)`

### T1.5 Migrate `GetProductByIdController`
- [ ] T1.5.1 Change signature: `(c echo.Context) error` → `(w http.ResponseWriter, r *http.Request)`
- [ ] T1.5.2 Replace `c.Param("id")` → `r.PathValue("id")`
- [ ] T1.5.3 Replace `return c.JSON(http.StatusBadRequest, ...)` (UUID error) → `writeJSON(w, http.StatusBadRequest, ...)`
- [ ] T1.5.4 Replace `return c.JSON(http.StatusInternalServerError, ...)` → `writeJSON(w, http.StatusInternalServerError, ...)`
- [ ] T1.5.5 Replace `return c.JSON(http.StatusNotFound, ...)` → `writeJSON(w, http.StatusNotFound, ...)`
- [ ] T1.5.6 Replace `return c.JSON(http.StatusOK, response)` → `writeJSON(w, http.StatusOK, response)`

### T1.6 Write tests for ProductController
- [ ] T1.6.1 Create `internal/interface/api/rest/product_controller_test.go`
- [ ] T1.6.2 Define `mockProductService` implementing `interfaces.ProductService` with controllable return values
- [ ] T1.6.3 Test `CreateProductController`: happy path → 201
- [ ] T1.6.4 Test `CreateProductController`: malformed JSON → 400
- [ ] T1.6.5 Test `CreateProductController`: service error → 500
- [ ] T1.6.6 Test `GetAllProductsController`: happy path → 200
- [ ] T1.6.7 Test `GetAllProductsController`: service error → 500
- [ ] T1.6.8 Test `GetProductByIdController`: happy path → 200
- [ ] T1.6.9 Test `GetProductByIdController`: invalid UUID → 400
- [ ] T1.6.10 Test `GetProductByIdController`: not found → 404
- [ ] T1.6.11 Test `GetProductByIdController`: service error → 500

---

## Phase 2 — SellerController

### T2.1 Migrate `NewSellerController` constructor
- [ ] T2.1.1 Change parameter: `e *echo.Echo` → `mux *http.ServeMux`
- [ ] T2.1.2 Replace `e.POST("/api/v1/sellers", ...)` → `mux.HandleFunc("POST /api/v1/sellers", ...)`
- [ ] T2.1.3 Replace `e.GET("/api/v1/sellers", ...)` → `mux.HandleFunc("GET /api/v1/sellers", ...)`
- [ ] T2.1.4 Replace `e.GET("/api/v1/sellers/:id", ...)` → `mux.HandleFunc("GET /api/v1/sellers/{id}", ...)`
- [ ] T2.1.5 Replace `e.PUT("/api/v1/sellers", ...)` → `mux.HandleFunc("PUT /api/v1/sellers", ...)`
- [ ] T2.1.6 Replace `e.DELETE("/api/v1/sellers/:id", ...)` → `mux.HandleFunc("DELETE /api/v1/sellers/{id}", ...)`
- [ ] T2.1.7 Remove import `"github.com/labstack/echo/v4"`

### T2.2 Migrate `CreateSellerController`
- [ ] T2.2.1 Change signature: `(c echo.Context) error` → `(w http.ResponseWriter, r *http.Request)`
- [ ] T2.2.2 Replace `c.Bind(&req)` → `json.NewDecoder(r.Body).Decode(&req)`
- [ ] T2.2.3 Replace all 3 `return c.JSON(...)` calls → `writeJSON(w, ...)`

### T2.3 Migrate `GetAllSellersController`
- [ ] T2.3.1 Change signature: `(c echo.Context) error` → `(w http.ResponseWriter, r *http.Request)`
- [ ] T2.3.2 Replace all 2 `return c.JSON(...)` calls → `writeJSON(w, ...)`

### T2.4 Migrate `GetSellerByIdController`
- [ ] T2.4.1 Change signature: `(c echo.Context) error` → `(w http.ResponseWriter, r *http.Request)`
- [ ] T2.4.2 Remove hack: delete `idRaw := c.Request().URL.Path[len("/api/v1/sellers/"):]`
- [ ] T2.4.3 Replace with `idRaw := r.PathValue("id")`
- [ ] T2.4.4 Replace all 3 `return c.JSON(...)` calls → `writeJSON(w, ...)`

### T2.5 Migrate `PutSellerController`
- [ ] T2.5.1 Change signature: `(c echo.Context) error` → `(w http.ResponseWriter, r *http.Request)`
- [ ] T2.5.2 Replace `c.Bind(&req)` → `json.NewDecoder(r.Body).Decode(&req)`
- [ ] T2.5.3 Replace all 3 `return c.JSON(...)` calls → `writeJSON(w, ...)`

### T2.6 Migrate `DeleteSellerController`
- [ ] T2.6.1 Change signature: `(c echo.Context) error` → `(w http.ResponseWriter, r *http.Request)`
- [ ] T2.6.2 Remove hack: delete `idRaw := c.Request().URL.Path[len("/api/v1/sellers/"):]`
- [ ] T2.6.3 Replace with `idRaw := r.PathValue("id")`
- [ ] T2.6.4 Replace `return c.JSON(http.StatusInternalServerError, ...)` → `writeJSON(w, http.StatusInternalServerError, ...)`
- [ ] T2.6.5 Replace `return c.NoContent(http.StatusNoContent)` → `w.WriteHeader(http.StatusNoContent)`

### T2.7 Write tests for SellerController
- [ ] T2.7.1 Create `internal/interface/api/rest/seller_controller_test.go`
- [ ] T2.7.2 Define `mockSellerService` implementing `interfaces.SellerService` with controllable return values
- [ ] T2.7.3 Test `CreateSellerController`: happy path → 201
- [ ] T2.7.4 Test `CreateSellerController`: malformed JSON → 400
- [ ] T2.7.5 Test `CreateSellerController`: service error → 500
- [ ] T2.7.6 Test `GetAllSellersController`: happy path → 200
- [ ] T2.7.7 Test `GetAllSellersController`: service error → 500
- [ ] T2.7.8 Test `GetSellerByIdController`: happy path → 200
- [ ] T2.7.9 Test `GetSellerByIdController`: invalid UUID → 400
- [ ] T2.7.10 Test `GetSellerByIdController`: not found → 404
- [ ] T2.7.11 Test `GetSellerByIdController`: service error → 500
- [ ] T2.7.12 Test `PutSellerController`: happy path → 200
- [ ] T2.7.13 Test `PutSellerController`: malformed JSON → 400
- [ ] T2.7.14 Test `PutSellerController`: service error → 500
- [ ] T2.7.15 Test `DeleteSellerController`: happy path → 204
- [ ] T2.7.16 Test `DeleteSellerController`: invalid UUID → 400
- [ ] T2.7.17 Test `DeleteSellerController`: service error → 500

---

## Phase 3 — Recovery Middleware

### T3.1 Implement `RecoveryMiddleware`
- [ ] T3.1.1 Create `internal/interface/api/rest/middleware.go`
- [ ] T3.1.2 Implement `RecoveryMiddleware(next http.Handler) http.Handler`
- [ ] T3.1.3 Use `defer func() { if r := recover(); r != nil { ... } }()` inside the handler
- [ ] T3.1.4 On panic: call `writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})`

---

## Phase 4 — Wiring & Cleanup

### T4.1 Update `cmd/marketplace/main.go`
- [ ] T4.1.1 Remove `"github.com/labstack/echo/v4"` import
- [ ] T4.1.2 Replace `e := echo.New()` → `mux := http.NewServeMux()`
- [ ] T4.1.3 Update `NewProductController(e, productService)` → `NewProductController(mux, productService)`
- [ ] T4.1.4 Update `NewSellerController(e, sellerService)` → `NewSellerController(mux, sellerService)`
- [ ] T4.1.5 Wrap mux: `handler := rest.RecoveryMiddleware(mux)`
- [ ] T4.1.6 Replace `e.Start(port)` → `log.Fatal(http.ListenAndServe(port, handler))`

### T4.2 Remove echo dependency
- [ ] T4.2.1 Run `go mod tidy`
- [ ] T4.2.2 Verify: `go mod graph | grep echo` returns nothing
- [ ] T4.2.3 Verify: `grep -r "labstack/echo" internal/` returns nothing

### T4.3 Verify build and tests
- [ ] T4.3.1 Run `go build ./...` — no errors
- [ ] T4.3.2 Run `go test ./internal/interface/api/rest/...` — all tests pass
- [ ] T4.3.3 Run `go test -cover ./internal/interface/api/rest/...` — coverage ≥80%

---

## Open Questions / Review Items
- [ ] Q1: Should `RecoveryMiddleware` log the panic before returning 500? If yes, which logger?
- [ ] Q2: `PutSellerController` ID comes from request body — should route stay `PUT /api/v1/sellers` or move to `PUT /api/v1/sellers/{id}`?
- [ ] Q3: Are there any integration or e2e tests outside `internal/` that wire `*echo.Echo` and also need updating?
