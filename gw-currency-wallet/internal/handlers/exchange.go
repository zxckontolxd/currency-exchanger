package handlers

import (
	pb "github.com/zxckontolxd/proto-exchange/exchange"
    log "github.com/sirupsen/logrus"
	. "gw-currency-wallet/internal/data"

    "google.golang.org/grpc"
    "net/http"
	"fmt"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

// вот тут таится босс говнокода
// добавлю, что в реальном проекте следует использовать библиотеку для больших чисел
// она пофиксит float и можно будет безопасно оперировать большими суммами

func Exchange(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		log.Errorf("Token not specified or incorrect format")
		log.Infof(authHeader)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token not specified or incorrect format"})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	query := `
	SELECT w.balance_usd, w.balance_eur, w.balance_rub
	FROM JWTTokens AS tk
	JOIN Users AS u ON tk.user_id = u.id
	JOIN Wallets AS w ON u.wallet_id = w.id
	WHERE tk.token = $1;`

	var balance_usd float64
	var balance_eur float64
	var balance_rub float64

	Err = DB.QueryRow(ctx, query, token).Scan(&balance_usd, &balance_eur, &balance_rub)
	if Err == pgx.ErrNoRows {
		log.Errorf("Token not found: %v", Err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "token not found"})
		return
	} else if Err != nil {
		log.Errorf("Database error: %v", Err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	currencyMap := make(map[string]float64)
	currencyMap["USD"] = balance_usd
	currencyMap["EUR"] = balance_eur
	currencyMap["RUB"] = balance_rub

	var exchanger Exchanger

	if Err := ctx.ShouldBindJSON(&exchanger); Err != nil {
		log.Errorf("Field to bind JSON: %v", Err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": Err.Error()})
		return
	}

	if currencyMap[exchanger.FromCurrency] < exchanger.Amount {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient funds or invalid currencies"})
		return
	}

	exchangerConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        log.Errorf("Cannot dial with grpc service (gw-exchanger): %v", err)
        //TODO json
        return
    }
    defer exchangerConn.Close()

    exchangerClient := pb.NewExchangeServiceClient(exchangerConn)

	request := &pb.CurrencyRequest{
        FromCurrency: exchanger.FromCurrency,
        ToCurrency:   exchanger.ToCurrency,
    }

    response, err := exchangerClient.GetExchangeRateForCurrency(ctx, request)
    if err != nil {
        log.Errorf("Cannot get rate: %v", err)
        //TODO json
        return
    }
    
	// запрос списать amount
	// запрос добавить обмененую валюту

	tx, err := DB.Begin(ctx)
	if err != nil {
		log.Errorf("Cannot start transaction: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	defer func() {
		if err != nil {
		log.Errorf("Transaction field: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		tx.Rollback(ctx)
		return
	} else {
		err = tx.Commit(ctx)
		if err != nil {
			log.Errorf("Field to commit transaction: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			tx.Rollback(ctx)
			return
			}
		}
	}()

	from := "balance_" + strings.ToLower(exchanger.FromCurrency)
	to := "balance_" + strings.ToLower(exchanger.ToCurrency)

	converted := exchanger.Amount * float64(response.Rate)

	var builder strings.Builder
	builder.WriteString("UPDATE wallets SET ")
	builder.WriteString(fmt.Sprintf("%s = %s - $1, ", from, from))
	builder.WriteString(fmt.Sprintf("%s = %s + $2 ", to, to))
	builder.WriteString("WHERE id = ( ")
	builder.WriteString("SELECT wallets.id ")
	builder.WriteString("FROM JWTTokens ")
	builder.WriteString("JOIN users ON JWTTokens.user_id = users.id ")
	builder.WriteString("JOIN wallets ON users.wallet_id = wallets.id ")
	builder.WriteString("WHERE JWTTokens.token = $3);")

	_, err = tx.Exec(ctx, builder.String(), exchanger.Amount, converted, token)
	if err != nil {
		log.Errorf("Error to update wallet: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// это можно исправить на мапу с валютами

	var balanceFrom float64
	var balanceTo float64

	query = fmt.Sprintf(`
	SELECT w.%s, w.%s
	FROM JWTTokens AS tk
	JOIN Users AS u ON tk.user_id = u.id
	JOIN Wallets AS w ON u.wallet_id = w.id
	WHERE tk.token = $1;
	`, from, to)

	Err = DB.QueryRow(ctx, query, token).Scan(&balanceFrom, &balanceTo)
	if Err == pgx.ErrNoRows {
		log.Errorf("Token not found: %v", Err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "token not found"})
		return
	} else if Err != nil {
		log.Errorf("Database error: %v", Err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Exchange successful",
		"exchanged_amount": converted,
		"new_balance": gin.H{
			exchanger.FromCurrency: balanceFrom,
			exchanger.ToCurrency: balanceTo,
		},
	})
}