package database

import (
	"context"
	"fmt"
	"log"
	"main/internal/config"

	"github.com/jackc/pgx/v5"
)

const dbConfigString = "postgres://%s:%s@%s:%s/%s"

type db struct {
	// logger *slog.Logger
	conn   *pgx.Conn
}

type DB interface {
	NewConn(ctx context.Context, config config.Config) DB
	GetConn() *pgx.Conn
}

func NewDB() DB {
	return &db{}
}

func (d *db) NewConn(ctx context.Context, config config.Config) DB {
	connString := fmt.Sprintf(dbConfigString, config.Database.User, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.Name)

	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		log.Fatalf("Failed to ping the database: %v\n", err)
	}

	// logger.Info("New Postgres connection opened")

	return &db{
		conn:   conn,
	}
}

func (db *db) GetConn() *pgx.Conn {
	return db.conn
}
