package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/noellimx/go-ddd/internal/application/services"
	"github.com/noellimx/go-ddd/internal/infrastructure/db/postgres"
	"github.com/noellimx/go-ddd/internal/interface/api/rest"
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

	mux := http.NewServeMux()

	const readTimeoutSeconds = 10
	(&http.Server{
		Addr:                         ":8080",
		Handler:                      mux,
		DisableGeneralOptionsHandler: false,
		TLSConfig:                    nil,
		ReadTimeout:                  readTimeoutSeconds * time.Second,
		ReadHeaderTimeout:            readTimeoutSeconds * time.Second,
		WriteTimeout:                 0,
		IdleTimeout:                  0,
		MaxHeaderBytes:               0,
		TLSNextProto:                 nil,
		ConnState:                    nil,
		ErrorLog:                     log.Default(),
		BaseContext:                  nil,
		ConnContext:                  nil,
		HTTP2:                        nil,
		Protocols:                    nil,
	}).ListenAndServe()

	rest.NewProductController(productService)

	e := echo.New()
	rest.NewSellerController(e, sellerService)

	if err := e.Start(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
