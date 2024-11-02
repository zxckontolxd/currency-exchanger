package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"strings"
	"net/http"
	"fmt"

	. "gw-currency-wallet/cmd/data"
	log "github.com/sirupsen/logrus"
)

func Deposit(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		log.Errorf("Token not specified or incorrect format")
		log.Infof(authHeader)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token not specified or incorrect format"})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	var dep Deposite

	if Err := ctx.ShouldBindJSON(&dep); Err != nil {
		log.Errorf("Field to bind JSON: %v", Err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": Err.Error()})
		return
	}

	switch dep.Currency {
	case "RUB":
		break
	case "EUR":
		break
	case "USD":
		break
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount or currency"})
		return
	}

	// проверки на каст к float нет, так как все равно он не сможет заанмаршалить что-то, помимо float

	tx, err := DB.Begin(ctx)
	if err != nil {
		log.Errorf("Cannot start transaction: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	valueStr := "balance_" + strings.ToLower(dep.Currency)

	var builder strings.Builder
	builder.WriteString("UPDATE wallets ")
	builder.WriteString(fmt.Sprintf("SET %s = %s + $1 ", valueStr, valueStr))
	builder.WriteString("WHERE id = ( ")
	builder.WriteString("SELECT wallets.id ")
	builder.WriteString("FROM JWTTokens ")
	builder.WriteString("JOIN users ON JWTTokens.user_id = users.id ")
	builder.WriteString("JOIN wallets ON users.wallet_id = wallets.id ")
	builder.WriteString("WHERE JWTTokens.token = $2);")

	_, err = tx.Exec(ctx, builder.String(), dep.Amount, token)
	if err != nil {
		log.Errorf("Error to update wallet: %v", err)
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


	// по хорошему, этот запрос можно выделить в отдельную функцию
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

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Account topped up successfully",
		"new_balance": gin.H{
			"USD": balance_usd,
			"RUB": balance_rub,
			"EUR": balance_eur,
		},
	})

	return
}