package data

import (
	"github.com/jackc/pgx/v4"
	"context"
)

var (
    DB *pgx.Conn
	Err error
	//REWRITE
    SecretKey = []byte("key")
)

func Reconnect(ctx context.Context, conn *pgx.Conn) (*pgx.Conn, error) {
    conn.Close(ctx)
    newConn, err := pgx.Connect(ctx, "postgres://postgres:postgres@localhost:5432/value_exchanger")
    if err != nil {
        return nil, err
    }
    return newConn, nil
}