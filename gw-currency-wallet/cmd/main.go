package main

import (
    "github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"context"
	"github.com/sirupsen/logrus"

	. "gw-currency-wallet/internal/data"
	. "gw-currency-wallet/internal/handlers"
	. "gw-currency-wallet/internal/pkg/logger"
)

// немного мешанины, так как изначально планировалось делать все в одном файле, но он слишком разросся

func main() {
	router := gin.Default()

	DB, Err = pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/value_exchanger")
	defer DB.Close(context.Background())

	loggerConfig := LoggerConfig{
		Level: logrus.InfoLevel,
		LogToFile: false,
		LogFilePath: "",
	}
	logger := NewLogger(loggerConfig)
	logger.Info("TEST")

	router.POST("/api/v1/register", Register)
	router.POST("/api/v1/login", Login)
	router.GET("/api/v1/balance", Balance)
	router.POST("/api/v1/deposit", Deposit)
	router.POST("/api/v1/withdraw", Withdraw)
	router.GET("/api/v1/rates", Rates)
	router.POST("/api/v1/exchange", Exchange)

	router.Run(":8080")
}