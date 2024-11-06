package data

import (
	"github.com/jackc/pgx/v4"
)

var (
    DB *pgx.Conn
	Err error
	//REWRITE
    SecretKey = []byte("key")
)