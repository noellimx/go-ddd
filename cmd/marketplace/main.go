package main

import (
	"context"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/sklinkert/go-ddd/internal/application/services"
	"github.com/sklinkert/go-ddd/internal/infrastructure/db/postgres"
	"github.com/sklinkert/go-ddd/internal/interface/api/rest"
)

func main() {
	dsn := "host=localhost user=marketplace password=marketplace dbname=marketplace port=5432 sslmode=disable"
	port := ":8080"

	ctx := context.Background()
	conn, err := postgres.NewConnection(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	queries := postgres.NewQueries(conn)

	productRepo := postgres.NewSqlcProductRepository(queries)
	sellerRepo := postgres.NewSqlcSellerRepository(queries)
	idempotencyRepo := postgres.NewSqlcIdempotencyRepository(queries)

	productService := services.NewProductService(productRepo, sellerRepo, idempotencyRepo)
	sellerService := services.NewSellerService(sellerRepo, idempotencyRepo)

	e := echo.New()
	rest.NewProductController(e, productService)
	rest.NewSellerController(e, sellerService)

	if err := e.Start(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
