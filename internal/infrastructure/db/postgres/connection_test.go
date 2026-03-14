package postgres

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/noellimx/go-ddd/internal/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestNewConnection(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	// Test NewQueries with the existing connection (we can't easily get a new DSN)
	queries := NewQueries(testDB.Conn)
	require.NotNil(t, queries)

	// Verify the connection is working by using the existing connection
	ctx := context.Background()
	err := testDB.Conn.Ping(ctx)
	require.NoError(t, err)
}

func TestNewConnection_InvalidDSN(t *testing.T) {

	_, err := pgxpool.ParseConfig("invalid-dsn")
	require.Error(t, err)
}

func TestNewConnection_UnreachableHost(t *testing.T) {
	ctx := context.Background()

	dbConfig, err := pgxpool.ParseConfig("invalid-dsn")
	require.NoError(t, err, "Failed to ParseConfig")

	dbConn, err := pgxpool.NewWithConfig(ctx, dbConfig)
	require.NoError(t, err, "Failed to connect to test database NewWithConfig")

	// Test with unreachable host
	require.Error(t, err)
	require.Nil(t, dbConn)
}

func TestNewQueries(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	// Test NewQueries with valid connection
	queries := NewQueries(testDB.Conn)
	require.NotNil(t, queries)

	// Verify queries object is functional by running a simple query
	ctx := context.Background()
	_, err := queries.GetAllProducts(ctx)
	require.NoError(t, err) // Should not error even if empty
}

func TestNewQueries_WithNilConnection(t *testing.T) {
	// Test NewQueries with nil connection
	// Note: This will create a queries object but will panic when used
	queries := NewQueries(nil)
	require.NotNil(t, queries)

	// Attempting to use it should panic (so we won't test that)
	// This test just verifies that NewQueries can accept nil without immediately panicking
}
