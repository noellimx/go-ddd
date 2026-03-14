package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/noellimx/go-ddd/internal/infrastructure/db/sqlc"
)

func NewQueries(conn *pgxpool.Pool) *db.Queries {
	return db.New(conn)
}
