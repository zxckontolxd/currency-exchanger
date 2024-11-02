package handlers

import (
	"strings"
	"net/http"

//	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"

	. "gw-currency-wallet/cmd/data"
	log "github.com/sirupsen/logrus"
)

func Balance(ctx *gin.Context) {
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

	ctx.JSON(http.StatusOK, gin.H{
		"balance": gin.H{
			"USD": balance_usd,
			"RUB": balance_rub,
			"EUR": balance_eur,
		},
	})

	return
}