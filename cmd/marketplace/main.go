package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/noellimx/go-ddd/cmd/config"
	"github.com/noellimx/go-ddd/internal/application/services"
	"github.com/noellimx/go-ddd/internal/infrastructure/db/postgres"
	"github.com/noellimx/go-ddd/internal/interface/api/rest"
)

func main() {
	ctx := context.Background()

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.LevelDebug,
		ReplaceAttr: nil,
	})

	defaultLogger := slog.New(handler)
	slog.SetDefault(defaultLogger)

	defaultConfig, dbConn, err := Init(ctx)
	if err != nil {
		panic(err)
	}
	defaultLogger.Info("Starting http-service::main().")

	queries := postgres.NewQueries(dbConn)

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

	if err := e.Start(defaultConfig.Server.Address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

const (
	maxConnIdleTimeMinute = 5
	minConn               = 2
	maxConn               = 30
)

func Init(ctx context.Context) (*config.App, *pgxpool.Pool, error) {
	appConfig, err := config.ReadDefaultConfig(ctx)
	if err != nil {
		return nil, nil, err
	}

	// DBInfo Connection Pool

	dbConfig, err := pgxpool.ParseConfig(appConfig.ConnString)
	if err != nil {
		return nil, nil, err
	}

	dbConfig.MaxConns = maxConn
	dbConfig.MinConns = minConn
	dbConfig.MaxConnIdleTime = maxConnIdleTimeMinute * time.Minute

	// Create the pool
	dbConnPool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, nil, err
	}

	if dbConnPool == nil {
		return nil, nil, errors.New("appConfig or dbConnPool is nil")
	}
	return &appConfig, dbConnPool, nil
}
